package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
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

// LikeableResourceType represents engagement.v1.LikeableResourceType.
type LikeableResourceType string

const (
	LikeableResourceTypeUnspecified LikeableResourceType = "LIKEABLE_RESOURCE_TYPE_UNSPECIFIED"
	LikeableResourceTypeDeck        LikeableResourceType = "LIKEABLE_RESOURCE_TYPE_DECK"
	LikeableResourceTypeCollection  LikeableResourceType = "LIKEABLE_RESOURCE_TYPE_COLLECTION"
)

// EngagementResourceRef mirrors engagement.v1.ResourceRef.
type EngagementResourceRef struct {
	ResourceType TrackableResourceType `json:"resourceType,omitempty"`
	ResourceID   string                `json:"resourceId,omitempty"`
}

// LikeableResourceRef mirrors engagement.v1.LikeableResourceRef.
type LikeableResourceRef struct {
	ResourceType LikeableResourceType `json:"resourceType,omitempty"`
	ResourceID   string               `json:"resourceId,omitempty"`
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

func (r *TrendingResource) UnmarshalJSON(data []byte) error {
	type rawTrendingResource struct {
		Resource     *EngagementResourceRef `json:"resource,omitempty"`
		ViewCount    json.RawMessage        `json:"viewCount,omitempty"`
		Score        float64                `json:"score,omitempty"`
		LastViewedAt *time.Time             `json:"lastViewedAt,omitempty"`
	}

	var raw rawTrendingResource
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	count, err := parseFlexibleInt64(raw.ViewCount)
	if err != nil {
		return fmt.Errorf("viewCount: %w", err)
	}

	r.Resource = raw.Resource
	r.ViewCount = count
	r.Score = raw.Score
	r.LastViewedAt = raw.LastViewedAt
	return nil
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

func (r *ResourceViewCount) UnmarshalJSON(data []byte) error {
	type rawResourceViewCount struct {
		Resource     *EngagementResourceRef `json:"resource,omitempty"`
		TotalViews   json.RawMessage        `json:"totalViews,omitempty"`
		LastViewedAt *time.Time             `json:"lastViewedAt,omitempty"`
	}

	var raw rawResourceViewCount
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	totalViews, err := parseFlexibleInt64(raw.TotalViews)
	if err != nil {
		return fmt.Errorf("totalViews: %w", err)
	}

	r.Resource = raw.Resource
	r.TotalViews = totalViews
	r.LastViewedAt = raw.LastViewedAt
	return nil
}

// BatchGetViewCountsResponse mirrors engagement.v1.BatchGetViewCountsResponse.
type BatchGetViewCountsResponse struct {
	Counts   []ResourceViewCount `json:"counts,omitempty"`
	Metadata ResponseMetadata    `json:"-"`
}

func (r *BatchGetViewCountsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// LikeResourceRequest likes a resource.
type LikeResourceRequest struct {
	Resource *LikeableResourceRef `json:"resource,omitempty"`
}

// LikeResourceResponse returns like state after a like operation.
type LikeResourceResponse struct {
	Resource   *LikeableResourceRef `json:"resource,omitempty"`
	Liked      bool                 `json:"liked,omitempty"`
	TotalLikes int64                `json:"totalLikes,omitempty"`
	Changed    bool                 `json:"changed,omitempty"`
	Metadata   ResponseMetadata     `json:"-"`
}

func (r *LikeResourceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

func (r *LikeResourceResponse) UnmarshalJSON(data []byte) error {
	type rawLikeResourceResponse struct {
		Resource   *LikeableResourceRef `json:"resource,omitempty"`
		Liked      bool                 `json:"liked,omitempty"`
		TotalLikes json.RawMessage      `json:"totalLikes,omitempty"`
		Changed    bool                 `json:"changed,omitempty"`
	}

	var raw rawLikeResourceResponse
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	totalLikes, err := parseFlexibleInt64(raw.TotalLikes)
	if err != nil {
		return fmt.Errorf("totalLikes: %w", err)
	}
	r.Resource = raw.Resource
	r.Liked = raw.Liked
	r.TotalLikes = totalLikes
	r.Changed = raw.Changed
	return nil
}

// UnlikeResourceRequest unlikes a resource.
type UnlikeResourceRequest struct {
	Resource *LikeableResourceRef `json:"resource,omitempty"`
}

// UnlikeResourceResponse returns like state after an unlike operation.
type UnlikeResourceResponse struct {
	Resource   *LikeableResourceRef `json:"resource,omitempty"`
	Liked      bool                 `json:"liked,omitempty"`
	TotalLikes int64                `json:"totalLikes,omitempty"`
	Changed    bool                 `json:"changed,omitempty"`
	Metadata   ResponseMetadata     `json:"-"`
}

func (r *UnlikeResourceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

func (r *UnlikeResourceResponse) UnmarshalJSON(data []byte) error {
	var like LikeResourceResponse
	if err := json.Unmarshal(data, &like); err != nil {
		return err
	}
	r.Resource = like.Resource
	r.Liked = like.Liked
	r.TotalLikes = like.TotalLikes
	r.Changed = like.Changed
	return nil
}

// BatchGetLikeCountsRequest fetches like counts for resources.
type BatchGetLikeCountsRequest struct {
	Resources []LikeableResourceRef `json:"resources,omitempty"`
}

// ResourceLikeCount mirrors engagement.v1.ResourceLikeCount.
type ResourceLikeCount struct {
	Resource   *LikeableResourceRef `json:"resource,omitempty"`
	TotalLikes int64                `json:"totalLikes,omitempty"`
}

func (r *ResourceLikeCount) UnmarshalJSON(data []byte) error {
	type rawResourceLikeCount struct {
		Resource   *LikeableResourceRef `json:"resource,omitempty"`
		TotalLikes json.RawMessage      `json:"totalLikes,omitempty"`
	}

	var raw rawResourceLikeCount
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	totalLikes, err := parseFlexibleInt64(raw.TotalLikes)
	if err != nil {
		return fmt.Errorf("totalLikes: %w", err)
	}
	r.Resource = raw.Resource
	r.TotalLikes = totalLikes
	return nil
}

// BatchGetLikeCountsResponse returns like counts.
type BatchGetLikeCountsResponse struct {
	Counts   []ResourceLikeCount `json:"counts,omitempty"`
	Metadata ResponseMetadata    `json:"-"`
}

func (r *BatchGetLikeCountsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetViewerLikesRequest fetches viewer like states for resources.
type BatchGetViewerLikesRequest struct {
	Resources []LikeableResourceRef `json:"resources,omitempty"`
}

// ViewerLike mirrors engagement.v1.ViewerLike.
type ViewerLike struct {
	Resource *LikeableResourceRef `json:"resource,omitempty"`
	Liked    bool                 `json:"liked,omitempty"`
	LikedAt  *time.Time           `json:"likedAt,omitempty"`
}

// BatchGetViewerLikesResponse returns viewer like states.
type BatchGetViewerLikesResponse struct {
	Likes    []ViewerLike     `json:"likes,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *BatchGetViewerLikesResponse) setMetadata(metadata ResponseMetadata) {
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

// LikeResource idempotently likes a resource.
func (c *Client) LikeResource(ctx context.Context, request *LikeResourceRequest, opts ...RequestOpt) (*LikeResourceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := validateLikeableResourceRef(request.Resource); err != nil {
		return nil, err
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/engagement/likes:like", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &LikeResourceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UnlikeResource idempotently removes a like from a resource.
func (c *Client) UnlikeResource(ctx context.Context, request *UnlikeResourceRequest, opts ...RequestOpt) (*UnlikeResourceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := validateLikeableResourceRef(request.Resource); err != nil {
		return nil, err
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/engagement/likes:unlike", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UnlikeResourceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BatchGetLikeCounts retrieves aggregate like counts for multiple resources.
func (c *Client) BatchGetLikeCounts(ctx context.Context, request *BatchGetLikeCountsRequest, opts ...RequestOpt) (*BatchGetLikeCountsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.Resources) == 0 {
		return nil, errors.New("resources must not be empty")
	}
	for i := range request.Resources {
		ref := request.Resources[i]
		if err := validateLikeableResourceRef(&ref); err != nil {
			return nil, fmt.Errorf("resources[%d]: %w", i, err)
		}
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/engagement/likes:batchGetCounts", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BatchGetLikeCountsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BatchGetViewerLikes retrieves whether the authenticated viewer liked resources.
func (c *Client) BatchGetViewerLikes(ctx context.Context, request *BatchGetViewerLikesRequest, opts ...RequestOpt) (*BatchGetViewerLikesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.Resources) == 0 {
		return nil, errors.New("resources must not be empty")
	}
	for i := range request.Resources {
		ref := request.Resources[i]
		if err := validateLikeableResourceRef(&ref); err != nil {
			return nil, fmt.Errorf("resources[%d]: %w", i, err)
		}
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/engagement/likes:batchGetViewerLikes", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BatchGetViewerLikesResponse{}
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

func validateLikeableResourceRef(resource *LikeableResourceRef) error {
	if resource == nil {
		return errors.New("resource must be provided")
	}
	if strings.TrimSpace(string(resource.ResourceType)) == "" || resource.ResourceType == LikeableResourceTypeUnspecified {
		return errors.New("resource.resourceType must be specified")
	}
	if strings.TrimSpace(resource.ResourceID) == "" {
		return errors.New("resource.resourceID must not be empty")
	}
	return nil
}

func parseFlexibleInt64(raw json.RawMessage) (int64, error) {
	if len(raw) == 0 || string(raw) == "null" {
		return 0, nil
	}

	var asInt int64
	if err := json.Unmarshal(raw, &asInt); err == nil {
		return asInt, nil
	}

	var asString string
	if err := json.Unmarshal(raw, &asString); err == nil {
		asString = strings.TrimSpace(asString)
		if asString == "" {
			return 0, nil
		}
		value, parseErr := strconv.ParseInt(asString, 10, 64)
		if parseErr != nil {
			return 0, parseErr
		}
		return value, nil
	}

	return 0, fmt.Errorf("expected int64 or string-encoded int64, got %s", string(raw))
}
