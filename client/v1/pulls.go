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

// PullScope represents pulls.v1.PullScope.
type PullScope string

const (
	PullScopeUnspecified PullScope = "PULL_SCOPE_UNSPECIFIED"
	PullScopePack        PullScope = "PULL_SCOPE_PACK"
	PullScopeBox         PullScope = "PULL_SCOPE_BOX"
	PullScopeCase        PullScope = "PULL_SCOPE_CASE"
)

// PullCollectionSyncStatus represents pulls.v1.PullCollectionSyncStatus.
type PullCollectionSyncStatus string

const (
	PullCollectionSyncStatusUnspecified PullCollectionSyncStatus = "PULL_COLLECTION_SYNC_STATUS_UNSPECIFIED"
	PullCollectionSyncStatusPending     PullCollectionSyncStatus = "PULL_COLLECTION_SYNC_STATUS_PENDING"
	PullCollectionSyncStatusApplied     PullCollectionSyncStatus = "PULL_COLLECTION_SYNC_STATUS_APPLIED"
	PullCollectionSyncStatusFailed      PullCollectionSyncStatus = "PULL_COLLECTION_SYNC_STATUS_FAILED"
)

// PullCollectionSyncFailure represents pulls.v1.PullCollectionSyncFailure.
type PullCollectionSyncFailure string

const (
	PullCollectionSyncFailureUnspecified           PullCollectionSyncFailure = "PULL_COLLECTION_SYNC_FAILURE_UNSPECIFIED"
	PullCollectionSyncFailureCollectionUnavailable PullCollectionSyncFailure = "PULL_COLLECTION_SYNC_FAILURE_COLLECTION_UNAVAILABLE"
	PullCollectionSyncFailureRejected              PullCollectionSyncFailure = "PULL_COLLECTION_SYNC_FAILURE_REJECTED"
	PullCollectionSyncFailureRetriesExhausted      PullCollectionSyncFailure = "PULL_COLLECTION_SYNC_FAILURE_RETRIES_EXHAUSTED"
)

// CurrencyAmount represents pulls.v1.CurrencyAmount.
type CurrencyAmount struct {
	Amount   float64 `json:"amount,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

// PullValueSnapshot represents pulls.v1.PullValueSnapshot.
type PullValueSnapshot struct {
	Currency         string     `json:"currency,omitempty"`
	Source           string     `json:"source,omitempty"`
	ValuedAt         *time.Time `json:"valuedAt,omitempty"`
	TotalValue       float64    `json:"totalValue,omitempty"`
	PricedQuantity   int32      `json:"pricedQuantity,omitempty"`
	UnpricedQuantity int32      `json:"unpricedQuantity,omitempty"`
}

// Pull mirrors pulls.v1.Pull.
type Pull struct {
	ID                    string                    `json:"id,omitempty"`
	OwnerID               string                    `json:"ownerId,omitempty"`
	SealedProductID       string                    `json:"sealedProductId,omitempty"`
	Scope                 PullScope                 `json:"scope,omitempty"`
	SetID                 string                    `json:"setId,omitempty"`
	SetName               string                    `json:"setName,omitempty"`
	UnitsOpened           int32                     `json:"unitsOpened,omitempty"`
	PulledAt              *time.Time                `json:"pulledAt,omitempty"`
	CreatedAt             *time.Time                `json:"createdAt,omitempty"`
	UpdatedAt             *time.Time                `json:"updatedAt,omitempty"`
	PricePaid             *CurrencyAmount           `json:"pricePaid,omitempty"`
	Value                 *PullValueSnapshot        `json:"value,omitempty"`
	CollectionID          string                    `json:"collectionId,omitempty"`
	CollectionSyncStatus  PullCollectionSyncStatus  `json:"collectionSyncStatus,omitempty"`
	CollectionSyncFailure PullCollectionSyncFailure `json:"collectionSyncFailure,omitempty"`
	CollectionSyncedAt    *time.Time                `json:"collectionSyncedAt,omitempty"`
}

// PulledCard mirrors pulls.v1.PulledCard.
type PulledCard struct {
	ProductID  string   `json:"productId,omitempty"`
	Quantity   int32    `json:"quantity,omitempty"`
	Rarity     Rarity   `json:"rarity,omitempty"`
	Foiling    Foiling  `json:"foiling,omitempty"`
	PackNumber *int32   `json:"packNumber,omitempty"`
	UnitValue  *float64 `json:"unitValue,omitempty"`
}

// RarityCount mirrors pulls.v1.RarityCount.
type RarityCount struct {
	Rarity   Rarity `json:"rarity,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
}

// FoilingCount mirrors pulls.v1.FoilingCount.
type FoilingCount struct {
	Foiling  Foiling `json:"foiling,omitempty"`
	Quantity int32   `json:"quantity,omitempty"`
}

// RarityFoilingCount mirrors pulls.v1.RarityFoilingCount.
type RarityFoilingCount struct {
	Rarity   Rarity  `json:"rarity,omitempty"`
	Foiling  Foiling `json:"foiling,omitempty"`
	Quantity int32   `json:"quantity,omitempty"`
}

// PulledCardInput mirrors pulls.v1.PulledCardInput.
type PulledCardInput struct {
	ProductID  string `json:"productId,omitempty"`
	Quantity   int32  `json:"quantity,omitempty"`
	PackNumber *int32 `json:"packNumber,omitempty"`
}

// PullContents mirrors pulls.v1.PullContents.
type PullContents struct {
	PulledCards  []PulledCardInput `json:"pulledCards,omitempty"`
	RarityTotals []RarityCount     `json:"rarityTotals,omitempty"`
}

// CreatePullRequest creates a pull entry.
type CreatePullRequest struct {
	SealedProductID string            `json:"sealedProductId,omitempty"`
	UnitsOpened     int32             `json:"unitsOpened,omitempty"`
	PulledAt        *time.Time        `json:"pulledAt,omitempty"`
	PullID          string            `json:"pullId,omitempty"`
	PricePaid       *CurrencyAmount   `json:"pricePaid,omitempty"`
	CollectionID    string            `json:"collectionId,omitempty"`
	PulledCards     []PulledCardInput `json:"pulledCards,omitempty"`
	RarityTotals    []RarityCount     `json:"rarityTotals,omitempty"`
}

// CreatePullResponse returns the created pull.
type CreatePullResponse struct {
	Pull     *Pull            `json:"pull,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CreatePullResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetPullRequest identifies a pull.
type GetPullRequest struct {
	PullID string `json:"-"`
}

// GetPullResponse returns a pull and optional details.
type GetPullResponse struct {
	Pull         *Pull            `json:"pull,omitempty"`
	PulledCards  []PulledCard     `json:"pulledCards,omitempty"`
	RarityTotals []RarityCount    `json:"rarityTotals,omitempty"`
	Metadata     ResponseMetadata `json:"-"`
}

func (r *GetPullResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListPullsRequest filters pull history.
type ListPullsRequest struct {
	SetID     string
	Scope     PullScope
	PageSize  *int32
	NextToken string
}

// ListPullsResponse returns pull history.
type ListPullsResponse struct {
	Pulls     []Pull           `json:"pulls,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListPullsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeletePullRequest identifies a pull deletion target.
// RemoveFromCollection also removes the cards the pull added to a collection.
type DeletePullRequest struct {
	PullID               string `json:"-"`
	RemoveFromCollection bool
}

// DeletePullResponse captures response metadata.
type DeletePullResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeletePullResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetPullStatsRequest filters pull stats.
type GetPullStatsRequest struct {
	SetID string
	Scope PullScope
}

// PullStatsByScope mirrors pulls.v1.PullStatsByScope.
type PullStatsByScope struct {
	Scope               PullScope            `json:"scope,omitempty"`
	PullsCount          int32                `json:"pullsCount,omitempty"`
	UnitsOpened         int32                `json:"unitsOpened,omitempty"`
	RarityTotals        []RarityCount        `json:"rarityTotals,omitempty"`
	DetailedPullsCount  int32                `json:"detailedPullsCount,omitempty"`
	DetailedUnitsOpened int32                `json:"detailedUnitsOpened,omitempty"`
	FoilingTotals       []FoilingCount       `json:"foilingTotals,omitempty"`
	RarityFoilingTotals []RarityFoilingCount `json:"rarityFoilingTotals,omitempty"`
	CurrencyTotals      []PullCurrencyTotals `json:"currencyTotals,omitempty"`
}

// PullCurrencyTotals mirrors pulls.v1.PullCurrencyTotals.
type PullCurrencyTotals struct {
	Currency          string  `json:"currency,omitempty"`
	TotalSpent        float64 `json:"totalSpent,omitempty"`
	SpentPullsCount   int32   `json:"spentPullsCount,omitempty"`
	SpentUnitsOpened  int32   `json:"spentUnitsOpened,omitempty"`
	TotalPulledValue  float64 `json:"totalPulledValue,omitempty"`
	ValuedPullsCount  int32   `json:"valuedPullsCount,omitempty"`
	ValuedUnitsOpened int32   `json:"valuedUnitsOpened,omitempty"`
	UnpricedQuantity  int32   `json:"unpricedQuantity,omitempty"`
}

// GetPullStatsResponse returns pull stats grouped by scope.
type GetPullStatsResponse struct {
	ByScope  []PullStatsByScope `json:"byScope,omitempty"`
	Metadata ResponseMetadata   `json:"-"`
}

func (r *GetPullStatsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreatePull creates a pull entry.
func (c *Client) CreatePull(ctx context.Context, request *CreatePullRequest, opts ...RequestOpt) (*CreatePullResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.SealedProductID) == "" {
		return nil, errors.New("sealedProductID must not be empty")
	}
	if request.UnitsOpened <= 0 {
		return nil, errors.New("unitsOpened must be greater than zero")
	}

	body, err := jsonBody(struct {
		SealedProductID string            `json:"sealedProductId,omitempty"`
		UnitsOpened     int32             `json:"unitsOpened,omitempty"`
		PulledAt        *time.Time        `json:"pulledAt,omitempty"`
		PullID          string            `json:"pullId,omitempty"`
		PricePaid       *CurrencyAmount   `json:"pricePaid,omitempty"`
		CollectionID    string            `json:"collectionId,omitempty"`
		PulledCards     []PulledCardInput `json:"pulledCards,omitempty"`
		RarityTotals    []RarityCount     `json:"rarityTotals,omitempty"`
	}{
		SealedProductID: strings.TrimSpace(request.SealedProductID),
		UnitsOpened:     request.UnitsOpened,
		PulledAt:        request.PulledAt,
		PullID:          strings.TrimSpace(request.PullID),
		PricePaid:       request.PricePaid,
		CollectionID:    strings.TrimSpace(request.CollectionID),
		PulledCards:     request.PulledCards,
		RarityTotals:    request.RarityTotals,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/pulls", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CreatePullResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetPull retrieves a pull entry.
func (c *Client) GetPull(ctx context.Context, request *GetPullRequest, opts ...RequestOpt) (*GetPullResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	pullID := strings.TrimSpace(request.PullID)
	if pullID == "" {
		return nil, errors.New("pullID must not be empty")
	}

	path := fmt.Sprintf("/v1/pulls/%s", url.PathEscape(pullID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetPullResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListPulls lists pull entries for the authenticated user.
func (c *Client) ListPulls(ctx context.Context, request *ListPullsRequest, opts ...RequestOpt) (*ListPullsResponse, error) {
	if request == nil {
		request = &ListPullsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/pulls", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "setId", request.SetID)
	if scope := strings.TrimSpace(string(request.Scope)); scope != "" && scope != string(PullScopeUnspecified) {
		query.Set("scope", scope)
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	setQueryString(query, "nextToken", request.NextToken)
	req.URL.RawQuery = query.Encode()

	response := &ListPullsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeletePull deletes a pull entry.
func (c *Client) DeletePull(ctx context.Context, request *DeletePullRequest, opts ...RequestOpt) (*DeletePullResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	pullID := strings.TrimSpace(request.PullID)
	if pullID == "" {
		return nil, errors.New("pullID must not be empty")
	}

	path := fmt.Sprintf("/v1/pulls/%s", url.PathEscape(pullID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	if request.RemoveFromCollection {
		query := req.URL.Query()
		query.Set("removeFromCollection", "true")
		req.URL.RawQuery = query.Encode()
	}

	response := &DeletePullResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetPullStats retrieves pull stats grouped by scope.
func (c *Client) GetPullStats(ctx context.Context, request *GetPullStatsRequest, opts ...RequestOpt) (*GetPullStatsResponse, error) {
	if request == nil {
		request = &GetPullStatsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/pulls:stats", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "setId", request.SetID)
	if scope := strings.TrimSpace(string(request.Scope)); scope != "" && scope != string(PullScopeUnspecified) {
		query.Set("scope", scope)
	}
	req.URL.RawQuery = query.Encode()

	response := &GetPullStatsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdatePullRequest updates a pull entry. Nil fields are left unchanged.
// Contents, when set, replaces the pull's cards and rarity totals.
// CollectionID starts or retries adding the pull's cards to a collection.
// ExpectedUpdatedAt guards against concurrent edits.
type UpdatePullRequest struct {
	PullID            string          `json:"-"`
	UnitsOpened       *int32          `json:"unitsOpened,omitempty"`
	PulledAt          *time.Time      `json:"pulledAt,omitempty"`
	PricePaid         *CurrencyAmount `json:"pricePaid,omitempty"`
	ClearPricePaid    bool            `json:"clearPricePaid,omitempty"`
	Contents          *PullContents   `json:"contents,omitempty"`
	CollectionID      *string         `json:"collectionId,omitempty"`
	ExpectedUpdatedAt *time.Time      `json:"expectedUpdatedAt,omitempty"`
}

// UpdatePullResponse returns the updated pull and its details.
type UpdatePullResponse struct {
	Pull         *Pull            `json:"pull,omitempty"`
	PulledCards  []PulledCard     `json:"pulledCards,omitempty"`
	RarityTotals []RarityCount    `json:"rarityTotals,omitempty"`
	Metadata     ResponseMetadata `json:"-"`
}

func (r *UpdatePullResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetCommunityPullRatesRequest selects the set, and optionally the sealed product and scope, to report on.
type GetCommunityPullRatesRequest struct {
	SetID           string
	SealedProductID string
	Scope           PullScope
}

// CommunityPullRatesThresholds mirrors pulls.v1.CommunityPullRatesThresholds.
type CommunityPullRatesThresholds struct {
	MinContributors int32 `json:"minContributors,omitempty"`
	MinUnitsOpened  int32 `json:"minUnitsOpened,omitempty"`
}

// CommunityPullRatesByScope mirrors pulls.v1.CommunityPullRatesByScope.
// Withheld scopes have not met the publication thresholds.
type CommunityPullRatesByScope struct {
	Scope               PullScope            `json:"scope,omitempty"`
	ContributorsCount   int32                `json:"contributorsCount,omitempty"`
	Withheld            bool                 `json:"withheld,omitempty"`
	PullsCount          int32                `json:"pullsCount,omitempty"`
	UnitsOpened         int32                `json:"unitsOpened,omitempty"`
	RarityTotals        []RarityCount        `json:"rarityTotals,omitempty"`
	FoilingTotals       []FoilingCount       `json:"foilingTotals,omitempty"`
	RarityFoilingTotals []RarityFoilingCount `json:"rarityFoilingTotals,omitempty"`
}

// GetCommunityPullRatesResponse returns aggregated community pull rates.
type GetCommunityPullRatesResponse struct {
	SetID           string                        `json:"setId,omitempty"`
	SealedProductID string                        `json:"sealedProductId,omitempty"`
	ByScope         []CommunityPullRatesByScope   `json:"byScope,omitempty"`
	Thresholds      *CommunityPullRatesThresholds `json:"thresholds,omitempty"`
	ComputedAt      *time.Time                    `json:"computedAt,omitempty"`
	Metadata        ResponseMetadata              `json:"-"`
}

func (r *GetCommunityPullRatesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// PullSettings mirrors pulls.v1.PullSettings.
type PullSettings struct {
	ContributeToCommunityRates bool       `json:"contributeToCommunityRates,omitempty"`
	UpdatedAt                  *time.Time `json:"updatedAt,omitempty"`
}

// GetPullSettingsResponse returns the caller's pull settings.
type GetPullSettingsResponse struct {
	Settings *PullSettings    `json:"settings,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetPullSettingsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdatePullSettingsRequest updates the caller's pull settings. Nil fields are left unchanged.
type UpdatePullSettingsRequest struct {
	ContributeToCommunityRates *bool `json:"contributeToCommunityRates,omitempty"`
}

// UpdatePullSettingsResponse returns the updated pull settings.
type UpdatePullSettingsResponse struct {
	Settings *PullSettings    `json:"settings,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdatePullSettingsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdatePull updates a pull entry.
func (c *Client) UpdatePull(ctx context.Context, request *UpdatePullRequest, opts ...RequestOpt) (*UpdatePullResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	pullID := strings.TrimSpace(request.PullID)
	if pullID == "" {
		return nil, errors.New("pullID must not be empty")
	}
	if request.ClearPricePaid && request.PricePaid != nil {
		return nil, errors.New("pricePaid and clearPricePaid must not both be set")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/pulls/%s", url.PathEscape(pullID))
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdatePullResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetCommunityPullRates retrieves aggregated community pull rates for a set.
func (c *Client) GetCommunityPullRates(ctx context.Context, request *GetCommunityPullRatesRequest, opts ...RequestOpt) (*GetCommunityPullRatesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.SetID) == "" {
		return nil, errors.New("setID must not be empty")
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/pulls:communityRates", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "setId", request.SetID)
	setQueryString(query, "sealedProductId", request.SealedProductID)
	setQueryEnum(query, "scope", string(request.Scope), string(PullScopeUnspecified))
	req.URL.RawQuery = query.Encode()

	response := &GetCommunityPullRatesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetPullSettings retrieves the caller's pull settings.
func (c *Client) GetPullSettings(ctx context.Context, opts ...RequestOpt) (*GetPullSettingsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/pulls:settings", nil)
	if err != nil {
		return nil, err
	}

	response := &GetPullSettingsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdatePullSettings updates the caller's pull settings.
func (c *Client) UpdatePullSettings(ctx context.Context, request *UpdatePullSettingsRequest, opts ...RequestOpt) (*UpdatePullSettingsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPatch, "/v1/pulls:settings", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdatePullSettingsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}
