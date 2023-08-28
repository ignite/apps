.PHONY: build mocks format lint test test-unit govet govulncheck

build:
	@echo Building...
	@go build -o ./cli-plugin-hermes .

## mocks: generate mocks
mocks:
	@echo Generating mocks
	@go install github.com/vektra/mockery/v2
	@go generate ./...

## format: Install and run goimports and gofumpt
format:
	@echo Formatting...
	@go run mvdan.cc/gofumpt -w .
	@go run golang.org/x/tools/cmd/goimports -w -local github.com/ignite/cli-plugin-hermes .

## lint: Run Golang CI Lint.
lint:
	@echo Running gocilint...
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint run --out-format=tab --issues-exit-code=0

test: govet test-unit

## govet: Run go vet.
govet:
	@echo Running go vet...
	@go vet ./...

test-unit:
	@echo Running unit tests...
	@go test -race -failfast -v ./...

## govulncheck: Run govulncheck
govulncheck:
	@echo Running govulncheck...
	@go run golang.org/x/vuln/cmd/govulncheck ./...
