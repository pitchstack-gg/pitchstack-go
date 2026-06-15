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

// PriceEntry mirrors v1PriceEntry.
type PriceEntry struct {
	EntryID   string     `json:"entryId,omitempty"`
	ProductID string     `json:"productId,omitempty"`
	Currency  string     `json:"currency,omitempty"`
	Price     float64    `json:"price,omitempty"`
	Price2    float64    `json:"price2,omitempty"`
	Price3    float64    `json:"price3,omitempty"`
	RecordAt  *time.Time `json:"recordAt,omitempty"`
	Source    string     `json:"source,omitempty"`
	SourceURL string     `json:"sourceUrl,omitempty"`
}

// GetProductPriceRequest identifies which product price to fetch.
type GetProductPriceRequest struct {
	ProductID string `json:"-"`
	Source    string
}

// GetProductPriceResponse contains the current price.
type GetProductPriceResponse struct {
	Entry    *PriceEntry      `json:"entry,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetProductPriceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetProductPriceHistoryRequest fetches historical pricing.
type GetProductPriceHistoryRequest struct {
	ProductID string `json:"-"`
	StartDate string
	EndDate   string
	Limit     *int32
	Source    string
}

// GetProductPriceHistoryResponse mirrors v1GetProductPriceHistoryResponse.
type GetProductPriceHistoryResponse struct {
	ProductID string           `json:"productId,omitempty"`
	Source    string           `json:"source,omitempty"`
	SourceURL string           `json:"sourceUrl,omitempty"`
	StartTime *time.Time       `json:"startTime,omitempty"`
	EndTime   *time.Time       `json:"endTime,omitempty"`
	Entries   []PriceEntry     `json:"entries,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *GetProductPriceHistoryResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetProductPricesRequest fetches prices for multiple products.
type BatchGetProductPricesRequest struct {
	ProductIDs []string `json:"productIds,omitempty"`
	Source     string   `json:"source,omitempty"`
}

// BatchGetProductPricesResponse mirrors v1BatchGetProductPricesResponse.
type BatchGetProductPricesResponse struct {
	Prices   []PriceEntry     `json:"prices,omitempty"`
	NotFound []string         `json:"notFound,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *BatchGetProductPricesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ProductPriceWatch represents a price watch for a product/source pair.
type ProductPriceWatch struct {
	WatchID                  string     `json:"watchId,omitempty"`
	ProductID                string     `json:"productId,omitempty"`
	Source                   string     `json:"source,omitempty"`
	Direction                string     `json:"direction,omitempty"`
	AbsoluteChange           *float64   `json:"absoluteChange,omitempty"`
	PercentChange            *float64   `json:"percentChange,omitempty"`
	Period                   string     `json:"period,omitempty"`
	Active                   bool       `json:"active,omitempty"`
	CreatedAt                *time.Time `json:"createdAt,omitempty"`
	UpdatedAt                *time.Time `json:"updatedAt,omitempty"`
	LastNotifiedAt           *time.Time `json:"lastNotifiedAt,omitempty"`
	LastNotifiedDirection    string     `json:"lastNotifiedDirection,omitempty"`
	LastNotifiedPriceEntryID string     `json:"lastNotifiedPriceEntryId,omitempty"`
}

// ProductPriceWatchList represents a named group of price watches.
type ProductPriceWatchList struct {
	ListID      string     `json:"listId,omitempty"`
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	IsDefault   bool       `json:"isDefault,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
}

// ProductPriceWatchListItem links a price watch to a watch list.
type ProductPriceWatchListItem struct {
	ItemID    string             `json:"itemId,omitempty"`
	ListID    string             `json:"listId,omitempty"`
	Watch     *ProductPriceWatch `json:"watch,omitempty"`
	CreatedAt *time.Time         `json:"createdAt,omitempty"`
}

// BatchAddProductPriceWatchFailure describes a product that could not be added.
type BatchAddProductPriceWatchFailure struct {
	ProductID string `json:"productId,omitempty"`
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
}

// CreateProductPriceWatchRequest creates or replaces a price watch.
type CreateProductPriceWatchRequest struct {
	ProductID      string   `json:"productId,omitempty"`
	Source         string   `json:"source,omitempty"`
	Direction      string   `json:"direction,omitempty"`
	AbsoluteChange *float64 `json:"absoluteChange,omitempty"`
	PercentChange  *float64 `json:"percentChange,omitempty"`
	Period         string   `json:"period,omitempty"`
}

// CreateProductPriceWatchResponse returns the created price watch.
type CreateProductPriceWatchResponse struct {
	Watch    *ProductPriceWatch `json:"watch,omitempty"`
	Metadata ResponseMetadata   `json:"-"`
}

func (r *CreateProductPriceWatchResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateProductPriceWatchRequest updates a price watch.
type UpdateProductPriceWatchRequest struct {
	WatchID        string   `json:"-"`
	Direction      *string  `json:"direction,omitempty"`
	AbsoluteChange *float64 `json:"absoluteChange,omitempty"`
	PercentChange  *float64 `json:"percentChange,omitempty"`
	Period         *string  `json:"period,omitempty"`
	Active         *bool    `json:"active,omitempty"`
}

// UpdateProductPriceWatchResponse returns the updated price watch.
type UpdateProductPriceWatchResponse struct {
	Watch    *ProductPriceWatch `json:"watch,omitempty"`
	Metadata ResponseMetadata   `json:"-"`
}

func (r *UpdateProductPriceWatchResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteProductPriceWatchRequest identifies a price watch to delete.
type DeleteProductPriceWatchRequest struct {
	WatchID string `json:"-"`
}

// DeleteProductPriceWatchResponse captures metadata for delete operations.
type DeleteProductPriceWatchResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteProductPriceWatchResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListProductPriceWatchesRequest filters price watches.
type ListProductPriceWatchesRequest struct {
	ActiveOnly *bool
	ProductIDs []string
}

// ListProductPriceWatchesResponse returns price watches.
type ListProductPriceWatchesResponse struct {
	Watches  []ProductPriceWatch `json:"watches,omitempty"`
	Metadata ResponseMetadata    `json:"-"`
}

func (r *ListProductPriceWatchesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListProductPriceWatchListsRequest lists the authenticated user's watch lists.
type ListProductPriceWatchListsRequest struct{}

// ListProductPriceWatchListsResponse returns price watch lists.
type ListProductPriceWatchListsResponse struct {
	Lists    []ProductPriceWatchList `json:"lists,omitempty"`
	Metadata ResponseMetadata        `json:"-"`
}

func (r *ListProductPriceWatchListsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateProductPriceWatchListRequest creates a price watch list.
type CreateProductPriceWatchListRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// CreateProductPriceWatchListResponse returns the created watch list.
type CreateProductPriceWatchListResponse struct {
	List     *ProductPriceWatchList `json:"list,omitempty"`
	Metadata ResponseMetadata       `json:"-"`
}

func (r *CreateProductPriceWatchListResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateProductPriceWatchListRequest updates a price watch list.
type UpdateProductPriceWatchListRequest struct {
	ListID      string  `json:"-"`
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// UpdateProductPriceWatchListResponse returns the updated watch list.
type UpdateProductPriceWatchListResponse struct {
	List     *ProductPriceWatchList `json:"list,omitempty"`
	Metadata ResponseMetadata       `json:"-"`
}

func (r *UpdateProductPriceWatchListResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteProductPriceWatchListRequest identifies a price watch list to delete.
type DeleteProductPriceWatchListRequest struct {
	ListID string `json:"-"`
}

// DeleteProductPriceWatchListResponse captures metadata for delete operations.
type DeleteProductPriceWatchListResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteProductPriceWatchListResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListProductPriceWatchListItemsRequest lists items for a price watch list.
type ListProductPriceWatchListItemsRequest struct {
	ListID     string `json:"-"`
	ActiveOnly *bool
}

// ListProductPriceWatchListItemsResponse returns price watch list items.
type ListProductPriceWatchListItemsResponse struct {
	Items    []ProductPriceWatchListItem `json:"items,omitempty"`
	Metadata ResponseMetadata            `json:"-"`
}

func (r *ListProductPriceWatchListItemsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// AddProductPriceWatchListItemRequest adds a watch or product watch config to a list.
type AddProductPriceWatchListItemRequest struct {
	ListID         string   `json:"-"`
	WatchID        string   `json:"watchId,omitempty"`
	ProductID      string   `json:"productId,omitempty"`
	Source         string   `json:"source,omitempty"`
	Direction      string   `json:"direction,omitempty"`
	AbsoluteChange *float64 `json:"absoluteChange,omitempty"`
	PercentChange  *float64 `json:"percentChange,omitempty"`
	Period         string   `json:"period,omitempty"`
}

// AddProductPriceWatchListItemResponse returns the added list item.
type AddProductPriceWatchListItemResponse struct {
	Item     *ProductPriceWatchListItem `json:"item,omitempty"`
	Metadata ResponseMetadata           `json:"-"`
}

func (r *AddProductPriceWatchListItemResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RemoveProductPriceWatchListItemRequest identifies a watch to remove from a list.
type RemoveProductPriceWatchListItemRequest struct {
	ListID  string `json:"-"`
	WatchID string `json:"-"`
}

// RemoveProductPriceWatchListItemResponse captures metadata for remove operations.
type RemoveProductPriceWatchListItemResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *RemoveProductPriceWatchListItemResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchAddProductsToProductPriceWatchListRequest adds many product watches to a list.
type BatchAddProductsToProductPriceWatchListRequest struct {
	ListID         string   `json:"-"`
	ProductIDs     []string `json:"productIds,omitempty"`
	Source         string   `json:"source,omitempty"`
	Direction      string   `json:"direction,omitempty"`
	AbsoluteChange *float64 `json:"absoluteChange,omitempty"`
	PercentChange  *float64 `json:"percentChange,omitempty"`
	Period         string   `json:"period,omitempty"`
}

// BatchAddProductsToProductPriceWatchListResponse returns added items and per-product failures.
type BatchAddProductsToProductPriceWatchListResponse struct {
	Items    []ProductPriceWatchListItem        `json:"items,omitempty"`
	Failures []BatchAddProductPriceWatchFailure `json:"failures,omitempty"`
	Metadata ResponseMetadata                   `json:"-"`
}

func (r *BatchAddProductsToProductPriceWatchListResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetPricingStatsResponse mirrors v1GetPricingStatsResponse.
type GetPricingStatsResponse struct {
	TotalCards        int32            `json:"totalCards,omitempty"`
	UpdatedLast24h    int32            `json:"updatedLast24h,omitempty"`
	AveragePrice      float64          `json:"averagePrice,omitempty"`
	MedianPrice       float64          `json:"medianPrice,omitempty"`
	TotalValue        float64          `json:"totalValue,omitempty"`
	HighestPricedCard *PriceEntry      `json:"highestPricedCard,omitempty"`
	BiggestGainer     *PriceEntry      `json:"biggestGainer,omitempty"`
	BiggestLoser      *PriceEntry      `json:"biggestLoser,omitempty"`
	LastUpdateTime    *time.Time       `json:"lastUpdateTime,omitempty"`
	Metadata          ResponseMetadata `json:"-"`
}

func (r *GetPricingStatsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetProductPrice retrieves the current price entry for a product.
func (c *Client) GetProductPrice(ctx context.Context, request *GetProductPriceRequest, opts ...RequestOpt) (*GetProductPriceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ProductID) == "" {
		return nil, errors.New("productID must not be empty")
	}

	path := fmt.Sprintf("/v1/prices/%s", url.PathEscape(request.ProductID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "source", request.Source)
	req.URL.RawQuery = query.Encode()

	response := &GetProductPriceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetProductPriceHistory retrieves pricing history for a product.
func (c *Client) GetProductPriceHistory(ctx context.Context, request *GetProductPriceHistoryRequest, opts ...RequestOpt) (*GetProductPriceHistoryResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ProductID) == "" {
		return nil, errors.New("productID must not be empty")
	}

	path := fmt.Sprintf("/v1/prices/%s/history", url.PathEscape(request.ProductID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "startDate", request.StartDate)
	setQueryString(query, "endDate", request.EndDate)
	if request.Limit != nil && *request.Limit > 0 {
		query.Set("limit", strconv.Itoa(int(*request.Limit)))
	}
	setQueryString(query, "source", request.Source)
	req.URL.RawQuery = query.Encode()

	response := &GetProductPriceHistoryResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchGetProductPrices retrieves prices for multiple products.
func (c *Client) BatchGetProductPrices(ctx context.Context, request *BatchGetProductPricesRequest, opts ...RequestOpt) (*BatchGetProductPricesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.ProductIDs) == 0 {
		return nil, errors.New("productIDs must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/prices:batchGet", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &BatchGetProductPricesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateProductPriceWatch creates or replaces a price watch.
func (c *Client) CreateProductPriceWatch(ctx context.Context, request *CreateProductPriceWatchRequest, opts ...RequestOpt) (*CreateProductPriceWatchResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ProductID) == "" {
		return nil, errors.New("productID must not be empty")
	}
	if strings.TrimSpace(request.Source) == "" {
		return nil, errors.New("source must not be empty")
	}
	if strings.TrimSpace(request.Direction) == "" {
		return nil, errors.New("direction must not be empty")
	}
	if strings.TrimSpace(request.Period) == "" {
		return nil, errors.New("period must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/price-watches", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CreateProductPriceWatchResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpdateProductPriceWatch updates a price watch.
func (c *Client) UpdateProductPriceWatch(ctx context.Context, request *UpdateProductPriceWatchRequest, opts ...RequestOpt) (*UpdateProductPriceWatchResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	watchID := strings.TrimSpace(request.WatchID)
	if watchID == "" {
		return nil, errors.New("watchID must not be empty")
	}

	body, err := jsonBody(struct {
		Direction      *string  `json:"direction,omitempty"`
		AbsoluteChange *float64 `json:"absoluteChange,omitempty"`
		PercentChange  *float64 `json:"percentChange,omitempty"`
		Period         *string  `json:"period,omitempty"`
		Active         *bool    `json:"active,omitempty"`
	}{
		Direction:      request.Direction,
		AbsoluteChange: request.AbsoluteChange,
		PercentChange:  request.PercentChange,
		Period:         request.Period,
		Active:         request.Active,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/price-watches/%s", url.PathEscape(watchID))
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdateProductPriceWatchResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// DeleteProductPriceWatch deletes a price watch.
func (c *Client) DeleteProductPriceWatch(ctx context.Context, request *DeleteProductPriceWatchRequest, opts ...RequestOpt) (*DeleteProductPriceWatchResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	watchID := strings.TrimSpace(request.WatchID)
	if watchID == "" {
		return nil, errors.New("watchID must not be empty")
	}

	path := fmt.Sprintf("/v1/price-watches/%s", url.PathEscape(watchID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteProductPriceWatchResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListProductPriceWatches lists price watches for the authenticated user.
func (c *Client) ListProductPriceWatches(ctx context.Context, request *ListProductPriceWatchesRequest, opts ...RequestOpt) (*ListProductPriceWatchesResponse, error) {
	if request == nil {
		request = &ListProductPriceWatchesRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/price-watches", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryBool(query, "activeOnly", request.ActiveOnly)
	for _, productID := range request.ProductIDs {
		if productID = strings.TrimSpace(productID); productID != "" {
			query.Add("productIds", productID)
		}
	}
	req.URL.RawQuery = query.Encode()

	response := &ListProductPriceWatchesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListProductPriceWatchLists lists price watch lists for the authenticated user.
func (c *Client) ListProductPriceWatchLists(ctx context.Context, request *ListProductPriceWatchListsRequest, opts ...RequestOpt) (*ListProductPriceWatchListsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/price-watch-lists", nil)
	if err != nil {
		return nil, err
	}

	response := &ListProductPriceWatchListsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CreateProductPriceWatchList creates a price watch list.
func (c *Client) CreateProductPriceWatchList(ctx context.Context, request *CreateProductPriceWatchListRequest, opts ...RequestOpt) (*CreateProductPriceWatchListResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/price-watch-lists", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CreateProductPriceWatchListResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpdateProductPriceWatchList updates a price watch list.
func (c *Client) UpdateProductPriceWatchList(ctx context.Context, request *UpdateProductPriceWatchListRequest, opts ...RequestOpt) (*UpdateProductPriceWatchListResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	listID := strings.TrimSpace(request.ListID)
	if listID == "" {
		return nil, errors.New("listID must not be empty")
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

	path := fmt.Sprintf("/v1/price-watch-lists/%s", url.PathEscape(listID))
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdateProductPriceWatchListResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// DeleteProductPriceWatchList deletes a price watch list.
func (c *Client) DeleteProductPriceWatchList(ctx context.Context, request *DeleteProductPriceWatchListRequest, opts ...RequestOpt) (*DeleteProductPriceWatchListResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	listID := strings.TrimSpace(request.ListID)
	if listID == "" {
		return nil, errors.New("listID must not be empty")
	}

	path := fmt.Sprintf("/v1/price-watch-lists/%s", url.PathEscape(listID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteProductPriceWatchListResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListProductPriceWatchListItems lists items for a price watch list.
func (c *Client) ListProductPriceWatchListItems(ctx context.Context, request *ListProductPriceWatchListItemsRequest, opts ...RequestOpt) (*ListProductPriceWatchListItemsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	listID := strings.TrimSpace(request.ListID)
	if listID == "" {
		return nil, errors.New("listID must not be empty")
	}

	path := fmt.Sprintf("/v1/price-watch-lists/%s/items", url.PathEscape(listID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryBool(query, "activeOnly", request.ActiveOnly)
	req.URL.RawQuery = query.Encode()

	response := &ListProductPriceWatchListItemsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// AddProductPriceWatchListItem adds a watch or product watch config to a list.
func (c *Client) AddProductPriceWatchListItem(ctx context.Context, request *AddProductPriceWatchListItemRequest, opts ...RequestOpt) (*AddProductPriceWatchListItemResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	listID := strings.TrimSpace(request.ListID)
	if listID == "" {
		return nil, errors.New("listID must not be empty")
	}

	body, err := jsonBody(struct {
		WatchID        string   `json:"watchId,omitempty"`
		ProductID      string   `json:"productId,omitempty"`
		Source         string   `json:"source,omitempty"`
		Direction      string   `json:"direction,omitempty"`
		AbsoluteChange *float64 `json:"absoluteChange,omitempty"`
		PercentChange  *float64 `json:"percentChange,omitempty"`
		Period         string   `json:"period,omitempty"`
	}{
		WatchID:        strings.TrimSpace(request.WatchID),
		ProductID:      strings.TrimSpace(request.ProductID),
		Source:         strings.TrimSpace(request.Source),
		Direction:      strings.TrimSpace(request.Direction),
		AbsoluteChange: request.AbsoluteChange,
		PercentChange:  request.PercentChange,
		Period:         strings.TrimSpace(request.Period),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/price-watch-lists/%s/items", url.PathEscape(listID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &AddProductPriceWatchListItemResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// RemoveProductPriceWatchListItem removes a watch from a price watch list.
func (c *Client) RemoveProductPriceWatchListItem(ctx context.Context, request *RemoveProductPriceWatchListItemRequest, opts ...RequestOpt) (*RemoveProductPriceWatchListItemResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	listID := strings.TrimSpace(request.ListID)
	if listID == "" {
		return nil, errors.New("listID must not be empty")
	}
	watchID := strings.TrimSpace(request.WatchID)
	if watchID == "" {
		return nil, errors.New("watchID must not be empty")
	}

	path := fmt.Sprintf("/v1/price-watch-lists/%s/items/%s", url.PathEscape(listID), url.PathEscape(watchID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &RemoveProductPriceWatchListItemResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BatchAddProductsToProductPriceWatchList adds many product watches to a list.
func (c *Client) BatchAddProductsToProductPriceWatchList(ctx context.Context, request *BatchAddProductsToProductPriceWatchListRequest, opts ...RequestOpt) (*BatchAddProductsToProductPriceWatchListResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	listID := strings.TrimSpace(request.ListID)
	if listID == "" {
		return nil, errors.New("listID must not be empty")
	}
	if len(request.ProductIDs) == 0 {
		return nil, errors.New("productIDs must not be empty")
	}

	body, err := jsonBody(struct {
		ProductIDs     []string `json:"productIds,omitempty"`
		Source         string   `json:"source,omitempty"`
		Direction      string   `json:"direction,omitempty"`
		AbsoluteChange *float64 `json:"absoluteChange,omitempty"`
		PercentChange  *float64 `json:"percentChange,omitempty"`
		Period         string   `json:"period,omitempty"`
	}{
		ProductIDs:     request.ProductIDs,
		Source:         strings.TrimSpace(request.Source),
		Direction:      strings.TrimSpace(request.Direction),
		AbsoluteChange: request.AbsoluteChange,
		PercentChange:  request.PercentChange,
		Period:         strings.TrimSpace(request.Period),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/price-watch-lists/%s/items:batchAddProducts", url.PathEscape(listID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BatchAddProductsToProductPriceWatchListResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
