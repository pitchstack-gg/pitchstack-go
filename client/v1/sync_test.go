package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestClientPullChanges(t *testing.T) {
	t.Run("when request valid, then payload encoded and response decoded", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/sync/changes:pull", r.URL.Path)

			var payload struct {
				Cursor struct {
					Overall string `json:"overall"`
				} `json:"cursor"`
				Filters struct {
					ResourceTypes  []ResourceType `json:"resourceTypes"`
					PinnedOnly     bool           `json:"pinnedOnly"`
					IncludeDeletes bool           `json:"includeDeletes"`
				} `json:"filters"`
				Limit int32 `json:"limit"`
			}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "42", payload.Cursor.Overall)
			require.Equal(t, []ResourceType{ResourceTypeDeck}, payload.Filters.ResourceTypes)
			require.True(t, payload.Filters.PinnedOnly)
			require.True(t, payload.Filters.IncludeDeletes)
			require.Equal(t, int32(25), payload.Limit)

			resp := PullChangesResponse{
				Changes: []Change{
					{
						EnvelopeID:   "1",
						ResourceType: ResourceTypeDeck,
						ResourceID:   "deck-1",
						Action:       SyncActionUpsert,
						Metadata:     map[string]string{"key": "value"},
					},
				},
				NextCursor: &Cursor{Overall: "43"},
				HasMore:    true,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

		limit := int32(25)
		resp, err := client.PullChanges(context.Background(), &PullChangesRequest{
			Cursor: &Cursor{Overall: "42"},
			Filters: &PullFilters{
				ResourceTypes:  []ResourceType{ResourceTypeDeck},
				PinnedOnly:     true,
				IncludeDeletes: true,
			},
			Limit: &limit,
		})
		require.NoError(t, err)
		require.True(t, resp.HasMore)
		require.Equal(t, "43", resp.NextCursor.Overall)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when request nil, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.PullChanges(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when filters invalid, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.PullChanges(context.Background(), &PullChangesRequest{
			Filters: &PullFilters{
				ResourceTypes: []ResourceType{ResourceTypeUnspecified},
			},
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientBatchApplyMutations(t *testing.T) {
	t.Run("when request valid, then payload sent and response decoded", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/sync/mutations:batchApply", r.URL.Path)

			var payload struct {
				ClientID  string     `json:"clientId"`
				Mutations []Mutation `json:"mutations"`
			}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "client-123", payload.ClientID)
			require.Len(t, payload.Mutations, 1)
			require.Equal(t, ResourceTypeCollection, payload.Mutations[0].ResourceType)
			require.Equal(t, SyncActionUpsert, payload.Mutations[0].Action)
			require.Equal(t, "resource-1", payload.Mutations[0].ResourceID)
			require.Equal(t, "Deck", payload.Mutations[0].Payload["name"])

			resp := BatchApplyMutationsResponse{
				Results: []MutationResult{
					{
						LocalChangeID: "local-1",
						Status:        SyncStatusOK,
						Snapshot: &ResourceSnapshot{
							ResourceID:   "resource-1",
							ResourceType: ResourceTypeCollection,
							Version:      "2",
						},
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		ts := time.Now().UTC()
		resp, err := client.BatchApplyMutations(context.Background(), &BatchApplyMutationsRequest{
			ClientID: "client-123",
			Mutations: []Mutation{
				{
					LocalChangeID:   "local-1",
					ResourceType:    ResourceTypeCollection,
					Action:          SyncActionUpsert,
					ResourceID:      "resource-1",
					Payload:         map[string]any{"name": "Deck"},
					BaseVersion:     "1",
					ClientTimestamp: &ts,
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.Results, 1)
		require.Equal(t, SyncStatusOK, resp.Results[0].Status)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when request nil, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchApplyMutations(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when mutations empty, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchApplyMutations(context.Background(), &BatchApplyMutationsRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when mutation invalid, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchApplyMutations(context.Background(), &BatchApplyMutationsRequest{
			Mutations: []Mutation{
				{
					ResourceType: ResourceTypeCollection,
					Action:       SyncActionUnspecified,
					ResourceID:   "resource-1",
				},
			},
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientAddPin(t *testing.T) {
	t.Run("when request valid, then pin added", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/sync/pins", r.URL.Path)
			var payload struct {
				ResourceID string `json:"resourceId"`
			}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "deck-1", payload.ResourceID)
			w.WriteHeader(http.StatusOK)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.AddPin(context.Background(), &AddPinRequest{ResourceID: "deck-1"})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when request nil, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.AddPin(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when resource id empty, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.AddPin(context.Background(), &AddPinRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientRemovePin(t *testing.T) {
	t.Run("when request valid, then pin removed", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/v1/sync/pins/deck-1", r.URL.Path)
			w.WriteHeader(http.StatusOK)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.RemovePin(context.Background(), &RemovePinRequest{ResourceID: "deck-1"})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when request nil, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.RemovePin(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when resource id empty, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.RemovePin(context.Background(), &RemovePinRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientListSubscriptions(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/sync/subscriptions", r.URL.Path)
		resp := ListSubscriptionsResponse{
			Subscriptions: []Subscription{
				{
					ResourceID:   "deck-1",
					ResourceType: ResourceTypeDeck,
					Reason:       SubscriptionReasonOwner,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ListSubscriptions(context.Background())
	require.NoError(t, err)
	require.Len(t, resp.Subscriptions, 1)
	require.Equal(t, SubscriptionReasonOwner, resp.Subscriptions[0].Reason)
	require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
}
