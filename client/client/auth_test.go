package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClientLogin(t *testing.T) {
	t.Run("when provided credentials, then issues POST and decodes response", func(t *testing.T) {
		expectedExpiry := time.Now().UTC().Truncate(time.Second)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/api/v1/auth/login", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload LoginRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user@example.com", payload.Email)
			require.Equal(t, "secret", payload.Password)
			require.Equal(t, "ios", payload.DeviceInfo)

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Request-Id", "req-login")
			require.NoError(t, json.NewEncoder(w).Encode(LoginResponse{
				UserID:               "user-1",
				Username:             "user",
				AccessToken:          "access",
				RefreshToken:         "refresh",
				AccessTokenExpiresAt: &expectedExpiry,
				Roles:                []string{"member"},
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

		resp, err := client.Login(context.Background(), &LoginRequest{
			Email:      "user@example.com",
			Password:   "secret",
			DeviceInfo: "ios",
		})
		require.NoError(t, err)
		require.Equal(t, "user-1", resp.UserID)
		require.Equal(t, "access", resp.AccessToken)
		require.Equal(t, "refresh", resp.RefreshToken)
		require.Equal(t, expectedExpiry, *resp.AccessTokenExpiresAt)
		require.Equal(t, "req-login", resp.Metadata.RequestID)
	})

	t.Run("when request missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.Login(context.Background(), nil)
		require.Nil(t, resp)
		require.Error(t, err)
	})
}

func TestClientRefreshToken(t *testing.T) {
	t.Run("when provided refresh token, then exchanges for access token", func(t *testing.T) {
		expectedExpiry := time.Now().UTC().Truncate(time.Second)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/api/v1/auth/token/refresh", r.URL.Path)

			var payload RefreshTokenRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "refresh", payload.RefreshToken)

			w.Header().Set("Content-Type", "application/json")
			require.NoError(t, json.NewEncoder(w).Encode(RefreshTokenResponse{
				AccessToken:          "new-access",
				RefreshToken:         "new-refresh",
				AccessTokenExpiresAt: &expectedExpiry,
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.RefreshToken(context.Background(), &RefreshTokenRequest{RefreshToken: "refresh"})
		require.NoError(t, err)
		require.Equal(t, "new-access", resp.AccessToken)
		require.Equal(t, "new-refresh", resp.RefreshToken)
		require.Equal(t, expectedExpiry, *resp.AccessTokenExpiresAt)
	})

	t.Run("when refresh token missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.RefreshToken(context.Background(), &RefreshTokenRequest{})
		require.Nil(t, resp)
		require.Error(t, err)
	})
}

func TestClientLogout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v1/auth/logout", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		require.NoError(t, json.NewDecoder(r.Body).Decode(&LogoutRequest{}))
		w.Header().Set("X-Request-Id", "req-logout")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	resp, err := client.Logout(context.Background(), &LogoutRequest{RefreshToken: "refresh"})
	require.NoError(t, err)
	require.Equal(t, "req-logout", resp.Metadata.RequestID)
}

func TestClientMe(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/v1/me", r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(MeResponse{
			User: &User{
				UserID:   "user-1",
				Username: "tester",
			},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.Me(context.Background())
	require.NoError(t, err)
	require.Equal(t, "user-1", resp.User.UserID)
}

func TestClientGetUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/v1/users/user-1", r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(GetUserResponse{
			User: &User{UserID: "user-1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	resp, err := client.GetUser(context.Background(), &GetUserRequest{UserID: "user-1"})
	require.NoError(t, err)
	require.Equal(t, "user-1", resp.User.UserID)

	resp, err = client.GetUser(context.Background(), &GetUserRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientUpdateUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/api/v1/users/user-1", r.URL.Path)

		var payload struct {
			Email string   `json:"email"`
			Roles []string `json:"roles"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "user@example.com", payload.Email)
		require.Equal(t, []string{"admin"}, payload.Roles)

		require.NoError(t, json.NewEncoder(w).Encode(UpdateUserResponse{
			User: &User{UserID: "user-1", Email: "user@example.com"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	email := "user@example.com"
	resp, err := client.UpdateUser(context.Background(), &UpdateUserRequest{
		UserID: "user-1",
		Email:  &email,
		Roles:  []string{"admin"},
	})
	require.NoError(t, err)
	require.Equal(t, "user@example.com", resp.User.Email)
}

func TestClientCreateGroup(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v1/groups", r.URL.Path)
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
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/api/v1/groups/group-1", r.URL.Path)
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
		require.Equal(t, "/api/v1/groups/group-1", r.URL.Path)
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
		require.Equal(t, "/api/v1/groups", r.URL.Path)
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
		require.Equal(t, "/api/v1/groups/group-1/members", r.URL.Path)

		var payload struct {
			UserID string `json:"userId"`
			Role   string `json:"role"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "user-1", payload.UserID)
		require.Equal(t, "member", payload.Role)

		w.Header().Set("X-Request-Id", "req-add")
		w.WriteHeader(http.StatusOK)
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
}

func TestClientRemoveGroupMember(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/api/v1/groups/group-1/members/user-1", r.URL.Path)
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
		require.Equal(t, "/api/v1/groups/group-1/members", r.URL.Path)
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
