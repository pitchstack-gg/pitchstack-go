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
