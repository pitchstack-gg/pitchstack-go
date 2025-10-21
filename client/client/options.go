package client

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WithBaseURL overrides the default API base URL.
func WithBaseURL(baseURL string) ClientOpt {
	return func(c *Client) error {
		baseURL = strings.TrimSpace(baseURL)
		if baseURL == "" {
			return fmt.Errorf("baseURL must not be empty")
		}
		if _, err := url.ParseRequestURI(baseURL); err != nil {
			return fmt.Errorf("invalid baseURL %q: %w", baseURL, err)
		}
		c.baseURL = baseURL
		return nil
	}
}

// WithHTTPClient supplies a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOpt {
	return func(c *Client) error {
		if httpClient == nil {
			return fmt.Errorf("httpClient must not be nil")
		}
		c.httpClient = httpClient
		return nil
	}
}

// WithHTTPTimeout sets the timeout on the underlying HTTP client, creating a new client if needed.
func WithHTTPTimeout(timeout time.Duration) ClientOpt {
	return func(c *Client) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must be positive")
		}
		if c.httpClient == nil {
			c.httpClient = &http.Client{}
		}
		c.httpClient.Timeout = timeout
		return nil
	}
}

// WithCredentialProvider configures how the client obtains credentials for each request.
func WithCredentialProvider(provider CredentialProvider) ClientOpt {
	return func(c *Client) error {
		c.credentialProvider = provider
		return nil
	}
}

// WithStaticCredential registers a static credential that is reused for every request.
func WithStaticCredential(header, value string) ClientOpt {
	header = http.CanonicalHeaderKey(strings.TrimSpace(header))
	value = strings.TrimSpace(value)

	return WithCredentialProvider(CredentialProviderFunc(func(context.Context) (*Credential, error) {
		if header == "" || value == "" {
			return nil, nil
		}
		return &Credential{
			Header: header,
			Value:  value,
		}, nil
	}))
}

// WithStaticBearerToken configures bearer token authentication for the client.
func WithStaticBearerToken(token string) ClientOpt {
	token = strings.TrimSpace(token)
	if token == "" {
		return WithStaticCredential("", "")
	}
	return WithStaticCredential("Authorization", "Bearer "+token)
}

// WithDefaultHeader adds a header that is applied to every request.
func WithDefaultHeader(key, value string) ClientOpt {
	return func(c *Client) error {
		key = http.CanonicalHeaderKey(strings.TrimSpace(key))
		value = strings.TrimSpace(value)
		if key == "" {
			return fmt.Errorf("header key must not be empty")
		}
		if c.defaultHeaders == nil {
			c.defaultHeaders = make(http.Header)
		}
		c.defaultHeaders.Add(key, value)
		return nil
	}
}
