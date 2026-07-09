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

func TestClientGetProductPrice(t *testing.T) {
	t.Run("when card id provided, then price entry returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v1/prices/prod-1", r.URL.Path)
			require.Equal(t, "TCGPlayer", r.URL.Query().Get("source"))

			resp := GetProductPriceResponse{
				Entry: &PriceEntry{
					EntryID:   "entry-1",
					ProductID: "prod-1",
					Price:     12.34,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetProductPrice(context.Background(), &GetProductPriceRequest{
			ProductID: "prod-1",
			Source:    "TCGPlayer",
		})
		require.NoError(t, err)
		require.Equal(t, "prod-1", resp.Entry.ProductID)
		require.Equal(t, 12.34, resp.Entry.Price)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetProductPrice(context.Background(), &GetProductPriceRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetProductPriceHistory(t *testing.T) {
	t.Run("when filters provided, then history returned", func(t *testing.T) {
		limit := int32(50)
		now := time.Now().UTC().Round(time.Second)

		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v1/prices/prod-1/history", r.URL.Path)
			require.Equal(t, "2024-01-01", r.URL.Query().Get("startDate"))
			require.Equal(t, "2024-02-01", r.URL.Query().Get("endDate"))
			require.Equal(t, "50", r.URL.Query().Get("limit"))
			require.Equal(t, "TCGPlayer", r.URL.Query().Get("source"))

			resp := GetProductPriceHistoryResponse{
				ProductID: "prod-1",
				Source:    "TCGPlayer",
				StartTime: &now,
				EndTime:   &now,
				Entries: []PriceEntry{
					{EntryID: "entry-1", Price: 10.0},
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetProductPriceHistory(context.Background(), &GetProductPriceHistoryRequest{
			ProductID: "prod-1",
			StartDate: "2024-01-01",
			EndDate:   "2024-02-01",
			Limit:     &limit,
			Source:    "TCGPlayer",
		})
		require.NoError(t, err)
		require.Equal(t, "prod-1", resp.ProductID)
		require.Len(t, resp.Entries, 1)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetProductPriceHistory(context.Background(), &GetProductPriceHistoryRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientBatchGetProductPrices(t *testing.T) {
	t.Run("when ids provided, then prices returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			var body BatchGetProductPricesRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.ElementsMatch(t, []string{"prod-1", "prod-2"}, body.ProductIDs)
			require.Equal(t, "TCGPlayer", body.Source)

			resp := BatchGetProductPricesResponse{
				Prices: []PriceEntry{
					{EntryID: "entry-1", ProductID: "prod-1"},
				},
				NotFound: []string{"prod-3"},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BatchGetProductPrices(context.Background(), &BatchGetProductPricesRequest{
			ProductIDs: []string{"prod-1", "prod-2"},
			Source:     "TCGPlayer",
		})
		require.NoError(t, err)
		require.Len(t, resp.Prices, 1)
		require.Equal(t, []string{"prod-3"}, resp.NotFound)
	})

	t.Run("when ids missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchGetProductPrices(context.Background(), &BatchGetProductPricesRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchGetProductPrices(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientProductPriceWatches(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/price-watches":
			require.Equal(t, "true", r.URL.Query().Get("activeOnly"))
			require.Equal(t, []string{"prod-1", "prod-2"}, r.URL.Query()["productIds"])
			require.NoError(t, json.NewEncoder(w).Encode(ListProductPriceWatchesResponse{
				Watches: []ProductPriceWatch{{WatchID: "watch-1", ProductID: "prod-1"}},
			}))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/price-watches":
			var payload CreateProductPriceWatchRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "prod-1", payload.ProductID)
			require.Equal(t, "down", payload.Direction)
			require.NoError(t, json.NewEncoder(w).Encode(CreateProductPriceWatchResponse{Watch: &ProductPriceWatch{WatchID: "watch-1"}}))
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/price-watches/watch-1":
			var payload UpdateProductPriceWatchRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.NotNil(t, payload.Active)
			require.NoError(t, json.NewEncoder(w).Encode(UpdateProductPriceWatchResponse{Watch: &ProductPriceWatch{WatchID: "watch-1", Active: *payload.Active}}))
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/price-watches/watch-1":
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	activeOnly := true
	listResp, err := client.ListProductPriceWatches(context.Background(), &ListProductPriceWatchesRequest{
		ActiveOnly: &activeOnly,
		ProductIDs: []string{"prod-1", "prod-2"},
	})
	require.NoError(t, err)
	require.Len(t, listResp.Watches, 1)

	createResp, err := client.CreateProductPriceWatch(context.Background(), &CreateProductPriceWatchRequest{
		ProductID: "prod-1",
		Source:    "TCGPlayer",
		Direction: "down",
		Period:    "7d",
	})
	require.NoError(t, err)
	require.Equal(t, "watch-1", createResp.Watch.WatchID)

	active := false
	updateResp, err := client.UpdateProductPriceWatch(context.Background(), &UpdateProductPriceWatchRequest{WatchID: "watch-1", Active: &active})
	require.NoError(t, err)
	require.False(t, updateResp.Watch.Active)

	deleteResp, err := client.DeleteProductPriceWatch(context.Background(), &DeleteProductPriceWatchRequest{WatchID: "watch-1"})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, deleteResp.Metadata.StatusCode)
}

func TestClientProductPriceWatchLists(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/price-watch-lists":
			require.NoError(t, json.NewEncoder(w).Encode(ListProductPriceWatchListsResponse{
				Lists: []ProductPriceWatchList{{ListID: "list-1", Name: "Deals"}},
			}))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/price-watch-lists":
			var payload CreateProductPriceWatchListRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "Deals", payload.Name)
			require.Equal(t, "cards to watch", payload.Description)
			require.NoError(t, json.NewEncoder(w).Encode(CreateProductPriceWatchListResponse{
				List: &ProductPriceWatchList{ListID: "list-1", Name: payload.Name},
			}))
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/price-watch-lists/list-1":
			var payload UpdateProductPriceWatchListRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.NotNil(t, payload.Name)
			require.Equal(t, "Updated", *payload.Name)
			require.NoError(t, json.NewEncoder(w).Encode(UpdateProductPriceWatchListResponse{
				List: &ProductPriceWatchList{ListID: "list-1", Name: *payload.Name},
			}))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/price-watch-lists/list-1/items":
			require.Equal(t, "true", r.URL.Query().Get("activeOnly"))
			require.NoError(t, json.NewEncoder(w).Encode(ListProductPriceWatchListItemsResponse{
				Items: []ProductPriceWatchListItem{{ItemID: "item-1", ListID: "list-1", Watch: &ProductPriceWatch{WatchID: "watch-1"}}},
			}))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/price-watch-lists/list-1/items":
			var payload AddProductPriceWatchListItemRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "watch-1", payload.WatchID)
			require.NoError(t, json.NewEncoder(w).Encode(AddProductPriceWatchListItemResponse{
				Item: &ProductPriceWatchListItem{ItemID: "item-1", ListID: "list-1"},
			}))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/price-watch-lists/list-1/items:batchAddProducts":
			var payload BatchAddProductsToProductPriceWatchListRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, []string{"prod-1", "prod-2"}, payload.ProductIDs)
			require.Equal(t, "TCGPlayer", payload.Source)
			require.NoError(t, json.NewEncoder(w).Encode(BatchAddProductsToProductPriceWatchListResponse{
				Items:    []ProductPriceWatchListItem{{ItemID: "item-2"}},
				Failures: []BatchAddProductPriceWatchFailure{{ProductID: "prod-3", Code: "not_found"}},
			}))
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/price-watch-lists/list-1/items/watch-1":
			w.Header().Set("X-Request-Id", "req-remove")
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/price-watch-lists/list-1":
			w.Header().Set("X-Request-Id", "req-delete")
			w.WriteHeader(http.StatusOK)
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	listsResp, err := client.ListProductPriceWatchLists(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, "list-1", listsResp.Lists[0].ListID)

	createResp, err := client.CreateProductPriceWatchList(context.Background(), &CreateProductPriceWatchListRequest{
		Name:        "Deals",
		Description: "cards to watch",
	})
	require.NoError(t, err)
	require.Equal(t, "Deals", createResp.List.Name)

	name := "Updated"
	updateResp, err := client.UpdateProductPriceWatchList(context.Background(), &UpdateProductPriceWatchListRequest{
		ListID: "list-1",
		Name:   &name,
	})
	require.NoError(t, err)
	require.Equal(t, "Updated", updateResp.List.Name)

	activeOnly := true
	itemsResp, err := client.ListProductPriceWatchListItems(context.Background(), &ListProductPriceWatchListItemsRequest{
		ListID:     "list-1",
		ActiveOnly: &activeOnly,
	})
	require.NoError(t, err)
	require.Equal(t, "watch-1", itemsResp.Items[0].Watch.WatchID)

	addResp, err := client.AddProductPriceWatchListItem(context.Background(), &AddProductPriceWatchListItemRequest{
		ListID:  "list-1",
		WatchID: "watch-1",
	})
	require.NoError(t, err)
	require.Equal(t, "item-1", addResp.Item.ItemID)

	batchResp, err := client.BatchAddProductsToProductPriceWatchList(context.Background(), &BatchAddProductsToProductPriceWatchListRequest{
		ListID:     "list-1",
		ProductIDs: []string{"prod-1", "prod-2"},
		Source:     "TCGPlayer",
	})
	require.NoError(t, err)
	require.Len(t, batchResp.Items, 1)
	require.Equal(t, "prod-3", batchResp.Failures[0].ProductID)

	removeResp, err := client.RemoveProductPriceWatchListItem(context.Background(), &RemoveProductPriceWatchListItemRequest{ListID: "list-1", WatchID: "watch-1"})
	require.NoError(t, err)
	require.Equal(t, "req-remove", removeResp.Metadata.RequestID)

	deleteResp, err := client.DeleteProductPriceWatchList(context.Background(), &DeleteProductPriceWatchListRequest{ListID: "list-1"})
	require.NoError(t, err)
	require.Equal(t, "req-delete", deleteResp.Metadata.RequestID)
}

func TestClientGetProductPriceWatchDigest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/v1/price-watch-digests/digest-1":
			require.Equal(t, http.MethodGet, r.Method)
			require.NoError(t, json.NewEncoder(w).Encode(GetProductPriceWatchDigestResponse{
				Digest: &ProductPriceWatchDigest{
					DigestID:       "digest-1",
					DigestDate:     "2026-07-08",
					ItemCount:      1,
					WatchListCount: 1,
					Groups: []ProductPriceWatchDigestGroup{{
						ListID: "list-1",
						Name:   "Deals",
						Items: []ProductPriceWatchDigestItem{{
							ItemID:        "item-1",
							ProductID:     "prod-1",
							CurrentPrice:  10,
							BaselinePrice: 8,
						}},
					}},
				},
			}))
		case "/v1/price-watch-digests/by-date/2026-07-08":
			require.Equal(t, http.MethodGet, r.Method)
			require.NoError(t, json.NewEncoder(w).Encode(GetProductPriceWatchDigestResponse{
				Digest: &ProductPriceWatchDigest{DigestID: "digest-2", DigestDate: "2026-07-08"},
			}))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	byID, err := client.GetProductPriceWatchDigest(context.Background(), &GetProductPriceWatchDigestRequest{DigestID: "digest-1"})
	require.NoError(t, err)
	require.Equal(t, "digest-1", byID.Digest.DigestID)
	require.Equal(t, "prod-1", byID.Digest.Groups[0].Items[0].ProductID)

	byDate, err := client.GetProductPriceWatchDigest(context.Background(), &GetProductPriceWatchDigestRequest{DigestDate: "2026-07-08"})
	require.NoError(t, err)
	require.Equal(t, "digest-2", byDate.Digest.DigestID)

	resp, err := client.GetProductPriceWatchDigest(context.Background(), &GetProductPriceWatchDigestRequest{})
	require.Error(t, err)
	require.Nil(t, resp)
}
