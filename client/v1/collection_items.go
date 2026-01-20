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

// Condition represents v1Condition.
type Condition string

const (
	ConditionUnspecified   Condition = "CONDITION_UNSPECIFIED"
	ConditionNearMint      Condition = "NEAR_MINT"
	ConditionLightlyPlayed Condition = "LIGHTLY_PLAYED"
	ConditionHeavilyPlayed Condition = "HEAVILY_PLAYED"
	ConditionDamaged       Condition = "DAMAGED"
)

// CollectionItem mirrors v1CollectionItem.
type CollectionItem struct {
	ID             string     `json:"id,omitempty"`
	CollectionID   string     `json:"collectionId,omitempty"`
	OwnerID        string     `json:"ownerId,omitempty"`
	ProductID      string     `json:"productId,omitempty"`
	Quantity       int32      `json:"quantity,omitempty"`
	Condition      Condition  `json:"condition,omitempty"`
	PinnedAt       *time.Time `json:"pinnedAt,omitempty"`
	Value          float64    `json:"value,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
	UpdatedAt      *time.Time `json:"updatedAt,omitempty"`
	CardID         string     `json:"cardId,omitempty"`
	PrintingID     string     `json:"printingId,omitempty"`
	BackCardID     string     `json:"backCardId,omitempty"`
	BackPrintingID string     `json:"backPrintingId,omitempty"`
}

// ListCollectionItemsRequest captures filters for listing items.
type ListCollectionItemsRequest struct {
	CollectionID string
	CardID       string
	PrintingID   string
	ProductID    string
	PageSize     *int32
	NextToken    string
}

// ListCollectionItemsResponse returns paginated items.
type ListCollectionItemsResponse struct {
	Items     []CollectionItem `json:"items,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListCollectionItemsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetCollectionItemRequest identifies a collection item.
type GetCollectionItemRequest struct {
	ItemID string `json:"-"`
}

// GetCollectionItemResponse returns a single collection item.
type GetCollectionItemResponse struct {
	Item     *CollectionItem  `json:"item,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetCollectionItemResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateCollectionItemRequest creates an item within a collection.
type CreateCollectionItemRequest struct {
	CollectionID string    `json:"collectionId,omitempty"`
	ProductID    string    `json:"productId,omitempty"`
	Quantity     int32     `json:"quantity,omitempty"`
	Condition    Condition `json:"condition,omitempty"`
	Value        *float64  `json:"value,omitempty"`
	ItemID       string    `json:"itemId,omitempty"`
}

// CreateCollectionItemResponse returns the created item.
type CreateCollectionItemResponse struct {
	Item     *CollectionItem  `json:"item,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CreateCollectionItemResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateCollectionItemRequest adjusts quantity or condition for an item.
type UpdateCollectionItemRequest struct {
	ItemID    string     `json:"-"`
	Quantity  *int32     `json:"quantity,omitempty"`
	Condition *Condition `json:"condition,omitempty"`
	Value     *float64   `json:"value,omitempty"`
	Pinned    *bool      `json:"pinned,omitempty"`
}

// UpdateCollectionItemResponse returns the updated item.
type UpdateCollectionItemResponse struct {
	Item     *CollectionItem  `json:"item,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateCollectionItemResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteCollectionItemRequest identifies the item to delete.
type DeleteCollectionItemRequest struct {
	ItemID string `json:"-"`
}

// DeleteCollectionItemResponse captures metadata for delete operations.
type DeleteCollectionItemResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteCollectionItemResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetCollectionItemsRequest fetches items in bulk.
type BatchGetCollectionItemsRequest struct {
	ItemIDs      []string `json:"itemIds,omitempty"`
	AllowPartial bool     `json:"allowPartial,omitempty"`
}

// BatchGetCollectionItemsResponse returns bulk-fetched items.
type BatchGetCollectionItemsResponse struct {
	Items       []CollectionItem `json:"items,omitempty"`
	NotFoundIDs []string         `json:"notFoundIds,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *BatchGetCollectionItemsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListCollectionItems enumerates collection items with optional filters.
func (c *Client) ListCollectionItems(ctx context.Context, request *ListCollectionItemsRequest, opts ...RequestOpt) (*ListCollectionItemsResponse, error) {
	if request == nil {
		request = &ListCollectionItemsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/collection_items", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if collectionID := strings.TrimSpace(request.CollectionID); collectionID != "" {
		query.Set("collectionId", collectionID)
	}
	if cardID := strings.TrimSpace(request.CardID); cardID != "" {
		query.Set("cardId", cardID)
	}
	if printingID := strings.TrimSpace(request.PrintingID); printingID != "" {
		query.Set("printingId", printingID)
	}
	if productID := strings.TrimSpace(request.ProductID); productID != "" {
		query.Set("productId", productID)
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	req.URL.RawQuery = query.Encode()

	response := &ListCollectionItemsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetCollectionItem fetches an item by ID.
func (c *Client) GetCollectionItem(ctx context.Context, request *GetCollectionItemRequest, opts ...RequestOpt) (*GetCollectionItemResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ItemID) == "" {
		return nil, errors.New("itemID must not be empty")
	}

	path := fmt.Sprintf("/v1/collection_items/%s", url.PathEscape(request.ItemID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetCollectionItemResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateCollectionItem adds a new item to a collection.
func (c *Client) CreateCollectionItem(ctx context.Context, request *CreateCollectionItemRequest, opts ...RequestOpt) (*CreateCollectionItemResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collection_items", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CreateCollectionItemResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateCollectionItem modifies an existing collection item.
func (c *Client) UpdateCollectionItem(ctx context.Context, request *UpdateCollectionItemRequest, opts ...RequestOpt) (*UpdateCollectionItemResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ItemID) == "" {
		return nil, errors.New("itemID must not be empty")
	}

	body, err := jsonBody(struct {
		Quantity  *int32     `json:"quantity,omitempty"`
		Condition *Condition `json:"condition,omitempty"`
		Value     *float64   `json:"value,omitempty"`
		Pinned    *bool      `json:"pinned,omitempty"`
	}{
		Quantity:  request.Quantity,
		Condition: request.Condition,
		Value:     request.Value,
		Pinned:    request.Pinned,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/collection_items/%s", url.PathEscape(request.ItemID))
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateCollectionItemResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteCollectionItem removes an item by ID.
func (c *Client) DeleteCollectionItem(ctx context.Context, request *DeleteCollectionItemRequest, opts ...RequestOpt) (*DeleteCollectionItemResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ItemID) == "" {
		return nil, errors.New("itemID must not be empty")
	}

	path := fmt.Sprintf("/v1/collection_items/%s", url.PathEscape(request.ItemID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteCollectionItemResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchGetCollectionItems retrieves items in bulk.
func (c *Client) BatchGetCollectionItems(ctx context.Context, request *BatchGetCollectionItemsRequest, opts ...RequestOpt) (*BatchGetCollectionItemsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collection_items:batchGet", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &BatchGetCollectionItemsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}
