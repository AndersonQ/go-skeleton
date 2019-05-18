.PHONY: all deps build test check

GO ?= go
GOLINT ?= golint

deps:
	$(GO) mod tidy

build:
	$(GO) build  -ldflags="-X 'github.com/AndersonQ/go-skeleton/handlers.version=$$(git rev-parse HEAD)' -X 'github.com/AndersonQ/go-skeleton/handlers.buildTime=$$(date -R)'"

test: deps
	$(GO) test -cover ./...

check:
	$(GO) vet $$($(GO) list ./...)
	$(GOLINT) $$($(GO) list ./...)