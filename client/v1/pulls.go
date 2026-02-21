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

// Pull mirrors pulls.v1.Pull.
type Pull struct {
	ID              string     `json:"id,omitempty"`
	OwnerID         string     `json:"ownerId,omitempty"`
	SealedProductID string     `json:"sealedProductId,omitempty"`
	Scope           PullScope  `json:"scope,omitempty"`
	SetID           string     `json:"setId,omitempty"`
	SetName         string     `json:"setName,omitempty"`
	UnitsOpened     int32      `json:"unitsOpened,omitempty"`
	PulledAt        *time.Time `json:"pulledAt,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
}

// PulledCard mirrors pulls.v1.PulledCard.
type PulledCard struct {
	ProductID string `json:"productId,omitempty"`
	Quantity  int32  `json:"quantity,omitempty"`
	Rarity    Rarity `json:"rarity,omitempty"`
}

// RarityCount mirrors pulls.v1.RarityCount.
type RarityCount struct {
	Rarity   Rarity `json:"rarity,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
}

// PulledCardInput mirrors pulls.v1.PulledCardInput.
type PulledCardInput struct {
	ProductID string `json:"productId,omitempty"`
	Quantity  int32  `json:"quantity,omitempty"`
}

// CreatePullRequest creates a pull entry.
type CreatePullRequest struct {
	SealedProductID string            `json:"sealedProductId,omitempty"`
	UnitsOpened     int32             `json:"unitsOpened,omitempty"`
	PulledAt        *time.Time        `json:"pulledAt,omitempty"`
	PullID          string            `json:"pullId,omitempty"`
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
type DeletePullRequest struct {
	PullID string `json:"-"`
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
	Scope        PullScope     `json:"scope,omitempty"`
	PullsCount   int32         `json:"pullsCount,omitempty"`
	UnitsOpened  int32         `json:"unitsOpened,omitempty"`
	RarityTotals []RarityCount `json:"rarityTotals,omitempty"`
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
		PulledCards     []PulledCardInput `json:"pulledCards,omitempty"`
		RarityTotals    []RarityCount     `json:"rarityTotals,omitempty"`
	}{
		SealedProductID: strings.TrimSpace(request.SealedProductID),
		UnitsOpened:     request.UnitsOpened,
		PulledAt:        request.PulledAt,
		PullID:          strings.TrimSpace(request.PullID),
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
