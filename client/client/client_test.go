package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type metadataAwareResponse struct {
	Message  string           `json:"message"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *metadataAwareResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

var _ metadataSetter = (*metadataAwareResponse)(nil)

func TestNewClient(t *testing.T) {
	t.Run("when provided no options, then returns client with defaults", func(t *testing.T) {
		client, err := NewClient()
		require.NoError(t, err)
		require.NotNil(t, client)
		require.Equal(t, defaultBaseURL, client.baseURL)
		require.NotNil(t, client.httpClient)
		require.Equal(t, 30*time.Second, client.httpClient.Timeout)
		require.NotNil(t, client.defaultHeaders)
	})

	t.Run("when option returns error, then construction fails", func(t *testing.T) {
		client, err := NewClient(func(c *Client) error {
			return errors.New("boom")
		})
		require.Error(t, err)
		require.Nil(t, client)
	})
}

func TestClientNewRequest(t *testing.T) {
	client := newTestClient(t, WithBaseURL("https://example.com/api"))

	t.Run("when context is nil, then returns error", func(t *testing.T) {
		var nilCtx context.Context
		req, err := client.newRequest(nilCtx, http.MethodGet, "/collections", nil)
		require.Error(t, err)
		require.Nil(t, req)
	})

	t.Run("when path lacks leading slash, then prefixes automatically", func(t *testing.T) {
		req, err := client.newRequest(context.Background(), http.MethodGet, "collections", nil)
		require.NoError(t, err)
		require.Equal(t, "https://example.com/collections", req.URL.String())
	})
}

func TestClientDo(t *testing.T) {
	t.Run("when response succeeds, then decodes body and captures metadata", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "application/json", r.Header.Get("Accept"))
			require.Equal(t, "value", r.Header.Get("X-Test"))
			require.Equal(t, "Bearer token", r.Header.Get("Authorization"))
			require.Equal(t, "opt", r.Header.Get("X-Opt"))

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Request-Id", "req-123")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"ok"}`))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(
			t,
			WithBaseURL(server.URL),
			WithHTTPClient(server.Client()),
			WithDefaultHeader("X-Test", "value"),
			WithStaticBearerToken("token"),
		)

		req, err := client.newRequest(context.Background(), http.MethodGet, "/path", nil)
		require.NoError(t, err)

		resp := &metadataAwareResponse{}
		err = client.do(req, resp, WithRequestHeader("X-Opt", "opt"))
		require.NoError(t, err)
		require.Equal(t, "ok", resp.Message)
		require.Equal(t, 200, resp.Metadata.StatusCode)
		require.Equal(t, "req-123", resp.Metadata.RequestID)
		require.Equal(t, "application/json", resp.Metadata.Header.Get("Content-Type"))
	})

	t.Run("when request option fails, then execution fails", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

		req, err := client.newRequest(context.Background(), http.MethodGet, "/path", nil)
		require.NoError(t, err)

		optErr := errors.New("opt failure")
		err = client.do(req, nil, func(*http.Request) error { return optErr })
		require.ErrorIs(t, err, optErr)
	})

	t.Run("when server returns non success, then api error surfaces", func(t *testing.T) {
		handler := func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"code":7,"message":"bad request"}`))
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

		req, err := client.newRequest(context.Background(), http.MethodGet, "/path", nil)
		require.NoError(t, err)

		err = client.do(req, nil)
		require.Error(t, err)

		var apiErr *APIError
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, 7, apiErr.Status.Code)
		require.Equal(t, "bad request", apiErr.Status.Message)
	})

	t.Run("when response body is invalid json, then decode fails", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("not json"))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

		req, err := client.newRequest(context.Background(), http.MethodGet, "/path", nil)
		require.NoError(t, err)

		var out struct{}
		err = client.do(req, &out)
		require.Error(t, err)
	})

	t.Run("when using credential provider, then credentials are fetched per request", func(t *testing.T) {
		var providerCalls int32
		var requestCount int32

		handler := func(w http.ResponseWriter, r *http.Request) {
			currentRequest := atomic.AddInt32(&requestCount, 1)
			expected := fmt.Sprintf("Bearer token-%d", currentRequest)
			require.Equal(t, expected, r.Header.Get("Authorization"))
			w.WriteHeader(http.StatusOK)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		provider := CredentialProviderFunc(func(context.Context) (*Credential, error) {
			call := atomic.AddInt32(&providerCalls, 1)
			return &Credential{
				Header: "Authorization",
				Value:  fmt.Sprintf("Bearer token-%d", call),
			}, nil
		})

		client := newTestClient(
			t,
			WithBaseURL(server.URL),
			WithHTTPClient(server.Client()),
			WithCredentialProvider(provider),
		)

		for i := 0; i < 2; i++ {
			req, err := client.newRequest(context.Background(), http.MethodGet, "/path", nil)
			require.NoError(t, err)

			err = client.do(req, nil)
			require.NoError(t, err)
		}

		require.Equal(t, int32(2), atomic.LoadInt32(&providerCalls))
		require.Equal(t, int32(2), atomic.LoadInt32(&requestCount))
	})
}

func newTestClient(t *testing.T, opts ...ClientOpt) *Client {
	t.Helper()

	client, err := NewClient(opts...)
	require.NoError(t, err)
	require.NotNil(t, client)

	return client
}
