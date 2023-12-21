#! /usr/bin/make -f

# Project variables.
PROJECT_NAME = 'ignite apps'

## govet: Run go vet.
govet:
	@echo Running go vet...
	@go list -f '{{.Dir}}/...' -m | xargs go vet 

## govulncheck: Run govulncheck
govulncheck:
	@echo Running govulncheck...
	@go list -f '{{.Dir}}/...' -m | grep -v integration |xargs go run golang.org/x/vuln/cmd/govulncheck

## format: Install and run goimports and gofumpt
format:
	@echo Formatting...
	@go run mvdan.cc/gofumpt -w .
	@go run golang.org/x/tools/cmd/goimports -w -local github.com/ignite/apps .

## lint: Run Golang CI Lint.
lint:
	@echo Running golangci-lint...
	@go list -f '{{.Dir}}/...' -m | xargs go run github.com/golangci/golangci-lint/cmd/golangci-lint run --out-format=tab --issues-exit-code=0

.PHONY: govet format lint

## test-unit: Run the unit tests.
test-unit:
	@echo Running unit tests...
	@go list -f '{{.Dir}}/...' -m | xargs go test -race -failfast -v

## test-integration: Run the integration tests.
test-integration: install
	@echo Running integration tests...
	@go test -race -failfast -v -timeout 60m ./integration/...

## test: Run unit and integration tests.
test: govet govulncheck test-unit test-integration

.PHONY: test-unit test-integration test

help: Makefile
	@echo
	@echo "\n Choose a command run in "$(PROJECT_NAME)", or just run 'make' for install\n"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: help
