<div align="center">

<a href="https://scalekit.com" target="_blank" rel="noopener noreferrer">
  <picture>
    <img src="../../scalekit-logo.svg" alt="Scalekit" height="64">
  </picture>
</a>

<p><strong>Official Go SDK for Scalekit — the auth stack for agents.</strong><br>
Authentication, authorization, and tool-calling for human-in-the-loop and autonomous agent flows.</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/scalekit-inc/scalekit-sdk-go/v2.svg)](https://pkg.go.dev/github.com/scalekit-inc/scalekit-sdk-go/v2)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/scalekit-inc/scalekit-sdk-go/v2)](https://goreportcard.com/report/github.com/scalekit-inc/scalekit-sdk-go/v2)

**[📖 Documentation](https://docs.scalekit.com)** · **[🚀 Quickstart](https://docs.scalekit.com/sso/quickstart/)** · **[💻 API Reference](https://docs.scalekit.com/apis)** · **[💬 Slack](https://join.slack.com/t/scalekit-community/shared_invite/zt-3gsxwr4hc-0tvhwT2b_qgVSIZQBQCWRw)**

</div>

---

This is the official Go SDK for [Scalekit](https://scalekit.com) — a complete auth stack for agents. Build secure AI products faster with authentication for humans (SSO, SCIM, passwordless, full-stack auth) and agents (MCP/APIs, delegated actions, tool-calling), all unified on one platform.

---

### Features

#### Agent-First

- **Agent Identity** — Agents as first-class actors with human ownership and org context
- **MCP-Native OAuth 2.1** — Purpose-built for Model Context Protocol with DCR/PKCE support
- **Ephemeral Credentials** — Time-bound, task-based authorization (minutes, not days)
- **Token Vault** — Per-user, per-tool token storage with rotation and progressive consent
- **Human-in-the-Loop** — Step-up authentication when risk crosses thresholds
- **Immutable Audit** — Track which user initiated, which agent acted, what resource was accessed

#### Human Authentication

- **Enterprise SSO** — Support for SAML and OIDC protocols
- **SCIM Provisioning** — Automated user provisioning and deprovisioning
- **Passwordless Authentication** — Magic links, OTP, and modern auth flows
- **Multi-tenant Architecture** — Organization-level authentication policies
- **Social Logins** — Support for popular social identity providers
- **Full-Stack Auth** — Complete IdP-of-record solution for B2B SaaS

---

### Getting started

#### Prerequisites

- **Go** ≥ 1.25.8 (required for security fixes in stdlib)
- [Scalekit account](https://scalekit.com) with `env_url`, `client_id`, and `client_secret`

#### Installation

```sh
go get -u github.com/scalekit-inc/scalekit-sdk-go/v2
```

#### Usage

```go
import "github.com/scalekit-inc/scalekit-sdk-go/v2"

func main() {
    scalekitClient := scalekit.NewScalekitClient(
        "<SCALEKIT_ENV_URL>",
        "<SCALEKIT_CLIENT_ID>",
        "<SCALEKIT_CLIENT_SECRET>",
    )

    // Use scalekitClient to interact with the Scalekit API
    authUrl, _ := scalekitClient.GetAuthorizationUrl(
        "https://acme-corp.com/redirect-uri",
        scalekit.AuthorizationUrlOptions{
            State:        "state",
            ConnectionId: "con_123456789",
        },
    )
}
```

---

### Example — SSO with Go HTTP Server

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/scalekit-inc/scalekit-sdk-go/v2"
)

func main() {
    scalekitClient := scalekit.NewScalekitClient(
        "<SCALEKIT_ENV_URL>",
        "<SCALEKIT_CLIENT_ID>",
        "<SCALEKIT_CLIENT_SECRET>",
    )

    redirectUri := "http://localhost:8080/auth/callback"

    http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
        authUrl, _ := scalekitClient.GetAuthorizationUrl(
            redirectUri,
            scalekit.AuthorizationUrlOptions{
                State:        "state",
                ConnectionId: "con_123456789",
            },
        )
        http.Redirect(w, r, authUrl.String(), http.StatusSeeOther)
    })

    http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
        code := r.URL.Query().Get("code")

        authResp, _ := scalekitClient.AuthenticateWithCode(r.Context(), code, redirectUri, scalekit.AuthenticationOptions{})

        http.SetCookie(w, &http.Cookie{
            Name:  "access_token",
            Value: authResp.AccessToken,
        })

        fmt.Fprintf(w, "Access token: %s", authResp.AccessToken)
    })

    fmt.Println("Server started at http://localhost:8080")
    http.ListenAndServe(":8080", nil)
}
```

---

### Example Apps

| Framework | Repository | Description |
|-----------|------------|-------------|
| **Go HTTP Server** | [scalekit-go-example](https://github.com/scalekit-developers/scalekit-go-example) | Basic HTTP server implementation |
| **Gin** | [scalekit-gin-example](https://github.com/scalekit-developers/scalekit-go-example) | Gin framework integration |

---

### Helpful Links

#### Quickstart Guides

- [SSO Integration](https://docs.scalekit.com/sso/quickstart/) — Implement enterprise Single Sign-on
- [Full Stack Auth](https://docs.scalekit.com/fsa/quickstart/) — Complete authentication solution
- [Passwordless Auth](https://docs.scalekit.com/passwordless/quickstart/) — Modern authentication flows
- [Social Logins](https://docs.scalekit.com/social-logins/quickstart/) — Popular social identity providers
- [Machine-to-Machine](https://docs.scalekit.com/m2m/quickstart/) — API authentication

#### Documentation & Reference

- [API Reference](https://docs.scalekit.com/apis) — Complete API documentation
- [Developer Kit](https://docs.scalekit.com/dev-kit/) — Tools and utilities
- [API Authentication Guide](https://docs.scalekit.com/guides/authenticate-scalekit-api/) — Secure API access

#### Additional Resources

- [Setup Guide](https://docs.scalekit.com/guides/setup-scalekit/) — Initial platform configuration
- [Code Examples](https://docs.scalekit.com/directory/code-examples/) — Ready-to-use code snippets
- [Admin Portal Guide](https://docs.scalekit.com/directory/guides/admin-portal/) — Administrative interface
- [Launch Checklist](https://docs.scalekit.com/directory/guides/launch-checklist/) — Pre-production checklist

---

### Contributing

Contributions are welcome! Coming soon: contribution guidelines.

For now:
1. Fork this repository
2. Create a branch — `git checkout -b fix/my-improvement`
3. Make your changes
4. Run tests — `make test`
5. Open a Pull Request

---

### License

This project is licensed under the **MIT license**. See the [LICENSE](LICENSE) file for more information.
