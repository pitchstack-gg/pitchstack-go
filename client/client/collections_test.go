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

func TestClientListCollections(t *testing.T) {
	t.Run("when request has filters, then query string and response are handled", func(t *testing.T) {
		pageSize := int32(50)
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "COLLECTION_LIST_SCOPE_SHARED", r.URL.Query().Get("scope"))
			require.Equal(t, "user-123", r.URL.Query().Get("userId"))
			require.Equal(t, "50", r.URL.Query().Get("pageSize"))
			require.Equal(t, "token-abc", r.URL.Query().Get("nextToken"))

			payload := ListCollectionsResponse{
				Collections: []Collection{
					{
						ID:         "col-1",
						Name:       "Collection",
						UserID:     "user-123",
						Visibility: VisibilityLevelShared,
					},
				},
				NextToken: "next",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(payload)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

		resp, err := client.ListCollections(context.Background(), &ListCollectionsRequest{
			Scope:     CollectionListScopeShared,
			UserID:    "user-123",
			PageSize:  &pageSize,
			NextToken: "token-abc",
		})
		require.NoError(t, err)
		require.Len(t, resp.Collections, 1)
		require.Equal(t, "col-1", resp.Collections[0].ID)
		require.Equal(t, "next", resp.NextToken)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when server returns error, then api error is surfaced", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"code":16,"message":"unauthorized"}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

		resp, err := client.ListCollections(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)

		var apiErr *APIError
		require.ErrorAs(t, err, &apiErr)
		require.Equal(t, 16, apiErr.Status.Code)
	})
}

func TestClientCreateCollection(t *testing.T) {
	t.Run("when payload is valid, then collection is created", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload CreateCollectionRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "My Collection", payload.Name)
			require.Equal(t, "A description", payload.Description)
			require.Equal(t, VisibilityLevelPrivate, payload.Visibility)

			response := CreateCollectionResponse{
				Collection: &Collection{ID: "col-1"},
			}

			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

		resp, err := client.CreateCollection(context.Background(), &CreateCollectionRequest{
			Name:        "My Collection",
			Description: "A description",
			Visibility:  VisibilityLevelPrivate,
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Collection)
		require.Equal(t, "col-1", resp.Collection.ID)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.CreateCollection(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetCollection(t *testing.T) {
	t.Run("when id provided, then collection fetched", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/collections/col-1", r.URL.Path)
			response := GetCollectionResponse{
				Collection: &Collection{ID: "col-1"},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetCollection(context.Background(), &GetCollectionRequest{CollectionID: "col-1"})
		require.NoError(t, err)
		require.Equal(t, "col-1", resp.Collection.ID)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetCollection(context.Background(), &GetCollectionRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientUpdateCollection(t *testing.T) {
	t.Run("when fields provided, then partial update is sent", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "/api/v1/collections/col-1", r.URL.Path)

			var body map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.Equal(t, "New Name", body["name"])
			require.Equal(t, "New Description", body["description"])

			resp := UpdateCollectionResponse{
				Collection: &Collection{ID: "col-1", Name: "New Name"},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		name := "New Name"
		description := "New Description"

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.UpdateCollection(context.Background(), &UpdateCollectionRequest{
			CollectionID: "col-1",
			Name:         &name,
			Description:  &description,
		})
		require.NoError(t, err)
		require.Equal(t, "New Name", resp.Collection.Name)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.UpdateCollection(context.Background(), &UpdateCollectionRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientDeleteCollection(t *testing.T) {
	t.Run("when id provided, then delete is issued", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/api/v1/collections/col-1", r.URL.Path)
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{}`))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.DeleteCollection(context.Background(), &DeleteCollectionRequest{CollectionID: "col-1"})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.DeleteCollection(context.Background(), &DeleteCollectionRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientBatchGetCollections(t *testing.T) {
	t.Run("when ids provided, then collections returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/collections:batchGet", r.URL.Path)

			var body BatchGetCollectionsRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.ElementsMatch(t, []string{"col-1", "col-2"}, body.CollectionIDs)
			require.True(t, body.AllowPartial)

			resp := BatchGetCollectionsResponse{
				Collections: []Collection{{ID: "col-1"}},
				NotFoundIDs: []string{"missing"},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BatchGetCollections(context.Background(), &BatchGetCollectionsRequest{
			CollectionIDs: []string{"col-1", "col-2"},
			AllowPartial:  true,
		})
		require.NoError(t, err)
		require.Len(t, resp.Collections, 1)
		require.Equal(t, []string{"missing"}, resp.NotFoundIDs)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchGetCollections(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientCollectionsSync(t *testing.T) {
	t.Run("when ops provided, then sync endpoint is invoked", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/collections:sync", r.URL.Path)

			var body CollectionsSyncRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.Len(t, body.Ops, 1)
			require.Equal(t, SyncActionUpsert, body.Ops[0].Action)

			resp := CollectionsSyncResponse{
				Results: []CollectionSyncResult{{OpID: "1", Status: SyncStatusOK}},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.CollectionsSync(context.Background(), &CollectionsSyncRequest{
			Ops: []CollectionSyncOp{
				{OpID: "1", Action: SyncActionUpsert, Name: "Name"},
			},
		})
		require.NoError(t, err)
		require.Len(t, resp.Results, 1)
		require.Equal(t, SyncStatusOK, resp.Results[0].Status)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.CollectionsSync(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetCollectionValuation(t *testing.T) {
	t.Run("when source provided, then valuation is returned", func(t *testing.T) {
		now := time.Now().UTC().Round(time.Second)
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/collections/col-1/valuation", r.URL.Path)
			require.Equal(t, "external", r.URL.Query().Get("source"))

			resp := GetCollectionValuationResponse{
				CollectionID:        "col-1",
				TotalEstimatedValue: 123.45,
				ComputedAt:          &now,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetCollectionValuation(context.Background(), &GetCollectionValuationRequest{
			CollectionID: "col-1",
			Source:       "external",
		})
		require.NoError(t, err)
		require.Equal(t, float64(123.45), resp.TotalEstimatedValue)
		require.NotNil(t, resp.ComputedAt)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetCollectionValuation(context.Background(), &GetCollectionValuationRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientListStarredCollections(t *testing.T) {
	t.Run("when user provided, then starred collections returned", func(t *testing.T) {
		limit := int32(10)
		offset := int32(5)
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/users/user-1/stars/collections", r.URL.Path)
			require.Equal(t, "10", r.URL.Query().Get("limit"))
			require.Equal(t, "5", r.URL.Query().Get("offset"))

			resp := ListStarredCollectionsResponse{
				Collections: []Collection{{ID: "col-1"}},
				Total:       42,
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListStarredCollections(context.Background(), &ListStarredCollectionsRequest{
			UserID: "user-1",
			Limit:  &limit,
			Offset: &offset,
		})
		require.NoError(t, err)
		require.Equal(t, int32(42), resp.Total)
		require.Len(t, resp.Collections, 1)
	})

	t.Run("when user id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.ListStarredCollections(context.Background(), &ListStarredCollectionsRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientStarCollection(t *testing.T) {
	t.Run("when id provided, then star request succeeds", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/api/v1/collections/col-1/stars", r.URL.Path)
			_, _ = w.Write([]byte(`{}`))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.StarCollection(context.Background(), &StarCollectionRequest{CollectionID: "col-1"})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.StarCollection(context.Background(), &StarCollectionRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientUnstarCollection(t *testing.T) {
	t.Run("when id provided, then unstar request succeeds", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/api/v1/collections/col-1/stars", r.URL.Path)
			_, _ = w.Write([]byte(`{}`))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.UnstarCollection(context.Background(), &UnstarCollectionRequest{CollectionID: "col-1"})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.UnstarCollection(context.Background(), &UnstarCollectionRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
