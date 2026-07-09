package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListAPIKeys(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/auth/api-keys", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "25", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(ListAPIKeysResponse{
			APIKeys: []APIKey{{APIKeyID: "key-1"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	pageSize := int32(25)
	resp, err := client.ListAPIKeys(context.Background(), &ListAPIKeysRequest{PageSize: &pageSize, NextToken: "token"})
	require.NoError(t, err)
	require.Equal(t, "key-1", resp.APIKeys[0].APIKeyID)
}

func TestClientCreateAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/auth/api-keys", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload CreateAPIKeyRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "my key", payload.Name)
		require.Equal(t, []string{"read"}, payload.Scopes)

		require.NoError(t, json.NewEncoder(w).Encode(CreateAPIKeyResponse{
			APIKey:       &APIKey{APIKeyID: "key-1"},
			PlaintextKey: "pt",
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.CreateAPIKey(context.Background(), &CreateAPIKeyRequest{Name: "my key", Scopes: []string{"read"}})
	require.NoError(t, err)
	require.Equal(t, "key-1", resp.APIKey.APIKeyID)
	require.Equal(t, "pt", resp.PlaintextKey)
}

func TestClientValidateAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/auth/api-keys/validate", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload ValidateAPIKeyRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "psk_test", payload.APIKey)

		require.NoError(t, json.NewEncoder(w).Encode(ValidateAPIKeyResponse{
			Valid:   true,
			Details: &APIKeyDetails{APIKeyID: "key-1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ValidateAPIKey(context.Background(), &ValidateAPIKeyRequest{APIKey: "psk_test"})
	require.NoError(t, err)
	require.True(t, resp.Valid)
	require.Equal(t, "key-1", resp.Details.APIKeyID)
}

func TestClientRevokeAPIKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/auth/api-keys/key-1/revoke", r.URL.Path)
		w.Header().Set("X-Request-Id", "req-revoke")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.RevokeAPIKey(context.Background(), &RevokeAPIKeyRequest{APIKeyID: "key-1"})
	require.NoError(t, err)
	require.Equal(t, "req-revoke", resp.Metadata.RequestID)
}

func TestClientChangePassword(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/auth/change-password", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload ChangePasswordRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "user-1", payload.UserID)
		require.Equal(t, "old", payload.CurrentPassword)
		require.Equal(t, "new", payload.NewPassword)

		w.Header().Set("X-Request-Id", "req-change")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ChangePassword(context.Background(), &ChangePasswordRequest{UserID: "user-1", CurrentPassword: "old", NewPassword: "new"})
	require.NoError(t, err)
	require.Equal(t, "req-change", resp.Metadata.RequestID)
}

func TestClientRegister(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/auth/register", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload RegisterRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "user@example.com", payload.Email)
		require.Equal(t, "pw", payload.Password)

		require.NoError(t, json.NewEncoder(w).Encode(RegisterResponse{UserID: "user-1"}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.Register(context.Background(), &RegisterRequest{Email: "user@example.com", Password: "pw"})
	require.NoError(t, err)
	require.Equal(t, "user-1", resp.UserID)
}

func TestClientPasswordResetFlow(t *testing.T) {
	var sawRequest, sawReset bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/request-password-reset":
			sawRequest = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload RequestPasswordResetRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user@example.com", payload.Email)
			w.WriteHeader(http.StatusOK)
		case "/v1/auth/reset-password":
			sawReset = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload ResetPasswordRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "token", payload.ResetToken)
			require.Equal(t, "new", payload.NewPassword)
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	_, err := client.RequestPasswordReset(context.Background(), &RequestPasswordResetRequest{Email: "user@example.com"})
	require.NoError(t, err)
	_, err = client.ResetPassword(context.Background(), &ResetPasswordRequest{ResetToken: "token", NewPassword: "new"})
	require.NoError(t, err)
	require.True(t, sawRequest)
	require.True(t, sawReset)
}

func TestClientEmailVerificationFlow(t *testing.T) {
	var sawResend, sawVerify bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/resend-verification-email":
			sawResend = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload ResendVerificationEmailRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user-1", payload.UserID)
			w.WriteHeader(http.StatusOK)
		case "/v1/auth/verify-email":
			sawVerify = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload VerifyEmailRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user-1", payload.UserID)
			require.Equal(t, "vtok", payload.VerificationToken)
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	_, err := client.ResendVerificationEmail(context.Background(), &ResendVerificationEmailRequest{UserID: "user-1"})
	require.NoError(t, err)
	_, err = client.VerifyEmail(context.Background(), &VerifyEmailRequest{UserID: "user-1", VerificationToken: "vtok"})
	require.NoError(t, err)
	require.True(t, sawResend)
	require.True(t, sawVerify)
}

func TestClientEmailChangeFlow(t *testing.T) {
	var sawStatus, sawRequest, sawResend, sawCancel, sawConfirm bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/email-change":
			switch r.Method {
			case http.MethodGet:
				sawStatus = true
				require.NoError(t, json.NewEncoder(w).Encode(GetEmailChangeStatusResponse{
					PendingChange: &EmailChangeRequest{NewEmail: "new@example.com"},
				}))
			case http.MethodPost:
				sawRequest = true
				var payload RequestEmailChangeRequest
				require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
				require.Equal(t, "new@example.com", payload.NewEmail)
				require.NoError(t, json.NewEncoder(w).Encode(RequestEmailChangeResponse{
					PendingChange: &EmailChangeRequest{NewEmail: "new@example.com"},
				}))
			case http.MethodDelete:
				sawCancel = true
				w.WriteHeader(http.StatusOK)
			default:
				t.Fatalf("unexpected method for email-change: %s", r.Method)
			}
		case "/v1/auth/email-change/resend":
			sawResend = true
			require.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(http.StatusOK)
		case "/v1/auth/email-change/confirm":
			sawConfirm = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload ConfirmEmailChangeRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user-1", payload.UserID)
			require.Equal(t, "change-token", payload.ChangeToken)
			require.NoError(t, json.NewEncoder(w).Encode(ConfirmEmailChangeResponse{
				User: &User{UserID: "user-1", Email: "new@example.com"},
			}))
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	statusResp, err := client.GetEmailChangeStatus(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, "new@example.com", statusResp.PendingChange.NewEmail)

	requestResp, err := client.RequestEmailChange(context.Background(), &RequestEmailChangeRequest{NewEmail: "new@example.com"})
	require.NoError(t, err)
	require.Equal(t, "new@example.com", requestResp.PendingChange.NewEmail)

	_, err = client.ResendEmailChangeConfirmation(context.Background(), nil)
	require.NoError(t, err)
	_, err = client.CancelEmailChange(context.Background(), nil)
	require.NoError(t, err)
	confirmResp, err := client.ConfirmEmailChange(context.Background(), &ConfirmEmailChangeRequest{
		UserID:      "user-1",
		ChangeToken: "change-token",
	})
	require.NoError(t, err)
	require.Equal(t, "new@example.com", confirmResp.User.Email)

	require.True(t, sawStatus)
	require.True(t, sawRequest)
	require.True(t, sawResend)
	require.True(t, sawCancel)
	require.True(t, sawConfirm)
}

func TestClientValidateToken(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/auth/token/validate", r.URL.Path)
		var payload ValidateTokenRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "at", payload.AccessToken)
		require.NoError(t, json.NewEncoder(w).Encode(ValidateTokenResponse{Valid: true}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ValidateToken(context.Background(), &ValidateTokenRequest{AccessToken: "at"})
	require.NoError(t, err)
	require.True(t, resp.Valid)
}

func TestClientOAuthFlow(t *testing.T) {
	var sawInitiate, sawComplete, sawLink, sawUnlink bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/oauth/google/initiate":
			sawInitiate = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "https://cb", payload["redirectUri"])
			require.NoError(t, json.NewEncoder(w).Encode(InitiateOAuthResponse{AuthorizationURL: "https://auth"}))
		case "/v1/auth/oauth/google/complete":
			sawComplete = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "code", payload["code"])
			require.NoError(t, json.NewEncoder(w).Encode(LoginResponse{UserID: "user-1"}))
		case "/v1/auth/oauth/google/link":
			sawLink = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user-1", payload["userId"])
			require.Equal(t, "code", payload["code"])
			w.WriteHeader(http.StatusOK)
		case "/v1/auth/oauth/google/unlink":
			sawUnlink = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user-1", payload["userId"])
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	_, err := client.InitiateOAuth(context.Background(), &InitiateOAuthRequest{Provider: "google", RedirectURI: "https://cb"})
	require.NoError(t, err)
	_, err = client.CompleteOAuth(context.Background(), &CompleteOAuthRequest{Provider: "google", Code: "code"})
	require.NoError(t, err)
	_, err = client.LinkOAuthProvider(context.Background(), &LinkOAuthProviderRequest{Provider: "google", UserID: "user-1", Code: "code"})
	require.NoError(t, err)
	_, err = client.UnlinkOAuthProvider(context.Background(), &UnlinkOAuthProviderRequest{Provider: "google", UserID: "user-1"})
	require.NoError(t, err)
	require.True(t, sawInitiate)
	require.True(t, sawComplete)
	require.True(t, sawLink)
	require.True(t, sawUnlink)
}

func TestClientAuthMethodManagement(t *testing.T) {
	var sawList, sawRemove, sawPreferred bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/users/user-1/methods":
			sawList = true
			require.Equal(t, http.MethodGet, r.Method)
			require.NoError(t, json.NewEncoder(w).Encode(ListAuthMethodsResponse{
				Methods: []AuthMethod{{MethodType: AuthMethodTypePassword, Preferred: true}},
			}))
		case "/v1/auth/users/user-1/methods/AUTH_METHOD_TYPE_PASSWORD":
			sawRemove = true
			require.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusOK)
		case "/v1/auth/users/user-1/methods/AUTH_METHOD_TYPE_PASSWORD/preferred":
			sawPreferred = true
			require.Equal(t, http.MethodPost, r.Method)
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ListAuthMethods(context.Background(), &ListAuthMethodsRequest{UserID: "user-1"})
	require.NoError(t, err)
	require.Equal(t, AuthMethodTypePassword, resp.Methods[0].MethodType)
	_, err = client.RemoveAuthMethod(context.Background(), &RemoveAuthMethodRequest{UserID: "user-1", MethodType: AuthMethodTypePassword})
	require.NoError(t, err)
	_, err = client.SetPreferredAuthMethod(context.Background(), &SetPreferredAuthMethodRequest{UserID: "user-1", MethodType: AuthMethodTypePassword})
	require.NoError(t, err)
	require.True(t, sawList)
	require.True(t, sawRemove)
	require.True(t, sawPreferred)
}
