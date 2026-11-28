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

// TourStop represents events.v1.TourStop.
// StartsOn and EndsOn are YYYY-MM-DD dates in the stop's local timezone.
type TourStop struct {
	TourStopID    string      `json:"tourStopId,omitempty"`
	Title         string      `json:"title,omitempty"`
	StartsAt      *time.Time  `json:"startsAt,omitempty"`
	EndsAt        *time.Time  `json:"endsAt,omitempty"`
	VenueName     string      `json:"venueName,omitempty"`
	Country       string      `json:"country,omitempty"`
	HostName      string      `json:"hostName,omitempty"`
	LivestreamURL string      `json:"livestreamUrl,omitempty"`
	CreatedAt     *time.Time  `json:"createdAt,omitempty"`
	UpdatedAt     *time.Time  `json:"updatedAt,omitempty"`
	LastSeenAt    *time.Time  `json:"lastSeenAt,omitempty"`
	LocationImage *EventImage `json:"locationImage,omitempty"`
	CatalogStopID string      `json:"catalogStopId,omitempty"`
	StartsOn      string      `json:"startsOn,omitempty"`
	EndsOn        string      `json:"endsOn,omitempty"`
	Status        EventStatus `json:"status,omitempty"`
	Timezone      string      `json:"timezone,omitempty"`
	ProgramURL    string      `json:"programUrl,omitempty"`
}

// ListTourStopsRequest configures tour stop listing and pagination.
type ListTourStopsRequest struct {
	PageSize  *int32
	NextToken string
	Year      *int32
}

// ListTourStopsResponse returns tour stops.
type ListTourStopsResponse struct {
	TourStops []TourStop       `json:"tourStops,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListTourStopsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetTourStopRequest identifies a tour stop.
type GetTourStopRequest struct {
	TourStopID string `json:"-"`
}

// GetTourStopResponse returns one tour stop and its events.
type GetTourStopResponse struct {
	TourStop *TourStop        `json:"tourStop,omitempty"`
	Events   []Event          `json:"events,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetTourStopResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListTourStops lists organized-play tour stops.
func (c *Client) ListTourStops(ctx context.Context, request *ListTourStopsRequest, opts ...RequestOpt) (*ListTourStopsResponse, error) {
	if request == nil {
		request = &ListTourStopsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/tour-stops", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryInt32(query, "pageSize", request.PageSize)
	setQueryString(query, "nextToken", request.NextToken)
	setQueryInt32(query, "year", request.Year)
	req.URL.RawQuery = query.Encode()

	response := &ListTourStopsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetTourStop fetches one tour stop with its events.
func (c *Client) GetTourStop(ctx context.Context, request *GetTourStopRequest, opts ...RequestOpt) (*GetTourStopResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	tourStopID := strings.TrimSpace(request.TourStopID)
	if tourStopID == "" {
		return nil, errors.New("tourStopID must not be empty")
	}

	path := fmt.Sprintf("/v1/tour-stops/%s", url.PathEscape(tourStopID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetTourStopResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
