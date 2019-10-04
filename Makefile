.PHONY: all deps build test check

GO ?= go
# "As of Go 1.13, the go command by default downloads and authenticates modules using the
# Go module mirror and Go checksum database run by Google"
# see go help module-private for detail
GOPRIVATE ?=
GOLINT ?= golint
ENV ?= local

GITURL ?= git@github.com

deps:
	GOPRIVATE=${GOPRIVATE} $(GO) mod download

# Use a github token or similar to access private repos
# In this example it's got no efffexct rather than changing from ssh to https
gitconfig:
	git config --global url.${GITURL}.insteadOf https://github.com


build: deps
	ENV=${ENV} CGO_ENABLED=0 $(GO) build -ldflags="-X 'github.com/AndersonQ/go-skeleton/handlers.version=$$(git rev-parse HEAD)' -X 'github.com/AndersonQ/go-skeleton/handlers.buildTime=$$(date -R)'"

test: deps
	$(GO) test -cover ./...

check:
	$(GO) vet $$($(GO) list ./...)
	$(GOLINT) $$($(GO) list ./...)