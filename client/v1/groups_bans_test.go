package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientBanGroupUser(t *testing.T) {
	t.Run("when given a group and user, then posts the ban and decodes it", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/groups/group-1/bans", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user-1", payload["userId"])
			require.Equal(t, "spam", payload["reason"])
			w.Header().Set("X-Request-Id", "req-ban")
			_, _ = w.Write([]byte(`{"ban":{"groupId":"group-1","userId":"user-1","bannedByUserId":"admin-1","reason":"spam","createdAt":"2026-01-02T03:04:05Z"}}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BanGroupUser(context.Background(), &BanGroupUserRequest{GroupID: "group-1", UserID: " user-1 ", Reason: "spam"})
		require.NoError(t, err)
		require.Equal(t, "user-1", resp.Ban.UserID)
		require.Equal(t, "admin-1", resp.Ban.BannedByUserID)
		require.NotNil(t, resp.Ban.CreatedAt)
		require.Equal(t, "req-ban", resp.Metadata.RequestID)
	})

	t.Run("when the user is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.BanGroupUser(context.Background(), &BanGroupUserRequest{GroupID: "group-1"})
		require.Error(t, err)
		_, err = client.BanGroupUser(context.Background(), nil)
		require.Error(t, err)
	})
}

func TestClientUnbanGroupUser(t *testing.T) {
	t.Run("when given a group and user, then deletes the ban", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/v1/groups/group-1/bans/user-1", r.URL.Path)
			w.Header().Set("X-Request-Id", "req-unban")
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.UnbanGroupUser(context.Background(), &UnbanGroupUserRequest{GroupID: "group-1", UserID: "user-1"})
		require.NoError(t, err)
		require.Equal(t, "req-unban", resp.Metadata.RequestID)
	})

	t.Run("when the group is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.UnbanGroupUser(context.Background(), &UnbanGroupUserRequest{UserID: "user-1"})
		require.Error(t, err)
	})
}

func TestClientListGroupBans(t *testing.T) {
	t.Run("when given paging, then lists bans with query parameters", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/groups/group-1/bans", r.URL.Path)
			require.Equal(t, "10", r.URL.Query().Get("pageSize"))
			require.Equal(t, "tok", r.URL.Query().Get("nextToken"))
			require.NoError(t, json.NewEncoder(w).Encode(ListGroupBansResponse{
				Bans:      []GroupBan{{UserID: "user-1"}},
				NextToken: "next",
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		pageSize := int32(10)
		resp, err := client.ListGroupBans(context.Background(), &ListGroupBansRequest{GroupID: "group-1", PageSize: &pageSize, NextToken: "tok"})
		require.NoError(t, err)
		require.Len(t, resp.Bans, 1)
		require.Equal(t, "next", resp.NextToken)
	})

	t.Run("when the group is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ListGroupBans(context.Background(), &ListGroupBansRequest{})
		require.Error(t, err)
	})
}
