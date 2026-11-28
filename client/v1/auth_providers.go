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

// ProviderAuthOperation selects whether a provider flow signs in or links.
type ProviderAuthOperation string

const (
	ProviderAuthOperationUnspecified ProviderAuthOperation = "PROVIDER_AUTH_OPERATION_UNSPECIFIED"
	ProviderAuthOperationSignIn      ProviderAuthOperation = "PROVIDER_AUTH_OPERATION_SIGN_IN"
	ProviderAuthOperationLink        ProviderAuthOperation = "PROVIDER_AUTH_OPERATION_LINK"
)

// ProviderAuthProofKind identifies the proof a provider challenge expects.
type ProviderAuthProofKind string

const (
	ProviderAuthProofKindUnspecified      ProviderAuthProofKind = "PROVIDER_AUTH_PROOF_KIND_UNSPECIFIED"
	ProviderAuthProofKindAppleNative      ProviderAuthProofKind = "PROVIDER_AUTH_PROOF_KIND_APPLE_NATIVE"
	ProviderAuthProofKindGoogleServerCode ProviderAuthProofKind = "PROVIDER_AUTH_PROOF_KIND_GOOGLE_SERVER_CODE"
)

// InitiateOAuthLinkRequest starts an OAuth flow that links a provider to the
// authenticated account.
type InitiateOAuthLinkRequest struct {
	Provider    string `json:"-"`
	RedirectURI string `json:"redirectUri,omitempty"`
	Prompt      string `json:"prompt,omitempty"`
	// ReauthenticationAccessToken is an access token from a recent sign-in.
	ReauthenticationAccessToken string `json:"reauthenticationAccessToken,omitempty"`
}

// CreateProviderAuthChallengeRequest starts a native provider sign-in or link.
type CreateProviderAuthChallengeRequest struct {
	ClientRegistrationID string                `json:"clientRegistrationId,omitempty"`
	Operation            ProviderAuthOperation `json:"operation,omitempty"`
	// ReauthenticationAccessToken is required only for link operations.
	ReauthenticationAccessToken string `json:"reauthenticationAccessToken,omitempty"`
}

// CreateProviderAuthChallengeResponse returns a short-lived provider challenge.
type CreateProviderAuthChallengeResponse struct {
	ChallengeID     string                `json:"challengeId,omitempty"`
	ChallengeSecret string                `json:"challengeSecret,omitempty"`
	ExpiresAt       *time.Time            `json:"expiresAt,omitempty"`
	ProofKind       ProviderAuthProofKind `json:"proofKind,omitempty"`
	// ProviderNonce is passed unchanged to the provider SDK when present.
	ProviderNonce string           `json:"providerNonce,omitempty"`
	Metadata      ResponseMetadata `json:"-"`
}

func (r *CreateProviderAuthChallengeResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// AppleNativeAuthProof carries the credentials returned by native Sign in with Apple.
type AppleNativeAuthProof struct {
	IdentityToken     string `json:"identityToken,omitempty"`
	AuthorizationCode string `json:"authorizationCode,omitempty"`
}

// GoogleServerAuthorizationCodeProof carries a Google server authorization code.
type GoogleServerAuthorizationCodeProof struct {
	AuthorizationCode string `json:"authorizationCode,omitempty"`
}

// CompleteProviderAuthChallengeRequest completes a provider challenge. Set
// exactly one of AppleNative or GoogleServerCode.
type CompleteProviderAuthChallengeRequest struct {
	ChallengeID      string                              `json:"challengeId,omitempty"`
	ChallengeSecret  string                              `json:"challengeSecret,omitempty"`
	AppleNative      *AppleNativeAuthProof               `json:"appleNative,omitempty"`
	GoogleServerCode *GoogleServerAuthorizationCodeProof `json:"googleServerCode,omitempty"`
}

// ProviderAuthLinkedResult indicates that a provider was linked.
type ProviderAuthLinkedResult struct{}

// CompleteProviderAuthChallengeResponse returns either a session (sign-in) or
// a linked result (link).
type CompleteProviderAuthChallengeResponse struct {
	Session  *LoginResponse            `json:"session,omitempty"`
	Linked   *ProviderAuthLinkedResult `json:"linked,omitempty"`
	Metadata ResponseMetadata          `json:"-"`
}

func (r *CompleteProviderAuthChallengeResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CancelProviderAuthChallengeRequest cancels an unfinished provider challenge.
type CancelProviderAuthChallengeRequest struct {
	ChallengeID     string `json:"challengeId,omitempty"`
	ChallengeSecret string `json:"challengeSecret,omitempty"`
}

// CancelProviderAuthChallengeResponse captures metadata for cancel operations.
type CancelProviderAuthChallengeResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *CancelProviderAuthChallengeResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BeginHostedAppleAuthRequest starts a browser-hosted Sign in with Apple flow.
type BeginHostedAppleAuthRequest struct {
	ClientRegistrationID string                `json:"clientRegistrationId,omitempty"`
	OnwardDestinationID  string                `json:"onwardDestinationId,omitempty"`
	Operation            ProviderAuthOperation `json:"operation,omitempty"`
	// ReauthenticationAccessToken is required only for link operations.
	ReauthenticationAccessToken string `json:"reauthenticationAccessToken,omitempty"`
}

// BeginHostedAppleAuthResponse returns the Apple authorization URL and the
// secret needed to complete the flow.
type BeginHostedAppleAuthResponse struct {
	ChallengeID      string           `json:"challengeId,omitempty"`
	CompletionSecret string           `json:"completionSecret,omitempty"`
	ExpiresAt        *time.Time       `json:"expiresAt,omitempty"`
	AuthorizationURL string           `json:"authorizationUrl,omitempty"`
	Metadata         ResponseMetadata `json:"-"`
}

func (r *BeginHostedAppleAuthResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CompleteHostedAppleAuthRequest completes a browser-hosted Sign in with Apple flow.
type CompleteHostedAppleAuthRequest struct {
	ChallengeID      string `json:"challengeId,omitempty"`
	CompletionSecret string `json:"completionSecret,omitempty"`
	// CallbackReceipt is the handle delivered to the onward destination.
	CallbackReceipt string `json:"callbackReceipt,omitempty"`
}

// InitiateOAuthLink starts an OAuth flow that links a provider to the
// authenticated account. It requires an access token from a recent sign-in.
func (c *Client) InitiateOAuthLink(ctx context.Context, request *InitiateOAuthLinkRequest, opts ...RequestOpt) (*InitiateOAuthResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	provider := strings.TrimSpace(request.Provider)
	if provider == "" {
		return nil, errors.New("provider must not be empty")
	}
	if strings.TrimSpace(request.ReauthenticationAccessToken) == "" {
		return nil, errors.New("reauthenticationAccessToken must not be empty")
	}

	body, err := jsonBody(struct {
		RedirectURI                 string `json:"redirectUri,omitempty"`
		Prompt                      string `json:"prompt,omitempty"`
		ReauthenticationAccessToken string `json:"reauthenticationAccessToken,omitempty"`
	}{
		RedirectURI:                 strings.TrimSpace(request.RedirectURI),
		Prompt:                      strings.TrimSpace(request.Prompt),
		ReauthenticationAccessToken: strings.TrimSpace(request.ReauthenticationAccessToken),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/oauth/%s/link:initiate", url.PathEscape(provider))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &InitiateOAuthResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CreateProviderAuthChallenge starts a native provider sign-in or link
// challenge.
func (c *Client) CreateProviderAuthChallenge(ctx context.Context, request *CreateProviderAuthChallengeRequest, opts ...RequestOpt) (*CreateProviderAuthChallengeResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ClientRegistrationID) == "" {
		return nil, errors.New("clientRegistrationID must not be empty")
	}
	if request.Operation == "" {
		return nil, errors.New("operation must not be empty")
	}

	response := &CreateProviderAuthChallengeResponse{}
	if err := c.postAuthJSON(ctx, "/v1/auth/provider-challenges", request, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CompleteProviderAuthChallenge completes a provider challenge with the
// provider's proof.
func (c *Client) CompleteProviderAuthChallenge(ctx context.Context, request *CompleteProviderAuthChallengeRequest, opts ...RequestOpt) (*CompleteProviderAuthChallengeResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := requireChallenge(request.ChallengeID, request.ChallengeSecret, "challengeSecret"); err != nil {
		return nil, err
	}
	if (request.AppleNative == nil) == (request.GoogleServerCode == nil) {
		return nil, errors.New("exactly one of appleNative or googleServerCode must be set")
	}

	response := &CompleteProviderAuthChallengeResponse{}
	if err := c.postAuthJSON(ctx, "/v1/auth/provider-challenges:complete", request, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CancelProviderAuthChallenge cancels an unfinished provider challenge.
func (c *Client) CancelProviderAuthChallenge(ctx context.Context, request *CancelProviderAuthChallengeRequest, opts ...RequestOpt) (*CancelProviderAuthChallengeResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := requireChallenge(request.ChallengeID, request.ChallengeSecret, "challengeSecret"); err != nil {
		return nil, err
	}

	response := &CancelProviderAuthChallengeResponse{}
	if err := c.postAuthJSON(ctx, "/v1/auth/provider-challenges:cancel", request, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BeginHostedAppleAuth starts a browser-hosted Sign in with Apple flow and
// returns the URL to open.
func (c *Client) BeginHostedAppleAuth(ctx context.Context, request *BeginHostedAppleAuthRequest, opts ...RequestOpt) (*BeginHostedAppleAuthResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ClientRegistrationID) == "" {
		return nil, errors.New("clientRegistrationID must not be empty")
	}
	if strings.TrimSpace(request.OnwardDestinationID) == "" {
		return nil, errors.New("onwardDestinationID must not be empty")
	}
	if request.Operation == "" {
		return nil, errors.New("operation must not be empty")
	}

	response := &BeginHostedAppleAuthResponse{}
	if err := c.postAuthJSON(ctx, "/v1/auth/apple/hosted:begin", request, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CompleteHostedAppleAuth completes a browser-hosted Sign in with Apple flow.
func (c *Client) CompleteHostedAppleAuth(ctx context.Context, request *CompleteHostedAppleAuthRequest, opts ...RequestOpt) (*CompleteProviderAuthChallengeResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := requireChallenge(request.ChallengeID, request.CompletionSecret, "completionSecret"); err != nil {
		return nil, err
	}
	if strings.TrimSpace(request.CallbackReceipt) == "" {
		return nil, errors.New("callbackReceipt must not be empty")
	}

	response := &CompleteProviderAuthChallengeResponse{}
	if err := c.postAuthJSON(ctx, "/v1/auth/apple/hosted:complete", request, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

func requireChallenge(challengeID, secret, secretName string) error {
	if strings.TrimSpace(challengeID) == "" {
		return errors.New("challengeID must not be empty")
	}
	if strings.TrimSpace(secret) == "" {
		return fmt.Errorf("%s must not be empty", secretName)
	}
	return nil
}

func (c *Client) postAuthJSON(ctx context.Context, path string, payload, response any, opts ...RequestOpt) error {
	body, err := jsonBody(payload)
	if err != nil {
		return fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	return c.do(req, response, opts...)
}
