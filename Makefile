# Makefile for scalekit-sdk-go: build, lint, test, and vulnerability check.
# Requires: go, golangci-lint (https://golangci-lint.run/).

.PHONY: build lint test vuln all

# Default target
all: build lint test vuln

# Build compiles all packages.
build:
	go build ./...

# Lint runs golangci-lint and go vet.
lint:
	golangci-lint run --new-from-rev=origin/main --color=always
	go vet ./...

test:
	go test $(TEST_FLAGS) ./...

# Vuln runs govulncheck for vulnerability scanning.
vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...