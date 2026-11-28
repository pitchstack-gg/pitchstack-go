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

func TestClientPulls(t *testing.T) {
	t.Run("create pull", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/pulls", r.URL.Path)

			var payload CreatePullRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "prod-sealed-1", payload.SealedProductID)
			require.Equal(t, int32(2), payload.UnitsOpened)
			require.Len(t, payload.RarityTotals, 1)

			require.NoError(t, json.NewEncoder(w).Encode(CreatePullResponse{
				Pull: &Pull{ID: "pl-1", SealedProductID: "prod-sealed-1"},
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.CreatePull(context.Background(), &CreatePullRequest{
			SealedProductID: "prod-sealed-1",
			UnitsOpened:     2,
			RarityTotals:    []RarityCount{{Rarity: RarityMajestic, Quantity: 1}},
		})
		require.NoError(t, err)
		require.Equal(t, "pl-1", resp.Pull.ID)

		_, err = client.CreatePull(context.Background(), &CreatePullRequest{})
		require.Error(t, err)
	})

	t.Run("get pull", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/pulls/pl-1", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(GetPullResponse{Pull: &Pull{ID: "pl-1"}}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetPull(context.Background(), &GetPullRequest{PullID: "pl-1"})
		require.NoError(t, err)
		require.Equal(t, "pl-1", resp.Pull.ID)

		_, err = client.GetPull(context.Background(), &GetPullRequest{})
		require.Error(t, err)
	})

	t.Run("list pulls", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/pulls", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "WTR", q.Get("setId"))
			require.Equal(t, string(PullScopePack), q.Get("scope"))
			require.Equal(t, "10", q.Get("pageSize"))
			require.Equal(t, "tok-1", q.Get("nextToken"))
			require.NoError(t, json.NewEncoder(w).Encode(ListPullsResponse{Pulls: []Pull{{ID: "pl-1"}}}))
		}))
		t.Cleanup(server.Close)

		pageSize := int32(10)
		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListPulls(context.Background(), &ListPullsRequest{
			SetID:     "WTR",
			Scope:     PullScopePack,
			PageSize:  &pageSize,
			NextToken: "tok-1",
		})
		require.NoError(t, err)
		require.Len(t, resp.Pulls, 1)
	})

	t.Run("delete pull", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/v1/pulls/pl-1", r.URL.Path)
			w.Header().Set("X-Request-Id", "req-del")
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.DeletePull(context.Background(), &DeletePullRequest{PullID: "pl-1"})
		require.NoError(t, err)
		require.Equal(t, "req-del", resp.Metadata.RequestID)

		_, err = client.DeletePull(context.Background(), &DeletePullRequest{})
		require.Error(t, err)
	})

	t.Run("get pull stats", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/pulls:stats", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "WTR", q.Get("setId"))
			require.Equal(t, string(PullScopeBox), q.Get("scope"))
			require.NoError(t, json.NewEncoder(w).Encode(GetPullStatsResponse{
				ByScope: []PullStatsByScope{{Scope: PullScopeBox, PullsCount: 3}},
			}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetPullStats(context.Background(), &GetPullStatsRequest{SetID: "WTR", Scope: PullScopeBox})
		require.NoError(t, err)
		require.Len(t, resp.ByScope, 1)
	})
}

func TestClientPullsContractUpdates(t *testing.T) {
	t.Run("when price paid and collection set, then create body includes them and pull value decoded", func(t *testing.T) {
		pack := int32(2)
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, map[string]any{"amount": 4.99, "currency": "USD"}, payload["pricePaid"])
			require.Equal(t, "col-1", payload["collectionId"])
			cards := payload["pulledCards"].([]any)
			require.Equal(t, float64(2), cards[0].(map[string]any)["packNumber"])
		}, `{"pull":{"id":"pl-1","pricePaid":{"amount":4.99,"currency":"USD"},"value":{"currency":"USD","totalValue":12.5,"pricedQuantity":10,"unpricedQuantity":1},"collectionId":"col-1","collectionSyncStatus":"PULL_COLLECTION_SYNC_STATUS_PENDING"}}`)

		resp, err := client.CreatePull(context.Background(), &CreatePullRequest{
			SealedProductID: "prod-sealed-1",
			UnitsOpened:     1,
			PricePaid:       &CurrencyAmount{Amount: 4.99, Currency: "USD"},
			CollectionID:    "col-1",
			PulledCards:     []PulledCardInput{{ProductID: "prod-1", Quantity: 1, PackNumber: &pack}},
		})
		require.NoError(t, err)
		require.Equal(t, PullCollectionSyncStatusPending, resp.Pull.CollectionSyncStatus)
		require.InDelta(t, 12.5, resp.Pull.Value.TotalValue, 0.0001)
		require.Equal(t, "USD", resp.Pull.PricePaid.Currency)
	})

	t.Run("when remove from collection set, then delete sends query", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/v1/pulls/pl-1", r.URL.Path)
			require.Equal(t, "true", r.URL.Query().Get("removeFromCollection"))
		}, `{}`)
		_, err := client.DeletePull(context.Background(), &DeletePullRequest{PullID: "pl-1", RemoveFromCollection: true})
		require.NoError(t, err)
	})

	t.Run("when stats returned, then foiling and currency totals decoded", func(t *testing.T) {
		client := newStubAPIClient(t, nil, `{"byScope":[{"scope":"PULL_SCOPE_PACK","detailedPullsCount":3,"foilingTotals":[{"foiling":"RAINBOW","quantity":2}],"rarityFoilingTotals":[{"rarity":"MAJESTIC","foiling":"COLD","quantity":1}],"currencyTotals":[{"currency":"USD","totalSpent":15,"unpricedQuantity":2}]}]}`)
		resp, err := client.GetPullStats(context.Background(), nil)
		require.NoError(t, err)
		scope := resp.ByScope[0]
		require.Equal(t, int32(3), scope.DetailedPullsCount)
		require.Equal(t, FoilingRainbow, scope.FoilingTotals[0].Foiling)
		require.Equal(t, FoilingCold, scope.RarityFoilingTotals[0].Foiling)
		require.InDelta(t, 15.0, scope.CurrencyTotals[0].TotalSpent, 0.0001)
	})
}

func TestClientUpdatePull(t *testing.T) {
	t.Run("when fields provided, then PATCH body sent and details decoded", func(t *testing.T) {
		units := int32(3)
		collectionID := "col-2"
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodPatch, r.Method)
			require.Equal(t, "/v1/pulls/pl-1", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, float64(3), payload["unitsOpened"])
			require.Equal(t, true, payload["clearPricePaid"])
			require.Equal(t, "col-2", payload["collectionId"])
			require.Equal(t, "2026-11-01T00:00:00Z", payload["expectedUpdatedAt"])
			contents := payload["contents"].(map[string]any)
			require.Len(t, contents["rarityTotals"], 1)
			require.NotContains(t, payload, "pullId")
		}, `{"pull":{"id":"pl-1","unitsOpened":3},"pulledCards":[{"productId":"prod-1","quantity":1,"foiling":"GOLD","packNumber":1,"unitValue":0}]}`)

		expected := time.Date(2026, 11, 1, 0, 0, 0, 0, time.UTC)
		resp, err := client.UpdatePull(context.Background(), &UpdatePullRequest{
			PullID:            "pl-1",
			UnitsOpened:       &units,
			ClearPricePaid:    true,
			Contents:          &PullContents{RarityTotals: []RarityCount{{Rarity: RarityRare, Quantity: 1}}},
			CollectionID:      &collectionID,
			ExpectedUpdatedAt: &expected,
		})
		require.NoError(t, err)
		require.Equal(t, int32(3), resp.Pull.UnitsOpened)
		require.Equal(t, FoilingGold, resp.PulledCards[0].Foiling)
		require.NotNil(t, resp.PulledCards[0].UnitValue)
		require.Zero(t, *resp.PulledCards[0].UnitValue)
	})

	t.Run("when price paid and clear both set, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.UpdatePull(context.Background(), &UpdatePullRequest{PullID: "pl-1", PricePaid: &CurrencyAmount{Amount: 1}, ClearPricePaid: true})
		require.Error(t, err)
	})

	t.Run("when pull id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.UpdatePull(context.Background(), &UpdatePullRequest{})
		require.Error(t, err)
	})
}

func TestClientGetCommunityPullRates(t *testing.T) {
	t.Run("when filters provided, then query sent and rates decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/pulls:communityRates", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "WTR", q.Get("setId"))
			require.Equal(t, "prod-box", q.Get("sealedProductId"))
			require.Equal(t, "PULL_SCOPE_BOX", q.Get("scope"))
		}, `{"setId":"WTR","byScope":[{"scope":"PULL_SCOPE_BOX","contributorsCount":4,"withheld":true}],"thresholds":{"minContributors":5,"minUnitsOpened":10},"computedAt":"2026-11-01T00:00:00Z"}`)

		resp, err := client.GetCommunityPullRates(context.Background(), &GetCommunityPullRatesRequest{SetID: "WTR", SealedProductID: "prod-box", Scope: PullScopeBox})
		require.NoError(t, err)
		require.True(t, resp.ByScope[0].Withheld)
		require.Equal(t, int32(5), resp.Thresholds.MinContributors)
		require.NotNil(t, resp.ComputedAt)
	})

	t.Run("when set id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.GetCommunityPullRates(context.Background(), &GetCommunityPullRatesRequest{})
		require.Error(t, err)
	})
}

func TestClientPullSettings(t *testing.T) {
	t.Run("when settings fetched, then decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/pulls:settings", r.URL.Path)
		}, `{"settings":{"contributeToCommunityRates":true}}`)
		resp, err := client.GetPullSettings(context.Background())
		require.NoError(t, err)
		require.True(t, resp.Settings.ContributeToCommunityRates)
	})

	t.Run("when opting out, then PATCH sends explicit false", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodPatch, r.Method)
			require.Equal(t, "/v1/pulls:settings", r.URL.Path)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, false, payload["contributeToCommunityRates"])
		}, `{"settings":{"updatedAt":"2026-11-01T00:00:00Z"}}`)
		optIn := false
		resp, err := client.UpdatePullSettings(context.Background(), &UpdatePullSettingsRequest{ContributeToCommunityRates: &optIn})
		require.NoError(t, err)
		require.False(t, resp.Settings.ContributeToCommunityRates)
		require.NotNil(t, resp.Settings.UpdatedAt)
	})

	t.Run("when update request nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.UpdatePullSettings(context.Background(), nil)
		require.Error(t, err)
	})
}
