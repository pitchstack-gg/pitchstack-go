package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListMyGroupInvites(t *testing.T) {
	t.Run("when given paging, then lists the caller's invites with group summaries", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/groups:myInvites", r.URL.Path)
			require.Equal(t, "20", r.URL.Query().Get("pageSize"))
			require.Equal(t, "tok", r.URL.Query().Get("nextToken"))
			_, _ = w.Write([]byte(`{"invites":[{"invite":{"inviteId":"inv-1","roleValue":"GROUP_ROLE_MEMBER"},"group":{"groupId":"group-1","slug":"g","name":"G","avatarUrl":"https://cdn/a.png"}}],"nextToken":"next"}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		pageSize := int32(20)
		resp, err := client.ListMyGroupInvites(context.Background(), &ListMyGroupInvitesRequest{PageSize: &pageSize, NextToken: "tok"})
		require.NoError(t, err)
		require.Equal(t, "next", resp.NextToken)
		require.Equal(t, "inv-1", resp.Invites[0].Invite.InviteID)
		require.Equal(t, "g", resp.Invites[0].Group.Slug)
		require.Equal(t, "https://cdn/a.png", resp.Invites[0].Group.AvatarURL)
	})

	t.Run("when request is nil, then lists without query parameters", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v1/groups:myInvites", r.URL.Path)
			require.Empty(t, r.URL.RawQuery)
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		_, err := client.ListMyGroupInvites(context.Background(), nil)
		require.NoError(t, err)
	})
}

func TestClientListMyGroupJoinRequests(t *testing.T) {
	t.Run("when called, then lists the caller's join requests with decline reasons", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/groups:myJoinRequests", r.URL.Path)
			require.Equal(t, "5", r.URL.Query().Get("pageSize"))
			_, _ = w.Write([]byte(`{"joinRequests":[{"joinRequest":{"requestId":"jr-1","status":"GROUP_JOIN_REQUEST_STATUS_DECLINED","declineReason":"full"},"group":{"groupId":"group-1"}}]}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		pageSize := int32(5)
		resp, err := client.ListMyGroupJoinRequests(context.Background(), &ListMyGroupJoinRequestsRequest{PageSize: &pageSize})
		require.NoError(t, err)
		require.Len(t, resp.JoinRequests, 1)
		require.Equal(t, GroupJoinRequestStatusDeclined, resp.JoinRequests[0].JoinRequest.Status)
		require.Equal(t, "full", resp.JoinRequests[0].JoinRequest.DeclineReason)
		require.Equal(t, "group-1", resp.JoinRequests[0].Group.GroupID)
	})
}

func TestClientAcceptMyGroupInvite(t *testing.T) {
	t.Run("when given group and invite ids, then posts to the accept path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/groups/group-1/invites/inv-1:accept", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			w.Header().Set("X-Request-Id", "req-accept")
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.AcceptMyGroupInvite(context.Background(), &MyGroupInviteRequest{GroupID: "group-1", InviteID: "inv-1"})
		require.NoError(t, err)
		require.Equal(t, "req-accept", resp.Metadata.RequestID)
	})

	t.Run("when the invite id is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.AcceptMyGroupInvite(context.Background(), &MyGroupInviteRequest{GroupID: "group-1"})
		require.Error(t, err)
		_, err = client.AcceptMyGroupInvite(context.Background(), nil)
		require.Error(t, err)
	})
}

func TestClientDeclineMyGroupInvite(t *testing.T) {
	t.Run("when given group and invite ids, then posts to the decline path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/groups/group-1/invites/inv-1:decline", r.URL.Path)
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		_, err := client.DeclineMyGroupInvite(context.Background(), &MyGroupInviteRequest{GroupID: "group-1", InviteID: "inv-1"})
		require.NoError(t, err)
	})

	t.Run("when the group id is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.DeclineMyGroupInvite(context.Background(), &MyGroupInviteRequest{InviteID: "inv-1"})
		require.Error(t, err)
	})
}
