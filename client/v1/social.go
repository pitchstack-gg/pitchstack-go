package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// FollowUserRequest initiates a follow action for the target user.
type FollowUserRequest struct {
	TargetUserID string
}

// FollowUserResponse captures metadata from a follow request.
type FollowUserResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *FollowUserResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnfollowUserRequest removes a follow relationship.
type UnfollowUserRequest struct {
	TargetUserID string
}

// UnfollowUserResponse captures metadata from an unfollow request.
type UnfollowUserResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnfollowUserResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListFollowersRequest fetches user IDs following a given user.
type ListFollowersRequest struct {
	UserID    string
	PageSize  *int32
	NextToken string
}

// ListFollowersResponse returns follower user IDs.
type ListFollowersResponse struct {
	UserIDs   []string         `json:"userIds,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListFollowersResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListFollowingRequest fetches user IDs that a user is following.
type ListFollowingRequest struct {
	UserID    string
	PageSize  *int32
	NextToken string
}

// ListFollowingResponse returns user IDs being followed.
type ListFollowingResponse struct {
	UserIDs   []string         `json:"userIds,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListFollowingResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetFollowStatsRequest fetches follower/following counts for a user.
type GetFollowStatsRequest struct {
	UserID string
}

// GetFollowStatsResponse returns follower/following totals.
type GetFollowStatsResponse struct {
	Followers int64            `json:"followers,omitempty"`
	Following int64            `json:"following,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *GetFollowStatsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// IsFollowingRequest checks whether followerID follows followeeID.
type IsFollowingRequest struct {
	FollowerID string
	FolloweeID string
}

// IsFollowingResponse indicates if the follower relationship exists.
type IsFollowingResponse struct {
	IsFollowing bool             `json:"isFollowing,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *IsFollowingResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// FollowUser begins following the specified target user.
func (c *Client) FollowUser(ctx context.Context, request *FollowUserRequest, opts ...RequestOpt) (*FollowUserResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	targetUserID := strings.TrimSpace(request.TargetUserID)
	if targetUserID == "" {
		return nil, errors.New("targetUserID must not be empty")
	}

	path := fmt.Sprintf("/v1/users/%s/followers", url.PathEscape(targetUserID))
	req, err := c.newRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	response := &FollowUserResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UnfollowUser stops following the specified target user.
func (c *Client) UnfollowUser(ctx context.Context, request *UnfollowUserRequest, opts ...RequestOpt) (*UnfollowUserResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	targetUserID := strings.TrimSpace(request.TargetUserID)
	if targetUserID == "" {
		return nil, errors.New("targetUserID must not be empty")
	}

	path := fmt.Sprintf("/v1/users/%s/followers", url.PathEscape(targetUserID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &UnfollowUserResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListFollowers retrieves user IDs following the specified user.
func (c *Client) ListFollowers(ctx context.Context, request *ListFollowersRequest, opts ...RequestOpt) (*ListFollowersResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	query := url.Values{}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}

	path := fmt.Sprintf("/v1/users/%s/followers", url.PathEscape(userID))
	if len(query) > 0 {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &ListFollowersResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListFollowing retrieves user IDs the specified user follows.
func (c *Client) ListFollowing(ctx context.Context, request *ListFollowingRequest, opts ...RequestOpt) (*ListFollowingResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	query := url.Values{}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}

	path := fmt.Sprintf("/v1/users/%s/following", url.PathEscape(userID))
	if len(query) > 0 {
		path = fmt.Sprintf("%s?%s", path, query.Encode())
	}

	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &ListFollowingResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetFollowStats retrieves follower/following counts for a user.
func (c *Client) GetFollowStats(ctx context.Context, request *GetFollowStatsRequest, opts ...RequestOpt) (*GetFollowStatsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/v1/users/%s/follow_stats", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetFollowStatsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// IsFollowing determines if followerID is following followeeID.
func (c *Client) IsFollowing(ctx context.Context, request *IsFollowingRequest, opts ...RequestOpt) (*IsFollowingResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	followerID := strings.TrimSpace(request.FollowerID)
	if followerID == "" {
		return nil, errors.New("followerID must not be empty")
	}
	followeeID := strings.TrimSpace(request.FolloweeID)
	if followeeID == "" {
		return nil, errors.New("followeeID must not be empty")
	}

	path := fmt.Sprintf("/v1/users/%s/following/%s", url.PathEscape(followerID), url.PathEscape(followeeID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &IsFollowingResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
