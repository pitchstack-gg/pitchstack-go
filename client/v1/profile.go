package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ProfileVisibilityLevel models profilev1VisibilityLevel.
type ProfileVisibilityLevel string

const (
	ProfileVisibilityLevelUnspecified ProfileVisibilityLevel = "VISIBILITY_LEVEL_UNSPECIFIED"
	ProfileVisibilityLevelPrivate     ProfileVisibilityLevel = "VISIBILITY_LEVEL_PRIVATE"
	ProfileVisibilityLevelFollowers   ProfileVisibilityLevel = "VISIBILITY_LEVEL_FOLLOWERS"
	ProfileVisibilityLevelPublic      ProfileVisibilityLevel = "VISIBILITY_LEVEL_PUBLIC"
)

// UserProfile represents a user's public profile details.
type UserProfile struct {
	Username  string `json:"username,omitempty"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Location  string `json:"location,omitempty"`
	Pronouns  string `json:"pronouns,omitempty"`
}

// ProfileSettings controls profile visibility preferences.
type ProfileSettings struct {
	ProfileVisibility        ProfileVisibilityLevel `json:"profileVisibility,omitempty"`
	SocialProfilesVisibility ProfileVisibilityLevel `json:"socialProfilesVisibility,omitempty"`
	AllowMessages            string                 `json:"allowMessages,omitempty"`
	AllowTradeOffers         string                 `json:"allowTradeOffers,omitempty"`
}

// SocialProfile represents a linked social account.
type SocialProfile struct {
	Platform string `json:"platform,omitempty"`
	Handle   string `json:"handle,omitempty"`
	URL      string `json:"url,omitempty"`
}

// SetAvatarURLRequest updates the authenticated user's avatar URL.
type SetAvatarURLRequest struct {
	AvatarURL string
}

// SetAvatarURLResponse captures metadata for avatar updates.
type SetAvatarURLResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *SetAvatarURLResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateProfileRequest modifies the authenticated user's profile.
type UpdateProfileRequest struct {
	Profile    *UserProfile
	UpdateMask string
}

// UpdateProfileResponse returns the updated profile.
type UpdateProfileResponse struct {
	Profile  *UserProfile     `json:"profile,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateProfileResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetProfileRequest fetches profile details for a user.
type GetProfileRequest struct {
	UserID string
}

// GetProfileResponse returns a user's profile.
type GetProfileResponse struct {
	Profile  *UserProfile     `json:"profile,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetProfileResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetProfileSettingsRequest retrieves profile settings for the caller.
type GetProfileSettingsRequest struct {
	UserID string
}

// GetProfileSettingsResponse returns profile settings.
type GetProfileSettingsResponse struct {
	Settings *ProfileSettings `json:"settings,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetProfileSettingsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateProfileSettingsRequest updates profile visibility settings.
type UpdateProfileSettingsRequest struct {
	Settings   *ProfileSettings
	UpdateMask string
}

// UpdateProfileSettingsResponse returns updated settings.
type UpdateProfileSettingsResponse struct {
	Settings *ProfileSettings `json:"settings,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateProfileSettingsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpsertSocialProfileRequest upserts a social profile for the current user.
type UpsertSocialProfileRequest struct {
	Platform string
	Handle   string
	URL      string
}

// UpsertSocialProfileResponse returns the upserted profile.
type UpsertSocialProfileResponse struct {
	SocialProfile *SocialProfile   `json:"socialProfile,omitempty"`
	Metadata      ResponseMetadata `json:"-"`
}

func (r *UpsertSocialProfileResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RemoveSocialProfileRequest removes a social profile by platform.
type RemoveSocialProfileRequest struct {
	Platform string
}

// RemoveSocialProfileResponse captures metadata for removal operations.
type RemoveSocialProfileResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *RemoveSocialProfileResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetSocialProfilesRequest retrieves social profiles for a user.
type GetSocialProfilesRequest struct {
	UserID string
}

// GetSocialProfilesResponse lists social profiles tied to a user.
type GetSocialProfilesResponse struct {
	SocialProfiles []SocialProfile  `json:"socialProfiles,omitempty"`
	Metadata       ResponseMetadata `json:"-"`
}

func (r *GetSocialProfilesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// SetAvatarURL updates the authenticated user's avatar URL.
func (c *Client) SetAvatarURL(ctx context.Context, request *SetAvatarURLRequest, opts ...RequestOpt) (*SetAvatarURLResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	avatar := strings.TrimSpace(request.AvatarURL)
	if avatar == "" {
		return nil, errors.New("avatarURL must not be empty")
	}

	body, err := jsonBody(struct {
		AvatarURL string `json:"avatarUrl,omitempty"`
	}{
		AvatarURL: avatar,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPut, "/api/v1/me/avatar", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &SetAvatarURLResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpdateProfile applies updates to the authenticated user's profile.
func (c *Client) UpdateProfile(ctx context.Context, request *UpdateProfileRequest, opts ...RequestOpt) (*UpdateProfileResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if request.Profile == nil {
		return nil, errors.New("profile must not be nil")
	}

	body, err := jsonBody(struct {
		Profile    *UserProfile `json:"profile,omitempty"`
		UpdateMask string       `json:"updateMask,omitempty"`
	}{
		Profile:    request.Profile,
		UpdateMask: strings.TrimSpace(request.UpdateMask),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPatch, "/api/v1/me/profile", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdateProfileResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetProfile retrieves a user's profile.
func (c *Client) GetProfile(ctx context.Context, request *GetProfileRequest, opts ...RequestOpt) (*GetProfileResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/api/v1/users/%s/profile", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetProfileResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetProfileSettings retrieves profile settings for the authenticated user.
func (c *Client) GetProfileSettings(ctx context.Context, request *GetProfileSettingsRequest, opts ...RequestOpt) (*GetProfileSettingsResponse, error) {
	if request == nil {
		request = &GetProfileSettingsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/me/profile/settings", nil)
	if err != nil {
		return nil, err
	}

	if userID := strings.TrimSpace(request.UserID); userID != "" {
		query := req.URL.Query()
		query.Set("userId", userID)
		req.URL.RawQuery = query.Encode()
	}

	response := &GetProfileSettingsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpdateProfileSettings updates profile settings for the authenticated user.
func (c *Client) UpdateProfileSettings(ctx context.Context, request *UpdateProfileSettingsRequest, opts ...RequestOpt) (*UpdateProfileSettingsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if request.Settings == nil {
		return nil, errors.New("settings must not be nil")
	}

	body, err := jsonBody(struct {
		Settings   *ProfileSettings `json:"settings,omitempty"`
		UpdateMask string           `json:"updateMask,omitempty"`
	}{
		Settings:   request.Settings,
		UpdateMask: strings.TrimSpace(request.UpdateMask),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPut, "/api/v1/me/profile/settings", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdateProfileSettingsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpsertSocialProfile creates or updates a social profile for the authenticated user.
func (c *Client) UpsertSocialProfile(ctx context.Context, request *UpsertSocialProfileRequest, opts ...RequestOpt) (*UpsertSocialProfileResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	platform := strings.TrimSpace(request.Platform)
	if platform == "" {
		return nil, errors.New("platform must not be empty")
	}

	body, err := jsonBody(struct {
		Platform string `json:"platform,omitempty"`
		Handle   string `json:"handle,omitempty"`
		URL      string `json:"url,omitempty"`
	}{
		Platform: platform,
		Handle:   strings.TrimSpace(request.Handle),
		URL:      strings.TrimSpace(request.URL),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/me/social_profiles", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpsertSocialProfileResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// RemoveSocialProfile deletes a social profile for the authenticated user.
func (c *Client) RemoveSocialProfile(ctx context.Context, request *RemoveSocialProfileRequest, opts ...RequestOpt) (*RemoveSocialProfileResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	platform := strings.TrimSpace(request.Platform)
	if platform == "" {
		return nil, errors.New("platform must not be empty")
	}

	path := fmt.Sprintf("/api/v1/me/social_profiles/%s", url.PathEscape(platform))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &RemoveSocialProfileResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetSocialProfiles retrieves social profiles for the specified user.
func (c *Client) GetSocialProfiles(ctx context.Context, request *GetSocialProfilesRequest, opts ...RequestOpt) (*GetSocialProfilesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/api/v1/users/%s/social_profiles", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetSocialProfilesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
