package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientInitiateOAuthLink(t *testing.T) {
	t.Run("when given a provider and proof, then posts to the link initiate path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/auth/oauth/google/link:initiate", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "https://app/cb", payload["redirectUri"])
			require.Equal(t, "consent", payload["prompt"])
			require.Equal(t, "recent-access", payload["reauthenticationAccessToken"])
			require.NotContains(t, payload, "provider")
			_, _ = w.Write([]byte(`{"authorizationUrl":"https://accounts/auth","state":"st"}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.InitiateOAuthLink(context.Background(), &InitiateOAuthLinkRequest{
			Provider:                    "google",
			RedirectURI:                 "https://app/cb",
			Prompt:                      "consent",
			ReauthenticationAccessToken: "recent-access",
		})
		require.NoError(t, err)
		require.Equal(t, "https://accounts/auth", resp.AuthorizationURL)
		require.Equal(t, "st", resp.State)
	})

	t.Run("when the proof is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.InitiateOAuthLink(context.Background(), &InitiateOAuthLinkRequest{Provider: "google"})
		require.Error(t, err)
		_, err = client.InitiateOAuthLink(context.Background(), &InitiateOAuthLinkRequest{ReauthenticationAccessToken: "x"})
		require.Error(t, err)
	})
}

func TestClientCreateProviderAuthChallenge(t *testing.T) {
	t.Run("when given a registration and operation, then decodes the challenge", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/auth/provider-challenges", r.URL.Path)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "ios-apple", payload["clientRegistrationId"])
			require.Equal(t, "PROVIDER_AUTH_OPERATION_SIGN_IN", payload["operation"])
			require.NotContains(t, payload, "reauthenticationAccessToken")
			_, _ = w.Write([]byte(`{"challengeId":"ch-1","challengeSecret":"sec","expiresAt":"2026-01-02T03:04:05Z","proofKind":"PROVIDER_AUTH_PROOF_KIND_APPLE_NATIVE","providerNonce":"abc"}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.CreateProviderAuthChallenge(context.Background(), &CreateProviderAuthChallengeRequest{
			ClientRegistrationID: "ios-apple",
			Operation:            ProviderAuthOperationSignIn,
		})
		require.NoError(t, err)
		require.Equal(t, "ch-1", resp.ChallengeID)
		require.Equal(t, "sec", resp.ChallengeSecret)
		require.Equal(t, ProviderAuthProofKindAppleNative, resp.ProofKind)
		require.Equal(t, "abc", resp.ProviderNonce)
		require.NotNil(t, resp.ExpiresAt)
	})

	t.Run("when the operation is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.CreateProviderAuthChallenge(context.Background(), &CreateProviderAuthChallengeRequest{ClientRegistrationID: "ios-apple"})
		require.Error(t, err)
		_, err = client.CreateProviderAuthChallenge(context.Background(), nil)
		require.Error(t, err)
	})
}

func TestClientCompleteProviderAuthChallenge(t *testing.T) {
	t.Run("when signing in with an Apple proof, then returns the session", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/auth/provider-challenges:complete", r.URL.Path)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "ch-1", payload["challengeId"])
			require.Equal(t, "sec", payload["challengeSecret"])
			apple, ok := payload["appleNative"].(map[string]any)
			require.True(t, ok)
			require.Equal(t, "id-token", apple["identityToken"])
			require.Equal(t, "auth-code", apple["authorizationCode"])
			require.NotContains(t, payload, "googleServerCode")
			_, _ = w.Write([]byte(`{"session":{"userId":"user-1","accessToken":"access","accountCreated":true,"independentRecoveryStatus":"INDEPENDENT_RECOVERY_STATUS_NOT_REQUIRED"}}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.CompleteProviderAuthChallenge(context.Background(), &CompleteProviderAuthChallengeRequest{
			ChallengeID:     "ch-1",
			ChallengeSecret: "sec",
			AppleNative:     &AppleNativeAuthProof{IdentityToken: "id-token", AuthorizationCode: "auth-code"},
		})
		require.NoError(t, err)
		require.Equal(t, "user-1", resp.Session.UserID)
		require.True(t, resp.Session.AccountCreated)
		require.Equal(t, IndependentRecoveryStatusNotRequired, resp.Session.IndependentRecoveryStatus)
		require.Nil(t, resp.Linked)
	})

	t.Run("when linking with a Google code, then returns the linked result", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			google, ok := payload["googleServerCode"].(map[string]any)
			require.True(t, ok)
			require.Equal(t, "server-code", google["authorizationCode"])
			_, _ = w.Write([]byte(`{"linked":{}}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.CompleteProviderAuthChallenge(context.Background(), &CompleteProviderAuthChallengeRequest{
			ChallengeID:      "ch-1",
			ChallengeSecret:  "sec",
			GoogleServerCode: &GoogleServerAuthorizationCodeProof{AuthorizationCode: "server-code"},
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Linked)
		require.Nil(t, resp.Session)
	})

	t.Run("when zero or both proofs are set, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.CompleteProviderAuthChallenge(context.Background(), &CompleteProviderAuthChallengeRequest{ChallengeID: "ch-1", ChallengeSecret: "sec"})
		require.Error(t, err)
		_, err = client.CompleteProviderAuthChallenge(context.Background(), &CompleteProviderAuthChallengeRequest{
			ChallengeID:      "ch-1",
			ChallengeSecret:  "sec",
			AppleNative:      &AppleNativeAuthProof{},
			GoogleServerCode: &GoogleServerAuthorizationCodeProof{},
		})
		require.Error(t, err)
	})

	t.Run("when the secret is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.CompleteProviderAuthChallenge(context.Background(), &CompleteProviderAuthChallengeRequest{ChallengeID: "ch-1", AppleNative: &AppleNativeAuthProof{}})
		require.Error(t, err)
	})
}

func TestClientCancelProviderAuthChallenge(t *testing.T) {
	t.Run("when given a challenge, then posts to the cancel path", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/auth/provider-challenges:cancel", r.URL.Path)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "ch-1", payload["challengeId"])
			require.Equal(t, "sec", payload["challengeSecret"])
			w.Header().Set("X-Request-Id", "req-cancel")
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.CancelProviderAuthChallenge(context.Background(), &CancelProviderAuthChallengeRequest{ChallengeID: "ch-1", ChallengeSecret: "sec"})
		require.NoError(t, err)
		require.Equal(t, "req-cancel", resp.Metadata.RequestID)
	})

	t.Run("when the challenge id is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.CancelProviderAuthChallenge(context.Background(), &CancelProviderAuthChallengeRequest{ChallengeSecret: "sec"})
		require.Error(t, err)
	})
}

func TestClientBeginHostedAppleAuth(t *testing.T) {
	t.Run("when given a registration, destination and operation, then decodes the flow", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/auth/apple/hosted:begin", r.URL.Path)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "web-apple", payload["clientRegistrationId"])
			require.Equal(t, "web-home", payload["onwardDestinationId"])
			require.Equal(t, "PROVIDER_AUTH_OPERATION_LINK", payload["operation"])
			require.Equal(t, "recent-access", payload["reauthenticationAccessToken"])
			_, _ = w.Write([]byte(`{"challengeId":"ch-1","completionSecret":"done","expiresAt":"2026-01-02T03:04:05Z","authorizationUrl":"https://appleid/auth"}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BeginHostedAppleAuth(context.Background(), &BeginHostedAppleAuthRequest{
			ClientRegistrationID:        "web-apple",
			OnwardDestinationID:         "web-home",
			Operation:                   ProviderAuthOperationLink,
			ReauthenticationAccessToken: "recent-access",
		})
		require.NoError(t, err)
		require.Equal(t, "ch-1", resp.ChallengeID)
		require.Equal(t, "done", resp.CompletionSecret)
		require.Equal(t, "https://appleid/auth", resp.AuthorizationURL)
	})

	t.Run("when the destination is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.BeginHostedAppleAuth(context.Background(), &BeginHostedAppleAuthRequest{ClientRegistrationID: "web-apple", Operation: ProviderAuthOperationSignIn})
		require.Error(t, err)
	})
}

func TestClientCompleteHostedAppleAuth(t *testing.T) {
	t.Run("when given the challenge, secret and receipt, then returns the session", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/auth/apple/hosted:complete", r.URL.Path)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "ch-1", payload["challengeId"])
			require.Equal(t, "done", payload["completionSecret"])
			require.Equal(t, "rcpt", payload["callbackReceipt"])
			_, _ = w.Write([]byte(`{"session":{"userId":"user-1"}}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.CompleteHostedAppleAuth(context.Background(), &CompleteHostedAppleAuthRequest{ChallengeID: "ch-1", CompletionSecret: "done", CallbackReceipt: "rcpt"})
		require.NoError(t, err)
		require.Equal(t, "user-1", resp.Session.UserID)
	})

	t.Run("when the receipt is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.CompleteHostedAppleAuth(context.Background(), &CompleteHostedAppleAuthRequest{ChallengeID: "ch-1", CompletionSecret: "done"})
		require.Error(t, err)
	})
}
