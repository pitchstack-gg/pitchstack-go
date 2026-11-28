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

func TestClientEvents(t *testing.T) {
	start := time.Date(2026, 5, 1, 18, 0, 0, 0, time.UTC)
	pageSize := int32(10)
	lat := 47.61
	lng := -122.33
	radius := 25.0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/events":
			q := r.URL.Query()
			require.Equal(t, "10", q.Get("pageSize"))
			require.Equal(t, "EVENT_SOURCE_COMMUNITY", q.Get("source"))
			require.Equal(t, "armory", q.Get("query"))
			require.Equal(t, "CC", q.Get("format"))
			require.Equal(t, "US", q.Get("country"))
			require.Equal(t, "47.61", q.Get("latitude"))
			require.NoError(t, json.NewEncoder(w).Encode(ListEventsResponse{
				Events: []Event{{EventID: "evt-1", Title: "Armory"}},
			}))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/events/evt-1":
			require.NoError(t, json.NewEncoder(w).Encode(GetEventResponse{Event: &Event{EventID: "evt-1"}}))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/events":
			var payload CreateCommunityEventRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "Community Armory", payload.Title)
			require.NoError(t, json.NewEncoder(w).Encode(CreateCommunityEventResponse{Event: &Event{EventID: "evt-2"}}))
		case r.Method == http.MethodPatch && r.URL.Path == "/v1/events/evt-2":
			var payload UpdateCommunityEventRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.NotNil(t, payload.Title)
			require.NoError(t, json.NewEncoder(w).Encode(UpdateCommunityEventResponse{Event: &Event{EventID: "evt-2"}}))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/events/evt-2:cancel":
			require.NoError(t, json.NewEncoder(w).Encode(CancelCommunityEventResponse{Event: &Event{EventID: "evt-2", Status: EventStatusCancelled}}))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	listResp, err := client.ListEvents(context.Background(), &ListEventsRequest{
		PageSize:  &pageSize,
		Source:    EventSourceCommunity,
		Query:     "armory",
		Format:    "CC",
		Country:   "US",
		Latitude:  &lat,
		Longitude: &lng,
		RadiusKM:  &radius,
	})
	require.NoError(t, err)
	require.Len(t, listResp.Events, 1)

	getResp, err := client.GetEvent(context.Background(), &GetEventRequest{EventID: "evt-1"})
	require.NoError(t, err)
	require.Equal(t, "evt-1", getResp.Event.EventID)

	createResp, err := client.CreateCommunityEvent(context.Background(), &CreateCommunityEventRequest{Title: "Community Armory", StartsAt: &start})
	require.NoError(t, err)
	require.Equal(t, "evt-2", createResp.Event.EventID)

	title := "Updated Armory"
	updateResp, err := client.UpdateCommunityEvent(context.Background(), &UpdateCommunityEventRequest{EventID: "evt-2", Title: &title})
	require.NoError(t, err)
	require.Equal(t, "evt-2", updateResp.Event.EventID)

	cancelResp, err := client.CancelCommunityEvent(context.Background(), &CancelCommunityEventRequest{EventID: "evt-2"})
	require.NoError(t, err)
	require.Equal(t, EventStatusCancelled, cancelResp.Event.Status)
}

func TestClientEventStoresAndAdmin(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v1/event-stores":
			require.Equal(t, "Monday", r.URL.Query().Get("armoryDay"))
			require.Equal(t, "true", r.URL.Query().Get("onlineStore"))
			require.NoError(t, json.NewEncoder(w).Encode(ListStoresResponse{
				Stores: []EventStore{{StoreID: "store-1", Name: "LGS"}},
			}))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/event-stores/store-1":
			require.NoError(t, json.NewEncoder(w).Encode(GetStoreResponse{Store: &EventStore{StoreID: "store-1"}}))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/events/filters":
			require.NoError(t, json.NewEncoder(w).Encode(ListEventFiltersResponse{Filters: []EventFilterValue{{Kind: "format", Value: "CC"}}}))
		case r.Method == http.MethodGet && r.URL.Path == "/v1/admin/events/gem-locator-scans/run-1":
			require.NoError(t, json.NewEncoder(w).Encode(GetGemLocatorScanResponse{Run: &GemLocatorScan{RunID: "run-1"}}))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/admin/events/evt-1:hide":
			require.NoError(t, json.NewEncoder(w).Encode(HideEventResponse{Event: &Event{EventID: "evt-1", Status: EventStatusHidden}}))
		case r.Method == http.MethodPost && r.URL.Path == "/v1/admin/events/evt-1:unhide":
			require.NoError(t, json.NewEncoder(w).Encode(UnhideEventResponse{Event: &Event{EventID: "evt-1", Status: EventStatusPublished}}))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	online := true
	storesResp, err := client.ListStores(context.Background(), &ListStoresRequest{ArmoryDay: "Monday", OnlineStore: &online})
	require.NoError(t, err)
	require.Len(t, storesResp.Stores, 1)

	storeResp, err := client.GetStore(context.Background(), &GetStoreRequest{StoreRef: "store-1"})
	require.NoError(t, err)
	require.Equal(t, "store-1", storeResp.Store.StoreID)

	filtersResp, err := client.ListEventFilters(context.Background())
	require.NoError(t, err)
	require.Len(t, filtersResp.Filters, 1)

	scanResp, err := client.GetGemLocatorScan(context.Background(), &GetGemLocatorScanRequest{RunID: "run-1"})
	require.NoError(t, err)
	require.Equal(t, "run-1", scanResp.Run.RunID)

	hideResp, err := client.HideEvent(context.Background(), &HideEventRequest{EventID: "evt-1", Reason: "duplicate"})
	require.NoError(t, err)
	require.Equal(t, EventStatusHidden, hideResp.Event.Status)

	unhideResp, err := client.UnhideEvent(context.Background(), &UnhideEventRequest{EventID: "evt-1"})
	require.NoError(t, err)
	require.Equal(t, EventStatusPublished, unhideResp.Event.Status)
}

func TestClientEventsContractUpdates(t *testing.T) {
	t.Run("when order set, then list events sends orderBy and decodes tournament fields", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, "/v1/events", r.URL.Path)
			require.Equal(t, "LIST_EVENTS_ORDER_DISTANCE_ASC", r.URL.Query().Get("orderBy"))
			require.Equal(t, "EVENT_SOURCE_TOURNAMENT", r.URL.Query().Get("source"))
		}, `{"events":[{"eventId":"evt-1","status":"EVENT_STATUS_COMPLETED","tourStopId":"ts-1","registrationUrl":"https://reg.example","canonicalEventId":"evt-0","timezone":"America/Chicago","distanceKm":12.5,"image":{"url":"https://cdn.example/e.jpg","kind":"EVENT_IMAGE_KIND_PRODUCT_SET"},"tournament":{"catalogEventId":"cat-1","days":[{"date":"2026-05-01","startTimeLocal":"09:00"}],"attendance":{"count":300,"unit":"players"},"prizePool":{"amount":10000,"currency":"USD"},"banner":{"accentColor":"#fff"}}}]}`)

		resp, err := client.ListEvents(context.Background(), &ListEventsRequest{Source: EventSourceTournament, OrderBy: ListEventsOrderDistanceAsc})
		require.NoError(t, err)
		require.Len(t, resp.Events, 1)
		event := resp.Events[0]
		require.Equal(t, EventStatusCompleted, event.Status)
		require.Equal(t, "ts-1", event.TourStopID)
		require.Equal(t, "America/Chicago", event.Timezone)
		require.InDelta(t, 12.5, *event.DistanceKM, 0.0001)
		require.Equal(t, EventImageKindProductSet, event.Image.Kind)
		require.Equal(t, int32(300), event.Tournament.Attendance.Count)
		require.Equal(t, "09:00", event.Tournament.Days[0].StartTimeLocal)
		require.Equal(t, "USD", event.Tournament.PrizePool.Currency)
	})

	t.Run("when order unspecified, then orderBy omitted", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.False(t, r.URL.Query().Has("orderBy"))
		}, `{}`)
		_, err := client.ListEvents(context.Background(), &ListEventsRequest{OrderBy: ListEventsOrderUnspecified})
		require.NoError(t, err)
	})

	t.Run("when store order set, then list stores sends orderBy and decodes distance", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, "/v1/event-stores", r.URL.Path)
			require.Equal(t, "LIST_STORES_ORDER_NAME_ASC", r.URL.Query().Get("orderBy"))
		}, `{"stores":[{"storeId":"store-1","distanceKm":3.2}]}`)

		resp, err := client.ListStores(context.Background(), &ListStoresRequest{OrderBy: ListStoresOrderNameAsc})
		require.NoError(t, err)
		require.InDelta(t, 3.2, *resp.Stores[0].DistanceKM, 0.0001)
	})

	t.Run("when timezone and update mask set, then body includes them", func(t *testing.T) {
		start := time.Date(2026, 5, 1, 18, 0, 0, 0, time.UTC)
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "Europe/London", payload["timezone"])
			if r.Method == http.MethodPatch {
				require.Equal(t, "title,timezone", payload["updateMask"])
			}
		}, `{"event":{"eventId":"evt-2"}}`)

		_, err := client.CreateCommunityEvent(context.Background(), &CreateCommunityEventRequest{Title: "Armory", StartsAt: &start, Timezone: "Europe/London"})
		require.NoError(t, err)

		title := "Armory"
		tz := "Europe/London"
		_, err = client.UpdateCommunityEvent(context.Background(), &UpdateCommunityEventRequest{EventID: "evt-2", Title: &title, Timezone: &tz, UpdateMask: "title,timezone"})
		require.NoError(t, err)
	})
}
