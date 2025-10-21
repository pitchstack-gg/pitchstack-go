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

// CollectionListScope matches the API's v1CollectionListScope enum.
type CollectionListScope string

const (
	CollectionListScopeUnspecified CollectionListScope = "COLLECTION_LIST_SCOPE_UNSPECIFIED"
	CollectionListScopeOwned       CollectionListScope = "COLLECTION_LIST_SCOPE_OWNED"
	CollectionListScopeShared      CollectionListScope = "COLLECTION_LIST_SCOPE_SHARED"
	CollectionListScopeAccessible  CollectionListScope = "COLLECTION_LIST_SCOPE_ACCESSIBLE"
)

// VisibilityLevel matches commonv1VisibilityLevel.
type VisibilityLevel string

const (
	VisibilityLevelUnspecified VisibilityLevel = "VISIBILITY_LEVEL_UNSPECIFIED"
	VisibilityLevelPrivate     VisibilityLevel = "VISIBILITY_LEVEL_PRIVATE"
	VisibilityLevelShared      VisibilityLevel = "VISIBILITY_LEVEL_SHARED"
	VisibilityLevelPublic      VisibilityLevel = "VISIBILITY_LEVEL_PUBLIC"
)

// SyncAction represents v1SyncAction.
type SyncAction string

const (
	SyncActionUnspecified SyncAction = "SYNC_ACTION_UNSPECIFIED"
	SyncActionUpsert      SyncAction = "UPSERT"
	SyncActionDelete      SyncAction = "DELETE"
)

// SyncStatus represents v1SyncStatus.
type SyncStatus string

const (
	SyncStatusUnspecified SyncStatus = "SYNC_STATUS_UNSPECIFIED"
	SyncStatusOK          SyncStatus = "OK"
	SyncStatusConflict    SyncStatus = "CONFLICT"
	SyncStatusError       SyncStatus = "ERROR"
)

// Collection mirrors v1Collection from the API definition.
type Collection struct {
	ID              string          `json:"id,omitempty"`
	Name            string          `json:"name,omitempty"`
	Description     string          `json:"description,omitempty"`
	UserID          string          `json:"userId,omitempty"`
	Visibility      VisibilityLevel `json:"visibility,omitempty"`
	CreatedAt       *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time      `json:"updatedAt,omitempty"`
	ItemsCount      int32           `json:"itemsCount,omitempty"`
	QuantityCount   int32           `json:"quantityCount,omitempty"`
	UniqueCardCount int32           `json:"uniqueCardCount,omitempty"`
}

// ListCollectionsRequest captures query parameters for ListCollections.
type ListCollectionsRequest struct {
	Scope     CollectionListScope
	UserID    string
	PageSize  *int32
	NextToken string
}

// ListCollectionsResponse is returned from ListCollections calls.
type ListCollectionsResponse struct {
	Collections []Collection     `json:"collections,omitempty"`
	NextToken   string           `json:"nextToken,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *ListCollectionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateCollectionRequest captures payload fields for collection creation.
type CreateCollectionRequest struct {
	Name        string          `json:"name,omitempty"`
	Description string          `json:"description,omitempty"`
	Visibility  VisibilityLevel `json:"visibility,omitempty"`
}

// CreateCollectionResponse contains the created collection.
type CreateCollectionResponse struct {
	Collection *Collection      `json:"collection,omitempty"`
	Metadata   ResponseMetadata `json:"-"`
}

func (r *CreateCollectionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetCollectionRequest identifies the collection to fetch.
type GetCollectionRequest struct {
	CollectionID string `json:"-"`
}

// GetCollectionResponse returns a single collection.
type GetCollectionResponse struct {
	Collection *Collection      `json:"collection,omitempty"`
	Metadata   ResponseMetadata `json:"-"`
}

func (r *GetCollectionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateCollectionRequest models a partial update for a collection.
type UpdateCollectionRequest struct {
	CollectionID string           `json:"-"`
	Name         *string          `json:"name,omitempty"`
	Description  *string          `json:"description,omitempty"`
	Visibility   *VisibilityLevel `json:"visibility,omitempty"`
}

// UpdateCollectionResponse returns the updated collection.
type UpdateCollectionResponse struct {
	Collection *Collection      `json:"collection,omitempty"`
	Metadata   ResponseMetadata `json:"-"`
}

func (r *UpdateCollectionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteCollectionRequest identifies the collection to delete.
type DeleteCollectionRequest struct {
	CollectionID string `json:"-"`
}

// DeleteCollectionResponse captures metadata for delete calls.
type DeleteCollectionResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteCollectionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetCollectionsRequest fetches collections in bulk.
type BatchGetCollectionsRequest struct {
	CollectionIDs []string `json:"collectionIds,omitempty"`
	AllowPartial  bool     `json:"allowPartial,omitempty"`
}

// BatchGetCollectionsResponse returns collections fetched in bulk.
type BatchGetCollectionsResponse struct {
	Collections []Collection     `json:"collections,omitempty"`
	NotFoundIDs []string         `json:"notFoundIds,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *BatchGetCollectionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CollectionSyncOp describes a collections sync operation.
type CollectionSyncOp struct {
	OpID            string          `json:"opId,omitempty"`
	Action          SyncAction      `json:"action,omitempty"`
	ID              string          `json:"id,omitempty"`
	ClientUpdatedAt *time.Time      `json:"clientUpdatedAt,omitempty"`
	Name            string          `json:"name,omitempty"`
	Description     string          `json:"description,omitempty"`
	Visibility      VisibilityLevel `json:"visibility,omitempty"`
}

// CollectionsSyncRequest carries a batch of sync operations.
type CollectionsSyncRequest struct {
	Ops []CollectionSyncOp `json:"ops,omitempty"`
}

// CollectionSyncResult provides per-operation outcomes.
type CollectionSyncResult struct {
	OpID     string      `json:"opId,omitempty"`
	Status   SyncStatus  `json:"status,omitempty"`
	Entity   *Collection `json:"entity,omitempty"`
	ServerID string      `json:"serverId,omitempty"`
	Message  string      `json:"message,omitempty"`
}

// CollectionsSyncResponse aggregates sync results.
type CollectionsSyncResponse struct {
	Results  []CollectionSyncResult `json:"results,omitempty"`
	Metadata ResponseMetadata       `json:"-"`
}

func (r *CollectionsSyncResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetCollectionValuationRequest identifies a collection valuation lookup.
type GetCollectionValuationRequest struct {
	CollectionID string `json:"-"`
	Source       string
}

// GetCollectionValuationResponse describes valuation metrics.
type GetCollectionValuationResponse struct {
	CollectionID        string           `json:"collectionId,omitempty"`
	Currency            string           `json:"currency,omitempty"`
	Source              string           `json:"source,omitempty"`
	TotalEstimatedValue float64          `json:"totalEstimatedValue,omitempty"`
	TotalItems          int32            `json:"totalItems,omitempty"`
	PricedItems         int32            `json:"pricedItems,omitempty"`
	MissingPriceItems   int32            `json:"missingPriceItems,omitempty"`
	ComputedAt          *time.Time       `json:"computedAt,omitempty"`
	Metadata            ResponseMetadata `json:"-"`
}

func (r *GetCollectionValuationResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListStarredCollectionsRequest enumerates starred collections for a user.
type ListStarredCollectionsRequest struct {
	UserID string `json:"-"`
	Limit  *int32
	Offset *int32
}

// ListStarredCollectionsResponse includes collections and totals.
type ListStarredCollectionsResponse struct {
	Collections []Collection     `json:"collections,omitempty"`
	Total       int32            `json:"total,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *ListStarredCollectionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// StarCollectionRequest toggles starring for a collection.
type StarCollectionRequest struct {
	CollectionID string `json:"-"`
}

// StarCollectionResponse carries response metadata.
type StarCollectionResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *StarCollectionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnstarCollectionRequest removes a star from a collection.
type UnstarCollectionRequest struct {
	CollectionID string `json:"-"`
}

// UnstarCollectionResponse carries response metadata.
type UnstarCollectionResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnstarCollectionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListCollections enumerates collections visible to the caller.
func (c *Client) ListCollections(ctx context.Context, request *ListCollectionsRequest, opts ...RequestOpt) (*ListCollectionsResponse, error) {
	if request == nil {
		request = &ListCollectionsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/collections", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()

	if scope := strings.TrimSpace(string(request.Scope)); scope != "" && scope != string(CollectionListScopeUnspecified) {
		query.Set("scope", scope)
	}
	if userID := strings.TrimSpace(request.UserID); userID != "" {
		query.Set("userId", userID)
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}

	req.URL.RawQuery = query.Encode()

	response := &ListCollectionsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateCollection provisions a new collection.
func (c *Client) CreateCollection(ctx context.Context, request *CreateCollectionRequest, opts ...RequestOpt) (*CreateCollectionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/collections", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CreateCollectionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetCollection retrieves a collection by ID.
func (c *Client) GetCollection(ctx context.Context, request *GetCollectionRequest, opts ...RequestOpt) (*GetCollectionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	path := fmt.Sprintf("/api/v1/collections/%s", url.PathEscape(request.CollectionID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetCollectionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateCollection applies a partial update to a collection.
func (c *Client) UpdateCollection(ctx context.Context, request *UpdateCollectionRequest, opts ...RequestOpt) (*UpdateCollectionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	body, err := jsonBody(struct {
		Name        *string          `json:"name,omitempty"`
		Description *string          `json:"description,omitempty"`
		Visibility  *VisibilityLevel `json:"visibility,omitempty"`
	}{
		Name:        request.Name,
		Description: request.Description,
		Visibility:  request.Visibility,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/api/v1/collections/%s", url.PathEscape(request.CollectionID))
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateCollectionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteCollection removes a collection by ID.
func (c *Client) DeleteCollection(ctx context.Context, request *DeleteCollectionRequest, opts ...RequestOpt) (*DeleteCollectionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	path := fmt.Sprintf("/api/v1/collections/%s", url.PathEscape(request.CollectionID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteCollectionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchGetCollections retrieves collections in bulk.
func (c *Client) BatchGetCollections(ctx context.Context, request *BatchGetCollectionsRequest, opts ...RequestOpt) (*BatchGetCollectionsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/collections:batchGet", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &BatchGetCollectionsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CollectionsSync synchronizes collections in bulk.
func (c *Client) CollectionsSync(ctx context.Context, request *CollectionsSyncRequest, opts ...RequestOpt) (*CollectionsSyncResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/collections:sync", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CollectionsSyncResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetCollectionValuation fetches valuation metrics for a collection.
func (c *Client) GetCollectionValuation(ctx context.Context, request *GetCollectionValuationRequest, opts ...RequestOpt) (*GetCollectionValuationResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	path := fmt.Sprintf("/api/v1/collections/%s/valuation", url.PathEscape(request.CollectionID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if source := strings.TrimSpace(request.Source); source != "" {
		query.Set("source", source)
	}
	req.URL.RawQuery = query.Encode()

	response := &GetCollectionValuationResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListStarredCollections returns starred collections for a user.
func (c *Client) ListStarredCollections(ctx context.Context, request *ListStarredCollectionsRequest, opts ...RequestOpt) (*ListStarredCollectionsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.UserID) == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/api/v1/users/%s/stars/collections", url.PathEscape(request.UserID))
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

	response := &ListStarredCollectionsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// StarCollection marks a collection as starred.
func (c *Client) StarCollection(ctx context.Context, request *StarCollectionRequest, opts ...RequestOpt) (*StarCollectionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	path := fmt.Sprintf("/api/v1/collections/%s/stars", url.PathEscape(request.CollectionID))
	req, err := c.newRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	response := &StarCollectionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UnstarCollection removes the starred state from a collection.
func (c *Client) UnstarCollection(ctx context.Context, request *UnstarCollectionRequest, opts ...RequestOpt) (*UnstarCollectionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	path := fmt.Sprintf("/api/v1/collections/%s/stars", url.PathEscape(request.CollectionID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &UnstarCollectionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}
