package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientModerationReports(t *testing.T) {
	tests := []struct {
		name       string
		path       string
		call       func(*Client) (*SubmitModerationReportResponse, error)
		assertBody func(*testing.T, map[string]any)
	}{
		{
			name: "profile", path: "/v1/users/user%2Fone/profile:report",
			call: func(client *Client) (*SubmitModerationReportResponse, error) {
				return client.ReportProfile(context.Background(), &ReportProfileRequest{UserID: "user/one", Reason: ModerationReportReasonImpersonation, Details: " copied me "})
			},
			assertBody: func(t *testing.T, body map[string]any) {
				require.Equal(t, string(ModerationReportReasonImpersonation), body["reason"])
				require.Equal(t, "copied me", body["details"])
			},
		},
		{
			name: "deck subresource", path: "/v1/decks/deck-1:report",
			call: func(client *Client) (*SubmitModerationReportResponse, error) {
				return client.ReportDeckContent(context.Background(), &ReportDeckContentRequest{DeckID: "deck-1", TargetType: DeckReportTargetTypeMatch, TargetID: "match-1", Reason: ModerationReportReasonHarassmentOrHate})
			},
			assertBody: func(t *testing.T, body map[string]any) {
				require.Equal(t, string(DeckReportTargetTypeMatch), body["targetType"])
				require.Equal(t, "match-1", body["targetId"])
			},
		},
		{
			name: "collection", path: "/v1/collections/collection-1:report",
			call: func(client *Client) (*SubmitModerationReportResponse, error) {
				return client.ReportCollection(context.Background(), &ReportCollectionRequest{CollectionID: "collection-1", Reason: ModerationReportReasonSpamOrScam})
			},
			assertBody: func(t *testing.T, body map[string]any) {
				require.Equal(t, string(ModerationReportReasonSpamOrScam), body["reason"])
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, test.path, r.URL.EscapedPath())
				var body map[string]any
				require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
				test.assertBody(t, body)
				w.Header().Set("X-Request-Id", "req-report")
				_, _ = w.Write([]byte(`{"receipt":{"caseId":"case-1","duplicate":false}}`))
			}))
			t.Cleanup(server.Close)
			client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
			response, err := test.call(client)
			require.NoError(t, err)
			require.Equal(t, "case-1", response.Receipt.CaseID)
			require.Equal(t, "req-report", response.Metadata.RequestID)
		})
	}
}

func TestClientModerationReportValidation(t *testing.T) {
	client := newTestClient(t)
	_, err := client.ReportProfile(context.Background(), &ReportProfileRequest{UserID: "user-1"})
	require.ErrorContains(t, err, "reason")
	_, err = client.ReportDeckContent(context.Background(), &ReportDeckContentRequest{DeckID: "deck-1", TargetType: DeckReportTargetTypeMatch, Reason: ModerationReportReasonOther})
	require.ErrorContains(t, err, "targetID")
	_, err = client.ReportCollection(context.Background(), nil)
	require.Error(t, err)
}

func TestSharedGroupReportReason(t *testing.T) {
	reason, ok := sharedGroupReportReason(ModerationReportReasonViolence)
	require.True(t, ok)
	require.Equal(t, GroupReportReasonViolence, reason)
	_, ok = sharedGroupReportReason(ModerationReportReasonUnspecified)
	require.False(t, ok)
}
