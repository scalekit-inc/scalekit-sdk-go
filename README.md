<p align="center">
  <a href="https://scalekit.com" target="_blank" rel="noopener noreferrer">
    <picture>
      <img src="https://cdn.scalekit.cloud/v1/scalekit-logo-dark.svg" height="64">
    </picture>
  </a>
  <br/>
</p>
<h1 align="center">
  Official Scalekit Go SDK
</h1>

Scalekit helps you in shipping enterprise auth in days.

This Go SDK is a wrapper around Scalekit's REST API to help you integrate Scalekit with your Go applications.

## Getting Started

1. [Sign up](https://scalekit.com) for a Scalekit account.
2. Get your ```env_url```, ```client_id``` and ```client_secret``` from the Scalekit dashboard.

## Installation

```sh
go get -u github.com/scalekit-inc/scalekit-sdk-go
```

## Usage

```go
import "github.com/scalekit-inc/scalekit-sdk-go"

func main() {
  sc := scalekit.NewScalekit(
    "<SCALEKIT_ENV_URL>",
    "<SCALEKIT_CLIENT_ID>",
    "<SCALEKIT_CLIENT_SECRET>",
  )

  // Use the sc object to interact with the Scalekit API
  authUrl, _ := sc.GetAuthorizationUrl("https://acme-corp.com/redirect-uri", AuthorizationUrlOptions{
		State: "state",
		ConnectionId:   "connection_id",
	})
}

```

## API Reference
See the [Scalekit API docs](https://docs.scalekit.com) for more information about the API and authentication.

## License
This project is licensed under the **MIT license**.
See the [LICENSE](LICENSE) file for more information.