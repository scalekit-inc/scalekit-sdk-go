# Changelog

All notable changes to this SDK are documented in this file. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and versions follow
[Semantic Versioning](https://semver.org/).

Sections up to and including 2.9.0 were imported from [GitHub Releases](https://github.com/scalekit-inc/scalekit-sdk-go/releases). They keep their original wording.

## [2.9.0] - 2026-10-05

### Changes

- chore: add saif-at-scalekit to CODEOWNERS ([#91](https://github.com/scalekit-inc/scalekit-sdk-go/pull/91))
- [SK-1980] Add resources client for listing and revoking user consents ([#90](https://github.com/scalekit-inc/scalekit-sdk-go/pull/90))
- [SK-2043] chore: update proto to v0.1.150.0 (v2.8.1) ([#94](https://github.com/scalekit-inc/scalekit-sdk-go/pull/94))
- feat: accept multiple issuers in ValidateTokenWithOptions (SK-2080) ([#95](https://github.com/scalekit-inc/scalekit-sdk-go/pull/95))
- Add Create/Update/Delete/List/Get to ResourceClient ([#92](https://github.com/scalekit-inc/scalekit-sdk-go/pull/92))

## [2.8.0] - 2026-07-28

### Changes

- [SK-1339] feat: typed UpdateLoginUserDetails response, issuer validation, listEventsPaginated ([#88](https://github.com/scalekit-inc/scalekit-sdk-go/pull/88))

## [2.7.0] - 2026-07-17

### Changes

- docs: Update README with agent-first positioning ([#68](https://github.com/scalekit-inc/scalekit-sdk-go/pull/68))
- Update go-jose and golang.org/x dependencies ([#71](https://github.com/scalekit-inc/scalekit-sdk-go/pull/71))
- feat: organization session policy SDK methods ([#70](https://github.com/scalekit-inc/scalekit-sdk-go/pull/70))
- feat: regenerate protos from v0.1.123.0 and add slug/logo_url to CreateOrganizationOptions ([#77](https://github.com/scalekit-inc/scalekit-sdk-go/pull/77))
- feat(users): add external_id lookup methods ([#83](https://github.com/scalekit-inc/scalekit-sdk-go/pull/83))
- chore: update proto to v0.1.137.0 ([#85](https://github.com/scalekit-inc/scalekit-sdk-go/pull/85))

## [2.6.0] - 2026-03-30

### Changes

- docs: expand REFERENCE.md for clients, API tokens, and missing methods ([#66](https://github.com/scalekit-inc/scalekit-sdk-go/pull/66))

## [2.5.0] - 2026-03-30

### Changes

- docs: rename reference.md to REFERENCE.md ([#64](https://github.com/scalekit-inc/scalekit-sdk-go/pull/64))
- [SK-2660] feat(token): add UpdateToken method ([#60](https://github.com/scalekit-inc/scalekit-sdk-go/pull/60))
- fix: method additions, CI hardening, and version bump to 2.3.0 ([#61](https://github.com/scalekit-inc/scalekit-sdk-go/pull/61))
- [SK-2664] feat(m2m): add M2MService (ClientService org-client CRUD) ([#62](https://github.com/scalekit-inc/scalekit-sdk-go/pull/62))
- [SK-2668] feat(users): add ListUserRoles and ListUserPermissions ([#63](https://github.com/scalekit-inc/scalekit-sdk-go/pull/63))
- [SK-2671] feat(roles): add UpdateDefaultRoles and ListDependentRoles ([#65](https://github.com/scalekit-inc/scalekit-sdk-go/pull/65))
- ci: add release workflow with dry-run support ([#67](https://github.com/scalekit-inc/scalekit-sdk-go/pull/67))

## [2.4.0] - 2026-03-12

### Changes

#### Added
- Multi-app client management support with `Client()` APIs for client CRUD and client secret management.
- PKCE helper support via `GeneratePKCEConfiguration(...)` for public client authorization flows.
- `WithSecret(...)` to derive a secret-enabled client without mutating the original client.
- `ValidateTokenWithOptions(...)` with audience and scope validation support.

#### Changed
- `NewScalekitClient(...)` now supports initialization without a client secret for public clients while preserving backward compatibility.
- `GenerateClientToken(...)` now accepts structured options and returns token metadata.
- `RefreshAccessToken(...)` now includes `id_token` in the response.
- SDK and API versions were updated, along with regenerated gRPC/Connect client code.

#### Fixed
- `AuthenticateWithCode(...)` and `RefreshAccessToken(...)` now send `client_secret` only when configured.
- Clear fail-fast errors are now returned when authenticated API calls require a client secret but none is configured.

## [2.3.0] - 2026-03-10

### Changes

- Document DeleteMembership cascade param
- feat: add DeleteRoleBase for env-level role inheritance removal
- feat: change ListPermissions signature to accept pageToken and pageSize
- feat add global ListUsers to Go UserService interface
- feat: add ValidateToken, GenerateClientToken, GetClientAccessToken to Scalekit interface
- chore: bump version to 2.3.0, api-version to 20260310
- fix: surface deleteEnvRole failures instead of swallowing them

## [2.2.0] - 2026-03-06

### Changes

#### Release Notes  v2.2.0

**Release date:** 2026-03-06
**Previous release:** v2.1.0 (2026-02-24)

---

#### PRs included in this release

| # | Title |
|---|-------|
| [#45](https://github.com/scalekit-inc/scalekit-sdk-go/pull/45) | API tokens |
| [#50](https://github.com/scalekit-inc/scalekit-sdk-go/pull/50) | chore: add CODEOWNERS |
| [#58](https://github.com/scalekit-inc/scalekit-sdk-go/pull/58) | Structured error handling, sentinel errors, and thread-safe client |

---

#### New Features

#### API Token Management (`token.go`)

A new `TokenClient` is available on `ScalekitClient` for managing API tokens programmatically:

- **Create tokens** - user-scoped or machine tokens with optional expiry
- **List tokens** - filter by user ID or token type
- **Revoke/delete tokens** - revoke individual tokens or all tokens for a user
- **Validate tokens** - validate a raw token string; returns structured claims

```go
// Example
token, err := client.Token().CreateToken(ctx, scalekit.CreateTokenOptions{...})
```

#### Structured Error Handling

Non-2xx HTTP responses are now surfaced as `*scalekit.Error` with an inspectable `StatusCode` field:

```go
var e *scalekit.Error
if errors.As(err, &e) {
    fmt.Println(e.StatusCode) // e.g. 403
}
```

#### Sentinel Errors

Callers can now use `errors.Is` instead of string matching for common failure modes:

| Sentinel | Condition |
|----------|-----------|
| `ErrTokenRequired` | empty token passed to ValidateToken |
| `ErrTokenValidationFailed` | token validation connect error |
| `ErrMissingExpClaim` | token lacks exp claim |
| `ErrCodeOrLinkTokenRequired` | missing code/link token |
| `ErrOrganizationIdRequired` | missing org ID |
| `ErrDirectoryNotFound` | directory lookup returned nothing |

---

#### Improvements

#### Thread Safety

- `accessToken` and `jsonWebKeySet` now use `atomic.Pointer` with singleflight groups — eliminates data races on concurrent requests.
- JWKS read lock removed in favour of atomic load.

#### HTTP Client Hardening

- A 10s default timeout is attached to every outbound request via `withDefaultTimeout`. Callers with their own `context.Deadline` are unaffected.
- Auth retry is now scoped to genuine 401 / `CodeUnauthenticated` responses only; `authenticateClient` failures are propagated instead of silently ignored.

#### API Version Header

- `x-api-version` header updated to `20260226`.

---

#### Breaking Changes

**None.** A potentially breaking change (removal of external ID field) was reverted before this release. All existing method signatures remain compatible with v2.1.0.

---

#### Internal / Chore

- Added `.github/CODEOWNERS` to enforce required reviews on all PRs (owners: @AkshayParihar33, @dhawani).
- Proto stubs regenerated from source; `ListTokensRequest.UserId` changed to `*string` pointer.
- `go.mod` updated.

---

#### Upgrade

```bash
go get github.com/scalekit-inc/scalekit-sdk-go/v2@v2.2.0
```

## [2.1.0] - 2026-02-24

### Changes

#### v2.1.0

#### ⚠️ Breaking Changes

All methods that make network calls now require a `context.Context` as the first argument to support cancellation or timeout propagation.

| Method | Old Signature | New Signature |
|--------|--------------|---------------|
| `AuthenticateWithCode` | `(code, redirectUri, options)` | `(ctx, code, redirectUri, options)` |
| `ValidateAccessToken` | `(accessToken)` | `(ctx, accessToken)` |
| `GetAccessTokenClaims` | `(accessToken)` | `(ctx, accessToken)` |
| `GetIdpInitiatedLoginClaims` | `(token)` | `(ctx, token)` |
| `RefreshAccessToken` | `(refreshToken)` | `(ctx, refreshToken)` |
| `ValidateToken` | `(token, jwksFn)` | `(ctx, token, jwksFn)` |

**Migration:** Pass `ctx` as the first argument.

```go
// Before
scalekitClient.AuthenticateWithCode(token)

// After
scalekitClient.AuthenticateWithCode(ctx, token)
```

---

#### ✨ New Features

- **`Connection.CreateConnection`** — programmatically create SSO connections for an organization
- **`Connection.DeleteConnection`** — delete an existing SSO connection
- **`Directory.CreateDirectory`** — programmatically create a directory for an organization
- **`Directory.DeleteDirectory`** — delete an existing directory

## [2.0.11] - 2026-02-21

### Changes

- add refresh token in Authenticate with code ([#47](https://github.com/scalekit-inc/scalekit-sdk-go/pull/47))

## [2.0.10] - 2026-01-14

### Changes

- Sk 2390/connection docs ([#42](https://github.com/scalekit-inc/scalekit-sdk-go/pull/42))

#### Domain List API - Release Summary

#### Overview
Enhanced the `ListDomains` API with filtering and pagination capabilities to provide more flexible and efficient domain management operations.

#### Changes

#### New Features

#### 1. Domain Type Filtering
The `ListDomains` API now supports filtering domains by type:
- **ALLOWED_EMAIL_DOMAIN**: Filter for allowed email domains only
- **ORGANIZATION_DOMAIN**: Filter for organization domains only
- **DOMAIN_TYPE_UNSPECIFIED**: Unspecified domain type

**Usage Example:**
```go
// List only organization domains
orgDomains, err := client.Domain().ListDomains(ctx, organizationId, &scalekit.ListDomainOptions{
    DomainType: scalekit.DomainTypeOrganization,
})

// List only allowed email domains
allowedEmailDomains, err := client.Domain().ListDomains(ctx, organizationId, &scalekit.ListDomainOptions{
    DomainType: scalekit.DomainTypeAllowedEmail,
})
```

#### 2. Pagination Support
Added pagination controls to manage large result sets:
- **PageSize**: Number of domains to return per page (default: 100)
- **PageNumber**: Page number to retrieve (1-indexed)

**Usage Example:**
```go
// List domains with custom pagination
domains, err := client.Domain().ListDomains(ctx, organizationId, &scalekit.ListDomainOptions{
    PageSize:   50,
    PageNumber: 1,
})
```

#### 3. Combined Filtering and Pagination
Filter and paginate simultaneously for efficient domain queries:

**Usage Example:**
```go
// List organization domains with pagination
orgDomains, err := client.Domain().ListDomains(ctx, organizationId, &scalekit.ListDomainOptions{
    DomainType: scalekit.DomainTypeOrganization,
    PageSize:   25,
    PageNumber: 1,
})
```

#### API Signature

**Before:**
```go
ListDomains(ctx context.Context, organizationId string) (*ListDomainResponse, error)
```

**After:**
```go
ListDomains(ctx context.Context, organizationId string, options ...*ListDomainOptions) (*ListDomainResponse, error)
```

#### Backward Compatibility
✅ **Fully backward compatible** - The API can still be called without options:
```go
// Still works - returns all domains with default page size of 10
allDomains, err := client.Domain().ListDomains(ctx, organizationId)
```

#### New Types

#### `ListDomainOptions`
```go
type ListDomainOptions struct {
    DomainType DomainType  // Optional: Filter by domain type
    PageSize   uint32      // Optional: Number of results per page (default: 10)
    PageNumber uint32      // Optional: Page number to retrieve
}
```

#### Implementation Details

#### Default Behavior
- When no options are provided, the API returns all domains with a default page size of 10
- Domain type filtering is optional and can be omitted
- Pagination parameters are optional; if not specified, default page size applies

#### Type Conversion
The SDK automatically converts string domain type constants to the appropriate gRPC enum values:
- `scalekit.DomainTypeAllowedEmail` → `domains.DomainType_ALLOWED_EMAIL_DOMAIN`
- `scalekit.DomainTypeOrganization` → `domains.DomainType_ORGANIZATION_DOMAIN`
- `scalekit.DomainTypeUnspecified` → `domains.DomainType_DOMAIN_TYPE_UNSPECIFIED`

#### Testing

Comprehensive test coverage has been added in `test/domain_test.go`:
- ✅ Test listing all domains (no filter)
- ✅ Test filtering by `ORGANIZATION_DOMAIN` type
- ✅ Test filtering by `ALLOWED_EMAIL_DOMAIN` type
- ✅ Test verification that filtered results only contain domains of the specified type
- ✅ Test that created domains appear in the list

#### Migration Guide

#### No Changes Required
Existing code continues to work without modification:
```go
// Existing code - no changes needed
domains, err := client.Domain().ListDomains(ctx, organizationId)
```

#### Optional: Add Filtering
To take advantage of new filtering capabilities:
```go
// New: Filter by domain type
orgDomains, err := client.Domain().ListDomains(ctx, organizationId, &scalekit.ListDomainOptions{
    DomainType: scalekit.DomainTypeOrganization,
})
```

#### Benefits

1. **Performance**: Filtering reduces unnecessary data transfer and processing
2. **Scalability**: Pagination enables handling of large domain lists efficiently
3. **Flexibility**: Combined filtering and pagination for precise queries
4. **Backward Compatibility**: No breaking changes to existing implementations

#### Related Files

- `domain.go`: Core implementation with `ListDomainOptions` and enhanced `ListDomains` method
- `test/domain_test.go`: Comprehensive test suite including `TestListDomains` function
- `pkg/grpc/scalekit/v1/domains/`: gRPC protocol buffer definitions

## [2.0.9] - 2025-12-23

### Changes

- Add WebAuthn sdk methods  in https://github.com/scalekit-inc/scalekit-sdk-go/pull/41

## [2.0.8] - 2025-12-01

### Changes

- Added support for Bring your own Auth for MCP Servers in https://github.com/scalekit-inc/scalekit-sdk-go/pull/39

## [2.0.7] - 2025-11-19

### Changes

**New sdk methods**

- Add user management settings to organization API ([#38](https://github.com/scalekit-inc/scalekit-sdk-go/pull/38))

## [2.0.6] - 2025-11-18

### Changes

Generate proto files to support given_name and family_name in user object  in https://github.com/scalekit-inc/scalekit-sdk-go/pull/37

## [2.0.5] - 2025-10-29

### Changes

- Added Session management sdk methods & interceptorpayload verification method https://github.com/scalekit-inc/scalekit-sdk-go/pull/34

## [2.0.4] - 2025-09-22

### Changes

- Add CreateDomainOptions to domain creation ([#30](https://github.com/scalekit-inc/scalekit-sdk-go/pull/30))
- Add roles and permissions support to SDK ([#35](https://github.com/scalekit-inc/scalekit-sdk-go/pull/35))

## [2.0.3] - 2025-08-26

### Changes

- Update Scalekit SDK import paths to v2 ([#29](https://github.com/scalekit-inc/scalekit-sdk-go/pull/29))

## [2.0.2] - 2025-08-21

### Changes

- Add ResendInvite method to UserService and Regenerate protobuf files ([#28](https://github.com/scalekit-inc/scalekit-sdk-go/pull/28))

## [2.0.1] - 2025-07-16

### Changes

- Added passwordless sdk methods ([#26](https://github.com/scalekit-inc/scalekit-sdk-go/pull/26))

## [2.0.0] - 2025-07-13

### Changes

- upgrade crypto version to v0.39.0 ([#25](https://github.com/scalekit-inc/scalekit-sdk-go/pull/25))
- Added required sdk methods for fullstack auth ([#22](https://github.com/scalekit-inc/scalekit-sdk-go/pull/22))

## [1.0.5] - 2025-06-04

### Changes

- Adding raw token claims and handle verify webhook token ([#23](https://github.com/scalekit-inc/scalekit-sdk-go/pull/23))

## [1.0.4] - 2025-03-13

### Changes

- Remove Get and Delete Portal Link Methods
- add Minimum requirements in readme
- Add Integration tests
- Update  go-jose version

## [1.0.3] - 2024-12-14

### Changes

- Add Directory Sync Support  https://github.com/scalekit-inc/scalekit-sdk-go/pull/11
- Fix crypto Vulnerabiliy  https://github.com/scalekit-inc/scalekit-sdk-go/pull/15

## [1.0.2] - 2024-08-09

### Changes

- Bump github.com/go-jose/go-jose/v4 from 4.0.3 to 4.0.4 ([#6](https://github.com/scalekit-inc/scalekit-sdk-go/pull/6))
- Bump github.com/grpc-ecosystem/grpc-gateway/v2 from 2.20.0 to 2.21.0 ([#5](https://github.com/scalekit-inc/scalekit-sdk-go/pull/5))
- Bump buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go from 1.34.2-20240508200655-46a4cf4ba109.2 to 1.34.2-20240717164558-a6c49f84cc0f.2 ([#4](https://github.com/scalekit-inc/scalekit-sdk-go/pull/4))
- IDP Initiated SSO DX v2 ([#7](https://github.com/scalekit-inc/scalekit-sdk-go/pull/7))

## [1.0.1] - 2024-07-18

### Changes

- Bump github.com/go-jose/go-jose/v4 from 4.0.2 to 4.0.3 ([#2](https://github.com/scalekit-inc/scalekit-sdk-go/pull/2))
- Add missing method and fix method signature ([#3](https://github.com/scalekit-inc/scalekit-sdk-go/pull/3))

## [1.0.0] - 2024-07-17

### Changes

- First Release of the official Scalekit Golang SDK

[2.9.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.9.0
[2.8.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.8.0
[2.7.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.7.0
[2.6.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.6.0
[2.5.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.5.0
[2.4.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.4.0
[2.3.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.3.0
[2.2.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.2.0
[2.1.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.1.0
[2.0.11]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.11
[2.0.10]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.10
[2.0.9]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.9
[2.0.8]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.8
[2.0.7]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.7
[2.0.6]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.6
[2.0.5]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.5
[2.0.4]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.4
[2.0.3]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.3
[2.0.2]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.2
[2.0.1]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.1
[2.0.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v2.0.0
[1.0.5]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v1.0.5
[1.0.4]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v1.0.4
[1.0.3]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v1.0.3
[1.0.2]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v1.0.2
[1.0.1]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v1.0.1
[1.0.0]: https://github.com/scalekit-inc/scalekit-sdk-go/releases/tag/v1.0.0
