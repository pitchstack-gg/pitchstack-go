package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientGetChangeSet(t *testing.T) {
	t.Run("when request valid, then response decoded", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/sync/changeSet", r.URL.Path)
			query := r.URL.Query()
			require.Equal(t, "cursor-1", query.Get("cursor"))
			require.Equal(t, "25", query.Get("pageSize"))
			require.Equal(t, "true", query.Get("includeDocuments"))

			resp := GetChangeSetResponse{
				Events: []SyncEvent{
					{
						EventID: "evt-1",
						Resource: &ResourceDescriptor{
							Type: ResourceTypeDeck,
							ID:   "deck-1",
						},
						Kind:    SyncEventKindUpdated,
						Version: "v2",
					},
				},
				NextCursor: "cursor-2",
				HasMore:    true,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		pageSize := int32(25)
		resp, err := client.GetChangeSet(context.Background(), &GetChangeSetRequest{
			Cursor:           "cursor-1",
			PageSize:         &pageSize,
			IncludeDocuments: true,
		})
		require.NoError(t, err)
		require.True(t, resp.HasMore)
		require.Equal(t, "cursor-2", resp.NextCursor)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when request nil, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetChangeSet(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientBatchApplyChanges(t *testing.T) {
	t.Run("when request valid, then payload sent and response decoded", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/sync:batchApply", r.URL.Path)

			var payload BatchApplyChangesRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "device-abc", payload.DeviceID)
			require.Len(t, payload.Changes, 1)
			require.Equal(t, "client-1", payload.Changes[0].ClientChangeID)
			require.Equal(t, ResourceTypeCollection, payload.Changes[0].Resource.Type)
			require.Equal(t, "resource-1", payload.Changes[0].Resource.ID)
			require.Equal(t, SyncActionUpsert, payload.Changes[0].Action)
			require.Equal(t, "Deck", payload.Changes[0].Document["name"])

			resp := BatchApplyChangesResponse{
				Results: []AppliedChangeResult{
					{
						ClientChangeID: "client-1",
						Status:         SyncStatusOK,
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BatchApplyChanges(context.Background(), &BatchApplyChangesRequest{
			DeviceID: "device-abc",
			Changes: []LocalChange{
				{
					ClientChangeID: "client-1",
					Resource: &ResourceDescriptor{
						Type: ResourceTypeCollection,
						ID:   "resource-1",
					},
					Action:   SyncActionUpsert,
					Document: map[string]any{"name": "Deck"},
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
		resp, err := client.BatchApplyChanges(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when changes empty, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchApplyChanges(context.Background(), &BatchApplyChangesRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when change invalid, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchApplyChanges(context.Background(), &BatchApplyChangesRequest{
			Changes: []LocalChange{
				{
					Resource: &ResourceDescriptor{
						Type: ResourceTypeCollection,
						ID:   "resource-1",
					},
					Action: SyncActionUnspecified,
				},
			},
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientUpdateSubscriptions(t *testing.T) {
	t.Run("when request valid, then payload sent", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/sync/subscriptions:update", r.URL.Path)

			var payload UpdateSubscriptionsRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Len(t, payload.Mutations, 1)
			require.Equal(t, SubscriptionMutationTypeSubscribe, payload.Mutations[0].Type)
			require.Equal(t, ResourceTypeDeck, payload.Mutations[0].Resource.Type)
			require.Equal(t, "deck-1", payload.Mutations[0].Resource.ID)

			resp := UpdateSubscriptionsResponse{
				Subscriptions: []ResourceSubscription{
					{
						SubscriptionID: "sub-1",
						Resource:       payload.Mutations[0].Resource,
						Source:         SubscriptionSourceManual,
					},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.UpdateSubscriptions(context.Background(), &UpdateSubscriptionsRequest{
			Mutations: []SubscriptionMutation{
				{
					Type: SubscriptionMutationTypeSubscribe,
					Resource: &ResourceDescriptor{
						Type: ResourceTypeDeck,
						ID:   "deck-1",
					},
				},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.Subscriptions, 1)
		require.Equal(t, "sub-1", resp.Subscriptions[0].SubscriptionID)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when request nil, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.UpdateSubscriptions(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when no mutations provided, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.UpdateSubscriptions(context.Background(), &UpdateSubscriptionsRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when mutation invalid, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.UpdateSubscriptions(context.Background(), &UpdateSubscriptionsRequest{
			Mutations: []SubscriptionMutation{
				{
					Type: SubscriptionMutationTypeUnspecified,
					Resource: &ResourceDescriptor{
						Type: ResourceTypeDeck,
						ID:   "deck-1",
					},
				},
			},
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientListSubscriptions(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/sync/subscriptions", r.URL.Path)
		resp := ListSubscriptionsResponse{
			Subscriptions: []ResourceSubscription{
				{
					SubscriptionID: "sub-1",
					Resource: &ResourceDescriptor{
						Type: ResourceTypeDeck,
						ID:   "deck-1",
					},
					Source: SubscriptionSourceOwned,
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

	server := httptest.NewServer(http.HandlerFunc(handler))
	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	t.Cleanup(server.Close)

	resp, err := client.ListSubscriptions(context.Background())
	require.NoError(t, err)
	require.Len(t, resp.Subscriptions, 1)
	require.Equal(t, SubscriptionSourceOwned, resp.Subscriptions[0].Source)
	require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
}
