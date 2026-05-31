# Pitchstack Go

Go client for the Pitchstack v1 API.

The client is intentionally small: it wraps the public HTTP API, handles JSON encoding and decoding, applies request credentials, and exposes typed request and response structs for the v1 resources.

## Install

```sh
go get github.com/pitchstack-gg/pitchstack-go@v0.1.0
```

Import the v1 client package:

```go
import client "github.com/pitchstack-gg/pitchstack-go/client/v1"
```

## Usage

By default the client targets `https://api.pitchstack.gg`.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	client "github.com/pitchstack-gg/pitchstack-go/client/v1"
)

func main() {
	c, err := client.NewClient(
		client.WithStaticBearerToken("YOUR_API_TOKEN"),
		client.WithHTTPTimeout(20*time.Second),
	)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pageSize := int32(10)

	resp, err := c.ListCards(ctx, &client.ListCardsRequest{
		PageSize: &pageSize,
	})
	if err != nil {
		log.Fatal(err)
	}

	for _, card := range resp.Cards {
		fmt.Println(card.Name)
	}
}
```

For non-production environments, override the API endpoint:

```go
c, err := client.NewClient(
	client.WithBaseURL("https://api-gamma.pitchstack.gg:8443"),
	client.WithStaticBearerToken("YOUR_API_TOKEN"),
)
```

## Authentication

Most authenticated endpoints expect a bearer token. Use `WithStaticBearerToken` for a fixed token or `WithCredentialProvider` when credentials need to be loaded or refreshed per request.

```go
c, err := client.NewClient(
	client.WithCredentialProvider(client.CredentialProviderFunc(func(ctx context.Context) (*client.Credential, error) {
		token, err := loadToken(ctx)
		if err != nil {
			return nil, err
		}

		return &client.Credential{
			Header: "Authorization",
			Value:  "Bearer " + token,
		}, nil
	})),
)
```

## Errors

Non-2xx API responses are returned as `*client.APIError`. The error includes HTTP response metadata and the structured API status payload when available.

```go
resp, err := c.GetCard(ctx, &client.GetCardRequest{Identifier: "missing-card"})
if err != nil {
	var apiErr *client.APIError
	if errors.As(err, &apiErr) {
		log.Printf("request failed: status=%d request_id=%s message=%s", apiErr.Metadata.StatusCode, apiErr.Metadata.RequestID, apiErr.Status.Message)
		return
	}
	log.Fatal(err)
}
_ = resp
```

## Development

```sh
make test-unit
make lint
make gosec
```

Integration tests require API credentials and are skipped when `PITCHSTACK_API_TOKEN` is not set:

```sh
PITCHSTACK_API_BASE_URL=https://api-gamma.pitchstack.gg:8443 \
PITCHSTACK_API_TOKEN=... \
make test-integration
```

## Releases

Public releases are tagged with semantic versions. The first public client release is `v0.1.0`; while the module remains pre-1.0, minor versions may include API-shaping changes and patch versions should be bug fixes only.

## License

See [LICENSE](LICENSE).
