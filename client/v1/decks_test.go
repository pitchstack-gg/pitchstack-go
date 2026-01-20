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
		require.Equal(t, "subject-1", query.Get("subjectId"))

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
		SubjectID: "subject-1",
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
		require.Equal(t, "desc", payload.Description)
		require.Equal(t, "author", payload.Author)
		require.Equal(t, VisibilityLevelShared, payload.Visibility)
		require.Equal(t, "deck-client-1", payload.DeckID)
		require.NotNil(t, payload.CreateInitialVersion)
		require.True(t, *payload.CreateInitialVersion)
		require.NotNil(t, payload.InitialVersion)
		require.Equal(t, "v1", payload.InitialVersion.Name)
		require.Equal(t, "dv-1", payload.InitialVersion.DeckVersionID)

		require.NoError(t, json.NewEncoder(w).Encode(CreateDeckResponse{
			Deck: &Deck{ID: "deck-1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	createInitial := true
	resp, err := client.CreateDeck(context.Background(), &CreateDeckRequest{
		Name:                 "Aggro",
		HeroID:               "hero-1",
		Format:               "cc",
		Description:          "desc",
		Author:               "author",
		Visibility:           VisibilityLevelShared,
		DeckID:               "deck-client-1",
		CreateInitialVersion: &createInitial,
		InitialVersion: &CreateDeckInitialVersion{
			Name:          "v1",
			DeckVersionID: "dv-1",
		},
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

func TestClientGetDeckAccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/decks/deck-1/access", r.URL.Path)
		require.NoError(t, json.NewEncoder(w).Encode(GetDeckAccessResponse{
			Permission: DeckPermissionReader,
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetDeckAccess(context.Background(), &GetDeckAccessRequest{DeckID: "deck-1"})
	require.NoError(t, err)
	require.Equal(t, DeckPermissionReader, resp.Permission)

	resp, err = client.GetDeckAccess(context.Background(), &GetDeckAccessRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientListDeckAccessGrants(t *testing.T) {
	pageSize := int32(5)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/decks/deck-1/permissions", r.URL.Path)
		require.Equal(t, "5", r.URL.Query().Get("pageSize"))
		require.Equal(t, "token", r.URL.Query().Get("nextToken"))

		require.NoError(t, json.NewEncoder(w).Encode(ListDeckAccessGrantsResponse{
			Grants: []DeckAccessGrant{{SubjectID: "user-1", Permission: DeckPermissionWriter}},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ListDeckAccessGrants(context.Background(), &ListDeckAccessGrantsRequest{
		DeckID:    "deck-1",
		PageSize:  &pageSize,
		NextToken: "token",
	})
	require.NoError(t, err)
	require.Len(t, resp.Grants, 1)
	require.Equal(t, "user-1", resp.Grants[0].SubjectID)
	require.Equal(t, DeckPermissionWriter, resp.Grants[0].Permission)

	resp, err = client.ListDeckAccessGrants(context.Background(), &ListDeckAccessGrantsRequest{})
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
			Name        string `json:"name"`
			Description string `json:"description"`
			Author      string `json:"author"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "Updated", payload.Name)
		require.Equal(t, "New desc", payload.Description)
		require.Equal(t, "New author", payload.Author)

		require.NoError(t, json.NewEncoder(w).Encode(UpdateDeckResponse{
			Deck: &Deck{ID: "deck-1", Name: "Updated"},
		}))

	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	name := "Updated"
	description := "New desc"
	author := "New author"
	resp, err := client.UpdateDeck(context.Background(), &UpdateDeckRequest{
		DeckID:      "deck-1",
		Name:        &name,
		Description: &description,
		Author:      &author,
	})
	require.NoError(t, err)
	require.Equal(t, "Updated", resp.Deck.Name)

	resp, err = client.UpdateDeck(context.Background(), &UpdateDeckRequest{DeckID: "deck-1"})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientUpdateDeckVisibility(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/v1/decks/deck-1/visibility", r.URL.Path)
		var payload struct {
			Visibility string `json:"visibility"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, string(VisibilityLevelShared), payload.Visibility)

		require.NoError(t, json.NewEncoder(w).Encode(UpdateDeckVisibilityResponse{
			Deck: &Deck{ID: "deck-1", Visibility: VisibilityLevelShared},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	visibility := VisibilityLevelShared
	resp, err := client.UpdateDeckVisibility(context.Background(), &UpdateDeckVisibilityRequest{
		DeckID:     "deck-1",
		Visibility: &visibility,
	})
	require.NoError(t, err)
	require.Equal(t, VisibilityLevelShared, resp.Deck.Visibility)

	resp, err = client.UpdateDeckVisibility(context.Background(), &UpdateDeckVisibilityRequest{DeckID: "deck-1"})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientGrantDeckAccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/decks/permissions:grant", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload struct {
			DeckID     string         `json:"resourceId"`
			SubjectID  string         `json:"subjectId"`
			Permission DeckPermission `json:"permission"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "deck-1", payload.DeckID)
		require.Equal(t, "user-1", payload.SubjectID)
		require.Equal(t, DeckPermissionReader, payload.Permission)

		w.Header().Set("X-Request-Id", "req-grant")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GrantDeckAccess(context.Background(), &GrantDeckAccessRequest{
		DeckID:     "deck-1",
		SubjectID:  "user-1",
		Permission: DeckPermissionReader,
	})
	require.NoError(t, err)
	require.Equal(t, "req-grant", resp.Metadata.RequestID)

	resp, err = client.GrantDeckAccess(context.Background(), &GrantDeckAccessRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientRevokeDeckAccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/decks/permissions:revoke", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload struct {
			DeckID     string         `json:"resourceId"`
			SubjectID  string         `json:"subjectId"`
			Permission DeckPermission `json:"permission"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "deck-1", payload.DeckID)
		require.Equal(t, "user-1", payload.SubjectID)
		require.Equal(t, DeckPermissionWriter, payload.Permission)

		w.Header().Set("X-Request-Id", "req-revoke")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.RevokeDeckAccess(context.Background(), &RevokeDeckAccessRequest{
		DeckID:     "deck-1",
		SubjectID:  "user-1",
		Permission: DeckPermissionWriter,
	})
	require.NoError(t, err)
	require.Equal(t, "req-revoke", resp.Metadata.RequestID)

	resp, err = client.RevokeDeckAccess(context.Background(), &RevokeDeckAccessRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
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
			DeckVersions: []DeckVersion{{ID: "dv-1", Name: "v1"}},
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
	require.Equal(t, []DeckVersion{{ID: "dv-1", Name: "v1"}}, resp.DeckVersions)
}

func TestClientCreateDeckVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/decks/deck-1/versions", r.URL.Path)

		var payload struct {
			Name                string `json:"name"`
			ImageURL            string `json:"imageUrl"`
			Description         string `json:"description"`
			DeckVersionID       string `json:"deckVersionId"`
			SourceDeckVersionID string `json:"sourceDeckVersionId"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "v2", payload.Name)
		require.Equal(t, "https://image", payload.ImageURL)
		require.Equal(t, "desc", payload.Description)
		require.Equal(t, "dv-2", payload.DeckVersionID)
		require.Equal(t, "dv-source", payload.SourceDeckVersionID)

		require.NoError(t, json.NewEncoder(w).Encode(CreateDeckVersionResponse{
			DeckVersion: &DeckVersion{ID: "dv-2", Name: "v2"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.CreateDeckVersion(context.Background(), &CreateDeckVersionRequest{
		DeckID:              "deck-1",
		Name:                "v2",
		ImageURL:            "https://image",
		Description:         "desc",
		DeckVersionID:       "dv-2",
		SourceDeckVersionID: "dv-source",
	})
	require.NoError(t, err)
	require.Equal(t, "v2", resp.DeckVersion.Name)
}

func TestClientGetDeckVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/deck_versions/dv-1", r.URL.Path)

		require.NoError(t, json.NewEncoder(w).Encode(GetDeckVersionResponse{
			DeckVersion: &DeckVersion{ID: "dv-1", Name: "v1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetDeckVersion(context.Background(), &GetDeckVersionRequest{
		DeckVersionID: "dv-1",
	})
	require.NoError(t, err)
	require.Equal(t, "v1", resp.DeckVersion.Name)

	resp, err = client.GetDeckVersion(context.Background(), &GetDeckVersionRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientDeleteDeckVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/deck_versions/dv-1", r.URL.Path)
		w.Header().Set("X-Request-Id", "req-delete-version")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.DeleteDeckVersion(context.Background(), &DeleteDeckVersionRequest{
		DeckVersionID: "dv-1",
	})

	require.NoError(t, err)
	require.Equal(t, "req-delete-version", resp.Metadata.RequestID)
}

func TestClientUpdateDeckVersion(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/deck_versions/dv-1", r.URL.Path)

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
		DeckVersionID: "dv-1",
		ImageURL:      &imageURL,
		Description:   &desc,
	})
	require.NoError(t, err)
	require.Equal(t, "https://image", resp.DeckVersion.ImageURL)
}

func TestClientListDeckVersionCards(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/deck_versions/dv-1/cards", r.URL.Path)

		require.NoError(t, json.NewEncoder(w).Encode(ListDeckVersionCardsResponse{
			MainboardCards:  []DeckCard{{CardID: "card-1", Quantity: 3}},
			SideboardCards:  []DeckCard{{CardID: "card-2", Quantity: 2}},
			MaybeboardCards: []DeckCard{{CardID: "card-3", Quantity: 1}},
		}))

	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ListDeckVersionCards(context.Background(), &ListDeckVersionCardsRequest{
		DeckVersionID: "dv-1",
	})
	require.NoError(t, err)
	require.Len(t, resp.MainboardCards, 1)
	require.Len(t, resp.SideboardCards, 1)
	require.Len(t, resp.MaybeboardCards, 1)
}

func TestClientModifyDeckVersionCard(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/deck_versions/dv-1/cards", r.URL.Path)
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
		DeckVersionID: "dv-1",
		CardID:        "card-1",
		Board:         BoardTypeMainboard,
		Quantity:      3,
	})
	require.NoError(t, err)
	require.Equal(t, "card-1", resp.Card.CardID)

	resp, err = client.ModifyDeckVersionCard(context.Background(), &ModifyDeckVersionCardRequest{
		DeckVersionID: "dv-1",
	})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientGetDeckVersionHistory(t *testing.T) {
	ts := time.Now().UTC()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/deck_versions/dv-1/history", r.URL.Path)

		require.NoError(t, json.NewEncoder(w).Encode(GetDeckVersionHistoryResponse{
			Changes: []DeckVersionChange{{Name: "v1", Timestamp: &ts}},
		}))

	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.GetDeckVersionHistory(context.Background(), &GetDeckVersionHistoryRequest{
		DeckVersionID: "dv-1",
	})
	require.NoError(t, err)
	require.Len(t, resp.Changes, 1)
}

func TestClientGetAndUpdateDeckVersionNotes(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			require.Equal(t, "/v1/deck_versions/dv-1/notes", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(GetDeckVersionNotesResponse{
				Notes: "Initial",
			}))
		case http.MethodPut:
			require.Equal(t, "/v1/deck_versions/dv-1/notes", r.URL.Path)
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
		DeckVersionID: "dv-1",
	})
	require.NoError(t, err)
	require.Equal(t, "Initial", getResp.Notes)

	updateResp, err := client.UpdateDeckVersionNotes(context.Background(), &UpdateDeckVersionNotesRequest{
		DeckVersionID: "dv-1",
		Notes:         "Updated",
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
