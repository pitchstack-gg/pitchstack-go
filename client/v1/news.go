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

// NewsSourceType represents news.v1.NewsSourceType.
type NewsSourceType string

const (
	NewsSourceTypeUnspecified NewsSourceType = "NEWS_SOURCE_TYPE_UNSPECIFIED"
	NewsSourceTypeSystem      NewsSourceType = "NEWS_SOURCE_TYPE_SYSTEM"
	NewsSourceTypeOfficial    NewsSourceType = "NEWS_SOURCE_TYPE_OFFICIAL"
	NewsSourceTypePartner     NewsSourceType = "NEWS_SOURCE_TYPE_PARTNER"
	NewsSourceTypeCommunity   NewsSourceType = "NEWS_SOURCE_TYPE_COMMUNITY"
)

// NewsSourceStatus represents news.v1.NewsSourceStatus.
type NewsSourceStatus string

const (
	NewsSourceStatusUnspecified NewsSourceStatus = "NEWS_SOURCE_STATUS_UNSPECIFIED"
	NewsSourceStatusActive      NewsSourceStatus = "NEWS_SOURCE_STATUS_ACTIVE"
	NewsSourceStatusDisabled    NewsSourceStatus = "NEWS_SOURCE_STATUS_DISABLED"
)

// NewsArticleOrigin represents news.v1.NewsArticleOrigin.
type NewsArticleOrigin string

const (
	NewsArticleOriginUnspecified NewsArticleOrigin = "NEWS_ARTICLE_ORIGIN_UNSPECIFIED"
	NewsArticleOriginRSS         NewsArticleOrigin = "NEWS_ARTICLE_ORIGIN_RSS"
	NewsArticleOriginManual      NewsArticleOrigin = "NEWS_ARTICLE_ORIGIN_MANUAL"
)

// NewsArticleStatus represents news.v1.NewsArticleStatus.
type NewsArticleStatus string

const (
	NewsArticleStatusUnspecified NewsArticleStatus = "NEWS_ARTICLE_STATUS_UNSPECIFIED"
	NewsArticleStatusDraft       NewsArticleStatus = "NEWS_ARTICLE_STATUS_DRAFT"
	NewsArticleStatusScheduled   NewsArticleStatus = "NEWS_ARTICLE_STATUS_SCHEDULED"
	NewsArticleStatusPublished   NewsArticleStatus = "NEWS_ARTICLE_STATUS_PUBLISHED"
	NewsArticleStatusUnpublished NewsArticleStatus = "NEWS_ARTICLE_STATUS_UNPUBLISHED"
	NewsArticleStatusArchived    NewsArticleStatus = "NEWS_ARTICLE_STATUS_ARCHIVED"
)

// NewsSource mirrors news.v1.NewsSource.
type NewsSource struct {
	SourceID            string           `json:"sourceId,omitempty"`
	Name                string           `json:"name,omitempty"`
	RSSURL              string           `json:"rssUrl,omitempty"`
	SourceType          NewsSourceType   `json:"sourceType,omitempty"`
	Status              NewsSourceStatus `json:"status,omitempty"`
	BaseWeight          float64          `json:"baseWeight,omitempty"`
	PollIntervalSeconds int32            `json:"pollIntervalSeconds,omitempty"`
	NextPollAt          *time.Time       `json:"nextPollAt,omitempty"`
	LastPolledAt        *time.Time       `json:"lastPolledAt,omitempty"`
	LastSuccessAt       *time.Time       `json:"lastSuccessAt,omitempty"`
	FailureCount        int32            `json:"failureCount,omitempty"`
	ETag                string           `json:"etag,omitempty"`
	LastModified        string           `json:"lastModified,omitempty"`
	CreatedAt           *time.Time       `json:"createdAt,omitempty"`
	UpdatedAt           *time.Time       `json:"updatedAt,omitempty"`
}

// NewsArticleOverride mirrors news.v1.NewsArticleOverride.
type NewsArticleOverride struct {
	PinRank     *int32     `json:"pinRank,omitempty"`
	ManualBoost *float64   `json:"manualBoost,omitempty"`
	StartsAt    *time.Time `json:"startsAt,omitempty"`
	EndsAt      *time.Time `json:"endsAt,omitempty"`
}

// NewsArticle mirrors news.v1.NewsArticle.
type NewsArticle struct {
	ArticleID    string               `json:"articleId,omitempty"`
	Origin       NewsArticleOrigin    `json:"origin,omitempty"`
	Status       NewsArticleStatus    `json:"status,omitempty"`
	SourceID     string               `json:"sourceId,omitempty"`
	SourceName   string               `json:"sourceName,omitempty"`
	SourceType   NewsSourceType       `json:"sourceType,omitempty"`
	Title        string               `json:"title,omitempty"`
	Summary      string               `json:"summary,omitempty"`
	CanonicalURL string               `json:"canonicalUrl,omitempty"`
	ImageURL     string               `json:"imageUrl,omitempty"`
	Author       string               `json:"author,omitempty"`
	Language     string               `json:"language,omitempty"`
	PublishedAt  *time.Time           `json:"publishedAt,omitempty"`
	ExpireAt     *time.Time           `json:"expireAt,omitempty"`
	BodyMarkdown string               `json:"bodyMarkdown,omitempty"`
	BodyHTML     string               `json:"bodyHtml,omitempty"`
	Metadata     map[string]any       `json:"metadata,omitempty"`
	Override     *NewsArticleOverride `json:"override,omitempty"`
	CreatedAt    *time.Time           `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time           `json:"updatedAt,omitempty"`
}

// RecommendedArticle mirrors news.v1.RecommendedArticle.
type RecommendedArticle struct {
	Article *NewsArticle `json:"article,omitempty"`
	Score   float64      `json:"score,omitempty"`
	PinRank *int32       `json:"pinRank,omitempty"`
}

// ListRecommendedArticlesRequest configures a recommended-news query.
type ListRecommendedArticlesRequest struct {
	PageSize  *int32
	NextToken string
	Locale    string
	Platform  string
}

// ListRecommendedArticlesResponse returns ranked feed entries.
type ListRecommendedArticlesResponse struct {
	Items     []RecommendedArticle `json:"items,omitempty"`
	NextToken string               `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata     `json:"-"`
}

func (r *ListRecommendedArticlesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetArticleRequest identifies an article to fetch.
type GetArticleRequest struct {
	ArticleID string `json:"-"`
}

// GetArticleResponse returns a single article.
type GetArticleResponse struct {
	Article  *NewsArticle     `json:"article,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetArticleResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// TrackArticleImpressionRequest records impression telemetry.
type TrackArticleImpressionRequest struct {
	ArticleID      string         `json:"articleId,omitempty"`
	SessionID      string         `json:"sessionId,omitempty"`
	ClientEventID  string         `json:"clientEventId,omitempty"`
	OccurredAt     *time.Time     `json:"occurredAt,omitempty"`
	MetadataFields map[string]any `json:"metadata,omitempty"`
}

// TrackArticleImpressionResponse indicates acceptance and dedupe.
type TrackArticleImpressionResponse struct {
	Accepted bool             `json:"accepted,omitempty"`
	Deduped  bool             `json:"deduped,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *TrackArticleImpressionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// TrackArticleClickRequest records click telemetry.
type TrackArticleClickRequest struct {
	ArticleID      string         `json:"articleId,omitempty"`
	SessionID      string         `json:"sessionId,omitempty"`
	ClientEventID  string         `json:"clientEventId,omitempty"`
	OccurredAt     *time.Time     `json:"occurredAt,omitempty"`
	DestinationURL string         `json:"destinationUrl,omitempty"`
	MetadataFields map[string]any `json:"metadata,omitempty"`
}

// TrackArticleClickResponse indicates acceptance and dedupe.
type TrackArticleClickResponse struct {
	Accepted bool             `json:"accepted,omitempty"`
	Deduped  bool             `json:"deduped,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *TrackArticleClickResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DisableSourceRequest disables a news source.
type DisableSourceRequest struct {
	SourceID string `json:"sourceId,omitempty"`
}

// DisableSourceResponse returns the updated source.
type DisableSourceResponse struct {
	Source   *NewsSource      `json:"source,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *DisableSourceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// EnableSourceRequest enables a news source.
type EnableSourceRequest struct {
	SourceID string `json:"sourceId,omitempty"`
}

// EnableSourceResponse returns the updated source.
type EnableSourceResponse struct {
	Source   *NewsSource      `json:"source,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *EnableSourceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RunSourceIngestionRequest requests manual ingestion for a source.
type RunSourceIngestionRequest struct {
	SourceID         string `json:"sourceId,omitempty"`
	ForceFullRefresh *bool  `json:"forceFullRefresh,omitempty"`
}

// RunSourceIngestionResponse describes queueing state.
type RunSourceIngestionResponse struct {
	SourceID string           `json:"sourceId,omitempty"`
	RunID    string           `json:"runId,omitempty"`
	Queued   bool             `json:"queued,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *RunSourceIngestionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateManualArticleRequest defines first-party article fields.
type CreateManualArticleRequest struct {
	Title         string             `json:"title,omitempty"`
	Summary       string             `json:"summary,omitempty"`
	CanonicalURL  string             `json:"canonicalUrl,omitempty"`
	ImageURL      string             `json:"imageUrl,omitempty"`
	Author        string             `json:"author,omitempty"`
	Language      string             `json:"language,omitempty"`
	BodyMarkdown  string             `json:"bodyMarkdown,omitempty"`
	BodyHTML      string             `json:"bodyHtml,omitempty"`
	Metadata      map[string]any     `json:"metadata,omitempty"`
	PublishAt     *time.Time         `json:"publishAt,omitempty"`
	ExpireAt      *time.Time         `json:"expireAt,omitempty"`
	InitialStatus *NewsArticleStatus `json:"initialStatus,omitempty"`
}

// CreateManualArticleResponse returns the created article.
type CreateManualArticleResponse struct {
	Article  *NewsArticle     `json:"article,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CreateManualArticleResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateManualArticleRequest defines mutable fields for a manual article.
type UpdateManualArticleRequest struct {
	ArticleID    string         `json:"articleId,omitempty"`
	Title        *string        `json:"title,omitempty"`
	Summary      *string        `json:"summary,omitempty"`
	CanonicalURL *string        `json:"canonicalUrl,omitempty"`
	ImageURL     *string        `json:"imageUrl,omitempty"`
	Author       *string        `json:"author,omitempty"`
	Language     *string        `json:"language,omitempty"`
	BodyMarkdown *string        `json:"bodyMarkdown,omitempty"`
	BodyHTML     *string        `json:"bodyHtml,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
	PublishAt    *time.Time     `json:"publishAt,omitempty"`
	ExpireAt     *time.Time     `json:"expireAt,omitempty"`
}

// UpdateManualArticleResponse returns the updated article.
type UpdateManualArticleResponse struct {
	Article  *NewsArticle     `json:"article,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateManualArticleResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// PublishManualArticleRequest publishes a manual article.
type PublishManualArticleRequest struct {
	ArticleID string     `json:"articleId,omitempty"`
	PublishAt *time.Time `json:"publishAt,omitempty"`
	ExpireAt  *time.Time `json:"expireAt,omitempty"`
}

// PublishManualArticleResponse returns the updated article.
type PublishManualArticleResponse struct {
	Article  *NewsArticle     `json:"article,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *PublishManualArticleResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnpublishManualArticleRequest unpublishes an article.
type UnpublishManualArticleRequest struct {
	ArticleID string `json:"articleId,omitempty"`
}

// UnpublishManualArticleResponse returns the updated article.
type UnpublishManualArticleResponse struct {
	Article  *NewsArticle     `json:"article,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnpublishManualArticleResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// SetArticleOverrideRequest applies editorial overrides.
type SetArticleOverrideRequest struct {
	ArticleID   string     `json:"articleId,omitempty"`
	PinRank     *int32     `json:"pinRank,omitempty"`
	ManualBoost *float64   `json:"manualBoost,omitempty"`
	StartsAt    *time.Time `json:"startsAt,omitempty"`
	EndsAt      *time.Time `json:"endsAt,omitempty"`
}

// SetArticleOverrideResponse returns the updated article.
type SetArticleOverrideResponse struct {
	Article  *NewsArticle     `json:"article,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *SetArticleOverrideResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListRecommendedArticles retrieves ranked recommended articles.
func (c *Client) ListRecommendedArticles(ctx context.Context, request *ListRecommendedArticlesRequest, opts ...RequestOpt) (*ListRecommendedArticlesResponse, error) {
	if request == nil {
		request = &ListRecommendedArticlesRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/news/recommendations", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	setQueryString(query, "nextToken", request.NextToken)
	setQueryString(query, "locale", request.Locale)
	setQueryString(query, "platform", request.Platform)
	req.URL.RawQuery = query.Encode()

	response := &ListRecommendedArticlesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetArticle fetches a single article.
func (c *Client) GetArticle(ctx context.Context, request *GetArticleRequest, opts ...RequestOpt) (*GetArticleResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	articleID := strings.TrimSpace(request.ArticleID)
	if articleID == "" {
		return nil, errors.New("articleID must not be empty")
	}

	path := fmt.Sprintf("/v1/news/articles/%s", url.PathEscape(articleID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetArticleResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// TrackArticleImpression records an article impression event.
func (c *Client) TrackArticleImpression(ctx context.Context, request *TrackArticleImpressionRequest, opts ...RequestOpt) (*TrackArticleImpressionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	articleID := strings.TrimSpace(request.ArticleID)
	if articleID == "" {
		return nil, errors.New("articleID must not be empty")
	}

	body, err := jsonBody(struct {
		ArticleID      string         `json:"articleId,omitempty"`
		SessionID      string         `json:"sessionId,omitempty"`
		ClientEventID  string         `json:"clientEventId,omitempty"`
		OccurredAt     *time.Time     `json:"occurredAt,omitempty"`
		MetadataFields map[string]any `json:"metadata,omitempty"`
	}{
		ArticleID:      articleID,
		SessionID:      strings.TrimSpace(request.SessionID),
		ClientEventID:  strings.TrimSpace(request.ClientEventID),
		OccurredAt:     request.OccurredAt,
		MetadataFields: request.MetadataFields,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/news/articles/%s:trackImpression", url.PathEscape(articleID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &TrackArticleImpressionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// TrackArticleClick records an article click event.
func (c *Client) TrackArticleClick(ctx context.Context, request *TrackArticleClickRequest, opts ...RequestOpt) (*TrackArticleClickResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	articleID := strings.TrimSpace(request.ArticleID)
	if articleID == "" {
		return nil, errors.New("articleID must not be empty")
	}

	body, err := jsonBody(struct {
		ArticleID      string         `json:"articleId,omitempty"`
		SessionID      string         `json:"sessionId,omitempty"`
		ClientEventID  string         `json:"clientEventId,omitempty"`
		OccurredAt     *time.Time     `json:"occurredAt,omitempty"`
		DestinationURL string         `json:"destinationUrl,omitempty"`
		MetadataFields map[string]any `json:"metadata,omitempty"`
	}{
		ArticleID:      articleID,
		SessionID:      strings.TrimSpace(request.SessionID),
		ClientEventID:  strings.TrimSpace(request.ClientEventID),
		OccurredAt:     request.OccurredAt,
		DestinationURL: strings.TrimSpace(request.DestinationURL),
		MetadataFields: request.MetadataFields,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/news/articles/%s:trackClick", url.PathEscape(articleID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &TrackArticleClickResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DisableSource disables a source.
func (c *Client) DisableSource(ctx context.Context, request *DisableSourceRequest, opts ...RequestOpt) (*DisableSourceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	sourceID := strings.TrimSpace(request.SourceID)
	if sourceID == "" {
		return nil, errors.New("sourceID must not be empty")
	}

	body, err := jsonBody(struct {
		SourceID string `json:"sourceId,omitempty"`
	}{SourceID: sourceID})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/admin/news/sources/%s:disable", url.PathEscape(sourceID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &DisableSourceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// EnableSource enables a source.
func (c *Client) EnableSource(ctx context.Context, request *EnableSourceRequest, opts ...RequestOpt) (*EnableSourceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	sourceID := strings.TrimSpace(request.SourceID)
	if sourceID == "" {
		return nil, errors.New("sourceID must not be empty")
	}

	body, err := jsonBody(struct {
		SourceID string `json:"sourceId,omitempty"`
	}{SourceID: sourceID})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/admin/news/sources/%s:enable", url.PathEscape(sourceID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &EnableSourceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// RunSourceIngestion triggers an immediate source ingestion run.
func (c *Client) RunSourceIngestion(ctx context.Context, request *RunSourceIngestionRequest, opts ...RequestOpt) (*RunSourceIngestionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	sourceID := strings.TrimSpace(request.SourceID)
	if sourceID == "" {
		return nil, errors.New("sourceID must not be empty")
	}

	body, err := jsonBody(struct {
		SourceID         string `json:"sourceId,omitempty"`
		ForceFullRefresh *bool  `json:"forceFullRefresh,omitempty"`
	}{
		SourceID:         sourceID,
		ForceFullRefresh: request.ForceFullRefresh,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/admin/news/sources/%s:ingestNow", url.PathEscape(sourceID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &RunSourceIngestionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateManualArticle creates a manual article.
func (c *Client) CreateManualArticle(ctx context.Context, request *CreateManualArticleRequest, opts ...RequestOpt) (*CreateManualArticleResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.Title) == "" {
		return nil, errors.New("title must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/admin/news/articles", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CreateManualArticleResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateManualArticle updates a manual article.
func (c *Client) UpdateManualArticle(ctx context.Context, request *UpdateManualArticleRequest, opts ...RequestOpt) (*UpdateManualArticleResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	articleID := strings.TrimSpace(request.ArticleID)
	if articleID == "" {
		return nil, errors.New("articleID must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/admin/news/articles/%s", url.PathEscape(articleID))
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdateManualArticleResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// PublishManualArticle publishes a manual article.
func (c *Client) PublishManualArticle(ctx context.Context, request *PublishManualArticleRequest, opts ...RequestOpt) (*PublishManualArticleResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	articleID := strings.TrimSpace(request.ArticleID)
	if articleID == "" {
		return nil, errors.New("articleID must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/admin/news/articles/%s:publish", url.PathEscape(articleID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &PublishManualArticleResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UnpublishManualArticle unpublishes a manual article.
func (c *Client) UnpublishManualArticle(ctx context.Context, request *UnpublishManualArticleRequest, opts ...RequestOpt) (*UnpublishManualArticleResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	articleID := strings.TrimSpace(request.ArticleID)
	if articleID == "" {
		return nil, errors.New("articleID must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/admin/news/articles/%s:unpublish", url.PathEscape(articleID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UnpublishManualArticleResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// SetArticleOverride applies an editorial override.
func (c *Client) SetArticleOverride(ctx context.Context, request *SetArticleOverrideRequest, opts ...RequestOpt) (*SetArticleOverrideResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	articleID := strings.TrimSpace(request.ArticleID)
	if articleID == "" {
		return nil, errors.New("articleID must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/admin/news/articles/%s:setOverride", url.PathEscape(articleID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &SetArticleOverrideResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}
