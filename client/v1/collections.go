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

// CollectionType matches the API's v1CollectionType enum.
type CollectionType string

const (
	CollectionTypeUnspecified CollectionType = "COLLECTION_TYPE_UNSPECIFIED"
	CollectionTypeBinder      CollectionType = "BINDER"
	CollectionTypeWantlist    CollectionType = "WANTLIST"
	CollectionTypeTradelist   CollectionType = "TRADELIST"
	CollectionTypeList        CollectionType = "LIST"
)

// CollectionPermission mirrors authzv1Permission.
type CollectionPermission string

const (
	CollectionPermissionUnspecified CollectionPermission = "PERMISSION_UNSPECIFIED"
	CollectionPermissionReader      CollectionPermission = "PERMISSION_READER"
	CollectionPermissionWriter      CollectionPermission = "PERMISSION_WRITER"
)

// Collection mirrors v1Collection from the API definition.
type Collection struct {
	ID             string          `json:"id,omitempty"`
	Name           string          `json:"name,omitempty"`
	Description    string          `json:"description,omitempty"`
	OwnerID        string          `json:"ownerId,omitempty"`
	CollectionType CollectionType  `json:"collectionType,omitempty"`
	Visibility     VisibilityLevel `json:"visibility,omitempty"`
	CreatedAt      *time.Time      `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time      `json:"updatedAt,omitempty"`
}

// CollectionStats captures aggregate metrics for a collection.
type CollectionStats struct {
	ItemsCount      int32 `json:"itemsCount,omitempty"`
	QuantityCount   int32 `json:"quantityCount,omitempty"`
	UniqueCardCount int32 `json:"uniqueCardCount,omitempty"`
}

// CollectionWithStats bundles a collection and its stats.
type CollectionWithStats struct {
	Collection *Collection      `json:"collection,omitempty"`
	Stats      *CollectionStats `json:"stats,omitempty"`
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
	Name           string          `json:"name,omitempty"`
	CollectionType CollectionType  `json:"collectionType,omitempty"`
	Description    string          `json:"description,omitempty"`
	Visibility     VisibilityLevel `json:"visibility,omitempty"`
	CollectionID   string          `json:"collectionId,omitempty"`
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
	Stats      *CollectionStats `json:"stats,omitempty"`
	Metadata   ResponseMetadata `json:"-"`
}

func (r *GetCollectionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateCollectionRequest models a partial update for a collection.
type UpdateCollectionRequest struct {
	CollectionID string  `json:"-"`
	Name         *string `json:"name,omitempty"`
	Description  *string `json:"description,omitempty"`
}

// UpdateCollectionResponse returns the updated collection.
type UpdateCollectionResponse struct {
	Collection *Collection      `json:"collection,omitempty"`
	Metadata   ResponseMetadata `json:"-"`
}

func (r *UpdateCollectionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateCollectionVisibilityRequest updates only the visibility field of a collection.
type UpdateCollectionVisibilityRequest struct {
	CollectionID string           `json:"-"`
	Visibility   *VisibilityLevel `json:"visibility,omitempty"`
}

// UpdateCollectionVisibilityResponse returns the updated collection.
type UpdateCollectionVisibilityResponse struct {
	Collection *Collection      `json:"collection,omitempty"`
	Metadata   ResponseMetadata `json:"-"`
}

func (r *UpdateCollectionVisibilityResponse) setMetadata(metadata ResponseMetadata) {
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
	Collections []CollectionWithStats `json:"collections,omitempty"`
	NotFoundIDs []string              `json:"notFoundIds,omitempty"`
	Metadata    ResponseMetadata      `json:"-"`
}

func (r *BatchGetCollectionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GrantCollectionAccessRequest assigns a permission for a collection.
type GrantCollectionAccessRequest struct {
	CollectionID string               `json:"resourceId,omitempty"`
	SubjectID    string               `json:"subjectId,omitempty"`
	Permission   CollectionPermission `json:"permission,omitempty"`
}

// GrantCollectionAccessResponse captures metadata for grant operations.
type GrantCollectionAccessResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *GrantCollectionAccessResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RevokeCollectionAccessRequest removes a permission for a collection.
type RevokeCollectionAccessRequest struct {
	CollectionID string               `json:"resourceId,omitempty"`
	SubjectID    string               `json:"subjectId,omitempty"`
	Permission   CollectionPermission `json:"permission,omitempty"`
}

// RevokeCollectionAccessResponse captures metadata for revoke operations.
type RevokeCollectionAccessResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *RevokeCollectionAccessResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CollectionAccessGrant describes an explicit permission grant for a collection.
type CollectionAccessGrant struct {
	SubjectID  string               `json:"subjectId,omitempty"`
	Permission CollectionPermission `json:"permission,omitempty"`
}

// GetCollectionAccessRequest identifies the collection access to retrieve.
type GetCollectionAccessRequest struct {
	CollectionID string `json:"-"`
}

// GetCollectionAccessResponse returns the caller's effective permission for a collection.
type GetCollectionAccessResponse struct {
	Permission CollectionPermission `json:"permission,omitempty"`
	Metadata   ResponseMetadata     `json:"-"`
}

func (r *GetCollectionAccessResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListCollectionAccessGrantsRequest lists subjects explicitly granted access to a collection.
type ListCollectionAccessGrantsRequest struct {
	CollectionID string `json:"-"`
	PageSize     *int32
	NextToken    string
}

// ListCollectionAccessGrantsResponse returns explicit access grants for a collection.
type ListCollectionAccessGrantsResponse struct {
	Grants    []CollectionAccessGrant `json:"grants,omitempty"`
	NextToken string                  `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata        `json:"-"`
}

func (r *ListCollectionAccessGrantsResponse) setMetadata(metadata ResponseMetadata) {
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

// ListCollections enumerates collections visible to the caller.
func (c *Client) ListCollections(ctx context.Context, request *ListCollectionsRequest, opts ...RequestOpt) (*ListCollectionsResponse, error) {
	if request == nil {
		request = &ListCollectionsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/collections", nil)
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

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collections", body)
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

	path := fmt.Sprintf("/v1/collections/%s", url.PathEscape(request.CollectionID))
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
		Name        *string `json:"name,omitempty"`
		Description *string `json:"description,omitempty"`
	}{
		Name:        request.Name,
		Description: request.Description,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/collections/%s", url.PathEscape(request.CollectionID))
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

// UpdateCollectionVisibility updates only the visibility of a collection.
func (c *Client) UpdateCollectionVisibility(ctx context.Context, request *UpdateCollectionVisibilityRequest, opts ...RequestOpt) (*UpdateCollectionVisibilityResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}
	if request.Visibility == nil {
		return nil, errors.New("visibility must not be nil")
	}

	body, err := jsonBody(struct {
		Visibility *VisibilityLevel `json:"visibility,omitempty"`
	}{
		Visibility: request.Visibility,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/collections/%s/visibility", url.PathEscape(request.CollectionID))
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateCollectionVisibilityResponse{}
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

	path := fmt.Sprintf("/v1/collections/%s", url.PathEscape(request.CollectionID))
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

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collections:batchGet", body)
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

// GetCollectionValuation fetches valuation metrics for a collection.
func (c *Client) GetCollectionValuation(ctx context.Context, request *GetCollectionValuationRequest, opts ...RequestOpt) (*GetCollectionValuationResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	path := fmt.Sprintf("/v1/collections/%s/valuation", url.PathEscape(request.CollectionID))
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

// GetCollectionAccess returns the caller's effective permission for a collection.
func (c *Client) GetCollectionAccess(ctx context.Context, request *GetCollectionAccessRequest, opts ...RequestOpt) (*GetCollectionAccessResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	path := fmt.Sprintf("/v1/collections/%s/access", url.PathEscape(request.CollectionID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetCollectionAccessResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListCollectionAccessGrants lists subjects explicitly granted access to a collection.
func (c *Client) ListCollectionAccessGrants(ctx context.Context, request *ListCollectionAccessGrantsRequest, opts ...RequestOpt) (*ListCollectionAccessGrantsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}

	path := fmt.Sprintf("/v1/collections/%s/permissions", url.PathEscape(request.CollectionID))
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

	response := &ListCollectionAccessGrantsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GrantCollectionAccess assigns a permission for a user to access a collection.
func (c *Client) GrantCollectionAccess(ctx context.Context, request *GrantCollectionAccessRequest, opts ...RequestOpt) (*GrantCollectionAccessResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	collectionID := strings.TrimSpace(request.CollectionID)
	if collectionID == "" {
		return nil, errors.New("collectionID must not be empty")
	}
	subjectID := strings.TrimSpace(request.SubjectID)
	if subjectID == "" {
		return nil, errors.New("subjectID must not be empty")
	}
	permission := CollectionPermission(strings.TrimSpace(string(request.Permission)))
	if permission == "" || permission == CollectionPermissionUnspecified {
		return nil, errors.New("permission must be specified")
	}

	body, err := jsonBody(struct {
		CollectionID string               `json:"resourceId,omitempty"`
		SubjectID    string               `json:"subjectId,omitempty"`
		Permission   CollectionPermission `json:"permission,omitempty"`
	}{
		CollectionID: collectionID,
		SubjectID:    subjectID,
		Permission:   permission,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collections/permissions:grant", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &GrantCollectionAccessResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// RevokeCollectionAccess removes a permission previously granted to access a collection.
func (c *Client) RevokeCollectionAccess(ctx context.Context, request *RevokeCollectionAccessRequest, opts ...RequestOpt) (*RevokeCollectionAccessResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	collectionID := strings.TrimSpace(request.CollectionID)
	if collectionID == "" {
		return nil, errors.New("collectionID must not be empty")
	}
	subjectID := strings.TrimSpace(request.SubjectID)
	if subjectID == "" {
		return nil, errors.New("subjectID must not be empty")
	}
	permission := CollectionPermission(strings.TrimSpace(string(request.Permission)))
	if permission == "" || permission == CollectionPermissionUnspecified {
		return nil, errors.New("permission must be specified")
	}

	body, err := jsonBody(struct {
		CollectionID string               `json:"resourceId,omitempty"`
		SubjectID    string               `json:"subjectId,omitempty"`
		Permission   CollectionPermission `json:"permission,omitempty"`
	}{
		CollectionID: collectionID,
		SubjectID:    subjectID,
		Permission:   permission,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collections/permissions:revoke", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &RevokeCollectionAccessResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}
