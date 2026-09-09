SHELL := /bin/bash

GO := go
TOOLS_BIN := $(CURDIR)/.tools/bin
GO_TOOLCHAIN := go1.25.13

PROTO_REPO_URL := https://github.com/scalekit-inc/scalekit.git
PROTO_REF := v0.1.147.0
PROTO_SUBDIR := proto
PROTO_REMOTE_INPUT := $(PROTO_REPO_URL)\#ref=$(PROTO_REF),subdir=$(PROTO_SUBDIR)
PROTO_LOCAL_PATH ?= ../scalekit/proto #or paste absolute path to protos

BUF := PATH="$(TOOLS_BIN):$$PATH" buf

.PHONY: setup generate generate_local lint test tools-check

setup:
	@mkdir -p "$(TOOLS_BIN)"
	GOTOOLCHAIN="$(GO_TOOLCHAIN)" GOBIN="$(TOOLS_BIN)" $(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.33.0
	GOTOOLCHAIN="$(GO_TOOLCHAIN)" GOBIN="$(TOOLS_BIN)" $(GO) install connectrpc.com/connect/cmd/protoc-gen-connect-go@v1.19.1
	GOTOOLCHAIN="$(GO_TOOLCHAIN)" GOBIN="$(TOOLS_BIN)" $(GO) install github.com/bufbuild/buf/cmd/buf@v1.50.1
	# govulncheck is version-pinned like every other tool here, and must stay that
	# way. It used to track @latest, which broke `make setup` the moment x/vuln
	# v1.8.0 raised its own go directive to 1.26.0 past this repo's toolchain.
	# Unsetting GOTOOLCHAIN does not help in CI: actions/setup-go exports
	# GOTOOLCHAIN=local whenever go-version-file is used, so toolchain switching
	# is already disabled by the environment. v1.7.0 is the last release whose go
	# directive is 1.25.x. Raise it only alongside GO_TOOLCHAIN.
	GOTOOLCHAIN="$(GO_TOOLCHAIN)" GOBIN="$(TOOLS_BIN)" $(GO) install golang.org/x/vuln/cmd/govulncheck@v1.7.0
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
