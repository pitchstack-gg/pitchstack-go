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
		require.Equal(t, "/v1/me/avatar", r.URL.Path)
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

func TestClientGetMyProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/me/profile", r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(GetProfileResponse{
			Profile: &UserProfile{Username: "me"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetMyProfile(context.Background())
	require.NoError(t, err)
	require.Equal(t, "me", resp.Profile.Username)
}

func TestClientSearchUsers(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/users/search", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "alex", query.Get("searchTerm"))
		require.Equal(t, "25", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))
		require.NoError(t, json.NewEncoder(w).Encode(SearchUsersResponse{
			Users: []UserSearchResult{{
				UserID:       "u-00000000-0000-0000-0000-0000000000a1",
				Username:     "alex",
				Name:         "Alex Search",
				AvatarURL:    "https://cdn.example.com/alex.png",
				UserIDSuffix: "00a1",
			}},
			NextToken: "next-page-token",
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	pageSize := int32(25)
	resp, err := client.SearchUsers(context.Background(), &SearchUsersRequest{
		SearchTerm: "alex",
		PageSize:   &pageSize,
		NextToken:  "token",
	})
	require.NoError(t, err)
	require.Len(t, resp.Users, 1)
	require.Equal(t, "u-00000000-0000-0000-0000-0000000000a1", resp.Users[0].UserID)
	require.Equal(t, "00a1", resp.Users[0].UserIDSuffix)
	require.Equal(t, "next-page-token", resp.NextToken)
}

func TestClientAvatarUploadFlow(t *testing.T) {
	var sawBegin, sawComplete bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/me/avatar:beginUpload":
			sawBegin = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload BeginAvatarUploadRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "image/png", payload.ContentType)
			require.NoError(t, json.NewEncoder(w).Encode(BeginAvatarUploadResponse{
				UploadID:  "up-1",
				UploadURL: "https://upload",
				MaxBytes:  1024,
			}))
		case "/v1/me/avatar:completeUpload":
			sawComplete = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload CompleteAvatarUploadRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "up-1", payload.UploadID)
			require.NoError(t, json.NewEncoder(w).Encode(CompleteAvatarUploadResponse{
				Profile: &UserProfile{AvatarURL: "https://cdn/avatar.png"},
			}))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	_, err := client.BeginAvatarUpload(context.Background(), &BeginAvatarUploadRequest{ContentType: "image/png"})
	require.NoError(t, err)
	resp, err := client.CompleteAvatarUpload(context.Background(), &CompleteAvatarUploadRequest{UploadID: "up-1"})
	require.NoError(t, err)
	require.Equal(t, "https://cdn/avatar.png", resp.Profile.AvatarURL)
	require.True(t, sawBegin)
	require.True(t, sawComplete)
}

func TestClientProfileBackgroundUploadFlow(t *testing.T) {
	var sawBegin, sawComplete, sawClear bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/me/profile/background:beginUpload":
			sawBegin = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload BeginProfileBackgroundUploadRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "image/jpeg", payload.ContentType)
			require.NoError(t, json.NewEncoder(w).Encode(BeginProfileBackgroundUploadResponse{
				UploadID:  "bg-up-1",
				UploadURL: "https://upload/background",
				MaxBytes:  2048,
			}))
		case "/v1/me/profile/background:completeUpload":
			sawComplete = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload CompleteProfileBackgroundUploadRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "bg-up-1", payload.UploadID)
			require.NoError(t, json.NewEncoder(w).Encode(CompleteProfileBackgroundUploadResponse{
				Profile: &UserProfile{ProfileBackgroundURL: "https://cdn/background.jpg"},
			}))
		case "/v1/me/profile/background:clear":
			sawClear = true
			require.Equal(t, http.MethodPost, r.Method)
			require.NoError(t, json.NewEncoder(w).Encode(ClearProfileBackgroundResponse{
				Profile: &UserProfile{ProfileBackgroundURL: ""},
			}))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	_, err := client.BeginProfileBackgroundUpload(context.Background(), &BeginProfileBackgroundUploadRequest{ContentType: "image/jpeg"})
	require.NoError(t, err)
	completeResp, err := client.CompleteProfileBackgroundUpload(context.Background(), &CompleteProfileBackgroundUploadRequest{UploadID: "bg-up-1"})
	require.NoError(t, err)
	require.Equal(t, "https://cdn/background.jpg", completeResp.Profile.ProfileBackgroundURL)
	clearResp, err := client.ClearProfileBackground(context.Background())
	require.NoError(t, err)
	require.NotNil(t, clearResp.Profile)
	require.True(t, sawBegin)
	require.True(t, sawComplete)
	require.True(t, sawClear)
}

func TestClientUpdateProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/v1/me/profile", r.URL.Path)

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

func TestClientPrivacyConsent(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			require.Equal(t, "/v1/me/privacy/consent", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(GetPrivacyConsentResponse{
				Consent: &PrivacyConsent{AnalyticsAllowed: true, ConsentVersion: 2},
			}))
		case http.MethodPut:
			require.Equal(t, "/v1/me/privacy/consent", r.URL.Path)
			var payload UpdatePrivacyConsentRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.True(t, payload.AnalyticsAllowed)
			require.Equal(t, int32(2), payload.ConsentVersion)
			require.NoError(t, json.NewEncoder(w).Encode(UpdatePrivacyConsentResponse{
				Consent: &PrivacyConsent{AnalyticsAllowed: payload.AnalyticsAllowed, ConsentVersion: payload.ConsentVersion},
			}))
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	getResp, err := client.GetPrivacyConsent(context.Background())
	require.NoError(t, err)
	require.True(t, getResp.Consent.AnalyticsAllowed)

	updateResp, err := client.UpdatePrivacyConsent(context.Background(), &UpdatePrivacyConsentRequest{
		AnalyticsAllowed: true,
		ConsentVersion:   2,
		Source:           "settings",
		Platform:         "ios",
	})
	require.NoError(t, err)
	require.Equal(t, int32(2), updateResp.Consent.ConsentVersion)
}

func TestClientGetProfile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/users/user-1/profile", r.URL.Path)
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

func TestClientPinnedResources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			require.Equal(t, "/v1/users/user-1/pins", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(GetPinnedResourcesResponse{
				PinnedCollections: []PinnedCollection{{CollectionID: "col-1"}},
				PinnedDecks:       []PinnedDeck{{DeckID: "deck-1"}},
			}))
		case http.MethodPost:
			switch r.URL.Path {
			case "/v1/me/pins/collections":
				var payload PinCollectionRequest
				require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
				require.Equal(t, "col-1", payload.CollectionID)
			case "/v1/me/pins/decks":
				var payload PinDeckRequest
				require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
				require.Equal(t, "deck-1", payload.DeckID)
			default:
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		case http.MethodDelete:
			switch r.URL.Path {
			case "/v1/me/pins/collections/col-1":
			case "/v1/me/pins/decks/deck-1":
			default:
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetPinnedResources(context.Background(), &GetPinnedResourcesRequest{UserID: "user-1"})
	require.NoError(t, err)
	require.Len(t, resp.PinnedCollections, 1)
	require.Len(t, resp.PinnedDecks, 1)

	_, err = client.PinCollection(context.Background(), &PinCollectionRequest{CollectionID: "col-1"})
	require.NoError(t, err)
	_, err = client.UnpinCollection(context.Background(), &UnpinCollectionRequest{CollectionID: "col-1"})
	require.NoError(t, err)
	_, err = client.PinDeck(context.Background(), &PinDeckRequest{DeckID: "deck-1"})
	require.NoError(t, err)
	_, err = client.UnpinDeck(context.Background(), &UnpinDeckRequest{DeckID: "deck-1"})
	require.NoError(t, err)
}

func TestClientGetProfileSettings(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/me/profile/settings", r.URL.Path)
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
		require.Equal(t, "/v1/me/profile/settings", r.URL.Path)

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
		require.Equal(t, "/v1/me/social_profiles", r.URL.Path)

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
		require.Equal(t, "/v1/me/social_profiles/bluesky", r.URL.Path)
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
		require.Equal(t, "/v1/users/user-1/social_profiles", r.URL.Path)
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
