SHELL := /bin/bash

GO := go
TOOLS_BIN := $(CURDIR)/.tools/bin
GO_TOOLCHAIN := $(shell awk '/^toolchain / { print $$2; exit } /^go / { print "go"$$2; exit }' go.mod)

PROTO_REPO_URL := https://github.com/scalekit-inc/scalekit.git
PROTO_REF := v0.1.103
PROTO_SUBDIR := proto
PROTO_REMOTE_INPUT := $(PROTO_REPO_URL)\#ref=$(PROTO_REF),subdir=$(PROTO_SUBDIR)
PROTO_LOCAL_PATH ?= /Users/akshayparihar/Documents/repos/scalekit/proto

BUF := PATH="$(TOOLS_BIN):$$PATH" buf

.PHONY: setup generate generate_local lint test tools-check

setup:
	@mkdir -p "$(TOOLS_BIN)"
	GOTOOLCHAIN="$(GO_TOOLCHAIN)" GOBIN="$(TOOLS_BIN)" $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
	GOTOOLCHAIN="$(GO_TOOLCHAIN)" GOBIN="$(TOOLS_BIN)" $(GO) install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.19.1
	GOTOOLCHAIN="$(GO_TOOLCHAIN)" GOBIN="$(TOOLS_BIN)" $(GO) install github.com/bufbuild/buf/cmd/buf@v1.50.1
	GOTOOLCHAIN="$(GO_TOOLCHAIN)" GOBIN="$(TOOLS_BIN)" $(GO) install golang.org/x/vuln/cmd/govulncheck@latest
	GOTOOLCHAIN="$(GO_TOOLCHAIN)" GOBIN="$(TOOLS_BIN)" $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

tools-check:
	@command -v "$(TOOLS_BIN)/buf" >/dev/null 2>&1 || (echo "missing buf. run 'make setup'" && exit 1)
	@command -v "$(TOOLS_BIN)/protoc-gen-go" >/dev/null 2>&1 || (echo "missing protoc-gen-go v1.33.0. run 'make setup'" && exit 1)
	@command -v "$(TOOLS_BIN)/protoc-gen-connect-go" >/dev/null 2>&1 || (echo "missing protoc-gen-connect-go. run 'make setup'" && exit 1)
	@command -v "$(TOOLS_BIN)/govulncheck" >/dev/null 2>&1 || (echo "missing govulncheck. run 'make setup'" && exit 1)

generate: tools-check
	$(BUF) generate "$(PROTO_REMOTE_INPUT)"

generate_local: tools-check
	$(BUF) generate "$(PROTO_LOCAL_PATH)"

lint:
	@if [ -x "$(TOOLS_BIN)/golangci-lint" ]; then \
		PATH="$(TOOLS_BIN):$$PATH" golangci-lint run --new-from-rev=origin/main --color=always; \
	else \
		echo "golangci-lint not installed; running gofmt + go vet fallback"; \
		test -z "$$(gofmt -l .)"; \
		$(GO) vet ./...; \
	fi

test: tools-check lint
	$(GO) test -count=1 -v ./...
	PATH="$(TOOLS_BIN):$$PATH" govulncheck ./...
