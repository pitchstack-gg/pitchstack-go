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
	EntryID        string     `json:"entryId,omitempty"`
	PhysicalCardID string     `json:"physicalCardId,omitempty"`
	Currency       string     `json:"currency,omitempty"`
	Price          float64    `json:"price,omitempty"`
	LowPrice       float64    `json:"lowPrice,omitempty"`
	HighPrice      float64    `json:"highPrice,omitempty"`
	RecordAt       *time.Time `json:"recordAt,omitempty"`
	Source         string     `json:"source,omitempty"`
	SourceURL      string     `json:"sourceUrl,omitempty"`
}

// GetPhysicalCardPriceRequest identifies which card price to fetch.
type GetPhysicalCardPriceRequest struct {
	PhysicalCardID string `json:"-"`
	Source         string
	Currency       string
}

// GetPhysicalCardPriceResponse contains the current price.
type GetPhysicalCardPriceResponse struct {
	Entry    *PriceEntry      `json:"entry,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetPhysicalCardPriceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetPhysicalCardPriceHistoryRequest fetches historical pricing.
type GetPhysicalCardPriceHistoryRequest struct {
	PhysicalCardID string `json:"-"`
	StartDate      string
	EndDate        string
	Limit          *int32
	Source         string
}

// GetPhysicalCardPriceHistoryResponse mirrors v1GetPhysicalCardPriceHistoryResponse.
type GetPhysicalCardPriceHistoryResponse struct {
	PhysicalCardID string           `json:"physicalCardId,omitempty"`
	Source         string           `json:"source,omitempty"`
	SourceURL      string           `json:"sourceUrl,omitempty"`
	StartTime      *time.Time       `json:"startTime,omitempty"`
	EndTime        *time.Time       `json:"endTime,omitempty"`
	Entries        []PriceEntry     `json:"entries,omitempty"`
	Metadata       ResponseMetadata `json:"-"`
}

func (r *GetPhysicalCardPriceHistoryResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetBulkPhysicalCardPricesRequest fetches prices for multiple cards.
type GetBulkPhysicalCardPricesRequest struct {
	PhysicalCardIDs []string `json:"physicalCardIds,omitempty"`
	Source          string   `json:"source,omitempty"`
}

// GetBulkPhysicalCardPricesResponse mirrors v1GetBulkPhysicalCardPricesResponse.
type GetBulkPhysicalCardPricesResponse struct {
	Prices   []PriceEntry     `json:"prices,omitempty"`
	NotFound []string         `json:"notFound,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetBulkPhysicalCardPricesResponse) setMetadata(metadata ResponseMetadata) {
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

// GetPhysicalCardPrice retrieves the current price entry for a physical card.
func (c *Client) GetPhysicalCardPrice(ctx context.Context, request *GetPhysicalCardPriceRequest, opts ...RequestOpt) (*GetPhysicalCardPriceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.PhysicalCardID) == "" {
		return nil, errors.New("physicalCardID must not be empty")
	}

	path := fmt.Sprintf("/api/v1/physical_cards/%s/price", url.PathEscape(request.PhysicalCardID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "source", request.Source)
	setQueryString(query, "currency", request.Currency)
	req.URL.RawQuery = query.Encode()

	response := &GetPhysicalCardPriceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetPhysicalCardPriceHistory retrieves pricing history for a card.
func (c *Client) GetPhysicalCardPriceHistory(ctx context.Context, request *GetPhysicalCardPriceHistoryRequest, opts ...RequestOpt) (*GetPhysicalCardPriceHistoryResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.PhysicalCardID) == "" {
		return nil, errors.New("physicalCardID must not be empty")
	}

	path := fmt.Sprintf("/api/v1/physical_cards/%s/price/history", url.PathEscape(request.PhysicalCardID))
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

	response := &GetPhysicalCardPriceHistoryResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetBulkPhysicalCardPrices retrieves prices for multiple cards.
func (c *Client) GetBulkPhysicalCardPrices(ctx context.Context, request *GetBulkPhysicalCardPricesRequest, opts ...RequestOpt) (*GetBulkPhysicalCardPricesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.PhysicalCardIDs) == 0 {
		return nil, errors.New("physicalCardIDs must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/api/v1/prices:batchGet", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &GetBulkPhysicalCardPricesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetPricingStats retrieves aggregate pricing statistics.
func (c *Client) GetPricingStats(ctx context.Context, opts ...RequestOpt) (*GetPricingStatsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/prices:stats", nil)
	if err != nil {
		return nil, err
	}

	response := &GetPricingStatsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}
