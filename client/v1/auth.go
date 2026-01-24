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

// LoginRequest carries credentials used to authenticate a user.
type LoginRequest struct {
	Email      string `json:"email,omitempty"`
	Password   string `json:"password,omitempty"`
	DeviceInfo string `json:"deviceInfo,omitempty"`
}

// LoginResponse contains tokens issued during authentication.
type LoginResponse struct {
	UserID               string           `json:"userId,omitempty"`
	AccessToken          string           `json:"accessToken,omitempty"`
	RefreshToken         string           `json:"refreshToken,omitempty"`
	AccessTokenExpiresAt *time.Time       `json:"accessTokenExpiresAt,omitempty"`
	Roles                []string         `json:"roles,omitempty"`
	Metadata             ResponseMetadata `json:"-"`
}

func (r *LoginResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CLILoginSessionStatus mirrors v1CLILoginSessionStatus.
type CLILoginSessionStatus string

const (
	CLILoginSessionStatusUnspecified CLILoginSessionStatus = "CLI_LOGIN_SESSION_STATUS_UNSPECIFIED"
	CLILoginSessionStatusPending     CLILoginSessionStatus = "CLI_LOGIN_SESSION_STATUS_PENDING"
	CLILoginSessionStatusComplete    CLILoginSessionStatus = "CLI_LOGIN_SESSION_STATUS_COMPLETE"
	CLILoginSessionStatusExpired     CLILoginSessionStatus = "CLI_LOGIN_SESSION_STATUS_EXPIRED"
	CLILoginSessionStatusCanceled    CLILoginSessionStatus = "CLI_LOGIN_SESSION_STATUS_CANCELED"
)

// CreateCLILoginSessionRequest creates a new browser-based OAuth login session for CLI clients.
type CreateCLILoginSessionRequest struct {
	BaseURL string `json:"baseUrl,omitempty"`
}

// CreateCLILoginSessionResponse returns the session secret and verification URL.
type CreateCLILoginSessionResponse struct {
	SessionID           string           `json:"sessionId,omitempty"`
	SessionSecret       string           `json:"sessionSecret,omitempty"`
	VerificationPath    string           `json:"verificationPath,omitempty"`
	VerificationURL     string           `json:"verificationUrl,omitempty"`
	ExpiresAt           *time.Time       `json:"expiresAt,omitempty"`
	PollIntervalSeconds int32            `json:"pollIntervalSeconds,omitempty"`
	Metadata            ResponseMetadata `json:"-"`
}

func (r *CreateCLILoginSessionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetCLILoginSessionRequest polls a login session for completion.
type GetCLILoginSessionRequest struct {
	SessionID     string `json:"-"`
	SessionSecret string `json:"sessionSecret,omitempty"`
}

// GetCLILoginSessionResponse includes the current session status and login credentials when complete.
type GetCLILoginSessionResponse struct {
	Status   CLILoginSessionStatus `json:"status,omitempty"`
	Login    *LoginResponse        `json:"login,omitempty"`
	Metadata ResponseMetadata      `json:"-"`
}

func (r *GetCLILoginSessionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CancelCLILoginSessionRequest cancels an in-progress login session.
type CancelCLILoginSessionRequest struct {
	SessionID     string `json:"-"`
	SessionSecret string `json:"sessionSecret,omitempty"`
}

// CancelCLILoginSessionResponse captures metadata for a cancellation operation.
type CancelCLILoginSessionResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *CancelCLILoginSessionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CompleteOAuthForCLILoginSessionRequest completes an OAuth flow for a specific CLI login session.
type CompleteOAuthForCLILoginSessionRequest struct {
	SessionID   string `json:"-"`
	Provider    string `json:"-"`
	Code        string `json:"code,omitempty"`
	State       string `json:"state,omitempty"`
	RedirectURI string `json:"redirectUri,omitempty"`
}

// CompleteOAuthForCLILoginSessionResponse captures metadata for a completion callback.
type CompleteOAuthForCLILoginSessionResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *CompleteOAuthForCLILoginSessionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RefreshTokenRequest exchanges a refresh token for new credentials.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken,omitempty"`
}

// RefreshTokenResponse returns renewed access and refresh tokens.
type RefreshTokenResponse struct {
	AccessToken          string           `json:"accessToken,omitempty"`
	RefreshToken         string           `json:"refreshToken,omitempty"`
	AccessTokenExpiresAt *time.Time       `json:"accessTokenExpiresAt,omitempty"`
	Metadata             ResponseMetadata `json:"-"`
}

func (r *RefreshTokenResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// LogoutRequest revokes authentication state.
type LogoutRequest struct {
	RefreshToken string `json:"refreshToken,omitempty"`
}

// LogoutResponse captures metadata from a logout operation.
type LogoutResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *LogoutResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// MeResponse returns information about the current user.
type MeResponse struct {
	User     *User            `json:"user,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *MeResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetUserRequest identifies a user by ID.
type GetUserRequest struct {
	UserID string `json:"-"`
}

// GetUserResponse returns a user.
type GetUserResponse struct {
	User     *User            `json:"user,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetUserResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateUserRequest updates mutable fields for a user.
type UpdateUserRequest struct {
	UserID string   `json:"-"`
	Email  *string  `json:"email,omitempty"`
	Roles  []string `json:"roles,omitempty"`
}

// UpdateUserResponse returns the updated user details.
type UpdateUserResponse struct {
	User     *User            `json:"user,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateUserResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteUserRequest identifies a user to delete.
type DeleteUserRequest struct {
	UserID string `json:"-"`
}

// DeleteUserResponse captures metadata for delete operations.
type DeleteUserResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteUserResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// User represents an account within Pitchstack.
type User struct {
	UserID           string     `json:"userId,omitempty"`
	Email            string     `json:"email,omitempty"`
	EmailVerified    bool       `json:"emailVerified,omitempty"`
	Phone            string     `json:"phone,omitempty"`
	PhoneVerified    bool       `json:"phoneVerified,omitempty"`
	IsActive         bool       `json:"isActive,omitempty"`
	IsSuspended      bool       `json:"isSuspended,omitempty"`
	SuspendedUntil   *time.Time `json:"suspendedUntil,omitempty"`
	SuspensionReason string     `json:"suspensionReason,omitempty"`
	Roles            []string   `json:"roles,omitempty"`
	CreatedAt        *time.Time `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
	LastLoginAt      *time.Time `json:"lastLoginAt,omitempty"`
}

// Login authenticates a user using supplied credentials.
func (c *Client) Login(ctx context.Context, request *LoginRequest, opts ...RequestOpt) (*LoginResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/login", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &LoginResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateCLILoginSession starts a CLI login session to be completed via OAuth in a browser.
func (c *Client) CreateCLILoginSession(ctx context.Context, request *CreateCLILoginSessionRequest, opts ...RequestOpt) (*CreateCLILoginSessionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/cli/sessions", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CreateCLILoginSessionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetCLILoginSession polls the CLI login session for completion.
func (c *Client) GetCLILoginSession(ctx context.Context, request *GetCLILoginSessionRequest, opts ...RequestOpt) (*GetCLILoginSessionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	sessionID := strings.TrimSpace(request.SessionID)
	if sessionID == "" {
		return nil, errors.New("sessionID must not be empty")
	}
	if strings.TrimSpace(request.SessionSecret) == "" {
		return nil, errors.New("sessionSecret must not be empty")
	}

	body, err := jsonBody(struct {
		SessionSecret string `json:"sessionSecret,omitempty"`
	}{
		SessionSecret: request.SessionSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/cli/sessions/%s:poll", url.PathEscape(sessionID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &GetCLILoginSessionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CancelCLILoginSession cancels an in-progress CLI login session.
func (c *Client) CancelCLILoginSession(ctx context.Context, request *CancelCLILoginSessionRequest, opts ...RequestOpt) (*CancelCLILoginSessionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	sessionID := strings.TrimSpace(request.SessionID)
	if sessionID == "" {
		return nil, errors.New("sessionID must not be empty")
	}
	if strings.TrimSpace(request.SessionSecret) == "" {
		return nil, errors.New("sessionSecret must not be empty")
	}

	body, err := jsonBody(struct {
		SessionSecret string `json:"sessionSecret,omitempty"`
	}{
		SessionSecret: request.SessionSecret,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/cli/sessions/%s:cancel", url.PathEscape(sessionID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CancelCLILoginSessionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CompleteOAuthForCLILoginSession completes OAuth for a CLI login session.
func (c *Client) CompleteOAuthForCLILoginSession(ctx context.Context, request *CompleteOAuthForCLILoginSessionRequest, opts ...RequestOpt) (*CompleteOAuthForCLILoginSessionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	sessionID := strings.TrimSpace(request.SessionID)
	if sessionID == "" {
		return nil, errors.New("sessionID must not be empty")
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
		Code:        request.Code,
		State:       request.State,
		RedirectURI: request.RedirectURI,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/auth/cli/sessions/%s/oauth/%s:complete", url.PathEscape(sessionID), url.PathEscape(provider))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CompleteOAuthForCLILoginSessionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// RefreshToken exchanges a refresh token for a fresh access token.
func (c *Client) RefreshToken(ctx context.Context, request *RefreshTokenRequest, opts ...RequestOpt) (*RefreshTokenResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.RefreshToken) == "" {
		return nil, errors.New("refreshToken must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/token/refresh", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &RefreshTokenResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// Logout invalidates the current session.
func (c *Client) Logout(ctx context.Context, request *LogoutRequest, opts ...RequestOpt) (*LogoutResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/auth/logout", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &LogoutResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// Me retrieves the current authenticated user.
func (c *Client) Me(ctx context.Context, opts ...RequestOpt) (*MeResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/me", nil)
	if err != nil {
		return nil, err
	}

	response := &MeResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetUser retrieves a user by ID.
func (c *Client) GetUser(ctx context.Context, request *GetUserRequest, opts ...RequestOpt) (*GetUserResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/v1/users/%s", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetUserResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateUser applies changes to user attributes.
func (c *Client) UpdateUser(ctx context.Context, request *UpdateUserRequest, opts ...RequestOpt) (*UpdateUserResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	body, err := jsonBody(struct {
		Email *string  `json:"email,omitempty"`
		Roles []string `json:"roles,omitempty"`
	}{
		Email: request.Email,
		Roles: request.Roles,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/users/%s", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateUserResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteUser removes a user by ID.
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
