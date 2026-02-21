package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientCreateDeckVersionMatch(t *testing.T) {
	t.Run("when request valid, then match created", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodPost, r.Method)
			require.Equal(t, "/v1/deck_versions/dv-1/matches", r.URL.Path)
			require.Equal(t, "application/json", r.Header.Get("Content-Type"))

			var payload map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "dv-1", payload["deckVersionId"])
			require.Equal(t, string(DeckVersionMatchResultWin), payload["result"])

			resp := CreateDeckVersionMatchResponse{
				Match: &DeckVersionMatch{
					ID:            "match-1",
					DeckVersionID: "dv-1",
					Result:        DeckVersionMatchResultWin,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.CreateDeckVersionMatch(context.Background(), &CreateDeckVersionMatchRequest{
			DeckVersionID: " dv-1 ",
			Result:        DeckVersionMatchResultWin,
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Match)
		require.Equal(t, "match-1", resp.Match.ID)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when request invalid, then error returned", func(t *testing.T) {
		client := newTestClient(t)

		resp, err := client.CreateDeckVersionMatch(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)

		resp, err = client.CreateDeckVersionMatch(context.Background(), &CreateDeckVersionMatchRequest{Result: DeckVersionMatchResultWin})
		require.Error(t, err)
		require.Nil(t, resp)

		resp, err = client.CreateDeckVersionMatch(context.Background(), &CreateDeckVersionMatchRequest{DeckVersionID: "dv-1"})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetDeckVersionMatch(t *testing.T) {
	t.Run("when ids provided, then match fetched", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/deck_versions/dv-1/matches/match-1", r.URL.Path)

			resp := GetDeckVersionMatchResponse{
				Match: &DeckVersionMatch{
					ID:            "match-1",
					DeckVersionID: "dv-1",
					Result:        DeckVersionMatchResultLoss,
				},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetDeckVersionMatch(context.Background(), &GetDeckVersionMatchRequest{
			DeckVersionID: " dv-1 ",
			MatchID:       " match-1 ",
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Match)
		require.Equal(t, DeckVersionMatchResultLoss, resp.Match.Result)
	})

	t.Run("when ids missing, then error returned", func(t *testing.T) {
		client := newTestClient(t)

		resp, err := client.GetDeckVersionMatch(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)

		resp, err = client.GetDeckVersionMatch(context.Background(), &GetDeckVersionMatchRequest{MatchID: "match-1"})
		require.Error(t, err)
		require.Nil(t, resp)

		resp, err = client.GetDeckVersionMatch(context.Background(), &GetDeckVersionMatchRequest{DeckVersionID: "dv-1"})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientListDeckVersionMatches(t *testing.T) {
	t.Run("when request valid, then list returned", func(t *testing.T) {
		size := int32(10)
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/deck_versions/dv-1/matches", r.URL.Path)
			require.Equal(t, "10", r.URL.Query().Get("pageSize"))
			require.Equal(t, "tok-1", r.URL.Query().Get("nextToken"))

			resp := ListDeckVersionMatchesResponse{
				Matches: []DeckVersionMatch{
					{ID: "match-1", DeckVersionID: "dv-1", Result: DeckVersionMatchResultWin},
				},
				NextToken: "tok-2",
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListDeckVersionMatches(context.Background(), &ListDeckVersionMatchesRequest{
			DeckVersionID: "dv-1",
			PageSize:      &size,
			NextToken:     " tok-1 ",
		})
		require.NoError(t, err)
		require.Len(t, resp.Matches, 1)
		require.Equal(t, "tok-2", resp.NextToken)
	})

	t.Run("when request invalid, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.ListDeckVersionMatches(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)

		resp, err = client.ListDeckVersionMatches(context.Background(), &ListDeckVersionMatchesRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientDeleteDeckVersionMatch(t *testing.T) {
	t.Run("when request valid, then match deleted", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodDelete, r.Method)
			require.Equal(t, "/v1/deck_versions/dv-1/matches/match-1", r.URL.Path)
			_, _ = w.Write([]byte(`{}`))
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.DeleteDeckVersionMatch(context.Background(), &DeleteDeckVersionMatchRequest{
			DeckVersionID: "dv-1",
			MatchID:       "match-1",
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.Metadata.StatusCode)
	})

	t.Run("when request invalid, then error returned", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.DeleteDeckVersionMatch(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)

		resp, err = client.DeleteDeckVersionMatch(context.Background(), &DeleteDeckVersionMatchRequest{MatchID: "match-1"})
		require.Error(t, err)
		require.Nil(t, resp)

		resp, err = client.DeleteDeckVersionMatch(context.Background(), &DeleteDeckVersionMatchRequest{DeckVersionID: "dv-1"})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
