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

// SyncAction represents common.v1.SyncAction.
type SyncAction string

const (
	SyncActionUnspecified SyncAction = "SYNC_ACTION_UNSPECIFIED"
	SyncActionUpsert      SyncAction = "UPSERT"
	SyncActionDelete      SyncAction = "DELETE"
)

// SyncStatus represents common.v1.SyncStatus.
type SyncStatus string

const (
	SyncStatusUnspecified SyncStatus = "SYNC_STATUS_UNSPECIFIED"
	SyncStatusOK          SyncStatus = "OK"
	SyncStatusConflict    SyncStatus = "CONFLICT"
	SyncStatusError       SyncStatus = "ERROR"
)

// SyncEventKind enumerates sync.v1.SyncEventKind values.
type SyncEventKind string

const (
	SyncEventKindUnspecified       SyncEventKind = "SYNC_EVENT_KIND_UNSPECIFIED"
	SyncEventKindCreated           SyncEventKind = "SYNC_EVENT_KIND_CREATED"
	SyncEventKindUpdated           SyncEventKind = "SYNC_EVENT_KIND_UPDATED"
	SyncEventKindDeleted           SyncEventKind = "SYNC_EVENT_KIND_DELETED"
	SyncEventKindPermissionGranted SyncEventKind = "SYNC_EVENT_KIND_PERMISSION_GRANTED"
	SyncEventKindPermissionRevoked SyncEventKind = "SYNC_EVENT_KIND_PERMISSION_REVOKED"
	SyncEventKindTagged            SyncEventKind = "SYNC_EVENT_KIND_TAGGED"
	SyncEventKindUntagged          SyncEventKind = "SYNC_EVENT_KIND_UNTAGGED"
	SyncEventKindSocial            SyncEventKind = "SYNC_EVENT_KIND_SOCIAL"
)

// TombstoneReason captures sync.v1.TombstoneReason values.
type TombstoneReason string

const (
	TombstoneReasonUnspecified       TombstoneReason = "TOMBSTONE_REASON_UNSPECIFIED"
	TombstoneReasonResourceDeleted   TombstoneReason = "TOMBSTONE_REASON_RESOURCE_DELETED"
	TombstoneReasonPermissionRevoked TombstoneReason = "TOMBSTONE_REASON_PERMISSION_REVOKED"
)

// SubscriptionSource mirrors sync.v1.SubscriptionSource.
type SubscriptionSource string

const (
	SubscriptionSourceUnspecified SubscriptionSource = "SUBSCRIPTION_SOURCE_UNSPECIFIED"
	SubscriptionSourceOwned       SubscriptionSource = "SUBSCRIPTION_SOURCE_OWNED"
	SubscriptionSourceManual      SubscriptionSource = "SUBSCRIPTION_SOURCE_MANUAL"
	SubscriptionSourceShared      SubscriptionSource = "SUBSCRIPTION_SOURCE_SHARED"
)

// SubscriptionMutationType mirrors sync.v1.SubscriptionMutation.MutationType.
type SubscriptionMutationType string

const (
	SubscriptionMutationTypeUnspecified SubscriptionMutationType = "MUTATION_TYPE_UNSPECIFIED"
	SubscriptionMutationTypeSubscribe   SubscriptionMutationType = "MUTATION_TYPE_SUBSCRIBE"
	SubscriptionMutationTypeUnsubscribe SubscriptionMutationType = "MUTATION_TYPE_UNSUBSCRIBE"
)

// ResourceDescriptor identifies a specific resource.
type ResourceDescriptor struct {
	Type ResourceType `json:"type,omitempty"`
	ID   string       `json:"id,omitempty"`
}

// GetChangeSetRequest configures SyncService.GetChangeSet.
type GetChangeSetRequest struct {
	Cursor           string
	PageSize         *int32
	IncludeDocuments bool
}

// SyncEvent mirrors sync.v1.SyncEvent.
type SyncEvent struct {
	EventID              string                      `json:"eventId,omitempty"`
	Resource             *ResourceDescriptor         `json:"resource,omitempty"`
	Kind                 SyncEventKind               `json:"kind,omitempty"`
	OccurredAt           *time.Time                  `json:"occurredAt,omitempty"`
	Version              string                      `json:"version,omitempty"`
	ProducerService      string                      `json:"producerService,omitempty"`
	Document             map[string]any              `json:"document,omitempty"`
	Tombstone            *Tombstone                  `json:"tombstone,omitempty"`
	PermissionRevocation *PermissionRevocationSignal `json:"permissionRevocation,omitempty"`
}

// Tombstone represents sync.v1.Tombstone.
type Tombstone struct {
	Reason TombstoneReason `json:"reason,omitempty"`
}

// PermissionRevocationSignal mirrors sync.v1.PermissionRevocationSignal.
type PermissionRevocationSignal struct {
	RevokedByUserID    string `json:"revokedByUserId,omitempty"`
	PreviousPermission string `json:"previousPermission,omitempty"`
	SubscriptionID     string `json:"subscriptionId,omitempty"`
}

// GetChangeSetResponse is returned from SyncService.GetChangeSet.
type GetChangeSetResponse struct {
	Events     []SyncEvent      `json:"events,omitempty"`
	NextCursor string           `json:"nextCursor,omitempty"`
	HasMore    bool             `json:"hasMore,omitempty"`
	Metadata   ResponseMetadata `json:"-"`
}

func (r *GetChangeSetResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// LocalChange mirrors sync.v1.LocalChange.
type LocalChange struct {
	ClientChangeID string              `json:"clientChangeId,omitempty"`
	Resource       *ResourceDescriptor `json:"resource,omitempty"`
	Action         SyncAction          `json:"action,omitempty"`
	Document       map[string]any      `json:"document,omitempty"`
	BaseVersion    string              `json:"baseVersion,omitempty"`
}

// BatchApplyChangesRequest configures SyncService.BatchApplyChanges.
type BatchApplyChangesRequest struct {
	Changes  []LocalChange `json:"changes,omitempty"`
	DeviceID string        `json:"deviceId,omitempty"`
}

// AppliedChangeResult mirrors sync.v1.AppliedChangeResult.
type AppliedChangeResult struct {
	ClientChangeID    string     `json:"clientChangeId,omitempty"`
	Status            SyncStatus `json:"status,omitempty"`
	LatestServerEvent *SyncEvent `json:"latestServerEvent,omitempty"`
	ErrorMessage      string     `json:"errorMessage,omitempty"`
}

// BatchApplyChangesResponse is returned from SyncService.BatchApplyChanges.
type BatchApplyChangesResponse struct {
	Results  []AppliedChangeResult `json:"results,omitempty"`
	Metadata ResponseMetadata      `json:"-"`
}

func (r *BatchApplyChangesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// SubscriptionMutation updates manual subscriptions.
type SubscriptionMutation struct {
	Type     SubscriptionMutationType `json:"type,omitempty"`
	Resource *ResourceDescriptor      `json:"resource,omitempty"`
}

// UpdateSubscriptionsRequest configures SyncService.UpdateSubscriptions.
type UpdateSubscriptionsRequest struct {
	Mutations []SubscriptionMutation `json:"mutations,omitempty"`
}

// ResourceSubscription mirrors sync.v1.ResourceSubscription.
type ResourceSubscription struct {
	SubscriptionID string              `json:"subscriptionId,omitempty"`
	Resource       *ResourceDescriptor `json:"resource,omitempty"`
	Source         SubscriptionSource  `json:"source,omitempty"`
	CreatedAt      *time.Time          `json:"createdAt,omitempty"`
}

// UpdateSubscriptionsResponse is returned from SyncService.UpdateSubscriptions.
type UpdateSubscriptionsResponse struct {
	Subscriptions []ResourceSubscription `json:"subscriptions,omitempty"`
	Metadata      ResponseMetadata       `json:"-"`
}

func (r *UpdateSubscriptionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListSubscriptionsResponse is returned from SyncService.ListSubscriptions.
type ListSubscriptionsResponse struct {
	Subscriptions []ResourceSubscription `json:"subscriptions,omitempty"`
	Metadata      ResponseMetadata       `json:"-"`
}

func (r *ListSubscriptionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetChangeSet retrieves ordered events after the provided cursor.
func (c *Client) GetChangeSet(ctx context.Context, request *GetChangeSetRequest, opts ...RequestOpt) (*GetChangeSetResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	values := url.Values{}
	if cursor := strings.TrimSpace(request.Cursor); cursor != "" {
		values.Set("cursor", cursor)
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		values.Set("pageSize", strconv.FormatInt(int64(*request.PageSize), 10))
	}
	if request.IncludeDocuments {
		values.Set("includeDocuments", strconv.FormatBool(true))
	}

	path := "/v1/sync/changeSet"
	if encoded := values.Encode(); encoded != "" {
		path = fmt.Sprintf("%s?%s", path, encoded)
	}

	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetChangeSetResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BatchApplyChanges replays local changes on the server.
func (c *Client) BatchApplyChanges(ctx context.Context, request *BatchApplyChangesRequest, opts ...RequestOpt) (*BatchApplyChangesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.Changes) == 0 {
		return nil, errors.New("changes must not be empty")
	}
	for i, change := range request.Changes {
		if err := validateLocalChange(change); err != nil {
			return nil, fmt.Errorf("changes[%d]: %w", i, err)
		}
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/sync:batchApply", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BatchApplyChangesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpdateSubscriptions mutates manual sync subscriptions.
func (c *Client) UpdateSubscriptions(ctx context.Context, request *UpdateSubscriptionsRequest, opts ...RequestOpt) (*UpdateSubscriptionsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.Mutations) == 0 {
		return nil, errors.New("mutations must not be empty")
	}
	for i, mutation := range request.Mutations {
		if err := validateSubscriptionMutation(mutation); err != nil {
			return nil, fmt.Errorf("mutations[%d]: %w", i, err)
		}
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/sync/subscriptions:update", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdateSubscriptionsResponse{}
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

func validateLocalChange(change LocalChange) error {
	if err := validateResourceDescriptor(change.Resource); err != nil {
		return err
	}
	if err := validateSyncAction(change.Action); err != nil {
		return err
	}
	if change.Action == SyncActionUpsert && len(change.Document) == 0 {
		return errors.New("document must be provided for UPSERT changes")
	}
	return nil
}

func validateResourceDescriptor(resource *ResourceDescriptor) error {
	if resource == nil {
		return errors.New("resource must be provided")
	}
	if err := validateResourceType(resource.Type); err != nil {
		return fmt.Errorf("resource.type: %w", err)
	}
	if strings.TrimSpace(resource.ID) == "" {
		return errors.New("resource.id must not be empty")
	}
	return nil
}

func validateSubscriptionMutation(mutation SubscriptionMutation) error {
	if err := validateSubscriptionMutationType(mutation.Type); err != nil {
		return err
	}
	return validateResourceDescriptor(mutation.Resource)
}

func validateSubscriptionMutationType(t SubscriptionMutationType) error {
	switch t {
	case SubscriptionMutationTypeSubscribe, SubscriptionMutationTypeUnsubscribe:
		return nil
	default:
		return errors.New("mutation type must be SUBSCRIBE or UNSUBSCRIBE")
	}
}

func validateSyncAction(action SyncAction) error {
	switch action {
	case SyncActionUpsert, SyncActionDelete:
		return nil
	default:
		return errors.New("action must be UPSERT or DELETE")
	}
}
