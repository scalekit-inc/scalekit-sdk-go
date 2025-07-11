<p align="left">
  <a href="https://scalekit.com" target="_blank" rel="noopener noreferrer">
    <picture>
      <img src="https://cdn.scalekit.cloud/v1/scalekit-logo-dark.svg" height="64">
    </picture>
  </a>
  <br/>
</p>

# Official Go SDK

<a href="https://scalekit.com" target="_blank" rel="noopener noreferrer">Scalekit</a> is an Enterprise Authentication Platform purpose built for B2B applications. This Go SDK helps implement Enterprise Capabilities like Single Sign-on via SAML or OIDC in your Golang applications within a few hours.

<div>
📚 <a target="_blank" href="https://docs.scalekit.com">Documentation</a> - 🚀 <a target="_blank" href="https://docs.scalekit.com">Quick-start Guide</a> - 💻 <a target="_blank" href="https://docs.scalekit.com/apis">API Reference</a>
</div>
<hr />

## Pre-requisites

1. [Sign up](https://scalekit.com) for a Scalekit account.
2. Get your ```env_url```, ```client_id``` and ```client_secret``` from the Scalekit dashboard.

## Installation

```sh
go get -u github.com/scalekit-inc/scalekit-sdk-go
```

## Usage

Initialize the Scalekit client using the appropriate credentials. Refer code sample below.

```go
import "github.com/scalekit-inc/scalekit-sdk-go"

func main() {
  scalekitClient := scalekit.NewScalekitClient(
    "<SCALEKIT_ENV_URL>",
    "<SCALEKIT_CLIENT_ID>",
    "<SCALEKIT_CLIENT_SECRET>",
  )

  // Use the sc object to interact with the Scalekit API
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
| Go        | 1.22+   |

> **Note:** Go 1.22+ provides the essential features required by this SDK. For optimal performance and security, consider using the latest stable release.


## Examples - SSO with Go HTTP Server

Below is a simple code sample that showcases how to implement Single Sign-on using Scalekit SDK

```go
package main

import (
  "fmt"
  "net/http"

  "github.com/scalekit-inc/scalekit-sdk-go"
)

func main() {
  sc := scalekit.NewScalekit(
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
    http.Redirect(w, r, authUrl, http.StatusSeeOther)
  })

  // Handle the callback from the Scalekit
  http.HandleFunc("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    state := r.URL.Query().Get("state")

    authResp, _ := scalekitClient.AuthenticateWithCode(code, redirectUri)

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

## Example Apps

Fully functional sample applications written using some popular web application frameworks and Scalekit SDK. Feel free to clone the repo and run them locally

- [Go Http](https://github.com/scalekit-inc/scalekit-go-example.git)


## API Reference

Refer to our [API reference docs](https://docs.scalekit.com/apis) for detailed information about all our API endpoints and their usage.

## More Information

- Quickstart Guide to implement Single Sign-on in your application: [SSO Quickstart Guide](https://docs.scalekit.com)
- Understand Single Sign-on basics: [SSO Basics](https://docs.scalekit.com/best-practices/single-sign-on)

## License

This project is licensed under the **MIT license**.
See the [LICENSE](LICENSE) file for more information.
