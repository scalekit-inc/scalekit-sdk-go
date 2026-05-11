# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Session start

At the beginning of each session, if changes to code need to be made, ask the user: "Do you want me to cut a new branch from `main`, or continue working on the existing branch (`<current-branch>`)?"

## Running commands

When running any of the following, delegate to a **Haiku subagent** (`model: "haiku"`) via the Agent tool:
- `make generate`, `make test`, `make lint`
- `go test ...`, `golangci-lint ...`, `govulncheck ...`

## Commands

```bash
# Install protoc generators and linting tools
make setup

# Run linter (golangci-lint or gofmt fallback)
make lint

# Run all tests + vulnerability check (requires env vars)
make test

# Run a single test
go test -count=1 -v ./test/ -run TestName

# Regenerate protobuf/gRPC code (requires PAT_TOKEN for private proto repo)
make generate

# Check for vulnerabilities
govulncheck ./...
```

### Test environment variables

Integration tests require a `.env` file in `test/` or the repo root, or these env vars set:
- `SCALEKIT_ENVIRONMENT_URL`
- `SCALEKIT_CLIENT_ID`
- `SCALEKIT_CLIENT_SECRET`

## Architecture

This is a Go SDK for Scalekit — a B2B authentication/authorization platform. The SDK wraps ConnectRPC (HTTP/2 gRPC-compatible) calls behind clean Go interfaces.

### Layer structure

1. **`core.go`** — HTTP client with OAuth client-credentials auth. Caches the access token (singleflight prevents stampede) and JWKS for token validation.

2. **`connect.go`** — Generic `connectExecuter[TRequest, TResponse]` that wraps every RPC call with: header injection (user-agent, sdk-version, api-version, Authorization), and automatic retry on 401 with token refresh.

3. **`scalekit.go`** — Root `Scalekit` interface and `scalekitClient` implementation. Houses OAuth flows (authorization code, PKCE), token validation, webhook verification, and the sub-service accessors.

4. **Service files** (`connection.go`, `directory.go`, `domain.go`, `organization.go`, `users.go`, `role.go`, `permission.go`, `clients.go`, `token.go`, `sessions.go`, `passwordless.go`, `webauthn.go`) — Each wraps a generated gRPC client and implements a single exported interface.

5. **`pkg/grpc/scalekit/v1/`** — Generated protobuf/gRPC stubs. Never edit by hand; regenerate via `make generate`.

### Key patterns

- Every service method body follows: `newConnectExecuter(client, method).exec(ctx, req)`
- Token validation uses a generic `ValidateToken[T]` function in `scalekit.go` for typed claim unmarshaling.
- Errors in `errors.go` use sentinel values (`errors.New`) with `errors.Is` for matching. Some errors are joined for backward-compatible error trees.

### Constructor

```go
client := scalekit.NewScalekitClient(envUrl, clientId, clientSecret)
// or with secret chaining
client := scalekit.NewScalekitClient(envUrl, clientId).WithSecret(clientSecret)
```

### Code generation

Protobuf definitions live in a separate private repo (`github.com/scalekit-inc/scalekit`, ref pinned in `buf.gen.yaml`). Generated code is checked in under `pkg/grpc/`. To regenerate, `PAT_TOKEN` must be set with access to the proto repo.
