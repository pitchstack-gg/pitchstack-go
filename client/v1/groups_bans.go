package client

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GroupBan records a user banned from a group.
type GroupBan struct {
	GroupID        string     `json:"groupId,omitempty"`
	UserID         string     `json:"userId,omitempty"`
	BannedByUserID string     `json:"bannedByUserId,omitempty"`
	Reason         string     `json:"reason,omitempty"`
	CreatedAt      *time.Time `json:"createdAt,omitempty"`
}

// BanGroupUserRequest bans a user from a group.
type BanGroupUserRequest struct {
	GroupID string `json:"-"`
	UserID  string `json:"userId,omitempty"`
	// Reason is an optional note visible to the group's admins only.
	Reason string `json:"reason,omitempty"`
}

// GroupBanResponse returns a group ban.
type GroupBanResponse struct {
	Ban      *GroupBan        `json:"ban,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GroupBanResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

// UnbanGroupUserRequest lifts a ban.
type UnbanGroupUserRequest struct {
	GroupID string
	UserID  string
}

// ListGroupBansRequest paginates a group's bans.
type ListGroupBansRequest struct {
	GroupID   string
	PageSize  *int32
	NextToken string
}

// ListGroupBansResponse lists a group's bans, newest first.
type ListGroupBansResponse struct {
	Bans      []GroupBan       `json:"bans,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListGroupBansResponse) setMetadata(m ResponseMetadata) { r.Metadata = m }

// BanGroupUser bans a user from a group, removing any membership they hold.
// Banning an already banned user returns the existing ban.
func (c *Client) BanGroupUser(ctx context.Context, request *BanGroupUserRequest, opts ...RequestOpt) (*GroupBanResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	group, err := required(request.GroupID, "groupID")
	if err != nil {
		return nil, err
	}
	user, err := required(request.UserID, "userID")
	if err != nil {
		return nil, err
	}
	body := struct {
		UserID string `json:"userId,omitempty"`
		Reason string `json:"reason,omitempty"`
	}{UserID: user, Reason: strings.TrimSpace(request.Reason)}
	response := &GroupBanResponse{}
	if err := c.doGroupJSON(ctx, http.MethodPost, "/v1/groups/"+url.PathEscape(group)+"/bans", body, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UnbanGroupUser lifts a user's ban. It does not restore a removed membership.
func (c *Client) UnbanGroupUser(ctx context.Context, request *UnbanGroupUserRequest, opts ...RequestOpt) (*GroupActionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	group, err := required(request.GroupID, "groupID")
	if err != nil {
		return nil, err
	}
	user, err := required(request.UserID, "userID")
	if err != nil {
		return nil, err
	}
	response := &GroupActionResponse{}
	if err := c.doGroupJSON(ctx, http.MethodDelete, "/v1/groups/"+url.PathEscape(group)+"/bans/"+url.PathEscape(user), nil, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListGroupBans lists a group's bans, newest first.
func (c *Client) ListGroupBans(ctx context.Context, request *ListGroupBansRequest, opts ...RequestOpt) (*ListGroupBansResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	group, err := required(request.GroupID, "groupID")
	if err != nil {
		return nil, err
	}
	response := &ListGroupBansResponse{}
	path := addPageQuery("/v1/groups/"+url.PathEscape(group)+"/bans", request.PageSize, request.NextToken)
	if err := c.doGroupJSON(ctx, http.MethodGet, path, nil, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
