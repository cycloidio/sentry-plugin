package mock

import (
	context "context"
	"log/slog"
	"os"

	"github.com/cycloidio/sentry-plugin/config"
	"github.com/cycloidio/sentry-plugin/service"
	gomock "go.uber.org/mock/gomock"
)

type MockRegistry struct {
	Organizations *OrganizationRepository
	Projects      *ProjectRepository
	Issues        *IssueRepository
	Sentry        *SentryService

	Config  *config.Config
	Started bool

	S service.Service
}

func NewService(ctx context.Context, ctrl *gomock.Controller) MockRegistry {
	or := NewOrganizationRepository(ctrl)
	pr := NewProjectRepository(ctrl)
	ir := NewIssueRepository(ctrl)
	ss := NewSentryService(ctrl)
	started := true

	cfg := &config.Config{}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	s := service.New(ctx, or, pr, ir, ss, started, cfg, logger)
	return MockRegistry{
		Organizations: or,
		Projects:      pr,
		Issues:        ir,
		Sentry:        ss,

		Config: cfg,

		Started: started,

		S: s,
	}
}
