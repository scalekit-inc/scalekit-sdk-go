# Introduction

The Scalekit API is a RESTful API that enables you to manage organizations, users, and authentication settings. All requests must use HTTPS.
All API requests use the following base URLs:

```
https://{environment}.scalekit.dev (Development)
https://{environment}.scalekit.com (Production)
https://auth.example.com (Custom domain)
```

Scalekit operates two separate environments: Development and Production. Resources cannot be moved between environments.

# Authentication

The Scalekit API uses OAuth 2.0 Client Credentials for authentication.

Copy your API credentials from the Scalekit dashboard's API Config section and set them as environment variables.

```sh
SCALEKIT_ENVIRONMENT_URL='<YOUR_ENVIRONMENT_URL>'
SCALEKIT_CLIENT_ID='<ENVIRONMENT_CLIENT_ID>'
SCALEKIT_CLIENT_SECRET='<ENVIRONMENT_CLIENT_SECRET>'
```

Getting an access token

1. Get your credentials from the [Scalekit Dashboard](https://app.scalekit.com)
2. Request an access token:

```sh
curl https://{SCALEKIT_ENVIRONMENT_URL}/oauth/token \
  -X POST \
  -H 'Content-Type: application/x-www-form-urlencoded' \
  -d 'client_id={client_id}' \
  -d 'client_secret={client_secret}' \
  -d 'grant_type=client_credentials'
```

3. Use the access token in API requests:

```sh
curl https://{SCALEKIT_ENVIRONMENT_URL}/api/v1/organizations \
  -H 'Content-Type: application/json' \
  -H 'Authorization: Bearer {access_token}'
```

The response includes an access token:

```json
{
	"access_token": "eyJhbGciOiJSUzI1NiIsImtpZCI6InNua181Ok4OTEyMjU2NiIsInR5cCI6IkpXVCJ9...",
	"token_type": "Bearer",
	"expires_in": 86399,
	"scope": "openid"
}
```

# SDKs

Scalekit provides official SDKs for multiple programming languages. Check the changelog at GitHub repositories for the latest updates.

### Node.js

```sh
npm install @scalekit-sdk/node
```

Create a new Scalekit client instance after initializing the environment variables

```js
import { Scalekit } from "@scalekit-sdk/node";

export let scalekit = new Scalekit(
	process.env.SCALEKIT_ENVIRONMENT_URL,
	process.env.SCALEKIT_CLIENT_ID,
	process.env.SCALEKIT_CLIENT_SECRET
);
```

[See the Node SDK changelog](https://github.com/scalekit-inc/scalekit-sdk-node/releases)

### Python

```sh
pip install scalekit-sdk-python
```

Create a new Scalekit client instance after initializing the environment variables.

```py
from scalekit import ScalekitClient
import os

scalekit_client = ScalekitClient(
    os.environ.get('SCALEKIT_ENVIRONMENT_URL'),
    os.environ.get('SCALEKIT_CLIENT_ID'),
    os.environ.get('SCALEKIT_CLIENT_SECRET')
)
```

[See the Python SDK changelog](https://github.com/scalekit-inc/scalekit-sdk-python/releases)

### Go

```sh
go get -u github.com/scalekit-inc/scalekit-sdk-go
```

Create a new Scalekit client instance after initializing the environment variables.

```go
package main

import (
    "os"
    "github.com/scalekit-inc/scalekit-sdk-go"
)

scalekitClient := scalekit.NewScalekitClient(
    os.Getenv("SCALEKIT_ENVIRONMENT_URL"),
    os.Getenv("SCALEKIT_CLIENT_ID"),
    os.Getenv("SCALEKIT_CLIENT_SECRET"),
)
```

[See the Go SDK changelog](https://github.com/scalekit-inc/scalekit-sdk-go/releases)

### Java

```gradle
/* Gradle users - add the following to your dependencies in build file */
implementation "com.scalekit:scalekit-sdk-java:2.0.11"
```

```xml
<!-- Maven users - add the following to your `pom.xml` -->
<dependency>
    <groupId>com.scalekit</groupId>
    <artifactId>scalekit-sdk-java</artifactId>
    <version>2.0.11</version>
</dependency>
```

[See the Java SDK changelog](https://github.com/scalekit-inc/scalekit-sdk-java/releases)

# Error handling

The API uses standard HTTP status codes:

| Code        | Description          |
| ----------- | -------------------- |
| 200/201     | Success              |
| 400         | Invalid request      |
| 401         | Authentication error |
| 404         | Resource not found   |
| 429         | Rate limit exceeded  |
| 500/501/504 | Server error         |

Error responses include detailed information:

```json
{
	"code": 16,
	"message": "Token empty",
	"details": [
		{
			"@type": "type.googleapis.com/scalekit.v1.errdetails.ErrorInfo",
			"error_code": "UNAUTHENTICATED"
		}
	]
}
```
