SHELL = /bin/bash

# Build settings
GO_ASMFLAGS =
GO_GCFLAGS =
GO_BUILD_ARGS = $(GO_GCFLAGS) $(GO_ASMFLAGS) -trimpath

# Tool settings
BIN_DIR = ./bin
GOLANGCI_LINT_VERSION = v2.8.0
GOVULNCHECK_VERSION = latest

export PATH := $(PWD)/$(BIN_DIR):$(PATH)

.PHONY: all
all: test lint vulncheck

.PHONY: help
help: ## Show available targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-15s %s\n", $$1, $$2}'


##@ Development

.PHONY: fix
fix: $(BIN_DIR)/golangci-lint ## Fixup files in the repo
	go mod tidy
	go fmt ./...
	$(BIN_DIR)/golangci-lint run --fix

.PHONY: lint
lint: $(BIN_DIR)/golangci-lint ## Run the lint check
	$(BIN_DIR)/golangci-lint run

.PHONY: vulncheck
vulncheck: $(BIN_DIR)/govulncheck ## Run vulnerability check
	$(BIN_DIR)/govulncheck ./...

.PHONY: clean
clean: ## Cleanup tool binaries
	rm -rvf $(BIN_DIR)


##@ Build

.PHONY: install
install: ## Install the library
	go install $(GO_BUILD_ARGS) ./...


##@ Test

.PHONY: test
test: ## Run unit tests
	go test -cover ./...

.PHONY: cover-report
cover-report: ## Generate and open coverage report
	go test -cover -coverprofile coverage.out ./...
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out


##@ Tools

$(BIN_DIR)/golangci-lint:
	@mkdir -p $(BIN_DIR)
	curl -sSfL https://golangci-lint.run/install.sh | sh -s -- $(GOLANGCI_LINT_VERSION)

$(BIN_DIR)/govulncheck:
	@mkdir -p $(BIN_DIR)
	GOBIN=$(PWD)/$(BIN_DIR) go install golang.org/x/vuln/cmd/govulncheck@$(GOVULNCHECK_VERSION)
