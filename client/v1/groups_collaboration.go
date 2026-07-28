package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type GroupStatus string

const (
	GroupStatusActive      GroupStatus = "GROUP_STATUS_ACTIVE"
	GroupStatusSuspended   GroupStatus = "GROUP_STATUS_SUSPENDED"
	GroupStatusDeactivated GroupStatus = "GROUP_STATUS_DEACTIVATED"
)

type GroupJoinPolicy string

const (
	GroupJoinPolicyOpen       GroupJoinPolicy = "GROUP_JOIN_POLICY_OPEN"
	GroupJoinPolicyRequest    GroupJoinPolicy = "GROUP_JOIN_POLICY_REQUEST"
	GroupJoinPolicyInviteOnly GroupJoinPolicy = "GROUP_JOIN_POLICY_INVITE_ONLY"
)

type GroupMemberDirectoryVisibility string

const (
	GroupMemberDirectoryVisibilityPublic  GroupMemberDirectoryVisibility = "GROUP_MEMBER_DIRECTORY_VISIBILITY_PUBLIC"
	GroupMemberDirectoryVisibilityMembers GroupMemberDirectoryVisibility = "GROUP_MEMBER_DIRECTORY_VISIBILITY_MEMBERS"
	GroupMemberDirectoryVisibilityAdmins  GroupMemberDirectoryVisibility = "GROUP_MEMBER_DIRECTORY_VISIBILITY_ADMINS"
)

type GroupRole string

const (
	GroupRoleMember GroupRole = "GROUP_ROLE_MEMBER"
	GroupRoleAdmin  GroupRole = "GROUP_ROLE_ADMIN"
	GroupRoleOwner  GroupRole = "GROUP_ROLE_OWNER"
)

type GroupJoinRequestStatus string

const (
	GroupJoinRequestStatusPending   GroupJoinRequestStatus = "GROUP_JOIN_REQUEST_STATUS_PENDING"
	GroupJoinRequestStatusApproved  GroupJoinRequestStatus = "GROUP_JOIN_REQUEST_STATUS_APPROVED"
	GroupJoinRequestStatusDeclined  GroupJoinRequestStatus = "GROUP_JOIN_REQUEST_STATUS_DECLINED"
	GroupJoinRequestStatusCancelled GroupJoinRequestStatus = "GROUP_JOIN_REQUEST_STATUS_CANCELLED"
)

type GroupResourceType string

const (
	GroupResourceTypeDeck       GroupResourceType = "GROUP_RESOURCE_TYPE_DECK"
	GroupResourceTypeCollection GroupResourceType = "GROUP_RESOURCE_TYPE_COLLECTION"
)

type GroupReportReason string

const (
	GroupReportReasonSpamOrScam       GroupReportReason = "GROUP_REPORT_REASON_SPAM_OR_SCAM"
	GroupReportReasonImpersonation    GroupReportReason = "GROUP_REPORT_REASON_IMPERSONATION"
	GroupReportReasonHarassmentOrHate GroupReportReason = "GROUP_REPORT_REASON_HARASSMENT_OR_HATE"
	GroupReportReasonSexualContent    GroupReportReason = "GROUP_REPORT_REASON_SEXUAL_CONTENT"
	GroupReportReasonViolence         GroupReportReason = "GROUP_REPORT_REASON_VIOLENCE"
	GroupReportReasonOther            GroupReportReason = "GROUP_REPORT_REASON_OTHER"
)

type GroupReportStatus string

const (
	GroupReportStatusOpen      GroupReportStatus = "GROUP_REPORT_STATUS_OPEN"
	GroupReportStatusInReview  GroupReportStatus = "GROUP_REPORT_STATUS_IN_REVIEW"
	GroupReportStatusResolved  GroupReportStatus = "GROUP_REPORT_STATUS_RESOLVED"
	GroupReportStatusDismissed GroupReportStatus = "GROUP_REPORT_STATUS_DISMISSED"
)

type GroupViewerContext struct {
	IsMember              bool      `json:"isMember,omitempty"`
	Role                  GroupRole `json:"role,omitempty"`
	HasPendingJoinRequest bool      `json:"hasPendingJoinRequest,omitempty"`
	CanManage             bool      `json:"canManage,omitempty"`
	CanInvite             bool      `json:"canInvite,omitempty"`
	CanJoin               bool      `json:"canJoin,omitempty"`
	CanReport             bool      `json:"canReport,omitempty"`
}
type GroupJoinRequest struct {
	RequestID        string                 `json:"requestId,omitempty"`
	GroupID          string                 `json:"groupId,omitempty"`
	UserID           string                 `json:"userId,omitempty"`
	Status           GroupJoinRequestStatus `json:"status,omitempty"`
	Message          string                 `json:"message,omitempty"`
	CreatedAt        *time.Time             `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time             `json:"updatedAt,omitempty"`
	ReviewedByUserID string                 `json:"reviewedByUserId,omitempty"`
}
type GroupInviteLink struct {
	InviteLinkID string     `json:"inviteLinkId,omitempty"`
	GroupID      string     `json:"groupId,omitempty"`
	Label        string     `json:"label,omitempty"`
	UseLimit     int32      `json:"useLimit,omitempty"`
	UseCount     int32      `json:"useCount,omitempty"`
	ExpiresAt    *time.Time `json:"expiresAt,omitempty"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	RevokedAt    *time.Time `json:"revokedAt,omitempty"`
}
type GroupPin struct {
	PinID          string            `json:"pinId,omitempty"`
	GroupID        string            `json:"groupId,omitempty"`
	ResourceType   GroupResourceType `json:"resourceType,omitempty"`
	ResourceID     string            `json:"resourceId,omitempty"`
	SortOrder      int32             `json:"sortOrder,omitempty"`
	PinnedByUserID string            `json:"pinnedByUserId,omitempty"`
	CreatedAt      *time.Time        `json:"createdAt,omitempty"`
}
type GroupReport struct {
	ReportID         string            `json:"reportId,omitempty"`
	GroupID          string            `json:"groupId,omitempty"`
	ReporterUserID   string            `json:"reporterUserId,omitempty"`
	Reason           GroupReportReason `json:"reason,omitempty"`
	Details          string            `json:"details,omitempty"`
	Status           GroupReportStatus `json:"status,omitempty"`
	AssignedToUserID string            `json:"assignedToUserId,omitempty"`
	ResolutionNotes  string            `json:"resolutionNotes,omitempty"`
	CreatedAt        *time.Time        `json:"createdAt,omitempty"`
	UpdatedAt        *time.Time        `json:"updatedAt,omitempty"`
}

type GroupActionResponse = DeleteGroupResponse
type GroupMembershipResponse struct {
	Membership    *GroupMember     `json:"membership,omitempty"`
	AlreadyMember bool             `json:"alreadyMember,omitempty"`
	Metadata      ResponseMetadata `json:"-"`
}

func (r *GroupMembershipResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

type GroupJoinRequestResponse struct {
	JoinRequest *GroupJoinRequest `json:"joinRequest,omitempty"`
	Metadata    ResponseMetadata  `json:"-"`
}

func (r *GroupJoinRequestResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

type ListGroupJoinRequestsResponse struct {
	JoinRequests []GroupJoinRequest `json:"joinRequests,omitempty"`
	NextToken    string             `json:"nextToken,omitempty"`
	Metadata     ResponseMetadata   `json:"-"`
}

func (r *ListGroupJoinRequestsResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

type ListGroupInvitesResponse struct {
	Invites   []GroupInvite    `json:"invites,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListGroupInvitesResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

type GroupInviteLinkResponse struct {
	InviteLink *GroupInviteLink `json:"inviteLink,omitempty"`
	Token      string           `json:"token,omitempty"`
	Metadata   ResponseMetadata `json:"-"`
}

func (r *GroupInviteLinkResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

type ListGroupInviteLinksResponse struct {
	InviteLinks []GroupInviteLink `json:"inviteLinks,omitempty"`
	NextToken   string            `json:"nextToken,omitempty"`
	Metadata    ResponseMetadata  `json:"-"`
}

func (r *ListGroupInviteLinksResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

type ReplaceGroupLinksResponse struct {
	Links    []GroupLink      `json:"links,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ReplaceGroupLinksResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

type ListGroupPinsResponse struct {
	Pins     []GroupPin       `json:"pins,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ListGroupPinsResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

type GroupPinResponse struct {
	Pin      *GroupPin        `json:"pin,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GroupPinResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

type GroupReportResponse struct {
	// Report is retained for compatibility during the group-report migration.
	Report   *GroupReport             `json:"report,omitempty"`
	Receipt  *ModerationReportReceipt `json:"receipt,omitempty"`
	Metadata ResponseMetadata         `json:"-"`
}

func (r *GroupReportResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

type ListGroupReportsResponse struct {
	Reports   []GroupReport    `json:"reports,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListGroupReportsResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

func (c *Client) doGroupJSON(ctx context.Context, method, path string, body any, response any, opts ...RequestOpt) error {
	var reader io.Reader
	if body != nil {
		encoded, err := jsonBody(body)
		if err != nil {
			return fmt.Errorf("encode body: %w", err)
		}
		reader = encoded
	}
	req, err := c.newRequest(ctx, method, path, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	return c.do(req, response, opts...)
}
func required(value, name string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", fmt.Errorf("%s must not be empty", name)
	}
	return value, nil
}
func addPageQuery(path string, pageSize *int32, nextToken string) string {
	values := url.Values{}
	if pageSize != nil && *pageSize > 0 {
		values.Set("pageSize", strconv.Itoa(int(*pageSize)))
	}
	if token := strings.TrimSpace(nextToken); token != "" {
		values.Set("nextToken", token)
	}
	if len(values) == 0 {
		return path
	}
	return path + "?" + values.Encode()
}

func (c *Client) GetGroupBySlug(ctx context.Context, slug string, opts ...RequestOpt) (*GetGroupResponse, error) {
	value, err := required(slug, "slug")
	if err != nil {
		return nil, err
	}
	response := &GetGroupResponse{}
	err = c.doGroupJSON(ctx, http.MethodGet, "/v1/groups/by-slug/"+url.PathEscape(value), nil, response, opts...)
	return response, err
}
func (c *Client) JoinGroup(ctx context.Context, groupID string, opts ...RequestOpt) (*GroupMembershipResponse, error) {
	value, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &GroupMembershipResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(value)+":join", map[string]any{}, response, opts...)
	return response, err
}
func (c *Client) LeaveGroup(ctx context.Context, groupID string, opts ...RequestOpt) (*GroupActionResponse, error) {
	value, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &GroupActionResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(value)+":leave", map[string]any{}, response, opts...)
	return response, err
}
func (c *Client) TransferGroupOwnership(ctx context.Context, groupID, newOwnerUserID string, opts ...RequestOpt) (*GroupActionResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	owner, err := required(newOwnerUserID, "newOwnerUserID")
	if err != nil {
		return nil, err
	}
	response := &GroupActionResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+":transferOwnership", map[string]string{"newOwnerUserId": owner}, response, opts...)
	return response, err
}
func (c *Client) RequestToJoinGroup(ctx context.Context, groupID, message string, opts ...RequestOpt) (*GroupJoinRequestResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &GroupJoinRequestResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/join-requests", map[string]string{"message": message}, response, opts...)
	return response, err
}
func (c *Client) CancelGroupJoinRequest(ctx context.Context, groupID string, opts ...RequestOpt) (*GroupActionResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &GroupActionResponse{}
	err = c.doGroupJSON(ctx, http.MethodDelete, "/v1/groups/"+url.PathEscape(group)+"/join-requests:mine", nil, response, opts...)
	return response, err
}
func (c *Client) ListGroupJoinRequests(ctx context.Context, groupID string, pageSize *int32, nextToken string, opts ...RequestOpt) (*ListGroupJoinRequestsResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &ListGroupJoinRequestsResponse{}
	path := addPageQuery("/v1/groups/"+url.PathEscape(group)+"/join-requests", pageSize, nextToken)
	err = c.doGroupJSON(ctx, http.MethodGet, path, nil, response, opts...)
	return response, err
}
func (c *Client) ReviewGroupJoinRequest(ctx context.Context, groupID, requestID string, decision GroupJoinRequestStatus, opts ...RequestOpt) (*GroupJoinRequestResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	requestValue, err := required(requestID, "requestID")
	if err != nil {
		return nil, err
	}
	response := &GroupJoinRequestResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/join-requests/"+url.PathEscape(requestValue)+":review", map[string]any{"decision": decision}, response, opts...)
	return response, err
}
func (c *Client) ListGroupInvites(ctx context.Context, groupID string, pageSize *int32, nextToken string, opts ...RequestOpt) (*ListGroupInvitesResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &ListGroupInvitesResponse{}
	err = c.doGroupJSON(ctx, http.MethodGet, addPageQuery("/v1/groups/"+url.PathEscape(group)+"/invites", pageSize, nextToken), nil, response, opts...)
	return response, err
}
func (c *Client) DeclineGroupInvite(ctx context.Context, token string, opts ...RequestOpt) (*GroupActionResponse, error) {
	value, err := required(token, "token")
	if err != nil {
		return nil, err
	}
	response := &GroupActionResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups:declineInvite", map[string]string{"token": value}, response, opts...)
	return response, err
}
func (c *Client) ResendGroupInvite(ctx context.Context, groupID, inviteID string, opts ...RequestOpt) (*GroupActionResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	invite, err := required(inviteID, "inviteID")
	if err != nil {
		return nil, err
	}
	response := &GroupActionResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/invites/"+url.PathEscape(invite)+":resend", map[string]any{}, response, opts...)
	return response, err
}
func (c *Client) CreateGroupInviteLink(ctx context.Context, groupID, label string, useLimit *int32, expiresAt *time.Time, opts ...RequestOpt) (*GroupInviteLinkResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &GroupInviteLinkResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/invite-links", map[string]any{"label": label, "useLimit": useLimit, "expiresAt": expiresAt}, response, opts...)
	return response, err
}
func (c *Client) ListGroupInviteLinks(ctx context.Context, groupID string, pageSize *int32, nextToken string, opts ...RequestOpt) (*ListGroupInviteLinksResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &ListGroupInviteLinksResponse{}
	err = c.doGroupJSON(ctx, http.MethodGet, addPageQuery("/v1/groups/"+url.PathEscape(group)+"/invite-links", pageSize, nextToken), nil, response, opts...)
	return response, err
}
func (c *Client) RevokeGroupInviteLink(ctx context.Context, groupID, linkID string, opts ...RequestOpt) (*GroupActionResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	link, err := required(linkID, "linkID")
	if err != nil {
		return nil, err
	}
	response := &GroupActionResponse{}
	err = c.doGroupJSON(ctx, http.MethodDelete, "/v1/groups/"+url.PathEscape(group)+"/invite-links/"+url.PathEscape(link), nil, response, opts...)
	return response, err
}
func (c *Client) RotateGroupInviteLink(ctx context.Context, groupID, linkID string, opts ...RequestOpt) (*GroupInviteLinkResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	link, err := required(linkID, "linkID")
	if err != nil {
		return nil, err
	}
	response := &GroupInviteLinkResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/invite-links/"+url.PathEscape(link)+":rotate", map[string]any{}, response, opts...)
	return response, err
}
func (c *Client) RedeemGroupInvite(ctx context.Context, token, groupID string, opts ...RequestOpt) (*GroupMembershipResponse, error) {
	tokenValue, err := required(token, "token")
	if err != nil {
		return nil, err
	}
	response := &GroupMembershipResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups:redeemInvite", map[string]string{"token": tokenValue, "groupId": groupID}, response, opts...)
	return response, err
}
func (c *Client) ReplaceGroupLinks(ctx context.Context, groupID string, links []GroupLink, opts ...RequestOpt) (*ReplaceGroupLinksResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &ReplaceGroupLinksResponse{}
	err = c.doGroupJSON(ctx, http.MethodPut, "/v1/groups/"+url.PathEscape(group)+"/links", map[string]any{"links": links}, response, opts...)
	return response, err
}
func (c *Client) ListGroupPins(ctx context.Context, groupID string, opts ...RequestOpt) (*ListGroupPinsResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &ListGroupPinsResponse{}
	err = c.doGroupJSON(ctx, http.MethodGet, "/v1/groups/"+url.PathEscape(group)+"/pins", nil, response, opts...)
	return response, err
}
func (c *Client) PinGroupResource(ctx context.Context, groupID string, resourceType GroupResourceType, resourceID string, opts ...RequestOpt) (*GroupPinResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	resourceValue, err := required(resourceID, "resourceID")
	if err != nil {
		return nil, err
	}
	response := &GroupPinResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/pins", map[string]any{"resourceType": resourceType, "resourceId": resourceValue}, response, opts...)
	return response, err
}
func (c *Client) UnpinGroupResource(ctx context.Context, groupID, pinID string, opts ...RequestOpt) (*GroupActionResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	pin, err := required(pinID, "pinID")
	if err != nil {
		return nil, err
	}
	response := &GroupActionResponse{}
	err = c.doGroupJSON(ctx, http.MethodDelete, "/v1/groups/"+url.PathEscape(group)+"/pins/"+url.PathEscape(pin), nil, response, opts...)
	return response, err
}
func (c *Client) ReorderGroupPins(ctx context.Context, groupID string, pinIDs []string, opts ...RequestOpt) (*ListGroupPinsResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &ListGroupPinsResponse{}
	err = c.doGroupJSON(ctx, http.MethodPut, "/v1/groups/"+url.PathEscape(group)+"/pins:reorder", map[string]any{"pinIds": pinIDs}, response, opts...)
	return response, err
}
func (c *Client) BeginGroupBackgroundUpload(ctx context.Context, request *BeginGroupAvatarUploadRequest, opts ...RequestOpt) (*BeginAvatarUploadResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	group, err := required(request.GroupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &BeginAvatarUploadResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/background:beginUpload", map[string]any{"contentType": request.ContentType, "contentLength": request.ContentLength}, response, opts...)
	return response, err
}
func (c *Client) CompleteGroupBackgroundUpload(ctx context.Context, request *CompleteGroupAvatarUploadRequest, opts ...RequestOpt) (*CompleteGroupAvatarUploadResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	group, err := required(request.GroupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &CompleteGroupAvatarUploadResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/background:completeUpload", map[string]string{"uploadId": request.UploadID}, response, opts...)
	return response, err
}
func (c *Client) ClearGroupBackground(ctx context.Context, groupID string, opts ...RequestOpt) (*UpdateGroupResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &UpdateGroupResponse{}
	err = c.doGroupJSON(ctx, http.MethodDelete, "/v1/groups/"+url.PathEscape(group)+"/background", nil, response, opts...)
	return response, err
}
func (c *Client) DeleteGroupAvatar(ctx context.Context, groupID string, opts ...RequestOpt) (*UpdateGroupResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &UpdateGroupResponse{}
	err = c.doGroupJSON(ctx, http.MethodDelete, "/v1/groups/"+url.PathEscape(group)+"/avatar", nil, response, opts...)
	return response, err
}
func (c *Client) ReportGroup(ctx context.Context, groupID string, reason GroupReportReason, details string, opts ...RequestOpt) (*GroupReportResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &GroupReportResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+":report", map[string]any{"reason": reason, "details": details}, response, opts...)
	return response, err
}
func (c *Client) ListGroupReports(ctx context.Context, statusFilter *GroupReportStatus, pageSize *int32, nextToken string, opts ...RequestOpt) (*ListGroupReportsResponse, error) {
	values := url.Values{}
	if statusFilter != nil {
		values.Set("status", string(*statusFilter))
	}
	if pageSize != nil {
		values.Set("pageSize", strconv.Itoa(int(*pageSize)))
	}
	if nextToken != "" {
		values.Set("nextToken", nextToken)
	}
	path := "/v1/admin/groups/reports"
	if len(values) > 0 {
		path += "?" + values.Encode()
	}
	response := &ListGroupReportsResponse{}
	err := c.doGroupJSON(ctx, http.MethodGet, path, nil, response, opts...)
	return response, err
}
func (c *Client) UpdateGroupReport(ctx context.Context, reportID string, reportStatus *GroupReportStatus, assignedToUserID, resolutionNotes *string, updateMask string, opts ...RequestOpt) (*GroupReportResponse, error) {
	report, err := required(reportID, "reportID")
	if err != nil {
		return nil, err
	}
	response := &GroupReportResponse{}
	err = c.doGroupJSON(ctx, http.MethodPatch, "/v1/admin/groups/reports/"+url.PathEscape(report), map[string]any{"status": reportStatus, "assignedToUserId": assignedToUserID, "resolutionNotes": resolutionNotes, "updateMask": updateMask}, response, opts...)
	return response, err
}
func (c *Client) SuspendGroup(ctx context.Context, groupID, reason string, opts ...RequestOpt) (*UpdateGroupResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &UpdateGroupResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/admin/groups/"+url.PathEscape(group)+":suspend", map[string]string{"reason": reason}, response, opts...)
	return response, err
}
func (c *Client) RestoreGroup(ctx context.Context, groupID string, opts ...RequestOpt) (*UpdateGroupResponse, error) {
	group, err := required(groupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &UpdateGroupResponse{}
	err = c.doGroupJSON(ctx, http.MethodPost, "/v1/admin/groups/"+url.PathEscape(group)+":restore", map[string]any{}, response, opts...)
	return response, err
}
