package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListCollectionItems(t *testing.T) {
	t.Run("when filters provided, then items are listed", func(t *testing.T) {
		pageSize := int32(25)
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			values := r.URL.Query()
			require.Equal(t, "col-1", values.Get("collectionId"))
			require.Equal(t, "card-1", values.Get("cardId"))
			require.Equal(t, "printing-1", values.Get("printingId"))
			require.Equal(t, "product-1", values.Get("productId"))
			require.Equal(t, "25", values.Get("pageSize"))
			require.Equal(t, "token", values.Get("nextToken"))

			resp := ListCollectionItemsResponse{
				Items: []CollectionItem{{ID: "item-1"}},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListCollectionItems(context.Background(), &ListCollectionItemsRequest{
			CollectionID: "col-1",
			CardID:       "card-1",
			PrintingID:   "printing-1",
			ProductID:    "product-1",
			PageSize:     &pageSize,
			NextToken:    "token",
		})
		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
	})
}

func TestClientGetCollectionItem(t *testing.T) {
	t.Run("when id provided, then item fetched", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v1/collection_items/item-1", r.URL.Path)
			resp := GetCollectionItemResponse{Item: &CollectionItem{ID: "item-1"}}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetCollectionItem(context.Background(), &GetCollectionItemRequest{ItemID: "item-1"})
		require.NoError(t, err)
		require.Equal(t, "item-1", resp.Item.ID)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetCollectionItem(context.Background(), &GetCollectionItemRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientCreateCollectionItem(t *testing.T) {
	t.Run("when request provided, then item created", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var body CreateCollectionItemRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.Equal(t, "col-1", body.CollectionID)
			require.Equal(t, "product-1", body.ProductID)
			require.Equal(t, int32(3), body.Quantity)
			require.Equal(t, ConditionNearMint, body.Condition)
			require.NotNil(t, body.Value)
			require.InEpsilon(t, 9.99, *body.Value, 1e-9)
			require.Equal(t, "item-client-1", body.ItemID)

			resp := CreateCollectionItemResponse{Item: &CollectionItem{ID: "item-1"}}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		value := 9.99
		resp, err := client.CreateCollectionItem(context.Background(), &CreateCollectionItemRequest{
			CollectionID: "col-1",
			ProductID:    "product-1",
			Quantity:     3,
			Condition:    ConditionNearMint,
			Value:        &value,
			ItemID:       "item-client-1",
		})
		require.NoError(t, err)
		require.Equal(t, "item-1", resp.Item.ID)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.CreateCollectionItem(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientUpdateCollectionItem(t *testing.T) {
	t.Run("when fields provided, then update request sent", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPut, r.Method)
			require.Equal(t, "/v1/collection_items/item-1", r.URL.Path)

			var body map[string]interface{}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.EqualValues(t, 5, body["quantity"])
			require.Equal(t, string(ConditionLightlyPlayed), body["condition"])
			require.EqualValues(t, 12.5, body["value"])

			resp := UpdateCollectionItemResponse{Item: &CollectionItem{ID: "item-1", Quantity: 5, Value: 12.5}}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		quantity := int32(5)
		condition := ConditionLightlyPlayed
		value := 12.5

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.UpdateCollectionItem(context.Background(), &UpdateCollectionItemRequest{
			ItemID:    "item-1",
			Quantity:  &quantity,
			Condition: &condition,
			Value:     &value,
		})
		require.NoError(t, err)
		require.Equal(t, int32(5), resp.Item.Quantity)
		require.InEpsilon(t, 12.5, resp.Item.Value, 1e-9)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.UpdateCollectionItem(context.Background(), &UpdateCollectionItemRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientDeleteCollectionItem(t *testing.T) {
	t.Run("when id provided, then delete succeeds", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/v1/collection_items/item-1", r.URL.Path)
			_, _ = w.Write([]byte(`{}`))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.DeleteCollectionItem(context.Background(), &DeleteCollectionItemRequest{ItemID: "item-1"})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.DeleteCollectionItem(context.Background(), &DeleteCollectionItemRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientBatchGetCollectionItems(t *testing.T) {
	t.Run("when ids provided, then items fetched", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v1/collection_items:batchGet", r.URL.Path)

			var body BatchGetCollectionItemsRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.ElementsMatch(t, []string{"item-1", "item-2"}, body.ItemIDs)
			require.True(t, body.AllowPartial)

			resp := BatchGetCollectionItemsResponse{
				Items:       []CollectionItem{{ID: "item-1"}},
				NotFoundIDs: []string{"missing"},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BatchGetCollectionItems(context.Background(), &BatchGetCollectionItemsRequest{
			ItemIDs:      []string{"item-1", "item-2"},
			AllowPartial: true,
		})
		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
		require.Equal(t, []string{"missing"}, resp.NotFoundIDs)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchGetCollectionItems(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
