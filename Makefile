SHELL := /bin/sh

GO ?= go
BINARY ?= tamcp
CMD ?= ./cmd/tamcp
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.4.0")
COMMIT ?= $(shell git rev-parse HEAD 2>/dev/null || echo "unknown")
BUILD_DATE ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ || echo "unknown")

.PHONY: help run build init server test fmt vet tidy clean bump

help: ## Show this help message
	@awk 'BEGIN {FS = ":.*## "}; /^[a-zA-Z0-9_.-]+:.*## / {printf "  %-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

run: ## Run the tamcp server (dev mode)
	$(GO) run $(CMD) server

init: ## Initialize tamcp configuration files
	$(GO) run $(CMD) init

build: ## Build the tamcp binary with version ldflags
	$(GO) build \
		-ldflags "\
			-X github.com/rangertaha/tamcp/internal.Version=$(VERSION) \
			-X github.com/rangertaha/tamcp/internal.Commit=$(COMMIT) \
			-X github.com/rangertaha/tamcp/internal.BuildDate=$(BUILD_DATE)" \
		-o ./bin/$(BINARY) \
		$(CMD)

test: ## Run all tests
	$(GO) test ./...

fmt: ## Format Go source files
	$(GO) fmt ./...

vet: ## Run go vet checks
	$(GO) vet ./...

tidy: ## Tidy module dependencies
	$(GO) mod tidy

clean: ## Remove build artifacts and init files
	$(GO) run $(CMD) init --clean || true
	rm -rf ./bin/

bump: ## Bump version: make bump [v=major|minor|patch] (default: patch)
	@v=$${v:-patch}; \
	latest=$$(git describe --tags --abbrev=0 2>/dev/null || echo "0.0.0"); \
	major=$$(echo $$latest | sed 's/^v//' | cut -d. -f1); \
	minor=$$(echo $$latest | sed 's/^v//' | cut -d. -f2); \
	patch=$$(echo $$latest | sed 's/^v//' | cut -d. -f3); \
	case $$v in \
		major) major=$$((major + 1)); minor=0; patch=0 ;; \
		minor) minor=$$((minor + 1)); patch=0 ;; \
		patch) patch=$$((patch + 1)) ;; \
		*) echo "usage: make bump v=major|minor|patch"; exit 1 ;; \
	esac; \
	next="v$$major.$$minor.$$patch"; \
	echo "$$latest -> $$next"; \
	git tag -a "$$next" -m "Release $$next: $$v bump from $$latest"; \
	echo "tagged $$next"

server: fmt vet test clean build ## Format, vet, test, build and run
	./bin/$(BINARY) server
