package client

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientInitiatePatreonLink(t *testing.T) {
	t.Run("when called without a request, then posts an empty body and decodes the URL", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/auth/patreon/link:initiate", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)
			require.JSONEq(t, `{}`, string(body))
			_, _ = w.Write([]byte(`{"authorizationUrl":"https://patreon/auth","state":"st"}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.InitiatePatreonLink(context.Background(), nil)
		require.NoError(t, err)
		require.Equal(t, "https://patreon/auth", resp.AuthorizationURL)
		require.Equal(t, "st", resp.State)
	})
}

func TestClientCompletePatreonLink(t *testing.T) {
	t.Run("when given a code and state, then returns the link and access profile", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/auth/patreon/link:complete", r.URL.Path)
			var payload CompletePatreonLinkRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "code", payload.Code)
			require.Equal(t, "st", payload.State)
			_, _ = w.Write([]byte(`{"link":{"patreonUserId":"p-1","tierIds":["t-1"]},"accessProfile":{"userId":"user-1","tier":"supporter","version":"3"}}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.CompletePatreonLink(context.Background(), &CompletePatreonLinkRequest{Code: " code ", State: "st"})
		require.NoError(t, err)
		require.Equal(t, "p-1", resp.Link.PatreonUserID)
		require.Equal(t, []string{"t-1"}, resp.Link.TierIDs)
		require.Equal(t, "supporter", resp.AccessProfile.Tier)
		require.Equal(t, int64(3), resp.AccessProfile.Version)
	})

	t.Run("when the code is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.CompletePatreonLink(context.Background(), &CompletePatreonLinkRequest{State: "st"})
		require.Error(t, err)
		_, err = client.CompletePatreonLink(context.Background(), nil)
		require.Error(t, err)
	})
}

func TestClientUnlinkPatreon(t *testing.T) {
	t.Run("when called, then deletes the Patreon link", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/v1/auth/patreon/link", r.URL.Path)
			w.Header().Set("X-Request-Id", "req-unlink")
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.UnlinkPatreon(context.Background(), nil)
		require.NoError(t, err)
		require.Equal(t, "req-unlink", resp.Metadata.RequestID)
	})
}
