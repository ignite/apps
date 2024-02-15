#! /usr/bin/make -f

# Project variables.
PROJECT_NAME = 'ignite apps'

## govet: Run go vet.
govet:
	@echo Running go vet...
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
        if [ -e "$$dir/go.mod" ]; then \
            echo "Running go vet in $$dir"; \
            cd "$$dir" && go vet ./...; \
        fi \
    done

## govulncheck: Run govulncheck
govulncheck:
	@command -v govulncheck >/dev/null 2>&1 || { \
        echo "Installing govulncheck..."; \
        go install golang.org/x/vuln/cmd/govulncheck@latest; \
    }
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
        if [ -e "$$dir/go.mod" ]; then \
            echo "Running go vet in $$dir"; \
            cd "$$dir" && govulncheck ./...; \
        fi \
    done
	@echo Running govulncheck...

## lint: Run Golang CI Lint.
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
        echo "Installing golangci-lint..."; \
        curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.42.1; \
    }
	@echo Running golangci-lint...
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
        if [ -e "$$dir/go.mod" ]; then \
            echo "Running golangci-lint in $$dir"; \
            cd "$$dir" && golangci-lint run --out-format=tab --issues-exit-code=0; \
        fi \
    done

## check: Run govet, govulncheck and lint
check: govet govulncheck lint

## format: Install and run goimports and gofumpt
format:
	@echo Formatting...
	@command -v gofumpt >/dev/null 2>&1 || { \
        echo "Installing gofumpt..."; \
        go install mvdan.cc/gofumpt@latest; \
    }
	@command -v goimports >/dev/null 2>&1 || { \
        echo "Installing goimports..."; \
        go install golang.org/x/tools/cmd/goimports@latest; \
    }
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
		if [ -e "$$dir/go.mod" ]; then \
			echo "Running format in $$dir"; \
			cd "$$dir" && gofumpt -w .; \
			cd "$$dir" && goimports -w -local github.com/ignite/apps .; \
		fi \
	done

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
