package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientEngagement(t *testing.T) {
	t.Run("track view", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/engagement/views:track", r.URL.Path)
			var payload TrackViewRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, TrackableResourceTypeDeck, payload.Resource.ResourceType)
			require.Equal(t, "d-1", payload.Resource.ResourceID)
			require.NoError(t, json.NewEncoder(w).Encode(TrackViewResponse{Counted: true}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.TrackView(context.Background(), &TrackViewRequest{
			Resource: &EngagementResourceRef{ResourceType: TrackableResourceTypeDeck, ResourceID: "d-1"},
		})
		require.NoError(t, err)
		require.True(t, resp.Counted)

		_, err = client.TrackView(context.Background(), &TrackViewRequest{})
		require.Error(t, err)
	})

	t.Run("batch track views", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/engagement/views:batchTrack", r.URL.Path)
			var payload BatchTrackViewsRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Len(t, payload.Views, 1)
			require.NoError(t, json.NewEncoder(w).Encode(BatchTrackViewsResponse{TotalCounted: 1}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BatchTrackViews(context.Background(), &BatchTrackViewsRequest{
			Views: []TrackViewRequest{{Resource: &EngagementResourceRef{ResourceType: TrackableResourceTypeDeck, ResourceID: "d-1"}}},
		})
		require.NoError(t, err)
		require.Equal(t, int32(1), resp.TotalCounted)
	})

	t.Run("list trending resources", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/engagement/trending:list", r.URL.Path)
			var payload ListTrendingResourcesRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, TrackableResourceTypeDeck, payload.ResourceType)
			require.Equal(t, TrendingWindow24H, payload.Window)
			require.NoError(t, json.NewEncoder(w).Encode(ListTrendingResourcesResponse{
				Resources: []TrendingResource{{Resource: &EngagementResourceRef{ResourceType: TrackableResourceTypeDeck, ResourceID: "d-1"}}},
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListTrendingResources(context.Background(), &ListTrendingResourcesRequest{
			ResourceType: TrackableResourceTypeDeck,
			Window:       TrendingWindow24H,
			PageSize:     10,
		})
		require.NoError(t, err)
		require.Len(t, resp.Resources, 1)

		_, err = client.ListTrendingResources(context.Background(), &ListTrendingResourcesRequest{})
		require.Error(t, err)
	})

	t.Run("batch get view counts", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/engagement/views:batchGet", r.URL.Path)
			var payload BatchGetViewCountsRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Len(t, payload.Resources, 1)
			require.NoError(t, json.NewEncoder(w).Encode(BatchGetViewCountsResponse{
				Counts: []ResourceViewCount{{Resource: &EngagementResourceRef{ResourceType: TrackableResourceTypeDeck, ResourceID: "d-1"}, TotalViews: 3}},
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BatchGetViewCounts(context.Background(), &BatchGetViewCountsRequest{
			Resources: []EngagementResourceRef{{ResourceType: TrackableResourceTypeDeck, ResourceID: "d-1"}},
		})
		require.NoError(t, err)
		require.Len(t, resp.Counts, 1)

		_, err = client.BatchGetViewCounts(context.Background(), &BatchGetViewCountsRequest{})
		require.Error(t, err)
	})

	t.Run("batch get view counts with string encoded totals", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/engagement/views:batchGet", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"counts": []map[string]any{
					{
						"resource": map[string]any{
							"resourceType": string(TrackableResourceTypeDeck),
							"resourceId":   "d-1",
						},
						"totalViews": "3",
					},
				},
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BatchGetViewCounts(context.Background(), &BatchGetViewCountsRequest{
			Resources: []EngagementResourceRef{{ResourceType: TrackableResourceTypeDeck, ResourceID: "d-1"}},
		})
		require.NoError(t, err)
		require.Len(t, resp.Counts, 1)
		require.Equal(t, int64(3), resp.Counts[0].TotalViews)
	})

	t.Run("list trending resources with string encoded view count", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/engagement/trending:list", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
				"resources": []map[string]any{
					{
						"resource": map[string]any{
							"resourceType": string(TrackableResourceTypeDeck),
							"resourceId":   "d-1",
						},
						"viewCount": "7",
						"score":     1.5,
					},
				},
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListTrendingResources(context.Background(), &ListTrendingResourcesRequest{
			ResourceType: TrackableResourceTypeDeck,
			Window:       TrendingWindow24H,
			PageSize:     10,
		})
		require.NoError(t, err)
		require.Len(t, resp.Resources, 1)
		require.Equal(t, int64(7), resp.Resources[0].ViewCount)
	})

	t.Run("likes", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/v1/engagement/likes:like":
				require.Equal(t, http.MethodPost, r.Method)
				require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
					"resource":   map[string]any{"resourceType": string(LikeableResourceTypeDeck), "resourceId": "d-1"},
					"liked":      true,
					"totalLikes": "4",
					"changed":    true,
				}))
			case "/v1/engagement/likes:unlike":
				require.Equal(t, http.MethodPost, r.Method)
				require.NoError(t, json.NewEncoder(w).Encode(map[string]any{"liked": false, "totalLikes": "3", "changed": true}))
			case "/v1/engagement/likes:batchGetCounts":
				require.Equal(t, http.MethodPost, r.Method)
				require.NoError(t, json.NewEncoder(w).Encode(map[string]any{
					"counts": []map[string]any{{"resource": map[string]any{"resourceType": string(LikeableResourceTypeDeck), "resourceId": "d-1"}, "totalLikes": "4"}},
				}))
			case "/v1/engagement/likes:batchGetViewerLikes":
				require.Equal(t, http.MethodPost, r.Method)
				require.NoError(t, json.NewEncoder(w).Encode(BatchGetViewerLikesResponse{
					Likes: []ViewerLike{{Resource: &LikeableResourceRef{ResourceType: LikeableResourceTypeDeck, ResourceID: "d-1"}, Liked: true}},
				}))
			default:
				t.Fatalf("unexpected path: %s", r.URL.Path)
			}
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resource := LikeableResourceRef{ResourceType: LikeableResourceTypeDeck, ResourceID: "d-1"}
		likeResp, err := client.LikeResource(context.Background(), &LikeResourceRequest{Resource: &resource})
		require.NoError(t, err)
		require.Equal(t, int64(4), likeResp.TotalLikes)

		unlikeResp, err := client.UnlikeResource(context.Background(), &UnlikeResourceRequest{Resource: &resource})
		require.NoError(t, err)
		require.Equal(t, int64(3), unlikeResp.TotalLikes)

		countsResp, err := client.BatchGetLikeCounts(context.Background(), &BatchGetLikeCountsRequest{Resources: []LikeableResourceRef{resource}})
		require.NoError(t, err)
		require.Equal(t, int64(4), countsResp.Counts[0].TotalLikes)

		viewerResp, err := client.BatchGetViewerLikes(context.Background(), &BatchGetViewerLikesRequest{Resources: []LikeableResourceRef{resource}})
		require.NoError(t, err)
		require.True(t, viewerResp.Likes[0].Liked)
	})
}

func TestClientListLikedResources(t *testing.T) {
	t.Run("when filters provided, then POST body sent and likes decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/engagement/likes:list", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "LIKEABLE_RESOURCE_TYPE_DECK", payload["resourceType"])
			require.Equal(t, float64(50), payload["pageSize"])
			require.Equal(t, "tok", payload["nextPageToken"])
		}, `{"resources":[{"resource":{"resourceType":"LIKEABLE_RESOURCE_TYPE_DECK","resourceId":"deck-1"},"likedAt":"2026-11-01T00:00:00Z"}],"nextPageToken":"next"}`)

		resp, err := client.ListLikedResources(context.Background(), &ListLikedResourcesRequest{
			ResourceType:  LikeableResourceTypeDeck,
			PageSize:      50,
			NextPageToken: "tok",
		})
		require.NoError(t, err)
		require.Len(t, resp.Resources, 1)
		require.Equal(t, "deck-1", resp.Resources[0].Resource.ResourceID)
		require.NotNil(t, resp.Resources[0].LikedAt)
		require.Equal(t, "next", resp.NextPageToken)
	})

	t.Run("when request is nil, then empty body sent", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Empty(t, payload)
		}, `{}`)
		_, err := client.ListLikedResources(context.Background(), nil)
		require.NoError(t, err)
	})

	t.Run("when page size negative, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ListLikedResources(context.Background(), &ListLikedResourcesRequest{PageSize: -1})
		require.Error(t, err)
	})
}
