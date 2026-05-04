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
	Username             string `json:"username,omitempty"`
	Name                 string `json:"name,omitempty"`
	AvatarURL            string `json:"avatarUrl,omitempty"`
	Bio                  string `json:"bio,omitempty"`
	Location             string `json:"location,omitempty"`
	Pronouns             string `json:"pronouns,omitempty"`
	GemID                string `json:"gemId,omitempty"`
	ProfileColor         string `json:"profileColor,omitempty"`
	ProfileBackgroundURL string `json:"profileBackgroundUrl,omitempty"`
}

// UserSearchResult represents a compact search result for a user.
type UserSearchResult struct {
	UserID       string `json:"userId,omitempty"`
	Username     string `json:"username,omitempty"`
	Name         string `json:"name,omitempty"`
	AvatarURL    string `json:"avatarUrl,omitempty"`
	UserIDSuffix string `json:"userIdSuffix,omitempty"`
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

// BeginAvatarUploadRequest starts a server-controlled avatar upload.
type BeginAvatarUploadRequest struct {
	ContentType   string `json:"contentType,omitempty"`
	ContentLength *int64 `json:"contentLength,omitempty,string"`
}

// BeginAvatarUploadResponse returns upload instructions.
type BeginAvatarUploadResponse struct {
	UploadID        string            `json:"uploadId,omitempty"`
	UploadURL       string            `json:"uploadUrl,omitempty"`
	RequiredHeaders map[string]string `json:"requiredHeaders,omitempty"`
	MaxBytes        int64             `json:"maxBytes,omitempty,string"`
	ExpiresAt       *time.Time        `json:"expiresAt,omitempty"`
	Metadata        ResponseMetadata  `json:"-"`
}

func (r *BeginAvatarUploadResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CompleteAvatarUploadRequest completes an avatar upload.
type CompleteAvatarUploadRequest struct {
	UploadID string `json:"uploadId,omitempty"`
}

// CompleteAvatarUploadResponse returns the updated profile after applying the avatar.
type CompleteAvatarUploadResponse struct {
	Profile  *UserProfile     `json:"profile,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CompleteAvatarUploadResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BeginProfileBackgroundUploadRequest starts a server-controlled profile background upload.
type BeginProfileBackgroundUploadRequest struct {
	ContentType   string `json:"contentType,omitempty"`
	ContentLength *int64 `json:"contentLength,omitempty,string"`
}

// BeginProfileBackgroundUploadResponse returns upload instructions.
type BeginProfileBackgroundUploadResponse struct {
	UploadID        string            `json:"uploadId,omitempty"`
	UploadURL       string            `json:"uploadUrl,omitempty"`
	RequiredHeaders map[string]string `json:"requiredHeaders,omitempty"`
	MaxBytes        int64             `json:"maxBytes,omitempty,string"`
	ExpiresAt       *time.Time        `json:"expiresAt,omitempty"`
	Metadata        ResponseMetadata  `json:"-"`
}

func (r *BeginProfileBackgroundUploadResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CompleteProfileBackgroundUploadRequest completes a profile background upload.
type CompleteProfileBackgroundUploadRequest struct {
	UploadID string `json:"uploadId,omitempty"`
}

// CompleteProfileBackgroundUploadResponse returns the updated profile.
type CompleteProfileBackgroundUploadResponse struct {
	Profile  *UserProfile     `json:"profile,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CompleteProfileBackgroundUploadResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ClearProfileBackgroundResponse returns the updated profile after clearing the background.
type ClearProfileBackgroundResponse struct {
	Profile  *UserProfile     `json:"profile,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ClearProfileBackgroundResponse) setMetadata(metadata ResponseMetadata) {
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

// SearchUsersRequest searches for users by username prefix.
type SearchUsersRequest struct {
	SearchTerm string
	PageSize   *int32
	NextToken  string
}

// SearchUsersResponse lists users matching the search term.
type SearchUsersResponse struct {
	Users     []UserSearchResult `json:"users,omitempty"`
	NextToken string             `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata   `json:"-"`
}

func (r *SearchUsersResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// PrivacyConsent represents account-level privacy consent state.
type PrivacyConsent struct {
	AnalyticsAllowed bool       `json:"analyticsAllowed,omitempty"`
	ConsentVersion   int32      `json:"consentVersion,omitempty"`
	Source           string     `json:"source,omitempty"`
	Platform         string     `json:"platform,omitempty"`
	AppVersion       string     `json:"appVersion,omitempty"`
	DeviceIDHash     string     `json:"deviceIdHash,omitempty"`
	UpdatedAt        *time.Time `json:"updatedAt,omitempty"`
	ClientActionAt   *time.Time `json:"clientActionAt,omitempty"`
}

// GetPrivacyConsentResponse returns account-level privacy consent.
type GetPrivacyConsentResponse struct {
	Consent  *PrivacyConsent  `json:"consent,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetPrivacyConsentResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdatePrivacyConsentRequest updates account-level privacy consent.
type UpdatePrivacyConsentRequest struct {
	AnalyticsAllowed bool       `json:"analyticsAllowed,omitempty"`
	ConsentVersion   int32      `json:"consentVersion,omitempty"`
	Source           string     `json:"source,omitempty"`
	Platform         string     `json:"platform,omitempty"`
	AppVersion       string     `json:"appVersion,omitempty"`
	DeviceIDHash     string     `json:"deviceIdHash,omitempty"`
	ClientActionAt   *time.Time `json:"clientActionAt,omitempty"`
}

// UpdatePrivacyConsentResponse returns updated account-level privacy consent.
type UpdatePrivacyConsentResponse struct {
	Consent  *PrivacyConsent  `json:"consent,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdatePrivacyConsentResponse) setMetadata(metadata ResponseMetadata) {
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

// PinCollectionRequest pins a collection for the current user.
type PinCollectionRequest struct {
	CollectionID string `json:"collectionId,omitempty"`
}

// PinCollectionResponse captures metadata for pin operations.
type PinCollectionResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *PinCollectionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnpinCollectionRequest removes a pinned collection for the current user.
type UnpinCollectionRequest struct {
	CollectionID string `json:"-"`
}

// UnpinCollectionResponse captures metadata for unpin operations.
type UnpinCollectionResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnpinCollectionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// PinDeckRequest pins a deck for the current user.
type PinDeckRequest struct {
	DeckID string `json:"deckId,omitempty"`
}

// PinDeckResponse captures metadata for pin operations.
type PinDeckResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *PinDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnpinDeckRequest removes a pinned deck for the current user.
type UnpinDeckRequest struct {
	DeckID string `json:"-"`
}

// UnpinDeckResponse captures metadata for unpin operations.
type UnpinDeckResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnpinDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// PinnedCollection captures pinned collection metadata.
type PinnedCollection struct {
	CollectionID string     `json:"collectionId,omitempty"`
	PinnedAt     *time.Time `json:"pinnedAt,omitempty"`
}

// PinnedDeck captures pinned deck metadata.
type PinnedDeck struct {
	DeckID   string     `json:"deckId,omitempty"`
	PinnedAt *time.Time `json:"pinnedAt,omitempty"`
}

// GetPinnedResourcesRequest retrieves pinned collections and decks for a user.
type GetPinnedResourcesRequest struct {
	UserID string
}

// GetPinnedResourcesResponse returns pinned collections and decks for a user.
type GetPinnedResourcesResponse struct {
	PinnedCollections []PinnedCollection `json:"pinnedCollections,omitempty"`
	PinnedDecks       []PinnedDeck       `json:"pinnedDecks,omitempty"`
	Metadata          ResponseMetadata   `json:"-"`
}

func (r *GetPinnedResourcesResponse) setMetadata(metadata ResponseMetadata) {
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

// GetPinnedResources retrieves pinned collections/decks for a user.
func (c *Client) GetPinnedResources(ctx context.Context, request *GetPinnedResourcesRequest, opts ...RequestOpt) (*GetPinnedResourcesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/v1/users/%s/pins", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetPinnedResourcesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// PinCollection pins a collection for the current user.
func (c *Client) PinCollection(ctx context.Context, request *PinCollectionRequest, opts ...RequestOpt) (*PinCollectionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	collectionID := strings.TrimSpace(request.CollectionID)
	if collectionID == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	body, err := jsonBody(struct {
		CollectionID string `json:"collectionId,omitempty"`
	}{
		CollectionID: collectionID,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/me/pins/collections", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &PinCollectionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UnpinCollection removes a pinned collection for the current user.
func (c *Client) UnpinCollection(ctx context.Context, request *UnpinCollectionRequest, opts ...RequestOpt) (*UnpinCollectionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	collectionID := strings.TrimSpace(request.CollectionID)
	if collectionID == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	path := fmt.Sprintf("/v1/me/pins/collections/%s", url.PathEscape(collectionID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &UnpinCollectionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// PinDeck pins a deck for the current user.
func (c *Client) PinDeck(ctx context.Context, request *PinDeckRequest, opts ...RequestOpt) (*PinDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	body, err := jsonBody(struct {
		DeckID string `json:"deckId,omitempty"`
	}{
		DeckID: deckID,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/me/pins/decks", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &PinDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UnpinDeck removes a pinned deck for the current user.
func (c *Client) UnpinDeck(ctx context.Context, request *UnpinDeckRequest, opts ...RequestOpt) (*UnpinDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/me/pins/decks/%s", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &UnpinDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetMyProfile reads the authenticated user's profile.
func (c *Client) GetMyProfile(ctx context.Context, opts ...RequestOpt) (*GetProfileResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/me/profile", nil)
	if err != nil {
		return nil, err
	}

	response := &GetProfileResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// SearchUsers searches for users by username prefix.
func (c *Client) SearchUsers(ctx context.Context, request *SearchUsersRequest, opts ...RequestOpt) (*SearchUsersResponse, error) {
	if request == nil {
		request = &SearchUsersRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/users/search", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if term := strings.TrimSpace(request.SearchTerm); term != "" {
		query.Set("searchTerm", term)
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	response := &SearchUsersResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
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

	req, err := c.newRequest(ctx, http.MethodPut, "/v1/me/avatar", body)
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

// SetAvatarUrl is an alias for SetAvatarURL.
func (c *Client) SetAvatarUrl(ctx context.Context, request *SetAvatarURLRequest, opts ...RequestOpt) (*SetAvatarURLResponse, error) {
	return c.SetAvatarURL(ctx, request, opts...)
}

// BeginAvatarUpload starts a server-controlled avatar upload session.
func (c *Client) BeginAvatarUpload(ctx context.Context, request *BeginAvatarUploadRequest, opts ...RequestOpt) (*BeginAvatarUploadResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ContentType) == "" {
		return nil, errors.New("contentType must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/me/avatar:beginUpload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BeginAvatarUploadResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CompleteAvatarUpload completes an avatar upload session and applies it to the current user.
func (c *Client) CompleteAvatarUpload(ctx context.Context, request *CompleteAvatarUploadRequest, opts ...RequestOpt) (*CompleteAvatarUploadResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.UploadID) == "" {
		return nil, errors.New("uploadID must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/me/avatar:completeUpload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CompleteAvatarUploadResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BeginProfileBackgroundUpload starts a server-controlled profile background upload session.
func (c *Client) BeginProfileBackgroundUpload(ctx context.Context, request *BeginProfileBackgroundUploadRequest, opts ...RequestOpt) (*BeginProfileBackgroundUploadResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ContentType) == "" {
		return nil, errors.New("contentType must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/me/profile/background:beginUpload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BeginProfileBackgroundUploadResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CompleteProfileBackgroundUpload completes a profile background upload session.
func (c *Client) CompleteProfileBackgroundUpload(ctx context.Context, request *CompleteProfileBackgroundUploadRequest, opts ...RequestOpt) (*CompleteProfileBackgroundUploadResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.UploadID) == "" {
		return nil, errors.New("uploadID must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/me/profile/background:completeUpload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CompleteProfileBackgroundUploadResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ClearProfileBackground clears the current user's profile background image.
func (c *Client) ClearProfileBackground(ctx context.Context, opts ...RequestOpt) (*ClearProfileBackgroundResponse, error) {
	body, err := jsonBody(struct{}{})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/me/profile/background:clear", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &ClearProfileBackgroundResponse{}
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

	req, err := c.newRequest(ctx, http.MethodPatch, "/v1/me/profile", body)
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

	path := fmt.Sprintf("/v1/users/%s/profile", url.PathEscape(userID))
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

// GetPrivacyConsent returns account-level privacy consent for the authenticated user.
func (c *Client) GetPrivacyConsent(ctx context.Context, opts ...RequestOpt) (*GetPrivacyConsentResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/me/privacy/consent", nil)
	if err != nil {
		return nil, err
	}

	response := &GetPrivacyConsentResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpdatePrivacyConsent updates account-level privacy consent for the authenticated user.
func (c *Client) UpdatePrivacyConsent(ctx context.Context, request *UpdatePrivacyConsentRequest, opts ...RequestOpt) (*UpdatePrivacyConsentResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPut, "/v1/me/privacy/consent", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdatePrivacyConsentResponse{}
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

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/me/profile/settings", nil)
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

	req, err := c.newRequest(ctx, http.MethodPut, "/v1/me/profile/settings", body)
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

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/me/social_profiles", body)
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

	path := fmt.Sprintf("/v1/me/social_profiles/%s", url.PathEscape(platform))
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

	path := fmt.Sprintf("/v1/users/%s/social_profiles", url.PathEscape(userID))
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
