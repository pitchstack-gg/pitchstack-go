package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// TrendingWindow represents engagement.v1.TrendingWindow.
type TrendingWindow string

const (
	TrendingWindowUnspecified TrendingWindow = "TRENDING_WINDOW_UNSPECIFIED"
	TrendingWindow24H         TrendingWindow = "TRENDING_WINDOW_24H"
	TrendingWindow7D          TrendingWindow = "TRENDING_WINDOW_7D"
	TrendingWindow30D         TrendingWindow = "TRENDING_WINDOW_30D"
)

// TrackableResourceType represents engagement.v1.TrackableResourceType.
type TrackableResourceType string

const (
	TrackableResourceTypeUnspecified TrackableResourceType = "TRACKABLE_RESOURCE_TYPE_UNSPECIFIED"
	TrackableResourceTypeCard        TrackableResourceType = "TRACKABLE_RESOURCE_TYPE_CARD"
	TrackableResourceTypeDeck        TrackableResourceType = "TRACKABLE_RESOURCE_TYPE_DECK"
	TrackableResourceTypeCollection  TrackableResourceType = "TRACKABLE_RESOURCE_TYPE_COLLECTION"
	TrackableResourceTypeUserProfile TrackableResourceType = "TRACKABLE_RESOURCE_TYPE_USER_PROFILE"
)

// EngagementResourceRef mirrors engagement.v1.ResourceRef.
type EngagementResourceRef struct {
	ResourceType TrackableResourceType `json:"resourceType,omitempty"`
	ResourceID   string                `json:"resourceId,omitempty"`
}

// TrackViewRequest mirrors engagement.v1.TrackViewRequest.
type TrackViewRequest struct {
	Resource       *EngagementResourceRef `json:"resource,omitempty"`
	ClientViewerID string                 `json:"clientViewerId,omitempty"`
	OccurredAt     *time.Time             `json:"occurredAt,omitempty"`
}

// TrackViewResponse mirrors engagement.v1.TrackViewResponse.
type TrackViewResponse struct {
	Counted         bool             `json:"counted,omitempty"`
	Deduped         bool             `json:"deduped,omitempty"`
	SkippedSelfView bool             `json:"skippedSelfView,omitempty"`
	Metadata        ResponseMetadata `json:"-"`
}

func (r *TrackViewResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchTrackViewsRequest mirrors engagement.v1.BatchTrackViewsRequest.
type BatchTrackViewsRequest struct {
	Views []TrackViewRequest `json:"views,omitempty"`
}

// BatchTrackViewResult mirrors engagement.v1.BatchTrackViewResult.
type BatchTrackViewResult struct {
	Resource        *EngagementResourceRef `json:"resource,omitempty"`
	Counted         bool                   `json:"counted,omitempty"`
	Deduped         bool                   `json:"deduped,omitempty"`
	SkippedSelfView bool                   `json:"skippedSelfView,omitempty"`
}

// BatchTrackViewsResponse mirrors engagement.v1.BatchTrackViewsResponse.
type BatchTrackViewsResponse struct {
	Results              []BatchTrackViewResult `json:"results,omitempty"`
	TotalCounted         int32                  `json:"totalCounted,omitempty"`
	TotalDeduped         int32                  `json:"totalDeduped,omitempty"`
	TotalSkippedSelfView int32                  `json:"totalSkippedSelfView,omitempty"`
	Metadata             ResponseMetadata       `json:"-"`
}

func (r *BatchTrackViewsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListTrendingResourcesRequest mirrors engagement.v1.ListTrendingResourcesRequest.
type ListTrendingResourcesRequest struct {
	ResourceType  TrackableResourceType `json:"resourceType,omitempty"`
	Window        TrendingWindow        `json:"window,omitempty"`
	PageSize      int32                 `json:"pageSize,omitempty"`
	NextPageToken string                `json:"nextPageToken,omitempty"`
}

// TrendingResource mirrors engagement.v1.TrendingResource.
type TrendingResource struct {
	Resource     *EngagementResourceRef `json:"resource,omitempty"`
	ViewCount    int64                  `json:"viewCount,omitempty"`
	Score        float64                `json:"score,omitempty"`
	LastViewedAt *time.Time             `json:"lastViewedAt,omitempty"`
}

// ListTrendingResourcesResponse mirrors engagement.v1.ListTrendingResourcesResponse.
type ListTrendingResourcesResponse struct {
	Resources     []TrendingResource `json:"resources,omitempty"`
	NextPageToken string             `json:"nextPageToken,omitempty"`
	Metadata      ResponseMetadata   `json:"-"`
}

func (r *ListTrendingResourcesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetViewCountsRequest mirrors engagement.v1.BatchGetViewCountsRequest.
type BatchGetViewCountsRequest struct {
	Resources []EngagementResourceRef `json:"resources,omitempty"`
}

// ResourceViewCount mirrors engagement.v1.ResourceViewCount.
type ResourceViewCount struct {
	Resource     *EngagementResourceRef `json:"resource,omitempty"`
	TotalViews   int64                  `json:"totalViews,omitempty"`
	LastViewedAt *time.Time             `json:"lastViewedAt,omitempty"`
}

// BatchGetViewCountsResponse mirrors engagement.v1.BatchGetViewCountsResponse.
type BatchGetViewCountsResponse struct {
	Counts   []ResourceViewCount `json:"counts,omitempty"`
	Metadata ResponseMetadata    `json:"-"`
}

func (r *BatchGetViewCountsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// TrackView records a single resource view event.
func (c *Client) TrackView(ctx context.Context, request *TrackViewRequest, opts ...RequestOpt) (*TrackViewResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := validateEngagementResourceRef(request.Resource); err != nil {
		return nil, err
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/engagement/views:track", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &TrackViewResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchTrackViews records multiple view events.
func (c *Client) BatchTrackViews(ctx context.Context, request *BatchTrackViewsRequest, opts ...RequestOpt) (*BatchTrackViewsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.Views) == 0 {
		return nil, errors.New("views must not be empty")
	}
	for i := range request.Views {
		if err := validateEngagementResourceRef(request.Views[i].Resource); err != nil {
			return nil, fmt.Errorf("views[%d]: %w", i, err)
		}
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/engagement/views:batchTrack", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BatchTrackViewsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListTrendingResources returns ranked resources for a requested window.
func (c *Client) ListTrendingResources(ctx context.Context, request *ListTrendingResourcesRequest, opts ...RequestOpt) (*ListTrendingResourcesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(string(request.ResourceType)) == "" || request.ResourceType == TrackableResourceTypeUnspecified {
		return nil, errors.New("resourceType must be specified")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/engagement/trending:list", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &ListTrendingResourcesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchGetViewCounts retrieves aggregate view counts for multiple resources.
func (c *Client) BatchGetViewCounts(ctx context.Context, request *BatchGetViewCountsRequest, opts ...RequestOpt) (*BatchGetViewCountsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.Resources) == 0 {
		return nil, errors.New("resources must not be empty")
	}
	for i := range request.Resources {
		ref := request.Resources[i]
		if err := validateEngagementResourceRef(&ref); err != nil {
			return nil, fmt.Errorf("resources[%d]: %w", i, err)
		}
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/engagement/views:batchGet", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BatchGetViewCountsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

func validateEngagementResourceRef(resource *EngagementResourceRef) error {
	if resource == nil {
		return errors.New("resource must be provided")
	}
	if strings.TrimSpace(string(resource.ResourceType)) == "" || resource.ResourceType == TrackableResourceTypeUnspecified {
		return errors.New("resource.resourceType must be specified")
	}
	if strings.TrimSpace(resource.ResourceID) == "" {
		return errors.New("resource.resourceID must not be empty")
	}
	return nil
}
