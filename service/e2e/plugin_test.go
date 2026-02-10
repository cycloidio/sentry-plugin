package e2e_test

import (
	"context"
	"database/sql"
	"log/slog"
	"os"
	"testing"

	"github.com/cycloidio/sentry-plugin/config"
	"github.com/cycloidio/sentry-plugin/sentry"
	"github.com/cycloidio/sentry-plugin/service"
	"github.com/cycloidio/sentry-plugin/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlugin(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load()
	require.NoError(t, err)

	// By default we use the 'memory' setting so for testing we can easily use it
	q := "file::memory:?cache=shared&_foreign_keys=true"
	db, err := sql.Open("sqlite3", q)
	require.NoError(t, err)
	started := true

	// TODO: Run the migrations!!
	b, err := os.ReadFile("../../schema.sql")
	require.NoError(t, err)

	_, err = db.Exec(string(b))
	require.NoError(t, err)

	or := sqlite.NewOrganizationRepository(db)
	pr := sqlite.NewProjectRepository(db)
	ir := sqlite.NewIssueRepository(db)
	ss, err := sentry.New(cfg.Sentry.APIKey, cfg.Sentry.Endpoint)
	require.NoError(t, err)

	// We do this because the Service starts and starts fetching,
	// so it may start when we are testing and there are random errors
	// so we start it with a cancel ctx
	sctx, cfn := context.WithCancel(context.TODO())
	cfn()
	s := service.New(sctx, or, pr, ir, ss, started, cfg, logger)

	ctx := context.Background()
	s.Resync(ctx)
	assert.Equal(t, service.Ok.String(), s.Ping(ctx).String())
}
