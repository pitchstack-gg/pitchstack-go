package client

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// PlayEnvironment represents meta.v1.PlayEnvironment.
type PlayEnvironment string

const (
	PlayEnvironmentUnspecified PlayEnvironment = "PLAY_ENVIRONMENT_UNSPECIFIED"
	PlayEnvironmentPaper       PlayEnvironment = "PLAY_ENVIRONMENT_PAPER"
	PlayEnvironmentDigital     PlayEnvironment = "PLAY_ENVIRONMENT_DIGITAL"
)

// SourceVisibility represents meta.v1.SourceVisibility.
type SourceVisibility string

const (
	SourceVisibilityUnspecified SourceVisibility = "SOURCE_VISIBILITY_UNSPECIFIED"
	SourceVisibilityPublic      SourceVisibility = "SOURCE_VISIBILITY_PUBLIC"
	SourceVisibilityPrivate     SourceVisibility = "SOURCE_VISIBILITY_PRIVATE"
	SourceVisibilityAll         SourceVisibility = "SOURCE_VISIBILITY_ALL"
)

// PlayDataType represents meta.v1.PlayDataType.
type PlayDataType string

const (
	PlayDataTypeUnspecified           PlayDataType = "PLAY_DATA_TYPE_UNSPECIFIED"
	PlayDataTypeMatchResult           PlayDataType = "PLAY_DATA_TYPE_MATCH_RESULT"
	PlayDataTypeEventEntrant          PlayDataType = "PLAY_DATA_TYPE_EVENT_ENTRANT"
	PlayDataTypeRegisteredDecklist    PlayDataType = "PLAY_DATA_TYPE_REGISTERED_DECKLIST"
	PlayDataTypeObservedCardTelemetry PlayDataType = "PLAY_DATA_TYPE_OBSERVED_CARD_TELEMETRY"
	PlayDataTypeEventProgression      PlayDataType = "PLAY_DATA_TYPE_EVENT_PROGRESSION"
)

// AnalyticsCapability represents meta.v1.AnalyticsCapability.
type AnalyticsCapability string

const (
	AnalyticsCapabilityUnspecified       AnalyticsCapability = "ANALYTICS_CAPABILITY_UNSPECIFIED"
	AnalyticsCapabilityHeroMatchups      AnalyticsCapability = "ANALYTICS_CAPABILITY_HERO_MATCHUPS"
	AnalyticsCapabilityHeroMetaShares    AnalyticsCapability = "ANALYTICS_CAPABILITY_HERO_META_SHARES"
	AnalyticsCapabilityCardUsage         AnalyticsCapability = "ANALYTICS_CAPABILITY_CARD_USAGE"
	AnalyticsCapabilityHeroLeaderboard   AnalyticsCapability = "ANALYTICS_CAPABILITY_HERO_LEADERBOARD"
	AnalyticsCapabilityHeroMatchupMatrix AnalyticsCapability = "ANALYTICS_CAPABILITY_HERO_MATCHUP_MATRIX"
	AnalyticsCapabilityEventStandings    AnalyticsCapability = "ANALYTICS_CAPABILITY_EVENT_STANDINGS"
	AnalyticsCapabilityEventConversions  AnalyticsCapability = "ANALYTICS_CAPABILITY_EVENT_CONVERSIONS"
)

// ObservationUnit represents meta.v1.ObservationUnit.
type ObservationUnit string

const (
	ObservationUnitUnspecified        ObservationUnit = "OBSERVATION_UNIT_UNSPECIFIED"
	ObservationUnitMatch              ObservationUnit = "OBSERVATION_UNIT_MATCH"
	ObservationUnitMatchParticipation ObservationUnit = "OBSERVATION_UNIT_MATCH_PARTICIPATION"
	ObservationUnitEventEntrant       ObservationUnit = "OBSERVATION_UNIT_EVENT_ENTRANT"
	ObservationUnitDeckSample         ObservationUnit = "OBSERVATION_UNIT_DECK_SAMPLE"
)

// DeckEvidenceType represents meta.v1.DeckEvidenceType.
type DeckEvidenceType string

const (
	DeckEvidenceTypeUnspecified        DeckEvidenceType = "DECK_EVIDENCE_TYPE_UNSPECIFIED"
	DeckEvidenceTypeTelemetryObserved  DeckEvidenceType = "DECK_EVIDENCE_TYPE_TELEMETRY_OBSERVED"
	DeckEvidenceTypeRegisteredDecklist DeckEvidenceType = "DECK_EVIDENCE_TYPE_REGISTERED_DECKLIST"
)

// HeroConversionBreakdownStatus represents meta.v1.HeroConversionBreakdownStatus.
type HeroConversionBreakdownStatus string

const (
	HeroConversionBreakdownStatusUnspecified                   HeroConversionBreakdownStatus = "HERO_CONVERSION_BREAKDOWN_STATUS_UNSPECIFIED"
	HeroConversionBreakdownStatusAvailable                     HeroConversionBreakdownStatus = "HERO_CONVERSION_BREAKDOWN_STATUS_AVAILABLE"
	HeroConversionBreakdownStatusUnavailableMixedOrLimited     HeroConversionBreakdownStatus = "HERO_CONVERSION_BREAKDOWN_STATUS_UNAVAILABLE_MIXED_OR_LIMITED"
	HeroConversionBreakdownStatusUnavailableIncompleteCoverage HeroConversionBreakdownStatus = "HERO_CONVERSION_BREAKDOWN_STATUS_UNAVAILABLE_INCOMPLETE_COVERAGE"
)

// EventConversionMilestoneKind represents meta.v1.EventConversionMilestoneKind.
type EventConversionMilestoneKind string

const (
	EventConversionMilestoneKindUnspecified   EventConversionMilestoneKind = "EVENT_CONVERSION_MILESTONE_KIND_UNSPECIFIED"
	EventConversionMilestoneKindObservedStage EventConversionMilestoneKind = "EVENT_CONVERSION_MILESTONE_KIND_OBSERVED_STAGE"
	EventConversionMilestoneKindPlacementCut  EventConversionMilestoneKind = "EVENT_CONVERSION_MILESTONE_KIND_PLACEMENT_CUT"
)

// HeroMatchupGrouping represents meta.v1.HeroMatchupGrouping.
type HeroMatchupGrouping string

const (
	HeroMatchupGroupingUnspecified   HeroMatchupGrouping = "HERO_MATCHUP_GROUPING_UNSPECIFIED"
	HeroMatchupGroupingSourcePairing HeroMatchupGrouping = "HERO_MATCHUP_GROUPING_SOURCE_PAIRING"
)

// HeroMatchupPairingOrder represents meta.v1.HeroMatchupPairingOrder.
type HeroMatchupPairingOrder string

const (
	HeroMatchupPairingOrderUnspecified HeroMatchupPairingOrder = "HERO_MATCHUP_PAIRING_ORDER_UNSPECIFIED"
	HeroMatchupPairingOrderGamesDesc   HeroMatchupPairingOrder = "HERO_MATCHUP_PAIRING_ORDER_GAMES_DESC"
	HeroMatchupPairingOrderWinRateDesc HeroMatchupPairingOrder = "HERO_MATCHUP_PAIRING_ORDER_WIN_RATE_DESC"
)

// HeroShareGrain represents meta.v1.HeroShareGrain.
type HeroShareGrain string

const (
	HeroShareGrainUnspecified HeroShareGrain = "HERO_SHARE_GRAIN_UNSPECIFIED"
	HeroShareGrainTotal       HeroShareGrain = "HERO_SHARE_GRAIN_TOTAL"
	HeroShareGrainDay         HeroShareGrain = "HERO_SHARE_GRAIN_DAY"
	HeroShareGrainEvent       HeroShareGrain = "HERO_SHARE_GRAIN_EVENT"
)

// CardUsageOrder represents meta.v1.CardUsageOrder.
type CardUsageOrder string

const (
	CardUsageOrderUnspecified        CardUsageOrder = "CARD_USAGE_ORDER_UNSPECIFIED"
	CardUsageOrderInclusionRateDesc  CardUsageOrder = "CARD_USAGE_ORDER_INCLUSION_RATE_DESC"
	CardUsageOrderDecksIncludingDesc CardUsageOrder = "CARD_USAGE_ORDER_DECKS_INCLUDING_DESC"
	CardUsageOrderTotalCopiesDesc    CardUsageOrder = "CARD_USAGE_ORDER_TOTAL_COPIES_DESC"
)

// MetaSource represents meta.v1.MetaSource.
type MetaSource struct {
	SourceID              string                `json:"sourceId,omitempty"`
	DisplayName           string                `json:"displayName,omitempty"`
	PlayEnvironment       PlayEnvironment       `json:"playEnvironment,omitempty"`
	DataTypes             []PlayDataType        `json:"dataTypes,omitempty"`
	Active                bool                  `json:"active,omitempty"`
	AnalyticsCapabilities []AnalyticsCapability `json:"analyticsCapabilities,omitempty"`
}

// PlayDataCohortSummary represents meta.v1.PlayDataCohortSummary.
type PlayDataCohortSummary struct {
	SourceID                     string           `json:"sourceId,omitempty"`
	PlayEnvironment              PlayEnvironment  `json:"playEnvironment,omitempty"`
	ObservationUnit              ObservationUnit  `json:"observationUnit,omitempty"`
	EligibleObservations         uint64           `json:"eligibleObservations,omitempty,string"`
	DataThrough                  *time.Time       `json:"dataThrough,omitempty"`
	MetricDefinitionVersion      string           `json:"metricDefinitionVersion,omitempty"`
	DeckEvidenceType             DeckEvidenceType `json:"deckEvidenceType,omitempty"`
	KnownFirstPlayerObservations uint64           `json:"knownFirstPlayerObservations,omitempty,string"`
	ObservedObservations         *uint64          `json:"observedObservations,omitempty,string"`
	ExcludedObservations         *uint64          `json:"excludedObservations,omitempty,string"`
}

// EventStanding represents meta.v1.EventStanding.
type EventStanding struct {
	Place          uint32 `json:"place,omitempty"`
	PlayerName     string `json:"playerName,omitempty"`
	HeroID         string `json:"heroId,omitempty"`
	SourceHeroName string `json:"sourceHeroName,omitempty"`
}

// EventHeroConversion represents meta.v1.EventHeroConversion.
type EventHeroConversion struct {
	HeroID            string  `json:"heroId,omitempty"`
	StartingEntrants  uint32  `json:"startingEntrants,omitempty"`
	ConvertedEntrants uint32  `json:"convertedEntrants,omitempty"`
	ConversionRate    float64 `json:"conversionRate,omitempty"`
}

// EventConversionMilestone represents meta.v1.EventConversionMilestone.
type EventConversionMilestone struct {
	Key                         string                       `json:"key,omitempty"`
	DisplayLabel                string                       `json:"displayLabel,omitempty"`
	Kind                        EventConversionMilestoneKind `json:"kind,omitempty"`
	Ordinal                     uint32                       `json:"ordinal,omitempty"`
	SourceRule                  string                       `json:"sourceRule,omitempty"`
	SourceRound                 *uint32                      `json:"sourceRound,omitempty"`
	PlacementCutoff             *uint32                      `json:"placementCutoff,omitempty"`
	ConvertedEntrants           uint32                       `json:"convertedEntrants,omitempty"`
	FieldConversionRate         float64                      `json:"fieldConversionRate,omitempty"`
	PreviousStageConversionRate *float64                     `json:"previousStageConversionRate,omitempty"`
	PreviousMilestoneKey        string                       `json:"previousMilestoneKey,omitempty"`
	Heroes                      []EventHeroConversion        `json:"heroes,omitempty"`
}

// HeroMatchupStatistic represents meta.v1.HeroMatchupStatistic.
type HeroMatchupStatistic struct {
	HeroID                string           `json:"heroId,omitempty"`
	OpponentHeroID        string           `json:"opponentHeroId,omitempty"`
	PlayEnvironment       PlayEnvironment  `json:"playEnvironment,omitempty"`
	Platform              string           `json:"platform,omitempty"`     // Deprecated: use SourceID.
	SourceSystem          string           `json:"sourceSystem,omitempty"` // Deprecated: use SourceID.
	Ruleset               string           `json:"ruleset,omitempty"`
	CompetitionMode       string           `json:"competitionMode,omitempty"`
	SourceVisibility      SourceVisibility `json:"sourceVisibility,omitempty"`
	Games                 uint64           `json:"games,omitempty,string"`
	Wins                  uint64           `json:"wins,omitempty,string"`
	Losses                uint64           `json:"losses,omitempty,string"`
	Draws                 uint64           `json:"draws,omitempty,string"`
	ConcededGames         uint64           `json:"concededGames,omitempty,string"`
	KnownFirstPlayerGames uint64           `json:"knownFirstPlayerGames,omitempty,string"`
	WinsGoingFirst        uint64           `json:"winsGoingFirst,omitempty,string"`
	WinsGoingSecond       uint64           `json:"winsGoingSecond,omitempty,string"`
	WinRate               float64          `json:"winRate,omitempty"`
	WinRateLower95        float64          `json:"winRateLower95,omitempty"`
	WinRateUpper95        float64          `json:"winRateUpper95,omitempty"`
	LowSample             bool             `json:"lowSample,omitempty"`
	KnownConcessionGames  uint64           `json:"knownConcessionGames,omitempty,string"`
	GamesGoingFirst       uint64           `json:"gamesGoingFirst,omitempty,string"`
	GamesGoingSecond      uint64           `json:"gamesGoingSecond,omitempty,string"`
	SourceID              string           `json:"sourceId,omitempty"`
}

// HeroLeaderboardMatchup represents meta.v1.HeroLeaderboardMatchup.
type HeroLeaderboardMatchup struct {
	OpponentHeroID string  `json:"opponentHeroId,omitempty"`
	Games          uint64  `json:"games,omitempty,string"`
	Wins           uint64  `json:"wins,omitempty,string"`
	WinRate        float64 `json:"winRate,omitempty"`
	WinRateLower95 float64 `json:"winRateLower95,omitempty"`
	WinRateUpper95 float64 `json:"winRateUpper95,omitempty"`
}

// HeroLeaderboardStatistic represents meta.v1.HeroLeaderboardStatistic.
type HeroLeaderboardStatistic struct {
	SourceID        string                  `json:"sourceId,omitempty"`
	HeroID          string                  `json:"heroId,omitempty"`
	PlayEnvironment PlayEnvironment         `json:"playEnvironment,omitempty"`
	Games           uint64                  `json:"games,omitempty,string"`
	Wins            uint64                  `json:"wins,omitempty,string"`
	Losses          uint64                  `json:"losses,omitempty,string"`
	Draws           uint64                  `json:"draws,omitempty,string"`
	OpponentCount   uint64                  `json:"opponentCount,omitempty,string"`
	WinRate         float64                 `json:"winRate,omitempty"`
	WinRateLower95  float64                 `json:"winRateLower95,omitempty"`
	WinRateUpper95  float64                 `json:"winRateUpper95,omitempty"`
	LowSample       bool                    `json:"lowSample,omitempty"`
	BestMatchup     *HeroLeaderboardMatchup `json:"bestMatchup,omitempty"`
	WorstMatchup    *HeroLeaderboardMatchup `json:"worstMatchup,omitempty"`
}

// HeroMatchupMatrixHero represents meta.v1.HeroMatchupMatrixHero.
type HeroMatchupMatrixHero struct {
	HeroID         string  `json:"heroId,omitempty"`
	Games          uint64  `json:"games,omitempty,string"`
	Wins           uint64  `json:"wins,omitempty,string"`
	Losses         uint64  `json:"losses,omitempty,string"`
	Draws          uint64  `json:"draws,omitempty,string"`
	OpponentCount  uint64  `json:"opponentCount,omitempty,string"`
	WinRate        float64 `json:"winRate,omitempty"`
	WinRateLower95 float64 `json:"winRateLower95,omitempty"`
	WinRateUpper95 float64 `json:"winRateUpper95,omitempty"`
	LowSample      bool    `json:"lowSample,omitempty"`
}

// HeroMatchupMatrixCell represents meta.v1.HeroMatchupMatrixCell.
type HeroMatchupMatrixCell struct {
	HeroID                string  `json:"heroId,omitempty"`
	OpponentHeroID        string  `json:"opponentHeroId,omitempty"`
	Games                 uint64  `json:"games,omitempty,string"`
	Wins                  uint64  `json:"wins,omitempty,string"`
	Losses                uint64  `json:"losses,omitempty,string"`
	Draws                 uint64  `json:"draws,omitempty,string"`
	ConcededGames         uint64  `json:"concededGames,omitempty,string"`
	KnownConcessionGames  uint64  `json:"knownConcessionGames,omitempty,string"`
	KnownFirstPlayerGames uint64  `json:"knownFirstPlayerGames,omitempty,string"`
	GamesGoingFirst       uint64  `json:"gamesGoingFirst,omitempty,string"`
	GamesGoingSecond      uint64  `json:"gamesGoingSecond,omitempty,string"`
	WinsGoingFirst        uint64  `json:"winsGoingFirst,omitempty,string"`
	WinsGoingSecond       uint64  `json:"winsGoingSecond,omitempty,string"`
	WinRate               float64 `json:"winRate,omitempty"`
	WinRateLower95        float64 `json:"winRateLower95,omitempty"`
	WinRateUpper95        float64 `json:"winRateUpper95,omitempty"`
	LowSample             bool    `json:"lowSample,omitempty"`
}

// HeroMatchupMatrix represents meta.v1.HeroMatchupMatrix.
type HeroMatchupMatrix struct {
	SourceID        string                  `json:"sourceId,omitempty"`
	PlayEnvironment PlayEnvironment         `json:"playEnvironment,omitempty"`
	Ruleset         string                  `json:"ruleset,omitempty"`
	Cohort          *PlayDataCohortSummary  `json:"cohort,omitempty"`
	Heroes          []HeroMatchupMatrixHero `json:"heroes,omitempty"`
	Cells           []HeroMatchupMatrixCell `json:"cells,omitempty"`
	TotalHeroes     uint32                  `json:"totalHeroes,omitempty"`
	Truncated       bool                    `json:"truncated,omitempty"`
}

// EventMatchupCoverageCohort represents meta.v1.EventMatchupCoverageCohort.
type EventMatchupCoverageCohort struct {
	SourceID                string     `json:"sourceId,omitempty"`
	Ruleset                 string     `json:"ruleset,omitempty"`
	EligibleMatches         uint64     `json:"eligibleMatches,omitempty,string"`
	KnownFirstPlayerMatches uint64     `json:"knownFirstPlayerMatches,omitempty,string"`
	DataThrough             *time.Time `json:"dataThrough,omitempty"`
}

// EventMatchupCoverage represents meta.v1.EventMatchupCoverage.
type EventMatchupCoverage struct {
	EventID string                       `json:"eventId,omitempty"`
	Cohorts []EventMatchupCoverageCohort `json:"cohorts,omitempty"`
}

// CardUsageStatistic represents meta.v1.CardUsageStatistic.
type CardUsageStatistic struct {
	CardID                       string           `json:"cardId,omitempty"`
	Zone                         string           `json:"zone,omitempty"`
	PlayEnvironment              PlayEnvironment  `json:"playEnvironment,omitempty"`
	SourceID                     string           `json:"sourceId,omitempty"`
	DecksIncluding               uint64           `json:"decksIncluding,omitempty,string"`
	EligibleDeckObservations     uint64           `json:"eligibleDeckObservations,omitempty,string"`
	InclusionRate                float64          `json:"inclusionRate,omitempty"`
	TotalCopies                  uint64           `json:"totalCopies,omitempty,string"`
	AverageCopiesWhenIncluded    float64          `json:"averageCopiesWhenIncluded,omitempty"`
	AverageCopiesPerEligibleDeck float64          `json:"averageCopiesPerEligibleDeck,omitempty"`
	LowSample                    bool             `json:"lowSample,omitempty"`
	EvidenceType                 DeckEvidenceType `json:"evidenceType,omitempty"`
}

// CardHeroUsageStatistic represents meta.v1.CardHeroUsageStatistic.
type CardHeroUsageStatistic struct {
	HeroID                       string           `json:"heroId,omitempty"`
	CardID                       string           `json:"cardId,omitempty"`
	Zone                         string           `json:"zone,omitempty"`
	PlayEnvironment              PlayEnvironment  `json:"playEnvironment,omitempty"`
	SourceID                     string           `json:"sourceId,omitempty"`
	DecksIncluding               uint64           `json:"decksIncluding,omitempty,string"`
	EligibleDeckObservations     uint64           `json:"eligibleDeckObservations,omitempty,string"`
	InclusionRate                float64          `json:"inclusionRate,omitempty"`
	TotalCopies                  uint64           `json:"totalCopies,omitempty,string"`
	AverageCopiesWhenIncluded    float64          `json:"averageCopiesWhenIncluded,omitempty"`
	AverageCopiesPerEligibleDeck float64          `json:"averageCopiesPerEligibleDeck,omitempty"`
	LowSample                    bool             `json:"lowSample,omitempty"`
	EvidenceType                 DeckEvidenceType `json:"evidenceType,omitempty"`
}

// HeroMetaShareStatistic represents meta.v1.HeroMetaShareStatistic.
type HeroMetaShareStatistic struct {
	MatchDate         string           `json:"matchDate,omitempty"`
	HeroID            string           `json:"heroId,omitempty"`
	PlayEnvironment   PlayEnvironment  `json:"playEnvironment,omitempty"`
	Platform          string           `json:"platform,omitempty"`     // Deprecated: use SourceID.
	SourceSystem      string           `json:"sourceSystem,omitempty"` // Deprecated: use SourceID.
	Ruleset           string           `json:"ruleset,omitempty"`
	CompetitionMode   string           `json:"competitionMode,omitempty"`
	SourceVisibility  SourceVisibility `json:"sourceVisibility,omitempty"`
	Games             uint64           `json:"games,omitempty,string"`
	TotalGames        uint64           `json:"totalGames,omitempty,string"`
	MetaShare         float64          `json:"metaShare,omitempty"`
	SourceID          string           `json:"sourceId,omitempty"`
	ObservationUnit   ObservationUnit  `json:"observationUnit,omitempty"`
	Observations      uint64           `json:"observations,omitempty,string"`
	TotalObservations uint64           `json:"totalObservations,omitempty,string"`
	PeriodStart       string           `json:"periodStart,omitempty"`
	PeriodEnd         string           `json:"periodEnd,omitempty"`
	EventID           string           `json:"eventId,omitempty"`
}

// ListMetaSourcesRequest configures meta source listing.
type ListMetaSourcesRequest struct {
	IncludeInactive bool
}

// ListMetaSourcesResponse returns meta data sources.
type ListMetaSourcesResponse struct {
	Sources  []MetaSource     `json:"sources,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ListMetaSourcesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListEventStandingsRequest identifies an event for meta standings.
type ListEventStandingsRequest struct {
	EventID  string `json:"-"`
	SourceID string
}

// ListEventStandingsResponse returns event standings.
type ListEventStandingsResponse struct {
	EventID      string           `json:"eventId,omitempty"`
	SourceID     string           `json:"sourceId,omitempty"`
	Standings    []EventStanding  `json:"standings,omitempty"`
	TotalPlayers uint32           `json:"totalPlayers,omitempty"`
	DataThrough  *time.Time       `json:"dataThrough,omitempty"`
	Metadata     ResponseMetadata `json:"-"`
}

func (r *ListEventStandingsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListEventConversionsRequest identifies an event for conversion analytics.
type ListEventConversionsRequest struct {
	EventID  string `json:"-"`
	SourceID string
}

// ListEventConversionsResponse returns event conversion milestones.
type ListEventConversionsResponse struct {
	EventID                 string                        `json:"eventId,omitempty"`
	SourceID                string                        `json:"sourceId,omitempty"`
	StartingEntrants        uint32                        `json:"startingEntrants,omitempty"`
	HeroEligibleEntrants    uint32                        `json:"heroEligibleEntrants,omitempty"`
	HeroExcludedEntrants    uint32                        `json:"heroExcludedEntrants,omitempty"`
	HeroBreakdownStatus     HeroConversionBreakdownStatus `json:"heroBreakdownStatus,omitempty"`
	DataThrough             *time.Time                    `json:"dataThrough,omitempty"`
	MetricDefinitionVersion string                        `json:"metricDefinitionVersion,omitempty"`
	Milestones              []EventConversionMilestone    `json:"milestones,omitempty"`
	Metadata                ResponseMetadata              `json:"-"`
}

func (r *ListEventConversionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListHeroMatchupsRequest filters hero matchup statistics.
// StartDate and EndDate are inclusive YYYY-MM-DD dates.
type ListHeroMatchupsRequest struct {
	PlayEnvironment  PlayEnvironment
	StartDate        string
	EndDate          string
	Ruleset          string
	CompetitionMode  string
	HeroID           string
	OpponentHeroID   string
	Platform         string // Deprecated: use SourceIDs.
	SourceSystem     string // Deprecated: use SourceIDs.
	SourceVisibility SourceVisibility
	PageSize         *int32
	NextPageToken    string
	EventID          string
	SourceIDs        []string
	Grouping         HeroMatchupGrouping
	OrderBy          HeroMatchupPairingOrder
	EventIDs         []string
}

// ListHeroMatchupsResponse returns hero matchup statistics.
type ListHeroMatchupsResponse struct {
	Matchups                []HeroMatchupStatistic  `json:"matchups,omitempty"`
	DataThrough             *time.Time              `json:"dataThrough,omitempty"`
	EligibleMatches         uint64                  `json:"eligibleMatches,omitempty,string"`
	KnownFirstPlayerMatches uint64                  `json:"knownFirstPlayerMatches,omitempty,string"`
	MetricDefinitionVersion string                  `json:"metricDefinitionVersion,omitempty"`
	NextPageToken           string                  `json:"nextPageToken,omitempty"`
	Cohorts                 []PlayDataCohortSummary `json:"cohorts,omitempty"`
	Metadata                ResponseMetadata        `json:"-"`
}

func (r *ListHeroMatchupsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListHeroLeaderboardRequest filters the hero leaderboard.
// StartDate and EndDate are inclusive YYYY-MM-DD dates.
type ListHeroLeaderboardRequest struct {
	PlayEnvironment  PlayEnvironment
	StartDate        string
	EndDate          string
	Ruleset          string
	CompetitionMode  string
	SourceVisibility SourceVisibility
	PageSize         *int32
	PageToken        string
	SourceIDs        []string
	EventIDs         []string
}

// ListHeroLeaderboardResponse returns hero leaderboard rows.
type ListHeroLeaderboardResponse struct {
	Heroes                  []HeroLeaderboardStatistic `json:"heroes,omitempty"`
	DataThrough             *time.Time                 `json:"dataThrough,omitempty"`
	MetricDefinitionVersion string                     `json:"metricDefinitionVersion,omitempty"`
	NextPageToken           string                     `json:"nextPageToken,omitempty"`
	Cohorts                 []PlayDataCohortSummary    `json:"cohorts,omitempty"`
	Metadata                ResponseMetadata           `json:"-"`
}

func (r *ListHeroLeaderboardResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListHeroMatchupMatrixRequest filters the hero matchup matrix.
// StartDate and EndDate are inclusive YYYY-MM-DD dates.
type ListHeroMatchupMatrixRequest struct {
	PlayEnvironment  PlayEnvironment
	StartDate        string
	EndDate          string
	Ruleset          string
	CompetitionMode  string
	SourceVisibility SourceVisibility
	SourceIDs        []string
	EventIDs         []string
	HeroLimit        *int32
	HeroIDs          []string
}

// ListHeroMatchupMatrixResponse returns hero matchup matrices.
type ListHeroMatchupMatrixResponse struct {
	Matrices                []HeroMatchupMatrix    `json:"matrices,omitempty"`
	EventCoverage           []EventMatchupCoverage `json:"eventCoverage,omitempty"`
	MetricDefinitionVersion string                 `json:"metricDefinitionVersion,omitempty"`
	Metadata                ResponseMetadata       `json:"-"`
}

func (r *ListHeroMatchupMatrixResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListCardUsageRequest filters card usage statistics.
// StartDate and EndDate are inclusive YYYY-MM-DD dates.
type ListCardUsageRequest struct {
	PlayEnvironment      PlayEnvironment
	StartDate            string
	EndDate              string
	Ruleset              string
	CompetitionModes     []string
	HeroIDs              []string
	HeroClasses          []string
	HeroTalents          []string
	CardIDs              []string
	CardClasses          []string
	CardTalents          []string
	Zones                []string
	SourceVisibility     SourceVisibility
	PageSize             *int32
	PageToken            string
	SourceIDs            []string
	OrderBy              CardUsageOrder
	CardTypes            []string
	CardSubtypes         []string
	ExcludedCardSubtypes []string
	EventIDs             []string
}

// ListCardUsageResponse returns card usage statistics.
type ListCardUsageResponse struct {
	Cards         []CardUsageStatistic    `json:"cards,omitempty"`
	Cohorts       []PlayDataCohortSummary `json:"cohorts,omitempty"`
	NextPageToken string                  `json:"nextPageToken,omitempty"`
	Metadata      ResponseMetadata        `json:"-"`
}

func (r *ListCardUsageResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListCardHeroUsageRequest filters per-hero usage for one card.
// StartDate and EndDate are inclusive YYYY-MM-DD dates.
type ListCardHeroUsageRequest struct {
	PlayEnvironment  PlayEnvironment
	StartDate        string
	EndDate          string
	Ruleset          string
	CardID           string
	CompetitionModes []string
	Zones            []string
	SourceVisibility SourceVisibility
	PageSize         *int32
	PageToken        string
	SourceIDs        []string
	EventIDs         []string
}

// ListCardHeroUsageResponse returns per-hero usage for one card.
type ListCardHeroUsageResponse struct {
	HeroUsage     []CardHeroUsageStatistic `json:"heroUsage,omitempty"`
	Cohorts       []PlayDataCohortSummary  `json:"cohorts,omitempty"`
	NextPageToken string                   `json:"nextPageToken,omitempty"`
	Metadata      ResponseMetadata         `json:"-"`
}

func (r *ListCardHeroUsageResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListHeroMetaSharesRequest filters hero meta share statistics.
// StartDate and EndDate are inclusive YYYY-MM-DD dates.
type ListHeroMetaSharesRequest struct {
	PlayEnvironment  PlayEnvironment
	StartDate        string
	EndDate          string
	Ruleset          string
	CompetitionMode  string
	HeroID           string
	Platform         string // Deprecated: use SourceIDs.
	SourceSystem     string // Deprecated: use SourceIDs.
	SourceVisibility SourceVisibility
	PageSize         *int32
	NextPageToken    string
	SourceIDs        []string
	EventID          string
	Grain            HeroShareGrain
	EventIDs         []string
}

// ListHeroMetaSharesResponse returns hero meta share statistics.
type ListHeroMetaSharesResponse struct {
	Shares                  []HeroMetaShareStatistic `json:"shares,omitempty"`
	DataThrough             *time.Time               `json:"dataThrough,omitempty"`
	EligibleMatches         uint64                   `json:"eligibleMatches,omitempty,string"`
	MetricDefinitionVersion string                   `json:"metricDefinitionVersion,omitempty"`
	NextPageToken           string                   `json:"nextPageToken,omitempty"`
	EligibleObservations    uint64                   `json:"eligibleObservations,omitempty,string"`
	ObservationUnit         ObservationUnit          `json:"observationUnit,omitempty"`
	Cohorts                 []PlayDataCohortSummary  `json:"cohorts,omitempty"`
	Metadata                ResponseMetadata         `json:"-"`
}

func (r *ListHeroMetaSharesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListMetaSources lists meta analytics data sources.
func (c *Client) ListMetaSources(ctx context.Context, request *ListMetaSourcesRequest, opts ...RequestOpt) (*ListMetaSourcesResponse, error) {
	if request == nil {
		request = &ListMetaSourcesRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/meta/sources", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.IncludeInactive {
		query.Set("includeInactive", "true")
	}
	req.URL.RawQuery = query.Encode()

	response := &ListMetaSourcesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListEventStandings lists final standings for an event.
func (c *Client) ListEventStandings(ctx context.Context, request *ListEventStandingsRequest, opts ...RequestOpt) (*ListEventStandingsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	eventID := strings.TrimSpace(request.EventID)
	if eventID == "" {
		return nil, errors.New("eventID must not be empty")
	}

	path := fmt.Sprintf("/v1/meta/events/%s/standings", url.PathEscape(eventID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "sourceId", request.SourceID)
	req.URL.RawQuery = query.Encode()

	response := &ListEventStandingsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListEventConversions lists stage conversion milestones for an event.
func (c *Client) ListEventConversions(ctx context.Context, request *ListEventConversionsRequest, opts ...RequestOpt) (*ListEventConversionsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	eventID := strings.TrimSpace(request.EventID)
	if eventID == "" {
		return nil, errors.New("eventID must not be empty")
	}

	path := fmt.Sprintf("/v1/meta/events/%s/conversions", url.PathEscape(eventID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "sourceId", request.SourceID)
	req.URL.RawQuery = query.Encode()

	response := &ListEventConversionsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListHeroMatchups lists hero-versus-hero matchup statistics.
func (c *Client) ListHeroMatchups(ctx context.Context, request *ListHeroMatchupsRequest, opts ...RequestOpt) (*ListHeroMatchupsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := requirePlayEnvironment(request.PlayEnvironment); err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/meta/hero-matchups", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryEnum(query, "playEnvironment", string(request.PlayEnvironment), string(PlayEnvironmentUnspecified))
	setQueryString(query, "startDate", request.StartDate)
	setQueryString(query, "endDate", request.EndDate)
	setQueryString(query, "ruleset", request.Ruleset)
	setQueryString(query, "competitionMode", request.CompetitionMode)
	setQueryString(query, "heroId", request.HeroID)
	setQueryString(query, "opponentHeroId", request.OpponentHeroID)
	setQueryString(query, "platform", request.Platform)
	setQueryString(query, "sourceSystem", request.SourceSystem)
	setQueryEnum(query, "sourceVisibility", string(request.SourceVisibility), string(SourceVisibilityUnspecified))
	setQueryInt32(query, "pageSize", request.PageSize)
	setQueryString(query, "nextPageToken", request.NextPageToken)
	setQueryString(query, "eventId", request.EventID)
	setQueryStringList(query, "sourceIds", request.SourceIDs)
	setQueryEnum(query, "grouping", string(request.Grouping), string(HeroMatchupGroupingUnspecified))
	setQueryEnum(query, "orderBy", string(request.OrderBy), string(HeroMatchupPairingOrderUnspecified))
	setQueryStringList(query, "eventIds", request.EventIDs)
	req.URL.RawQuery = query.Encode()

	response := &ListHeroMatchupsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListHeroLeaderboard lists hero win-rate leaderboard rows.
func (c *Client) ListHeroLeaderboard(ctx context.Context, request *ListHeroLeaderboardRequest, opts ...RequestOpt) (*ListHeroLeaderboardResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := requirePlayEnvironment(request.PlayEnvironment); err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/meta/hero-leaderboard", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryEnum(query, "playEnvironment", string(request.PlayEnvironment), string(PlayEnvironmentUnspecified))
	setQueryString(query, "startDate", request.StartDate)
	setQueryString(query, "endDate", request.EndDate)
	setQueryString(query, "ruleset", request.Ruleset)
	setQueryString(query, "competitionMode", request.CompetitionMode)
	setQueryEnum(query, "sourceVisibility", string(request.SourceVisibility), string(SourceVisibilityUnspecified))
	setQueryInt32(query, "pageSize", request.PageSize)
	setQueryString(query, "pageToken", request.PageToken)
	setQueryStringList(query, "sourceIds", request.SourceIDs)
	setQueryStringList(query, "eventIds", request.EventIDs)
	req.URL.RawQuery = query.Encode()

	response := &ListHeroLeaderboardResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListHeroMatchupMatrix lists hero matchup matrices.
func (c *Client) ListHeroMatchupMatrix(ctx context.Context, request *ListHeroMatchupMatrixRequest, opts ...RequestOpt) (*ListHeroMatchupMatrixResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := requirePlayEnvironment(request.PlayEnvironment); err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/meta/hero-matchup-matrix", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryEnum(query, "playEnvironment", string(request.PlayEnvironment), string(PlayEnvironmentUnspecified))
	setQueryString(query, "startDate", request.StartDate)
	setQueryString(query, "endDate", request.EndDate)
	setQueryString(query, "ruleset", request.Ruleset)
	setQueryString(query, "competitionMode", request.CompetitionMode)
	setQueryEnum(query, "sourceVisibility", string(request.SourceVisibility), string(SourceVisibilityUnspecified))
	setQueryStringList(query, "sourceIds", request.SourceIDs)
	setQueryStringList(query, "eventIds", request.EventIDs)
	setQueryInt32(query, "heroLimit", request.HeroLimit)
	setQueryStringList(query, "heroIds", request.HeroIDs)
	req.URL.RawQuery = query.Encode()

	response := &ListHeroMatchupMatrixResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListCardUsage lists card inclusion statistics.
func (c *Client) ListCardUsage(ctx context.Context, request *ListCardUsageRequest, opts ...RequestOpt) (*ListCardUsageResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := requirePlayEnvironment(request.PlayEnvironment); err != nil {
		return nil, err
	}
	if strings.TrimSpace(request.Ruleset) == "" {
		return nil, errors.New("ruleset must not be empty")
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/meta/card-usage", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryEnum(query, "playEnvironment", string(request.PlayEnvironment), string(PlayEnvironmentUnspecified))
	setQueryString(query, "startDate", request.StartDate)
	setQueryString(query, "endDate", request.EndDate)
	setQueryString(query, "ruleset", request.Ruleset)
	setQueryStringList(query, "competitionModes", request.CompetitionModes)
	setQueryStringList(query, "heroIds", request.HeroIDs)
	setQueryStringList(query, "heroClasses", request.HeroClasses)
	setQueryStringList(query, "heroTalents", request.HeroTalents)
	setQueryStringList(query, "cardIds", request.CardIDs)
	setQueryStringList(query, "cardClasses", request.CardClasses)
	setQueryStringList(query, "cardTalents", request.CardTalents)
	setQueryStringList(query, "zones", request.Zones)
	setQueryEnum(query, "sourceVisibility", string(request.SourceVisibility), string(SourceVisibilityUnspecified))
	setQueryInt32(query, "pageSize", request.PageSize)
	setQueryString(query, "pageToken", request.PageToken)
	setQueryStringList(query, "sourceIds", request.SourceIDs)
	setQueryEnum(query, "orderBy", string(request.OrderBy), string(CardUsageOrderUnspecified))
	setQueryStringList(query, "cardTypes", request.CardTypes)
	setQueryStringList(query, "cardSubtypes", request.CardSubtypes)
	setQueryStringList(query, "excludedCardSubtypes", request.ExcludedCardSubtypes)
	setQueryStringList(query, "eventIds", request.EventIDs)
	req.URL.RawQuery = query.Encode()

	response := &ListCardUsageResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListCardHeroUsage lists per-hero inclusion statistics for one card.
func (c *Client) ListCardHeroUsage(ctx context.Context, request *ListCardHeroUsageRequest, opts ...RequestOpt) (*ListCardHeroUsageResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := requirePlayEnvironment(request.PlayEnvironment); err != nil {
		return nil, err
	}
	if strings.TrimSpace(request.Ruleset) == "" {
		return nil, errors.New("ruleset must not be empty")
	}
	if strings.TrimSpace(request.CardID) == "" {
		return nil, errors.New("cardID must not be empty")
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/meta/card-hero-usage", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryEnum(query, "playEnvironment", string(request.PlayEnvironment), string(PlayEnvironmentUnspecified))
	setQueryString(query, "startDate", request.StartDate)
	setQueryString(query, "endDate", request.EndDate)
	setQueryString(query, "ruleset", request.Ruleset)
	setQueryString(query, "cardId", request.CardID)
	setQueryStringList(query, "competitionModes", request.CompetitionModes)
	setQueryStringList(query, "zones", request.Zones)
	setQueryEnum(query, "sourceVisibility", string(request.SourceVisibility), string(SourceVisibilityUnspecified))
	setQueryInt32(query, "pageSize", request.PageSize)
	setQueryString(query, "pageToken", request.PageToken)
	setQueryStringList(query, "sourceIds", request.SourceIDs)
	setQueryStringList(query, "eventIds", request.EventIDs)
	req.URL.RawQuery = query.Encode()

	response := &ListCardHeroUsageResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListHeroMetaShares lists hero meta share statistics.
func (c *Client) ListHeroMetaShares(ctx context.Context, request *ListHeroMetaSharesRequest, opts ...RequestOpt) (*ListHeroMetaSharesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if err := requirePlayEnvironment(request.PlayEnvironment); err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/meta/hero-meta-shares", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryEnum(query, "playEnvironment", string(request.PlayEnvironment), string(PlayEnvironmentUnspecified))
	setQueryString(query, "startDate", request.StartDate)
	setQueryString(query, "endDate", request.EndDate)
	setQueryString(query, "ruleset", request.Ruleset)
	setQueryString(query, "competitionMode", request.CompetitionMode)
	setQueryString(query, "heroId", request.HeroID)
	setQueryString(query, "platform", request.Platform)
	setQueryString(query, "sourceSystem", request.SourceSystem)
	setQueryEnum(query, "sourceVisibility", string(request.SourceVisibility), string(SourceVisibilityUnspecified))
	setQueryInt32(query, "pageSize", request.PageSize)
	setQueryString(query, "nextPageToken", request.NextPageToken)
	setQueryStringList(query, "sourceIds", request.SourceIDs)
	setQueryString(query, "eventId", request.EventID)
	setQueryEnum(query, "grain", string(request.Grain), string(HeroShareGrainUnspecified))
	setQueryStringList(query, "eventIds", request.EventIDs)
	req.URL.RawQuery = query.Encode()

	response := &ListHeroMetaSharesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

func requirePlayEnvironment(value PlayEnvironment) error {
	v := strings.TrimSpace(string(value))
	if v == "" || v == string(PlayEnvironmentUnspecified) {
		return errors.New("playEnvironment must be specified")
	}
	return nil
}

func setQueryEnum(values url.Values, key, value, unspecified string) {
	if v := strings.TrimSpace(value); v != "" && v != unspecified {
		values.Set(key, v)
	}
}

func setQueryStringList(values url.Values, key string, items []string) {
	for _, item := range items {
		if v := strings.TrimSpace(item); v != "" {
			values.Add(key, v)
		}
	}
}
