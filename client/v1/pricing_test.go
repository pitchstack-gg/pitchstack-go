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

func TestClientGetPhysicalCardPrice(t *testing.T) {
	t.Run("when card id provided, then price entry returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v1/prices/phys-1", r.URL.Path)
			require.Equal(t, "TCGPlayer", r.URL.Query().Get("source"))
			require.Equal(t, "USD", r.URL.Query().Get("currency"))

			resp := GetPhysicalCardPriceResponse{
				Entry: &PriceEntry{
					EntryID:        "entry-1",
					PhysicalCardID: "phys-1",
					Price:          12.34,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetPhysicalCardPrice(context.Background(), &GetPhysicalCardPriceRequest{
			PhysicalCardID: "phys-1",
			Source:         "TCGPlayer",
			Currency:       "USD",
		})
		require.NoError(t, err)
		require.Equal(t, "phys-1", resp.Entry.PhysicalCardID)
		require.Equal(t, 12.34, resp.Entry.Price)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetPhysicalCardPrice(context.Background(), &GetPhysicalCardPriceRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetPhysicalCardPriceHistory(t *testing.T) {
	t.Run("when filters provided, then history returned", func(t *testing.T) {
		limit := int32(50)
		now := time.Now().UTC().Round(time.Second)

		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v1/prices/phys-1/history", r.URL.Path)
			require.Equal(t, "2024-01-01", r.URL.Query().Get("startDate"))
			require.Equal(t, "2024-02-01", r.URL.Query().Get("endDate"))
			require.Equal(t, "50", r.URL.Query().Get("limit"))
			require.Equal(t, "TCGPlayer", r.URL.Query().Get("source"))

			resp := GetPhysicalCardPriceHistoryResponse{
				PhysicalCardID: "phys-1",
				Source:         "TCGPlayer",
				StartTime:      &now,
				EndTime:        &now,
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
		resp, err := client.GetPhysicalCardPriceHistory(context.Background(), &GetPhysicalCardPriceHistoryRequest{
			PhysicalCardID: "phys-1",
			StartDate:      "2024-01-01",
			EndDate:        "2024-02-01",
			Limit:          &limit,
			Source:         "TCGPlayer",
		})
		require.NoError(t, err)
		require.Equal(t, "phys-1", resp.PhysicalCardID)
		require.Len(t, resp.Entries, 1)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetPhysicalCardPriceHistory(context.Background(), &GetPhysicalCardPriceHistoryRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetBulkPhysicalCardPrices(t *testing.T) {
	t.Run("when ids provided, then prices returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			var body GetBulkPhysicalCardPricesRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.ElementsMatch(t, []string{"phys-1", "phys-2"}, body.PhysicalCardIDs)
			require.Equal(t, "TCGPlayer", body.Source)

			resp := GetBulkPhysicalCardPricesResponse{
				Prices: []PriceEntry{
					{EntryID: "entry-1", PhysicalCardID: "phys-1"},
				},
				NotFound: []string{"phys-3"},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetBulkPhysicalCardPrices(context.Background(), &GetBulkPhysicalCardPricesRequest{
			PhysicalCardIDs: []string{"phys-1", "phys-2"},
			Source:          "TCGPlayer",
		})
		require.NoError(t, err)
		require.Len(t, resp.Prices, 1)
		require.Equal(t, []string{"phys-3"}, resp.NotFound)
	})

	t.Run("when ids missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetBulkPhysicalCardPrices(context.Background(), &GetBulkPhysicalCardPricesRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetBulkPhysicalCardPrices(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
