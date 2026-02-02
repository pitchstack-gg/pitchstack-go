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

func TestClientCloneDeck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/decks:clone", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload CloneDeckRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "dv-source", payload.SourceDeckVersionID)
		require.Equal(t, "New Deck", payload.Name)
		require.NotNil(t, payload.Visibility)
		require.Equal(t, VisibilityLevelShared, *payload.Visibility)
		require.Equal(t, "deck-1", payload.DeckID)
		require.Equal(t, "v1", payload.InitialVersionName)
		require.Equal(t, "dv-1", payload.InitialDeckVersionID)

		require.NoError(t, json.NewEncoder(w).Encode(CloneDeckResponse{
			Deck:           &Deck{ID: "deck-1"},
			InitialVersion: &DeckVersion{ID: "dv-1"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	visibility := VisibilityLevelShared
	resp, err := client.CloneDeck(context.Background(), &CloneDeckRequest{
		SourceDeckVersionID:  "dv-source",
		Name:                 "New Deck",
		Visibility:           &visibility,
		DeckID:               "deck-1",
		InitialVersionName:   "v1",
		InitialDeckVersionID: "dv-1",
	})
	require.NoError(t, err)
	require.Equal(t, "deck-1", resp.Deck.ID)
	require.Equal(t, "dv-1", resp.InitialVersion.ID)
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
			Name                string `json:"name"`
			Author              string `json:"author"`
			ActiveDeckVersionID string `json:"activeDeckVersionId"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "Updated", payload.Name)
		require.Equal(t, "New author", payload.Author)
		require.Equal(t, "dv-active", payload.ActiveDeckVersionID)

		require.NoError(t, json.NewEncoder(w).Encode(UpdateDeckResponse{
			Deck: &Deck{ID: "deck-1", Name: "Updated"},
		}))

	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	name := "Updated"
	author := "New author"
	activeVersion := "dv-active"
	resp, err := client.UpdateDeck(context.Background(), &UpdateDeckRequest{
		DeckID:              "deck-1",
		Name:                &name,
		Author:              &author,
		ActiveDeckVersionID: &activeVersion,
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
			DeckVersionID       string `json:"deckVersionId"`
			SourceDeckVersionID string `json:"sourceDeckVersionId"`
			SetActive           *bool  `json:"setActive"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "v2", payload.Name)
		require.Equal(t, "dv-2", payload.DeckVersionID)
		require.Equal(t, "dv-source", payload.SourceDeckVersionID)
		require.NotNil(t, payload.SetActive)
		require.False(t, *payload.SetActive)

		require.NoError(t, json.NewEncoder(w).Encode(CreateDeckVersionResponse{
			DeckVersion: &DeckVersion{ID: "dv-2", Name: "v2"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	setActive := false
	resp, err := client.CreateDeckVersion(context.Background(), &CreateDeckVersionRequest{
		DeckID:              "deck-1",
		Name:                "v2",
		DeckVersionID:       "dv-2",
		SourceDeckVersionID: "dv-source",
		SetActive:           &setActive,
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

func TestClientListDeckVersionSideboardGuides(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/deck_versions/dv-1/sideboard_guides", r.URL.Path)

		require.NoError(t, json.NewEncoder(w).Encode(ListDeckVersionSideboardGuidesResponse{
			SideboardGuides: []SideboardGuide{
				{
					TargetType: SideboardGuideTargetTypeHero,
					Target:     "Dorinthea",
					Guide:      "Swap in 3x Sink Below.",
				},
			},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ListDeckVersionSideboardGuides(context.Background(), &ListDeckVersionSideboardGuidesRequest{
		DeckVersionID: "dv-1",
	})
	require.NoError(t, err)
	require.Len(t, resp.SideboardGuides, 1)
	require.Equal(t, SideboardGuideTargetTypeHero, resp.SideboardGuides[0].TargetType)
}

func TestClientUpsertDeckVersionSideboardGuide(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPut, r.Method)
		require.Equal(t, "/v1/deck_versions/dv-1/sideboard_guides", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload struct {
			TargetType string `json:"targetType"`
			Target     string `json:"target"`
			Guide      string `json:"guide"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, string(SideboardGuideTargetTypeClass), payload.TargetType)
		require.Equal(t, "Warrior", payload.Target)
		require.Equal(t, "Prioritize armor.", payload.Guide)

		require.NoError(t, json.NewEncoder(w).Encode(UpsertDeckVersionSideboardGuideResponse{
			SideboardGuide: &SideboardGuide{
				TargetType: SideboardGuideTargetTypeClass,
				Target:     "Warrior",
				Guide:      "Prioritize armor.",
			},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.UpsertDeckVersionSideboardGuide(context.Background(), &UpsertDeckVersionSideboardGuideRequest{
		DeckVersionID: "dv-1",
		TargetType:    SideboardGuideTargetTypeClass,
		Target:        "Warrior",
		Guide:         "Prioritize armor.",
	})
	require.NoError(t, err)
	require.Equal(t, "Warrior", resp.SideboardGuide.Target)
}

func TestClientDeleteDeckVersionSideboardGuide(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/v1/deck_versions/dv-1/sideboard_guides", r.URL.Path)
		require.Equal(t, string(SideboardGuideTargetTypeArchetype), r.URL.Query().Get("targetType"))
		require.Equal(t, "Fatigue", r.URL.Query().Get("target"))

		w.Header().Set("X-Request-Id", "req-delete-guide")
		_, _ = w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.DeleteDeckVersionSideboardGuide(context.Background(), &DeleteDeckVersionSideboardGuideRequest{
		DeckVersionID: "dv-1",
		TargetType:    SideboardGuideTargetTypeArchetype,
		Target:        "Fatigue",
	})
	require.NoError(t, err)
	require.Equal(t, "req-delete-guide", resp.Metadata.RequestID)
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
		require.True(t, payload.AllowPartial)

		require.NoError(t, json.NewEncoder(w).Encode(BatchGetDecksResponse{
			Decks:       []Deck{{ID: "deck-1"}, {ID: "deck-2"}},
			NotFoundIDs: []string{"deck-3"},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.BatchGetDecks(context.Background(), &BatchGetDecksRequest{
		DeckIDs:      []string{"deck-1", "deck-2"},
		AllowPartial: true,
	})
	require.NoError(t, err)
	require.Len(t, resp.Decks, 2)
	require.Equal(t, []string{"deck-3"}, resp.NotFoundIDs)

	resp, err = client.BatchGetDecks(context.Background(), &BatchGetDecksRequest{})
	require.Nil(t, resp)
	require.Error(t, err)
}

func TestClientExportDeck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)
		require.Equal(t, "/v1/decks/deck-1:export", r.URL.Path)

		require.NoError(t, json.NewEncoder(w).Encode(ExportDeckResponse{
			Deck: &Deck{ID: "deck-1"},
			Versions: []ExportDeckVersion{
				{
					DeckVersion:     &DeckVersion{ID: "dv-1"},
					Notes:           "notes",
					MainboardCards:  []DeckCard{{CardID: "card-1", Quantity: 3}},
					SideboardCards:  []DeckCard{{CardID: "card-2", Quantity: 1}},
					MaybeboardCards: []DeckCard{{CardID: "card-3", Quantity: 2}},
					SideboardGuides: []SideboardGuide{{
						TargetType: SideboardGuideTargetTypeHero,
						Target:     "Dorinthea",
						Guide:      "Swap in 3x Sink Below.",
					}},
				},
			},
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	resp, err := client.ExportDeck(context.Background(), &ExportDeckRequest{DeckID: "deck-1"})
	require.NoError(t, err)
	require.Equal(t, "deck-1", resp.Deck.ID)
	require.Len(t, resp.Versions, 1)
	require.Equal(t, "dv-1", resp.Versions[0].DeckVersion.ID)
	require.Len(t, resp.Versions[0].SideboardGuides, 1)
}

func TestClientImportDeck(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/v1/decks:import", r.URL.Path)
		require.Equal(t, "application/json", r.Header.Get("Content-Type"))

		var payload ImportDeckRequest
		require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		require.Equal(t, "deck-1", payload.DeckID)
		require.Equal(t, "Imported", payload.Name)
		require.Equal(t, "hero-1", payload.HeroID)
		require.Equal(t, "cc", payload.Format)
		require.Equal(t, "author", payload.Author)
		require.NotNil(t, payload.Visibility)
		require.Equal(t, VisibilityLevelPrivate, *payload.Visibility)
		require.Equal(t, "dv-active", payload.ActiveDeckVersionID)
		require.Len(t, payload.Versions, 1)
		require.Equal(t, "dv-1", payload.Versions[0].DeckVersionID)
		require.Equal(t, "v1", payload.Versions[0].Name)
		require.Equal(t, "notes", payload.Versions[0].Notes)
		require.Len(t, payload.Versions[0].MainboardCards, 1)
		require.Len(t, payload.Versions[0].SideboardGuides, 1)

		require.NoError(t, json.NewEncoder(w).Encode(ImportDeckResponse{
			Deck:             &Deck{ID: "deck-1"},
			ImportedVersions: 1,
		}))
	}))
	t.Cleanup(server.Close)

	client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
	visibility := VisibilityLevelPrivate
	resp, err := client.ImportDeck(context.Background(), &ImportDeckRequest{
		DeckID:              "deck-1",
		Name:                "Imported",
		HeroID:              "hero-1",
		Format:              "cc",
		Author:              "author",
		Visibility:          &visibility,
		ActiveDeckVersionID: "dv-active",
		Versions: []ImportDeckVersion{
			{
				DeckVersionID:  "dv-1",
				Name:           "v1",
				Notes:          "notes",
				MainboardCards: []DeckCard{{CardID: "card-1", Quantity: 3}},
				SideboardGuides: []SideboardGuide{{
					TargetType: SideboardGuideTargetTypeHero,
					Target:     "Dorinthea",
					Guide:      "Swap in 3x Sink Below.",
				}},
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "deck-1", resp.Deck.ID)
	require.Equal(t, int32(1), resp.ImportedVersions)
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
