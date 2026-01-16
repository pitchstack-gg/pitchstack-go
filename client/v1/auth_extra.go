package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// APIKey represents v1APIKey.
type APIKey struct {
	APIKeyID           string     `json:"apiKeyId,omitempty"`
	UserID             string     `json:"userId,omitempty"`
	Name               string     `json:"name,omitempty"`
	KeyPrefix          string     `json:"keyPrefix,omitempty"`
	Scopes             []string   `json:"scopes,omitempty"`
	RateLimitPerMinute int32      `json:"rateLimitPerMinute,omitempty"`
	ExpiresAt          *time.Time `json:"expiresAt,omitempty"`
	LastUsedAt         *time.Time `json:"lastUsedAt,omitempty"`
	CreatedAt          *time.Time `json:"createdAt,omitempty"`
}

// APIKeyDetails represents v1APIKeyDetails.
type APIKeyDetails struct {
	APIKeyID           string   `json:"apiKeyId,omitempty"`
	UserID             string   `json:"userId,omitempty"`
	Scopes             []string `json:"scopes,omitempty"`
	RateLimitPerMinute int32    `json:"rateLimitPerMinute,omitempty"`
}

// ListAPIKeysRequest paginates ListAPIKeys.
type ListAPIKeysRequest struct {
	PageSize  *int32
	NextToken string
}

// ListAPIKeysResponse lists API keys for the authenticated user.
type ListAPIKeysResponse struct {
	APIKeys   []APIKey         `json:"apiKeys,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListAPIKeysResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateAPIKeyRequest creates an API key.
type CreateAPIKeyRequest struct {
	Name               string     `json:"name,omitempty"`
	Scopes             []string   `json:"scopes,omitempty"`
	RateLimitPerMinute *int32     `json:"rateLimitPerMinute,omitempty"`
	ExpiresAt          *time.Time `json:"expiresAt,omitempty"`
}

// CreateAPIKeyResponse returns the created API key and plaintext key (once).
type CreateAPIKeyResponse struct {
	APIKey       *APIKey          `json:"apiKey,omitempty"`
	PlaintextKey string           `json:"plaintextKey,omitempty"`
	Metadata     ResponseMetadata `json:"-"`
}

func (r *CreateAPIKeyResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ValidateAPIKeyRequest validates a provided API key.
type ValidateAPIKeyRequest struct {
	APIKey string `json:"apiKey,omitempty"`
}

// ValidateAPIKeyResponse returns validation details.
type ValidateAPIKeyResponse struct {
	Valid    bool             `json:"valid,omitempty"`
	Details  *APIKeyDetails   `json:"details,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ValidateAPIKeyResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RevokeAPIKeyRequest identifies the API key to revoke.
type RevokeAPIKeyRequest struct {
	APIKeyID string
}

// RevokeAPIKeyResponse captures metadata for revoke operations.
type RevokeAPIKeyResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *RevokeAPIKeyResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ChangePasswordRequest changes a user's password.
type ChangePasswordRequest struct {
	UserID          string `json:"userId,omitempty"`
	CurrentPassword string `json:"currentPassword,omitempty"`
	NewPassword     string `json:"newPassword,omitempty"`
}

// ChangePasswordResponse captures metadata for change-password operations.
type ChangePasswordResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *ChangePasswordResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RequestPasswordResetRequest starts a password reset flow.
type RequestPasswordResetRequest struct {
	Email string `json:"email,omitempty"`
}

// RequestPasswordResetResponse captures metadata for password reset requests.
type RequestPasswordResetResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *RequestPasswordResetResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ResetPasswordRequest completes a password reset flow.
type ResetPasswordRequest struct {
	ResetToken  string `json:"resetToken,omitempty"`
	NewPassword string `json:"newPassword,omitempty"`
}

// ResetPasswordResponse captures metadata for password reset completion.
type ResetPasswordResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *ResetPasswordResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ResendVerificationEmailRequest resends a verification email.
type ResendVerificationEmailRequest struct {
	UserID string `json:"userId,omitempty"`
}

// ResendVerificationEmailResponse captures metadata for resend operations.
type ResendVerificationEmailResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *ResendVerificationEmailResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// VerifyEmailRequest verifies a user's email using a token.
type VerifyEmailRequest struct {
	UserID            string `json:"userId,omitempty"`
	VerificationToken string `json:"verificationToken,omitempty"`
}

// VerifyEmailResponse captures metadata for verification operations.
type VerifyEmailResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *VerifyEmailResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ValidateTokenRequest validates an access token.
type ValidateTokenRequest struct {
	AccessToken string `json:"accessToken,omitempty"`
}

// TokenClaims represents v1TokenClaims.
type TokenClaims struct {
	UserID    string     `json:"userId,omitempty"`
	Roles     []string   `json:"roles,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

// ValidateTokenResponse returns token validation results.
type ValidateTokenResponse struct {
	Valid    bool             `json:"valid,omitempty"`
	Claims   *TokenClaims     `json:"claims,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ValidateTokenResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RegisterRequest registers a new account.
type RegisterRequest struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Phone    string `json:"phone,omitempty"`
}

// RegisterResponse returns credentials created during registration.
type RegisterResponse struct {
	UserID               string           `json:"userId,omitempty"`
	AccessToken          string           `json:"accessToken,omitempty"`
	RefreshToken         string           `json:"refreshToken,omitempty"`
	AccessTokenExpiresAt *time.Time       `json:"accessTokenExpiresAt,omitempty"`
	Roles                []string         `json:"roles,omitempty"`
	Metadata             ResponseMetadata `json:"-"`
}

func (r *RegisterResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// InitiateOAuthRequest initiates an OAuth flow.
type InitiateOAuthRequest struct {
	Provider    string `json:"-"`
	RedirectURI string `json:"redirectUri,omitempty"`
	Prompt      string `json:"prompt,omitempty"`
}

// InitiateOAuthResponse returns the provider authorization URL.
type InitiateOAuthResponse struct {
	AuthorizationURL string           `json:"authorizationUrl,omitempty"`
	State            string           `json:"state,omitempty"`
	Metadata         ResponseMetadata `json:"-"`
}

func (r *InitiateOAuthResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CompleteOAuthRequest completes an OAuth flow with a provider.
type CompleteOAuthRequest struct {
	Provider    string `json:"-"`
	Code        string `json:"code,omitempty"`
	State       string `json:"state,omitempty"`
	RedirectURI string `json:"redirectUri,omitempty"`
}

// LinkOAuthProviderRequest links an OAuth provider to a user.
type LinkOAuthProviderRequest struct {
	Provider    string `json:"-"`
	UserID      string `json:"userId,omitempty"`
	Code        string `json:"code,omitempty"`
	State       string `json:"state,omitempty"`
	RedirectURI string `json:"redirectUri,omitempty"`
}

// LinkOAuthProviderResponse captures metadata for link operations.
type LinkOAuthProviderResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *LinkOAuthProviderResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnlinkOAuthProviderRequest unlinks an OAuth provider from a user.
type UnlinkOAuthProviderRequest struct {
	Provider string `json:"-"`
	UserID   string `json:"userId,omitempty"`
}

// UnlinkOAuthProviderResponse captures metadata for unlink operations.
type UnlinkOAuthProviderResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnlinkOAuthProviderResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// AuthMethodType represents v1AuthMethodType.
type AuthMethodType string

const (
	AuthMethodTypeUnspecified  AuthMethodType = "AUTH_METHOD_TYPE_UNSPECIFIED"
	AuthMethodTypePassword     AuthMethodType = "AUTH_METHOD_TYPE_PASSWORD"
	AuthMethodTypeOAuthGoogle  AuthMethodType = "AUTH_METHOD_TYPE_OAUTH_GOOGLE"
	AuthMethodTypeOAuthDiscord AuthMethodType = "AUTH_METHOD_TYPE_OAUTH_DISCORD"
	AuthMethodTypeOAuthApple   AuthMethodType = "AUTH_METHOD_TYPE_OAUTH_APPLE"
	AuthMethodTypePasskey      AuthMethodType = "AUTH_METHOD_TYPE_PASSKEY"
)

// AuthMethod represents v1AuthMethod.
type AuthMethod struct {
	MethodType  AuthMethodType `json:"methodType,omitempty"`
	DisplayName string         `json:"displayName,omitempty"`
	Preferred   bool           `json:"preferred,omitempty"`
}

// ListAuthMethodsRequest identifies the user to inspect.
type ListAuthMethodsRequest struct {
	UserID string
}

// ListAuthMethodsResponse lists auth methods for a user.
type ListAuthMethodsResponse struct {
	Methods  []AuthMethod     `json:"methods,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ListAuthMethodsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RemoveAuthMethodRequest removes an authentication method from a user.
type RemoveAuthMethodRequest struct {
	UserID     string
	MethodType AuthMethodType
}

// RemoveAuthMethodResponse captures metadata for removal operations.
type RemoveAuthMethodResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *RemoveAuthMethodResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// SetPreferredAuthMethodRequest marks an authentication method as preferred.
type SetPreferredAuthMethodRequest struct {
	UserID     string
	MethodType AuthMethodType
}

// SetPreferredAuthMethodResponse captures metadata for set-preferred operations.
type SetPreferredAuthMethodResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *SetPreferredAuthMethodResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteUserRequest identifies the user to delete.
type DeleteUserRequest struct {
	UserID string
}

// DeleteUserResponse captures metadata for delete operations.
type DeleteUserResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteUserResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

func setQueryInt32(values url.Values, key string, value *int32) {
	if value == nil || *value <= 0 {
		return
	}
	values.Set(key, strconv.Itoa(int(*value)))
}

// ListAPIKeys lists API keys for the authenticated user.
func (c *Client) ListAPIKeys(ctx context.Context, request *ListAPIKeysRequest, opts ...RequestOpt) (*ListAPIKeysResponse, error) {
	if request == nil {
		request = &ListAPIKeysRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/auth/api-keys", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryInt32(query, "pageSize", request.PageSize)
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	response := &ListAPIKeysResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CreateAPIKey creates a new API key for the authenticated user.
func (c *Client) CreateAPIKey(ctx context.Context, request *CreateAPIKeyRequest, opts ...RequestOpt) (*CreateAPIKeyResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/api-keys", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CreateAPIKeyResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ValidateAPIKey validates a provided API key.
func (c *Client) ValidateAPIKey(ctx context.Context, request *ValidateAPIKeyRequest, opts ...RequestOpt) (*ValidateAPIKeyResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.APIKey) == "" {
		return nil, errors.New("apiKey must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/api-keys/validate", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &ValidateAPIKeyResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// RevokeAPIKey revokes an API key by ID.
func (c *Client) RevokeAPIKey(ctx context.Context, request *RevokeAPIKeyRequest, opts ...RequestOpt) (*RevokeAPIKeyResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	apiKeyID := strings.TrimSpace(request.APIKeyID)
	if apiKeyID == "" {
		return nil, errors.New("apiKeyID must not be empty")
	}

	path := fmt.Sprintf("/v1/auth/api-keys/%s/revoke", url.PathEscape(apiKeyID))
	req, err := c.newRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	response := &RevokeAPIKeyResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ChangePassword changes a user's password.
func (c *Client) ChangePassword(ctx context.Context, request *ChangePasswordRequest, opts ...RequestOpt) (*ChangePasswordResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.UserID) == "" {
		return nil, errors.New("userID must not be empty")
	}
	if strings.TrimSpace(request.CurrentPassword) == "" {
		return nil, errors.New("currentPassword must not be empty")
	}
	if strings.TrimSpace(request.NewPassword) == "" {
		return nil, errors.New("newPassword must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/change-password", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &ChangePasswordResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// Register creates a new account using email/password.
func (c *Client) Register(ctx context.Context, request *RegisterRequest, opts ...RequestOpt) (*RegisterResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.Email) == "" {
		return nil, errors.New("email must not be empty")
	}
	if strings.TrimSpace(request.Password) == "" {
		return nil, errors.New("password must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/register", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &RegisterResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// RequestPasswordReset begins a password reset flow.
func (c *Client) RequestPasswordReset(ctx context.Context, request *RequestPasswordResetRequest, opts ...RequestOpt) (*RequestPasswordResetResponse, error) {
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

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/request-password-reset", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &RequestPasswordResetResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ResetPassword completes a password reset flow.
func (c *Client) ResetPassword(ctx context.Context, request *ResetPasswordRequest, opts ...RequestOpt) (*ResetPasswordResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ResetToken) == "" {
		return nil, errors.New("resetToken must not be empty")
	}
	if strings.TrimSpace(request.NewPassword) == "" {
		return nil, errors.New("newPassword must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/reset-password", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &ResetPasswordResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ResendVerificationEmail resends a verification email for a user.
func (c *Client) ResendVerificationEmail(ctx context.Context, request *ResendVerificationEmailRequest, opts ...RequestOpt) (*ResendVerificationEmailResponse, error) {
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

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/resend-verification-email", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &ResendVerificationEmailResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// VerifyEmail verifies a user's email address with a token.
func (c *Client) VerifyEmail(ctx context.Context, request *VerifyEmailRequest, opts ...RequestOpt) (*VerifyEmailResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.UserID) == "" {
		return nil, errors.New("userID must not be empty")
	}
	if strings.TrimSpace(request.VerificationToken) == "" {
		return nil, errors.New("verificationToken must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/verify-email", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &VerifyEmailResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ValidateToken validates an access token.
func (c *Client) ValidateToken(ctx context.Context, request *ValidateTokenRequest, opts ...RequestOpt) (*ValidateTokenResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.AccessToken) == "" {
		return nil, errors.New("accessToken must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/token/validate", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &ValidateTokenResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// InitiateOAuth starts an OAuth flow with the requested provider.
func (c *Client) InitiateOAuth(ctx context.Context, request *InitiateOAuthRequest, opts ...RequestOpt) (*InitiateOAuthResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	provider := strings.TrimSpace(request.Provider)
	if provider == "" {
		return nil, errors.New("provider must not be empty")
	}

	body, err := jsonBody(struct {
		RedirectURI string `json:"redirectUri,omitempty"`
		Prompt      string `json:"prompt,omitempty"`
	}{
		RedirectURI: strings.TrimSpace(request.RedirectURI),
		Prompt:      strings.TrimSpace(request.Prompt),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/oauth/%s/initiate", url.PathEscape(provider))
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

// CompleteOAuth completes an OAuth code exchange and returns login credentials.
func (c *Client) CompleteOAuth(ctx context.Context, request *CompleteOAuthRequest, opts ...RequestOpt) (*LoginResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	provider := strings.TrimSpace(request.Provider)
	if provider == "" {
		return nil, errors.New("provider must not be empty")
	}
	if strings.TrimSpace(request.Code) == "" {
		return nil, errors.New("code must not be empty")
	}

	body, err := jsonBody(struct {
		Code        string `json:"code,omitempty"`
		State       string `json:"state,omitempty"`
		RedirectURI string `json:"redirectUri,omitempty"`
	}{
		Code:        strings.TrimSpace(request.Code),
		State:       strings.TrimSpace(request.State),
		RedirectURI: strings.TrimSpace(request.RedirectURI),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/oauth/%s/complete", url.PathEscape(provider))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
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

// LinkOAuthProvider links a provider account to a user.
func (c *Client) LinkOAuthProvider(ctx context.Context, request *LinkOAuthProviderRequest, opts ...RequestOpt) (*LinkOAuthProviderResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	provider := strings.TrimSpace(request.Provider)
	if provider == "" {
		return nil, errors.New("provider must not be empty")
	}
	if strings.TrimSpace(request.UserID) == "" {
		return nil, errors.New("userID must not be empty")
	}
	if strings.TrimSpace(request.Code) == "" {
		return nil, errors.New("code must not be empty")
	}

	body, err := jsonBody(struct {
		UserID      string `json:"userId,omitempty"`
		Code        string `json:"code,omitempty"`
		State       string `json:"state,omitempty"`
		RedirectURI string `json:"redirectUri,omitempty"`
	}{
		UserID:      strings.TrimSpace(request.UserID),
		Code:        strings.TrimSpace(request.Code),
		State:       strings.TrimSpace(request.State),
		RedirectURI: strings.TrimSpace(request.RedirectURI),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/oauth/%s/link", url.PathEscape(provider))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &LinkOAuthProviderResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UnlinkOAuthProvider unlinks a provider account from a user.
func (c *Client) UnlinkOAuthProvider(ctx context.Context, request *UnlinkOAuthProviderRequest, opts ...RequestOpt) (*UnlinkOAuthProviderResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	provider := strings.TrimSpace(request.Provider)
	if provider == "" {
		return nil, errors.New("provider must not be empty")
	}
	if strings.TrimSpace(request.UserID) == "" {
		return nil, errors.New("userID must not be empty")
	}

	body, err := jsonBody(struct {
		UserID string `json:"userId,omitempty"`
	}{
		UserID: strings.TrimSpace(request.UserID),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/oauth/%s/unlink", url.PathEscape(provider))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UnlinkOAuthProviderResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListAuthMethods lists available authentication methods for a user.
func (c *Client) ListAuthMethods(ctx context.Context, request *ListAuthMethodsRequest, opts ...RequestOpt) (*ListAuthMethodsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/v1/auth/users/%s/methods", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &ListAuthMethodsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// RemoveAuthMethod removes an authentication method from a user.
func (c *Client) RemoveAuthMethod(ctx context.Context, request *RemoveAuthMethodRequest, opts ...RequestOpt) (*RemoveAuthMethodResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}
	methodType := strings.TrimSpace(string(request.MethodType))
	if methodType == "" || methodType == string(AuthMethodTypeUnspecified) {
		return nil, errors.New("methodType must not be empty")
	}

	path := fmt.Sprintf("/v1/auth/users/%s/methods/%s", url.PathEscape(userID), url.PathEscape(methodType))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &RemoveAuthMethodResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// SetPreferredAuthMethod marks an authentication method as preferred for a user.
func (c *Client) SetPreferredAuthMethod(ctx context.Context, request *SetPreferredAuthMethodRequest, opts ...RequestOpt) (*SetPreferredAuthMethodResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}
	methodType := strings.TrimSpace(string(request.MethodType))
	if methodType == "" || methodType == string(AuthMethodTypeUnspecified) {
		return nil, errors.New("methodType must not be empty")
	}

	path := fmt.Sprintf("/v1/auth/users/%s/methods/%s/preferred", url.PathEscape(userID), url.PathEscape(methodType))
	req, err := c.newRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	response := &SetPreferredAuthMethodResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// DeleteUser deletes a user by ID.
func (c *Client) DeleteUser(ctx context.Context, request *DeleteUserRequest, opts ...RequestOpt) (*DeleteUserResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/v1/users/%s", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteUserResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
