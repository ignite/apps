#! /usr/bin/make -f

# Project variables.
PROJECT_NAME = 'ignite apps'

## mocks: generate mocks
mocks:
	@echo Generating mocks
	@go install github.com/vektra/mockery/v2@latest
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
        if [ -e "$$dir/go.mod" ]; then \
            echo "Running go generate in $$dir"; \
			cd "$$dir" && mockery; \
        fi \
    done


## goget: Run go get for all apps.
goget:
	@echo Running go get $(REPO)...
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
        if [ -e "$$dir/go.mod" ]; then \
            echo "Running go get $(REPO) in $$dir"; \
            cd "$$dir" && go get $(REPO); \
        fi \
    done

## modtidy: Run go mod tidy for all apps.
modtidy:
	@echo Running go mod tidy...
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
        if [ -e "$$dir/go.mod" ]; then \
            echo "Running go mod tidy in $$dir"; \
            cd "$$dir" && go mod tidy; \
        fi \
    done

## govet: Run go vet for all apps.
govet:
	@echo Running go vet...
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
        if [ -e "$$dir/go.mod" ]; then \
            echo "Running go vet in $$dir"; \
            cd "$$dir" && go vet ./...; \
        fi \
    done

## govulncheck: Run govulncheck for all apps.
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

## lint: Run Golang Lint for all apps.
lint:
	@command -v golangci-lint >/dev/null 2>&1 || { \
        echo "Installing golangci-lint..."; \
        curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.3; \
    }
	@echo Running golangci-lint...
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
        if [ -e "$$dir/go.mod" ]; then \
            echo "Running golangci-lint in $$dir"; \
            cd "$$dir" && golangci-lint run; \
        fi \
    done

## lint-ci: Run Golang CI Lint for all apps.
lint-ci:
	@command -v golangci-lint >/dev/null 2>&1 || { \
        echo "Installing golangci-lint..."; \
        curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.3; \
    }
	@echo Running golangci-lint...
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
        if [ -e "$$dir/go.mod" ]; then \
            echo "Running golangci-lint in $$dir"; \
            cd "$$dir" && golangci-lint run --out-format=tab --issues-exit-code=0; \
        fi \
    done

## ci: Run CI pipeline govet, govulncheck and lint
ci: govet govulncheck lint-ci

## format: Install and run goimports and gofumpt for all apps.
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

## test-unit: Run unit tests for all apps.
test-unit:
	@echo Running unit tests...
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d -not -path '*/integration*'); do \
        if [ -e "$$dir/go.mod" ]; then \
            echo "Running unit tests in $$dir"; \
            cd "$$dir" && go test -race -failfast -v -coverpkg=./... $(go list ./... | grep -v integration); \
        fi \
    done

## test-integration: Run the integration tests.
test-integration:
	@for dir in $$(find $$(pwd -P) -mindepth 1 -maxdepth 4 -type d); do \
        if [ -e "$$dir/go.mod" ] && [ -d "$$dir/integration" ]; then \
			echo "Running integration tests in $$dir"; \
			cd "$$dir" && go test -race -failfast -v -timeout 60m ./integration/...; \
		fi \
	done

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
