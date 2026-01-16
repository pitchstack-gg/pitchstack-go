package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// InitiatePasskeyRegistrationRequest starts passkey registration.
type InitiatePasskeyRegistrationRequest struct {
	UserID      string `json:"userId,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// InitiatePasskeyRegistrationResponse contains WebAuthn creation options.
type InitiatePasskeyRegistrationResponse struct {
	PublicKeyCreationOptions string           `json:"publicKeyCreationOptions,omitempty"`
	SessionID                string           `json:"sessionId,omitempty"`
	Metadata                 ResponseMetadata `json:"-"`
}

func (r *InitiatePasskeyRegistrationResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CompletePasskeyRegistrationRequest completes passkey registration.
type CompletePasskeyRegistrationRequest struct {
	UserID            string `json:"userId,omitempty"`
	SessionID         string `json:"sessionId,omitempty"`
	ClientDataJSON    string `json:"clientDataJson,omitempty"`
	AttestationObject string `json:"attestationObject,omitempty"`
}

// CompletePasskeyRegistrationResponse returns the registered credential ID.
type CompletePasskeyRegistrationResponse struct {
	CredentialID string           `json:"credentialId,omitempty"`
	Metadata     ResponseMetadata `json:"-"`
}

func (r *CompletePasskeyRegistrationResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// InitiatePasskeyAuthenticationRequest starts passkey authentication.
type InitiatePasskeyAuthenticationRequest struct {
	Email string `json:"email,omitempty"`
}

// InitiatePasskeyAuthenticationResponse contains WebAuthn request options.
type InitiatePasskeyAuthenticationResponse struct {
	PublicKeyRequestOptions string           `json:"publicKeyRequestOptions,omitempty"`
	SessionID               string           `json:"sessionId,omitempty"`
	Metadata                ResponseMetadata `json:"-"`
}

func (r *InitiatePasskeyAuthenticationResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CompletePasskeyAuthenticationRequest completes passkey authentication.
type CompletePasskeyAuthenticationRequest struct {
	SessionID         string `json:"sessionId,omitempty"`
	CredentialID      string `json:"credentialId,omitempty"`
	ClientDataJSON    string `json:"clientDataJson,omitempty"`
	AuthenticatorData string `json:"authenticatorData,omitempty"`
	Signature         string `json:"signature,omitempty"`
	UserHandle        string `json:"userHandle,omitempty"`
}

// Passkey represents v1Passkey.
type Passkey struct {
	CredentialID string     `json:"credentialId,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	LastUsedAt   *time.Time `json:"lastUsedAt,omitempty"`
	Transports   []string   `json:"transports,omitempty"`
}

// ListUserPasskeysRequest identifies the user to inspect.
type ListUserPasskeysRequest struct {
	UserID string
}

// ListUserPasskeysResponse lists passkeys for the user.
type ListUserPasskeysResponse struct {
	Credentials []Passkey        `json:"credentials,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *ListUserPasskeysResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeletePasskeyRequest identifies the passkey to delete.
type DeletePasskeyRequest struct {
	UserID       string
	CredentialID string
}

// DeletePasskeyResponse captures metadata for delete operations.
type DeletePasskeyResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeletePasskeyResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// InitiatePasskeyRegistration initiates WebAuthn passkey registration.
func (c *Client) InitiatePasskeyRegistration(ctx context.Context, request *InitiatePasskeyRegistrationRequest, opts ...RequestOpt) (*InitiatePasskeyRegistrationResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.UserID) == "" {
		return nil, errors.New("userID must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/webauthn/registration/initiate", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &InitiatePasskeyRegistrationResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CompletePasskeyRegistration completes WebAuthn passkey registration.
func (c *Client) CompletePasskeyRegistration(ctx context.Context, request *CompletePasskeyRegistrationRequest, opts ...RequestOpt) (*CompletePasskeyRegistrationResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.UserID) == "" {
		return nil, errors.New("userID must not be empty")
	}
	if strings.TrimSpace(request.SessionID) == "" {
		return nil, errors.New("sessionID must not be empty")
	}
	if strings.TrimSpace(request.ClientDataJSON) == "" {
		return nil, errors.New("clientDataJson must not be empty")
	}
	if strings.TrimSpace(request.AttestationObject) == "" {
		return nil, errors.New("attestationObject must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/webauthn/registration/complete", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CompletePasskeyRegistrationResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// InitiatePasskeyAuthentication initiates WebAuthn passkey authentication.
func (c *Client) InitiatePasskeyAuthentication(ctx context.Context, request *InitiatePasskeyAuthenticationRequest, opts ...RequestOpt) (*InitiatePasskeyAuthenticationResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.Email) == "" {
		return nil, errors.New("email must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/webauthn/authentication/initiate", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &InitiatePasskeyAuthenticationResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CompletePasskeyAuthentication completes WebAuthn passkey authentication and returns login credentials.
func (c *Client) CompletePasskeyAuthentication(ctx context.Context, request *CompletePasskeyAuthenticationRequest, opts ...RequestOpt) (*LoginResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.SessionID) == "" {
		return nil, errors.New("sessionID must not be empty")
	}
	if strings.TrimSpace(request.CredentialID) == "" {
		return nil, errors.New("credentialId must not be empty")
	}
	if strings.TrimSpace(request.ClientDataJSON) == "" {
		return nil, errors.New("clientDataJson must not be empty")
	}
	if strings.TrimSpace(request.AuthenticatorData) == "" {
		return nil, errors.New("authenticatorData must not be empty")
	}
	if strings.TrimSpace(request.Signature) == "" {
		return nil, errors.New("signature must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/webauthn/authentication/complete", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &LoginResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListUserPasskeys lists passkeys for a user.
func (c *Client) ListUserPasskeys(ctx context.Context, request *ListUserPasskeysRequest, opts ...RequestOpt) (*ListUserPasskeysResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/v1/auth/webauthn/users/%s/credentials", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &ListUserPasskeysResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// DeletePasskey deletes a passkey credential for a user.
func (c *Client) DeletePasskey(ctx context.Context, request *DeletePasskeyRequest, opts ...RequestOpt) (*DeletePasskeyResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}
	credentialID := strings.TrimSpace(request.CredentialID)
	if credentialID == "" {
		return nil, errors.New("credentialID must not be empty")
	}

	path := fmt.Sprintf("/v1/auth/webauthn/users/%s/credentials/%s", url.PathEscape(userID), url.PathEscape(credentialID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeletePasskeyResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
