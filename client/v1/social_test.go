package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientFollowUser(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/users/user-1/followers", r.URL.Path)
		called = true
		w.Header().Set("X-Request-Id", "req-follow")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.FollowUser(context.Background(), &FollowUserRequest{TargetUserID: " user-1 "})
	require.NoError(t, err)
	require.True(t, called)
	require.Equal(t, "req-follow", resp.Metadata.RequestID)

	_, err = client.FollowUser(context.Background(), &FollowUserRequest{})
	require.Error(t, err)
}

func TestClientUnfollowUser(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/users/user-1/followers", r.URL.Path)
		called = true
		w.Header().Set("X-Request-Id", "req-unfollow")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.UnfollowUser(context.Background(), &UnfollowUserRequest{TargetUserID: " user-1 "})
	require.NoError(t, err)
	require.True(t, called)
	require.Equal(t, "req-unfollow", resp.Metadata.RequestID)

	_, err = client.UnfollowUser(context.Background(), &UnfollowUserRequest{})
	require.Error(t, err)
}

func TestClientListFollowers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/users/user-1/followers", r.URL.Path)
		require.Equal(t, "25", r.URL.Query().Get("pageSize"))
		require.Equal(t, "tok-1", r.URL.Query().Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(ListFollowersResponse{
			UserIDs:   []string{"user-2", "user-3"},
			NextToken: "tok-2",
		}))
	}))
	t.Cleanup(server.Close)

	size := int32(25)
	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ListFollowers(context.Background(), &ListFollowersRequest{
		UserID:    " user-1 ",
		PageSize:  &size,
		NextToken: " tok-1 ",
	})
	require.NoError(t, err)
	require.Equal(t, []string{"user-2", "user-3"}, resp.UserIDs)
	require.Equal(t, "tok-2", resp.NextToken)

	_, err = client.ListFollowers(context.Background(), &ListFollowersRequest{})
	require.Error(t, err)
}

func TestClientListFollowing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/users/user-1/following", r.URL.Path)
		require.Equal(t, "50", r.URL.Query().Get("pageSize"))
		require.Equal(t, "token", r.URL.Query().Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(ListFollowingResponse{
			UserIDs:   []string{"user-4"},
			NextToken: "",
		}))
	}))
	t.Cleanup(server.Close)

	size := int32(50)
	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ListFollowing(context.Background(), &ListFollowingRequest{
		UserID:    "user-1",
		PageSize:  &size,
		NextToken: "token",
	})
	require.NoError(t, err)
	require.Equal(t, []string{"user-4"}, resp.UserIDs)

	_, err = client.ListFollowing(context.Background(), &ListFollowingRequest{})
	require.Error(t, err)
}

func TestClientGetFollowStats(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/users/user-1/follow_stats", r.URL.Path)

		require.NoError(t, json.NewEncoder(w).Encode(GetFollowStatsResponse{
			Followers: 10,
			Following: 5,
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetFollowStats(context.Background(), &GetFollowStatsRequest{UserID: "user-1"})
	require.NoError(t, err)
	require.Equal(t, int64(10), resp.Followers)
	require.Equal(t, int64(5), resp.Following)

	_, err = client.GetFollowStats(context.Background(), &GetFollowStatsRequest{})
	require.Error(t, err)
}

func TestClientIsFollowing(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/users/follower-1/following/followee-1", r.URL.Path)

		require.NoError(t, json.NewEncoder(w).Encode(IsFollowingResponse{
			IsFollowing: true,
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.IsFollowing(context.Background(), &IsFollowingRequest{
		FollowerID: " follower-1 ",
		FolloweeID: " followee-1 ",
	})
	require.NoError(t, err)
	require.True(t, resp.IsFollowing)

	_, err = client.IsFollowing(context.Background(), &IsFollowingRequest{})
	require.Error(t, err)
	_, err = client.IsFollowing(context.Background(), &IsFollowingRequest{FollowerID: "follower"})
	require.Error(t, err)
}

func TestClientListActivityFeed(t *testing.T) {
	t.Run("when request valid, then feed returned", func(t *testing.T) {
		size := int32(20)
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/activity", r.URL.Path)
			require.Equal(t, "20", r.URL.Query().Get("pageSize"))
			require.Equal(t, "tok-1", r.URL.Query().Get("nextToken"))
			require.Equal(t, []string{
				string(ActivityScopeFollowing),
				string(ActivityScopeSystem),
			}, r.URL.Query()["scopes"])

			resp := ListActivityFeedResponse{
				Items: []ActivityItem{
					{
						ActivityID: "act-1",
						Kind:       ActivityKindSocial,
						ActorID:    "user-1",
						Verb:       "followed",
					},
				},
				NextToken: "tok-2",
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListActivityFeed(context.Background(), &ListActivityFeedRequest{
			PageSize:  &size,
			NextToken: " tok-1 ",
			Scopes: []ActivityScope{
				ActivityScopeFollowing,
				ActivityScopeUnspecified,
				ActivityScopeSystem,
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
		require.Equal(t, "act-1", resp.Items[0].ActivityID)
		require.Equal(t, "tok-2", resp.NextToken)
	})

	t.Run("when request nil, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.ListActivityFeed(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
