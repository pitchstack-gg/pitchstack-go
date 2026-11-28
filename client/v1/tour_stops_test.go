package client

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListTourStops(t *testing.T) {
	t.Run("when filters provided, then query encoded and tour stops decoded", func(t *testing.T) {
		pageSize := int32(10)
		year := int32(2026)
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/tour-stops", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "10", q.Get("pageSize"))
			require.Equal(t, "tok", q.Get("nextToken"))
			require.Equal(t, "2026", q.Get("year"))
		}, `{"tourStops":[{"tourStopId":"ts-1","title":"Pro Tour Example","startsOn":"2026-05-01","endsOn":"2026-05-03","status":"EVENT_STATUS_COMPLETED","timezone":"Europe/Amsterdam","locationImage":{"url":"https://cdn.example/img.jpg","kind":"EVENT_IMAGE_KIND_LOCATION","attribution":{"creator":"A","shareAlike":true}}}],"nextToken":"next"}`)

		resp, err := client.ListTourStops(context.Background(), &ListTourStopsRequest{PageSize: &pageSize, NextToken: "tok", Year: &year})
		require.NoError(t, err)
		require.Len(t, resp.TourStops, 1)
		stop := resp.TourStops[0]
		require.Equal(t, EventStatusCompleted, stop.Status)
		require.Equal(t, "2026-05-01", stop.StartsOn)
		require.Equal(t, EventImageKindLocation, stop.LocationImage.Kind)
		require.True(t, stop.LocationImage.Attribution.ShareAlike)
		require.Equal(t, "next", resp.NextToken)
	})

	t.Run("when request is nil, then no query is sent", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Empty(t, r.URL.RawQuery)
		}, `{}`)
		_, err := client.ListTourStops(context.Background(), nil)
		require.NoError(t, err)
	})
}

func TestClientGetTourStop(t *testing.T) {
	t.Run("when id provided, then tour stop and events decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/tour-stops/ts-1", r.URL.Path)
		}, `{"tourStop":{"tourStopId":"ts-1"},"events":[{"eventId":"evt-1","source":"EVENT_SOURCE_TOURNAMENT","tourStopId":"ts-1","tournament":{"catalogEventId":"cat-1","formats":["CC"],"entryFee":{"amount":50,"currency":"EUR"},"prizes":[{"place":"1st","worldTourPoints":0}]}}]}`)

		resp, err := client.GetTourStop(context.Background(), &GetTourStopRequest{TourStopID: "ts-1"})
		require.NoError(t, err)
		require.Equal(t, "ts-1", resp.TourStop.TourStopID)
		require.Len(t, resp.Events, 1)
		event := resp.Events[0]
		require.Equal(t, EventSourceTournament, event.Source)
		require.Equal(t, "cat-1", event.Tournament.CatalogEventID)
		require.InDelta(t, 50.0, event.Tournament.EntryFee.Amount, 0.0001)
		require.NotNil(t, event.Tournament.Prizes[0].WorldTourPoints)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetTourStop(context.Background(), &GetTourStopRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
