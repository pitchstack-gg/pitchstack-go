package client

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// newStubAPIClient starts an httptest server that runs assert against each
// request and replies with the raw JSON body, mirroring gateway encoding.
func newStubAPIClient(t *testing.T, assert func(t *testing.T, r *http.Request), body string) *Client {
	t.Helper()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if assert != nil {
			assert(t, r)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(body))
	}))
	t.Cleanup(server.Close)
	return newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
}

func TestClientListMetaSources(t *testing.T) {
	t.Run("when includeInactive set, then query sent and sources decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/meta/sources", r.URL.Path)
			require.Equal(t, "true", r.URL.Query().Get("includeInactive"))
		}, `{"sources":[{"sourceId":"src-1","displayName":"Paper","playEnvironment":"PLAY_ENVIRONMENT_PAPER","dataTypes":["PLAY_DATA_TYPE_MATCH_RESULT"],"active":true,"analyticsCapabilities":["ANALYTICS_CAPABILITY_HERO_MATCHUPS"]}]}`)

		resp, err := client.ListMetaSources(context.Background(), &ListMetaSourcesRequest{IncludeInactive: true})
		require.NoError(t, err)
		require.Len(t, resp.Sources, 1)
		require.Equal(t, PlayEnvironmentPaper, resp.Sources[0].PlayEnvironment)
		require.Equal(t, []PlayDataType{PlayDataTypeMatchResult}, resp.Sources[0].DataTypes)
		require.Equal(t, []AnalyticsCapability{AnalyticsCapabilityHeroMatchups}, resp.Sources[0].AnalyticsCapabilities)
		require.True(t, resp.Sources[0].Active)
	})

	t.Run("when request is nil, then no query is sent", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Empty(t, r.URL.RawQuery)
		}, `{}`)

		resp, err := client.ListMetaSources(context.Background(), nil)
		require.NoError(t, err)
		require.NotNil(t, resp)
	})
}

func TestClientListEventStandings(t *testing.T) {
	t.Run("when event id provided, then standings decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/meta/events/evt-1/standings", r.URL.Path)
			require.Equal(t, "src-1", r.URL.Query().Get("sourceId"))
		}, `{"eventId":"evt-1","sourceId":"src-1","standings":[{"place":1,"playerName":"Ada","heroId":"hero-1"}],"totalPlayers":64,"dataThrough":"2026-11-01T00:00:00Z"}`)

		resp, err := client.ListEventStandings(context.Background(), &ListEventStandingsRequest{EventID: "evt-1", SourceID: "src-1"})
		require.NoError(t, err)
		require.Equal(t, uint32(64), resp.TotalPlayers)
		require.Len(t, resp.Standings, 1)
		require.Equal(t, uint32(1), resp.Standings[0].Place)
		require.NotNil(t, resp.DataThrough)
	})

	t.Run("when event id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.ListEventStandings(context.Background(), &ListEventStandingsRequest{})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}

func TestClientListEventConversions(t *testing.T) {
	t.Run("when event id provided, then milestones decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, "/v1/meta/events/evt-1/conversions", r.URL.Path)
			require.Empty(t, r.URL.Query().Get("sourceId"))
		}, `{"eventId":"evt-1","startingEntrants":128,"heroBreakdownStatus":"HERO_CONVERSION_BREAKDOWN_STATUS_AVAILABLE","milestones":[{"key":"top8","kind":"EVENT_CONVERSION_MILESTONE_KIND_PLACEMENT_CUT","placementCutoff":8,"convertedEntrants":8,"fieldConversionRate":0.0625,"previousStageConversionRate":0,"heroes":[{"heroId":"hero-1","startingEntrants":16,"convertedEntrants":2,"conversionRate":0.125}]}]}`)

		resp, err := client.ListEventConversions(context.Background(), &ListEventConversionsRequest{EventID: "evt-1"})
		require.NoError(t, err)
		require.Equal(t, uint32(128), resp.StartingEntrants)
		require.Equal(t, HeroConversionBreakdownStatusAvailable, resp.HeroBreakdownStatus)
		require.Len(t, resp.Milestones, 1)
		milestone := resp.Milestones[0]
		require.Equal(t, EventConversionMilestoneKindPlacementCut, milestone.Kind)
		require.NotNil(t, milestone.PlacementCutoff)
		require.Equal(t, uint32(8), *milestone.PlacementCutoff)
		require.Nil(t, milestone.SourceRound)
		require.NotNil(t, milestone.PreviousStageConversionRate)
		require.Zero(t, *milestone.PreviousStageConversionRate)
		require.Len(t, milestone.Heroes, 1)
	})

	t.Run("when event id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ListEventConversions(context.Background(), nil)
		require.Error(t, err)
	})
}

func TestClientListHeroMatchups(t *testing.T) {
	t.Run("when filters provided, then query encoded and uint64 strings decoded", func(t *testing.T) {
		pageSize := int32(50)
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/meta/hero-matchups", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "PLAY_ENVIRONMENT_DIGITAL", q.Get("playEnvironment"))
			require.Equal(t, "2026-10-01", q.Get("startDate"))
			require.Equal(t, "2026-10-31", q.Get("endDate"))
			require.Equal(t, "cc", q.Get("ruleset"))
			require.Equal(t, "ranked", q.Get("competitionMode"))
			require.Equal(t, "hero-1", q.Get("heroId"))
			require.Equal(t, "hero-2", q.Get("opponentHeroId"))
			require.Equal(t, "SOURCE_VISIBILITY_PUBLIC", q.Get("sourceVisibility"))
			require.Equal(t, "50", q.Get("pageSize"))
			require.Equal(t, "tok", q.Get("nextPageToken"))
			require.Equal(t, []string{"src-1", "src-2"}, q["sourceIds"])
			require.Equal(t, "HERO_MATCHUP_GROUPING_SOURCE_PAIRING", q.Get("grouping"))
			require.Equal(t, "HERO_MATCHUP_PAIRING_ORDER_GAMES_DESC", q.Get("orderBy"))
			require.Equal(t, []string{"evt-1"}, q["eventIds"])
			require.False(t, q.Has("platform"))
		}, `{"matchups":[{"heroId":"hero-1","opponentHeroId":"hero-2","games":"120","wins":"70","winRate":0.58,"sourceVisibility":"SOURCE_VISIBILITY_PUBLIC"}],"eligibleMatches":"500","nextPageToken":"next","cohorts":[{"sourceId":"src-1","eligibleObservations":"500","observedObservations":"0"}]}`)

		resp, err := client.ListHeroMatchups(context.Background(), &ListHeroMatchupsRequest{
			PlayEnvironment:  PlayEnvironmentDigital,
			StartDate:        "2026-10-01",
			EndDate:          "2026-10-31",
			Ruleset:          "cc",
			CompetitionMode:  "ranked",
			HeroID:           "hero-1",
			OpponentHeroID:   "hero-2",
			SourceVisibility: SourceVisibilityPublic,
			PageSize:         &pageSize,
			NextPageToken:    "tok",
			SourceIDs:        []string{"src-1", " ", "src-2"},
			Grouping:         HeroMatchupGroupingSourcePairing,
			OrderBy:          HeroMatchupPairingOrderGamesDesc,
			EventIDs:         []string{"evt-1"},
		})
		require.NoError(t, err)
		require.Len(t, resp.Matchups, 1)
		require.Equal(t, uint64(120), resp.Matchups[0].Games)
		require.Equal(t, uint64(70), resp.Matchups[0].Wins)
		require.Equal(t, uint64(500), resp.EligibleMatches)
		require.Equal(t, "next", resp.NextPageToken)
		require.Len(t, resp.Cohorts, 1)
		require.NotNil(t, resp.Cohorts[0].ObservedObservations)
		require.Zero(t, *resp.Cohorts[0].ObservedObservations)
		require.Nil(t, resp.Cohorts[0].ExcludedObservations)
	})

	t.Run("when play environment missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ListHeroMatchups(context.Background(), &ListHeroMatchupsRequest{PlayEnvironment: PlayEnvironmentUnspecified})
		require.Error(t, err)
	})
}

func TestClientListHeroLeaderboard(t *testing.T) {
	t.Run("when filters provided, then query encoded and leaderboard decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, "/v1/meta/hero-leaderboard", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "PLAY_ENVIRONMENT_PAPER", q.Get("playEnvironment"))
			require.Equal(t, "blitz", q.Get("ruleset"))
			require.Equal(t, "page-2", q.Get("pageToken"))
			require.Equal(t, []string{"evt-1", "evt-2"}, q["eventIds"])
			require.False(t, q.Has("sourceVisibility"))
		}, `{"heroes":[{"heroId":"hero-1","games":"40","opponentCount":"9","winRate":0.6,"bestMatchup":{"opponentHeroId":"hero-3","games":"5","wins":"4"}}],"metricDefinitionVersion":"v3"}`)

		resp, err := client.ListHeroLeaderboard(context.Background(), &ListHeroLeaderboardRequest{
			PlayEnvironment: PlayEnvironmentPaper,
			Ruleset:         "blitz",
			PageToken:       "page-2",
			EventIDs:        []string{"evt-1", "evt-2"},
		})
		require.NoError(t, err)
		require.Len(t, resp.Heroes, 1)
		require.Equal(t, uint64(9), resp.Heroes[0].OpponentCount)
		require.NotNil(t, resp.Heroes[0].BestMatchup)
		require.Equal(t, uint64(4), resp.Heroes[0].BestMatchup.Wins)
		require.Nil(t, resp.Heroes[0].WorstMatchup)
		require.Equal(t, "v3", resp.MetricDefinitionVersion)
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ListHeroLeaderboard(context.Background(), nil)
		require.Error(t, err)
	})
}

func TestClientListHeroMatchupMatrix(t *testing.T) {
	t.Run("when filters provided, then query encoded and matrix decoded", func(t *testing.T) {
		limit := int32(12)
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, "/v1/meta/hero-matchup-matrix", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "PLAY_ENVIRONMENT_DIGITAL", q.Get("playEnvironment"))
			require.Equal(t, "12", q.Get("heroLimit"))
			require.Equal(t, []string{"hero-1", "hero-2"}, q["heroIds"])
			require.Equal(t, []string{"src-1"}, q["sourceIds"])
		}, `{"matrices":[{"sourceId":"src-1","cohort":{"eligibleObservations":"10"},"heroes":[{"heroId":"hero-1","games":"10"}],"cells":[{"heroId":"hero-1","opponentHeroId":"hero-2","games":"4","winsGoingFirst":"2"}],"totalHeroes":30,"truncated":true}],"eventCoverage":[{"eventId":"evt-1","cohorts":[{"sourceId":"src-1","eligibleMatches":"7"}]}]}`)

		resp, err := client.ListHeroMatchupMatrix(context.Background(), &ListHeroMatchupMatrixRequest{
			PlayEnvironment: PlayEnvironmentDigital,
			HeroLimit:       &limit,
			HeroIDs:         []string{"hero-1", "hero-2"},
			SourceIDs:       []string{"src-1"},
		})
		require.NoError(t, err)
		require.Len(t, resp.Matrices, 1)
		matrix := resp.Matrices[0]
		require.Equal(t, uint32(30), matrix.TotalHeroes)
		require.True(t, matrix.Truncated)
		require.Equal(t, uint64(10), matrix.Cohort.EligibleObservations)
		require.Equal(t, uint64(2), matrix.Cells[0].WinsGoingFirst)
		require.Equal(t, uint64(7), resp.EventCoverage[0].Cohorts[0].EligibleMatches)
	})

	t.Run("when play environment missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ListHeroMatchupMatrix(context.Background(), &ListHeroMatchupMatrixRequest{})
		require.Error(t, err)
	})
}

func TestClientListCardUsage(t *testing.T) {
	t.Run("when filters provided, then repeated query values encoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, "/v1/meta/card-usage", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "PLAY_ENVIRONMENT_PAPER", q.Get("playEnvironment"))
			require.Equal(t, "cc", q.Get("ruleset"))
			require.Equal(t, []string{"ranked", "casual"}, q["competitionModes"])
			require.Equal(t, []string{"hero-1"}, q["heroIds"])
			require.Equal(t, []string{"Warrior"}, q["heroClasses"])
			require.Equal(t, []string{"Light"}, q["heroTalents"])
			require.Equal(t, []string{"card-1"}, q["cardIds"])
			require.Equal(t, []string{"Generic"}, q["cardClasses"])
			require.Equal(t, []string{"Shadow"}, q["cardTalents"])
			require.Equal(t, []string{"main", "side"}, q["zones"])
			require.Equal(t, "CARD_USAGE_ORDER_TOTAL_COPIES_DESC", q.Get("orderBy"))
			require.Equal(t, []string{"Action"}, q["cardTypes"])
			require.Equal(t, []string{"Attack"}, q["cardSubtypes"])
			require.Equal(t, []string{"Aura"}, q["excludedCardSubtypes"])
			require.Equal(t, "SOURCE_VISIBILITY_ALL", q.Get("sourceVisibility"))
		}, `{"cards":[{"cardId":"card-1","zone":"main","decksIncluding":"30","eligibleDeckObservations":"60","inclusionRate":0.5,"totalCopies":"90","evidenceType":"DECK_EVIDENCE_TYPE_REGISTERED_DECKLIST"}],"nextPageToken":"n"}`)

		resp, err := client.ListCardUsage(context.Background(), &ListCardUsageRequest{
			PlayEnvironment:      PlayEnvironmentPaper,
			Ruleset:              "cc",
			CompetitionModes:     []string{"ranked", "casual"},
			HeroIDs:              []string{"hero-1"},
			HeroClasses:          []string{"Warrior"},
			HeroTalents:          []string{"Light"},
			CardIDs:              []string{"card-1"},
			CardClasses:          []string{"Generic"},
			CardTalents:          []string{"Shadow"},
			Zones:                []string{"main", "side"},
			SourceVisibility:     SourceVisibilityAll,
			OrderBy:              CardUsageOrderTotalCopiesDesc,
			CardTypes:            []string{"Action"},
			CardSubtypes:         []string{"Attack"},
			ExcludedCardSubtypes: []string{"Aura"},
		})
		require.NoError(t, err)
		require.Len(t, resp.Cards, 1)
		require.Equal(t, uint64(30), resp.Cards[0].DecksIncluding)
		require.Equal(t, uint64(90), resp.Cards[0].TotalCopies)
		require.Equal(t, DeckEvidenceTypeRegisteredDecklist, resp.Cards[0].EvidenceType)
	})

	t.Run("when ruleset missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ListCardUsage(context.Background(), &ListCardUsageRequest{PlayEnvironment: PlayEnvironmentPaper})
		require.Error(t, err)
	})
}

func TestClientListCardHeroUsage(t *testing.T) {
	t.Run("when card provided, then query encoded and hero usage decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, "/v1/meta/card-hero-usage", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "card-1", q.Get("cardId"))
			require.Equal(t, "cc", q.Get("ruleset"))
			require.Equal(t, []string{"main"}, q["zones"])
		}, `{"heroUsage":[{"heroId":"hero-1","cardId":"card-1","decksIncluding":"3","averageCopiesWhenIncluded":2.5}]}`)

		resp, err := client.ListCardHeroUsage(context.Background(), &ListCardHeroUsageRequest{
			PlayEnvironment: PlayEnvironmentPaper,
			Ruleset:         "cc",
			CardID:          "card-1",
			Zones:           []string{"main"},
		})
		require.NoError(t, err)
		require.Len(t, resp.HeroUsage, 1)
		require.Equal(t, uint64(3), resp.HeroUsage[0].DecksIncluding)
		require.InDelta(t, 2.5, resp.HeroUsage[0].AverageCopiesWhenIncluded, 0.0001)
	})

	t.Run("when card id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ListCardHeroUsage(context.Background(), &ListCardHeroUsageRequest{PlayEnvironment: PlayEnvironmentPaper, Ruleset: "cc"})
		require.Error(t, err)
	})
}

func TestClientListHeroMetaShares(t *testing.T) {
	t.Run("when filters provided, then query encoded and shares decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, "/v1/meta/hero-meta-shares", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "PLAY_ENVIRONMENT_PAPER", q.Get("playEnvironment"))
			require.Equal(t, "HERO_SHARE_GRAIN_EVENT", q.Get("grain"))
			require.Equal(t, "evt-1", q.Get("eventId"))
			require.Equal(t, "hero-1", q.Get("heroId"))
		}, `{"shares":[{"heroId":"hero-1","games":"8","totalGames":"80","metaShare":0.1,"observationUnit":"OBSERVATION_UNIT_EVENT_ENTRANT","periodStart":"2026-10-01","periodEnd":"2026-10-02","eventId":"evt-1"}],"eligibleObservations":"80","observationUnit":"OBSERVATION_UNIT_EVENT_ENTRANT"}`)

		resp, err := client.ListHeroMetaShares(context.Background(), &ListHeroMetaSharesRequest{
			PlayEnvironment: PlayEnvironmentPaper,
			HeroID:          "hero-1",
			EventID:         "evt-1",
			Grain:           HeroShareGrainEvent,
		})
		require.NoError(t, err)
		require.Len(t, resp.Shares, 1)
		require.Equal(t, uint64(80), resp.Shares[0].TotalGames)
		require.Equal(t, ObservationUnitEventEntrant, resp.ObservationUnit)
		require.Equal(t, uint64(80), resp.EligibleObservations)
	})

	t.Run("when play environment missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		_, err := client.ListHeroMetaShares(context.Background(), &ListHeroMetaSharesRequest{})
		require.Error(t, err)
	})
}

func TestClientGetPaperEventResults(t *testing.T) {
	t.Run("when event id provided, then results decoded", func(t *testing.T) {
		client := newStubAPIClient(t, func(t *testing.T, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/meta/events/evt-1/paper-results", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "src-1", q.Get("sourceId"))
			require.Equal(t, []string{"scope-1", "scope-2"}, q["scopeIds"])
		}, `{"results":{"schemaVersion":1,"eventId":"evt-1","revision":{"sourceId":"src-1","sourceRevision":"r1","capturedAt":"2026-11-01T10:00:00Z"},
			"scopes":[{"scopeId":"scope-1","kind":"PAPER_RESULTS_SCOPE_KIND_EVENT_FINAL_STANDINGS","ordinal":0}],
			"standings":[{"standingId":"st-1","scopeId":"scope-1","competitorKind":"PAPER_COMPETITOR_KIND_INDIVIDUAL","reportedRank":1,"tied":false,"attributionStatus":"PAPER_ATTRIBUTION_STATUS_VERIFIED"}],
			"decks":[{"deckId":"deck-1","cards":[{"sourceCardId":"x","cardId":"card-1","copies":3,"zone":"main"}]}],
			"analyticsPopulations":[{"populationId":"pop-1","unit":"PAPER_ANALYTICS_POPULATION_UNIT_GAME","coverage":{"observedObservations":"12","publishedCoverage":"PAPER_PUBLISHED_COVERAGE_NONE"},"wins":"6","winRate":0.5}],
			"availability":[{"capability":"PAPER_RESULTS_CAPABILITY_STANDINGS","processing":"PAPER_PROCESSING_STATE_COMPLETE","stale":false,"standingExtent":{"kind":"PAPER_STANDING_EXTENT_KIND_TOP_N","rankLimit":8},"counts":{"observedRows":8,"publishedRows":8}}]}}`)

		resp, err := client.GetPaperEventResults(context.Background(), &GetPaperEventResultsRequest{
			EventID:  "evt-1",
			SourceID: "src-1",
			ScopeIDs: []string{"scope-1", "scope-2"},
		})
		require.NoError(t, err)
		require.NotNil(t, resp.Results)
		results := resp.Results
		require.Equal(t, uint32(1), results.SchemaVersion)
		require.Equal(t, "r1", results.Revision.SourceRevision)
		require.NotNil(t, results.Scopes[0].Ordinal)
		require.Equal(t, PaperResultsScopeKindEventFinalStandings, results.Scopes[0].Kind)
		require.Equal(t, uint32(1), results.Standings[0].ReportedRank)
		require.NotNil(t, results.Standings[0].Tied)
		require.False(t, *results.Standings[0].Tied)
		require.Equal(t, int32(3), results.Decks[0].Cards[0].Copies)
		require.Equal(t, uint64(12), results.AnalyticsPopulations[0].Coverage.ObservedObservations)
		require.Equal(t, uint64(6), results.AnalyticsPopulations[0].Wins)
		require.NotNil(t, results.AnalyticsPopulations[0].WinRate)
		availability := results.Availability[0]
		require.Equal(t, PaperProcessingStateComplete, availability.Processing)
		require.Equal(t, uint32(8), *availability.StandingExtent.RankLimit)
		require.Equal(t, uint32(8), availability.Counts.PublishedRows)
	})

	t.Run("when event id missing, then returns error", func(t *testing.T) {
		client := newTestClient(t)
		resp, err := client.GetPaperEventResults(context.Background(), &GetPaperEventResultsRequest{EventID: " "})
		require.Error(t, err)
		require.Nil(t, resp)
	})
}
