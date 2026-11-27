package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// GroupAuditAction is the kind of change a group audit entry records.
type GroupAuditAction string

const (
	GroupAuditActionUnspecified          GroupAuditAction = "GROUP_AUDIT_ACTION_UNSPECIFIED"
	GroupAuditActionGroupCreated         GroupAuditAction = "GROUP_AUDIT_ACTION_GROUP_CREATED"
	GroupAuditActionGroupUpdated         GroupAuditAction = "GROUP_AUDIT_ACTION_GROUP_UPDATED"
	GroupAuditActionMemberAdded          GroupAuditAction = "GROUP_AUDIT_ACTION_MEMBER_ADDED"
	GroupAuditActionMemberRemoved        GroupAuditAction = "GROUP_AUDIT_ACTION_MEMBER_REMOVED"
	GroupAuditActionMemberLeft           GroupAuditAction = "GROUP_AUDIT_ACTION_MEMBER_LEFT"
	GroupAuditActionMemberRoleChanged    GroupAuditAction = "GROUP_AUDIT_ACTION_MEMBER_ROLE_CHANGED"
	GroupAuditActionOwnershipTransferred GroupAuditAction = "GROUP_AUDIT_ACTION_OWNERSHIP_TRANSFERRED"
	GroupAuditActionJoinRequestDeclined  GroupAuditAction = "GROUP_AUDIT_ACTION_JOIN_REQUEST_DECLINED"
	GroupAuditActionInviteCreated        GroupAuditAction = "GROUP_AUDIT_ACTION_INVITE_CREATED"
	GroupAuditActionInviteRevoked        GroupAuditAction = "GROUP_AUDIT_ACTION_INVITE_REVOKED"
	GroupAuditActionInviteLinkCreated    GroupAuditAction = "GROUP_AUDIT_ACTION_INVITE_LINK_CREATED"
	GroupAuditActionInviteLinkRevoked    GroupAuditAction = "GROUP_AUDIT_ACTION_INVITE_LINK_REVOKED"
	GroupAuditActionInviteLinkRotated    GroupAuditAction = "GROUP_AUDIT_ACTION_INVITE_LINK_ROTATED"
	GroupAuditActionLinksReplaced        GroupAuditAction = "GROUP_AUDIT_ACTION_LINKS_REPLACED"
	GroupAuditActionPinAdded             GroupAuditAction = "GROUP_AUDIT_ACTION_PIN_ADDED"
	GroupAuditActionPinRemoved           GroupAuditAction = "GROUP_AUDIT_ACTION_PIN_REMOVED"
	GroupAuditActionPinsReordered        GroupAuditAction = "GROUP_AUDIT_ACTION_PINS_REORDERED"
	GroupAuditActionAvatarUpdated        GroupAuditAction = "GROUP_AUDIT_ACTION_AVATAR_UPDATED"
	GroupAuditActionAvatarRemoved        GroupAuditAction = "GROUP_AUDIT_ACTION_AVATAR_REMOVED"
	GroupAuditActionBannerUpdated        GroupAuditAction = "GROUP_AUDIT_ACTION_BANNER_UPDATED"
	GroupAuditActionBannerRemoved        GroupAuditAction = "GROUP_AUDIT_ACTION_BANNER_REMOVED"
	GroupAuditActionGroupSuspended       GroupAuditAction = "GROUP_AUDIT_ACTION_GROUP_SUSPENDED"
	GroupAuditActionGroupRestored        GroupAuditAction = "GROUP_AUDIT_ACTION_GROUP_RESTORED"
	GroupAuditActionMemberBanned         GroupAuditAction = "GROUP_AUDIT_ACTION_MEMBER_BANNED"
	GroupAuditActionMemberUnbanned       GroupAuditAction = "GROUP_AUDIT_ACTION_MEMBER_UNBANNED"
)

// GroupAuditActorType identifies who performed an audited action.
type GroupAuditActorType string

const (
	GroupAuditActorTypeUnspecified GroupAuditActorType = "GROUP_AUDIT_ACTOR_TYPE_UNSPECIFIED"
	GroupAuditActorTypeUser        GroupAuditActorType = "GROUP_AUDIT_ACTOR_TYPE_USER"
	GroupAuditActorTypeStaff       GroupAuditActorType = "GROUP_AUDIT_ACTOR_TYPE_STAFF"
	GroupAuditActorTypeService     GroupAuditActorType = "GROUP_AUDIT_ACTOR_TYPE_SERVICE"
)

// GroupAuditTargetType is the kind of object an audit entry acted on.
type GroupAuditTargetType string

const (
	GroupAuditTargetTypeUnspecified GroupAuditTargetType = "GROUP_AUDIT_TARGET_TYPE_UNSPECIFIED"
	GroupAuditTargetTypeGroup       GroupAuditTargetType = "GROUP_AUDIT_TARGET_TYPE_GROUP"
	GroupAuditTargetTypeMember      GroupAuditTargetType = "GROUP_AUDIT_TARGET_TYPE_MEMBER"
	GroupAuditTargetTypeInvite      GroupAuditTargetType = "GROUP_AUDIT_TARGET_TYPE_INVITE"
	GroupAuditTargetTypeInviteLink  GroupAuditTargetType = "GROUP_AUDIT_TARGET_TYPE_INVITE_LINK"
	GroupAuditTargetTypeJoinRequest GroupAuditTargetType = "GROUP_AUDIT_TARGET_TYPE_JOIN_REQUEST"
	GroupAuditTargetTypePin         GroupAuditTargetType = "GROUP_AUDIT_TARGET_TYPE_PIN"
)

// GroupAuditMembershipSource is how a member was added to a group.
type GroupAuditMembershipSource string

const (
	GroupAuditMembershipSourceUnspecified         GroupAuditMembershipSource = "GROUP_AUDIT_MEMBERSHIP_SOURCE_UNSPECIFIED"
	GroupAuditMembershipSourceGroupCreation       GroupAuditMembershipSource = "GROUP_AUDIT_MEMBERSHIP_SOURCE_GROUP_CREATION"
	GroupAuditMembershipSourceOpenJoin            GroupAuditMembershipSource = "GROUP_AUDIT_MEMBERSHIP_SOURCE_OPEN_JOIN"
	GroupAuditMembershipSourceJoinRequestApproval GroupAuditMembershipSource = "GROUP_AUDIT_MEMBERSHIP_SOURCE_JOIN_REQUEST_APPROVAL"
	GroupAuditMembershipSourceInvite              GroupAuditMembershipSource = "GROUP_AUDIT_MEMBERSHIP_SOURCE_INVITE"
	GroupAuditMembershipSourceInviteLink          GroupAuditMembershipSource = "GROUP_AUDIT_MEMBERSHIP_SOURCE_INVITE_LINK"
)

// GroupAuditCause explains an audit entry the actor did not choose directly.
type GroupAuditCause string

const (
	GroupAuditCauseUnspecified           GroupAuditCause = "GROUP_AUDIT_CAUSE_UNSPECIFIED"
	GroupAuditCauseAccountErasure        GroupAuditCause = "GROUP_AUDIT_CAUSE_ACCOUNT_ERASURE"
	GroupAuditCauseResourceAccessRevoked GroupAuditCause = "GROUP_AUDIT_CAUSE_RESOURCE_ACCESS_REVOKED"
	GroupAuditCauseResourceDeleted       GroupAuditCause = "GROUP_AUDIT_CAUSE_RESOURCE_DELETED"
	GroupAuditCauseInviteResent          GroupAuditCause = "GROUP_AUDIT_CAUSE_INVITE_RESENT"
	GroupAuditCauseMembershipRolledBack  GroupAuditCause = "GROUP_AUDIT_CAUSE_MEMBERSHIP_ROLLED_BACK"
)

// GroupAuditActor identifies who performed an audited action.
type GroupAuditActor struct {
	Type GroupAuditActorType `json:"type,omitempty"`
	// UserID is set only for user actors.
	UserID string `json:"userId,omitempty"`
}

// GroupAuditFieldChange names a group field changed by an update.
type GroupAuditFieldChange struct {
	Field         string `json:"field,omitempty"`
	PreviousValue string `json:"previousValue,omitempty"`
	NewValue      string `json:"newValue,omitempty"`
}

// GroupAuditDetail carries action-specific audit detail.
type GroupAuditDetail struct {
	PreviousRole     *GroupRole                 `json:"previousRole,omitempty"`
	NewRole          *GroupRole                 `json:"newRole,omitempty"`
	FieldChanges     []GroupAuditFieldChange    `json:"fieldChanges,omitempty"`
	MembershipSource GroupAuditMembershipSource `json:"membershipSource,omitempty"`
	Cause            GroupAuditCause            `json:"cause,omitempty"`
	SuspensionReason string                     `json:"suspensionReason,omitempty"`
	LinkCount        *int32                     `json:"linkCount,omitempty"`
	ResourceType     GroupResourceType          `json:"resourceType,omitempty"`
	ResourceID       string                     `json:"resourceId,omitempty"`
}

// GroupAuditEntry is one entry in a group's audit log.
type GroupAuditEntry struct {
	EntryID      string               `json:"entryId,omitempty"`
	GroupID      string               `json:"groupId,omitempty"`
	Action       GroupAuditAction     `json:"action,omitempty"`
	Actor        *GroupAuditActor     `json:"actor,omitempty"`
	TargetUserID string               `json:"targetUserId,omitempty"`
	TargetType   GroupAuditTargetType `json:"targetType,omitempty"`
	TargetID     string               `json:"targetId,omitempty"`
	Detail       *GroupAuditDetail    `json:"detail,omitempty"`
	OccurredAt   *time.Time           `json:"occurredAt,omitempty"`
}

// ListGroupAuditLogRequest pages through a group's audit log.
type ListGroupAuditLogRequest struct {
	GroupID   string
	PageSize  *int32
	NextToken string
	// Actions limits results to these actions; empty means all.
	Actions []GroupAuditAction
	// TargetUserID limits results to entries about this user.
	TargetUserID string
}

// ListGroupAuditLogResponse lists audit entries, newest first.
type ListGroupAuditLogResponse struct {
	Entries   []GroupAuditEntry `json:"entries,omitempty"`
	NextToken string            `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata  `json:"-"`
}

func (r *ListGroupAuditLogResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

// ListGroupAuditLog lists a group's audit log, newest first. Available to the
// group's admins and owner.
func (c *Client) ListGroupAuditLog(ctx context.Context, request *ListGroupAuditLogRequest, opts ...RequestOpt) (*ListGroupAuditLogResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	group, err := required(request.GroupID, "groupID")
	if err != nil {
		return nil, err
	}

	values := url.Values{}
	if request.PageSize != nil && *request.PageSize > 0 {
		values.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		values.Set("nextToken", token)
	}
	for _, action := range request.Actions {
		if action != "" {
			values.Add("actions", string(action))
		}
	}
	if target := strings.TrimSpace(request.TargetUserID); target != "" {
		values.Set("targetUserId", target)
	}
	path := "/v1/groups/" + url.PathEscape(group) + "/audit-log"
	if len(values) > 0 {
		path += "?" + values.Encode()
	}

	response := &ListGroupAuditLogResponse{}
	if err := c.doGroupJSON(ctx, http.MethodGet, path, nil, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
