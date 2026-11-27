package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientListGroupAuditLog(t *testing.T) {
	t.Run("when given filters, then sends them as query parameters and decodes entries", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/groups/group-1/audit-log", r.URL.Path)
			query := r.URL.Query()
			require.Equal(t, "5", query.Get("pageSize"))
			require.Equal(t, "tok", query.Get("nextToken"))
			require.Equal(t, []string{"GROUP_AUDIT_ACTION_MEMBER_BANNED", "GROUP_AUDIT_ACTION_MEMBER_ADDED"}, query["actions"])
			require.Equal(t, "user-2", query.Get("targetUserId"))
			_, _ = w.Write([]byte(`{
				"entries": [{
					"entryId": "e-1",
					"groupId": "group-1",
					"action": "GROUP_AUDIT_ACTION_MEMBER_ROLE_CHANGED",
					"actor": {"type": "GROUP_AUDIT_ACTOR_TYPE_USER", "userId": "user-1"},
					"targetUserId": "user-2",
					"targetType": "GROUP_AUDIT_TARGET_TYPE_MEMBER",
					"detail": {
						"previousRole": "GROUP_ROLE_MEMBER",
						"newRole": "GROUP_ROLE_ADMIN",
						"fieldChanges": [{"field": "join_policy", "previousValue": "a", "newValue": "b"}],
						"membershipSource": "GROUP_AUDIT_MEMBERSHIP_SOURCE_INVITE",
						"cause": "GROUP_AUDIT_CAUSE_INVITE_RESENT",
						"linkCount": 3,
						"resourceType": "GROUP_RESOURCE_TYPE_EVENT"
					},
					"occurredAt": "2026-01-02T03:04:05Z"
				}],
				"nextToken": "next"
			}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		pageSize := int32(5)
		resp, err := client.ListGroupAuditLog(context.Background(), &ListGroupAuditLogRequest{
			GroupID:      "group-1",
			PageSize:     &pageSize,
			NextToken:    "tok",
			Actions:      []GroupAuditAction{GroupAuditActionMemberBanned, GroupAuditActionMemberAdded},
			TargetUserID: "user-2",
		})
		require.NoError(t, err)
		require.Equal(t, "next", resp.NextToken)
		require.Len(t, resp.Entries, 1)
		entry := resp.Entries[0]
		require.Equal(t, GroupAuditActionMemberRoleChanged, entry.Action)
		require.Equal(t, GroupAuditActorTypeUser, entry.Actor.Type)
		require.Equal(t, GroupAuditTargetTypeMember, entry.TargetType)
		require.Equal(t, GroupRoleAdmin, *entry.Detail.NewRole)
		require.Equal(t, GroupRoleMember, *entry.Detail.PreviousRole)
		require.Equal(t, "join_policy", entry.Detail.FieldChanges[0].Field)
		require.Equal(t, GroupAuditMembershipSourceInvite, entry.Detail.MembershipSource)
		require.Equal(t, GroupAuditCauseInviteResent, entry.Detail.Cause)
		require.Equal(t, int32(3), *entry.Detail.LinkCount)
		require.Equal(t, GroupResourceTypeEvent, entry.Detail.ResourceType)
		require.NotNil(t, entry.OccurredAt)
	})

	t.Run("when no filters are given, then sends no query string", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/v1/groups/group-1/audit-log", r.URL.Path)
			require.Empty(t, r.URL.RawQuery)
			_, _ = w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		_, err := client.ListGroupAuditLog(context.Background(), &ListGroupAuditLogRequest{GroupID: "group-1"})
		require.NoError(t, err)
	})

	t.Run("when the group is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ListGroupAuditLog(context.Background(), &ListGroupAuditLogRequest{})
		require.Error(t, err)
		_, err = client.ListGroupAuditLog(context.Background(), nil)
		require.Error(t, err)
	})
}
