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

func TestClientListDecks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/decks", r.URL.Path)

		query := r.URL.Query()
		require.Equal(t, "DECK_LIST_SCOPE_SHARED", query.Get("scope"))
		require.Equal(t, "user-1", query.Get("userId"))
		require.Equal(t, "hero-1", query.Get("heroId"))
		require.Equal(t, "cc", query.Get("format"))
		require.Equal(t, "aggro", query.Get("name"))
		require.Equal(t, "25", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(ListDecksResponse{
			Decks: []Deck{{ID: "deck-1"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	pageSize := int32(25)
	resp, err := client.ListDecks(context.Background(), &ListDecksRequest{
		Scope:     DeckListScopeShared,
		UserID:    "user-1",
		HeroID:    "hero-1",
		Format:    "cc",
		Name:      "aggro",
		PageSize:  &pageSize,
		NextToken: "token",
	})
	require.NoError(t, err)
	require.Len(t, resp.Decks, 1)
}

func TestClientCreateDeck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/decks", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload CreateDeckRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "Aggro", payload.Name)
		require.Equal(t, "hero-1", payload.HeroID)
		require.Equal(t, "cc", payload.Format)
		require.Equal(t, VisibilityLevelShared, payload.Visibility)

		require.NoError(t, json.NewEncoder(w).Encode(CreateDeckResponse{
			Deck: &Deck{ID: "deck-1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.CreateDeck(context.Background(), &CreateDeckRequest{
		Name:       "Aggro",
		HeroID:     "hero-1",
		Format:     "cc",
		Visibility: VisibilityLevelShared,
	})
	require.NoError(t, err)
	require.Equal(t, "deck-1", resp.Deck.ID)
}

func TestClientSearchDecks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/decks/search", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "hero-1", query.Get("heroId"))
		require.Equal(t, "cc", query.Get("format"))
		require.Equal(t, "aggro", query.Get("searchTerm"))
		require.Equal(t, "10", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(SearchDecksResponse{
			Decks: []Deck{{ID: "deck-1"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	pageSize := int32(10)
	resp, err := client.SearchDecks(context.Background(), &SearchDecksRequest{
		HeroID:     "hero-1",
		Format:     "cc",
		SearchTerm: "aggro",
		PageSize:   &pageSize,
		NextToken:  "token",
	})
	require.NoError(t, err)
	require.Len(t, resp.Decks, 1)
}

func TestClientGetDeck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/decks/deck-1", r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(GetDeckResponse{
			Deck: &Deck{ID: "deck-1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetDeck(context.Background(), &GetDeckRequest{DeckID: "deck-1"})
	require.NoError(t, err)
	require.Equal(t, "deck-1", resp.Deck.ID)

	resp, err = client.GetDeck(context.Background(), &GetDeckRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientDeleteDeck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/decks/deck-1", r.URL.Path)
		w.Header().Set("X-Request-Id", "req-delete")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.DeleteDeck(context.Background(), &DeleteDeckRequest{DeckID: "deck-1"})
	require.NoError(t, err)
	require.Equal(t, "req-delete", resp.Metadata.RequestID)
}

func TestClientUpdateDeck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/decks/deck-1", r.URL.Path)

		var payload struct {
			Name       string `json:"name"`
			Visibility string `json:"visibility"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "Updated", payload.Name)
		require.Equal(t, string(VisibilityLevelPublic), payload.Visibility)

		require.NoError(t, json.NewEncoder(w).Encode(UpdateDeckResponse{
			Deck: &Deck{ID: "deck-1", Name: "Updated"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	name := "Updated"
	visibility := VisibilityLevelPublic
	resp, err := client.UpdateDeck(context.Background(), &UpdateDeckRequest{
		DeckID:     "deck-1",
		Name:       &name,
		Visibility: &visibility,
	})
	require.NoError(t, err)
	require.Equal(t, "Updated", resp.Deck.Name)
}

func TestClientStarAndUnstarDeck(t *testing.T) {
	var starCalled, unstarCalled bool

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			starCalled = true
			require.Equal(t, "/v1/decks/deck-1/stars", r.URL.Path)
		case http.MethodDelete:
			unstarCalled = true
			require.Equal(t, "/v1/decks/deck-1/stars", r.URL.Path)
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

	_, err := client.StarDeck(context.Background(), &StarDeckRequest{DeckID: "deck-1"})
	require.NoError(t, err)
	require.True(t, starCalled)

	_, err = client.UnstarDeck(context.Background(), &UnstarDeckRequest{DeckID: "deck-1"})
	require.NoError(t, err)
	require.True(t, unstarCalled)
}

func TestClientListDeckVersions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/decks/deck-1/versions", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "15", query.Get("pageSize"))
		require.Equal(t, "token", query.Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(ListDeckVersionsResponse{
			Versions: []string{"v1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	pageSize := int32(15)
	resp, err := client.ListDeckVersions(context.Background(), &ListDeckVersionsRequest{
		DeckID:    "deck-1",
		PageSize:  &pageSize,
		NextToken: "token",
	})
	require.NoError(t, err)
	require.Equal(t, []string{"v1"}, resp.Versions)
}

func TestClientCreateDeckVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/decks/deck-1/versions", r.URL.Path)

		var payload struct {
			Version     string `json:"version"`
			ImageURL    string `json:"imageUrl"`
			Description string `json:"description"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "v2", payload.Version)
		require.Equal(t, "https://image", payload.ImageURL)
		require.Equal(t, "desc", payload.Description)

		require.NoError(t, json.NewEncoder(w).Encode(CreateDeckVersionResponse{
			DeckVersion: &DeckVersion{Version: "v2"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.CreateDeckVersion(context.Background(), &CreateDeckVersionRequest{
		DeckID:      "deck-1",
		Version:     "v2",
		ImageURL:    "https://image",
		Description: "desc",
	})
	require.NoError(t, err)
	require.Equal(t, "v2", resp.DeckVersion.Version)
}

func TestClientGetDeckVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/decks/deck-1/versions/v1", r.URL.Path)
		require.Equal(t, "token", r.URL.Query().Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(GetDeckVersionResponse{
			DeckVersion: &DeckVersion{Version: "v1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetDeckVersion(context.Background(), &GetDeckVersionRequest{
		DeckID:    "deck-1",
		Version:   "v1",
		NextToken: "token",
	})
	require.NoError(t, err)
	require.Equal(t, "v1", resp.DeckVersion.Version)

	resp, err = client.GetDeckVersion(context.Background(), &GetDeckVersionRequest{DeckID: "deck-1"})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientDeleteDeckVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/decks/deck-1/versions/v1", r.URL.Path)
		w.Header().Set("X-Request-Id", "req-delete-version")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.DeleteDeckVersion(context.Background(), &DeleteDeckVersionRequest{
		DeckID:  "deck-1",
		Version: "v1",
	})
	require.NoError(t, err)
	require.Equal(t, "req-delete-version", resp.Metadata.RequestID)
}

func TestClientUpdateDeckVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/decks/deck-1/versions/v1", r.URL.Path)

		var payload struct {
			ImageURL    string `json:"imageUrl"`
			Description string `json:"description"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "https://image", payload.ImageURL)
		require.Equal(t, "desc", payload.Description)

		require.NoError(t, json.NewEncoder(w).Encode(UpdateDeckVersionResponse{
			DeckVersion: &DeckVersion{ImageURL: "https://image"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	imageURL := "https://image"
	desc := "desc"
	resp, err := client.UpdateDeckVersion(context.Background(), &UpdateDeckVersionRequest{
		DeckID:      "deck-1",
		Version:     "v1",
		ImageURL:    &imageURL,
		Description: &desc,
	})
	require.NoError(t, err)
	require.Equal(t, "https://image", resp.DeckVersion.ImageURL)
}

func TestClientListDeckVersionCards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/decks/deck-1/versions/v1/cards", r.URL.Path)

		require.NoError(t, json.NewEncoder(w).Encode(ListDeckVersionCardsResponse{
			MainboardCards:  []DeckCard{{CardID: "card-1", Quantity: 3}},
			SideboardCards:  []DeckCard{{CardID: "card-2", Quantity: 2}},
			MaybeboardCards: []DeckCard{{CardID: "card-3", Quantity: 1}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ListDeckVersionCards(context.Background(), &ListDeckVersionCardsRequest{
		DeckID:  "deck-1",
		Version: "v1",
	})
	require.NoError(t, err)
	require.Len(t, resp.MainboardCards, 1)
	require.Len(t, resp.SideboardCards, 1)
	require.Len(t, resp.MaybeboardCards, 1)
}

func TestClientModifyDeckVersionCard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/decks/deck-1/versions/v1/cards", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload struct {
			CardID   string `json:"cardId"`
			Board    string `json:"boardType"`
			Quantity int32  `json:"quantity"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "card-1", payload.CardID)
		require.Equal(t, string(BoardTypeMainboard), payload.Board)
		require.Equal(t, int32(3), payload.Quantity)

		require.NoError(t, json.NewEncoder(w).Encode(ModifyDeckVersionCardResponse{
			Card:  &DeckCard{CardID: "card-1", Quantity: 3},
			Board: BoardTypeMainboard,
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ModifyDeckVersionCard(context.Background(), &ModifyDeckVersionCardRequest{
		DeckID:   "deck-1",
		Version:  "v1",
		CardID:   "card-1",
		Board:    BoardTypeMainboard,
		Quantity: 3,
	})
	require.NoError(t, err)
	require.Equal(t, "card-1", resp.Card.CardID)

	resp, err = client.ModifyDeckVersionCard(context.Background(), &ModifyDeckVersionCardRequest{
		DeckID:  "deck-1",
		Version: "v1",
	})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientGetDeckVersionHistory(t *testing.T) {
	ts := time.Now().UTC()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/decks/deck-1/versions/v1/history", r.URL.Path)

		require.NoError(t, json.NewEncoder(w).Encode(GetDeckVersionHistoryResponse{
			Changes: []DeckVersionChange{{Version: "v1", Timestamp: &ts}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetDeckVersionHistory(context.Background(), &GetDeckVersionHistoryRequest{
		DeckID:  "deck-1",
		Version: "v1",
	})
	require.NoError(t, err)
	require.Len(t, resp.Changes, 1)
}

func TestClientGetAndUpdateDeckVersionNotes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			require.Equal(t, "/v1/decks/deck-1/versions/v1/notes", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(GetDeckVersionNotesResponse{
				Notes: "Initial",
			}))
		case http.MethodPut:
			require.Equal(t, "/v1/decks/deck-1/versions/v1/notes", r.URL.Path)
			var payload struct {
				Notes string `json:"notes"`
			}
			require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
			require.Equal(t, "Updated", payload.Notes)
			require.NoError(t, json.NewEncoder(w).Encode(UpdateDeckVersionNotesResponse{
				Notes: "Updated",
			}))
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	getResp, err := client.GetDeckVersionNotes(context.Background(), &GetDeckVersionNotesRequest{
		DeckID:  "deck-1",
		Version: "v1",
	})
	require.NoError(t, err)
	require.Equal(t, "Initial", getResp.Notes)

	updateResp, err := client.UpdateDeckVersionNotes(context.Background(), &UpdateDeckVersionNotesRequest{
		DeckID:  "deck-1",
		Version: "v1",
		Notes:   "Updated",
	})
	require.NoError(t, err)
	require.Equal(t, "Updated", updateResp.Notes)
}

func TestClientBatchGetDecks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/decks:batchGet", r.URL.Path)

		var payload BatchGetDecksRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, []string{"deck-1", "deck-2"}, payload.DeckIDs)

		require.NoError(t, json.NewEncoder(w).Encode(BatchGetDecksResponse{
			Decks: []Deck{{ID: "deck-1"}, {ID: "deck-2"}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.BatchGetDecks(context.Background(), &BatchGetDecksRequest{
		DeckIDs: []string{"deck-1", "deck-2"},
	})
	require.NoError(t, err)
	require.Len(t, resp.Decks, 2)

	resp, err = client.BatchGetDecks(context.Background(), &BatchGetDecksRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientListStarredDecks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/users/user-1/stars/decks", r.URL.Path)
		query := r.URL.Query()
		require.Equal(t, "20", query.Get("limit"))
		require.Equal(t, "5", query.Get("offset"))

		require.NoError(t, json.NewEncoder(w).Encode(ListStarredDecksResponse{
			Decks: []Deck{{ID: "deck-1"}},
			Total: 1,
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	limit := int32(20)
	offset := int32(5)
	resp, err := client.ListStarredDecks(context.Background(), &ListStarredDecksRequest{
		UserID: "user-1",
		Limit:  &limit,
		Offset: &offset,
	})
	require.NoError(t, err)
	require.Equal(t, int32(1), resp.Total)
}
