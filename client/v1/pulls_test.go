package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
