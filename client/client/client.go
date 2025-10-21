package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.pitchstack.gg"

// Client is a lightweight REST client for the Pitchstack v1 API.
type Client struct {
	baseURL            string
	httpClient         *http.Client
	defaultHeaders     http.Header
	credentialProvider CredentialProvider
}

// ClientOpt configures a Client during construction.
type ClientOpt func(*Client) error

// NewClient constructs a Client that targets the Pitchstack API.
// Callers can customize behavior using ClientOpt modifiers.
func NewClient(opts ...ClientOpt) (*Client, error) {
	c := &Client{
		baseURL: defaultBaseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		defaultHeaders: make(http.Header),
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("apply client option: %w", err)
		}
	}

	if c.httpClient == nil {
		return nil, errors.New("client httpClient must not be nil")
	}
	if c.baseURL == "" {
		return nil, errors.New("client baseURL must not be empty")
	}

	return c, nil
}

// newRequest creates a request bound to the configured base URL.
func (c *Client) newRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	if ctx == nil {
		return nil, errors.New("context must not be nil")
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	base, err := url.Parse(c.baseURL)
	if err != nil {
		return nil, fmt.Errorf("parse base url: %w", err)
	}
	rel, err := url.Parse(path)
	if err != nil {
		return nil, fmt.Errorf("parse request path: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, base.ResolveReference(rel).String(), body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	return req, nil
}

// do executes an HTTP request and decodes the response into out.
func (c *Client) do(req *http.Request, out any, opts ...RequestOpt) error {
	if req == nil {
		return errors.New("request must not be nil")
	}

	copyHeaders(req, c.defaultHeaders)

	if req.Header.Get("Accept") == "" {
		req.Header.Set("Accept", "application/json")
	}

	if err := c.applyCredentials(req); err != nil {
		return fmt.Errorf("apply credentials: %w", err)
	}

	for _, opt := range opts {
		if opt == nil {
			continue
		}
		if err := opt(req); err != nil {
			return fmt.Errorf("apply request option: %w", err)
		}
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("execute http request: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response body: %w", err)
	}

	metadata := ResponseMetadata{
		StatusCode: resp.StatusCode,
		Header:     cloneHeader(resp.Header),
		RequestID:  resp.Header.Get("X-Request-Id"),
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return newAPIError(metadata, body)
	}

	if out != nil && len(body) > 0 {
		if err := json.Unmarshal(body, out); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}

	if setter, ok := out.(metadataSetter); ok && setter != nil {
		setter.setMetadata(metadata)
	}

	return nil
}

func (c *Client) applyCredentials(req *http.Request) error {
	if c.credentialProvider == nil {
		return nil
	}

	credential, err := c.credentialProvider.Credential(req.Context())
	if err != nil {
		return err
	}
	if credential == nil {
		return nil
	}

	header := http.CanonicalHeaderKey(strings.TrimSpace(credential.Header))
	value := strings.TrimSpace(credential.Value)
	if header == "" || value == "" {
		return nil
	}

	if req.Header.Get(header) == "" {
		req.Header.Set(header, value)
	}

	return nil
}

func copyHeaders(dst *http.Request, headers http.Header) {
	for key, values := range headers {
		for _, value := range values {
			dst.Header.Add(key, value)
		}
	}
}

func cloneHeader(h http.Header) http.Header {
	out := make(http.Header, len(h))
	for key, values := range h {
		cp := make([]string, len(values))
		copy(cp, values)
		out[key] = cp
	}
	return out
}

func jsonBody(payload any) (io.Reader, error) {
	if payload == nil {
		return nil, nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(data), nil
}
