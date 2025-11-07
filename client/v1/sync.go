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

// SyncAction represents v1SyncAction.
type SyncAction string

const (
	SyncActionUnspecified SyncAction = "SYNC_ACTION_UNSPECIFIED"
	SyncActionUpsert      SyncAction = "UPSERT"
	SyncActionDelete      SyncAction = "DELETE"
)

// SyncStatus represents v1SyncStatus.
type SyncStatus string

const (
	SyncStatusUnspecified SyncStatus = "SYNC_STATUS_UNSPECIFIED"
	SyncStatusOK          SyncStatus = "OK"
	SyncStatusConflict    SyncStatus = "CONFLICT"
	SyncStatusError       SyncStatus = "ERROR"
)

// SubscriptionReason matches v1SubscriptionReason.
type SubscriptionReason string

const (
	SubscriptionReasonUnspecified SubscriptionReason = "SUBSCRIPTION_REASON_UNSPECIFIED"
	SubscriptionReasonOwner       SubscriptionReason = "OWNER"
	SubscriptionReasonPinned      SubscriptionReason = "PINNED"
	SubscriptionReasonPermission  SubscriptionReason = "PERMISSION"
)

// Cursor models v1Cursor.
type Cursor struct {
	Overall     string            `json:"overall,omitempty"`
	PerResource map[string]string `json:"perResource,omitempty"`
}

// PullFilters represent v1PullFilters.
type PullFilters struct {
	ResourceTypes  []ResourceType `json:"resourceTypes,omitempty"`
	PinnedOnly     bool           `json:"pinnedOnly,omitempty"`
	IncludeDeletes bool           `json:"includeDeletes,omitempty"`
}

// PullChangesRequest configures SyncService.PullChanges.
type PullChangesRequest struct {
	Cursor  *Cursor
	Filters *PullFilters
	Limit   *int32
}

// Change mirrors v1Change.
type Change struct {
	EnvelopeID   string            `json:"envelopeId,omitempty"`
	ResourceType ResourceType      `json:"resourceType,omitempty"`
	ResourceID   string            `json:"resourceId,omitempty"`
	Action       SyncAction        `json:"action,omitempty"`
	Version      string            `json:"version,omitempty"`
	Payload      map[string]any    `json:"payload,omitempty"`
	EmittedAt    *time.Time        `json:"emittedAt,omitempty"`
	RecordedAt   *time.Time        `json:"recordedAt,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// PullChangesResponse is returned from SyncService.PullChanges.
type PullChangesResponse struct {
	Changes    []Change         `json:"changes,omitempty"`
	NextCursor *Cursor          `json:"nextCursor,omitempty"`
	HasMore    bool             `json:"hasMore,omitempty"`
	Metadata   ResponseMetadata `json:"-"`
}

func (r *PullChangesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// Mutation mirrors v1Mutation.
type Mutation struct {
	LocalChangeID   string         `json:"localChangeId,omitempty"`
	ResourceType    ResourceType   `json:"resourceType,omitempty"`
	Action          SyncAction     `json:"action,omitempty"`
	ResourceID      string         `json:"resourceId,omitempty"`
	Payload         map[string]any `json:"payload,omitempty"`
	BaseVersion     string         `json:"baseVersion,omitempty"`
	ClientTimestamp *time.Time     `json:"clientTimestamp,omitempty"`
}

// BatchApplyMutationsRequest configures SyncService.BatchApplyMutations.
type BatchApplyMutationsRequest struct {
	ClientID  string     `json:"clientId,omitempty"`
	Mutations []Mutation `json:"mutations,omitempty"`
}

// ResourceSnapshot mirrors v1ResourceSnapshot.
type ResourceSnapshot struct {
	ResourceID   string         `json:"resourceId,omitempty"`
	ResourceType ResourceType   `json:"resourceType,omitempty"`
	Version      string         `json:"version,omitempty"`
	Payload      map[string]any `json:"payload,omitempty"`
	UpdatedAt    *time.Time     `json:"updatedAt,omitempty"`
}

// MutationResult mirrors v1MutationResult.
type MutationResult struct {
	LocalChangeID string            `json:"localChangeId,omitempty"`
	Status        SyncStatus        `json:"status,omitempty"`
	Message       string            `json:"message,omitempty"`
	Snapshot      *ResourceSnapshot `json:"snapshot,omitempty"`
}

// BatchApplyMutationsResponse is returned from SyncService.BatchApplyMutations.
type BatchApplyMutationsResponse struct {
	Results  []MutationResult `json:"results,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *BatchApplyMutationsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// AddPinRequest registers a sync pin.
type AddPinRequest struct {
	ResourceID string
}

// AddPinResponse captures metadata for pin operations.
type AddPinResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *AddPinResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RemovePinRequest removes a sync pin.
type RemovePinRequest struct {
	ResourceID string
}

// RemovePinResponse captures metadata for pin removals.
type RemovePinResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *RemovePinResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// Subscription mirrors v1Subscription.
type Subscription struct {
	ResourceID   string             `json:"resourceId,omitempty"`
	ResourceType ResourceType       `json:"resourceType,omitempty"`
	Reason       SubscriptionReason `json:"reason,omitempty"`
	GrantedAt    *time.Time         `json:"grantedAt,omitempty"`
}

// ListSubscriptionsResponse is returned from SyncService.ListSubscriptions.
type ListSubscriptionsResponse struct {
	Subscriptions []Subscription   `json:"subscriptions,omitempty"`
	Metadata      ResponseMetadata `json:"-"`
}

func (r *ListSubscriptionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// PullChanges retrieves changes since the supplied cursor for subscribed resources.
func (c *Client) PullChanges(ctx context.Context, request *PullChangesRequest, opts ...RequestOpt) (*PullChangesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := validatePullFilters(request.Filters); err != nil {
		return nil, fmt.Errorf("filters: %w", err)
	}

	payload := struct {
		Cursor  *Cursor      `json:"cursor,omitempty"`
		Filters *PullFilters `json:"filters,omitempty"`
		Limit   *int32       `json:"limit,omitempty"`
	}{
		Cursor:  request.Cursor,
		Filters: request.Filters,
	}

	if request.Limit != nil && *request.Limit > 0 {
		payload.Limit = request.Limit
	}

	body, err := jsonBody(payload)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/sync/changes:pull", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &PullChangesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BatchApplyMutations applies mutations recorded locally to the authoritative services.
func (c *Client) BatchApplyMutations(ctx context.Context, request *BatchApplyMutationsRequest, opts ...RequestOpt) (*BatchApplyMutationsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.Mutations) == 0 {
		return nil, errors.New("mutations must not be empty")
	}
	for i, mutation := range request.Mutations {
		if err := validateMutation(mutation); err != nil {
			return nil, fmt.Errorf("mutations[%d]: %w", i, err)
		}
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/sync/mutations:batchApply", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BatchApplyMutationsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// AddPin registers a pin so resources always sync for the caller.
func (c *Client) AddPin(ctx context.Context, request *AddPinRequest, opts ...RequestOpt) (*AddPinResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	resourceID := strings.TrimSpace(request.ResourceID)
	if resourceID == "" {
		return nil, errors.New("resourceID must not be empty")
	}

	body, err := jsonBody(struct {
		ResourceID string `json:"resourceId,omitempty"`
	}{
		ResourceID: resourceID,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/sync/pins", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &AddPinResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// RemovePin deletes a previously registered pin.
func (c *Client) RemovePin(ctx context.Context, request *RemovePinRequest, opts ...RequestOpt) (*RemovePinResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	resourceID := strings.TrimSpace(request.ResourceID)
	if resourceID == "" {
		return nil, errors.New("resourceID must not be empty")
	}

	path := fmt.Sprintf("/v1/sync/pins/%s", url.PathEscape(resourceID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &RemovePinResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListSubscriptions enumerates sync subscriptions for the caller.
func (c *Client) ListSubscriptions(ctx context.Context, opts ...RequestOpt) (*ListSubscriptionsResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/sync/subscriptions", nil)
	if err != nil {
		return nil, err
	}

	response := &ListSubscriptionsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

func validatePullFilters(filters *PullFilters) error {
	if filters == nil {
		return nil
	}
	for i, resourceType := range filters.ResourceTypes {
		if err := validateResourceType(resourceType); err != nil {
			return fmt.Errorf("resourceTypes[%d]: %w", i, err)
		}
	}
	return nil
}

func validateMutation(mutation Mutation) error {
	if err := validateResourceType(mutation.ResourceType); err != nil {
		return fmt.Errorf("resourceType: %w", err)
	}
	if err := validateSyncAction(mutation.Action); err != nil {
		return err
	}
	if strings.TrimSpace(mutation.ResourceID) == "" {
		return errors.New("resourceId must not be empty")
	}
	return nil
}

func validateSyncAction(action SyncAction) error {
	switch action {
	case SyncActionUpsert, SyncActionDelete:
		return nil
	default:
		return errors.New("action must be UPSERT or DELETE")
	}
}
