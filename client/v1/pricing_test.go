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
