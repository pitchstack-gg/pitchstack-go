package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientPasskeysFlow(t *testing.T) {
	var sawInitReg, sawCompleteReg, sawInitAuth, sawCompleteAuth, sawList, sawDelete bool
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/auth/webauthn/registration/initiate":
			sawInitReg = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload InitiatePasskeyRegistrationRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user-1", payload.UserID)
			require.NoError(t, json.NewEncoder(w).Encode(InitiatePasskeyRegistrationResponse{SessionID: "sess"}))
		case "/v1/auth/webauthn/registration/complete":
			sawCompleteReg = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload CompletePasskeyRegistrationRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "sess", payload.SessionID)
			require.NoError(t, json.NewEncoder(w).Encode(CompletePasskeyRegistrationResponse{CredentialID: "cred"}))
		case "/v1/auth/webauthn/authentication/initiate":
			sawInitAuth = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload InitiatePasskeyAuthenticationRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "user@example.com", payload.Email)
			require.NoError(t, json.NewEncoder(w).Encode(InitiatePasskeyAuthenticationResponse{SessionID: "sess"}))
		case "/v1/auth/webauthn/authentication/complete":
			sawCompleteAuth = true
			require.Equal(t, http.MethodPost, r.Method)
			var payload CompletePasskeyAuthenticationRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "sess", payload.SessionID)
			require.Equal(t, "cred", payload.CredentialID)
			require.NoError(t, json.NewEncoder(w).Encode(LoginResponse{UserID: "user-1"}))
		case "/v1/auth/webauthn/users/user-1/credentials":
			sawList = true
			require.Equal(t, http.MethodGet, r.Method)
			require.NoError(t, json.NewEncoder(w).Encode(ListUserPasskeysResponse{Credentials: []Passkey{{CredentialID: "cred"}}}))
		case "/v1/auth/webauthn/users/user-1/credentials/cred":
			sawDelete = true
			require.Equal(t, http.MethodDelete, r.Method)
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	_, err := client.InitiatePasskeyRegistration(context.Background(), &InitiatePasskeyRegistrationRequest{UserID: "user-1"})
	require.NoError(t, err)
	_, err = client.CompletePasskeyRegistration(context.Background(), &CompletePasskeyRegistrationRequest{
		UserID:            "user-1",
		SessionID:         "sess",
		ClientDataJSON:    "cd",
		AttestationObject: "ao",
	})
	require.NoError(t, err)
	_, err = client.InitiatePasskeyAuthentication(context.Background(), &InitiatePasskeyAuthenticationRequest{Email: "user@example.com"})
	require.NoError(t, err)
	_, err = client.CompletePasskeyAuthentication(context.Background(), &CompletePasskeyAuthenticationRequest{
		SessionID:         "sess",
		CredentialID:      "cred",
		ClientDataJSON:    "cd",
		AuthenticatorData: "ad",
		Signature:         "sig",
	})
	require.NoError(t, err)
	_, err = client.ListUserPasskeys(context.Background(), &ListUserPasskeysRequest{UserID: "user-1"})
	require.NoError(t, err)
	_, err = client.DeletePasskey(context.Background(), &DeletePasskeyRequest{UserID: "user-1", CredentialID: "cred"})
	require.NoError(t, err)

	require.True(t, sawInitReg)
	require.True(t, sawCompleteReg)
	require.True(t, sawInitAuth)
	require.True(t, sawCompleteAuth)
	require.True(t, sawList)
	require.True(t, sawDelete)
}
