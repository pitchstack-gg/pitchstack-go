package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientCreateGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/groups", r.URL.Path)
		var payload CreateGroupRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "Testers", payload.Name)
		require.Equal(t, VisibilityLevelPrivate, payload.Visibility)

		require.NoError(t, json.NewEncoder(w).Encode(CreateGroupResponse{
			Group: &UserGroup{GroupID: "group-1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.CreateGroup(context.Background(), &CreateGroupRequest{
		Name:       "Testers",
		Visibility: VisibilityLevelPrivate,
	})
	require.NoError(t, err)
	require.Equal(t, "group-1", resp.Group.GroupID)
}

func TestClientUpdateGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/v1/groups/group-1", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload struct {
			Name string `json:"name"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "Renamed", payload.Name)

		require.NoError(t, json.NewEncoder(w).Encode(UpdateGroupResponse{
			Group: &UserGroup{GroupID: "group-1", Name: "Renamed"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	name := "Renamed"
	resp, err := client.UpdateGroup(context.Background(), &UpdateGroupRequest{
		GroupID: "group-1",
		Name:    &name,
	})
	require.NoError(t, err)
	require.Equal(t, "Renamed", resp.Group.Name)
}

func TestClientDeleteGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/groups/group-1", r.URL.Path)
		w.Header().Set("X-Request-Id", "req-delete")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.DeleteGroup(context.Background(), &DeleteGroupRequest{GroupID: "group-1"})
	require.NoError(t, err)
	require.Equal(t, "req-delete", resp.Metadata.RequestID)
}

func TestClientListGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/groups:mine", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "25", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))
		require.NoError(t, json.NewEncoder(w).Encode(ListGroupsResponse{
			Groups: []UserGroup{{GroupID: "group-1"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	pageSize := int32(25)
	resp, err := client.ListGroups(context.Background(), &ListGroupsRequest{
		PageSize:  &pageSize,
		NextToken: "token",
	})
	require.NoError(t, err)
	require.Len(t, resp.Groups, 1)
}

func TestClientAddGroupMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/groups/group-1/invites", r.URL.Path)

		var payload struct {
			InvitedUserID string `json:"invitedUserId"`
			Role          string `json:"role"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "user-1", payload.InvitedUserID)
		require.Equal(t, "member", payload.Role)

		w.Header().Set("X-Request-Id", "req-add")
		require.NoError(t, json.NewEncoder(w).Encode(AddGroupMemberResponse{
			Invite: &GroupInvite{InviteID: "inv-1"},
			Token:  "tok",
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.AddGroupMember(context.Background(), &AddGroupMemberRequest{
		GroupID: "group-1",
		UserID:  "user-1",
		Role:    "member",
	})
	require.NoError(t, err)
	require.Equal(t, "req-add", resp.Metadata.RequestID)
	require.Equal(t, "inv-1", resp.Invite.InviteID)
	require.Equal(t, "tok", resp.Token)
}

func TestClientRemoveGroupMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/groups/group-1/members/user-1", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.RemoveGroupMember(context.Background(), &RemoveGroupMemberRequest{
		GroupID: "group-1",
		UserID:  "user-1",
	})
	require.NoError(t, err)
	require.NotNil(t, resp)
}

func TestClientListGroupMembers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/groups/group-1/members", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "10", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))
		require.NoError(t, json.NewEncoder(w).Encode(ListGroupMembersResponse{
			Members: []GroupMember{{UserID: "user-1"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	pageSize := int32(10)
	resp, err := client.ListGroupMembers(context.Background(), &ListGroupMembersRequest{
		GroupID:   "group-1",
		PageSize:  &pageSize,
		NextToken: "token",
	})
	require.NoError(t, err)
	require.Len(t, resp.Members, 1)
}

func TestClientGetGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/groups/group-1", r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(GetGroupResponse{
			Group: &UserGroup{GroupID: "group-1", Name: "Testers"},
			Links: []GroupLink{{LinkID: "link-1"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetGroup(context.Background(), &GetGroupRequest{GroupID: "group-1"})
	require.NoError(t, err)
	require.Equal(t, "group-1", resp.Group.GroupID)
	require.Equal(t, "link-1", resp.Links[0].LinkID)
}

func TestClientSearchGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/groups/search", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "test", query.Get("searchTerm"))
		require.Equal(t, "25", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))
		require.NoError(t, json.NewEncoder(w).Encode(SearchGroupsResponse{
			Groups: []UserGroup{{GroupID: "group-1"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	pageSize := int32(25)
	resp, err := client.SearchGroups(context.Background(), &SearchGroupsRequest{
		SearchTerm: "test",
		PageSize:   &pageSize,
		NextToken:  "token",
	})
	require.NoError(t, err)
	require.Equal(t, "group-1", resp.Groups[0].GroupID)
}

func TestClientListMyGroups(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/groups:mine", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "25", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))
		require.NoError(t, json.NewEncoder(w).Encode(ListMyGroupsResponse{
			Groups: []UserGroup{{GroupID: "group-1"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	pageSize := int32(25)
	resp, err := client.ListMyGroups(context.Background(), &ListMyGroupsRequest{PageSize: &pageSize, NextToken: "token"})
	require.NoError(t, err)
	require.Equal(t, "group-1", resp.Groups[0].GroupID)
}

func TestClientGroupInvites(t *testing.T) {
	var sawCreate, sawRevoke, sawAccept bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/groups/group-1/invites":
			sawCreate = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user-1", payload["invitedUserId"])
			require.Equal(t, "member", payload["role"])
			require.NoError(t, json.NewEncoder(w).Encode(CreateGroupInviteResponse{
				Invite: &GroupInvite{InviteID: "inv-1"},
				Token:  "tok",
			}))
		case "/v1/groups/group-1/invites/inv-1":
			sawRevoke = true
			require.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusOK)
		case "/v1/groups:acceptInvite":
			sawAccept = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload AcceptGroupInviteRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "tok", payload.Token)
			require.Equal(t, "group-1", payload.GroupID)
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	_, err := client.CreateGroupInvite(context.Background(), &CreateGroupInviteRequest{GroupID: "group-1", InvitedUserID: "user-1", Role: "member"})
	require.NoError(t, err)
	_, err = client.RevokeGroupInvite(context.Background(), &RevokeGroupInviteRequest{GroupID: "group-1", InviteID: "inv-1"})
	require.NoError(t, err)
	_, err = client.AcceptGroupInvite(context.Background(), &AcceptGroupInviteRequest{Token: "tok", GroupID: "group-1"})
	require.NoError(t, err)
	require.True(t, sawCreate)
	require.True(t, sawRevoke)
	require.True(t, sawAccept)
}

func TestClientUpdateGroupMemberRole(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/v1/groups/group-1/members/user-1", r.URL.Path)
		var payload map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "admin", payload["role"])
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	_, err := client.UpdateGroupMemberRole(context.Background(), &UpdateGroupMemberRoleRequest{GroupID: "group-1", UserID: "user-1", Role: "admin"})
	require.NoError(t, err)
}

func TestClientGroupAvatarUploadFlow(t *testing.T) {
	var sawBegin, sawComplete bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/groups/group-1/avatar:beginUpload":
			sawBegin = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "image/png", payload["contentType"])
			require.NoError(t, json.NewEncoder(w).Encode(BeginAvatarUploadResponse{
				UploadID: "up-1",
				MaxBytes: 1024,
			}))
		case "/v1/groups/group-1/avatar:completeUpload":
			sawComplete = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "up-1", payload["uploadId"])
			require.NoError(t, json.NewEncoder(w).Encode(CompleteGroupAvatarUploadResponse{
				Group: &UserGroup{GroupID: "group-1", AvatarURL: "https://cdn/group.png"},
			}))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	_, err := client.BeginGroupAvatarUpload(context.Background(), &BeginGroupAvatarUploadRequest{GroupID: "group-1", ContentType: "image/png"})
	require.NoError(t, err)
	resp, err := client.CompleteGroupAvatarUpload(context.Background(), &CompleteGroupAvatarUploadRequest{GroupID: "group-1", UploadID: "up-1"})
	require.NoError(t, err)
	require.Equal(t, "https://cdn/group.png", resp.Group.AvatarURL)
	require.True(t, sawBegin)
	require.True(t, sawComplete)
}
