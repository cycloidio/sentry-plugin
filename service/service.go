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
	return
}

func (p *Plugin) DeletePlugin(ctx context.Context) {
	// NOTE: Nothing to do here
	return
}

func (p *Plugin) Resync(ctx context.Context) {
	select {
	case <-ctx.Done():
		return
	default:
	}
	if !p.started {
		return
	}
	// This will delete everything as all have FK to organizations
	err := p.organizations.DeleteAll(ctx)
	if err != nil {
		ferr := fmt.Errorf("failed to delete all Organizations: %w", err)
		p.logger.Error(ferr.Error())
		p.setStatus(Error)
		return
	}
	sorgs := make([]sentryAPI.Organization, 0)

	if p.config.Sentry.OrganizationSlug != "" {
		o, err := p.sentry.GetOrganization(p.config.Sentry.OrganizationSlug)
		if err != nil {
			ferr := fmt.Errorf("failed to get Sentry Organization: %w", err)
			p.logger.Error(ferr.Error())
			p.setStatus(Error)
			return
		}
		sorgs = append(sorgs, o)
	} else {
		sorgs, _, err = p.sentry.GetOrganizations()
		if err != nil {
			ferr := fmt.Errorf("failed to get Sentry Organizations: %w", err)
			p.logger.Error(ferr.Error())
			p.setStatus(Error)
			return
		}
	}

	for _, o := range sorgs {
		_, err := p.organizations.Create(ctx, sentry.ToOrganization(o))
		if err != nil {
			ferr := fmt.Errorf("failed to create Organization: %w", err)
			p.logger.Error(ferr.Error())
			p.setStatus(Error)
			continue
		}

		sprojs, _, err := p.sentry.GetOrgProjects(o)
		if err != nil {
			ferr := fmt.Errorf("failed to get Sentry Projects: %w", err)
			p.logger.Error(ferr.Error())
			p.setStatus(Error)
			continue
		}

		for _, prj := range sprojs {
			_, err := p.projects.Create(ctx, *o.Slug, sentry.ToProject(prj))
			if err != nil {
				ferr := fmt.Errorf("failed to create Project: %w", err)
				p.logger.Error(ferr.Error())
				p.setStatus(Error)
				continue
			}

			var (
				statsPeriod   *string = nil
				shortIDLookup *bool   = nil
				query         *string = nil
			)
			issues, _, err := p.sentry.GetIssues(o, prj, statsPeriod, shortIDLookup, query)
			if err != nil {
				ferr := fmt.Errorf("failed to get Sentry Issues: %w", err)
				p.logger.Error(ferr.Error())
				p.setStatus(Error)
				continue
			}

			for _, is := range issues {
				_, err := p.issues.Create(ctx, *o.Slug, *prj.Slug, sentry.ToIssue(is))
				if err != nil {
					ferr := fmt.Errorf("failed to create Issue: %w", err)
					p.logger.Error(ferr.Error())
					p.setStatus(Error)
					continue
				}
			}
		}
	}

	return
}

func (p *Plugin) setStatus(s Status) {
	p.mxStatus.Lock()
	defer p.mxStatus.Unlock()

	p.status = s
}
