package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientReviewGroupJoinRequest(t *testing.T) {
	t.Run("when approving without a reason, then sends only the decision", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/groups/group-1/join-requests/jr-1:review", r.URL.Path)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "GROUP_JOIN_REQUEST_STATUS_APPROVED", payload["decision"])
			require.NotContains(t, payload, "declineReason")
			_, _ = w.Write([]byte(`{"joinRequest":{"requestId":"jr-1","status":"GROUP_JOIN_REQUEST_STATUS_APPROVED"}}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ReviewGroupJoinRequest(context.Background(), "group-1", "jr-1", GroupJoinRequestStatusApproved)
		require.NoError(t, err)
		require.Equal(t, GroupJoinRequestStatusApproved, resp.JoinRequest.Status)
	})

	t.Run("when the request id is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ReviewGroupJoinRequest(context.Background(), "group-1", " ", GroupJoinRequestStatusApproved)
		require.Error(t, err)
	})
}

func TestClientReviewGroupJoinRequestWithReason(t *testing.T) {
	t.Run("when declining with a reason, then sends the decline reason", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/groups/group-1/join-requests/jr-1:review", r.URL.Path)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "GROUP_JOIN_REQUEST_STATUS_DECLINED", payload["decision"])
			require.Equal(t, "group is full", payload["declineReason"])
			_, _ = w.Write([]byte(`{"joinRequest":{"requestId":"jr-1","status":"GROUP_JOIN_REQUEST_STATUS_DECLINED","declineReason":"group is full"}}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ReviewGroupJoinRequestWithReason(context.Background(), "group-1", "jr-1", GroupJoinRequestStatusDeclined, " group is full ")
		require.NoError(t, err)
		require.Equal(t, "group is full", resp.JoinRequest.DeclineReason)
	})
}

func TestClientListGroupReports(t *testing.T) {
	t.Run("when the server returns a string-encoded version, then decodes it", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			_, _ = w.Write([]byte(`{"reports":[{"reportId":"rep-1","version":"7"}]}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListGroupReports(context.Background(), nil, nil, "")
		require.NoError(t, err)
		require.Equal(t, int64(7), resp.Reports[0].Version)
	})
}

func TestClientPinGroupResource(t *testing.T) {
	t.Run("when pinning an event, then sends the event resource type", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/groups/group-1/pins", r.URL.Path)
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "GROUP_RESOURCE_TYPE_EVENT", payload["resourceType"])
			require.Equal(t, "event-1", payload["resourceId"])
			_, _ = w.Write([]byte(`{"pin":{"pinId":"pin-1","resourceType":"GROUP_RESOURCE_TYPE_EVENT"}}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.PinGroupResource(context.Background(), "group-1", GroupResourceTypeEvent, "event-1")
		require.NoError(t, err)
		require.Equal(t, GroupResourceTypeEvent, resp.Pin.ResourceType)
	})
}
