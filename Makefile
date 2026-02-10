ifneq ("$(wildcard .env)","")
	include .env
endif
export SENTRY_API_KEY
export SENTRY_ENDPOINT


GOCACHE := $(shell go env GOCACHE)
export GOCACHE

ifeq ($(ARGS),)
	ARGS := ./...
endif

VERSION=0.0.1

.PHONY: help
help: Makefile ## This help dialog
	@IFS=$$'\n' ; \
	help_lines=(`grep -F -h "##" $(MAKEFILE_LIST) | grep -F -v grep -F | sed -e 's/\\$$//'`); \
	for help_line in $${help_lines[@]}; do \
		IFS=$$'#' ; \
		help_split=($$help_line) ; \
		help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
		printf "%-30s %s\n" $$help_command $$help_info ; \
	done

.PHONY: gen
gen: ## Generates code
	@go generate ./...

.PHONY: test
test: ## Tests the Plugin
	@docker-compose -f docker/docker-compose.yml -f docker/develop.yml run -e "SENTRY_API_KEY=$(SENTRY_API_KEY)" -e "SENTRY_ENDPOINT=$(SENTRY_ENDPOINT)" --name sentry_plugin --rm -p 4000:4000 sentry_plugin go test $(ARGS)

.PHONY: docker-release
docker-release: ## Builds the base Docker image for the registry
	@docker build -f ./docker/Dockerfile -t cycloid/sentry-plugin:$(VERSION) .
	@docker push cycloid/sentry-plugin:$(VERSION)
