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
	TradeQuantity  int32      `json:"tradeQuantity,omitempty"`
	Notes          string     `json:"notes,omitempty"`
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
	CollectionID  string    `json:"collectionId,omitempty"`
	ProductID     string    `json:"productId,omitempty"`
	Quantity      int32     `json:"quantity,omitempty"`
	Condition     Condition `json:"condition,omitempty"`
	Value         *float64  `json:"value,omitempty"`
	ItemID        string    `json:"itemId,omitempty"`
	TradeQuantity int32     `json:"tradeQuantity,omitempty"`
	Notes         string    `json:"notes,omitempty"`
}

// CreateCollectionItemResponse returns the created item.
type CreateCollectionItemResponse struct {
	Item     *CollectionItem  `json:"item,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CreateCollectionItemResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// AdjustCollectionItemQuantityRequest adjusts a collection item quantity by product.
type AdjustCollectionItemQuantityRequest struct {
	CollectionID     string    `json:"collectionId,omitempty"`
	ProductID        string    `json:"productId,omitempty"`
	Condition        Condition `json:"condition,omitempty"`
	QuantityDelta    int32     `json:"quantityDelta,omitempty"`
	ItemID           string    `json:"itemId,omitempty"`
	ClientMutationID string    `json:"clientMutationId,omitempty"`
}

// AdjustCollectionItemQuantityResponse returns the adjusted item.
type AdjustCollectionItemQuantityResponse struct {
	Item     *CollectionItem  `json:"item,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *AdjustCollectionItemQuantityResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateCollectionItemRequest adjusts quantity or condition for an item.
type UpdateCollectionItemRequest struct {
	ItemID            string     `json:"-"`
	Quantity          *int32     `json:"quantity,omitempty"`
	Condition         *Condition `json:"condition,omitempty"`
	Value             *float64   `json:"value,omitempty"`
	Pinned            *bool      `json:"pinned,omitempty"`
	ExpectedUpdatedAt *time.Time `json:"expectedUpdatedAt,omitempty"`
	ClientMutationID  string     `json:"clientMutationId,omitempty"`
	TradeQuantity     *int32     `json:"tradeQuantity,omitempty"`
	Notes             *string    `json:"notes,omitempty"`
}

// UpdateCollectionItemResponse returns the updated item.
type UpdateCollectionItemResponse struct {
	Item     *CollectionItem  `json:"item,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateCollectionItemResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// TransferCollectionItemRequest moves an item to another collection.
type TransferCollectionItemRequest struct {
	ItemID                  string `json:"-"`
	DestinationCollectionID string `json:"destinationCollectionId,omitempty"`
}

// TransferCollectionItemResponse returns the transferred item.
type TransferCollectionItemResponse struct {
	Item     *CollectionItem  `json:"item,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *TransferCollectionItemResponse) setMetadata(metadata ResponseMetadata) {
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

// BatchCollectionItemFailure describes one item that was not mutated.
type BatchCollectionItemFailure struct {
	Index   int32  `json:"index,omitempty"`
	ItemID  string `json:"itemId,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

type BatchUpdateCollectionItemsRequest struct {
	CollectionID string                        `json:"collectionId,omitempty"`
	Requests     []UpdateCollectionItemRequest `json:"-"`
}

type BatchUpdateCollectionItemsResponse struct {
	Items    []CollectionItem             `json:"items,omitempty"`
	Failures []BatchCollectionItemFailure `json:"failures,omitempty"`
	Metadata ResponseMetadata             `json:"-"`
}

func (r *BatchUpdateCollectionItemsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

type BatchDeleteCollectionItemsRequest struct {
	CollectionID string   `json:"collectionId,omitempty"`
	ItemIDs      []string `json:"itemIds,omitempty"`
}

type BatchDeleteCollectionItemsResponse struct {
	DeletedItemIDs []string                     `json:"deletedItemIds,omitempty"`
	Failures       []BatchCollectionItemFailure `json:"failures,omitempty"`
	Metadata       ResponseMetadata             `json:"-"`
}

func (r *BatchDeleteCollectionItemsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

type BatchTransferCollectionItemsRequest struct {
	SourceCollectionID      string   `json:"sourceCollectionId,omitempty"`
	DestinationCollectionID string   `json:"destinationCollectionId,omitempty"`
	ItemIDs                 []string `json:"itemIds,omitempty"`
}

type BatchTransferCollectionItemsResponse struct {
	Items    []CollectionItem             `json:"items,omitempty"`
	Failures []BatchCollectionItemFailure `json:"failures,omitempty"`
	Metadata ResponseMetadata             `json:"-"`
}

func (r *BatchTransferCollectionItemsResponse) setMetadata(metadata ResponseMetadata) {
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

// AdjustCollectionItemQuantity adjusts a collection item quantity by product.
func (c *Client) AdjustCollectionItemQuantity(ctx context.Context, request *AdjustCollectionItemQuantityRequest, opts ...RequestOpt) (*AdjustCollectionItemQuantityResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}
	if strings.TrimSpace(request.ProductID) == "" {
		return nil, errors.New("productID must not be empty")
	}
	if request.QuantityDelta == 0 {
		return nil, errors.New("quantityDelta must not be zero")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collection_items:adjustQuantity", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &AdjustCollectionItemQuantityResponse{}
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
		Quantity          *int32     `json:"quantity,omitempty"`
		Condition         *Condition `json:"condition,omitempty"`
		Value             *float64   `json:"value,omitempty"`
		Pinned            *bool      `json:"pinned,omitempty"`
		ExpectedUpdatedAt *time.Time `json:"expectedUpdatedAt,omitempty"`
		ClientMutationID  string     `json:"clientMutationId,omitempty"`
		TradeQuantity     *int32     `json:"tradeQuantity,omitempty"`
		Notes             *string    `json:"notes,omitempty"`
	}{
		Quantity:          request.Quantity,
		Condition:         request.Condition,
		Value:             request.Value,
		Pinned:            request.Pinned,
		ExpectedUpdatedAt: request.ExpectedUpdatedAt,
		ClientMutationID:  strings.TrimSpace(request.ClientMutationID),
		TradeQuantity:     request.TradeQuantity,
		Notes:             request.Notes,
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

// TransferCollectionItem moves an existing item into another collection.
func (c *Client) TransferCollectionItem(ctx context.Context, request *TransferCollectionItemRequest, opts ...RequestOpt) (*TransferCollectionItemResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	itemID := strings.TrimSpace(request.ItemID)
	if itemID == "" {
		return nil, errors.New("itemID must not be empty")
	}
	destinationCollectionID := strings.TrimSpace(request.DestinationCollectionID)
	if destinationCollectionID == "" {
		return nil, errors.New("destinationCollectionID must not be empty")
	}

	body, err := jsonBody(struct {
		DestinationCollectionID string `json:"destinationCollectionId,omitempty"`
	}{
		DestinationCollectionID: destinationCollectionID,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/collection_items/%s:transfer", url.PathEscape(itemID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &TransferCollectionItemResponse{}
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

// BatchUpdateCollectionItems patches collection items from one source collection.
func (c *Client) BatchUpdateCollectionItems(ctx context.Context, request *BatchUpdateCollectionItemsRequest, opts ...RequestOpt) (*BatchUpdateCollectionItemsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	collectionID := strings.TrimSpace(request.CollectionID)
	if collectionID == "" {
		return nil, errors.New("collectionID must not be empty")
	}
	if len(request.Requests) == 0 {
		return nil, errors.New("requests must not be empty")
	}

	type updateWire struct {
		ItemID            string     `json:"itemId,omitempty"`
		Quantity          *int32     `json:"quantity,omitempty"`
		Condition         *Condition `json:"condition,omitempty"`
		Value             *float64   `json:"value,omitempty"`
		Pinned            *bool      `json:"pinned,omitempty"`
		ExpectedUpdatedAt *time.Time `json:"expectedUpdatedAt,omitempty"`
		ClientMutationID  string     `json:"clientMutationId,omitempty"`
		TradeQuantity     *int32     `json:"tradeQuantity,omitempty"`
		Notes             *string    `json:"notes,omitempty"`
	}
	wires := make([]updateWire, len(request.Requests))
	for index, item := range request.Requests {
		wires[index] = updateWire{ItemID: strings.TrimSpace(item.ItemID), Quantity: item.Quantity, Condition: item.Condition, Value: item.Value, Pinned: item.Pinned, ExpectedUpdatedAt: item.ExpectedUpdatedAt, ClientMutationID: strings.TrimSpace(item.ClientMutationID), TradeQuantity: item.TradeQuantity, Notes: item.Notes}
	}
	body, err := jsonBody(struct {
		CollectionID string       `json:"collectionId"`
		Requests     []updateWire `json:"requests"`
	}{CollectionID: collectionID, Requests: wires})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collection_items:batchUpdate", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	response := &BatchUpdateCollectionItemsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BatchDeleteCollectionItems deletes collection items from one source collection.
func (c *Client) BatchDeleteCollectionItems(ctx context.Context, request *BatchDeleteCollectionItemsRequest, opts ...RequestOpt) (*BatchDeleteCollectionItemsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CollectionID) == "" {
		return nil, errors.New("collectionID must not be empty")
	}
	if len(request.ItemIDs) == 0 {
		return nil, errors.New("itemIDs must not be empty")
	}
	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collection_items:batchDelete", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	response := &BatchDeleteCollectionItemsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BatchTransferCollectionItems transfers collection items between two collections.
func (c *Client) BatchTransferCollectionItems(ctx context.Context, request *BatchTransferCollectionItemsRequest, opts ...RequestOpt) (*BatchTransferCollectionItemsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.SourceCollectionID) == "" {
		return nil, errors.New("sourceCollectionID must not be empty")
	}
	if strings.TrimSpace(request.DestinationCollectionID) == "" {
		return nil, errors.New("destinationCollectionID must not be empty")
	}
	if len(request.ItemIDs) == 0 {
		return nil, errors.New("itemIDs must not be empty")
	}
	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}
	req, err := c.newRequest(ctx, http.MethodPost, "/v1/collection_items:batchTransfer", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	response := &BatchTransferCollectionItemsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
