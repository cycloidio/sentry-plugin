package service

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/cycloidio/sentry-plugin/config"
	"github.com/cycloidio/sentry-plugin/event"
	"github.com/cycloidio/sentry-plugin/issue"
	"github.com/cycloidio/sentry-plugin/organization"
	"github.com/cycloidio/sentry-plugin/project"
	"github.com/cycloidio/sentry-plugin/sentry"

	sentryAPI "github.com/atlassian/go-sentry-api"
)

type Service interface {
	Ping(ctx context.Context) Status
	Event(ctx context.Context, e event.Event)
	DeletePlugin(ctx context.Context)
	Resync(ctx context.Context)
	ListAllIssues(ctx context.Context) ([]issue.IssueWithRelations, error)
}

type Plugin struct {
	organizations organization.Repository
	projects      project.Repository
	issues        issue.Repository

	sentry sentry.Service

	mxStatus sync.RWMutex
	status   Status

	started bool

	config *config.Config

	logger *slog.Logger
}

func New(ctx context.Context, or organization.Repository, pr project.Repository, ir issue.Repository, ss sentry.Service, started bool, cfg *config.Config, logger *slog.Logger) *Plugin {
	p := &Plugin{
		organizations: or,
		projects:      pr,
		issues:        ir,

		sentry: ss,

		started: started,

		config: cfg,

		logger: logger,
	}

	// Once the Plugin get's initialized we have to pull everything.
	// The 'Resync' is what it does so we call it on the BG to pull
	// all the info
	go p.Resync(ctx)

	return p
}

func (p *Plugin) Ping(ctx context.Context) Status {
	p.mxStatus.RLock()
	defer p.mxStatus.RUnlock()

	return p.status
}

func (p *Plugin) Event(ctx context.Context, e event.Event) {
	// NOTE: If it's of type project:create we could create directly
	// a new Project
}

func (p *Plugin) DeletePlugin(ctx context.Context) {
	// NOTE: Nothing to do here
}

func (p *Plugin) Resync(ctx context.Context) {
	p.logger.Info("resync started")
	p.setStatus(Syncthing)

	select {
	case <-ctx.Done():
		p.logger.Info("resync cancelled: context done")
		return
	default:
	}
	if !p.started {
		p.logger.Info("resync skipped: plugin not started")
		return
	}
	// This will delete everything as all have FK to organizations
	p.logger.Info("deleting all organizations")
	err := p.organizations.DeleteAll(ctx)
	if err != nil {
		ferr := fmt.Errorf("failed to delete all Organizations: %w", err)
		p.logger.Error(ferr.Error())
		p.setStatus(Error)
		return
	}
	sorgs := make([]sentryAPI.Organization, 0)

	if p.config.Sentry.OrganizationSlug != "" {
		slug := p.config.Sentry.OrganizationSlug
		p.logger.Info("using configured organization", "slug", slug)
		sorgs = append(sorgs, sentryAPI.Organization{
			ID:   &slug,
			Slug: &slug,
			Name: slug,
		})
	} else {
		p.logger.Info("fetching all organizations")
		sorgs, _, err = p.sentry.GetOrganizations()
		if err != nil {
			ferr := fmt.Errorf("failed to get Sentry Organizations: %w", err)
			p.logger.Error(ferr.Error())
			p.setStatus(Error)
			return
		}
	}
	p.logger.Info("organizations fetched", "count", len(sorgs))

	hasErrors := false
	for _, o := range sorgs {
		p.logger.Info("syncing organization", "slug", *o.Slug)
		_, err := p.organizations.Create(ctx, sentry.ToOrganization(o))
		if err != nil {
			ferr := fmt.Errorf("failed to create Organization: %w", err)
			p.logger.Error(ferr.Error())
			hasErrors = true
			continue
		}

		p.logger.Info("fetching projects", "organization", *o.Slug)
		var sprojs []sentryAPI.Project
		if p.config.Sentry.OrganizationSlug != "" {
			pageProjs, link, perr := p.sentry.GetProjects()
			if perr != nil {
				ferr := fmt.Errorf("failed to get Sentry Projects: %w", perr)
				p.logger.Error(ferr.Error())
				hasErrors = true
				continue
			}
			for _, prj := range pageProjs {
				if prj.Organization != nil && prj.Organization.Slug != nil && *prj.Organization.Slug == *o.Slug {
					sprojs = append(sprojs, prj)
				}
			}
			for link != nil && link.Next.Results {
				pageProjs = nil
				link, perr = p.sentry.GetPage(link.Next, &pageProjs)
				if perr != nil {
					ferr := fmt.Errorf("failed to get Sentry Projects page: %w", perr)
					p.logger.Error(ferr.Error())
					hasErrors = true
					break
				}
				for _, prj := range pageProjs {
					if prj.Organization != nil && prj.Organization.Slug != nil && *prj.Organization.Slug == *o.Slug {
						sprojs = append(sprojs, prj)
					}
				}
			}
		} else {
			var perr error
			var link *sentryAPI.Link
			sprojs, link, perr = p.sentry.GetOrgProjects(o)
			if perr != nil {
				ferr := fmt.Errorf("failed to get Sentry Projects: %w", perr)
				p.logger.Error(ferr.Error())
				hasErrors = true
				continue
			}
			for link != nil && link.Next.Results {
				var pageProjs []sentryAPI.Project
				link, perr = p.sentry.GetPage(link.Next, &pageProjs)
				if perr != nil {
					ferr := fmt.Errorf("failed to get Sentry Projects page: %w", perr)
					p.logger.Error(ferr.Error())
					hasErrors = true
					break
				}
				sprojs = append(sprojs, pageProjs...)
			}
		}
		p.logger.Info("projects fetched", "organization", *o.Slug, "count", len(sprojs))

		for _, prj := range sprojs {
			p.logger.Info("syncing project", "organization", *o.Slug, "project", *prj.Slug)
			_, err := p.projects.Create(ctx, *o.Slug, sentry.ToProject(prj))
			if err != nil {
				ferr := fmt.Errorf("failed to create Project: %w", err)
				p.logger.Error(ferr.Error())
				hasErrors = true
				continue
			}

			var (
				statsPeriod   *string = nil
				shortIDLookup *bool   = nil
				query         *string = nil
			)
			p.logger.Info("fetching issues", "organization", *o.Slug, "project", *prj.Slug)
			issues, link, err := p.sentry.GetIssues(o, prj, statsPeriod, shortIDLookup, query)
			if err != nil {
				ferr := fmt.Errorf("failed to get Sentry Issues: %w", err)
				p.logger.Error(ferr.Error())
				hasErrors = true
				continue
			}
			for link != nil && link.Next.Results {
				var pageIssues []sentryAPI.Issue
				link, err = p.sentry.GetPage(link.Next, &pageIssues)
				if err != nil {
					ferr := fmt.Errorf("failed to get Sentry Issues page: %w", err)
					p.logger.Error(ferr.Error())
					hasErrors = true
					break
				}
				issues = append(issues, pageIssues...)
			}
			p.logger.Info("issues fetched", "organization", *o.Slug, "project", *prj.Slug, "count", len(issues))

			for _, is := range issues {
				_, err := p.issues.Create(ctx, *o.Slug, *prj.Slug, sentry.ToIssue(is))
				if err != nil {
					ferr := fmt.Errorf("failed to create Issue: %w", err)
					p.logger.Error(ferr.Error())
					hasErrors = true
					continue
				}
			}
		}
	}

	if hasErrors {
		p.logger.Info("resync completed with errors")
		p.setStatus(Error)
	} else {
		p.logger.Info("resync completed")
		p.setStatus(Ok)
	}
}

func (p *Plugin) ListAllIssues(ctx context.Context) ([]issue.IssueWithRelations, error) {
	return p.issues.ListAll(ctx)
}

func (p *Plugin) setStatus(s Status) {
	p.mxStatus.Lock()
	defer p.mxStatus.Unlock()

	p.status = s
}
