# Makefile for scalekit-sdk-go: build, lint, test, and vulnerability check.
# Requires: go, golangci-lint (https://golangci-lint.run/).

PROTO_REF := v0.1.120.2

.PHONY: build lint test vuln all generate

# Default target
all: build lint test vuln

# generate regenerates gRPC stubs from the published BSR module.
generate:
	buf generate buf.build/scalekit/scalekit:$(PROTO_REF) --include-imports

# Build compiles all packages.
build:
	go build ./...

# Lint runs golangci-lint and go vet.
lint:
	golangci-lint run --new-from-rev=origin/main --color=always
	go vet ./...

test:
	go test -race -v ./...

# Vuln runs govulncheck for vulnerability scanning.
vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...
