# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What This Is

A Cycloid plugin that ETLs data from Sentry (organizations → projects → issues) into a local SQLite database and exposes HTTP endpoints for the Cycloid platform to interact with.

## Commands

```bash
# Build
go build -o sentry main.go

# Run tests (requires Docker — tests run inside container)
make test
make test ARGS="./service"   # specific package

# Regenerate mocks and enum string methods
make gen
```

## Environment Variables

| Variable | Default | Required |
|----------|---------|----------|
| `SENTRY_API_KEY` | — | Yes |
| `SENTRY_ENDPOINT` | `http://sentry.io/api/0/` | No |
| `SENTRY_ORGANIZATION_SLUG` | — | No (fetches all orgs if unset) |
| `DB_FILE` | — | No (in-memory SQLite if unset) |
| `PORT` | `8080` | No |

## Architecture

### Data Flow

`main.go` initializes DB → applies `schema.sql` (via `//go:embed`) → creates `service.Plugin` → starts HTTP server + background goroutine that calls `Resync()` immediately.

`Resync()` in `service/service.go`:
1. Deletes all organizations (CASCADE deletes projects and issues via FK constraints)
2. Fetches orgs from Sentry API (or one specific org if slug configured)
3. For each org: fetches projects; for each project: fetches issues
4. Writes everything via repository interfaces

### HTTP Endpoints (`service/transport/http/handler.go`)

- `GET /_cy/ping` — returns status JSON
- `POST /_cy/resync` — manually triggers `Resync()`
- `POST /_cy/events` — stub
- `DELETE /_cy/plugin` — stub

### Layering

- **`sentry/`** — wraps `atlassian/go-sentry-api` client; contains `ToOrganization()`, `ToProject()`, `ToIssue()` converters
- **`organization/`, `project/`, `issue/`** — domain models + repository interfaces
- **`sqlite/`** — SQLite implementations of those repository interfaces
- **`mock/`** — generated mocks (do not edit manually; regenerate with `make gen`)
- **`service/`** — orchestration, status management, HTTP transport

### Key Patterns

**Repository pattern with DI**: `Plugin` struct receives all repositories and the Sentry client via `New()`. No global state.

**Status enum** (`service/status.go`): `Ok`, `Syncthing`, `Error` — protected by `sync.RWMutex`. Status is set to `Error` on Resync failures but the loop continues to next item.

**Enums with codegen**: `service.Status`, `event.Severity`, `event.Type`, `event.Color` use `github.com/dmarkham/enumer`. Add `//go:generate` directives and run `make gen`.

**NULL helpers in `sqlite/`**: `toNullString()`, `toNullBool()`, `toNullInt64()`, `toNullTime()` convert zero/empty values to SQL NULLs.

### Platform Integration Files

- `manifest.yaml` — plugin metadata, config options, and relation definitions for Cycloid
- `widgets.yaml` — widget queries/selections rendered in the Cycloid UI
- `schema.sql` — embedded in the binary and applied at startup
