package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientSearchProfileCollectionItems(t *testing.T) {
	t.Run("when given filters, then posts them and decodes results", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/users/user-1/collection_items:search", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))
			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "fyendal", payload["searchTerm"])
			require.Equal(t, []any{"card-1"}, payload["cardIds"])
			require.Equal(t, []any{"BINDER"}, payload["collectionTypes"])
			require.Equal(t, float64(25), payload["pageSize"])
			require.Equal(t, "tok", payload["nextToken"])
			require.NotContains(t, payload, "userId")
			_, _ = w.Write([]byte(`{"items":[{"itemId":"item-1","productId":"prod-1","quantity":2,"condition":"NEAR_MINT","sourceCollection":{"id":"col-1","name":"Main","collectionType":"BINDER"},"catalogMatches":[{"cardId":"card-1","printingId":"pr-1","side":"PROFILE_COLLECTION_ITEM_SEARCH_SIDE_FRONT"}]}],"nextToken":"next"}`))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		pageSize := int32(25)
		resp, err := client.SearchProfileCollectionItems(context.Background(), &SearchProfileCollectionItemsRequest{
			UserID:          " user-1 ",
			SearchTerm:      " fyendal ",
			CardIDs:         []string{"card-1"},
			CollectionTypes: []CollectionType{CollectionTypeBinder},
			PageSize:        &pageSize,
			NextToken:       "tok",
		})
		require.NoError(t, err)
		require.Equal(t, "next", resp.NextToken)
		item := resp.Items[0]
		require.Equal(t, int32(2), item.Quantity)
		require.Equal(t, ConditionNearMint, item.Condition)
		require.Equal(t, CollectionTypeBinder, item.SourceCollection.CollectionType)
		require.Equal(t, ProfileCollectionItemSearchSideFront, item.CatalogMatches[0].Side)
	})

	t.Run("when the user id is missing, then returns an error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.SearchProfileCollectionItems(context.Background(), &SearchProfileCollectionItemsRequest{})
		require.Error(t, err)
		_, err = client.SearchProfileCollectionItems(context.Background(), nil)
		require.Error(t, err)
	})
}
