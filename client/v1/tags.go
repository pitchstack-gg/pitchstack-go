package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// ResourceType represents v1ResourceType.
type ResourceType string

const (
	ResourceTypeUnspecified    ResourceType = "RESOURCE_TYPE_UNSPECIFIED"
	ResourceTypeCollection     ResourceType = "RESOURCE_TYPE_COLLECTION"
	ResourceTypeDeck           ResourceType = "RESOURCE_TYPE_DECK"
	ResourceTypeUserProfile    ResourceType = "RESOURCE_TYPE_USER_PROFILE"
	ResourceTypeAdminPanel     ResourceType = "RESOURCE_TYPE_ADMIN_PANEL"
	ResourceTypeCollectionItem ResourceType = "RESOURCE_TYPE_COLLECTION_ITEM"
)

// MatchMode represents v1MatchMode.
type MatchMode string

const (
	MatchModeUnspecified MatchMode = "MATCH_MODE_UNSPECIFIED"
	MatchModeAny         MatchMode = "MATCH_MODE_ANY"
	MatchModeAll         MatchMode = "MATCH_MODE_ALL"
)

// MatchOperator represents v1MatchOperator.
type MatchOperator string

const (
	MatchOperatorUnspecified MatchOperator = "MATCH_OPERATOR_UNSPECIFIED"
	MatchOperatorExact       MatchOperator = "MATCH_OPERATOR_EXACT"
	MatchOperatorPrefix      MatchOperator = "MATCH_OPERATOR_PREFIX"
)

// Resource identifies a tagged entity.
type Resource struct {
	ID   string       `json:"id,omitempty"`
	Type ResourceType `json:"type,omitempty"`
}

// TagKeyValue represents a key/value tag specification.
type TagKeyValue struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

// Tag represents a label or key/value tag on a resource.
type Tag struct {
	Label    string       `json:"label,omitempty"`
	KeyValue *TagKeyValue `json:"keyValue,omitempty"`
}

// TagFilter filters resources by tag values.
type TagFilter struct {
	Key      string        `json:"key,omitempty"`
	Values   []string      `json:"values,omitempty"`
	Operator MatchOperator `json:"operator,omitempty"`
}

// ResourceTags associates a resource with its tags.
type ResourceTags struct {
	Resource  Resource `json:"resource,omitempty"`
	Tags      []Tag    `json:"tags,omitempty"`
	NextToken string   `json:"nextToken,omitempty"`
}

// ListResourceTagsRequest fetches tags for a resource.
type ListResourceTagsRequest struct {
	ResourceID   string
	ResourceType ResourceType
	NextToken    string
}

// ListResourceTagsResponse returns tags assigned to a resource.
type ListResourceTagsResponse struct {
	Tags      []Tag            `json:"tags,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListResourceTagsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// TagResourceRequest assigns a tag to a resource.
type TagResourceRequest struct {
	ResourceID   string
	ResourceType ResourceType
	Tag          Tag
}

// TagResourceResponse captures metadata for tag operations.
type TagResourceResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *TagResourceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UntagResourceRequest removes a tag from a resource.
type UntagResourceRequest struct {
	ResourceID   string
	ResourceType ResourceType
	Tag          Tag
}

// UntagResourceResponse captures metadata for untag operations.
type UntagResourceResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UntagResourceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UntagAllForResourceRequest removes all tags from a resource.
type UntagAllForResourceRequest struct {
	ResourceID   string
	ResourceType ResourceType
}

// UntagAllForResourceResponse captures metadata for bulk untag operations.
type UntagAllForResourceResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UntagAllForResourceResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchListResourceTagsRequest fetches tags for multiple resources.
type BatchListResourceTagsRequest struct {
	Resources []Resource
}

// BatchListResourceTagsResponse returns tags for requested resources.
type BatchListResourceTagsResponse struct {
	Results  []ResourceTags   `json:"results,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *BatchListResourceTagsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchTagResourcesRequest applies tags to multiple resources.
type BatchTagResourcesRequest struct {
	Items []TagResourceRequest
}

// BatchTagResourcesResponse captures metadata for batch tag operations.
type BatchTagResourcesResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *BatchTagResourcesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchUntagResourcesRequest removes tags from multiple resources.
type BatchUntagResourcesRequest struct {
	Items []UntagResourceRequest
}

// BatchUntagResourcesResponse captures metadata for batch untag operations.
type BatchUntagResourcesResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *BatchUntagResourcesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// QueryResourcesByTagsRequest searches resources by tags.
type QueryResourcesByTagsRequest struct {
	Filters      []TagFilter
	ResourceType ResourceType
	PageSize     *int32
	NextToken    string
	Mode         MatchMode
}

// QueryResourcesByTagsResponse contains resources matching the filters.
type QueryResourcesByTagsResponse struct {
	Resources []Resource       `json:"resources,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *QueryResourcesByTagsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListResourceTags retrieves tags for a specific resource.
func (c *Client) ListResourceTags(ctx context.Context, request *ListResourceTagsRequest, opts ...RequestOpt) (*ListResourceTagsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	resourceID := strings.TrimSpace(request.ResourceID)
	if resourceID == "" {
		return nil, errors.New("resourceID must not be empty")
	}

	path := fmt.Sprintf("/v1/tags/%s", url.PathEscape(resourceID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if t := strings.TrimSpace(string(request.ResourceType)); t != "" && t != string(ResourceTypeUnspecified) {
		query.Set("resource.type", t)
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	response := &ListResourceTagsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// TagResource assigns a tag to a resource.
func (c *Client) TagResource(ctx context.Context, request *TagResourceRequest, opts ...RequestOpt) (*TagResourceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	resourceID := strings.TrimSpace(request.ResourceID)
	if resourceID == "" {
		return nil, errors.New("resourceID must not be empty")
	}
	if err := validateResourceType(request.ResourceType); err != nil {
		return nil, err
	}
	if err := validateTag(request.Tag); err != nil {
		return nil, err
	}

	body, err := jsonBody(struct {
		Resource struct {
			Type ResourceType `json:"type,omitempty"`
		} `json:"resource,omitempty"`
		Tag Tag `json:"tag,omitempty"`
	}{
		Resource: struct {
			Type ResourceType `json:"type,omitempty"`
		}{Type: request.ResourceType},
		Tag: request.Tag,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/tags/%s", url.PathEscape(resourceID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &TagResourceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UntagResource removes a tag from a resource.
func (c *Client) UntagResource(ctx context.Context, request *UntagResourceRequest, opts ...RequestOpt) (*UntagResourceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	resourceID := strings.TrimSpace(request.ResourceID)
	if resourceID == "" {
		return nil, errors.New("resourceID must not be empty")
	}
	if err := validateResourceType(request.ResourceType); err != nil {
		return nil, err
	}
	if err := validateTag(request.Tag); err != nil {
		return nil, err
	}

	body, err := jsonBody(struct {
		Resource struct {
			Type ResourceType `json:"type,omitempty"`
		} `json:"resource,omitempty"`
		Tag Tag `json:"tag,omitempty"`
	}{
		Resource: struct {
			Type ResourceType `json:"type,omitempty"`
		}{Type: request.ResourceType},
		Tag: request.Tag,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/tags/%s:untag", url.PathEscape(resourceID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UntagResourceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UntagAllForResource removes all tags from a resource.
func (c *Client) UntagAllForResource(ctx context.Context, request *UntagAllForResourceRequest, opts ...RequestOpt) (*UntagAllForResourceResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	resourceID := strings.TrimSpace(request.ResourceID)
	if resourceID == "" {
		return nil, errors.New("resourceID must not be empty")
	}
	if err := validateResourceType(request.ResourceType); err != nil {
		return nil, err
	}

	body, err := jsonBody(struct {
		Type ResourceType `json:"type,omitempty"`
	}{
		Type: request.ResourceType,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/tags/%s:untagAll", url.PathEscape(resourceID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UntagAllForResourceResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BatchListResourceTags retrieves tags for multiple resources.
func (c *Client) BatchListResourceTags(ctx context.Context, request *BatchListResourceTagsRequest, opts ...RequestOpt) (*BatchListResourceTagsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.Resources) == 0 {
		return nil, errors.New("resources must not be empty")
	}
	for _, res := range request.Resources {
		if strings.TrimSpace(res.ID) == "" {
			return nil, errors.New("resource id must not be empty")
		}
		if err := validateResourceType(res.Type); err != nil {
			return nil, err
		}
	}

	body, err := jsonBody(struct {
		Resources []Resource `json:"resources,omitempty"`
	}{
		Resources: request.Resources,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/tags:batchList", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BatchListResourceTagsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BatchTagResources applies tags to multiple resources.
func (c *Client) BatchTagResources(ctx context.Context, request *BatchTagResourcesRequest, opts ...RequestOpt) (*BatchTagResourcesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.Items) == 0 {
		return nil, errors.New("items must not be empty")
	}

	type bodyItem struct {
		Resource Resource `json:"resource,omitempty"`
		Tag      Tag      `json:"tag,omitempty"`
	}
	payload := struct {
		Items []bodyItem `json:"items,omitempty"`
	}{Items: make([]bodyItem, len(request.Items))}

	for i, item := range request.Items {
		if strings.TrimSpace(item.ResourceID) == "" {
			return nil, errors.New("resourceID must not be empty")
		}
		if err := validateResourceType(item.ResourceType); err != nil {
			return nil, err
		}
		if err := validateTag(item.Tag); err != nil {
			return nil, err
		}
		payload.Items[i] = bodyItem{
			Resource: Resource{
				ID:   item.ResourceID,
				Type: item.ResourceType,
			},
			Tag: item.Tag,
		}
	}

	body, err := jsonBody(payload)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/tags:batchTag", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BatchTagResourcesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// BatchUntagResources removes tags from multiple resources.
func (c *Client) BatchUntagResources(ctx context.Context, request *BatchUntagResourcesRequest, opts ...RequestOpt) (*BatchUntagResourcesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.Items) == 0 {
		return nil, errors.New("items must not be empty")
	}

	type bodyItem struct {
		Resource Resource `json:"resource,omitempty"`
		Tag      Tag      `json:"tag,omitempty"`
	}
	payload := struct {
		Items []bodyItem `json:"items,omitempty"`
	}{Items: make([]bodyItem, len(request.Items))}

	for i, item := range request.Items {
		if strings.TrimSpace(item.ResourceID) == "" {
			return nil, errors.New("resourceID must not be empty")
		}
		if err := validateResourceType(item.ResourceType); err != nil {
			return nil, err
		}
		if err := validateTag(item.Tag); err != nil {
			return nil, err
		}
		payload.Items[i] = bodyItem{
			Resource: Resource{
				ID:   item.ResourceID,
				Type: item.ResourceType,
			},
			Tag: item.Tag,
		}
	}

	body, err := jsonBody(payload)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/tags:batchUntag", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &BatchUntagResourcesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// QueryResourcesByTags searches for resources matching tags.
func (c *Client) QueryResourcesByTags(ctx context.Context, request *QueryResourcesByTagsRequest, opts ...RequestOpt) (*QueryResourcesByTagsResponse, error) {
	if request == nil {
		request = &QueryResourcesByTagsRequest{}
	}

	type bodyFilter struct {
		Key      string        `json:"key,omitempty"`
		Values   []string      `json:"values,omitempty"`
		Operator MatchOperator `json:"operator,omitempty"`
	}

	payload := struct {
		Filters      []bodyFilter `json:"filters,omitempty"`
		ResourceType ResourceType `json:"resourceType,omitempty"`
		PageSize     *int32       `json:"pageSize,omitempty"`
		NextToken    string       `json:"nextToken,omitempty"`
		Mode         MatchMode    `json:"mode,omitempty"`
	}{}

	if len(request.Filters) > 0 {
		payload.Filters = make([]bodyFilter, len(request.Filters))
		for i, f := range request.Filters {
			payload.Filters[i] = bodyFilter(f)
		}
	}
	if t := strings.TrimSpace(string(request.ResourceType)); t != "" && t != string(ResourceTypeUnspecified) {
		payload.ResourceType = request.ResourceType
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		payload.PageSize = request.PageSize
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		payload.NextToken = token
	}
	if m := strings.TrimSpace(string(request.Mode)); m != "" && m != string(MatchModeUnspecified) {
		payload.Mode = request.Mode
	}

	body, err := jsonBody(payload)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/tags:queryResources", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &QueryResourcesByTagsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

func validateResourceType(resourceType ResourceType) error {
	t := strings.TrimSpace(string(resourceType))
	if t == "" || t == string(ResourceTypeUnspecified) {
		return errors.New("resourceType must be specified")
	}
	return nil
}

func validateTag(tag Tag) error {
	hasLabel := strings.TrimSpace(tag.Label) != ""
	hasKeyValue := tag.KeyValue != nil

	if !hasLabel && !hasKeyValue {
		return errors.New("tag must include label or keyValue")
	}
	if hasLabel && hasKeyValue {
		return errors.New("tag must not include both label and keyValue")
	}
	if hasKeyValue {
		if strings.TrimSpace(tag.KeyValue.Key) == "" {
			return errors.New("tag key must not be empty")
		}
	}
	return nil
}
