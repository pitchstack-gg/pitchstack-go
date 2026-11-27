package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListLivingLegendStandings(t *testing.T) {
	t.Run("when filters provided, then query encoded and standings decoded", func(t *testing.T) {
		pageSize := int32(25)
		includeRetired := false
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/living-legend/standings", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "cc", q.Get("format"))
			require.Equal(t, "2026-11-01", q.Get("snapshotDate"))
			require.Equal(t, "false", q.Get("includeRetired"))
			require.Equal(t, "25", q.Get("pageSize"))
			require.Equal(t, "tok", q.Get("nextToken"))
		}, `{"standings":[{"snapshotDate":"2026-11-01","format":"cc","heroName":"Dorinthea","heroCardId":"hero-1","points":850,"rank":1,"seasonPoints":0,"isLivingLegend":false,"source":"official","capturedAt":"2026-11-01T12:00:00Z"}],"nextToken":"next"}`)

		resp, err := client.ListLivingLegendStandings(context.Background(), &ListLivingLegendStandingsRequest{
			Format:         "cc",
			SnapshotDate:   "2026-11-01",
			IncludeRetired: &includeRetired,
			PageSize:       &pageSize,
			NextToken:      "tok",
		})
		require.NoError(t, err)
		require.Len(t, resp.Standings, 1)
		standing := resp.Standings[0]
		require.Equal(t, int32(850), standing.Points)
		require.NotNil(t, standing.Rank)
		require.Equal(t, int32(1), *standing.Rank)
		require.NotNil(t, standing.SeasonPoints)
		require.Zero(t, *standing.SeasonPoints)
		require.NotNil(t, standing.CapturedAt)
		require.Equal(t, "next", resp.NextToken)
	})

	t.Run("when request is nil, then no query is sent", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Empty(t, r.URL.RawQuery)
		}, `{}`)
		resp, err := client.ListLivingLegendStandings(context.Background(), nil)
		require.NoError(t, err)
		require.Empty(t, resp.Standings)
	})
}

func TestClientListHeroLivingLegendHistory(t *testing.T) {
	t.Run("when hero provided, then path and query encoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/living-legend/heroes/hero-1/history", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "blitz", q.Get("format"))
			require.Equal(t, "2026-01-01", q.Get("startDate"))
			require.Equal(t, "2026-06-30", q.Get("endDate"))
		}, `{"standings":[{"heroCardId":"hero-1","points":100}]}`)

		resp, err := client.ListHeroLivingLegendHistory(context.Background(), &ListHeroLivingLegendHistoryRequest{
			HeroCardID: "hero-1",
			Format:     "blitz",
			StartDate:  "2026-01-01",
			EndDate:    "2026-06-30",
		})
		require.NoError(t, err)
		require.Len(t, resp.Standings, 1)
		require.Nil(t, resp.Standings[0].Rank)
	})

	t.Run("when hero id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.ListHeroLivingLegendHistory(context.Background(), &ListHeroLivingLegendHistoryRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
