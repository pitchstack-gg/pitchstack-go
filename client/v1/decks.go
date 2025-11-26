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

// DeckListScope matches v1DeckListScope.
type DeckListScope string

const (
	DeckListScopeUnspecified DeckListScope = "DECK_LIST_SCOPE_UNSPECIFIED"
	DeckListScopeOwned       DeckListScope = "DECK_LIST_SCOPE_OWNED"
	DeckListScopeShared      DeckListScope = "DECK_LIST_SCOPE_SHARED"
	DeckListScopeAccessible  DeckListScope = "DECK_LIST_SCOPE_ACCESSIBLE"
)

// BoardType represents v1BoardType.
type BoardType string

const (
	BoardTypeUnspecified BoardType = "BOARD_TYPE_UNSPECIFIED"
	BoardTypeMainboard   BoardType = "BOARD_TYPE_MAINBOARD"
	BoardTypeSideboard   BoardType = "BOARD_TYPE_SIDEBOARD"
	BoardTypeMaybeboard  BoardType = "BOARD_TYPE_MAYBEBOARD"
)

// Deck models v1Deck.
type Deck struct {
	ID          string          `json:"id,omitempty"`
	UserID      string          `json:"userId,omitempty"`
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	HeroID      string          `json:"heroId,omitempty"`
	Format      string          `json:"format,omitempty"`
	Versions    []string        `json:"versions,omitempty"`
	Visibility  VisibilityLevel `json:"visibility,omitempty"`
	CreatedAt   *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time      `json:"updatedAt,omitempty"`
}

// DeckVersion represents v1DeckVersion.
type DeckVersion struct {
	DeckID      string     `json:"deckId,omitempty"`
	Version     string     `json:"version,omitempty"`
	ImageURL    string     `json:"imageUrl,omitempty"`
	Description string     `json:"description,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// DeckVersionChange mirrors v1DeckVersionChange.
type DeckVersionChange struct {
	DeckID      string     `json:"deckId,omitempty"`
	Version     string     `json:"version,omitempty"`
	Description string     `json:"description,omitempty"`
	Timestamp   *time.Time `json:"timestamp,omitempty"`
}

// DeckCard mirrors v1Card used within deck operations.
type DeckCard struct {
	CardID   string `json:"cardId,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
}

// ListDecksRequest configures filters for ListDecks.
type ListDecksRequest struct {
	Scope     DeckListScope
	UserID    string
	HeroID    string
	Format    string
	Name      string
	PageSize  *int32
	NextToken string
}

// ListDecksResponse contains decks visible to the caller.
type ListDecksResponse struct {
	Decks     []Deck           `json:"decks,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListDecksResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateDeckRequest captures payload fields for deck creation.
type CreateDeckRequest struct {
	Name       string          `json:"name,omitempty"`
	HeroID     string          `json:"heroId,omitempty"`
	Format     string          `json:"format,omitempty"`
	Visibility VisibilityLevel `json:"visibility,omitempty"`
}

// CreateDeckResponse returns the created deck.
type CreateDeckResponse struct {
	Deck     *Deck            `json:"deck,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CreateDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// SearchDecksRequest configures deck search parameters.
type SearchDecksRequest struct {
	HeroID     string
	Format     string
	SearchTerm string
	PageSize   *int32
	NextToken  string
}

// SearchDecksResponse returns decks matching the search query.
type SearchDecksResponse struct {
	Decks     []Deck           `json:"decks,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *SearchDecksResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetDeckRequest identifies the deck to retrieve.
type GetDeckRequest struct {
	DeckID string `json:"-"`
}

// GetDeckResponse wraps the retrieved deck.
type GetDeckResponse struct {
	Deck     *Deck            `json:"deck,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteDeckRequest identifies the deck to delete.
type DeleteDeckRequest struct {
	DeckID string `json:"-"`
}

// DeleteDeckResponse captures metadata for delete calls.
type DeleteDeckResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateDeckRequest applies changes to a deck.
type UpdateDeckRequest struct {
	DeckID     string           `json:"-"`
	Name       *string          `json:"name,omitempty"`
	Visibility *VisibilityLevel `json:"visibility,omitempty"`
}

// UpdateDeckResponse contains the updated deck.
type UpdateDeckResponse struct {
	Deck     *Deck            `json:"deck,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// StarDeckRequest marks a deck as starred.
type StarDeckRequest struct {
	DeckID string `json:"-"`
}

// StarDeckResponse captures metadata for star operations.
type StarDeckResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *StarDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnstarDeckRequest removes a star from a deck.
type UnstarDeckRequest struct {
	DeckID string `json:"-"`
}

// UnstarDeckResponse captures metadata for unstar operations.
type UnstarDeckResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnstarDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListDeckVersionsRequest enumerates versions for a deck.
type ListDeckVersionsRequest struct {
	DeckID    string
	PageSize  *int32
	NextToken string
}

// ListDeckVersionsResponse lists deck versions.
type ListDeckVersionsResponse struct {
	Versions  []string         `json:"versions,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListDeckVersionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateDeckVersionRequest provisions a deck version.
type CreateDeckVersionRequest struct {
	DeckID      string `json:"-"`
	Version     string `json:"version,omitempty"`
	ImageURL    string `json:"imageUrl,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateDeckVersionResponse includes the created version.
type CreateDeckVersionResponse struct {
	DeckVersion *DeckVersion     `json:"deckVersion,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *CreateDeckVersionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetDeckVersionRequest identifies a deck version to retrieve.
type GetDeckVersionRequest struct {
	DeckID    string `json:"-"`
	Version   string `json:"-"`
	NextToken string
}

// GetDeckVersionResponse returns a deck version.
type GetDeckVersionResponse struct {
	DeckVersion *DeckVersion     `json:"deckVersion,omitempty"`
	NextToken   string           `json:"nextToken,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *GetDeckVersionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteDeckVersionRequest removes a deck version.
type DeleteDeckVersionRequest struct {
	DeckID  string `json:"-"`
	Version string `json:"-"`
}

// DeleteDeckVersionResponse captures metadata for deletions.
type DeleteDeckVersionResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteDeckVersionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateDeckVersionRequest updates attributes on a deck version.
type UpdateDeckVersionRequest struct {
	DeckID      string  `json:"-"`
	Version     string  `json:"-"`
	ImageURL    *string `json:"imageUrl,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateDeckVersionResponse returns the updated version.
type UpdateDeckVersionResponse struct {
	DeckVersion *DeckVersion     `json:"deckVersion,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *UpdateDeckVersionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListDeckVersionCardsRequest lists cards within a deck version.
type ListDeckVersionCardsRequest struct {
	DeckID  string
	Version string
}

// ListDeckVersionCardsResponse enumerates cards within a deck version.
type ListDeckVersionCardsResponse struct {
	MainboardCards  []DeckCard       `json:"mainboardCards,omitempty"`
	SideboardCards  []DeckCard       `json:"sideboardCards,omitempty"`
	MaybeboardCards []DeckCard       `json:"maybeboardCards,omitempty"`
	Metadata        ResponseMetadata `json:"-"`
}

func (r *ListDeckVersionCardsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ModifyDeckVersionCardRequest adjusts a card within a deck version.
type ModifyDeckVersionCardRequest struct {
	DeckID   string    `json:"-"`
	Version  string    `json:"-"`
	CardID   string    `json:"cardId,omitempty"`
	Board    BoardType `json:"boardType,omitempty"`
	Quantity int32     `json:"quantity,omitempty"`
}

// ModifyDeckVersionCardResponse returns the affected card.
type ModifyDeckVersionCardResponse struct {
	Card     *DeckCard        `json:"card,omitempty"`
	Board    BoardType        `json:"boardType,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ModifyDeckVersionCardResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetDeckVersionHistoryRequest retrieves version history for a deck.
type GetDeckVersionHistoryRequest struct {
	DeckID  string
	Version string
}

// GetDeckVersionHistoryResponse contains version changes.
type GetDeckVersionHistoryResponse struct {
	Changes  []DeckVersionChange `json:"changes,omitempty"`
	Metadata ResponseMetadata    `json:"-"`
}

func (r *GetDeckVersionHistoryResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetDeckVersionNotesRequest fetches notes for a deck version.
type GetDeckVersionNotesRequest struct {
	DeckID  string
	Version string
}

// GetDeckVersionNotesResponse returns notes text.
type GetDeckVersionNotesResponse struct {
	Notes    string           `json:"notes,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetDeckVersionNotesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateDeckVersionNotesRequest updates notes for a deck version.
type UpdateDeckVersionNotesRequest struct {
	DeckID  string
	Version string
	Notes   string
}

// UpdateDeckVersionNotesResponse returns updated notes.
type UpdateDeckVersionNotesResponse struct {
	Notes    string           `json:"notes,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateDeckVersionNotesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetDecksRequest fetches multiple decks by ID.
type BatchGetDecksRequest struct {
	DeckIDs []string `json:"deckIds,omitempty"`
}

// BatchGetDecksResponse returns decks requested in batch.
type BatchGetDecksResponse struct {
	Decks    []Deck           `json:"decks,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *BatchGetDecksResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListStarredDecksRequest lists decks starred by a user.
type ListStarredDecksRequest struct {
	UserID string
	Limit  *int32
	Offset *int32
}

// ListStarredDecksResponse returns starred decks for a user.
type ListStarredDecksResponse struct {
	Decks    []Deck           `json:"decks,omitempty"`
	Total    int32            `json:"total,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ListStarredDecksResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListDecks enumerates decks visible to the caller.
func (c *Client) ListDecks(ctx context.Context, request *ListDecksRequest, opts ...RequestOpt) (*ListDecksResponse, error) {
	if request == nil {
		request = &ListDecksRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/decks", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if scope := strings.TrimSpace(string(request.Scope)); scope != "" && scope != string(DeckListScopeUnspecified) {
		query.Set("scope", scope)
	}
	if userID := strings.TrimSpace(request.UserID); userID != "" {
		query.Set("userId", userID)
	}
	if heroID := strings.TrimSpace(request.HeroID); heroID != "" {
		query.Set("heroId", heroID)
	}
	if format := strings.TrimSpace(request.Format); format != "" {
		query.Set("format", format)
	}
	if name := strings.TrimSpace(request.Name); name != "" {
		query.Set("name", name)
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	req.URL.RawQuery = query.Encode()

	response := &ListDecksResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateDeck provisions a new deck.
func (c *Client) CreateDeck(ctx context.Context, request *CreateDeckRequest, opts ...RequestOpt) (*CreateDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/decks", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CreateDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// SearchDecks searches for decks matching provided filters.
func (c *Client) SearchDecks(ctx context.Context, request *SearchDecksRequest, opts ...RequestOpt) (*SearchDecksResponse, error) {
	if request == nil {
		request = &SearchDecksRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/decks/search", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if heroID := strings.TrimSpace(request.HeroID); heroID != "" {
		query.Set("heroId", heroID)
	}
	if format := strings.TrimSpace(request.Format); format != "" {
		query.Set("format", format)
	}
	if term := strings.TrimSpace(request.SearchTerm); term != "" {
		query.Set("searchTerm", term)
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	req.URL.RawQuery = query.Encode()

	response := &SearchDecksResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeck retrieves a deck by ID.
func (c *Client) GetDeck(ctx context.Context, request *GetDeckRequest, opts ...RequestOpt) (*GetDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteDeck removes a deck by ID.
func (c *Client) DeleteDeck(ctx context.Context, request *DeleteDeckRequest, opts ...RequestOpt) (*DeleteDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateDeck applies changes to an existing deck.
func (c *Client) UpdateDeck(ctx context.Context, request *UpdateDeckRequest, opts ...RequestOpt) (*UpdateDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	body, err := jsonBody(struct {
		Name       *string          `json:"name,omitempty"`
		Visibility *VisibilityLevel `json:"visibility,omitempty"`
	}{
		Name:       request.Name,
		Visibility: request.Visibility,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/decks/%s", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// StarDeck marks a deck as starred for the current user.
func (c *Client) StarDeck(ctx context.Context, request *StarDeckRequest, opts ...RequestOpt) (*StarDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/stars", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	response := &StarDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UnstarDeck removes the starred marker from a deck for the current user.
func (c *Client) UnstarDeck(ctx context.Context, request *UnstarDeckRequest, opts ...RequestOpt) (*UnstarDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/stars", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &UnstarDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListDeckVersions enumerates versions for a deck.
func (c *Client) ListDeckVersions(ctx context.Context, request *ListDeckVersionsRequest, opts ...RequestOpt) (*ListDeckVersionsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/versions", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	req.URL.RawQuery = query.Encode()

	response := &ListDeckVersionsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateDeckVersion provisions a new deck version.
func (c *Client) CreateDeckVersion(ctx context.Context, request *CreateDeckVersionRequest, opts ...RequestOpt) (*CreateDeckVersionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	body, err := jsonBody(struct {
		Version     string `json:"version,omitempty"`
		ImageURL    string `json:"imageUrl,omitempty"`
		Description string `json:"description,omitempty"`
	}{
		Version:     strings.TrimSpace(request.Version),
		ImageURL:    strings.TrimSpace(request.ImageURL),
		Description: strings.TrimSpace(request.Description),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/decks/%s/versions", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CreateDeckVersionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeckVersion retrieves a specific deck version.
func (c *Client) GetDeckVersion(ctx context.Context, request *GetDeckVersionRequest, opts ...RequestOpt) (*GetDeckVersionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	version := strings.TrimSpace(request.Version)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	if version == "" {
		return nil, errors.New("version must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/versions/%s", url.PathEscape(deckID), url.PathEscape(version))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	if token := strings.TrimSpace(request.NextToken); token != "" {
		query := req.URL.Query()
		query.Set("nextToken", token)
		req.URL.RawQuery = query.Encode()
	}

	response := &GetDeckVersionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteDeckVersion removes a deck version.
func (c *Client) DeleteDeckVersion(ctx context.Context, request *DeleteDeckVersionRequest, opts ...RequestOpt) (*DeleteDeckVersionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	version := strings.TrimSpace(request.Version)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	if version == "" {
		return nil, errors.New("version must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/versions/%s", url.PathEscape(deckID), url.PathEscape(version))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteDeckVersionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateDeckVersion updates attributes for a deck version.
func (c *Client) UpdateDeckVersion(ctx context.Context, request *UpdateDeckVersionRequest, opts ...RequestOpt) (*UpdateDeckVersionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	version := strings.TrimSpace(request.Version)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	if version == "" {
		return nil, errors.New("version must not be empty")
	}

	body, err := jsonBody(struct {
		ImageURL    *string `json:"imageUrl,omitempty"`
		Description *string `json:"description,omitempty"`
	}{
		ImageURL:    request.ImageURL,
		Description: request.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/decks/%s/versions/%s", url.PathEscape(deckID), url.PathEscape(version))
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateDeckVersionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListDeckVersionCards returns cards contained within a deck version.
func (c *Client) ListDeckVersionCards(ctx context.Context, request *ListDeckVersionCardsRequest, opts ...RequestOpt) (*ListDeckVersionCardsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	version := strings.TrimSpace(request.Version)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	if version == "" {
		return nil, errors.New("version must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/versions/%s/cards", url.PathEscape(deckID), url.PathEscape(version))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &ListDeckVersionCardsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ModifyDeckVersionCard adjusts a card entry within a deck version.
func (c *Client) ModifyDeckVersionCard(ctx context.Context, request *ModifyDeckVersionCardRequest, opts ...RequestOpt) (*ModifyDeckVersionCardResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	version := strings.TrimSpace(request.Version)
	cardID := strings.TrimSpace(request.CardID)
	board := strings.TrimSpace(string(request.Board))
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	if version == "" {
		return nil, errors.New("version must not be empty")
	}
	if cardID == "" {
		return nil, errors.New("cardID must not be empty")
	}
	if board == "" || board == string(BoardTypeUnspecified) {
		return nil, errors.New("boardType must not be empty")
	}

	body, err := jsonBody(struct {
		CardID   string `json:"cardId,omitempty"`
		Board    string `json:"boardType,omitempty"`
		Quantity int32  `json:"quantity,omitempty"`
	}{
		CardID:   cardID,
		Board:    board,
		Quantity: request.Quantity,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/decks/%s/versions/%s/cards", url.PathEscape(deckID), url.PathEscape(version))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &ModifyDeckVersionCardResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeckVersionHistory retrieves change history for a deck version.
func (c *Client) GetDeckVersionHistory(ctx context.Context, request *GetDeckVersionHistoryRequest, opts ...RequestOpt) (*GetDeckVersionHistoryResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	version := strings.TrimSpace(request.Version)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	if version == "" {
		return nil, errors.New("version must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/versions/%s/history", url.PathEscape(deckID), url.PathEscape(version))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetDeckVersionHistoryResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeckVersionNotes fetches notes for a deck version.
func (c *Client) GetDeckVersionNotes(ctx context.Context, request *GetDeckVersionNotesRequest, opts ...RequestOpt) (*GetDeckVersionNotesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	version := strings.TrimSpace(request.Version)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	if version == "" {
		return nil, errors.New("version must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/versions/%s/notes", url.PathEscape(deckID), url.PathEscape(version))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetDeckVersionNotesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateDeckVersionNotes sets notes for a deck version.
func (c *Client) UpdateDeckVersionNotes(ctx context.Context, request *UpdateDeckVersionNotesRequest, opts ...RequestOpt) (*UpdateDeckVersionNotesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	version := strings.TrimSpace(request.Version)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	if version == "" {
		return nil, errors.New("version must not be empty")
	}

	body, err := jsonBody(struct {
		Notes string `json:"notes,omitempty"`
	}{
		Notes: request.Notes,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/decks/%s/versions/%s/notes", url.PathEscape(deckID), url.PathEscape(version))
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateDeckVersionNotesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchGetDecks fetches decks in bulk.
func (c *Client) BatchGetDecks(ctx context.Context, request *BatchGetDecksRequest, opts ...RequestOpt) (*BatchGetDecksResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.DeckIDs) == 0 {
		return nil, errors.New("deckIds must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/decks:batchGet", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &BatchGetDecksResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListStarredDecks lists decks starred by a specific user.
func (c *Client) ListStarredDecks(ctx context.Context, request *ListStarredDecksRequest, opts ...RequestOpt) (*ListStarredDecksResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/v1/users/%s/stars/decks", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.Limit != nil && *request.Limit >= 0 {
		query.Set("limit", strconv.Itoa(int(*request.Limit)))
	}
	if request.Offset != nil && *request.Offset >= 0 {
		query.Set("offset", strconv.Itoa(int(*request.Offset)))
	}
	req.URL.RawQuery = query.Encode()

	response := &ListStarredDecksResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}
