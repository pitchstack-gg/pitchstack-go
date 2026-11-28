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

// PaperResultsScopeKind represents meta.v1.PaperResultsScopeKind.
type PaperResultsScopeKind string

const (
	PaperResultsScopeKindUnspecified         PaperResultsScopeKind = "PAPER_RESULTS_SCOPE_KIND_UNSPECIFIED"
	PaperResultsScopeKindEvent               PaperResultsScopeKind = "PAPER_RESULTS_SCOPE_KIND_EVENT"
	PaperResultsScopeKindEventFinalStandings PaperResultsScopeKind = "PAPER_RESULTS_SCOPE_KIND_EVENT_FINAL_STANDINGS"
	PaperResultsScopeKindStage               PaperResultsScopeKind = "PAPER_RESULTS_SCOPE_KIND_STAGE"
	PaperResultsScopeKindRound               PaperResultsScopeKind = "PAPER_RESULTS_SCOPE_KIND_ROUND"
	PaperResultsScopeKindMatch               PaperResultsScopeKind = "PAPER_RESULTS_SCOPE_KIND_MATCH"
	PaperResultsScopeKindGame                PaperResultsScopeKind = "PAPER_RESULTS_SCOPE_KIND_GAME"
	PaperResultsScopeKindTeamSeatContest     PaperResultsScopeKind = "PAPER_RESULTS_SCOPE_KIND_TEAM_SEAT_CONTEST"
)

// PaperAttributionStatus represents meta.v1.PaperAttributionStatus.
type PaperAttributionStatus string

const (
	PaperAttributionStatusUnspecified PaperAttributionStatus = "PAPER_ATTRIBUTION_STATUS_UNSPECIFIED"
	PaperAttributionStatusVerified    PaperAttributionStatus = "PAPER_ATTRIBUTION_STATUS_VERIFIED"
	PaperAttributionStatusMissing     PaperAttributionStatus = "PAPER_ATTRIBUTION_STATUS_MISSING"
	PaperAttributionStatusAmbiguous   PaperAttributionStatus = "PAPER_ATTRIBUTION_STATUS_AMBIGUOUS"
	PaperAttributionStatusConflicting PaperAttributionStatus = "PAPER_ATTRIBUTION_STATUS_CONFLICTING"
)

// PaperCompetitorKind represents meta.v1.PaperCompetitorKind.
type PaperCompetitorKind string

const (
	PaperCompetitorKindUnspecified PaperCompetitorKind = "PAPER_COMPETITOR_KIND_UNSPECIFIED"
	PaperCompetitorKindIndividual  PaperCompetitorKind = "PAPER_COMPETITOR_KIND_INDIVIDUAL"
	PaperCompetitorKindTeam        PaperCompetitorKind = "PAPER_COMPETITOR_KIND_TEAM"
)

// PaperOutcomeKind represents meta.v1.PaperOutcomeKind.
type PaperOutcomeKind string

const (
	PaperOutcomeKindUnspecified PaperOutcomeKind = "PAPER_OUTCOME_KIND_UNSPECIFIED"
	PaperOutcomeKindWin         PaperOutcomeKind = "PAPER_OUTCOME_KIND_WIN"
	PaperOutcomeKindLoss        PaperOutcomeKind = "PAPER_OUTCOME_KIND_LOSS"
	PaperOutcomeKindDraw        PaperOutcomeKind = "PAPER_OUTCOME_KIND_DRAW"
	PaperOutcomeKindBye         PaperOutcomeKind = "PAPER_OUTCOME_KIND_BYE"
	PaperOutcomeKindUnfinished  PaperOutcomeKind = "PAPER_OUTCOME_KIND_UNFINISHED"
)

// PaperAnalyticsPopulationUnit represents meta.v1.PaperAnalyticsPopulationUnit.
type PaperAnalyticsPopulationUnit string

const (
	PaperAnalyticsPopulationUnitUnspecified     PaperAnalyticsPopulationUnit = "PAPER_ANALYTICS_POPULATION_UNIT_UNSPECIFIED"
	PaperAnalyticsPopulationUnitIndividualMatch PaperAnalyticsPopulationUnit = "PAPER_ANALYTICS_POPULATION_UNIT_INDIVIDUAL_MATCH"
	PaperAnalyticsPopulationUnitGame            PaperAnalyticsPopulationUnit = "PAPER_ANALYTICS_POPULATION_UNIT_GAME"
	PaperAnalyticsPopulationUnitTeamMatch       PaperAnalyticsPopulationUnit = "PAPER_ANALYTICS_POPULATION_UNIT_TEAM_MATCH"
	PaperAnalyticsPopulationUnitTeamSeatContest PaperAnalyticsPopulationUnit = "PAPER_ANALYTICS_POPULATION_UNIT_TEAM_SEAT_CONTEST"
	PaperAnalyticsPopulationUnitProgression     PaperAnalyticsPopulationUnit = "PAPER_ANALYTICS_POPULATION_UNIT_PROGRESSION"
)

// PaperPublishedCoverage represents meta.v1.PaperPublishedCoverage.
type PaperPublishedCoverage string

const (
	PaperPublishedCoverageUnspecified                 PaperPublishedCoverage = "PAPER_PUBLISHED_COVERAGE_UNSPECIFIED"
	PaperPublishedCoverageNone                        PaperPublishedCoverage = "PAPER_PUBLISHED_COVERAGE_NONE"
	PaperPublishedCoveragePartialSourceBackedSubset   PaperPublishedCoverage = "PAPER_PUBLISHED_COVERAGE_PARTIAL_SOURCE_BACKED_SUBSET"
	PaperPublishedCoverageVerifiedObservedSourceScope PaperPublishedCoverage = "PAPER_PUBLISHED_COVERAGE_VERIFIED_OBSERVED_SOURCE_SCOPE"
)

// PaperResultsCapability represents meta.v1.PaperResultsCapability.
type PaperResultsCapability string

const (
	PaperResultsCapabilityUnspecified      PaperResultsCapability = "PAPER_RESULTS_CAPABILITY_UNSPECIFIED"
	PaperResultsCapabilityStandings        PaperResultsCapability = "PAPER_RESULTS_CAPABILITY_STANDINGS"
	PaperResultsCapabilityRoundsAndMatches PaperResultsCapability = "PAPER_RESULTS_CAPABILITY_ROUNDS_AND_MATCHES"
	PaperResultsCapabilityDecks            PaperResultsCapability = "PAPER_RESULTS_CAPABILITY_DECKS"
	PaperResultsCapabilityProgression      PaperResultsCapability = "PAPER_RESULTS_CAPABILITY_PROGRESSION"
	PaperResultsCapabilityAnalytics        PaperResultsCapability = "PAPER_RESULTS_CAPABILITY_ANALYTICS"
)

// PaperOfficialEvidenceState represents meta.v1.PaperOfficialEvidenceState.
type PaperOfficialEvidenceState string

const (
	PaperOfficialEvidenceStateUnspecified         PaperOfficialEvidenceState = "PAPER_OFFICIAL_EVIDENCE_STATE_UNSPECIFIED"
	PaperOfficialEvidenceStateVerifiedPublished   PaperOfficialEvidenceState = "PAPER_OFFICIAL_EVIDENCE_STATE_VERIFIED_PUBLISHED"
	PaperOfficialEvidenceStateVerifiedUnpublished PaperOfficialEvidenceState = "PAPER_OFFICIAL_EVIDENCE_STATE_VERIFIED_UNPUBLISHED"
	PaperOfficialEvidenceStateInaccessible        PaperOfficialEvidenceState = "PAPER_OFFICIAL_EVIDENCE_STATE_INACCESSIBLE"
	PaperOfficialEvidenceStateUnverified          PaperOfficialEvidenceState = "PAPER_OFFICIAL_EVIDENCE_STATE_UNVERIFIED"
	PaperOfficialEvidenceStateUnknown             PaperOfficialEvidenceState = "PAPER_OFFICIAL_EVIDENCE_STATE_UNKNOWN"
)

// PaperImplementationSupport represents meta.v1.PaperImplementationSupport.
type PaperImplementationSupport string

const (
	PaperImplementationSupportUnspecified PaperImplementationSupport = "PAPER_IMPLEMENTATION_SUPPORT_UNSPECIFIED"
	PaperImplementationSupportSupported   PaperImplementationSupport = "PAPER_IMPLEMENTATION_SUPPORT_SUPPORTED"
	PaperImplementationSupportUnsupported PaperImplementationSupport = "PAPER_IMPLEMENTATION_SUPPORT_UNSUPPORTED"
	PaperImplementationSupportExcluded    PaperImplementationSupport = "PAPER_IMPLEMENTATION_SUPPORT_EXCLUDED"
	PaperImplementationSupportUnknown     PaperImplementationSupport = "PAPER_IMPLEMENTATION_SUPPORT_UNKNOWN"
)

// PaperProcessingState represents meta.v1.PaperProcessingState.
type PaperProcessingState string

const (
	PaperProcessingStateUnspecified PaperProcessingState = "PAPER_PROCESSING_STATE_UNSPECIFIED"
	PaperProcessingStateNotStarted  PaperProcessingState = "PAPER_PROCESSING_STATE_NOT_STARTED"
	PaperProcessingStateQueued      PaperProcessingState = "PAPER_PROCESSING_STATE_QUEUED"
	PaperProcessingStateRunning     PaperProcessingState = "PAPER_PROCESSING_STATE_RUNNING"
	PaperProcessingStateFailed      PaperProcessingState = "PAPER_PROCESSING_STATE_FAILED"
	PaperProcessingStateComplete    PaperProcessingState = "PAPER_PROCESSING_STATE_COMPLETE"
	PaperProcessingStateUnknown     PaperProcessingState = "PAPER_PROCESSING_STATE_UNKNOWN"
)

// PaperStandingExtentKind represents meta.v1.PaperStandingExtentKind.
type PaperStandingExtentKind string

const (
	PaperStandingExtentKindUnspecified  PaperStandingExtentKind = "PAPER_STANDING_EXTENT_KIND_UNSPECIFIED"
	PaperStandingExtentKindTopN         PaperStandingExtentKind = "PAPER_STANDING_EXTENT_KIND_TOP_N"
	PaperStandingExtentKindFullField    PaperStandingExtentKind = "PAPER_STANDING_EXTENT_KIND_FULL_FIELD"
	PaperStandingExtentKindSourceStated PaperStandingExtentKind = "PAPER_STANDING_EXTENT_KIND_SOURCE_STATED"
	PaperStandingExtentKindUnknown      PaperStandingExtentKind = "PAPER_STANDING_EXTENT_KIND_UNKNOWN"
)

// PaperSourceRevision represents meta.v1.PaperSourceRevision.
type PaperSourceRevision struct {
	SourceID                       string     `json:"sourceId,omitempty"`
	SourceRevision                 string     `json:"sourceRevision,omitempty"`
	CapturedAt                     *time.Time `json:"capturedAt,omitempty"`
	SourceUpdatedAt                *time.Time `json:"sourceUpdatedAt,omitempty"`
	VerifiedAt                     *time.Time `json:"verifiedAt,omitempty"`
	SourceURL                      string     `json:"sourceUrl,omitempty"`
	ContentSHA256                  string     `json:"contentSha256,omitempty"`
	ExpectedPriorConfirmedRevision string     `json:"expectedPriorConfirmedRevision,omitempty"`
}

// PaperPublicResultsScope represents meta.v1.PaperPublicResultsScope.
type PaperPublicResultsScope struct {
	ScopeID       string                `json:"scopeId,omitempty"`
	EventID       string                `json:"eventId,omitempty"`
	Kind          PaperResultsScopeKind `json:"kind,omitempty"`
	ParentScopeID string                `json:"parentScopeId,omitempty"`
	Ordinal       *uint32               `json:"ordinal,omitempty"`
	DisplayLabel  string                `json:"displayLabel,omitempty"`
}

// PaperPublicEventIndividual represents meta.v1.PaperPublicEventIndividual.
type PaperPublicEventIndividual struct {
	IndividualID string               `json:"individualId,omitempty"`
	EventID      string               `json:"eventId,omitempty"`
	DisplayName  string               `json:"displayName,omitempty"`
	Revision     *PaperSourceRevision `json:"revision,omitempty"`
}

// PaperPublicEventTeam represents meta.v1.PaperPublicEventTeam.
type PaperPublicEventTeam struct {
	TeamID      string               `json:"teamId,omitempty"`
	EventID     string               `json:"eventId,omitempty"`
	DisplayName string               `json:"displayName,omitempty"`
	Revision    *PaperSourceRevision `json:"revision,omitempty"`
}

// PaperPublicTeamSeat represents meta.v1.PaperPublicTeamSeat.
type PaperPublicTeamSeat struct {
	SeatID   string               `json:"seatId,omitempty"`
	TeamID   string               `json:"teamId,omitempty"`
	ScopeID  string               `json:"scopeId,omitempty"`
	Ordinal  *uint32              `json:"ordinal,omitempty"`
	Revision *PaperSourceRevision `json:"revision,omitempty"`
}

// PaperTeamMembership represents meta.v1.PaperTeamMembership.
type PaperTeamMembership struct {
	MembershipID      string                 `json:"membershipId,omitempty"`
	TeamID            string                 `json:"teamId,omitempty"`
	IndividualID      string                 `json:"individualId,omitempty"`
	ScopeID           string                 `json:"scopeId,omitempty"`
	SeatID            string                 `json:"seatId,omitempty"`
	AttributionStatus PaperAttributionStatus `json:"attributionStatus,omitempty"`
	Revision          *PaperSourceRevision   `json:"revision,omitempty"`
}

// PaperPublicStage represents meta.v1.PaperPublicStage.
type PaperPublicStage struct {
	StageID      string               `json:"stageId,omitempty"`
	ScopeID      string               `json:"scopeId,omitempty"`
	DisplayLabel string               `json:"displayLabel,omitempty"`
	Ordinal      *uint32              `json:"ordinal,omitempty"`
	Revision     *PaperSourceRevision `json:"revision,omitempty"`
}

// PaperPublicRound represents meta.v1.PaperPublicRound.
type PaperPublicRound struct {
	RoundID      string               `json:"roundId,omitempty"`
	StageID      string               `json:"stageId,omitempty"`
	ScopeID      string               `json:"scopeId,omitempty"`
	DisplayLabel string               `json:"displayLabel,omitempty"`
	Ordinal      *uint32              `json:"ordinal,omitempty"`
	Revision     *PaperSourceRevision `json:"revision,omitempty"`
}

// PaperPublicMatchCompetitor represents meta.v1.PaperPublicMatchCompetitor.
type PaperPublicMatchCompetitor struct {
	CompetitorID      string                 `json:"competitorId,omitempty"`
	Kind              PaperCompetitorKind    `json:"kind,omitempty"`
	AttributionStatus PaperAttributionStatus `json:"attributionStatus,omitempty"`
}

// PaperPublicMatch represents meta.v1.PaperPublicMatch.
type PaperPublicMatch struct {
	MatchID            string                       `json:"matchId,omitempty"`
	RoundID            string                       `json:"roundId,omitempty"`
	ScopeID            string                       `json:"scopeId,omitempty"`
	CompetitorKind     PaperCompetitorKind          `json:"competitorKind,omitempty"`
	Competitors        []PaperPublicMatchCompetitor `json:"competitors,omitempty"`
	WinnerCompetitorID string                       `json:"winnerCompetitorId,omitempty"`
	Outcome            PaperOutcomeKind             `json:"outcome,omitempty"`
	Revision           *PaperSourceRevision         `json:"revision,omitempty"`
}

// PaperPublicGame represents meta.v1.PaperPublicGame.
type PaperPublicGame struct {
	GameID             string               `json:"gameId,omitempty"`
	MatchID            string               `json:"matchId,omitempty"`
	ScopeID            string               `json:"scopeId,omitempty"`
	Ordinal            *uint32              `json:"ordinal,omitempty"`
	WinnerCompetitorID string               `json:"winnerCompetitorId,omitempty"`
	Outcome            PaperOutcomeKind     `json:"outcome,omitempty"`
	Revision           *PaperSourceRevision `json:"revision,omitempty"`
}

// PaperSeatContestCompetitor represents meta.v1.PaperSeatContestCompetitor.
type PaperSeatContestCompetitor struct {
	SeatID            string                 `json:"seatId,omitempty"`
	IndividualID      string                 `json:"individualId,omitempty"`
	AttributionStatus PaperAttributionStatus `json:"attributionStatus,omitempty"`
}

// PaperPublicSeatContest represents meta.v1.PaperPublicSeatContest.
type PaperPublicSeatContest struct {
	SeatContestID      string                       `json:"seatContestId,omitempty"`
	MatchID            string                       `json:"matchId,omitempty"`
	ScopeID            string                       `json:"scopeId,omitempty"`
	Competitors        []PaperSeatContestCompetitor `json:"competitors,omitempty"`
	WinnerIndividualID string                       `json:"winnerIndividualId,omitempty"`
	Outcome            PaperOutcomeKind             `json:"outcome,omitempty"`
	Revision           *PaperSourceRevision         `json:"revision,omitempty"`
}

// PaperFormatObservation represents meta.v1.PaperFormatObservation.
type PaperFormatObservation struct {
	ObservationID     string                 `json:"observationId,omitempty"`
	ScopeID           string                 `json:"scopeId,omitempty"`
	Format            string                 `json:"format,omitempty"`
	AttributionStatus PaperAttributionStatus `json:"attributionStatus,omitempty"`
	Revision          *PaperSourceRevision   `json:"revision,omitempty"`
}

// PaperHeroObservation represents meta.v1.PaperHeroObservation.
type PaperHeroObservation struct {
	ObservationID     string                 `json:"observationId,omitempty"`
	ScopeID           string                 `json:"scopeId,omitempty"`
	IndividualID      string                 `json:"individualId,omitempty"`
	SourceHeroName    string                 `json:"sourceHeroName,omitempty"`
	HeroID            string                 `json:"heroId,omitempty"`
	AttributionStatus PaperAttributionStatus `json:"attributionStatus,omitempty"`
	Revision          *PaperSourceRevision   `json:"revision,omitempty"`
}

// PaperDeckCard represents meta.v1.DeckCard, one card line of a published decklist.
type PaperDeckCard struct {
	SourceCardID string `json:"sourceCardId,omitempty"`
	CardID       string `json:"cardId,omitempty"`
	Copies       int32  `json:"copies,omitempty"`
	Zone         string `json:"zone,omitempty"`
}

// PaperPublishedDecklist represents meta.v1.PaperPublishedDecklist.
type PaperPublishedDecklist struct {
	DeckID    string               `json:"deckId,omitempty"`
	EventID   string               `json:"eventId,omitempty"`
	SourceURL string               `json:"sourceUrl,omitempty"`
	Cards     []PaperDeckCard      `json:"cards,omitempty"`
	Revision  *PaperSourceRevision `json:"revision,omitempty"`
}

// PaperDeckOwnershipObservation represents meta.v1.PaperDeckOwnershipObservation.
type PaperDeckOwnershipObservation struct {
	ObservationID     string                 `json:"observationId,omitempty"`
	DeckID            string                 `json:"deckId,omitempty"`
	ScopeID           string                 `json:"scopeId,omitempty"`
	IndividualID      string                 `json:"individualId,omitempty"`
	AttributionStatus PaperAttributionStatus `json:"attributionStatus,omitempty"`
	Revision          *PaperSourceRevision   `json:"revision,omitempty"`
}

// PaperDeckUseObservation represents meta.v1.PaperDeckUseObservation.
type PaperDeckUseObservation struct {
	ObservationID               string                 `json:"observationId,omitempty"`
	DeckID                      string                 `json:"deckId,omitempty"`
	ScopeID                     string                 `json:"scopeId,omitempty"`
	IndividualID                string                 `json:"individualId,omitempty"`
	DeckAttributionStatus       PaperAttributionStatus `json:"deckAttributionStatus,omitempty"`
	IndividualAttributionStatus PaperAttributionStatus `json:"individualAttributionStatus,omitempty"`
	Revision                    *PaperSourceRevision   `json:"revision,omitempty"`
}

// PaperStanding represents meta.v1.PaperStanding.
type PaperStanding struct {
	StandingID             string                 `json:"standingId,omitempty"`
	ScopeID                string                 `json:"scopeId,omitempty"`
	CompetitorKind         PaperCompetitorKind    `json:"competitorKind,omitempty"`
	ReportedRank           uint32                 `json:"reportedRank,omitempty"`
	CompetitorID           string                 `json:"competitorId,omitempty"`
	ReportedCompetitorName string                 `json:"reportedCompetitorName,omitempty"`
	Tied                   *bool                  `json:"tied,omitempty"`
	AttributionStatus      PaperAttributionStatus `json:"attributionStatus,omitempty"`
	Revision               *PaperSourceRevision   `json:"revision,omitempty"`
}

// PaperProgressionFact represents meta.v1.PaperProgressionFact.
type PaperProgressionFact struct {
	ProgressionID      string               `json:"progressionId,omitempty"`
	ScopeID            string               `json:"scopeId,omitempty"`
	Key                string               `json:"key,omitempty"`
	DisplayLabel       string               `json:"displayLabel,omitempty"`
	ObservedPopulation *uint32              `json:"observedPopulation,omitempty"`
	ObservedProgressed *uint32              `json:"observedProgressed,omitempty"`
	Revision           *PaperSourceRevision `json:"revision,omitempty"`
}

// PaperAnalyticsDimensions represents meta.v1.PaperAnalyticsDimensions.
type PaperAnalyticsDimensions struct {
	Format         string `json:"format,omitempty"`
	SubjectHeroID  string `json:"subjectHeroId,omitempty"`
	OpposingHeroID string `json:"opposingHeroId,omitempty"`
	TeamID         string `json:"teamId,omitempty"`
	OpposingTeamID string `json:"opposingTeamId,omitempty"`
	SeatID         string `json:"seatId,omitempty"`
	OpposingSeatID string `json:"opposingSeatId,omitempty"`
	ProgressionID  string `json:"progressionId,omitempty"`
}

// PaperAnalyticsCoverage represents meta.v1.PaperAnalyticsCoverage.
type PaperAnalyticsCoverage struct {
	ObservedObservations uint64                 `json:"observedObservations,omitempty,string"`
	EligibleObservations uint64                 `json:"eligibleObservations,omitempty,string"`
	ExcludedObservations uint64                 `json:"excludedObservations,omitempty,string"`
	PublishedCoverage    PaperPublishedCoverage `json:"publishedCoverage,omitempty"`
}

// PaperAnalyticsPopulation represents meta.v1.PaperAnalyticsPopulation.
type PaperAnalyticsPopulation struct {
	PopulationID                   string                       `json:"populationId,omitempty"`
	ScopeID                        string                       `json:"scopeId,omitempty"`
	Unit                           PaperAnalyticsPopulationUnit `json:"unit,omitempty"`
	Dimensions                     *PaperAnalyticsDimensions    `json:"dimensions,omitempty"`
	Coverage                       *PaperAnalyticsCoverage      `json:"coverage,omitempty"`
	CompletedObservations          uint64                       `json:"completedObservations,omitempty,string"`
	Wins                           uint64                       `json:"wins,omitempty,string"`
	Losses                         uint64                       `json:"losses,omitempty,string"`
	Draws                          uint64                       `json:"draws,omitempty,string"`
	CompleteAttributableComparable bool                         `json:"completeAttributableComparable,omitempty"`
	WinRate                        *float64                     `json:"winRate,omitempty"`
	Revision                       *PaperSourceRevision         `json:"revision,omitempty"`
}

// PaperStandingExtent represents meta.v1.PaperStandingExtent.
type PaperStandingExtent struct {
	Kind        PaperStandingExtentKind `json:"kind,omitempty"`
	RankLimit   *uint32                 `json:"rankLimit,omitempty"`
	SourceLabel string                  `json:"sourceLabel,omitempty"`
}

// PaperStandingCounts represents meta.v1.PaperStandingCounts.
type PaperStandingCounts struct {
	ObservedRows       uint32  `json:"observedRows,omitempty"`
	PublishedRows      uint32  `json:"publishedRows,omitempty"`
	RejectedRows       uint32  `json:"rejectedRows,omitempty"`
	UnresolvedRows     uint32  `json:"unresolvedRows,omitempty"`
	EvidencedFieldSize *uint32 `json:"evidencedFieldSize,omitempty"`
}

// PaperResultsAvailability represents meta.v1.PaperResultsAvailability.
type PaperResultsAvailability struct {
	Capability            PaperResultsCapability     `json:"capability,omitempty"`
	ScopeID               string                     `json:"scopeId,omitempty"`
	OfficialEvidence      PaperOfficialEvidenceState `json:"officialEvidence,omitempty"`
	ImplementationSupport PaperImplementationSupport `json:"implementationSupport,omitempty"`
	Processing            PaperProcessingState       `json:"processing,omitempty"`
	PublishedCoverage     PaperPublishedCoverage     `json:"publishedCoverage,omitempty"`
	Stale                 *bool                      `json:"stale,omitempty"`
	StandingExtent        *PaperStandingExtent       `json:"standingExtent,omitempty"`
	Counts                *PaperStandingCounts       `json:"counts,omitempty"`
	Revision              *PaperSourceRevision       `json:"revision,omitempty"`
}

// PaperEventResults represents meta.v1.PaperEventResults.
type PaperEventResults struct {
	SchemaVersion        uint32                          `json:"schemaVersion,omitempty"`
	EventID              string                          `json:"eventId,omitempty"`
	Revision             *PaperSourceRevision            `json:"revision,omitempty"`
	Scopes               []PaperPublicResultsScope       `json:"scopes,omitempty"`
	Individuals          []PaperPublicEventIndividual    `json:"individuals,omitempty"`
	Teams                []PaperPublicEventTeam          `json:"teams,omitempty"`
	Seats                []PaperPublicTeamSeat           `json:"seats,omitempty"`
	TeamMemberships      []PaperTeamMembership           `json:"teamMemberships,omitempty"`
	Stages               []PaperPublicStage              `json:"stages,omitempty"`
	Rounds               []PaperPublicRound              `json:"rounds,omitempty"`
	Matches              []PaperPublicMatch              `json:"matches,omitempty"`
	Games                []PaperPublicGame               `json:"games,omitempty"`
	SeatContests         []PaperPublicSeatContest        `json:"seatContests,omitempty"`
	Formats              []PaperFormatObservation        `json:"formats,omitempty"`
	Heroes               []PaperHeroObservation          `json:"heroes,omitempty"`
	Decks                []PaperPublishedDecklist        `json:"decks,omitempty"`
	DeckOwnership        []PaperDeckOwnershipObservation `json:"deckOwnership,omitempty"`
	DeckUses             []PaperDeckUseObservation       `json:"deckUses,omitempty"`
	Standings            []PaperStanding                 `json:"standings,omitempty"`
	Progression          []PaperProgressionFact          `json:"progression,omitempty"`
	AnalyticsPopulations []PaperAnalyticsPopulation      `json:"analyticsPopulations,omitempty"`
	Availability         []PaperResultsAvailability      `json:"availability,omitempty"`
}

// GetPaperEventResultsRequest identifies a paper event's published results.
// An empty SourceID selects the default source; empty ScopeIDs returns every published scope.
type GetPaperEventResultsRequest struct {
	EventID  string `json:"-"`
	SourceID string
	ScopeIDs []string
}

// GetPaperEventResultsResponse returns a paper event's published results.
type GetPaperEventResultsResponse struct {
	Results  *PaperEventResults `json:"results,omitempty"`
	Metadata ResponseMetadata   `json:"-"`
}

func (r *GetPaperEventResultsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetPaperEventResults fetches the published results of a paper event.
func (c *Client) GetPaperEventResults(ctx context.Context, request *GetPaperEventResultsRequest, opts ...RequestOpt) (*GetPaperEventResultsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	eventID := strings.TrimSpace(request.EventID)
	if eventID == "" {
		return nil, errors.New("eventID must not be empty")
	}

	path := fmt.Sprintf("/v1/meta/events/%s/paper-results", url.PathEscape(eventID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "sourceId", request.SourceID)
	setQueryStringList(query, "scopeIds", request.ScopeIDs)
	req.URL.RawQuery = query.Encode()

	response := &GetPaperEventResultsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}
