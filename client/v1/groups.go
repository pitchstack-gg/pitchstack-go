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

// CreateGroupRequest provisions a new user group.
type CreateGroupRequest struct {
	Slug                      string                         `json:"slug,omitempty"`
	Name                      string                         `json:"name,omitempty"`
	Description               string                         `json:"description,omitempty"`
	Visibility                VisibilityLevel                `json:"visibility,omitempty"`
	JoinPolicy                GroupJoinPolicy                `json:"joinPolicy,omitempty"`
	MemberDirectoryVisibility GroupMemberDirectoryVisibility `json:"memberDirectoryVisibility,omitempty"`
	BackgroundColor           string                         `json:"backgroundColor,omitempty"`
}

// CreateGroupResponse includes the created group.
type CreateGroupResponse struct {
	Group    *UserGroup       `json:"group,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CreateGroupResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateGroupRequest modifies an existing group.
type UpdateGroupRequest struct {
	GroupID                   string                          `json:"-"`
	Name                      *string                         `json:"name,omitempty"`
	Description               *string                         `json:"description,omitempty"`
	Visibility                *VisibilityLevel                `json:"visibility,omitempty"`
	JoinPolicy                *GroupJoinPolicy                `json:"joinPolicy,omitempty"`
	MemberDirectoryVisibility *GroupMemberDirectoryVisibility `json:"memberDirectoryVisibility,omitempty"`
	BackgroundColor           *string                         `json:"backgroundColor,omitempty"`
	UpdateMask                string                          `json:"updateMask,omitempty"`
}

// UpdateGroupResponse contains the updated group.
type UpdateGroupResponse struct {
	Group    *UserGroup       `json:"group,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateGroupResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteGroupRequest identifies the group to remove.
type DeleteGroupRequest struct {
	GroupID string `json:"-"`
}

// DeleteGroupResponse records metadata for the delete call.
type DeleteGroupResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteGroupResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListGroupsRequest enumerates groups with pagination.
type ListGroupsRequest struct {
	PageSize  *int32
	NextToken string
}

// ListGroupsResponse lists groups visible to the caller.
type ListGroupsResponse struct {
	Groups    []UserGroup      `json:"groups,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListGroupsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// AddGroupMemberRequest adds a member to a group.
type AddGroupMemberRequest struct {
	GroupID   string     `json:"-"`
	UserID    string     `json:"userId,omitempty"`
	Role      string     `json:"role,omitempty"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

// AddGroupMemberResponse records metadata for the add operation.
type AddGroupMemberResponse struct {
	Invite   *GroupInvite     `json:"invite,omitempty"`
	Token    string           `json:"token,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *AddGroupMemberResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RemoveGroupMemberRequest removes a user from a group.
type RemoveGroupMemberRequest struct {
	GroupID string `json:"-"`
	UserID  string `json:"-"`
}

// RemoveGroupMemberResponse captures metadata for the removal.
type RemoveGroupMemberResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *RemoveGroupMemberResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListGroupMembersRequest enumerates members within a group.
type ListGroupMembersRequest struct {
	GroupID   string
	PageSize  *int32
	NextToken string
}

// ListGroupMembersResponse returns group membership.
type ListGroupMembersResponse struct {
	Members   []GroupMember    `json:"members,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListGroupMembersResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UserGroup represents a user-defined group construct.
type UserGroup struct {
	GroupID                   string                         `json:"groupId,omitempty"`
	Slug                      string                         `json:"slug,omitempty"`
	Name                      string                         `json:"name,omitempty"`
	Description               string                         `json:"description,omitempty"`
	OwnerID                   string                         `json:"ownerId,omitempty"`
	Visibility                VisibilityLevel                `json:"visibility,omitempty"`
	IsActive                  bool                           `json:"isActive,omitempty"`
	CreatedByUserID           string                         `json:"createdByUserId,omitempty"`
	CreatedAt                 *time.Time                     `json:"createdAt,omitempty"`
	UpdatedAt                 *time.Time                     `json:"updatedAt,omitempty"`
	AvatarURL                 string                         `json:"avatarUrl,omitempty"`
	OwnerUserID               string                         `json:"ownerUserId,omitempty"`
	Status                    GroupStatus                    `json:"status,omitempty"`
	JoinPolicy                GroupJoinPolicy                `json:"joinPolicy,omitempty"`
	MemberDirectoryVisibility GroupMemberDirectoryVisibility `json:"memberDirectoryVisibility,omitempty"`
	BackgroundColor           string                         `json:"backgroundColor,omitempty"`
	BackgroundURL             string                         `json:"backgroundUrl,omitempty"`
	MemberCount               int32                          `json:"memberCount,omitempty"`
}

// GroupLink represents v1GroupLink.
type GroupLink struct {
	LinkID    string `json:"linkId,omitempty"`
	LinkType  string `json:"linkType,omitempty"`
	URL       string `json:"url,omitempty"`
	Label     string `json:"label,omitempty"`
	SortOrder int32  `json:"sortOrder,omitempty"`
}

// GroupMember captures membership details for a group.
type GroupMember struct {
	GroupID   string     `json:"groupId,omitempty"`
	UserID    string     `json:"userId,omitempty"`
	Role      string     `json:"role,omitempty"`
	JoinedAt  *time.Time `json:"joinedAt,omitempty"`
	RoleValue GroupRole  `json:"roleValue,omitempty"`
}

// GetGroupRequest identifies the group to retrieve.
type GetGroupRequest struct {
	GroupID string
}

// GetGroupResponse returns a group and related links.
type GetGroupResponse struct {
	Group         *UserGroup          `json:"group,omitempty"`
	Links         []GroupLink         `json:"links,omitempty"`
	ViewerContext *GroupViewerContext `json:"viewerContext,omitempty"`
	Metadata      ResponseMetadata    `json:"-"`
}

func (r *GetGroupResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// SearchGroupsRequest searches for groups.
type SearchGroupsRequest struct {
	SearchTerm string
	PageSize   *int32
	NextToken  string
}

// SearchGroupsResponse lists groups matching the query.
type SearchGroupsResponse struct {
	Groups    []UserGroup      `json:"groups,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *SearchGroupsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListMyGroupsRequest paginates ListMyGroups.
type ListMyGroupsRequest struct {
	PageSize  *int32
	NextToken string
}

// ListMyGroupsResponse lists groups the authenticated user belongs to.
type ListMyGroupsResponse struct {
	Groups    []UserGroup      `json:"groups,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListMyGroupsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GroupInvite represents v1GroupInvite.
type GroupInvite struct {
	InviteID      string     `json:"inviteId,omitempty"`
	GroupID       string     `json:"groupId,omitempty"`
	InvitedUserID string     `json:"invitedUserId,omitempty"`
	Role          string     `json:"role,omitempty"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
	CreatedAt     *time.Time `json:"createdAt,omitempty"`
	RoleValue     GroupRole  `json:"roleValue,omitempty"`
}

// CreateGroupInviteRequest creates a group invite.
type CreateGroupInviteRequest struct {
	GroupID       string     `json:"-"`
	InvitedUserID string     `json:"invitedUserId,omitempty"`
	Role          string     `json:"role,omitempty"`
	ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
}

// CreateGroupInviteResponse returns the created invite and token (once).
type CreateGroupInviteResponse struct {
	Invite   *GroupInvite     `json:"invite,omitempty"`
	Token    string           `json:"token,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CreateGroupInviteResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RevokeGroupInviteRequest identifies the invite to revoke.
type RevokeGroupInviteRequest struct {
	GroupID  string
	InviteID string
}

// RevokeGroupInviteResponse captures metadata for revoke operations.
type RevokeGroupInviteResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *RevokeGroupInviteResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// AcceptGroupInviteRequest accepts an invite token into a group.
type AcceptGroupInviteRequest struct {
	Token   string `json:"token,omitempty"`
	GroupID string `json:"groupId,omitempty"`
}

// AcceptGroupInviteResponse captures metadata for accept operations.
type AcceptGroupInviteResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *AcceptGroupInviteResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateGroupMemberRoleRequest updates a group member's role.
type UpdateGroupMemberRoleRequest struct {
	GroupID string `json:"-"`
	UserID  string `json:"-"`
	Role    string `json:"role,omitempty"`
}

// UpdateGroupMemberRoleResponse captures metadata for role updates.
type UpdateGroupMemberRoleResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateGroupMemberRoleResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BeginGroupAvatarUploadRequest starts a server-controlled avatar upload for a group.
type BeginGroupAvatarUploadRequest struct {
	GroupID       string `json:"-"`
	ContentType   string `json:"contentType,omitempty"`
	ContentLength *int64 `json:"contentLength,omitempty,string"`
}

// CompleteGroupAvatarUploadRequest completes an avatar upload for a group.
type CompleteGroupAvatarUploadRequest struct {
	GroupID  string `json:"-"`
	UploadID string `json:"uploadId,omitempty"`
}

// CompleteGroupAvatarUploadResponse returns the updated group.
type CompleteGroupAvatarUploadResponse struct {
	Group    *UserGroup       `json:"group,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CompleteGroupAvatarUploadResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateGroup creates a new user group.
func (c *Client) CreateGroup(ctx context.Context, request *CreateGroupRequest, opts ...RequestOpt) (*CreateGroupResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/groups", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CreateGroupResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateGroup updates an existing user group.
func (c *Client) UpdateGroup(ctx context.Context, request *UpdateGroupRequest, opts ...RequestOpt) (*UpdateGroupResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return nil, errors.New("groupID must not be empty")
	}

	body, err := jsonBody(struct {
		Name                      *string                         `json:"name,omitempty"`
		Description               *string                         `json:"description,omitempty"`
		Visibility                *VisibilityLevel                `json:"visibility,omitempty"`
		JoinPolicy                *GroupJoinPolicy                `json:"joinPolicy,omitempty"`
		MemberDirectoryVisibility *GroupMemberDirectoryVisibility `json:"memberDirectoryVisibility,omitempty"`
		BackgroundColor           *string                         `json:"backgroundColor,omitempty"`
		UpdateMask                string                          `json:"updateMask,omitempty"`
	}{
		Name:                      request.Name,
		Description:               request.Description,
		Visibility:                request.Visibility,
		JoinPolicy:                request.JoinPolicy,
		MemberDirectoryVisibility: request.MemberDirectoryVisibility,
		BackgroundColor:           request.BackgroundColor,
		UpdateMask:                request.UpdateMask,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/groups/%s", url.PathEscape(groupID))
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateGroupResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteGroup removes a group by ID.
func (c *Client) DeleteGroup(ctx context.Context, request *DeleteGroupRequest, opts ...RequestOpt) (*DeleteGroupResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return nil, errors.New("groupID must not be empty")
	}

	path := fmt.Sprintf("/v1/groups/%s", url.PathEscape(groupID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteGroupResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListGroups enumerates groups available to the caller.
func (c *Client) ListGroups(ctx context.Context, request *ListGroupsRequest, opts ...RequestOpt) (*ListGroupsResponse, error) {
	if request == nil {
		request = &ListGroupsRequest{}
	}
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/groups", nil)
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	req.URL.RawQuery = query.Encode()
	response := &ListGroupsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// AddGroupMember grants membership to a user within a group.
func (c *Client) AddGroupMember(ctx context.Context, request *AddGroupMemberRequest, opts ...RequestOpt) (*AddGroupMemberResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	createResp, err := c.CreateGroupInvite(ctx, &CreateGroupInviteRequest{
		GroupID:       request.GroupID,
		InvitedUserID: request.UserID,
		Role:          request.Role,
		ExpiresAt:     request.ExpiresAt,
	}, opts...)
	if err != nil {
		return nil, err
	}

	return &AddGroupMemberResponse{
		Invite:   createResp.Invite,
		Token:    createResp.Token,
		Metadata: createResp.Metadata,
	}, nil
}

// RemoveGroupMember revokes membership for a user within a group.
func (c *Client) RemoveGroupMember(ctx context.Context, request *RemoveGroupMemberRequest, opts ...RequestOpt) (*RemoveGroupMemberResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return nil, errors.New("groupID must not be empty")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/v1/groups/%s/members/%s", url.PathEscape(groupID), url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &RemoveGroupMemberResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListGroupMembers returns members within a group.
func (c *Client) ListGroupMembers(ctx context.Context, request *ListGroupMembersRequest, opts ...RequestOpt) (*ListGroupMembersResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return nil, errors.New("groupID must not be empty")
	}

	path := fmt.Sprintf("/v1/groups/%s/members", url.PathEscape(groupID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	req.URL.RawQuery = query.Encode()

	response := &ListGroupMembersResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetGroup retrieves a group by ID.
func (c *Client) GetGroup(ctx context.Context, request *GetGroupRequest, opts ...RequestOpt) (*GetGroupResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return nil, errors.New("groupID must not be empty")
	}

	path := fmt.Sprintf("/v1/groups/%s", url.PathEscape(groupID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetGroupResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// SearchGroups searches for groups by search term.
func (c *Client) SearchGroups(ctx context.Context, request *SearchGroupsRequest, opts ...RequestOpt) (*SearchGroupsResponse, error) {
	if request == nil {
		request = &SearchGroupsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/groups/search", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if term := strings.TrimSpace(request.SearchTerm); term != "" {
		query.Set("searchTerm", term)
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

	response := &SearchGroupsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListMyGroups lists groups the authenticated user belongs to.
func (c *Client) ListMyGroups(ctx context.Context, request *ListMyGroupsRequest, opts ...RequestOpt) (*ListMyGroupsResponse, error) {
	if request == nil {
		request = &ListMyGroupsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/groups:mine", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	response := &ListMyGroupsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CreateGroupInvite creates an invite for a user to join a group.
func (c *Client) CreateGroupInvite(ctx context.Context, request *CreateGroupInviteRequest, opts ...RequestOpt) (*CreateGroupInviteResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return nil, errors.New("groupID must not be empty")
	}
	invitedUserID := strings.TrimSpace(request.InvitedUserID)
	if invitedUserID == "" {
		return nil, errors.New("invitedUserID must not be empty")
	}
	role := strings.TrimSpace(request.Role)
	if role == "" {
		return nil, errors.New("role must not be empty")
	}

	body, err := jsonBody(struct {
		InvitedUserID string     `json:"invitedUserId,omitempty"`
		Role          string     `json:"role,omitempty"`
		ExpiresAt     *time.Time `json:"expiresAt,omitempty"`
	}{
		InvitedUserID: invitedUserID,
		Role:          role,
		ExpiresAt:     request.ExpiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/groups/%s/invites", url.PathEscape(groupID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CreateGroupInviteResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// RevokeGroupInvite revokes an existing group invite.
func (c *Client) RevokeGroupInvite(ctx context.Context, request *RevokeGroupInviteRequest, opts ...RequestOpt) (*RevokeGroupInviteResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return nil, errors.New("groupID must not be empty")
	}
	inviteID := strings.TrimSpace(request.InviteID)
	if inviteID == "" {
		return nil, errors.New("inviteID must not be empty")
	}

	path := fmt.Sprintf("/v1/groups/%s/invites/%s", url.PathEscape(groupID), url.PathEscape(inviteID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &RevokeGroupInviteResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// AcceptGroupInvite accepts an invite token.
func (c *Client) AcceptGroupInvite(ctx context.Context, request *AcceptGroupInviteRequest, opts ...RequestOpt) (*AcceptGroupInviteResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.Token) == "" {
		return nil, errors.New("token must not be empty")
	}
	if strings.TrimSpace(request.GroupID) == "" {
		return nil, errors.New("groupID must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/groups:acceptInvite", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &AcceptGroupInviteResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpdateGroupMemberRole updates a member's role within a group.
func (c *Client) UpdateGroupMemberRole(ctx context.Context, request *UpdateGroupMemberRoleRequest, opts ...RequestOpt) (*UpdateGroupMemberRoleResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return nil, errors.New("groupID must not be empty")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}
	role := strings.TrimSpace(request.Role)
	if role == "" {
		return nil, errors.New("role must not be empty")
	}

	body, err := jsonBody(struct {
		Role string `json:"role,omitempty"`
	}{Role: role})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/groups/%s/members/%s", url.PathEscape(groupID), url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdateGroupMemberRoleResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BeginGroupAvatarUpload starts a server-controlled avatar upload session for a group.
func (c *Client) BeginGroupAvatarUpload(ctx context.Context, request *BeginGroupAvatarUploadRequest, opts ...RequestOpt) (*BeginAvatarUploadResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return nil, errors.New("groupID must not be empty")
	}
	if strings.TrimSpace(request.ContentType) == "" {
		return nil, errors.New("contentType must not be empty")
	}

	body, err := jsonBody(struct {
		ContentType   string `json:"contentType,omitempty"`
		ContentLength *int64 `json:"contentLength,omitempty,string"`
	}{
		ContentType:   strings.TrimSpace(request.ContentType),
		ContentLength: request.ContentLength,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/groups/%s/avatar:beginUpload", url.PathEscape(groupID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BeginAvatarUploadResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CompleteGroupAvatarUpload completes an avatar upload session for a group.
func (c *Client) CompleteGroupAvatarUpload(ctx context.Context, request *CompleteGroupAvatarUploadRequest, opts ...RequestOpt) (*CompleteGroupAvatarUploadResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	groupID := strings.TrimSpace(request.GroupID)
	if groupID == "" {
		return nil, errors.New("groupID must not be empty")
	}
	if strings.TrimSpace(request.UploadID) == "" {
		return nil, errors.New("uploadID must not be empty")
	}

	body, err := jsonBody(struct {
		UploadID string `json:"uploadId,omitempty"`
	}{UploadID: strings.TrimSpace(request.UploadID)})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/groups/%s/avatar:completeUpload", url.PathEscape(groupID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CompleteGroupAvatarUploadResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
