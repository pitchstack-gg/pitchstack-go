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
			require.Equal(t, "subject-1", r.URL.Query().Get("subjectId"))

			payload := ListCollectionsResponse{
				Collections: []Collection{
					{
						ID:             "col-1",
						Name:           "Collection",
						OwnerID:        "user-123",
						CollectionType: CollectionTypeBinder,
						Visibility:     VisibilityLevelShared,
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
			SubjectID: "subject-1",
		})
		require.NoError(t, err)
		require.Len(t, resp.Collections, 1)
		require.Equal(t, "col-1", resp.Collections[0].ID)
		require.Equal(t, "next", resp.NextToken)
		require.Equal(t, "user-123", resp.Collections[0].OwnerID)
		require.Equal(t, CollectionTypeBinder, resp.Collections[0].CollectionType)
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
			require.Equal(t, CollectionTypeWantlist, payload.CollectionType)
			require.Equal(t, VisibilityLevelPrivate, payload.Visibility)
			require.Equal(t, "col-client-1", payload.CollectionID)

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
			Name:           "My Collection",
			CollectionType: CollectionTypeWantlist,
			Description:    "A description",
			Visibility:     VisibilityLevelPrivate,
			CollectionID:   "col-client-1",
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
			require.Equal(t, "/v1/collections/col-1", r.URL.Path)
			response := GetCollectionResponse{
				Collection: &Collection{ID: "col-1"},
				Stats: &CollectionStats{
					ItemsCount:    10,
					QuantityCount: 15,
				},
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
		require.NotNil(t, resp.Stats)
		require.Equal(t, int32(10), resp.Stats.ItemsCount)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetCollection(context.Background(), &GetCollectionRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetCollectionHistory(t *testing.T) {
	t.Run("when id provided, then history fetched", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/collections/col-1/history", r.URL.Path)
			response := GetCollectionHistoryResponse{
				Changes: []CollectionHistoryChange{{
					ID:           "change-1",
					CollectionID: "col-1",
					EventType:    CollectionHistoryEventTypeItemAdded,
					ItemChanges: []CollectionHistoryItemChange{{
						ItemID:    "item-1",
						ProductID: "product-1",
						Operation: CollectionHistoryItemChangeOperationAdd,
					}},
				}},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(response)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetCollectionHistory(context.Background(), &GetCollectionHistoryRequest{CollectionID: "col-1"})
		require.NoError(t, err)
		require.Len(t, resp.Changes, 1)
		require.Equal(t, CollectionHistoryEventTypeItemAdded, resp.Changes[0].EventType)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetCollectionHistory(context.Background(), &GetCollectionHistoryRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientUpdateCollection(t *testing.T) {
	t.Run("when fields provided, then partial update is sent", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "/v1/collections/col-1", r.URL.Path)

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

func TestClientUpdateCollectionVisibility(t *testing.T) {
	t.Run("when visibility provided, then targeted update is sent", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPatch, r.Method)
			require.Equal(t, "/v1/collections/col-1/visibility", r.URL.Path)

			var body map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.Equal(t, VisibilityLevelPublic, VisibilityLevel(body["visibility"].(string)))

			resp := UpdateCollectionVisibilityResponse{
				Collection: &Collection{ID: "col-1", Visibility: VisibilityLevelPublic},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		visibility := VisibilityLevelPublic
		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.UpdateCollectionVisibility(context.Background(), &UpdateCollectionVisibilityRequest{
			CollectionID: "col-1",
			Visibility:   &visibility,
		})
		require.NoError(t, err)
		require.Equal(t, VisibilityLevelPublic, resp.Collection.Visibility)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		visibility := VisibilityLevelShared
		resp, err := client.UpdateCollectionVisibility(context.Background(), &UpdateCollectionVisibilityRequest{
			Visibility: &visibility,
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when visibility missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.UpdateCollectionVisibility(context.Background(), &UpdateCollectionVisibilityRequest{
			CollectionID: "col-1",
		})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientUpdateCollectionArt(t *testing.T) {
	t.Run("when art update is provided, then targeted update is sent", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "/v1/collections/col-1/art", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload UpdateCollectionArtRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "printing-1", payload.SelectedArtPrintingID)
			require.False(t, payload.ClearSelectedArt)

			require.NoError(t, json.NewEncoder(w).Encode(UpdateCollectionArtResponse{
				Collection: &Collection{ID: "col-1", SelectedArtPrintingID: "printing-1"},
			}))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.UpdateCollectionArt(context.Background(), &UpdateCollectionArtRequest{
			CollectionID:          "col-1",
			SelectedArtPrintingID: "printing-1",
		})
		require.NoError(t, err)
		require.Equal(t, "printing-1", resp.Collection.SelectedArtPrintingID)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.UpdateCollectionArt(context.Background(), &UpdateCollectionArtRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientDeleteCollection(t *testing.T) {
	t.Run("when id provided, then delete is issued", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/v1/collections/col-1", r.URL.Path)
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
			require.Equal(t, "/v1/collections:batchGet", r.URL.Path)

			var body BatchGetCollectionsRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.ElementsMatch(t, []string{"col-1", "col-2"}, body.CollectionIDs)
			require.True(t, body.AllowPartial)

			resp := BatchGetCollectionsResponse{
				Collections: []CollectionWithStats{
					{
						Collection: &Collection{ID: "col-1"},
						Stats:      &CollectionStats{ItemsCount: 7},
					},
				},
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
		require.Equal(t, "col-1", resp.Collections[0].Collection.ID)
		require.NotNil(t, resp.Collections[0].Stats)
		require.Equal(t, int32(7), resp.Collections[0].Stats.ItemsCount)
		require.Equal(t, []string{"missing"}, resp.NotFoundIDs)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchGetCollections(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetCollectionValuation(t *testing.T) {
	t.Run("when source provided, then valuation is returned", func(t *testing.T) {
		now := time.Now().UTC().Round(time.Second)
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v1/collections/col-1/valuation", r.URL.Path)
			require.Equal(t, "external", r.URL.Query().Get("source"))

			resp := GetCollectionValuationResponse{
				CollectionID:        "col-1",
				TotalEstimatedValue: 123.45,
				ComputedAt:          &now,
				TotalQuantity:       17,
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
		require.Equal(t, int32(17), resp.TotalQuantity)
		require.NotNil(t, resp.ComputedAt)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetCollectionValuation(context.Background(), &GetCollectionValuationRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientListTradeItems(t *testing.T) {
	pageSize := int32(25)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/trade_items", r.URL.Path)
		require.Equal(t, "user-1", r.URL.Query().Get("userId"))
		require.Equal(t, "25", r.URL.Query().Get("pageSize"))
		require.Equal(t, "token", r.URL.Query().Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(ListTradeItemsResponse{
			Items: []TradeItem{{
				Item:             &CollectionItem{ID: "item-1", TradeQuantity: 2, Notes: "foil"},
				SourceCollection: &Collection{ID: "col-1"},
			}},
			NextToken: "next",
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ListTradeItems(context.Background(), &ListTradeItemsRequest{
		UserID:    "user-1",
		PageSize:  &pageSize,
		NextToken: "token",
	})
	require.NoError(t, err)
	require.Equal(t, "item-1", resp.Items[0].Item.ID)
	require.Equal(t, int32(2), resp.Items[0].Item.TradeQuantity)
	require.Equal(t, "next", resp.NextToken)
}

func TestClientGrantCollectionAccess(t *testing.T) {
	t.Run("when payload valid, then permission is granted", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/collections/permissions:grant", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "col-1", payload["resourceId"])
			require.Equal(t, "user-42", payload["subjectId"])
			require.Equal(t, string(CollectionPermissionReader), payload["permission"])

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GrantCollectionAccess(context.Background(), &GrantCollectionAccessRequest{
			CollectionID: "col-1",
			SubjectID:    "user-42",
			Permission:   CollectionPermissionReader,
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when required fields missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GrantCollectionAccess(context.Background(), &GrantCollectionAccessRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientRevokeCollectionAccess(t *testing.T) {
	t.Run("when payload valid, then permission is revoked", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/collections/permissions:revoke", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "col-1", payload["resourceId"])
			require.Equal(t, "user-42", payload["subjectId"])
			require.Equal(t, string(CollectionPermissionWriter), payload["permission"])

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{}`))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.RevokeCollectionAccess(context.Background(), &RevokeCollectionAccessRequest{
			CollectionID: "col-1",
			SubjectID:    "user-42",
			Permission:   CollectionPermissionWriter,
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when required fields missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.RevokeCollectionAccess(context.Background(), &RevokeCollectionAccessRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetCollectionAccess(t *testing.T) {
	t.Run("when id provided, then access is returned", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/collections/col-1/access", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(GetCollectionAccessResponse{
				Permission: CollectionPermissionReader,
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetCollectionAccess(context.Background(), &GetCollectionAccessRequest{CollectionID: "col-1"})
		require.NoError(t, err)
		require.Equal(t, CollectionPermissionReader, resp.Permission)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetCollectionAccess(context.Background(), &GetCollectionAccessRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientStopCollectionShare(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/collections/col-1/access:stop", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.StopCollectionShare(context.Background(), &StopCollectionShareRequest{CollectionID: "col-1"})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)

	resp, err = client.StopCollectionShare(context.Background(), &StopCollectionShareRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientListCollectionAccessGrants(t *testing.T) {
	t.Run("when request includes pagination, then query and grants are returned", func(t *testing.T) {
		pageSize := int32(10)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/collections/col-1/permissions", r.URL.Path)
			require.Equal(t, "10", r.URL.Query().Get("pageSize"))
			require.Equal(t, "token", r.URL.Query().Get("nextToken"))

			require.NoError(t, json.NewEncoder(w).Encode(ListCollectionAccessGrantsResponse{
				Grants: []CollectionAccessGrant{
					{SubjectID: "user-1", Permission: CollectionPermissionWriter},
				},
				NextToken: "next",
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListCollectionAccessGrants(context.Background(), &ListCollectionAccessGrantsRequest{
			CollectionID: "col-1",
			PageSize:     &pageSize,
			NextToken:    "token",
		})
		require.NoError(t, err)
		require.Len(t, resp.Grants, 1)
		require.Equal(t, "user-1", resp.Grants[0].SubjectID)
		require.Equal(t, CollectionPermissionWriter, resp.Grants[0].Permission)
		require.Equal(t, "next", resp.NextToken)
	})

	t.Run("when request missing id, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.ListCollectionAccessGrants(context.Background(), &ListCollectionAccessGrantsRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientExportCollection(t *testing.T) {
	t.Run("when request includes pagination, then query and response are handled", func(t *testing.T) {
		pageSize := int32(25)
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/collections/col-1:export", r.URL.Path)
			require.Equal(t, "25", r.URL.Query().Get("pageSize"))
			require.Equal(t, "token", r.URL.Query().Get("nextToken"))

			require.NoError(t, json.NewEncoder(w).Encode(ExportCollectionResponse{
				Collection: &Collection{ID: "col-1"},
				Items: []CollectionItem{
					{ID: "item-1"},
				},
				NextToken: "next",
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ExportCollection(context.Background(), &ExportCollectionRequest{
			CollectionID: "col-1",
			PageSize:     &pageSize,
			NextToken:    "token",
		})
		require.NoError(t, err)
		require.Equal(t, "col-1", resp.Collection.ID)
		require.Len(t, resp.Items, 1)
		require.Equal(t, "item-1", resp.Items[0].ID)
		require.Equal(t, "next", resp.NextToken)
	})

	t.Run("when request missing id, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.ExportCollection(context.Background(), &ExportCollectionRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientImportCollection(t *testing.T) {
	t.Run("when request is valid, then payload is sent and response parsed", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/collections:import", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload ImportCollectionRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "col-1", payload.CollectionID)
			require.Equal(t, "Imported Collection", payload.Name)
			require.Equal(t, CollectionTypeBinder, payload.CollectionType)
			require.Equal(t, VisibilityLevelPrivate, payload.Visibility)
			require.Len(t, payload.Items, 1)
			require.Equal(t, "item-1", payload.Items[0].ItemID)

			require.NoError(t, json.NewEncoder(w).Encode(ImportCollectionResponse{
				Collection:    &Collection{ID: "col-1"},
				Stats:         &CollectionStats{ItemsCount: 1},
				ImportedItems: 1,
			}))
		}))
		t.Cleanup(server.Close)

		value := 12.5
		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ImportCollection(context.Background(), &ImportCollectionRequest{
			CollectionID:   "col-1",
			Name:           "Imported Collection",
			CollectionType: CollectionTypeBinder,
			Visibility:     VisibilityLevelPrivate,
			Items: []ImportCollectionItem{
				{
					ItemID:    "item-1",
					ProductID: "prod-1",
					Quantity:  2,
					Condition: ConditionNearMint,
					Value:     &value,
				},
			},
		})
		require.NoError(t, err)
		require.Equal(t, "col-1", resp.Collection.ID)
		require.Equal(t, int32(1), resp.Stats.ItemsCount)
		require.Equal(t, int32(1), resp.ImportedItems)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.ImportCollection(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
