package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ProfileCollectionItemSearchSide identifies which face of a card matched.
type ProfileCollectionItemSearchSide string

const (
	ProfileCollectionItemSearchSideUnspecified ProfileCollectionItemSearchSide = "PROFILE_COLLECTION_ITEM_SEARCH_SIDE_UNSPECIFIED"
	ProfileCollectionItemSearchSideFront       ProfileCollectionItemSearchSide = "PROFILE_COLLECTION_ITEM_SEARCH_SIDE_FRONT"
	ProfileCollectionItemSearchSideBack        ProfileCollectionItemSearchSide = "PROFILE_COLLECTION_ITEM_SEARCH_SIDE_BACK"
)

// SearchProfileCollectionItemsRequest searches the collection items a user's
// profile exposes to the caller.
type SearchProfileCollectionItemsRequest struct {
	UserID          string           `json:"-"`
	SearchTerm      string           `json:"searchTerm,omitempty"`
	CardIDs         []string         `json:"cardIds,omitempty"`
	PrintingIDs     []string         `json:"printingIds,omitempty"`
	ProductIDs      []string         `json:"productIds,omitempty"`
	CollectionTypes []CollectionType `json:"collectionTypes,omitempty"`
	PageSize        *int32           `json:"pageSize,omitempty"`
	NextToken       string           `json:"nextToken,omitempty"`
}

// ProfileCollectionItemSearchSource identifies the collection an item belongs to.
type ProfileCollectionItemSearchSource struct {
	ID             string         `json:"id,omitempty"`
	Name           string         `json:"name,omitempty"`
	CollectionType CollectionType `json:"collectionType,omitempty"`
}

// ProfileCollectionItemSearchCatalogMatch describes the catalog entry an item matched.
type ProfileCollectionItemSearchCatalogMatch struct {
	CardID     string                          `json:"cardId,omitempty"`
	PrintingID string                          `json:"printingId,omitempty"`
	Side       ProfileCollectionItemSearchSide `json:"side,omitempty"`
}

// ProfileCollectionItemSearchResult is one matching collection item.
type ProfileCollectionItemSearchResult struct {
	ItemID           string                                    `json:"itemId,omitempty"`
	ProductID        string                                    `json:"productId,omitempty"`
	Quantity         int32                                     `json:"quantity,omitempty"`
	Condition        Condition                                 `json:"condition,omitempty"`
	SourceCollection *ProfileCollectionItemSearchSource        `json:"sourceCollection,omitempty"`
	CatalogMatches   []ProfileCollectionItemSearchCatalogMatch `json:"catalogMatches,omitempty"`
}

// SearchProfileCollectionItemsResponse lists matching collection items.
type SearchProfileCollectionItemsResponse struct {
	Items     []ProfileCollectionItemSearchResult `json:"items,omitempty"`
	NextToken string                              `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata                    `json:"-"`
}

func (r *SearchProfileCollectionItemsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// SearchProfileCollectionItems searches the collection items visible on a
// user's profile.
func (c *Client) SearchProfileCollectionItems(ctx context.Context, request *SearchProfileCollectionItemsRequest, opts ...RequestOpt) (*SearchProfileCollectionItemsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	payload := *request
	payload.SearchTerm = strings.TrimSpace(payload.SearchTerm)
	payload.NextToken = strings.TrimSpace(payload.NextToken)
	if payload.PageSize != nil && *payload.PageSize <= 0 {
		payload.PageSize = nil
	}
	body, err := jsonBody(payload)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/users/%s/collection_items:search", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &SearchProfileCollectionItemsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
