package main

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cycloidio/sentry-plugin/config"
	"github.com/cycloidio/sentry-plugin/issue"
	"github.com/cycloidio/sentry-plugin/organization"
	"github.com/cycloidio/sentry-plugin/project"
	"github.com/cycloidio/sentry-plugin/sentry"
	"github.com/cycloidio/sentry-plugin/service"
	thttp "github.com/cycloidio/sentry-plugin/service/transport/http"
	"github.com/cycloidio/sentry-plugin/sqlite"
	"github.com/gorilla/handlers"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	var started = true

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	cfg, err := config.Load()
	if err != nil {
		started = false
		logger.Error(fmt.Errorf("failed to load config: %w", err).Error())
	}

	// By default we use the 'memory' setting so for testing we can easily use it
	q := "file::memory:?cache=shared&_foreign_keys=true"
	if cfg.DB.File != "" {
		q = cfg.DB.File + "?_foreign_keys=true"
	}
	db, err := sql.Open("sqlite3", q)
	if err != nil {
		started = false
		logger.Error(fmt.Errorf("could not connect to the SQLite database: %w", err).Error())
	}

	var (
		or organization.Repository
		pr project.Repository
		ir issue.Repository

		ss sentry.Service
	)
	if started {
		or = sqlite.NewOrganizationRepository(db)
		pr = sqlite.NewProjectRepository(db)
		ir = sqlite.NewIssueRepository(db)

		ss, err = sentry.New(cfg.Sentry.APIKey, cfg.Sentry.Endpoint)
		if err != nil {
			started = false
		}
	}
	sctx, cfn := context.WithCancel(context.TODO())
	defer cfn()

	s := service.New(sctx, or, pr, ir, ss, started, cfg, logger)

	handler := thttp.Handler(s)

	mux := http.NewServeMux()
	mux.Handle("/", handler)

	svr := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: handlers.LoggingHandler(os.Stdout, mux),
	}

	errs := make(chan error)

	go func() {
		logger.Info("started server", "port", cfg.Port)
		errs <- svr.ListenAndServe()
	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Info("exit", "reason", <-errs)
}
