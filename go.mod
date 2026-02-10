module github.com/cycloidio/sentry-plugin

go 1.25.1

require (
	github.com/atlassian/go-sentry-api v1.0.0
	github.com/cycloidio/sqlr v1.0.0
	github.com/gorilla/handlers v1.5.2
	github.com/mattn/go-sqlite3 v1.14.33
	github.com/stretchr/testify v1.9.0
	go-simpler.org/env v0.12.0
	go.uber.org/mock v0.6.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dmarkham/enumer v1.6.3 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/pascaldekloe/name v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.27.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/tools v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

tool (
	github.com/dmarkham/enumer
	go.uber.org/mock/mockgen
)

replace github.com/atlassian/go-sentry-api v1.0.0 => github.com/cycloidio/go-sentry-api v1.0.0-cy
