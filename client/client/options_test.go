package client

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type stubCredentialProvider struct {
	credential *Credential
}

func (s stubCredentialProvider) Credential(context.Context) (*Credential, error) {
	return s.credential, nil
}

func TestWithBaseURL(t *testing.T) {
	t.Run("when provided a valid url, then client base url updates", func(t *testing.T) {
		client := &Client{}
		err := WithBaseURL(" https://example.com ")(client)
		require.NoError(t, err)
		require.Equal(t, "https://example.com", client.baseURL)
	})

	t.Run("when provided an invalid url, then returns error", func(t *testing.T) {
		client := &Client{}
		err := WithBaseURL("://bad")(client)
		require.Error(t, err)
	})
}

func TestWithHTTPClient(t *testing.T) {
	t.Run("when provided a client, then stored on target", func(t *testing.T) {
		client := &Client{}
		httpClient := &http.Client{Timeout: time.Minute}
		err := WithHTTPClient(httpClient)(client)
		require.NoError(t, err)
		require.Equal(t, httpClient, client.httpClient)
	})

	t.Run("when provided nil, then returns error", func(t *testing.T) {
		client := &Client{}
		err := WithHTTPClient(nil)(client)
		require.Error(t, err)
	})
}

func TestWithHTTPTimeout(t *testing.T) {
	t.Run("when provided positive duration, then sets timeout", func(t *testing.T) {
		client := &Client{}
		err := WithHTTPTimeout(2 * time.Second)(client)
		require.NoError(t, err)
		require.NotNil(t, client.httpClient)
		require.Equal(t, 2*time.Second, client.httpClient.Timeout)
	})

	t.Run("when provided non positive duration, then returns error", func(t *testing.T) {
		client := &Client{}
		err := WithHTTPTimeout(0)(client)
		require.Error(t, err)
	})
}

func TestWithCredentialProvider(t *testing.T) {
	t.Run("when provided provider, then stored on client", func(t *testing.T) {
		client := &Client{}
		expected := &Credential{Header: "Authorization", Value: "Bearer token"}
		provider := stubCredentialProvider{credential: expected}

		err := WithCredentialProvider(provider)(client)
		require.NoError(t, err)
		require.NotNil(t, client.credentialProvider)

		actual, err := client.credentialProvider.Credential(context.Background())
		require.NoError(t, err)
		require.Equal(t, expected, actual)
	})
}

func TestWithCredential(t *testing.T) {
	t.Run("when provided header and value, then canonicalizes and stores", func(t *testing.T) {
		client := &Client{}
		err := WithStaticCredential(" x-auth ", " value ")(client)
		require.NoError(t, err)

		credential, err := client.credentialProvider.Credential(context.Background())
		require.NoError(t, err)
		require.NotNil(t, credential)
		require.Equal(t, "X-Auth", credential.Header)
		require.Equal(t, "value", credential.Value)
	})

	t.Run("when provided empty header, then provider yields nil", func(t *testing.T) {
		client := &Client{}
		err := WithStaticCredential(" ", "value")(client)
		require.NoError(t, err)

		credential, err := client.credentialProvider.Credential(context.Background())
		require.NoError(t, err)
		require.Nil(t, credential)
	})
}

func TestWithBearerToken(t *testing.T) {
	t.Run("when provided token, then stores authorization header", func(t *testing.T) {
		client := &Client{}
		err := WithStaticBearerToken(" token ")(client)
		require.NoError(t, err)

		credential, err := client.credentialProvider.Credential(context.Background())
		require.NoError(t, err)
		require.NotNil(t, credential)
		require.Equal(t, "Authorization", credential.Header)
		require.Equal(t, "Bearer token", credential.Value)
	})

	t.Run("when provided empty token, then provider yields nil", func(t *testing.T) {
		client := &Client{}
		err := WithStaticBearerToken(" ")(client)
		require.NoError(t, err)

		credential, err := client.credentialProvider.Credential(context.Background())
		require.NoError(t, err)
		require.Nil(t, credential)
	})
}

func TestWithDefaultHeader(t *testing.T) {
	t.Run("when provided header, then canonicalizes and stores", func(t *testing.T) {
		client := &Client{}
		err := WithDefaultHeader("x-test", " value ")(client)
		require.NoError(t, err)
		require.Equal(t, "value", client.defaultHeaders.Get("X-Test"))
	})

	t.Run("when provided empty key, then returns error", func(t *testing.T) {
		client := &Client{}
		err := WithDefaultHeader(" ", "value")(client)
		require.Error(t, err)
	})
}
