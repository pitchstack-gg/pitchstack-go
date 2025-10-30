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
