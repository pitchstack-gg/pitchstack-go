package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// DeckVersionMatchResult represents decks.v1.DeckVersionMatchResult.
type DeckVersionMatchResult string

const (
	DeckVersionMatchResultUnspecified DeckVersionMatchResult = "DECK_VERSION_MATCH_RESULT_UNSPECIFIED"
	DeckVersionMatchResultWin         DeckVersionMatchResult = "DECK_VERSION_MATCH_RESULT_WIN"
	DeckVersionMatchResultLoss        DeckVersionMatchResult = "DECK_VERSION_MATCH_RESULT_LOSS"
	DeckVersionMatchResultDraw        DeckVersionMatchResult = "DECK_VERSION_MATCH_RESULT_DRAW"
)

// DeckVersionMatchCardChange mirrors decks.v1.DeckVersionMatchCardChange.
type DeckVersionMatchCardChange struct {
	CardID   string `json:"cardId,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
}

// DeckVersionMatchCoordinates mirrors decks.v1.DeckVersionMatchCoordinates.
type DeckVersionMatchCoordinates struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

// DeckVersionMatch mirrors decks.v1.DeckVersionMatch.
type DeckVersionMatch struct {
	ID                  string                       `json:"id,omitempty"`
	DeckVersionID       string                       `json:"deckVersionId,omitempty"`
	PlayedAt            *time.Time                   `json:"playedAt,omitempty"`
	Result              DeckVersionMatchResult       `json:"result,omitempty"`
	OpponentHeroID      string                       `json:"opponentHeroId,omitempty"`
	Opponent            string                       `json:"opponent,omitempty"`
	WentFirst           *bool                        `json:"wentFirst,omitempty"`
	EventType           string                       `json:"eventType,omitempty"`
	Notes               string                       `json:"notes,omitempty"`
	Location            string                       `json:"location,omitempty"`
	LocationCoordinates *DeckVersionMatchCoordinates `json:"locationCoordinates,omitempty"`
	CardsIn             []DeckVersionMatchCardChange `json:"cardsIn,omitempty"`
	CardsOut            []DeckVersionMatchCardChange `json:"cardsOut,omitempty"`
	CreatedAt           *time.Time                   `json:"createdAt,omitempty"`
	UpdatedAt           *time.Time                   `json:"updatedAt,omitempty"`
}

// CreateDeckVersionMatchRequest creates a match record for a deck version.
type CreateDeckVersionMatchRequest struct {
	DeckVersionID       string                       `json:"-"`
	MatchID             string                       `json:"matchId,omitempty"`
	PlayedAt            *time.Time                   `json:"playedAt,omitempty"`
	Result              DeckVersionMatchResult       `json:"result,omitempty"`
	OpponentHeroID      string                       `json:"opponentHeroId,omitempty"`
	Opponent            string                       `json:"opponent,omitempty"`
	WentFirst           *bool                        `json:"wentFirst,omitempty"`
	EventType           string                       `json:"eventType,omitempty"`
	Notes               string                       `json:"notes,omitempty"`
	Location            string                       `json:"location,omitempty"`
	LocationCoordinates *DeckVersionMatchCoordinates `json:"locationCoordinates,omitempty"`
	CardsIn             []DeckVersionMatchCardChange `json:"cardsIn,omitempty"`
	CardsOut            []DeckVersionMatchCardChange `json:"cardsOut,omitempty"`
}

// CreateDeckVersionMatchResponse returns the created match.
type CreateDeckVersionMatchResponse struct {
	Match    *DeckVersionMatch `json:"match,omitempty"`
	Metadata ResponseMetadata  `json:"-"`
}

func (r *CreateDeckVersionMatchResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetDeckVersionMatchRequest identifies a specific deck-version match.
type GetDeckVersionMatchRequest struct {
	DeckVersionID string `json:"-"`
	MatchID       string `json:"-"`
}

// GetDeckVersionMatchResponse returns one deck-version match.
type GetDeckVersionMatchResponse struct {
	Match    *DeckVersionMatch `json:"match,omitempty"`
	Metadata ResponseMetadata  `json:"-"`
}

func (r *GetDeckVersionMatchResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListDeckVersionMatchesRequest filters match listing for a deck version.
type ListDeckVersionMatchesRequest struct {
	DeckVersionID string
	PageSize      *int32
	NextToken     string
}

// ListDeckVersionMatchesResponse returns paginated deck-version matches.
type ListDeckVersionMatchesResponse struct {
	Matches   []DeckVersionMatch `json:"matches,omitempty"`
	NextToken string             `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata   `json:"-"`
}

func (r *ListDeckVersionMatchesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteDeckVersionMatchRequest identifies a match to delete.
type DeleteDeckVersionMatchRequest struct {
	DeckVersionID string `json:"-"`
	MatchID       string `json:"-"`
}

// DeleteDeckVersionMatchResponse captures response metadata.
type DeleteDeckVersionMatchResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteDeckVersionMatchResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateDeckVersionMatch creates a match record under a deck version.
func (c *Client) CreateDeckVersionMatch(ctx context.Context, request *CreateDeckVersionMatchRequest, opts ...RequestOpt) (*CreateDeckVersionMatchResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}
	if strings.TrimSpace(string(request.Result)) == "" || request.Result == DeckVersionMatchResultUnspecified {
		return nil, errors.New("result must be specified")
	}

	body, err := jsonBody(struct {
		DeckVersionID       string                       `json:"deckVersionId,omitempty"`
		MatchID             string                       `json:"matchId,omitempty"`
		PlayedAt            *time.Time                   `json:"playedAt,omitempty"`
		Result              DeckVersionMatchResult       `json:"result,omitempty"`
		OpponentHeroID      string                       `json:"opponentHeroId,omitempty"`
		Opponent            string                       `json:"opponent,omitempty"`
		WentFirst           *bool                        `json:"wentFirst,omitempty"`
		EventType           string                       `json:"eventType,omitempty"`
		Notes               string                       `json:"notes,omitempty"`
		Location            string                       `json:"location,omitempty"`
		LocationCoordinates *DeckVersionMatchCoordinates `json:"locationCoordinates,omitempty"`
		CardsIn             []DeckVersionMatchCardChange `json:"cardsIn,omitempty"`
		CardsOut            []DeckVersionMatchCardChange `json:"cardsOut,omitempty"`
	}{
		DeckVersionID:       deckVersionID,
		MatchID:             strings.TrimSpace(request.MatchID),
		PlayedAt:            request.PlayedAt,
		Result:              request.Result,
		OpponentHeroID:      strings.TrimSpace(request.OpponentHeroID),
		Opponent:            strings.TrimSpace(request.Opponent),
		WentFirst:           request.WentFirst,
		EventType:           strings.TrimSpace(request.EventType),
		Notes:               request.Notes,
		Location:            strings.TrimSpace(request.Location),
		LocationCoordinates: request.LocationCoordinates,
		CardsIn:             request.CardsIn,
		CardsOut:            request.CardsOut,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/matches", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CreateDeckVersionMatchResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeckVersionMatch fetches a match record by ID.
func (c *Client) GetDeckVersionMatch(ctx context.Context, request *GetDeckVersionMatchRequest, opts ...RequestOpt) (*GetDeckVersionMatchResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}
	matchID := strings.TrimSpace(request.MatchID)
	if matchID == "" {
		return nil, errors.New("matchID must not be empty")
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/matches/%s", url.PathEscape(deckVersionID), url.PathEscape(matchID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetDeckVersionMatchResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListDeckVersionMatches lists match records for a deck version.
func (c *Client) ListDeckVersionMatches(ctx context.Context, request *ListDeckVersionMatchesRequest, opts ...RequestOpt) (*ListDeckVersionMatchesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/matches", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	setQueryString(query, "nextToken", request.NextToken)
	req.URL.RawQuery = query.Encode()

	response := &ListDeckVersionMatchesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteDeckVersionMatch deletes a match record.
func (c *Client) DeleteDeckVersionMatch(ctx context.Context, request *DeleteDeckVersionMatchRequest, opts ...RequestOpt) (*DeleteDeckVersionMatchResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}
	matchID := strings.TrimSpace(request.MatchID)
	if matchID == "" {
		return nil, errors.New("matchID must not be empty")
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/matches/%s", url.PathEscape(deckVersionID), url.PathEscape(matchID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteDeckVersionMatchResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}
