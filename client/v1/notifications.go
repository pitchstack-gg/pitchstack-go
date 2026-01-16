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
