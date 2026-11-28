package client

import (
	"context"
	"errors"
	"net/http"
	"strings"
)

// PatreonLinkSummary describes the Patreon account linked to the caller.
type PatreonLinkSummary struct {
	PatreonUserID   string   `json:"patreonUserId,omitempty"`
	PatreonMemberID string   `json:"patreonMemberId,omitempty"`
	CampaignID      string   `json:"campaignId,omitempty"`
	Email           string   `json:"email,omitempty"`
	FullName        string   `json:"fullName,omitempty"`
	TierIDs         []string `json:"tierIds,omitempty"`
}

// InitiatePatreonLinkRequest starts linking a Patreon account.
type InitiatePatreonLinkRequest struct{}

// InitiatePatreonLinkResponse returns the Patreon authorization URL.
type InitiatePatreonLinkResponse struct {
	AuthorizationURL string           `json:"authorizationUrl,omitempty"`
	State            string           `json:"state,omitempty"`
	Metadata         ResponseMetadata `json:"-"`
}

func (r *InitiatePatreonLinkResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CompletePatreonLinkRequest completes linking a Patreon account.
type CompletePatreonLinkRequest struct {
	Code  string `json:"code,omitempty"`
	State string `json:"state,omitempty"`
}

// CompletePatreonLinkResponse returns the linked account and refreshed access profile.
type CompletePatreonLinkResponse struct {
	Link          *PatreonLinkSummary `json:"link,omitempty"`
	AccessProfile *AccessProfile      `json:"accessProfile,omitempty"`
	Metadata      ResponseMetadata    `json:"-"`
}

func (r *CompletePatreonLinkResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnlinkPatreonRequest unlinks the caller's Patreon account.
type UnlinkPatreonRequest struct{}

// UnlinkPatreonResponse captures metadata for unlink operations.
type UnlinkPatreonResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnlinkPatreonResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// InitiatePatreonLink starts linking a Patreon account to the authenticated user.
func (c *Client) InitiatePatreonLink(ctx context.Context, request *InitiatePatreonLinkRequest, opts ...RequestOpt) (*InitiatePatreonLinkResponse, error) {
	if request == nil {
		request = &InitiatePatreonLinkRequest{}
	}

	response := &InitiatePatreonLinkResponse{}
	if err := c.postAuthJSON(ctx, "/v1/auth/patreon/link:initiate", request, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CompletePatreonLink completes a Patreon link with the provider callback code.
func (c *Client) CompletePatreonLink(ctx context.Context, request *CompletePatreonLinkRequest, opts ...RequestOpt) (*CompletePatreonLinkResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.Code) == "" {
		return nil, errors.New("code must not be empty")
	}

	payload := CompletePatreonLinkRequest{
		Code:  strings.TrimSpace(request.Code),
		State: strings.TrimSpace(request.State),
	}
	response := &CompletePatreonLinkResponse{}
	if err := c.postAuthJSON(ctx, "/v1/auth/patreon/link:complete", payload, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UnlinkPatreon removes the authenticated user's Patreon link.
func (c *Client) UnlinkPatreon(ctx context.Context, request *UnlinkPatreonRequest, opts ...RequestOpt) (*UnlinkPatreonResponse, error) {
	req, err := c.newRequest(ctx, http.MethodDelete, "/v1/auth/patreon/link", nil)
	if err != nil {
		return nil, err
	}

	response := &UnlinkPatreonResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
