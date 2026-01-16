package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListInbox(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/notifications", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "true", query.Get("unreadOnly"))
		require.Equal(t, "true", query.Get("includeArchived"))
		require.Equal(t, "false", query.Get("includeExpired"))
		require.Equal(t, []string{"system", "social"}, query["categories"])
		require.Equal(t, "25", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(ListInboxResponse{
			Notifications: []Notification{{MessageID: "msg-1"}},
			NextToken:     "next",
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	unreadOnly := true
	includeArchived := true
	includeExpired := false
	pageSize := int32(25)
	resp, err := client.ListInbox(context.Background(), &ListInboxRequest{
		UnreadOnly:      &unreadOnly,
		IncludeArchived: &includeArchived,
		IncludeExpired:  &includeExpired,
		Categories:      []string{" system ", "social"},
		PageSize:        &pageSize,
		NextToken:       "token",
	})
	require.NoError(t, err)
	require.Equal(t, "msg-1", resp.Notifications[0].MessageID)
	require.Equal(t, "next", resp.NextToken)
}

func TestClientGetMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/notifications/msg-1", r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(GetMessageResponse{
			Notification: &Notification{MessageID: "msg-1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetMessage(context.Background(), &GetMessageRequest{MessageID: "msg-1"})
	require.NoError(t, err)
	require.Equal(t, "msg-1", resp.Notification.MessageID)
}

func TestClientDeleteMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/notifications/msg-1", r.URL.Path)
		w.Header().Set("X-Request-Id", "req-delete")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.DeleteMessage(context.Background(), &DeleteMessageRequest{MessageID: "msg-1"})
	require.NoError(t, err)
	require.Equal(t, "req-delete", resp.Metadata.RequestID)
}

func TestClientMarkRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/notifications/msg-1:markRead", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var payload map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Len(t, payload, 0)
		w.Header().Set("X-Request-Id", "req-mark")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.MarkRead(context.Background(), &MarkReadRequest{MessageID: "msg-1"})
	require.NoError(t, err)
	require.Equal(t, "req-mark", resp.Metadata.RequestID)
}

func TestClientArchiveMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/notifications/msg-1:archive", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var payload map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Len(t, payload, 0)
		w.Header().Set("X-Request-Id", "req-archive")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ArchiveMessage(context.Background(), &ArchiveMessageRequest{MessageID: "msg-1"})
	require.NoError(t, err)
	require.Equal(t, "req-archive", resp.Metadata.RequestID)
}
