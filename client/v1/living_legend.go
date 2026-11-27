package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// LivingLegendStanding represents cards.v1.LivingLegendStanding.
type LivingLegendStanding struct {
	SnapshotDate   string     `json:"snapshotDate,omitempty"`
	Format         string     `json:"format,omitempty"`
	HeroName       string     `json:"heroName,omitempty"`
	HeroCardID     string     `json:"heroCardId,omitempty"`
	Points         int32      `json:"points,omitempty"`
	Rank           *int32     `json:"rank,omitempty"`
	SeasonCode     string     `json:"seasonCode,omitempty"`
	SeasonPoints   *int32     `json:"seasonPoints,omitempty"`
	IsLivingLegend bool       `json:"isLivingLegend,omitempty"`
	Source         string     `json:"source,omitempty"`
	CapturedAt     *time.Time `json:"capturedAt,omitempty"`
}

// ListLivingLegendStandingsRequest filters Living Legend standings.
// SnapshotDate is a YYYY-MM-DD date; empty selects the latest snapshot.
type ListLivingLegendStandingsRequest struct {
	Format         string
	SnapshotDate   string
	IncludeRetired *bool
	PageSize       *int32
	NextToken      string
}

// ListLivingLegendStandingsResponse returns Living Legend standings.
type ListLivingLegendStandingsResponse struct {
	Standings []LivingLegendStanding `json:"standings,omitempty"`
	NextToken string                 `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata       `json:"-"`
}

func (r *ListLivingLegendStandingsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListHeroLivingLegendHistoryRequest filters one hero's Living Legend history.
// StartDate and EndDate are YYYY-MM-DD dates.
type ListHeroLivingLegendHistoryRequest struct {
	HeroCardID string `json:"-"`
	Format     string
	StartDate  string
	EndDate    string
	PageSize   *int32
	NextToken  string
}

// ListHeroLivingLegendHistoryResponse returns one hero's Living Legend history.
type ListHeroLivingLegendHistoryResponse struct {
	Standings []LivingLegendStanding `json:"standings,omitempty"`
	NextToken string                 `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata       `json:"-"`
}

func (r *ListHeroLivingLegendHistoryResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListLivingLegendStandings lists Living Legend point standings.
func (c *Client) ListLivingLegendStandings(ctx context.Context, request *ListLivingLegendStandingsRequest, opts ...RequestOpt) (*ListLivingLegendStandingsResponse, error) {
	if request == nil {
		request = &ListLivingLegendStandingsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/living-legend/standings", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "format", request.Format)
	setQueryString(query, "snapshotDate", request.SnapshotDate)
	setQueryBool(query, "includeRetired", request.IncludeRetired)
	setQueryInt32(query, "pageSize", request.PageSize)
	setQueryString(query, "nextToken", request.NextToken)
	req.URL.RawQuery = query.Encode()

	response := &ListLivingLegendStandingsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListHeroLivingLegendHistory lists Living Legend point history for one hero.
func (c *Client) ListHeroLivingLegendHistory(ctx context.Context, request *ListHeroLivingLegendHistoryRequest, opts ...RequestOpt) (*ListHeroLivingLegendHistoryResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	heroCardID := strings.TrimSpace(request.HeroCardID)
	if heroCardID == "" {
		return nil, errors.New("heroCardID must not be empty")
	}

	path := fmt.Sprintf("/v1/living-legend/heroes/%s/history", url.PathEscape(heroCardID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "format", request.Format)
	setQueryString(query, "startDate", request.StartDate)
	setQueryString(query, "endDate", request.EndDate)
	setQueryInt32(query, "pageSize", request.PageSize)
	setQueryString(query, "nextToken", request.NextToken)
	req.URL.RawQuery = query.Encode()

	response := &ListHeroLivingLegendHistoryResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
