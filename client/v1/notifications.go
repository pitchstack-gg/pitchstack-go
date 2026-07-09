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

// NotificationAction represents v1Action.
type NotificationAction struct {
	Label string `json:"label,omitempty"`
	URL   string `json:"url,omitempty"`
	Style string `json:"style,omitempty"`
}

// Notification represents v1Notification.
type Notification struct {
	MessageID    string               `json:"messageId,omitempty"`
	Category     string               `json:"category,omitempty"`
	Severity     string               `json:"severity,omitempty"`
	Title        string               `json:"title,omitempty"`
	BodyMarkdown string               `json:"bodyMarkdown,omitempty"`
	Blocks       map[string]any       `json:"blocks,omitempty"`
	Actions      []NotificationAction `json:"actions,omitempty"`
	CreatedAt    *time.Time           `json:"createdAt,omitempty"`
	ExpiresAt    *time.Time           `json:"expiresAt,omitempty"`
	ReadAt       *time.Time           `json:"readAt,omitempty"`
	ArchivedAt   *time.Time           `json:"archivedAt,omitempty"`
}

// PushDevice represents a registered mobile push target.
type PushDevice struct {
	DeviceID         string     `json:"deviceId,omitempty"`
	Platform         string     `json:"platform,omitempty"`
	ExpoPushToken    string     `json:"expoPushToken,omitempty"`
	AppVersion       string     `json:"appVersion,omitempty"`
	Active           bool       `json:"active,omitempty"`
	LastRegisteredAt *time.Time `json:"lastRegisteredAt,omitempty"`
}

// NotificationPreference controls delivery channels for a category.
type NotificationPreference struct {
	Category     string `json:"category,omitempty"`
	InAppEnabled bool   `json:"inAppEnabled,omitempty"`
	PushEnabled  bool   `json:"pushEnabled,omitempty"`
	EmailEnabled bool   `json:"emailEnabled,omitempty"`
}

// NotificationTopicSubscription represents a category/topic subscription.
type NotificationTopicSubscription struct {
	Category  string     `json:"category,omitempty"`
	TopicType string     `json:"topicType,omitempty"`
	TopicID   string     `json:"topicId,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
}

// CreateMessageRequest submits a producer message for a target user inbox.
type CreateMessageRequest struct {
	TargetUserID   string               `json:"targetUserId,omitempty"`
	Source         string               `json:"source,omitempty"`
	IdempotencyKey string               `json:"idempotencyKey,omitempty"`
	Category       string               `json:"category,omitempty"`
	Severity       string               `json:"severity,omitempty"`
	Title          string               `json:"title,omitempty"`
	BodyMarkdown   *string              `json:"bodyMarkdown,omitempty"`
	Blocks         map[string]any       `json:"blocks,omitempty"`
	Actions        []NotificationAction `json:"actions,omitempty"`
	ExpiresAt      *time.Time           `json:"expiresAt,omitempty"`
}

// CreateMessageResponse returns the created (or deduplicated existing) message id.
type CreateMessageResponse struct {
	MessageID string           `json:"messageId,omitempty"`
	Created   bool             `json:"created,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *CreateMessageResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RegisterPushDeviceRequest registers or updates a mobile push device.
type RegisterPushDeviceRequest struct {
	DeviceID      string `json:"deviceId,omitempty"`
	Platform      string `json:"platform,omitempty"`
	ExpoPushToken string `json:"expoPushToken,omitempty"`
	AppVersion    string `json:"appVersion,omitempty"`
}

// RegisterPushDeviceResponse returns the registered push device.
type RegisterPushDeviceResponse struct {
	Device   *PushDevice      `json:"device,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *RegisterPushDeviceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnregisterPushDeviceRequest disables a mobile push device.
type UnregisterPushDeviceRequest struct {
	DeviceID string `json:"-"`
}

// UnregisterPushDeviceResponse captures metadata for unregister operations.
type UnregisterPushDeviceResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnregisterPushDeviceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetNotificationPreferencesResponse returns notification channel preferences.
type GetNotificationPreferencesResponse struct {
	Preferences []NotificationPreference `json:"preferences,omitempty"`
	Metadata    ResponseMetadata         `json:"-"`
}

func (r *GetNotificationPreferencesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateNotificationPreferencesRequest updates notification channel preferences.
type UpdateNotificationPreferencesRequest struct {
	Preferences []NotificationPreference `json:"preferences,omitempty"`
}

// UpdateNotificationPreferencesResponse returns updated notification channel preferences.
type UpdateNotificationPreferencesResponse struct {
	Preferences []NotificationPreference `json:"preferences,omitempty"`
	Metadata    ResponseMetadata         `json:"-"`
}

func (r *UpdateNotificationPreferencesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnsubscribeEmailRequest disables an email notification category from an unsubscribe link.
type UnsubscribeEmailRequest struct {
	Token string
}

// UnsubscribeEmailResponse captures metadata for email unsubscribe operations.
type UnsubscribeEmailResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnsubscribeEmailResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetNotificationTopicSubscriptionsRequest filters topic subscriptions.
type GetNotificationTopicSubscriptionsRequest struct {
	Category  string
	TopicType string
}

// GetNotificationTopicSubscriptionsResponse returns topic subscriptions.
type GetNotificationTopicSubscriptionsResponse struct {
	TopicIDs      []string                        `json:"topicIds,omitempty"`
	Subscriptions []NotificationTopicSubscription `json:"subscriptions,omitempty"`
	Metadata      ResponseMetadata                `json:"-"`
}

func (r *GetNotificationTopicSubscriptionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateNotificationTopicSubscriptionsRequest replaces topic subscriptions.
type UpdateNotificationTopicSubscriptionsRequest struct {
	Category  string   `json:"category,omitempty"`
	TopicType string   `json:"topicType,omitempty"`
	TopicIDs  []string `json:"topicIds,omitempty"`
}

// UpdateNotificationTopicSubscriptionsResponse returns updated topic subscriptions.
type UpdateNotificationTopicSubscriptionsResponse struct {
	TopicIDs      []string                        `json:"topicIds,omitempty"`
	Subscriptions []NotificationTopicSubscription `json:"subscriptions,omitempty"`
	Metadata      ResponseMetadata                `json:"-"`
}

func (r *UpdateNotificationTopicSubscriptionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListInboxRequest controls ListInbox filtering and pagination.
type ListInboxRequest struct {
	UnreadOnly      *bool
	IncludeArchived *bool
	IncludeExpired  *bool
	Categories      []string
	PageSize        *int32
	NextToken       string
}

// ListInboxResponse lists notifications for the authenticated user.
type ListInboxResponse struct {
	Notifications []Notification   `json:"notifications,omitempty"`
	NextToken     string           `json:"nextToken,omitempty"`
	Metadata      ResponseMetadata `json:"-"`
}

func (r *ListInboxResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetMessageRequest identifies a notification by message ID.
type GetMessageRequest struct {
	MessageID string
}

// GetMessageResponse returns a specific notification.
type GetMessageResponse struct {
	Notification *Notification    `json:"notification,omitempty"`
	Metadata     ResponseMetadata `json:"-"`
}

func (r *GetMessageResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteMessageRequest identifies the notification to delete.
type DeleteMessageRequest struct {
	MessageID string
}

// DeleteMessageResponse captures metadata for delete operations.
type DeleteMessageResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteMessageResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// MarkReadRequest identifies the notification to mark as read.
type MarkReadRequest struct {
	MessageID string
}

// MarkReadResponse captures metadata for mark-read operations.
type MarkReadResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *MarkReadResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// MarkAllReadRequest marks all matching notifications as read.
type MarkAllReadRequest struct{}

// MarkAllReadResponse returns the number of notifications marked read.
type MarkAllReadResponse struct {
	MarkedReadCount int64            `json:"markedReadCount,omitempty,string"`
	Metadata        ResponseMetadata `json:"-"`
}

func (r *MarkAllReadResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ArchiveMessageRequest identifies the notification to archive.
type ArchiveMessageRequest struct {
	MessageID string
}

// ArchiveMessageResponse captures metadata for archive operations.
type ArchiveMessageResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *ArchiveMessageResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateMessage creates a notification message via the producer service.
func (c *Client) CreateMessage(ctx context.Context, request *CreateMessageRequest, opts ...RequestOpt) (*CreateMessageResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.TargetUserID) == "" {
		return nil, errors.New("targetUserID must not be empty")
	}
	if strings.TrimSpace(request.Source) == "" {
		return nil, errors.New("source must not be empty")
	}
	if strings.TrimSpace(request.IdempotencyKey) == "" {
		return nil, errors.New("idempotencyKey must not be empty")
	}
	if strings.TrimSpace(request.Category) == "" {
		return nil, errors.New("category must not be empty")
	}
	if strings.TrimSpace(request.Title) == "" {
		return nil, errors.New("title must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/notifications.v1.NotificationsProducerService/CreateMessage", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CreateMessageResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// RegisterPushDevice registers or updates a mobile push device.
func (c *Client) RegisterPushDevice(ctx context.Context, request *RegisterPushDeviceRequest, opts ...RequestOpt) (*RegisterPushDeviceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.DeviceID) == "" {
		return nil, errors.New("deviceID must not be empty")
	}
	if strings.TrimSpace(request.Platform) == "" {
		return nil, errors.New("platform must not be empty")
	}
	if strings.TrimSpace(request.ExpoPushToken) == "" {
		return nil, errors.New("expoPushToken must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/notifications/devices", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &RegisterPushDeviceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UnregisterPushDevice disables a mobile push device.
func (c *Client) UnregisterPushDevice(ctx context.Context, request *UnregisterPushDeviceRequest, opts ...RequestOpt) (*UnregisterPushDeviceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	deviceID := strings.TrimSpace(request.DeviceID)
	if deviceID == "" {
		return nil, errors.New("deviceID must not be empty")
	}

	path := fmt.Sprintf("/v1/notifications/devices/%s", url.PathEscape(deviceID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &UnregisterPushDeviceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetNotificationPreferences returns notification channel preferences.
func (c *Client) GetNotificationPreferences(ctx context.Context, opts ...RequestOpt) (*GetNotificationPreferencesResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/notifications/preferences", nil)
	if err != nil {
		return nil, err
	}

	response := &GetNotificationPreferencesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpdateNotificationPreferences updates notification channel preferences.
func (c *Client) UpdateNotificationPreferences(ctx context.Context, request *UpdateNotificationPreferencesRequest, opts ...RequestOpt) (*UpdateNotificationPreferencesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPut, "/v1/notifications/preferences", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdateNotificationPreferencesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UnsubscribeEmail disables one email notification category using an unsubscribe token.
func (c *Client) UnsubscribeEmail(ctx context.Context, request *UnsubscribeEmailRequest, opts ...RequestOpt) (*UnsubscribeEmailResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	token := strings.TrimSpace(request.Token)
	if token == "" {
		return nil, errors.New("token must not be empty")
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/notifications/email:unsubscribe", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("token", token)
	req.URL.RawQuery = query.Encode()

	response := &UnsubscribeEmailResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetNotificationTopicSubscriptions returns topic subscriptions.
func (c *Client) GetNotificationTopicSubscriptions(ctx context.Context, request *GetNotificationTopicSubscriptionsRequest, opts ...RequestOpt) (*GetNotificationTopicSubscriptionsResponse, error) {
	if request == nil {
		request = &GetNotificationTopicSubscriptionsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/notifications/topic-subscriptions", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "category", request.Category)
	setQueryString(query, "topicType", request.TopicType)
	req.URL.RawQuery = query.Encode()

	response := &GetNotificationTopicSubscriptionsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpdateNotificationTopicSubscriptions replaces topic subscriptions.
func (c *Client) UpdateNotificationTopicSubscriptions(ctx context.Context, request *UpdateNotificationTopicSubscriptionsRequest, opts ...RequestOpt) (*UpdateNotificationTopicSubscriptionsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.Category) == "" {
		return nil, errors.New("category must not be empty")
	}
	if strings.TrimSpace(request.TopicType) == "" {
		return nil, errors.New("topicType must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPut, "/v1/notifications/topic-subscriptions", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdateNotificationTopicSubscriptionsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListInbox lists notifications for the authenticated user.
func (c *Client) ListInbox(ctx context.Context, request *ListInboxRequest, opts ...RequestOpt) (*ListInboxResponse, error) {
	if request == nil {
		request = &ListInboxRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/notifications", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.UnreadOnly != nil {
		query.Set("unreadOnly", strconv.FormatBool(*request.UnreadOnly))
	}
	if request.IncludeArchived != nil {
		query.Set("includeArchived", strconv.FormatBool(*request.IncludeArchived))
	}
	if request.IncludeExpired != nil {
		query.Set("includeExpired", strconv.FormatBool(*request.IncludeExpired))
	}
	for _, category := range request.Categories {
		if c := strings.TrimSpace(category); c != "" {
			query.Add("categories", c)
		}
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	response := &ListInboxResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetMessage fetches a notification for the authenticated user.
func (c *Client) GetMessage(ctx context.Context, request *GetMessageRequest, opts ...RequestOpt) (*GetMessageResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	messageID := strings.TrimSpace(request.MessageID)
	if messageID == "" {
		return nil, errors.New("messageID must not be empty")
	}

	path := fmt.Sprintf("/v1/notifications/%s", url.PathEscape(messageID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetMessageResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// DeleteMessage removes a notification for the authenticated user.
func (c *Client) DeleteMessage(ctx context.Context, request *DeleteMessageRequest, opts ...RequestOpt) (*DeleteMessageResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	messageID := strings.TrimSpace(request.MessageID)
	if messageID == "" {
		return nil, errors.New("messageID must not be empty")
	}

	path := fmt.Sprintf("/v1/notifications/%s", url.PathEscape(messageID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteMessageResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// MarkRead marks a notification as read for the authenticated user.
func (c *Client) MarkRead(ctx context.Context, request *MarkReadRequest, opts ...RequestOpt) (*MarkReadResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	messageID := strings.TrimSpace(request.MessageID)
	if messageID == "" {
		return nil, errors.New("messageID must not be empty")
	}

	body, err := jsonBody(struct{}{})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/notifications/%s:markRead", url.PathEscape(messageID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &MarkReadResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// MarkAllRead marks all notifications as read for the authenticated user.
func (c *Client) MarkAllRead(ctx context.Context, request *MarkAllReadRequest, opts ...RequestOpt) (*MarkAllReadResponse, error) {
	if request == nil {
		request = &MarkAllReadRequest{}
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/notifications:markAllRead", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &MarkAllReadResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ArchiveMessage archives a notification for the authenticated user.
func (c *Client) ArchiveMessage(ctx context.Context, request *ArchiveMessageRequest, opts ...RequestOpt) (*ArchiveMessageResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	messageID := strings.TrimSpace(request.MessageID)
	if messageID == "" {
		return nil, errors.New("messageID must not be empty")
	}

	body, err := jsonBody(struct{}{})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/notifications/%s:archive", url.PathEscape(messageID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &ArchiveMessageResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
