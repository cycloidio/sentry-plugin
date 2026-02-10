package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/cycloidio/sentry-plugin/issue"
	"github.com/cycloidio/sentry-plugin/mock"
	"github.com/cycloidio/sentry-plugin/organization"
	"github.com/cycloidio/sentry-plugin/project"
	"github.com/cycloidio/sentry-plugin/service"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	sentryAPI "github.com/atlassian/go-sentry-api"
)

func TestImplements(t *testing.T) {
	assert.Implements(t, (*service.Service)(nil), new(service.Plugin))
}

func TestPing(t *testing.T) {
	// Cannot change it live so we can only test the OK
	t.Run("Ok", func(t *testing.T) {
		// We do this because the Service starts and starts fetching,
		// so it may start when we are testing and there are random errors
		// so we start it with a cancel ctx
		sctx, cfn := context.WithCancel(context.TODO())
		cfn()

		var (
			ctrl = gomock.NewController(t)
			s    = mock.NewService(sctx, ctrl)

			ctx = context.Background()
		)

		ok := s.S.Ping(ctx)
		assert.Equal(t, ok, service.Ok)
	})
}

func TestResync(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// We do this because the Service starts and starts fetching,
		// so it may start when we are testing and there are random errors
		// so we start it with a cancel ctx
		sctx, cfn := context.WithCancel(context.TODO())
		cfn()
		var (
			ctrl = gomock.NewController(t)
			s    = mock.NewService(sctx, ctrl)

			ctx = context.Background()
			org = organization.Organization{
				ID:   "123",
				Name: "My Org",
				Slug: "my_org",
			}
			sorg = sentryAPI.Organization{
				ID:   &org.ID,
				Name: org.Name,
				Slug: &org.Slug,
			}

			prj = project.Project{
				ID:     "456",
				Name:   "My Proj",
				Slug:   "my_proj",
				Status: "some_status",
			}
			sprj = sentryAPI.Project{
				ID:     prj.ID,
				Name:   prj.Name,
				Slug:   &prj.Slug,
				Status: prj.Status,
			}

			iss = issue.Issue{
				ID:        "789",
				Title:     "My Issue",
				Permalink: "https://google.com",
				HasSeen:   false,
				FirstSeen: time.Now(),
				LastSeen:  time.Now(),
				UserCount: 10,
				Level:     "level",
				Status:    "unresolved",
				Type:      "type",
			}
			issstatus = sentryAPI.Status(iss.Status)
			siss      = sentryAPI.Issue{
				ID:        &iss.ID,
				Title:     &iss.Title,
				Permalink: &iss.Permalink,
				HasSeen:   &iss.HasSeen,
				FirstSeen: &iss.FirstSeen,
				LastSeen:  &iss.LastSeen,
				UserCount: &iss.UserCount,
				Level:     &iss.Level,
				Status:    &issstatus,
				Type:      &iss.Type,
			}
		)

		s.Organizations.EXPECT().DeleteAll(ctx).Return(nil)
		s.Sentry.EXPECT().GetOrganizations().Return([]sentryAPI.Organization{sorg}, nil, nil)
		s.Organizations.EXPECT().Create(ctx, org).Return(uint32(1), nil)

		s.Sentry.EXPECT().GetOrgProjects(sorg).Return([]sentryAPI.Project{sprj}, nil, nil)
		s.Projects.EXPECT().Create(ctx, org.Slug, prj).Return(uint32(1), nil)

		s.Sentry.EXPECT().GetIssues(sorg, sprj, nil, nil, nil).Return([]sentryAPI.Issue{siss}, nil, nil)
		s.Issues.EXPECT().Create(ctx, org.Slug, prj.Slug, iss).Return(uint32(1), nil)

		s.S.Resync(ctx)
	})
	t.Run("OneConfigOrganization", func(t *testing.T) {
		// We do this because the Service starts and starts fetching,
		// so it may start when we are testing and there are random errors
		// so we start it with a cancel ctx
		sctx, cfn := context.WithCancel(context.TODO())
		cfn()
		var (
			ctrl = gomock.NewController(t)
			s    = mock.NewService(sctx, ctrl)

			ctx = context.Background()
			org = organization.Organization{
				ID:   "123",
				Name: "My Org",
				Slug: "my_org",
			}
			sorg = sentryAPI.Organization{
				ID:   &org.ID,
				Name: org.Name,
				Slug: &org.Slug,
			}

			prj = project.Project{
				ID:     "456",
				Name:   "My Proj",
				Slug:   "my_proj",
				Status: "some_status",
			}
			sprj = sentryAPI.Project{
				ID:     prj.ID,
				Name:   prj.Name,
				Slug:   &prj.Slug,
				Status: prj.Status,
			}

			iss = issue.Issue{
				ID:        "789",
				Title:     "My Issue",
				Permalink: "https://google.com",
				HasSeen:   false,
				FirstSeen: time.Now(),
				LastSeen:  time.Now(),
				UserCount: 10,
				Level:     "level",
				Status:    "unresolved",
				Type:      "type",
			}
			issstatus = sentryAPI.Status(iss.Status)
			siss      = sentryAPI.Issue{
				ID:        &iss.ID,
				Title:     &iss.Title,
				Permalink: &iss.Permalink,
				HasSeen:   &iss.HasSeen,
				FirstSeen: &iss.FirstSeen,
				LastSeen:  &iss.LastSeen,
				UserCount: &iss.UserCount,
				Level:     &iss.Level,
				Status:    &issstatus,
				Type:      &iss.Type,
			}
		)

		// We force the OrganizationSlug so we only fetch that one
		s.Config.Sentry.OrganizationSlug = org.Slug

		s.Organizations.EXPECT().DeleteAll(ctx).Return(nil)
		s.Sentry.EXPECT().GetOrganization(org.Slug).Return(sorg, nil)
		s.Organizations.EXPECT().Create(ctx, org).Return(uint32(1), nil)

		s.Sentry.EXPECT().GetOrgProjects(sorg).Return([]sentryAPI.Project{sprj}, nil, nil)
		s.Projects.EXPECT().Create(ctx, org.Slug, prj).Return(uint32(1), nil)

		s.Sentry.EXPECT().GetIssues(sorg, sprj, nil, nil, nil).Return([]sentryAPI.Issue{siss}, nil, nil)
		s.Issues.EXPECT().Create(ctx, org.Slug, prj.Slug, iss).Return(uint32(1), nil)

		s.S.Resync(ctx)
	})
	t.Run("WithCancel", func(t *testing.T) {
		// We do this because the Service starts and starts fetching,
		// so it may start when we are testing and there are random errors
		// so we start it with a cancel ctx
		sctx, cfn := context.WithCancel(context.TODO())
		cfn()
		var (
			ctrl = gomock.NewController(t)
			s    = mock.NewService(sctx, ctrl)

			ctx = context.Background()
		)
		nctx, ncfn := context.WithCancel(ctx)
		ncfn()

		s.S.Resync(nctx)
	})
	t.Run("!Started", func(t *testing.T) {
		// We do this because the Service starts and starts fetching,
		// so it may start when we are testing and there are random errors
		// so we start it with a cancel ctx
		sctx, cfn := context.WithCancel(context.TODO())
		cfn()
		var (
			ctrl = gomock.NewController(t)
			s    = mock.NewService(sctx, ctrl)

			ctx = context.Background()
		)
		nctx, ncfn := context.WithCancel(ctx)
		ncfn()

		s.S.Resync(nctx)
	})
}
