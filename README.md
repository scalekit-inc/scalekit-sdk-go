<p align="left">
  <a href="https://scalekit.com" target="_blank" rel="noopener noreferrer">
    <picture>
      <img src="https://cdn.scalekit.cloud/v1/scalekit-logo-dark.svg" height="64">
    </picture>
  </a>
  <br/>
</p>

# Official Go SDK

[![Go Reference](https://pkg.go.dev/badge/github.com/scalekit-inc/scalekit-sdk-go/v2.svg)](https://pkg.go.dev/github.com/scalekit-inc/scalekit-sdk-go/v2)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/scalekit-inc/scalekit-sdk-go/v2)](https://goreportcard.com/report/github.com/scalekit-inc/scalekit-sdk-go/v2)

<a href="https://scalekit.com" target="_blank" rel="noopener noreferrer">Scalekit</a> is the **auth stack for AI apps** - from human authentication to agent authorization. Build secure AI products faster with authentication for humans (SSO, passwordless, full-stack auth) and agents (MCP/APIs, delegated actions), all unified on one platform. This Go SDK enables both traditional B2B authentication and cutting-edge agentic workflows.

## 🤖 Agent-First Features

- **🔐 Agent Identity**: Agents as first-class actors with human ownership and org context
- **🎯 MCP-Native OAuth 2.1**: Purpose-built for Model Context Protocol with DCR/PKCE support
- **⏰ Ephemeral Credentials**: Time-bound, task-based authorization (minutes, not days)
- **🔒 Token Vault**: Per-user, per-tool token storage with rotation and progressive consent
- **👥 Human-in-the-Loop**: Step-up authentication when risk crosses thresholds
- **📊 Immutable Audit**: Track which user initiated, which agent acted, what resource was accessed

## 👨‍💼 Human Authentication

- **🔐 Enterprise SSO**: Support for SAML and OIDC protocols
- **👥 SCIM Provisioning**: Automated user provisioning and deprovisioning  
- **🚀 Passwordless Authentication**: Magic links, OTP, and modern auth flows
- **🏢 Multi-tenant Architecture**: Organization-level authentication policies
- **📱 Social Logins**: Support for popular social identity providers
- **🛡️ Full-Stack Auth**: Complete IdP-of-record solution for B2B SaaS

<div>
📚 <a target="_blank" href="https://docs.scalekit.com">Documentation</a> • 🚀 <a target="_blank" href="https://docs.scalekit.com/sso/quickstart/">SSO Quickstart</a> • 💻 <a target="_blank" href="https://docs.scalekit.com/apis">API Reference</a>
</div>
<hr />

## Pre-requisites

1. [Sign up](https://scalekit.com) for a Scalekit account.
2. Get your ```env_url```, ```client_id``` and ```client_secret``` from the Scalekit dashboard.

## Installation

```sh
go get -u github.com/scalekit-inc/scalekit-sdk-go/v2
```

## Usage

Initialize the Scalekit client using the appropriate credentials. Refer code sample below.

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
      State: "state",
      ConnectionId: "con_123456789",
    },
  )
}

```

##### Minimum Requirements

Before integrating the Scalekit Go SDK, ensure your development environment meets these requirements:

| Component | Version |
| --------- | ------- |
| Go        | 1.24+   |

> **Note:** Go 1.24+ provides the essential features required by this SDK. For optimal performance and security, consider using the latest stable release.


## Examples - SSO with Go HTTP Server

Below is a simple code sample that showcases how to implement Single Sign-on using Scalekit SDK

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

  // Get the authorization URL and redirect the user to the IdP login page
  http.HandleFunc("/auth/login", func(w http.ResponseWriter, r *http.Request) {
    authUrl, _ := scalekitClient.GetAuthorizationUrl(
      redirectUri,
      scalekit.AuthorizationUrlOptions{
        State: "state",
        ConnectionId: "con_123456789",
      },
    )
    http.Redirect(w, r, authUrl.String(), http.StatusSeeOther)
  })

  // Handle the callback from Scalekit
  http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    _ = r.URL.Query().Get("state") // validate state against the value set in /auth/login

    authResp, _ := scalekitClient.AuthenticateWithCode(r.Context(), code, redirectUri, scalekit.AuthenticationOptions{})

    http.SetCookie(w, &http.Cookie{
      Name: "access_token",
      Value: authResp.AccessToken,
    })

    fmt.Fprintf(w, "Access token: %s", authResp.AccessToken)
  })

  fmt.Println("Server started at http://localhost:8080")
  http.ListenAndServe(":8080", nil)
}
```

## 📱 Example Apps

Explore fully functional sample applications built with popular Go frameworks and the Scalekit SDK:

| Framework | Repository | Description |
|-----------|------------|-------------|
| **Go HTTP Server** | [scalekit-go-example](https://github.com/scalekit-developers/scalekit-go-example) | Basic HTTP server implementation |

## 🔗 Helpful Links

### 📖 Quickstart Guides
- [**SSO Integration**](https://docs.scalekit.com/sso/quickstart/) - Implement enterprise Single Sign-on
- [**Full Stack Auth**](https://docs.scalekit.com/fsa/quickstart/) - Complete authentication solution
- [**Passwordless Auth**](https://docs.scalekit.com/passwordless/quickstart/) - Modern authentication flows
- [**Social Logins**](https://docs.scalekit.com/social-logins/quickstart/) - Popular social identity providers
- [**Machine-to-Machine**](https://docs.scalekit.com/m2m/quickstart/) - API authentication

### 📚 Documentation & Reference
- [**API Reference**](https://docs.scalekit.com/apis) - Complete API documentation
- [**Developer Kit**](https://docs.scalekit.com/dev-kit/) - Tools and utilities
- [**API Authentication Guide**](https://docs.scalekit.com/guides/authenticate-scalekit-api/) - Secure API access

### 🛠️ Additional Resources
- [**Setup Guide**](https://docs.scalekit.com/guides/setup-scalekit/) - Initial platform configuration
- [**Code Examples**](https://docs.scalekit.com/directory/code-examples/) - Ready-to-use code snippets
- [**Admin Portal Guide**](https://docs.scalekit.com/directory/guides/admin-portal/) - Administrative interface
- [**Launch Checklist**](https://docs.scalekit.com/directory/guides/launch-checklist/) - Pre-production checklist

## License

This project is licensed under the **MIT license**.
See the [LICENSE](LICENSE) file for more information.
