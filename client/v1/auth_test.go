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
			require.Equal(t, "/v1/auth/login", r.URL.Path)
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
			require.Equal(t, "/v1/auth/token/refresh", r.URL.Path)

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

func TestClientCreateCLILoginSession(t *testing.T) {
	t.Run("when provided base url, then creates session", func(t *testing.T) {
		expectedExpiry := time.Now().UTC().Truncate(time.Second)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/auth/cli/sessions", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload CreateCLILoginSessionRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "https://gateway.pitchstack.gg", payload.BaseURL)

			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Request-Id", "req-cli-session")
			require.NoError(t, json.NewEncoder(w).Encode(CreateCLILoginSessionResponse{
				SessionID:           "session-1",
				SessionSecret:       "secret-1",
				VerificationPath:    "/v1/auth/cli/sessions/session-1/login",
				VerificationURL:     "https://gateway.pitchstack.gg/v1/auth/cli/sessions/session-1/login",
				ExpiresAt:           &expectedExpiry,
				PollIntervalSeconds: 2,
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

		resp, err := client.CreateCLILoginSession(context.Background(), &CreateCLILoginSessionRequest{
			BaseURL: "https://gateway.pitchstack.gg",
		})
		require.NoError(t, err)
		require.Equal(t, "session-1", resp.SessionID)
		require.Equal(t, "secret-1", resp.SessionSecret)
		require.Equal(t, "/v1/auth/cli/sessions/session-1/login", resp.VerificationPath)
		require.Equal(t, "https://gateway.pitchstack.gg/v1/auth/cli/sessions/session-1/login", resp.VerificationURL)
		require.Equal(t, expectedExpiry, *resp.ExpiresAt)
		require.Equal(t, int32(2), resp.PollIntervalSeconds)
		require.Equal(t, "req-cli-session", resp.Metadata.RequestID)
	})

	t.Run("when request missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.CreateCLILoginSession(context.Background(), nil)
		require.Nil(t, resp)
		require.Error(t, err)
	})
}

func TestClientGetCLILoginSession(t *testing.T) {
	expectedExpiry := time.Now().UTC().Truncate(time.Second)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/auth/cli/sessions/session-1:poll", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload struct {
			SessionSecret string `json:"sessionSecret"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "secret-1", payload.SessionSecret)

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Request-Id", "req-cli-poll")
		require.NoError(t, json.NewEncoder(w).Encode(GetCLILoginSessionResponse{
			Status: CLILoginSessionStatusComplete,
			Login: &LoginResponse{
				UserID:               "user-1",
				AccessToken:          "access",
				RefreshToken:         "refresh",
				AccessTokenExpiresAt: &expectedExpiry,
				Roles:                []string{"member"},
			},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetCLILoginSession(context.Background(), &GetCLILoginSessionRequest{
		SessionID:     "session-1",
		SessionSecret: "secret-1",
	})
	require.NoError(t, err)
	require.Equal(t, CLILoginSessionStatusComplete, resp.Status)
	require.NotNil(t, resp.Login)
	require.Equal(t, "user-1", resp.Login.UserID)
	require.Equal(t, "access", resp.Login.AccessToken)
	require.Equal(t, "refresh", resp.Login.RefreshToken)
	require.Equal(t, expectedExpiry, *resp.Login.AccessTokenExpiresAt)
	require.Equal(t, "req-cli-poll", resp.Metadata.RequestID)

	resp, err = client.GetCLILoginSession(context.Background(), &GetCLILoginSessionRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientCancelCLILoginSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/auth/cli/sessions/session-1:cancel", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload struct {
			SessionSecret string `json:"sessionSecret"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "secret-1", payload.SessionSecret)

		w.Header().Set("X-Request-Id", "req-cli-cancel")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	resp, err := client.CancelCLILoginSession(context.Background(), &CancelCLILoginSessionRequest{
		SessionID:     "session-1",
		SessionSecret: "secret-1",
	})
	require.NoError(t, err)
	require.Equal(t, "req-cli-cancel", resp.Metadata.RequestID)
}

func TestClientCompleteOAuthForCLILoginSession(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/auth/cli/sessions/session-1/oauth/google:complete", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload struct {
			Code        string `json:"code"`
			State       string `json:"state"`
			RedirectURI string `json:"redirectUri"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "oauth-code", payload.Code)
		require.Equal(t, "oauth-state", payload.State)
		require.Equal(t, "https://gateway.pitchstack.gg/callback", payload.RedirectURI)

		w.Header().Set("X-Request-Id", "req-cli-complete")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	resp, err := client.CompleteOAuthForCLILoginSession(context.Background(), &CompleteOAuthForCLILoginSessionRequest{
		SessionID:   "session-1",
		Provider:    "google",
		Code:        "oauth-code",
		State:       "oauth-state",
		RedirectURI: "https://gateway.pitchstack.gg/callback",
	})
	require.NoError(t, err)
	require.Equal(t, "req-cli-complete", resp.Metadata.RequestID)
}

func TestClientLogout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/auth/logout", r.URL.Path)
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
		require.Equal(t, "/v1/me", r.URL.Path)
		_, _ = w.Write([]byte(`{"user":{"userId":"user-1"},"accessProfile":{"userId":"user-1","limits":{"groups":"10","decks":25},"version":"7"}}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.Me(context.Background())
	require.NoError(t, err)
	require.Equal(t, "user-1", resp.User.UserID)
	require.Equal(t, map[string]int64{"groups": 10, "decks": 25}, resp.AccessProfile.Limits)
	require.EqualValues(t, 7, resp.AccessProfile.Version)
}

func TestClientGetUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/users/user-1", r.URL.Path)
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
		require.Equal(t, "/v1/users/user-1", r.URL.Path)

		var payload map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "user@example.com", payload["email"])
		require.Equal(t, []any{"admin"}, payload["roles"])
		_, hasUsername := payload["username"]
		require.False(t, hasUsername)

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

func TestClientDeleteUser(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/users/user-1", r.URL.Path)
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.DeleteUser(context.Background(), &DeleteUserRequest{UserID: "user-1"})
	require.NoError(t, err)
	require.NotNil(t, resp)

	resp, err = client.DeleteUser(context.Background(), &DeleteUserRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}
