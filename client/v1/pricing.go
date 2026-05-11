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
