package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
)

// GroupSummary is the minimal group detail needed to render a row that refers
// to a group.
type GroupSummary struct {
	GroupID   string `json:"groupId,omitempty"`
	Slug      string `json:"slug,omitempty"`
	Name      string `json:"name,omitempty"`
	AvatarURL string `json:"avatarUrl,omitempty"`
}

// MyGroupInvite is a pending invite addressed to the caller.
type MyGroupInvite struct {
	Invite *GroupInvite  `json:"invite,omitempty"`
	Group  *GroupSummary `json:"group,omitempty"`
}

// MyGroupJoinRequest is one of the caller's own join requests.
type MyGroupJoinRequest struct {
	JoinRequest *GroupJoinRequest `json:"joinRequest,omitempty"`
	Group       *GroupSummary     `json:"group,omitempty"`
}

// ListMyGroupInvitesRequest paginates ListMyGroupInvites.
type ListMyGroupInvitesRequest struct {
	PageSize  *int32
	NextToken string
}

// ListMyGroupInvitesResponse lists the caller's pending invites.
type ListMyGroupInvitesResponse struct {
	Invites   []MyGroupInvite  `json:"invites,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListMyGroupInvitesResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

// ListMyGroupJoinRequestsRequest paginates ListMyGroupJoinRequests.
type ListMyGroupJoinRequestsRequest struct {
	PageSize  *int32
	NextToken string
}

// ListMyGroupJoinRequestsResponse lists the caller's join requests.
type ListMyGroupJoinRequestsResponse struct {
	JoinRequests []MyGroupJoinRequest `json:"joinRequests,omitempty"`
	NextToken    string               `json:"nextToken,omitempty"`
	Metadata     ResponseMetadata     `json:"-"`
}

func (r *ListMyGroupJoinRequestsResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

// MyGroupInviteRequest identifies one of the caller's invites by id.
type MyGroupInviteRequest struct {
	GroupID  string
	InviteID string
}

// ListMyGroupInvites lists the caller's pending targeted invites.
func (c *Client) ListMyGroupInvites(ctx context.Context, request *ListMyGroupInvitesRequest, opts ...RequestOpt) (*ListMyGroupInvitesResponse, error) {
	if request == nil {
		request = &ListMyGroupInvitesRequest{}
	}
	response := &ListMyGroupInvitesResponse{}
	path := addPageQuery("/v1/groups:myInvites", request.PageSize, request.NextToken)
	if err := c.doGroupJSON(ctx, http.MethodGet, path, nil, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListMyGroupJoinRequests lists the caller's pending join requests and those
// declined recently.
func (c *Client) ListMyGroupJoinRequests(ctx context.Context, request *ListMyGroupJoinRequestsRequest, opts ...RequestOpt) (*ListMyGroupJoinRequestsResponse, error) {
	if request == nil {
		request = &ListMyGroupJoinRequestsRequest{}
	}
	response := &ListMyGroupJoinRequestsResponse{}
	path := addPageQuery("/v1/groups:myJoinRequests", request.PageSize, request.NextToken)
	if err := c.doGroupJSON(ctx, http.MethodGet, path, nil, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// AcceptMyGroupInvite accepts one of the caller's targeted invites by id.
func (c *Client) AcceptMyGroupInvite(ctx context.Context, request *MyGroupInviteRequest, opts ...RequestOpt) (*GroupActionResponse, error) {
	group, invite, err := myGroupInviteIDs(request)
	if err != nil {
		return nil, err
	}
	response := &GroupActionResponse{}
	if err := c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/invites/"+url.PathEscape(invite)+":accept", map[string]any{}, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// DeclineMyGroupInvite declines one of the caller's targeted invites by id.
func (c *Client) DeclineMyGroupInvite(ctx context.Context, request *MyGroupInviteRequest, opts ...RequestOpt) (*GroupActionResponse, error) {
	group, invite, err := myGroupInviteIDs(request)
	if err != nil {
		return nil, err
	}
	response := &GroupActionResponse{}
	if err := c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/invites/"+url.PathEscape(invite)+":decline", map[string]any{}, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

func myGroupInviteIDs(request *MyGroupInviteRequest) (string, string, error) {
	if request == nil {
		return "", "", errors.New("request must not be nil")
	}
	group, err := required(request.GroupID, "groupID")
	if err != nil {
		return "", "", err
	}
	invite, err := required(request.InviteID, "inviteID")
	if err != nil {
		return "", "", err
	}
	return group, invite, nil
}
