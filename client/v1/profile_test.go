package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientSetAvatarURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/api/v1/me/avatar", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload struct {
			AvatarURL string `json:"avatarUrl"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "https://example.com/avatar.png", payload.AvatarURL)

		w.Header().Set("X-Request-Id", "req-avatar")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.SetAvatarURL(context.Background(), &SetAvatarURLRequest{AvatarURL: " https://example.com/avatar.png "})
	require.NoError(t, err)
	require.Equal(t, "req-avatar", resp.Metadata.RequestID)

	_, err = client.SetAvatarURL(context.Background(), &SetAvatarURLRequest{})
	require.Error(t, err)
}

func TestClientUpdateProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/api/v1/me/profile", r.URL.Path)

		var payload struct {
			Profile struct {
				Name string `json:"name"`
				Bio  string `json:"bio"`
			} `json:"profile"`
			UpdateMask string `json:"updateMask"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "Test User", payload.Profile.Name)
		require.Equal(t, "Bio", payload.Profile.Bio)
		require.Equal(t, "name,bio", payload.UpdateMask)

		require.NoError(t, json.NewEncoder(w).Encode(UpdateProfileResponse{
			Profile: &UserProfile{Name: "Test User", Bio: "Bio"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.UpdateProfile(context.Background(), &UpdateProfileRequest{
		Profile: &UserProfile{
			Name: "Test User",
			Bio:  "Bio",
		},
		UpdateMask: "name,bio",
	})
	require.NoError(t, err)
	require.Equal(t, "Test User", resp.Profile.Name)

	_, err = client.UpdateProfile(context.Background(), &UpdateProfileRequest{})
	require.Error(t, err)
}

func TestClientGetProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/v1/users/user-1/profile", r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(GetProfileResponse{
			Profile: &UserProfile{Username: "tester"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetProfile(context.Background(), &GetProfileRequest{UserID: "user-1"})
	require.NoError(t, err)
	require.Equal(t, "tester", resp.Profile.Username)

	_, err = client.GetProfile(context.Background(), &GetProfileRequest{})
	require.Error(t, err)
}

func TestClientGetProfileSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/v1/me/profile/settings", r.URL.Path)
		require.Equal(t, "userId=user-1", r.URL.RawQuery)
		require.NoError(t, json.NewEncoder(w).Encode(GetProfileSettingsResponse{
			Settings: &ProfileSettings{AllowMessages: "ANYONE"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetProfileSettings(context.Background(), &GetProfileSettingsRequest{UserID: "user-1"})
	require.NoError(t, err)
	require.Equal(t, "ANYONE", resp.Settings.AllowMessages)
}

func TestClientUpdateProfileSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/api/v1/me/profile/settings", r.URL.Path)

		var payload struct {
			Settings struct {
				AllowMessages string `json:"allowMessages"`
			} `json:"settings"`
			UpdateMask string `json:"updateMask"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "FOLLOWERS", payload.Settings.AllowMessages)
		require.Equal(t, "allow_messages", payload.UpdateMask)

		require.NoError(t, json.NewEncoder(w).Encode(UpdateProfileSettingsResponse{
			Settings: &ProfileSettings{AllowMessages: "FOLLOWERS"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.UpdateProfileSettings(context.Background(), &UpdateProfileSettingsRequest{
		Settings:   &ProfileSettings{AllowMessages: "FOLLOWERS"},
		UpdateMask: "allow_messages",
	})
	require.NoError(t, err)
	require.Equal(t, "FOLLOWERS", resp.Settings.AllowMessages)

	_, err = client.UpdateProfileSettings(context.Background(), &UpdateProfileSettingsRequest{})
	require.Error(t, err)
}

func TestClientUpsertSocialProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v1/me/social_profiles", r.URL.Path)

		var payload UpsertSocialProfileRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "bluesky", payload.Platform)
		require.Equal(t, "user", payload.Handle)

		require.NoError(t, json.NewEncoder(w).Encode(UpsertSocialProfileResponse{
			SocialProfile: &SocialProfile{Platform: "bluesky", Handle: "user"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.UpsertSocialProfile(context.Background(), &UpsertSocialProfileRequest{
		Platform: "bluesky",
		Handle:   "user",
		URL:      "https://bsky.app/user",
	})
	require.NoError(t, err)
	require.Equal(t, "bluesky", resp.SocialProfile.Platform)

	_, err = client.UpsertSocialProfile(context.Background(), &UpsertSocialProfileRequest{})
	require.Error(t, err)
}

func TestClientRemoveSocialProfile(t *testing.T) {
	var called bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/api/v1/me/social_profiles/bluesky", r.URL.Path)
		called = true
		w.Header().Set("X-Request-Id", "req-remove-social")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.RemoveSocialProfile(context.Background(), &RemoveSocialProfileRequest{Platform: "bluesky"})
	require.NoError(t, err)
	require.True(t, called)
	require.Equal(t, "req-remove-social", resp.Metadata.RequestID)

	_, err = client.RemoveSocialProfile(context.Background(), &RemoveSocialProfileRequest{})
	require.Error(t, err)
}

func TestClientGetSocialProfiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/v1/users/user-1/social_profiles", r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(GetSocialProfilesResponse{
			SocialProfiles: []SocialProfile{{Platform: "bluesky"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetSocialProfiles(context.Background(), &GetSocialProfilesRequest{UserID: "user-1"})
	require.NoError(t, err)
	require.Len(t, resp.SocialProfiles, 1)

	_, err = client.GetSocialProfiles(context.Background(), &GetSocialProfilesRequest{})
	require.Error(t, err)
}
