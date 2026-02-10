package sentry

import (
	"fmt"

	sentry "github.com/atlassian/go-sentry-api"
	"github.com/cycloidio/sentry-plugin/issue"
	"github.com/cycloidio/sentry-plugin/organization"
	"github.com/cycloidio/sentry-plugin/project"
)

//go:generate go tool mockgen -destination=../mock/sentry_service.go -mock_names=Service=SentryService -package mock github.com/cycloidio/sentry-plugin/sentry Service

type Service interface {
	GetOrganizations() ([]sentry.Organization, *sentry.Link, error)
	GetOrganization(orgslug string) (sentry.Organization, error)

	GetOrgProjects(o sentry.Organization) ([]sentry.Project, *sentry.Link, error)

	GetIssues(o sentry.Organization, p sentry.Project, statsPeriod *string, shortIDLookup *bool, query *string) ([]sentry.Issue, *sentry.Link, error)
}

func New(apik, ep string) (*sentry.Client, error) {
	pep := &ep
	if ep == "" {
		pep = nil
	}
	c, err := sentry.NewClient(apik, pep, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Sentry client: %w", err)
	}
	return c, nil
}

func ToOrganization(o sentry.Organization) organization.Organization {
	return organization.Organization{
		ID:   *o.ID,
		Name: o.Name,
		Slug: *o.Slug,
	}
}

func ToProject(p sentry.Project) project.Project {
	return project.Project{
		ID:     p.ID,
		Name:   p.Name,
		Slug:   *p.Slug,
		Status: p.Status,
	}
}

func ToIssue(i sentry.Issue) issue.Issue {
	return issue.Issue{
		ID:        *i.ID,
		Title:     *i.Title,
		Permalink: *i.Permalink,

		HasSeen:   *i.HasSeen,
		FirstSeen: *i.FirstSeen,
		LastSeen:  *i.LastSeen,
		UserCount: *i.UserCount,

		Level:  *i.Level,
		Status: string(*i.Status),
		Type:   *i.Type,
	}
}
