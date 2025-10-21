package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListResourceTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/api/v1/tags/resource-1", r.URL.Path)
		require.Equal(t, string(ResourceTypeDeck), r.URL.Query().Get("resource.type"))

		require.NoError(t, json.NewEncoder(w).Encode(ListResourceTagsResponse{
			Tags: []Tag{{Key: "format", Value: "cc"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ListResourceTags(context.Background(), &ListResourceTagsRequest{
		ResourceID:   "resource-1",
		ResourceType: ResourceTypeDeck,
	})
	require.NoError(t, err)
	require.Len(t, resp.Tags, 1)

	resp, err = client.ListResourceTags(context.Background(), &ListResourceTagsRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientTagResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v1/tags/resource-1", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload struct {
			Resource struct {
				Type ResourceType `json:"type"`
			} `json:"resource"`
			Tag Tag `json:"tag"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, ResourceTypeCollection, payload.Resource.Type)
		require.Equal(t, "category", payload.Tag.Key)
		require.Equal(t, "guardian", payload.Tag.Value)

		w.Header().Set("X-Request-Id", "req-tag")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.TagResource(context.Background(), &TagResourceRequest{
		ResourceID:   "resource-1",
		ResourceType: ResourceTypeCollection,
		Tag: Tag{
			Key:   "category",
			Value: "guardian",
		},
	})
	require.NoError(t, err)
	require.Equal(t, "req-tag", resp.Metadata.RequestID)
}

func TestClientUntagResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v1/tags/resource-1:untag", r.URL.Path)

		var payload struct {
			Resource struct {
				Type ResourceType `json:"type"`
			} `json:"resource"`
			Tag Tag `json:"tag"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, ResourceTypeDeck, payload.Resource.Type)
		require.Equal(t, "format", payload.Tag.Key)

		w.Header().Set("X-Request-Id", "req-untag")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.UntagResource(context.Background(), &UntagResourceRequest{
		ResourceID:   "resource-1",
		ResourceType: ResourceTypeDeck,
		Tag: Tag{
			Key: "format",
		},
	})
	require.NoError(t, err)
	require.Equal(t, "req-untag", resp.Metadata.RequestID)
}

func TestClientUntagAllForResource(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v1/tags/resource-1:untagAll", r.URL.Path)

		var payload struct {
			Type ResourceType `json:"type"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, ResourceTypeDeck, payload.Type)

		w.Header().Set("X-Request-Id", "req-untag-all")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.UntagAllForResource(context.Background(), &UntagAllForResourceRequest{
		ResourceID:   "resource-1",
		ResourceType: ResourceTypeDeck,
	})
	require.NoError(t, err)
	require.Equal(t, "req-untag-all", resp.Metadata.RequestID)
}

func TestClientBatchListResourceTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v1/tags:batchList", r.URL.Path)

		var payload struct {
			Resources []Resource `json:"resources"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Len(t, payload.Resources, 1)
		require.Equal(t, "resource-1", payload.Resources[0].ID)
		require.Equal(t, ResourceTypeDeck, payload.Resources[0].Type)

		require.NoError(t, json.NewEncoder(w).Encode(BatchListResourceTagsResponse{
			Results: []ResourceTags{
				{
					Resource: Resource{ID: "resource-1"},
					Tags:     []Tag{{Key: "format", Value: "cc"}},
				},
			},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.BatchListResourceTags(context.Background(), &BatchListResourceTagsRequest{
		Resources: []Resource{
			{ID: "resource-1", Type: ResourceTypeDeck},
		},
	})
	require.NoError(t, err)
	require.Len(t, resp.Results, 1)
}

func TestClientBatchTagResources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v1/tags:batchTag", r.URL.Path)

		var payload struct {
			Items []struct {
				Resource Resource `json:"resource"`
				Tag      Tag      `json:"tag"`
			} `json:"items"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Len(t, payload.Items, 1)
		require.Equal(t, "resource-1", payload.Items[0].Resource.ID)
		require.Equal(t, ResourceTypeDeck, payload.Items[0].Resource.Type)
		require.Equal(t, "format", payload.Items[0].Tag.Key)

		w.Header().Set("X-Request-Id", "req-batch-tag")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.BatchTagResources(context.Background(), &BatchTagResourcesRequest{
		Items: []TagResourceRequest{
			{
				ResourceID:   "resource-1",
				ResourceType: ResourceTypeDeck,
				Tag: Tag{
					Key:   "format",
					Value: "cc",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "req-batch-tag", resp.Metadata.RequestID)
}

func TestClientBatchUntagResources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v1/tags:batchUntag", r.URL.Path)

		var payload struct {
			Items []struct {
				Resource Resource `json:"resource"`
				Tag      Tag      `json:"tag"`
			} `json:"items"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Len(t, payload.Items, 1)
		require.Equal(t, ResourceTypeDeck, payload.Items[0].Resource.Type)
		require.Equal(t, "format", payload.Items[0].Tag.Key)

		w.Header().Set("X-Request-Id", "req-batch-untag")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.BatchUntagResources(context.Background(), &BatchUntagResourcesRequest{
		Items: []UntagResourceRequest{
			{
				ResourceID:   "resource-1",
				ResourceType: ResourceTypeDeck,
				Tag: Tag{
					Key: "format",
				},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "req-batch-untag", resp.Metadata.RequestID)
}

func TestClientQueryResourcesByTags(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v1/tags:queryResources", r.URL.Path)

		var payload struct {
			Filters      []TagFilter  `json:"filters"`
			ResourceType ResourceType `json:"resourceType"`
			PageSize     int32        `json:"pageSize"`
			NextToken    string       `json:"nextToken"`
			Mode         MatchMode    `json:"mode"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Len(t, payload.Filters, 1)
		require.Equal(t, "format", payload.Filters[0].Key)
		require.Equal(t, MatchOperatorExact, payload.Filters[0].Operator)
		require.Equal(t, ResourceTypeDeck, payload.ResourceType)
		require.Equal(t, int32(10), payload.PageSize)
		require.Equal(t, "token", payload.NextToken)
		require.Equal(t, MatchModeAny, payload.Mode)

		require.NoError(t, json.NewEncoder(w).Encode(QueryResourcesByTagsResponse{
			Resources: []Resource{{ID: "resource-1"}},
			NextToken: "more",
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	pageSize := int32(10)
	resp, err := client.QueryResourcesByTags(context.Background(), &QueryResourcesByTagsRequest{
		Filters: []TagFilter{
			{Key: "format", Values: []string{"cc"}, Operator: MatchOperatorExact},
		},
		ResourceType: ResourceTypeDeck,
		PageSize:     &pageSize,
		NextToken:    "token",
		Mode:         MatchModeAny,
	})
	require.NoError(t, err)
	require.Equal(t, "more", resp.NextToken)
	require.Len(t, resp.Resources, 1)
}
