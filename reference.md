# Reference

## ScalekitClient

<details><summary><code>client := <a href="scalekit.go">scalekit.NewScalekitClient</a>(envUrl, clientId, clientSecret) -> scalekit.Scalekit</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a new Scalekit Go SDK client.

The returned client provides high-level OAuth helpers (authorization URL, token exchange, token validation, webhook verification) and typed sub-clients for resource APIs (organizations, users, sessions, etc.).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
import "github.com/scalekit-inc/scalekit-sdk-go/v2"

client := scalekit.NewScalekitClient(
  "<SCALEKIT_ENV_URL>",
  "<SCALEKIT_CLIENT_ID>",
  "<SCALEKIT_CLIENT_SECRET>",
)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**envUrl:** `string` - Your Scalekit environment URL (from the dashboard)

</dd>
</dl>

<dl>
<dd>

**clientId:** `string` - Scalekit client ID

</dd>
</dl>

<dl>
<dd>

**clientSecret:** `string` - Scalekit client secret

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">GetAuthorizationUrl</a>(redirectUri, options) -> (*url.URL, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Utility method to generate the OAuth 2.0 authorization URL to initiate the SSO authentication flow.

This method doesn't make any network calls. It returns a fully formed authorization URL that you can redirect users to.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
authURL, err := client.GetAuthorizationUrl(
  "https://yourapp.com/auth/callback",
  scalekit.AuthorizationUrlOptions{
    State:          "random-state-value",
    OrganizationId: "org_123",
  },
)
if err != nil {
  // handle
}

// Redirect user to authURL.String()
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**redirectUri:** `string` - The URL where users will be redirected after authentication (must match a configured redirect URI)

</dd>
</dl>

<dl>
<dd>

**options:** `AuthorizationUrlOptions` - Configuration for the authorization request
- `Scopes []string` - OAuth scopes to request (default: `openid profile email`)
- `State string` - Opaque value to maintain state between request and callback
- `Nonce string` - String value used to associate a client session with an ID Token
- `LoginHint string` - Hint about the login identifier the user might use
- `DomainHint string` - Domain hint to identify which organization's IdP to use
- `ConnectionId string` - Specific SSO connection ID to use for authentication
- `OrganizationId string` - Organization ID to authenticate against
- `Provider string` - Social login provider (e.g., `google`, `github`, `microsoft`)
- `CodeChallenge string` - PKCE code challenge
- `CodeChallengeMethod string` - Method used to generate the code challenge (e.g., `S256`)
- `Prompt string` - Controls authentication behavior (e.g., `login`, `consent`)

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">AuthenticateWithCode</a>(code, redirectUri, options) -> (*AuthenticationResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Exchanges an authorization code for tokens and user information.

Call this in your redirect handler after receiving the `code` query parameter.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.AuthenticateWithCode(
  code,
  "https://yourapp.com/auth/callback",
  scalekit.AuthenticationOptions{},
)
if err != nil {
  // handle
}

accessToken := resp.AccessToken
user := resp.User
_ = accessToken
_ = user
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**code:** `string` - The authorization code received in the callback URL

</dd>
</dl>

<dl>
<dd>

**redirectUri:** `string` - The same redirect URI used in `GetAuthorizationUrl` (must match exactly)

</dd>
</dl>

<dl>
<dd>

**options:** `AuthenticationOptions`
- `CodeVerifier string` - PKCE code verifier (required if PKCE was used)

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">GetIdpInitiatedLoginClaims</a>(idpInitiatedLoginToken) -> (*IdpInitiatedLoginClaims, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Extracts and validates claims from an IdP-initiated login token.

Use this method when handling IdP-initiated SSO flows, where authentication is initiated from the identity provider’s portal instead of your application.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
claims, err := client.GetIdpInitiatedLoginClaims(idpInitiatedLoginToken)
if err != nil {
  // handle
}

// claims.ConnectionID, claims.OrganizationID, claims.LoginHint, claims.RelayState
_ = claims
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**idpInitiatedLoginToken:** `string` - The token received via IdP-initiated login

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">GetAccessTokenClaims</a>(accessToken) -> (*AccessTokenClaims, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Parses and validates an access token and returns its claims.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
claims, err := client.GetAccessTokenClaims(accessToken)
if err != nil {
  // handle
}
_ = claims
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**accessToken:** `string` - The JWT access token

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">ValidateAccessToken</a>(accessToken) -> (bool, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Validates an access token (including expiration checks) and returns whether it is valid.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
ok, err := client.ValidateAccessToken(accessToken)
if err != nil {
  // invalid
}
_ = ok
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**accessToken:** `string` - The JWT access token

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">RefreshAccessToken</a>(refreshToken) -> (*TokenResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Exchanges a refresh token for a new access token (and optionally a new refresh token).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
tokens, err := client.RefreshAccessToken(refreshToken)
if err != nil {
  // handle
}
_ = tokens.AccessToken
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**refreshToken:** `string` - The refresh token

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">GetLogoutUrl</a>(options) -> (*url.URL, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Generates a logout URL for OIDC logout flows.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
logoutURL, err := client.GetLogoutUrl(scalekit.LogoutUrlOptions{
  IdTokenHint:           idToken,
  PostLogoutRedirectUri: "https://yourapp.com/",
  State:                "state",
})
if err != nil {
  // handle
}
_ = logoutURL
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**options:** `LogoutUrlOptions`
- `IdTokenHint string`
- `PostLogoutRedirectUri string`
- `State string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">VerifyWebhookPayload</a>(secret, headers, payload) -> (bool, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Verifies a Scalekit webhook payload signature using `webhook-id`, `webhook-timestamp`, and `webhook-signature` headers.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
valid, err := client.VerifyWebhookPayload(
  "whsec_...",
  map[string]string{
    "webhook-id":        "webhook_123",
    "webhook-timestamp": "1730000000",
    "webhook-signature": "v1,base64sig",
  },
  []byte(`{"event":"user.created","data":{"id":"123"}}`),
)
if err != nil {
  // handle
}
_ = valid
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**secret:** `string` - Your webhook signing secret (e.g. `whsec_...`)

</dd>
</dl>

<dl>
<dd>

**headers:** `map[string]string` - Request headers containing webhook signature fields

</dd>
</dl>

<dl>
<dd>

**payload:** `[]byte` - Raw request body

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">VerifyInterceptorPayload</a>(secret, headers, payload) -> (bool, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Verifies an interceptor payload signature. Uses the same signature format as webhooks.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
valid, err := client.VerifyInterceptorPayload(secret, headers, payload)
if err != nil {
  // handle
}
_ = valid
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**secret:** `string`

</dd>
</dl>

<dl>
<dd>

**headers:** `map[string]string`

</dd>
</dl>

<dl>
<dd>

**payload:** `[]byte`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">Connection</a>() -> scalekit.Connection</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the Connections client (`client.Connection()`), used to manage and query SSO connections.
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">Organization</a>() -> scalekit.Organization</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the Organizations client (`client.Organization()`), used to manage organizations (tenants).
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">User</a>() -> scalekit.UserService</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the Users client (`client.User()`), used to manage users and memberships.
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">Domain</a>() -> scalekit.Domain</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the Domains client (`client.Domain()`), used to manage and query domains for organizations.
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">Directory</a>() -> scalekit.Directory</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the Directories client (`client.Directory()`), used to list directories and directory users/groups (SCIM/Directory Sync).
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">Session</a>() -> scalekit.SessionService</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the Sessions client (`client.Session()`), used to list and revoke sessions.
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">Role</a>() -> scalekit.RoleService</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the Roles client (`client.Role()`), used to manage roles and organization roles.
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">Permission</a>() -> scalekit.PermissionService</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the Permissions client (`client.Permission()`), used to manage permissions and role-permission relationships.
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">Passwordless</a>() -> scalekit.PasswordlessService</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the Passwordless client (`client.Passwordless()`), used for passwordless email flows (OTP / magic link).
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">WebAuthn</a>() -> scalekit.WebAuthnService</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the WebAuthn client (`client.WebAuthn()`), used to manage passkey credentials.
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.<a href="scalekit.go">Auth</a>() -> scalekit.AuthService</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Returns the Auth client (`client.Auth()`), used for Auth gRPC helper methods (e.g. update login user details).
</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## Organizations

<details><summary><code>client.Organization().<a href="organization.go">CreateOrganization</a>(ctx, name, options) -> (*CreateOrganizationResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a new organization (tenant).

Organizations represent your B2B customers. Use `ExternalId` to map Scalekit organizations to your internal identifiers.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
import (
  "context"
  "fmt"

  "github.com/scalekit-inc/scalekit-sdk-go/v2"
)

ctx := context.Background()
org, err := client.Organization().CreateOrganization(ctx, "Acme Corporation", scalekit.CreateOrganizationOptions{
  ExternalId: "customer_12345",
  Metadata: map[string]string{
    "source": "signup",
  },
})
if err != nil {
  // handle
}

fmt.Println("Organization ID:", org.Organization.Id)
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context` - Request context

</dd>
</dl>

<dl>
<dd>

**name:** `string` - Display name for the organization

</dd>
</dl>

<dl>
<dd>

**options:** `CreateOrganizationOptions`
- `ExternalId string`
- `Metadata map[string]string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Organization().<a href="organization.go">ListOrganization</a>(ctx, options) -> (*ListOrganizationsResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Retrieves a paginated list of organizations in your environment.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
orgs, err := client.Organization().ListOrganization(ctx, &scalekit.ListOrganizationOptions{
  PageSize:  10,
  PageToken: "",
})
if err != nil {
  // handle
}

for _, org := range orgs.Organizations {
  _ = org.Id
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**options:** `*ListOrganizationOptions` (alias of `organizationsv1.ListOrganizationsRequest`)
- `PageSize uint32`
- `PageToken string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Organization().<a href="organization.go">GetOrganization</a>(ctx, id) -> (*GetOrganizationResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches an organization by Scalekit organization ID.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
org, err := client.Organization().GetOrganization(ctx, "org_123")
if err != nil {
  // handle
}
_ = org.Organization
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**id:** `string` - Scalekit organization ID

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Organization().<a href="organization.go">GetOrganizationByExternalId</a>(ctx, externalId) -> (*GetOrganizationResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches an organization by your external ID (if set).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
org, err := client.Organization().GetOrganizationByExternalId(ctx, "customer_12345")
if err != nil {
  // handle
}
_ = org.Organization
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**externalId:** `string` - Your system’s organization identifier

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Organization().<a href="organization.go">UpdateOrganization</a>(ctx, id, organization) -> (*UpdateOrganizationResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates an organization by Scalekit organization ID.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
updated, err := client.Organization().UpdateOrganization(
  ctx,
  "org_123",
  &organizations.UpdateOrganization{
    DisplayName: func() *string { s := "Updated name"; return &s }(),
  },
)
if err != nil {
  // handle
}
_ = updated.Organization
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**id:** `string` - Scalekit organization ID

</dd>
</dl>

<dl>
<dd>

**organization:** `*organizationsv1.UpdateOrganization` - Fields to update

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Organization().<a href="organization.go">UpdateOrganizationByExternalId</a>(ctx, externalId, organization) -> (*UpdateOrganizationResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates an organization by your external ID.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
updated, err := client.Organization().UpdateOrganizationByExternalId(
  ctx,
  "customer_12345",
  &organizations.UpdateOrganization{
    DisplayName: func() *string { s := "Updated name"; return &s }(),
  },
)
if err != nil {
  // handle
}
_ = updated.Organization
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**externalId:** `string`

</dd>
</dl>

<dl>
<dd>

**organization:** `*organizationsv1.UpdateOrganization`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Organization().<a href="organization.go">DeleteOrganization</a>(ctx, id) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes an organization by Scalekit organization ID.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
if err := client.Organization().DeleteOrganization(ctx, "org_123"); err != nil {
  // handle
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**id:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Organization().<a href="organization.go">GeneratePortalLink</a>(ctx, organizationId) -> (*organizationsv1.Link, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Generates an admin portal link for an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
link, err := client.Organization().GeneratePortalLink(ctx, "org_123")
if err != nil {
  // handle
}
_ = link.Url
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Organization().<a href="organization.go">UpdateOrganizationSettings</a>(ctx, id, settings) -> (*GetOrganizationResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates organization settings (feature toggles).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.Organization().UpdateOrganizationSettings(ctx, "org_123", scalekit.OrganizationSettings{
  Features: []scalekit.Feature{
    {Name: "sso", Enabled: true},
    {Name: "dir_sync", Enabled: true},
  },
})
if err != nil {
  // handle
}
_ = resp.Organization.Settings
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**id:** `string` - Scalekit organization ID

</dd>
</dl>

<dl>
<dd>

**settings:** `OrganizationSettings`
- `Features []Feature` where `Feature` is `{ Name string; Enabled bool }`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Organization().<a href="organization.go">UpsertUserManagementSettings</a>(ctx, organizationId, settings) -> (*organizationsv1.OrganizationUserManagementSettings, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates or updates user management settings for an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
maxUsers := int32(150)
settings, err := client.Organization().UpsertUserManagementSettings(
  ctx,
  "org_123",
  scalekit.OrganizationUserManagementSettings{
    MaxAllowedUsers: &maxUsers,
  },
)
if err != nil {
  // handle
}
_ = settings.MaxAllowedUsers
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**settings:** `OrganizationUserManagementSettings`
- `MaxAllowedUsers *int32`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## Connections

<details><summary><code>client.Connection().<a href="connection.go">GetConnection</a>(ctx, organizationId, id) -> (*GetConnectionResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches a connection by ID within an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
conn, err := client.Connection().GetConnection(ctx, "org_123", "conn_123")
if err != nil {
  // handle
}
_ = conn.Connection
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**id:** `string` - Connection ID

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Connection().<a href="connection.go">ListConnectionsByDomain</a>(ctx, domain) -> (*ListConnectionsResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists connections that match a given domain (e.g. to support domain discovery for SSO).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
conns, err := client.Connection().ListConnectionsByDomain(ctx, "acme.com")
if err != nil {
  // handle
}
_ = conns.Connections
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**domain:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Connection().<a href="connection.go">ListConnections</a>(ctx, organizationId) -> (*ListConnectionsResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists all connections for an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
conns, err := client.Connection().ListConnections(ctx, "org_123")
if err != nil {
  // handle
}
_ = conns.Connections
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Connection().<a href="connection.go">EnableConnection</a>(ctx, organizationId, id) -> (*ToggleConnectionResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Enables a connection for an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.Connection().EnableConnection(ctx, "org_123", "conn_123")
if err != nil {
  // handle
}
_ = resp.Enabled
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**id:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Connection().<a href="connection.go">DisableConnection</a>(ctx, organizationId, id) -> (*ToggleConnectionResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Disables a connection for an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.Connection().DisableConnection(ctx, "org_123", "conn_123")
if err != nil {
  // handle
}
_ = resp.Enabled
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**id:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## Users

<details><summary><code>client.User().<a href="users.go">ListOrganizationUsers</a>(ctx, organizationId, options) -> (*ListOrganizationUsersResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists users in an organization (paginated).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
users, err := client.User().ListOrganizationUsers(ctx, "org_123", &scalekit.ListUsersOptions{
  PageSize:  10,
  PageToken: "",
})
if err != nil {
  // handle
}
_ = users.Users
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**options:** `*ListUsersOptions`
- `PageSize uint32`
- `PageToken string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.User().<a href="users.go">GetUser</a>(ctx, userId) -> (*GetUserResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches a user by user ID.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
user, err := client.User().GetUser(ctx, "usr_123")
if err != nil {
  // handle
}
_ = user.User
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**userId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.User().<a href="users.go">UpdateUser</a>(ctx, userId, updateUser) -> (*UpdateUserResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates a user by user ID.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
firstName := "Test"
lastName := "User"
name := "Test User"
locale := "en-US"

updated, err := client.User().UpdateUser(ctx, "usr_123", &users.UpdateUser{
  UserProfile: &users.UpdateUserProfile{
    FirstName: &firstName,
    LastName:  &lastName,
    Name:      &name,
    Locale:    &locale,
  },
})
if err != nil {
  // handle
}
_ = updated.User
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**userId:** `string`

</dd>
</dl>

<dl>
<dd>

**updateUser:** `*usersv1.UpdateUser`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.User().<a href="users.go">CreateUserAndMembership</a>(ctx, organizationId, user, sendInvitationEmail) -> (*CreateUserAndMembershipResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a user and adds them to an organization with a membership.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
created, err := client.User().CreateUserAndMembership(ctx, "org_123", &users.CreateUser{
  Email: "test.user@example.com",
  Metadata: map[string]string{
    "source": "test",
  },
}, true)
if err != nil {
  // handle
}
_ = created.User.Id
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**user:** `*usersv1.CreateUser`

</dd>
</dl>

<dl>
<dd>

**sendInvitationEmail:** `bool` - Whether to send an invitation email for the created user

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.User().<a href="users.go">DeleteUser</a>(ctx, userId) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes a user by user ID.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
if err := client.User().DeleteUser(ctx, "usr_123"); err != nil {
  // handle
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**userId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.User().<a href="users.go">CreateMembership</a>(ctx, organizationId, userId, membership, sendInvitationEmail) -> (*CreateMembershipResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a membership for an existing user in an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.User().CreateMembership(ctx, "org_123", "usr_123", &users.CreateMembership{
  Roles: []*commons.Role{{Name: "admin"}},
  Metadata: map[string]string{
    "membership_type": "test",
  },
}, false)
if err != nil {
  // handle
}
_ = resp.User
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**userId:** `string`

</dd>
</dl>

<dl>
<dd>

**membership:** `*usersv1.CreateMembership`

</dd>
</dl>

<dl>
<dd>

**sendInvitationEmail:** `bool`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.User().<a href="users.go">UpdateMembership</a>(ctx, organizationId, userId, membership) -> (*UpdateMembershipResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates a membership for a user within an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.User().UpdateMembership(ctx, "org_123", "usr_123", &users.UpdateMembership{
  Roles: []*commons.Role{{Name: "member"}},
})
if err != nil {
  // handle
}
_ = resp.User
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**userId:** `string`

</dd>
</dl>

<dl>
<dd>

**membership:** `*usersv1.UpdateMembership`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.User().<a href="users.go">DeleteMembership</a>(ctx, organizationId, userId, cascade) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes a membership for a user from an organization.

If `cascade` is true, the API may also delete related resources (behavior depends on backend).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
if err := client.User().DeleteMembership(ctx, "org_123", "usr_123", false); err != nil {
  // handle
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**userId:** `string`

</dd>
</dl>

<dl>
<dd>

**cascade:** `bool`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.User().<a href="users.go">ResendInvite</a>(ctx, organizationId, userId) -> (*usersv1.ResendInviteResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Resends a pending invite for a user in an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.User().ResendInvite(ctx, "org_123", "usr_123")
if err != nil {
  // handle
}
_ = resp.Invite
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**userId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## Domains

<details><summary><code>client.Domain().<a href="domain.go">CreateDomain</a>(ctx, organizationId, name, options?) -> (*CreateDomainResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a domain for an organization.

The SDK supports backward-compatible signatures:
- `CreateDomain(ctx, orgId, domain)` (no options)
- `CreateDomain(ctx, orgId, domain, options)` (with options)
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
// Without options (backward compatible)
created, err := client.Domain().CreateDomain(ctx, "org_123", "acme.com")
if err != nil {
  // handle
}
_ = created.Domain

// With options
created2, err := client.Domain().CreateDomain(ctx, "org_123", "acme.com", &scalekit.CreateDomainOptions{
  DomainType: scalekit.DomainTypeOrganization,
})
if err != nil {
  // handle
}
_ = created2.Domain
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**name:** `string` - Domain name (e.g. `acme.com`)

</dd>
</dl>

<dl>
<dd>

**options?:** `*CreateDomainOptions`
- `DomainType DomainType` - `DOMAIN_TYPE_UNSPECIFIED`, `ALLOWED_EMAIL_DOMAIN`, or `ORGANIZATION_DOMAIN`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Domain().<a href="domain.go">GetDomain</a>(ctx, id, organizationId) -> (*GetDomainResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches a domain by ID within an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
domain, err := client.Domain().GetDomain(ctx, "dom_123", "org_123")
if err != nil {
  // handle
}
_ = domain.Domain
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**id:** `string` - Domain ID

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Domain().<a href="domain.go">ListDomains</a>(ctx, organizationId, options?) -> (*ListDomainResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists domains for an organization (supports optional filtering and pagination).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
// List all domains
all, err := client.Domain().ListDomains(ctx, "org_123")
if err != nil {
  // handle
}
_ = all.Domains

// Filter by domain type
orgDomains, err := client.Domain().ListDomains(ctx, "org_123", &scalekit.ListDomainOptions{
  DomainType: scalekit.DomainTypeOrganization,
})
if err != nil {
  // handle
}
_ = orgDomains.Domains
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**options?:** `*ListDomainOptions`
- `DomainType DomainType`
- `PageSize uint32`
- `PageNumber uint32`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Domain().<a href="domain.go">DeleteDomain</a>(ctx, id, organizationId) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes a domain by ID within an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
if err := client.Domain().DeleteDomain(ctx, "dom_123", "org_123"); err != nil {
  // handle
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**id:** `string`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## Directories

<details><summary><code>client.Directory().<a href="directory.go">ListDirectories</a>(ctx, organizationId) -> (*ListDirectoriesResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists directories for an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
dirs, err := client.Directory().ListDirectories(ctx, "org_123")
if err != nil {
  // handle
}
_ = dirs.Directories
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Directory().<a href="directory.go">GetDirectory</a>(ctx, organizationId, directoryId) -> (*GetDirectoryResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches a directory by ID within an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
dir, err := client.Directory().GetDirectory(ctx, "org_123", "dir_123")
if err != nil {
  // handle
}
_ = dir.Directory
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**directoryId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Directory().<a href="directory.go">GetPrimaryDirectoryByOrganizationId</a>(ctx, organizationId) -> (*GetDirectoryResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Convenience helper to return the first directory for an organization (if any).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
dir, err := client.Directory().GetPrimaryDirectoryByOrganizationId(ctx, "org_123")
if err != nil {
  // handle
}
_ = dir.Directory
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Directory().<a href="directory.go">ListDirectoryUsers</a>(ctx, organizationId, directoryId, options?) -> (*ListDirectoryUsersResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists users from a directory (paginated). Supports optional `UpdatedAfter`, group filtering, and detail inclusion.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
includeDetail := true
updatedAfter := time.Unix(1729851960, 0)

resp, err := client.Directory().ListDirectoryUsers(ctx, "org_123", "dir_123", &scalekit.ListDirectoryUsersOptions{
  PageSize:      10,
  PageToken:     "",
  IncludeDetail: &includeDetail,
  UpdatedAfter:  &updatedAfter,
})
if err != nil {
  // handle
}
_ = resp.Users
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**directoryId:** `string`

</dd>
</dl>

<dl>
<dd>

**options?:** `*ListDirectoryUsersOptions`
- `PageSize uint32`
- `PageToken string`
- `IncludeDetail *bool`
- `DirectoryGroupId *string`
- `UpdatedAfter *time.Time`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Directory().<a href="directory.go">ListDirectoryGroups</a>(ctx, organizationId, directoryId, options?) -> (*ListDirectoryGroupsResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists groups from a directory (paginated).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
includeDetail := true
resp, err := client.Directory().ListDirectoryGroups(ctx, "org_123", "dir_123", &scalekit.ListDirectoryGroupsOptions{
  PageSize:      10,
  PageToken:     "",
  IncludeDetail: &includeDetail,
})
if err != nil {
  // handle
}
_ = resp.Groups
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**directoryId:** `string`

</dd>
</dl>

<dl>
<dd>

**options?:** `*ListDirectoryGroupsOptions`
- `PageSize uint32`
- `PageToken string`
- `IncludeDetail *bool`
- `UpdatedAfter *time.Time`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Directory().<a href="directory.go">EnableDirectory</a>(ctx, organizationId, directoryId) -> (*ToggleDirectoryResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Enables a directory for an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.Directory().EnableDirectory(ctx, "org_123", "dir_123")
if err != nil {
  // handle
}
_ = resp.Enabled
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**directoryId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Directory().<a href="directory.go">DisableDirectory</a>(ctx, organizationId, directoryId) -> (*ToggleDirectoryResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Disables a directory for an organization.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.Directory().DisableDirectory(ctx, "org_123", "dir_123")
if err != nil {
  // handle
}
_ = resp.Enabled
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**organizationId:** `string`

</dd>
</dl>

<dl>
<dd>

**directoryId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## Sessions

<details><summary><code>client.Session().<a href="sessions.go">GetSession</a>(ctx, sessionId) -> (*SessionDetails, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches session details by session ID.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
session, err := client.Session().GetSession(ctx, "ses_123")
if err != nil {
  // handle
}
_ = session
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**sessionId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Session().<a href="sessions.go">GetUserSessions</a>(ctx, userId, pageSize, pageToken, filter?) -> (*UserSessionDetails, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists session details for a user (paginated).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.Session().GetUserSessions(ctx, "usr_123", 10, "", nil)
if err != nil {
  // handle
}
_ = resp.Sessions
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**userId:** `string`

</dd>
</dl>

<dl>
<dd>

**pageSize:** `uint32`

</dd>
</dl>

<dl>
<dd>

**pageToken:** `string`

</dd>
</dl>

<dl>
<dd>

**filter?:** `*sessionsv1.UserSessionFilter`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Session().<a href="sessions.go">RevokeSession</a>(ctx, sessionId) -> (*RevokeSessionResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Revokes a specific session by session ID.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.Session().RevokeSession(ctx, "ses_123")
if err != nil {
  // handle
}
_ = resp
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**sessionId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Session().<a href="sessions.go">RevokeAllUserSessions</a>(ctx, userId) -> (*RevokeAllUserSessionsResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Revokes all sessions for a user.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.Session().RevokeAllUserSessions(ctx, "usr_123")
if err != nil {
  // handle
}
_ = resp
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**userId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## Roles

<details><summary><code>client.Role().<a href="role.go">CreateRole</a>(ctx, role) -> (*CreateRoleResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a new environment-level role.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**role:** `*rolesv1.CreateRole`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">GetRole</a>(ctx, roleName) -> (*GetRoleResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches an environment-level role by name.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">ListRoles</a>(ctx) -> (*ListRolesResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists all environment-level roles.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">UpdateRole</a>(ctx, roleName, role) -> (*UpdateRoleResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates an environment-level role.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>

<dl>
<dd>

**role:** `*rolesv1.UpdateRole`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">DeleteRole</a>(ctx, roleName, reassignRoleName?) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes an environment-level role. Optionally provide `reassignRoleName` to reassign users.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>

<dl>
<dd>

**reassignRoleName?:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">GetRoleUsersCount</a>(ctx, roleName) -> (*GetRoleUsersCountResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Gets the count of users associated with an environment-level role.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">CreateOrganizationRole</a>(ctx, orgId, role) -> (*CreateOrganizationRoleResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates an organization-level role.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**orgId:** `string`

</dd>
</dl>

<dl>
<dd>

**role:** `*rolesv1.CreateOrganizationRole`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">GetOrganizationRole</a>(ctx, orgId, roleName) -> (*GetOrganizationRoleResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches an organization-level role by name.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**orgId:** `string`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">ListOrganizationRoles</a>(ctx, orgId) -> (*ListOrganizationRolesResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists organization-level roles.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**orgId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">UpdateOrganizationRole</a>(ctx, orgId, roleName, role) -> (*UpdateOrganizationRoleResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates an organization-level role by name.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**orgId:** `string`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>

<dl>
<dd>

**role:** `*rolesv1.UpdateRole`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">DeleteOrganizationRole</a>(ctx, orgId, roleName, reassignRoleName?) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes an organization-level role by name. Optionally provide `reassignRoleName`.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**orgId:** `string`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>

<dl>
<dd>

**reassignRoleName?:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">GetOrganizationRoleUsersCount</a>(ctx, orgId, roleName) -> (*GetOrganizationRoleUsersCountResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Gets the count of users associated with an organization-level role.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**orgId:** `string`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">UpdateDefaultOrganizationRoles</a>(ctx, orgId, defaultMemberRole) -> (*UpdateDefaultOrganizationRolesResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates the default member role for an organization.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**orgId:** `string`

</dd>
</dl>

<dl>
<dd>

**defaultMemberRole:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Role().<a href="role.go">DeleteOrganizationRoleBase</a>(ctx, orgId, roleName) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes the base relationship for an organization role.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**orgId:** `string`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## Permissions

<details><summary><code>client.Permission().<a href="permission.go">CreatePermission</a>(ctx, permission) -> (*CreatePermissionResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Creates a new permission.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**permission:** `*rolesv1.CreatePermission`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Permission().<a href="permission.go">GetPermission</a>(ctx, permissionName) -> (*GetPermissionResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Fetches a permission by name.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**permissionName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Permission().<a href="permission.go">ListPermissions</a>(ctx, pageToken?) -> (*ListPermissionsResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists permissions with optional pagination.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**pageToken?:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Permission().<a href="permission.go">UpdatePermission</a>(ctx, permissionName, permission) -> (*UpdatePermissionResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates an existing permission by name.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**permissionName:** `string`

</dd>
</dl>

<dl>
<dd>

**permission:** `*rolesv1.CreatePermission`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Permission().<a href="permission.go">DeletePermission</a>(ctx, permissionName) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes a permission by name.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**permissionName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Permission().<a href="permission.go">ListRolePermissions</a>(ctx, roleName) -> (*ListRolePermissionsResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists permissions associated with a role.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Permission().<a href="permission.go">AddPermissionsToRole</a>(ctx, roleName, permissionNames) -> (*AddPermissionsToRoleResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Adds permissions to a role.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>

<dl>
<dd>

**permissionNames:** `[]string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Permission().<a href="permission.go">RemovePermissionFromRole</a>(ctx, roleName, permissionName) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Removes a permission from a role.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>

<dl>
<dd>

**permissionName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Permission().<a href="permission.go">ListEffectiveRolePermissions</a>(ctx, roleName) -> (*ListEffectiveRolePermissionsResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists effective permissions for a role (including inherited permissions).
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**roleName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## Passwordless

<details><summary><code>client.Passwordless().<a href="passwordless.go">SendPasswordlessEmail</a>(ctx, email, options?) -> (*SendPasswordlessResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Sends a passwordless authentication email (OTP, magic link, etc. depending on configuration).
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
template := scalekit.TemplateTypeSignin
resp, err := client.Passwordless().SendPasswordlessEmail(ctx, "user@example.com", &scalekit.SendPasswordlessOptions{
  Template:         &template,
  MagiclinkAuthUri: "https://myapp.com/auth/callback",
  State:            "state",
  ExpiresIn:        1800,
  TemplateVariables: map[string]string{
    "app_name": "My App",
  },
})
if err != nil {
  // handle
}
_ = resp.AuthRequestId
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**email:** `string`

</dd>
</dl>

<dl>
<dd>

**options?:** `*SendPasswordlessOptions`
- `Template *TemplateType` (`SIGNIN`, `SIGNUP`)
- `MagiclinkAuthUri string`
- `State string`
- `ExpiresIn uint32`
- `TemplateVariables map[string]string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Passwordless().<a href="passwordless.go">VerifyPasswordlessEmail</a>(ctx, options) -> (*VerifyPasswordLessResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Verifies a passwordless authentication attempt using an OTP code or link token.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
verified, err := client.Passwordless().VerifyPasswordlessEmail(ctx, &scalekit.VerifyPasswordlessOptions{
  Code:          "123456",
  AuthRequestId: "auth_req_123",
})
if err != nil {
  // handle
}
_ = verified
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**options:** `*VerifyPasswordlessOptions`
- `Code string` - OTP code
- `LinkToken string` - Magic link token
- `AuthRequestId string` - Required in some flows

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.Passwordless().<a href="passwordless.go">ResendPasswordlessEmail</a>(ctx, authRequestId) -> (*SendPasswordlessResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Resends a passwordless authentication email for an existing auth request.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
resp, err := client.Passwordless().ResendPasswordlessEmail(ctx, "auth_req_123")
if err != nil {
  // handle
}
_ = resp.AuthRequestId
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**authRequestId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## WebAuthn

<details><summary><code>client.WebAuthn().<a href="webauthn.go">ListCredentials</a>(ctx, userId) -> (*ListCredentialsResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Lists passkey credentials for a user. If `userId` is empty, the API may list credentials for the current authenticated user.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**userId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.WebAuthn().<a href="webauthn.go">UpdateCredential</a>(ctx, credentialId, displayName) -> (*UpdateCredentialResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates the display name of a passkey credential.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**credentialId:** `string`

</dd>
</dl>

<dl>
<dd>

**displayName:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

<details><summary><code>client.WebAuthn().<a href="webauthn.go">DeleteCredential</a>(ctx, credentialId) -> (*DeleteCredentialResponse, error)</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Deletes a passkey credential by credential ID.
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**credentialId:** `string`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>

## Auth

<details><summary><code>client.Auth().<a href="auth_service.go">UpdateLoginUserDetails</a>(ctx, req) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Updates login user details associated with an authentication flow.

This method uses the Auth gRPC surface and expects a fully populated request.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
req := &scalekit.UpdateLoginUserDetailsRequest{
  ConnectionId:   "conn_123",
  LoginRequestId: "login_req_123",
  User: &scalekit.LoggedInUserDetails{
    Sub:   "sub",
    Email: "user@example.com",
  },
}

if err := client.Auth().UpdateLoginUserDetails(ctx, req); err != nil {
  // handle
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**ctx:** `context.Context`

</dd>
</dl>

<dl>
<dd>

**req:** `*UpdateLoginUserDetailsRequest`

</dd>
</dl>
</dd>
</dl>

</dd>
</dl>
</details>
