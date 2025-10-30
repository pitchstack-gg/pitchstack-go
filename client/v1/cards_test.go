package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientSearchCards(t *testing.T) {
	t.Run("when filters provided, then query parameters are set and response decoded", func(t *testing.T) {
		pageSize := int32(20)
		boolTrue := true
		boolFalse := false

		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			values := r.URL.Query()

			require.Equal(t, "Azalea", values.Get("searchTerm"))
			require.Equal(t, "Ranger", values.Get("class"))
			require.Equal(t, "Action", values.Get("type"))
			require.Equal(t, "Attack", values.Get("subtype"))
			require.Equal(t, "Shadow", values.Get("talent"))
			require.Equal(t, "2", values.Get("cost"))
			require.Equal(t, "3", values.Get("defense"))
			require.Equal(t, "1", values.Get("pitch"))
			require.Equal(t, "5", values.Get("power"))
			require.Equal(t, "4", values.Get("health"))
			require.Equal(t, "3", values.Get("intelligence"))
			require.Equal(t, "1", values.Get("arcane"))
			require.Equal(t, "COLOR_IDENTITY_RED", values.Get("colorIdentity"))
			require.Equal(t, strconv.FormatBool(true), values.Get("blitzLegal"))
			require.Equal(t, strconv.FormatBool(false), values.Get("ccBanned"))
			require.Equal(t, "20", values.Get("pageSize"))
			require.Equal(t, "token-123", values.Get("nextToken"))

			resp := SearchCardsResponse{
				Summaries: []CardSummary{{Identifier: "card-1"}},
				NextToken: "next-token",
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))

		resp, err := client.SearchCards(context.Background(), &SearchCardsRequest{
			SearchTerm:    "Azalea",
			Class:         "Ranger",
			Type:          "Action",
			Subtype:       "Attack",
			Talent:        "Shadow",
			Cost:          "2",
			Defense:       "3",
			Pitch:         "1",
			Power:         "5",
			Health:        "4",
			Intelligence:  "3",
			Arcane:        "1",
			ColorIdentity: "COLOR_IDENTITY_RED",
			BlitzLegal:    &boolTrue,
			CCBanned:      &boolFalse,
			PageSize:      &pageSize,
			NextToken:     "token-123",
		})
		require.NoError(t, err)
		require.Len(t, resp.Summaries, 1)
		require.Equal(t, "next-token", resp.NextToken)
	})

	t.Run("when request is nil, then defaults applied and response decoded", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Empty(t, r.URL.RawQuery)
			resp := SearchCardsResponse{}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.SearchCards(context.Background(), nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}

func TestClientGetCard(t *testing.T) {
	t.Run("when id provided, then card summary returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/cards/card-1", r.URL.Path)
			resp := GetCardResponse{Summary: &CardSummary{Identifier: "card-1"}}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetCard(context.Background(), &GetCardRequest{CardID: "card-1"})
		require.NoError(t, err)
		require.Equal(t, "card-1", resp.Summary.Identifier)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetCard(context.Background(), &GetCardRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientBatchGetCards(t *testing.T) {
	t.Run("when ids provided, then cards returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			var body BatchGetCardsRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.ElementsMatch(t, []string{"card-1", "card-2"}, body.CardIDs)
			require.True(t, body.AllowPartial)

			resp := BatchGetCardsResponse{
				Cards:       map[string]CardSummary{"card-1": {Identifier: "card-1"}},
				NotFoundIDs: []string{"card-3"},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BatchGetCards(context.Background(), &BatchGetCardsRequest{
			CardIDs:      []string{"card-1", "card-2"},
			AllowPartial: true,
		})
		require.NoError(t, err)
		require.Len(t, resp.Cards, 1)
		require.Equal(t, []string{"card-3"}, resp.NotFoundIDs)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchGetCards(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientListCardPrintings(t *testing.T) {
	t.Run("when card id provided, then printings listed", func(t *testing.T) {
		pageSize := int32(15)
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/cards/card-1/printings", r.URL.Path)
			require.Equal(t, "15", r.URL.Query().Get("pageSize"))
			require.Equal(t, "token", r.URL.Query().Get("nextToken"))

			resp := ListCardPrintingsResponse{
				Summaries: []CardPrintingSummary{{Identifier: "printing-1"}},
				NextToken: "next",
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListCardPrintings(context.Background(), &ListCardPrintingsRequest{
			CardID:    "card-1",
			PageSize:  &pageSize,
			NextToken: "token",
		})
		require.NoError(t, err)
		require.Len(t, resp.Summaries, 1)
		require.Equal(t, "next", resp.NextToken)
	})

	t.Run("when card id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.ListCardPrintings(context.Background(), &ListCardPrintingsRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientListCardPrintingsForSetNumber(t *testing.T) {
	t.Run("when set number provided, then summaries returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/printings/set/XYZ123", r.URL.Path)
			resp := ListCardPrintingsForSetNumberResponse{
				Summaries: []CardPrintingSummary{{Identifier: "printing-1"}},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListCardPrintingsForSetNumber(context.Background(), &ListCardPrintingsForSetNumberRequest{
			SetNumber: "XYZ123",
		})
		require.NoError(t, err)
		require.Len(t, resp.Summaries, 1)
	})

	t.Run("when set number missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.ListCardPrintingsForSetNumber(context.Background(), &ListCardPrintingsForSetNumberRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetCardPrinting(t *testing.T) {
	t.Run("when id provided, then printing returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/printings/printing-1", r.URL.Path)
			resp := GetCardPrintingResponse{Summary: &CardPrintingSummary{Identifier: "printing-1"}}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetCardPrinting(context.Background(), &GetCardPrintingRequest{PrintingID: "printing-1"})
		require.NoError(t, err)
		require.Equal(t, "printing-1", resp.Summary.Identifier)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetCardPrinting(context.Background(), &GetCardPrintingRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientBatchGetCardPrintings(t *testing.T) {
	t.Run("when ids provided, then printings returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			var body BatchGetCardPrintingsRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.ElementsMatch(t, []string{"printing-1", "printing-2"}, body.PrintingIDs)

			resp := BatchGetCardPrintingsResponse{
				Printings:   map[string]CardPrintingSummary{"printing-1": {Identifier: "printing-1"}},
				NotFoundIDs: []string{"printing-3"},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BatchGetCardPrintings(context.Background(), &BatchGetCardPrintingsRequest{
			PrintingIDs: []string{"printing-1", "printing-2"},
		})
		require.NoError(t, err)
		require.Len(t, resp.Printings, 1)
		require.Equal(t, []string{"printing-3"}, resp.NotFoundIDs)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchGetCardPrintings(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetProduct(t *testing.T) {
	t.Run("when id provided, then product returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/products/prod-1", r.URL.Path)
			resp := GetProductResponse{Summary: &ProductSummary{Identifier: "prod-1"}}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetProduct(context.Background(), &GetProductRequest{ProductID: "prod-1"})
		require.NoError(t, err)
		require.Equal(t, "prod-1", resp.Summary.Identifier)
	})

	t.Run("when id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetProduct(context.Background(), &GetProductRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientBatchGetProducts(t *testing.T) {
	t.Run("when ids provided, then summaries returned", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			var body BatchGetProductsRequest
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			require.ElementsMatch(t, []string{"prod-1", "prod-2"}, body.ProductIDs)
			require.True(t, body.AllowPartial)

			resp := BatchGetProductsResponse{
				Products: map[string]ProductSummary{"prod-1": {Identifier: "prod-1"}},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.BatchGetProducts(context.Background(), &BatchGetProductsRequest{
			ProductIDs:   []string{"prod-1", "prod-2"},
			AllowPartial: true,
		})
		require.NoError(t, err)
		require.Len(t, resp.Products, 1)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.BatchGetProducts(context.Background(), nil)
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientGetDataSnapshot(t *testing.T) {
	t.Run("when version filters provided, then query string set", func(t *testing.T) {
		schemaVersion := int32(3)
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/cards/data-snapshots", r.URL.Path)
			require.Equal(t, "3", r.URL.Query().Get("schemaVersion"))
			require.Equal(t, "2024-01-01", r.URL.Query().Get("version"))

			resp := GetDataSnapshotResponse{
				Manifest: &DataSnapshotManifest{Version: "2024-01-01"},
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}

		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetDataSnapshot(context.Background(), &GetDataSnapshotRequest{
			SchemaVersion: &schemaVersion,
			Version:       "2024-01-01",
		})
		require.NoError(t, err)
		require.Equal(t, "2024-01-01", resp.Manifest.Version)
	})

	t.Run("when request is nil, then defaults applied", func(t *testing.T) {
		handler := func(w http.ResponseWriter, r *http.Request) {
			require.Empty(t, r.URL.Query())
			resp := GetDataSnapshotResponse{}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
		}
		server := httptest.NewServer(http.HandlerFunc(handler))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetDataSnapshot(context.Background(), nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}
