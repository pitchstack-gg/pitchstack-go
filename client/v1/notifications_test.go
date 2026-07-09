package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientCreateMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/notifications.v1.NotificationsProducerService/CreateMessage", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload CreateMessageRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "u-11111111-1111-1111-1111-111111111111", payload.TargetUserID)
		require.Equal(t, "contract-tests", payload.Source)
		require.Equal(t, "idem-1", payload.IdempotencyKey)
		require.Equal(t, "contract", payload.Category)
		require.Equal(t, "message title", payload.Title)

		require.NoError(t, json.NewEncoder(w).Encode(CreateMessageResponse{
			MessageID: "n-11111111-1111-1111-1111-111111111111",
			Created:   true,
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.CreateMessage(context.Background(), &CreateMessageRequest{
		TargetUserID:   "u-11111111-1111-1111-1111-111111111111",
		Source:         "contract-tests",
		IdempotencyKey: "idem-1",
		Category:       "contract",
		Title:          "message title",
	})
	require.NoError(t, err)
	require.Equal(t, "n-11111111-1111-1111-1111-111111111111", resp.MessageID)
	require.True(t, resp.Created)
}

func TestClientPushDevicesAndPreferences(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/notifications/devices":
			var payload RegisterPushDeviceRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "device-1", payload.DeviceID)
			require.Equal(t, "ios", payload.Platform)
			require.Equal(t, "expo-token", payload.ExpoPushToken)
			require.NoError(t, json.NewEncoder(w).Encode(RegisterPushDeviceResponse{
				Device: &PushDevice{DeviceID: "device-1", Active: true},
			}))
		case r.Method == http.MethodDelete && r.URL.Path == "/v1/notifications/devices/device-1":
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodGet && r.URL.Path == "/v1/notifications/preferences":
			require.NoError(t, json.NewEncoder(w).Encode(GetNotificationPreferencesResponse{
				Preferences: []NotificationPreference{{Category: "social", PushEnabled: true}},
			}))
		case r.Method == http.MethodPut && r.URL.Path == "/v1/notifications/preferences":
			var payload UpdateNotificationPreferencesRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Len(t, payload.Preferences, 1)
			require.Equal(t, "social", payload.Preferences[0].Category)
			require.NoError(t, json.NewEncoder(w).Encode(UpdateNotificationPreferencesResponse{
				Preferences: payload.Preferences,
			}))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	registerResp, err := client.RegisterPushDevice(context.Background(), &RegisterPushDeviceRequest{
		DeviceID:      "device-1",
		Platform:      "ios",
		ExpoPushToken: "expo-token",
	})
	require.NoError(t, err)
	require.True(t, registerResp.Device.Active)

	unregisterResp, err := client.UnregisterPushDevice(context.Background(), &UnregisterPushDeviceRequest{DeviceID: "device-1"})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, unregisterResp.Metadata.StatusCode)

	getResp, err := client.GetNotificationPreferences(context.Background())
	require.NoError(t, err)
	require.True(t, getResp.Preferences[0].PushEnabled)

	updateResp, err := client.UpdateNotificationPreferences(context.Background(), &UpdateNotificationPreferencesRequest{
		Preferences: []NotificationPreference{{Category: "social", PushEnabled: true}},
	})
	require.NoError(t, err)
	require.Equal(t, "social", updateResp.Preferences[0].Category)
}

func TestClientUnsubscribeEmail(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/notifications/email:unsubscribe", r.URL.Path)
		require.Equal(t, "unsubscribe-token", r.URL.Query().Get("token"))
		w.Header().Set("X-Request-Id", "req-unsubscribe")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.UnsubscribeEmail(context.Background(), &UnsubscribeEmailRequest{Token: " unsubscribe-token "})
	require.NoError(t, err)
	require.Equal(t, "req-unsubscribe", resp.Metadata.RequestID)

	resp, err = client.UnsubscribeEmail(context.Background(), &UnsubscribeEmailRequest{})
	require.Error(t, err)
	require.Nil(t, resp)
}

func TestClientNotificationTopicSubscriptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			require.Equal(t, "/v1/notifications/topic-subscriptions", r.URL.Path)
			require.Equal(t, "news", r.URL.Query().Get("category"))
			require.Equal(t, "source", r.URL.Query().Get("topicType"))
			require.NoError(t, json.NewEncoder(w).Encode(GetNotificationTopicSubscriptionsResponse{
				TopicIDs: []string{"source-1"},
				Subscriptions: []NotificationTopicSubscription{{
					Category:  "news",
					TopicType: "source",
					TopicID:   "source-1",
				}},
			}))
		case http.MethodPut:
			require.Equal(t, "/v1/notifications/topic-subscriptions", r.URL.Path)
			var payload UpdateNotificationTopicSubscriptionsRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "news", payload.Category)
			require.Equal(t, []string{"source-1"}, payload.TopicIDs)
			require.NoError(t, json.NewEncoder(w).Encode(UpdateNotificationTopicSubscriptionsResponse{
				TopicIDs: payload.TopicIDs,
			}))
		default:
			t.Fatalf("unexpected method: %s", r.Method)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	getResp, err := client.GetNotificationTopicSubscriptions(context.Background(), &GetNotificationTopicSubscriptionsRequest{
		Category:  "news",
		TopicType: "source",
	})
	require.NoError(t, err)
	require.Equal(t, []string{"source-1"}, getResp.TopicIDs)

	updateResp, err := client.UpdateNotificationTopicSubscriptions(context.Background(), &UpdateNotificationTopicSubscriptionsRequest{
		Category:  "news",
		TopicType: "source",
		TopicIDs:  []string{"source-1"},
	})
	require.NoError(t, err)
	require.Equal(t, []string{"source-1"}, updateResp.TopicIDs)
}

func TestClientListInbox(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/notifications", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "true", query.Get("unreadOnly"))
		require.Equal(t, "true", query.Get("includeArchived"))
		require.Equal(t, "false", query.Get("includeExpired"))
		require.Equal(t, []string{"system", "social"}, query["categories"])
		require.Equal(t, "25", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(ListInboxResponse{
			Notifications: []Notification{{MessageID: "msg-1"}},
			NextToken:     "next",
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	unreadOnly := true
	includeArchived := true
	includeExpired := false
	pageSize := int32(25)
	resp, err := client.ListInbox(context.Background(), &ListInboxRequest{
		UnreadOnly:      &unreadOnly,
		IncludeArchived: &includeArchived,
		IncludeExpired:  &includeExpired,
		Categories:      []string{" system ", "social"},
		PageSize:        &pageSize,
		NextToken:       "token",
	})
	require.NoError(t, err)
	require.Equal(t, "msg-1", resp.Notifications[0].MessageID)
	require.Equal(t, "next", resp.NextToken)
}

func TestClientGetMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/notifications/msg-1", r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(GetMessageResponse{
			Notification: &Notification{MessageID: "msg-1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetMessage(context.Background(), &GetMessageRequest{MessageID: "msg-1"})
	require.NoError(t, err)
	require.Equal(t, "msg-1", resp.Notification.MessageID)
}

func TestClientDeleteMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/notifications/msg-1", r.URL.Path)
		w.Header().Set("X-Request-Id", "req-delete")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.DeleteMessage(context.Background(), &DeleteMessageRequest{MessageID: "msg-1"})
	require.NoError(t, err)
	require.Equal(t, "req-delete", resp.Metadata.RequestID)
}

func TestClientMarkRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/notifications/msg-1:markRead", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var payload map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Len(t, payload, 0)
		w.Header().Set("X-Request-Id", "req-mark")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.MarkRead(context.Background(), &MarkReadRequest{MessageID: "msg-1"})
	require.NoError(t, err)
	require.Equal(t, "req-mark", resp.Metadata.RequestID)
}

func TestClientMarkAllRead(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/notifications:markAllRead", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var payload map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Len(t, payload, 0)
		_, _ = w.Write([]byte(`{"markedReadCount":"3"}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.MarkAllRead(context.Background(), nil)
	require.NoError(t, err)
	require.Equal(t, int64(3), resp.MarkedReadCount)
}

func TestClientArchiveMessage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/notifications/msg-1:archive", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var payload map[string]any
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Len(t, payload, 0)
		w.Header().Set("X-Request-Id", "req-archive")
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ArchiveMessage(context.Background(), &ArchiveMessageRequest{MessageID: "msg-1"})
	require.NoError(t, err)
	require.Equal(t, "req-archive", resp.Metadata.RequestID)
}
