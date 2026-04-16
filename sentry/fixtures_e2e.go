package sentry

import (
	"time"

	sentry "github.com/atlassian/go-sentry-api"
)

// E2EFixturesKey is the sentinel value for SENTRY_API_KEY that activates fixture mode.
const E2EFixturesKey = "e2e"

// FixtureService implements Service and returns hardcoded data for e2e testing.
type FixtureService struct{}

func ptr[T any](v T) *T { return &v }

var fixtureOrg = sentry.Organization{
	ID:   ptr("1"),
	Name: "Youdeploy Org",
	Slug: ptr("youdeploy-org"),
}

var fixtureProjects = []sentry.Project{
	{ID: "101", Name: "YouDeploy API", Slug: ptr("youdeploy-api"), Status: "active"},
	{ID: "102", Name: "YouDeploy Frontend", Slug: ptr("youdeploy-frontend"), Status: "active"},
	{ID: "103", Name: "YouDeploy Worker", Slug: ptr("youdeploy-worker"), Status: "active"},
}

// fixtureIssueDef holds the static parts of a fixture issue definition.
type fixtureIssueDef struct {
	suffix    string
	title     string
	level     string
	hasSeen   bool
	userCount int
}

// fixtureIssuesByProjectID maps project ID to its specific set of issues.
var fixtureIssuesByProjectID = map[string][]fixtureIssueDef{
	"101": {
		{"001", "NullPointerException in PaymentService", "error", false, 42},
		{"002", "Timeout connecting to database", "warning", true, 7},
		{"003", "HTTP 502 Bad Gateway from upstream", "error", false, 13},
	},
	"102": {
		{"001", "Unhandled promise rejection in auth flow", "error", false, 28},
		{"002", "React render error in Dashboard component", "error", true, 5},
	},
	"103": {
		{"001", "Memory usage exceeded threshold", "warning", true, 3},
		{"002", "Job queue backed up: retries exhausted", "error", false, 19},
		{"003", "Deadlock detected in task scheduler", "error", false, 8},
		{"004", "Worker failed to connect to Redis", "warning", true, 11},
	},
}

func (f *FixtureService) GetOrganizations() ([]sentry.Organization, *sentry.Link, error) {
	return []sentry.Organization{fixtureOrg}, nil, nil
}

func (f *FixtureService) GetOrganization(orgSlug string) (sentry.Organization, error) {
	org := fixtureOrg
	org.Slug = ptr(orgSlug)
	return org, nil
}

func (f *FixtureService) GetOrgProjects(o sentry.Organization) ([]sentry.Project, *sentry.Link, error) {
	return fixtureProjects, nil, nil
}

func (f *FixtureService) GetIssues(o sentry.Organization, p sentry.Project, statsPeriod *string, shortIDLookup *bool, query *string) ([]sentry.Issue, *sentry.Link, error) {
	now := time.Now()
	yesterday := now.Add(-24 * time.Hour)
	weekAgo := now.Add(-7 * 24 * time.Hour)
	issueStatus := sentry.Status("unresolved")

	defs := fixtureIssuesByProjectID[p.ID]
	issues := make([]sentry.Issue, 0, len(defs))
	for _, def := range defs {
		// Prefix ID with project ID to keep issues unique across projects.
		id := p.ID + "-" + def.suffix
		firstSeen := weekAgo
		if def.hasSeen {
			firstSeen = yesterday
		}
		issues = append(issues, sentry.Issue{
			ID:        ptr(id),
			Title:     ptr("Fixture: " + def.title),
			Permalink: ptr("https://sentry.io/organizations/" + *o.Slug + "/issues/" + id + "/"),
			HasSeen:   ptr(def.hasSeen),
			FirstSeen: &firstSeen,
			LastSeen:  &now,
			UserCount: ptr(def.userCount),
			Level:     ptr(def.level),
			Status:    &issueStatus,
			Type:      ptr("error"),
		})
	}
	return issues, nil, nil
}
