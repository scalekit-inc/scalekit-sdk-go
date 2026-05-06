# Makefile for scalekit-sdk-go: build, lint, test, and vulnerability check.
# Requires: go, golangci-lint (https://golangci-lint.run/).

LOCAL_PROTO_DIR ?= ../scalekit/proto

.PHONY: build lint test vuln all generate generate-local

# Default target
all: build lint test vuln

# generate regenerates gRPC stubs from the published BSR module.
generate:
	buf generate buf.build/scalekit/scalekit --include-imports

# generate-local regenerates gRPC stubs from a local proto checkout.
generate-local:
	buf generate ../scalekit
	rm -rf pkg/grpc/google pkg/grpc/buf pkg/grpc/protoc-gen-openapiv2

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
