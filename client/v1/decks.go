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

// DeckListScope matches v1DeckListScope.
type DeckListScope string

const (
	DeckListScopeUnspecified DeckListScope = "DECK_LIST_SCOPE_UNSPECIFIED"
	DeckListScopeOwned       DeckListScope = "DECK_LIST_SCOPE_OWNED"
	DeckListScopeShared      DeckListScope = "DECK_LIST_SCOPE_SHARED"
	DeckListScopeAccessible  DeckListScope = "DECK_LIST_SCOPE_ACCESSIBLE"
)

// DeckKind matches v1DeckKind.
type DeckKind string

const (
	DeckKindUnspecified DeckKind = "DECK_KIND_UNSPECIFIED"
	DeckKindUser        DeckKind = "DECK_KIND_USER"
	DeckKindReference   DeckKind = "DECK_KIND_REFERENCE"
)

// DeckSourceKind matches v1DeckSourceKind.
type DeckSourceKind string

const (
	DeckSourceKindUnspecified      DeckSourceKind = "DECK_SOURCE_KIND_UNSPECIFIED"
	DeckSourceKindPromoArticle     DeckSourceKind = "DECK_SOURCE_KIND_PROMO_ARTICLE"
	DeckSourceKindPrecon           DeckSourceKind = "DECK_SOURCE_KIND_PRECON"
	DeckSourceKindTournamentResult DeckSourceKind = "DECK_SOURCE_KIND_TOURNAMENT_RESULT"
)

// BoardType represents v1BoardType.
type BoardType string

const (
	BoardTypeUnspecified BoardType = "BOARD_TYPE_UNSPECIFIED"
	BoardTypeMainboard   BoardType = "BOARD_TYPE_MAINBOARD"
	BoardTypeSideboard   BoardType = "BOARD_TYPE_SIDEBOARD"
	BoardTypeMaybeboard  BoardType = "BOARD_TYPE_MAYBEBOARD"
)

// SideboardGuideTargetType represents v1SideboardGuideTargetType.
type SideboardGuideTargetType string

const (
	SideboardGuideTargetTypeUnspecified SideboardGuideTargetType = "SIDEBOARD_GUIDE_TARGET_TYPE_UNSPECIFIED"
	SideboardGuideTargetTypeHero        SideboardGuideTargetType = "SIDEBOARD_GUIDE_TARGET_TYPE_HERO"
	SideboardGuideTargetTypeClass       SideboardGuideTargetType = "SIDEBOARD_GUIDE_TARGET_TYPE_CLASS"
	SideboardGuideTargetTypeArchetype   SideboardGuideTargetType = "SIDEBOARD_GUIDE_TARGET_TYPE_ARCHETYPE"
)

// DeckPermission mirrors authzv1Permission for decks.
type DeckPermission string

const (
	DeckPermissionUnspecified DeckPermission = "PERMISSION_UNSPECIFIED"
	DeckPermissionReader      DeckPermission = "PERMISSION_READER"
	DeckPermissionWriter      DeckPermission = "PERMISSION_WRITER"
)

// Deck models v1Deck.
type Deck struct {
	ID                  string               `json:"id,omitempty"`
	UserID              string               `json:"userId,omitempty"`
	Name                string               `json:"name,omitempty"`
	Author              string               `json:"author,omitempty"`
	HeroID              string               `json:"heroId,omitempty"`
	Format              string               `json:"format,omitempty"`
	Visibility          VisibilityLevel      `json:"visibility,omitempty"`
	DeckVersions        []DeckVersionSummary `json:"deckVersions,omitempty"`
	ActiveDeckVersionID string               `json:"activeDeckVersionId,omitempty"`
	DeckKind            DeckKind             `json:"deckKind,omitempty"`
	SourceKind          DeckSourceKind       `json:"sourceKind,omitempty"`
	SourceReference     string               `json:"sourceReference,omitempty"`
	CreatedAt           *time.Time           `json:"createdAt,omitempty"`
	UpdatedAt           *time.Time           `json:"updatedAt,omitempty"`
}

// DeckVersionSummary represents v1DeckVersionSummary.
type DeckVersionSummary struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

// DeckVersion represents v1DeckVersion.
type DeckVersion struct {
	DeckID    string     `json:"deckId,omitempty"`
	Name      string     `json:"name,omitempty"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	ID        string     `json:"id,omitempty"`
}

// DeckVersionChange mirrors v1DeckVersionChange.
type DeckVersionChange struct {
	DeckID      string                     `json:"deckId,omitempty"`
	Name        string                     `json:"name,omitempty"`
	Description string                     `json:"description,omitempty"`
	Timestamp   *time.Time                 `json:"timestamp,omitempty"`
	ID          string                     `json:"id,omitempty"`
	EventType   DeckVersionChangeEventType `json:"eventType,omitempty"`
	CardChanges []DeckVersionCardChange    `json:"cardChanges,omitempty"`
}

// DeckVersionChangeEventType matches v1DeckVersionChangeEventType.
type DeckVersionChangeEventType string

const (
	DeckVersionChangeEventTypeUnspecified          DeckVersionChangeEventType = "DECK_VERSION_CHANGE_EVENT_TYPE_UNSPECIFIED"
	DeckVersionChangeEventTypeDeckCardChange       DeckVersionChangeEventType = "DECK_VERSION_CHANGE_EVENT_TYPE_DECK_CARD_CHANGE"
	DeckVersionChangeEventTypeVersionCreated       DeckVersionChangeEventType = "DECK_VERSION_CHANGE_EVENT_TYPE_VERSION_CREATED"
	DeckVersionChangeEventTypeVersionCloned        DeckVersionChangeEventType = "DECK_VERSION_CHANGE_EVENT_TYPE_VERSION_CLONED"
	DeckVersionChangeEventTypeNotesUpdated         DeckVersionChangeEventType = "DECK_VERSION_CHANGE_EVENT_TYPE_NOTES_UPDATED"
	DeckVersionChangeEventTypeSideboardGuideUpsert DeckVersionChangeEventType = "DECK_VERSION_CHANGE_EVENT_TYPE_SIDEBOARD_GUIDE_UPSERT"
	DeckVersionChangeEventTypeSideboardGuideDelete DeckVersionChangeEventType = "DECK_VERSION_CHANGE_EVENT_TYPE_SIDEBOARD_GUIDE_DELETE"
	DeckVersionChangeEventTypeMatchCreated         DeckVersionChangeEventType = "DECK_VERSION_CHANGE_EVENT_TYPE_MATCH_CREATED"
	DeckVersionChangeEventTypeMatchDeleted         DeckVersionChangeEventType = "DECK_VERSION_CHANGE_EVENT_TYPE_MATCH_DELETED"
)

// DeckVersionCardChangeOperation matches v1DeckVersionCardChangeOperation.
type DeckVersionCardChangeOperation string

const (
	DeckVersionCardChangeOperationUnspecified    DeckVersionCardChangeOperation = "DECK_VERSION_CARD_CHANGE_OPERATION_UNSPECIFIED"
	DeckVersionCardChangeOperationAdd            DeckVersionCardChangeOperation = "DECK_VERSION_CARD_CHANGE_OPERATION_ADD"
	DeckVersionCardChangeOperationRemove         DeckVersionCardChangeOperation = "DECK_VERSION_CARD_CHANGE_OPERATION_REMOVE"
	DeckVersionCardChangeOperationQuantityChange DeckVersionCardChangeOperation = "DECK_VERSION_CARD_CHANGE_OPERATION_QUANTITY_CHANGE"
	DeckVersionCardChangeOperationMove           DeckVersionCardChangeOperation = "DECK_VERSION_CARD_CHANGE_OPERATION_MOVE"
)

// DeckVersionCardChange describes a card change in deck version history.
type DeckVersionCardChange struct {
	CardID           string                         `json:"cardId,omitempty"`
	FromBoard        BoardType                      `json:"fromBoard,omitempty"`
	ToBoard          BoardType                      `json:"toBoard,omitempty"`
	QuantityDelta    int32                          `json:"quantityDelta,omitempty"`
	PreviousQuantity int32                          `json:"previousQuantity,omitempty"`
	NewQuantity      int32                          `json:"newQuantity,omitempty"`
	Operation        DeckVersionCardChangeOperation `json:"operation,omitempty"`
}

// DeckCard mirrors v1Card used within deck operations.
type DeckCard struct {
	CardID   string `json:"cardId,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
}

// ListDecksRequest configures filters for ListDecks.
type ListDecksRequest struct {
	Scope     DeckListScope
	UserID    string
	HeroID    string
	Format    string
	Name      string
	PageSize  *int32
	NextToken string
	SubjectID string
}

// ListDecksResponse contains decks visible to the caller.
type ListDecksResponse struct {
	Decks     []Deck           `json:"decks,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListDecksResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateDeckRequest captures payload fields for deck creation.
type CreateDeckRequest struct {
	Name                 string                    `json:"name,omitempty"`
	HeroID               string                    `json:"heroId,omitempty"`
	Format               string                    `json:"format,omitempty"`
	Author               string                    `json:"author,omitempty"`
	Visibility           VisibilityLevel           `json:"visibility,omitempty"`
	DeckID               string                    `json:"deckId,omitempty"`
	CreateInitialVersion *bool                     `json:"createInitialVersion,omitempty"`
	InitialVersion       *CreateDeckInitialVersion `json:"initialVersion,omitempty"`
}

// CreateDeckInitialVersion configures the initial deck version.
type CreateDeckInitialVersion struct {
	Name          string `json:"name,omitempty"`
	DeckVersionID string `json:"deckVersionId,omitempty"`
}

// CreateDeckResponse returns the created deck.
type CreateDeckResponse struct {
	Deck           *Deck            `json:"deck,omitempty"`
	InitialVersion *DeckVersion     `json:"initialVersion,omitempty"`
	Metadata       ResponseMetadata `json:"-"`
}

func (r *CreateDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CloneDeckRequest clones a deck from an existing deck version.
type CloneDeckRequest struct {
	SourceDeckVersionID  string           `json:"sourceDeckVersionId,omitempty"`
	Name                 string           `json:"name,omitempty"`
	Visibility           *VisibilityLevel `json:"visibility,omitempty"`
	DeckID               string           `json:"deckId,omitempty"`
	InitialVersionName   string           `json:"initialVersionName,omitempty"`
	InitialDeckVersionID string           `json:"initialDeckVersionId,omitempty"`
}

// CloneDeckResponse returns the newly created deck and initial version.
type CloneDeckResponse struct {
	Deck           *Deck            `json:"deck,omitempty"`
	InitialVersion *DeckVersion     `json:"initialVersion,omitempty"`
	Metadata       ResponseMetadata `json:"-"`
}

func (r *CloneDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// SearchDecksRequest configures deck search parameters.
type SearchDecksRequest struct {
	HeroID          string
	Format          string
	SearchTerm      string
	PageSize        *int32
	NextToken       string
	DeckKind        DeckKind
	SourceKind      DeckSourceKind
	SourceReference string
}

// SearchDecksResponse returns decks matching the search query.
type SearchDecksResponse struct {
	Decks     []Deck           `json:"decks,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *SearchDecksResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetDeckRequest identifies the deck to retrieve.
type GetDeckRequest struct {
	DeckID string `json:"-"`
}

// GetDeckResponse wraps the retrieved deck.
type GetDeckResponse struct {
	Deck     *Deck            `json:"deck,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteDeckRequest identifies the deck to delete.
type DeleteDeckRequest struct {
	DeckID string `json:"-"`
}

// DeleteDeckResponse captures metadata for delete calls.
type DeleteDeckResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateDeckRequest applies changes to a deck.
type UpdateDeckRequest struct {
	DeckID              string  `json:"-"`
	Name                *string `json:"name,omitempty"`
	Author              *string `json:"author,omitempty"`
	ActiveDeckVersionID *string `json:"activeDeckVersionId,omitempty"`
}

// UpdateDeckResponse contains the updated deck.
type UpdateDeckResponse struct {
	Deck     *Deck            `json:"deck,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateDeckVisibilityRequest updates only the visibility field for a deck.
type UpdateDeckVisibilityRequest struct {
	DeckID     string           `json:"-"`
	Visibility *VisibilityLevel `json:"visibility,omitempty"`
}

// UpdateDeckVisibilityResponse returns the updated deck.
type UpdateDeckVisibilityResponse struct {
	Deck     *Deck            `json:"deck,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateDeckVisibilityResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GrantDeckAccessRequest assigns a permission for a deck.
type GrantDeckAccessRequest struct {
	DeckID     string         `json:"resourceId,omitempty"`
	SubjectID  string         `json:"subjectId,omitempty"`
	Permission DeckPermission `json:"permission,omitempty"`
}

// GrantDeckAccessResponse captures metadata for grant operations.
type GrantDeckAccessResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *GrantDeckAccessResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// RevokeDeckAccessRequest removes a permission for a deck.
type RevokeDeckAccessRequest struct {
	DeckID     string         `json:"resourceId,omitempty"`
	SubjectID  string         `json:"subjectId,omitempty"`
	Permission DeckPermission `json:"permission,omitempty"`
}

// RevokeDeckAccessResponse captures metadata for revoke operations.
type RevokeDeckAccessResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *RevokeDeckAccessResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeckAccessGrant describes an explicit permission grant for a deck.
type DeckAccessGrant struct {
	SubjectID  string         `json:"subjectId,omitempty"`
	Permission DeckPermission `json:"permission,omitempty"`
}

// GetDeckAccessRequest identifies the deck access to retrieve.
type GetDeckAccessRequest struct {
	DeckID string `json:"-"`
}

// GetDeckAccessResponse returns the caller's effective permission for a deck.
type GetDeckAccessResponse struct {
	Permission DeckPermission   `json:"permission,omitempty"`
	Metadata   ResponseMetadata `json:"-"`
}

func (r *GetDeckAccessResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// StopDeckShareRequest removes all explicit shares for a deck.
type StopDeckShareRequest struct {
	DeckID string `json:"-"`
}

// StopDeckShareResponse captures metadata for stop-share operations.
type StopDeckShareResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *StopDeckShareResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListDeckAccessGrantsRequest lists subjects explicitly granted access to a deck.
type ListDeckAccessGrantsRequest struct {
	DeckID    string `json:"-"`
	PageSize  *int32
	NextToken string
}

// ListDeckAccessGrantsResponse returns explicit access grants for a deck.
type ListDeckAccessGrantsResponse struct {
	Grants    []DeckAccessGrant `json:"grants,omitempty"`
	NextToken string            `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata  `json:"-"`
}

func (r *ListDeckAccessGrantsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// StarDeckRequest marks a deck as starred.
type StarDeckRequest struct {
	DeckID string `json:"-"`
}

// StarDeckResponse captures metadata for star operations.
type StarDeckResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *StarDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnstarDeckRequest removes a star from a deck.
type UnstarDeckRequest struct {
	DeckID string `json:"-"`
}

// UnstarDeckResponse captures metadata for unstar operations.
type UnstarDeckResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnstarDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListDeckVersionsRequest enumerates versions for a deck.
type ListDeckVersionsRequest struct {
	DeckID    string
	PageSize  *int32
	NextToken string
}

// ListDeckVersionsResponse lists deck versions.
type ListDeckVersionsResponse struct {
	DeckVersions []DeckVersion    `json:"deckVersions,omitempty"`
	NextToken    string           `json:"nextToken,omitempty"`
	Metadata     ResponseMetadata `json:"-"`
}

func (r *ListDeckVersionsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateDeckVersionRequest provisions a deck version.
type CreateDeckVersionRequest struct {
	DeckID              string `json:"-"`
	Name                string `json:"name,omitempty"`
	DeckVersionID       string `json:"deckVersionId,omitempty"`
	SourceDeckVersionID string `json:"sourceDeckVersionId,omitempty"`
	SetActive           *bool  `json:"setActive,omitempty"`
}

// CreateDeckVersionResponse includes the created version.
type CreateDeckVersionResponse struct {
	DeckVersion *DeckVersion     `json:"deckVersion,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *CreateDeckVersionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetDeckVersionRequest identifies a deck version to retrieve.
type GetDeckVersionRequest struct {
	DeckVersionID string `json:"-"`
}

// GetDeckVersionResponse returns a deck version.
type GetDeckVersionResponse struct {
	DeckVersion *DeckVersion     `json:"deckVersion,omitempty"`
	NextToken   string           `json:"nextToken,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *GetDeckVersionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteDeckVersionRequest removes a deck version.
type DeleteDeckVersionRequest struct {
	DeckVersionID string `json:"-"`
}

// DeleteDeckVersionResponse captures metadata for deletions.
type DeleteDeckVersionResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteDeckVersionResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListDeckVersionCardsRequest lists cards within a deck version.
type ListDeckVersionCardsRequest struct {
	DeckVersionID string
}

// ListDeckVersionCardsResponse enumerates cards within a deck version.
type ListDeckVersionCardsResponse struct {
	MainboardCards  []DeckCard       `json:"mainboardCards,omitempty"`
	SideboardCards  []DeckCard       `json:"sideboardCards,omitempty"`
	MaybeboardCards []DeckCard       `json:"maybeboardCards,omitempty"`
	Metadata        ResponseMetadata `json:"-"`
}

func (r *ListDeckVersionCardsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ModifyDeckVersionCardRequest adjusts a card within a deck version.
type ModifyDeckVersionCardRequest struct {
	DeckVersionID string    `json:"-"`
	CardID        string    `json:"cardId,omitempty"`
	Board         BoardType `json:"boardType,omitempty"`
	Quantity      int32     `json:"quantity,omitempty"`
}

// ModifyDeckVersionCardResponse returns the affected card.
type ModifyDeckVersionCardResponse struct {
	Card     *DeckCard        `json:"card,omitempty"`
	Board    BoardType        `json:"boardType,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ModifyDeckVersionCardResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetDeckVersionHistoryRequest retrieves version history for a deck.
type GetDeckVersionHistoryRequest struct {
	DeckVersionID string
}

// GetDeckVersionHistoryResponse contains version changes.
type GetDeckVersionHistoryResponse struct {
	Changes  []DeckVersionChange `json:"changes,omitempty"`
	Metadata ResponseMetadata    `json:"-"`
}

func (r *GetDeckVersionHistoryResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetDeckVersionNotesRequest fetches notes for a deck version.
type GetDeckVersionNotesRequest struct {
	DeckVersionID string
}

// GetDeckVersionNotesResponse returns notes text.
type GetDeckVersionNotesResponse struct {
	Notes    string           `json:"notes,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetDeckVersionNotesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateDeckVersionNotesRequest updates notes for a deck version.
type UpdateDeckVersionNotesRequest struct {
	DeckVersionID string
	Notes         string
}

// UpdateDeckVersionNotesResponse returns updated notes.
type UpdateDeckVersionNotesResponse struct {
	Notes    string           `json:"notes,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateDeckVersionNotesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// SideboardGuide mirrors v1SideboardGuide.
type SideboardGuide struct {
	TargetType    SideboardGuideTargetType   `json:"targetType,omitempty"`
	Target        string                     `json:"target,omitempty"`
	Guide         string                     `json:"guide,omitempty"`
	CardsToAdd    []SideboardGuideCardChange `json:"cardsToAdd,omitempty"`
	CardsToRemove []SideboardGuideCardChange `json:"cardsToRemove,omitempty"`
}

// SideboardGuideCardChange mirrors v1SideboardGuideCardChange.
type SideboardGuideCardChange struct {
	CardID   string `json:"cardId,omitempty"`
	Quantity int32  `json:"quantity,omitempty"`
}

// ListDeckVersionSideboardGuidesRequest lists sideboard guides for a deck version.
type ListDeckVersionSideboardGuidesRequest struct {
	DeckVersionID string
}

// ListDeckVersionSideboardGuidesResponse returns sideboard guides for a deck version.
type ListDeckVersionSideboardGuidesResponse struct {
	SideboardGuides []SideboardGuide `json:"sideboardGuides,omitempty"`
	Metadata        ResponseMetadata `json:"-"`
}

func (r *ListDeckVersionSideboardGuidesResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpsertDeckVersionSideboardGuideRequest creates or updates a sideboard guide.
type UpsertDeckVersionSideboardGuideRequest struct {
	DeckVersionID string                     `json:"-"`
	TargetType    SideboardGuideTargetType   `json:"targetType,omitempty"`
	Target        string                     `json:"target,omitempty"`
	Guide         string                     `json:"guide,omitempty"`
	CardsToAdd    []SideboardGuideCardChange `json:"cardsToAdd,omitempty"`
	CardsToRemove []SideboardGuideCardChange `json:"cardsToRemove,omitempty"`
}

// UpsertDeckVersionSideboardGuideResponse returns the upserted guide.
type UpsertDeckVersionSideboardGuideResponse struct {
	SideboardGuide *SideboardGuide  `json:"sideboardGuide,omitempty"`
	Metadata       ResponseMetadata `json:"-"`
}

func (r *UpsertDeckVersionSideboardGuideResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// DeleteDeckVersionSideboardGuideRequest removes a sideboard guide.
type DeleteDeckVersionSideboardGuideRequest struct {
	DeckVersionID string                   `json:"-"`
	TargetType    SideboardGuideTargetType `json:"targetType,omitempty"`
	Target        string                   `json:"target,omitempty"`
}

// DeleteDeckVersionSideboardGuideResponse captures metadata for deletions.
type DeleteDeckVersionSideboardGuideResponse struct {
	Metadata ResponseMetadata `json:"-"`
}

func (r *DeleteDeckVersionSideboardGuideResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetDecksRequest fetches multiple decks by ID.
type BatchGetDecksRequest struct {
	DeckIDs      []string `json:"deckIds,omitempty"`
	AllowPartial bool     `json:"allowPartial,omitempty"`
}

// BatchGetDecksResponse returns decks requested in batch.
type BatchGetDecksResponse struct {
	Decks       []Deck           `json:"decks,omitempty"`
	NotFoundIDs []string         `json:"notFoundIds,omitempty"`
	Metadata    ResponseMetadata `json:"-"`
}

func (r *BatchGetDecksResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ExportDeckRequest retrieves a deck export snapshot.
type ExportDeckRequest struct {
	DeckID string `json:"-"`
}

// ExportDeckVersion captures a deck version export.
type ExportDeckVersion struct {
	DeckVersion     *DeckVersion     `json:"deckVersion,omitempty"`
	Notes           string           `json:"notes,omitempty"`
	MainboardCards  []DeckCard       `json:"mainboardCards,omitempty"`
	SideboardCards  []DeckCard       `json:"sideboardCards,omitempty"`
	MaybeboardCards []DeckCard       `json:"maybeboardCards,omitempty"`
	SideboardGuides []SideboardGuide `json:"sideboardGuides,omitempty"`
}

// ExportDeckResponse returns a deck export snapshot.
type ExportDeckResponse struct {
	Deck     *Deck               `json:"deck,omitempty"`
	Versions []ExportDeckVersion `json:"versions,omitempty"`
	Metadata ResponseMetadata    `json:"-"`
}

func (r *ExportDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ImportDeckVersion describes an imported deck version snapshot.
type ImportDeckVersion struct {
	DeckVersionID   string           `json:"deckVersionId,omitempty"`
	Name            string           `json:"name,omitempty"`
	Notes           string           `json:"notes,omitempty"`
	MainboardCards  []DeckCard       `json:"mainboardCards,omitempty"`
	SideboardCards  []DeckCard       `json:"sideboardCards,omitempty"`
	MaybeboardCards []DeckCard       `json:"maybeboardCards,omitempty"`
	CreatedAt       *time.Time       `json:"createdAt,omitempty"`
	SideboardGuides []SideboardGuide `json:"sideboardGuides,omitempty"`
}

// ImportDeckRequest imports a deck snapshot.
type ImportDeckRequest struct {
	DeckID              string              `json:"deckId,omitempty"`
	Name                string              `json:"name,omitempty"`
	HeroID              string              `json:"heroId,omitempty"`
	Format              string              `json:"format,omitempty"`
	Author              string              `json:"author,omitempty"`
	Visibility          *VisibilityLevel    `json:"visibility,omitempty"`
	CreatedAt           *time.Time          `json:"createdAt,omitempty"`
	Versions            []ImportDeckVersion `json:"versions,omitempty"`
	ActiveDeckVersionID string              `json:"activeDeckVersionId,omitempty"`
}

// ImportDeckResponse returns the imported deck details.
type ImportDeckResponse struct {
	Deck             *Deck            `json:"deck,omitempty"`
	ImportedVersions int32            `json:"importedVersions,omitempty"`
	Metadata         ResponseMetadata `json:"-"`
}

func (r *ImportDeckResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListStarredDecksRequest lists decks starred by a user.
type ListStarredDecksRequest struct {
	UserID string
	Limit  *int32
	Offset *int32
}

// ListStarredDecksResponse returns starred decks for a user.
type ListStarredDecksResponse struct {
	Decks    []Deck           `json:"decks,omitempty"`
	Total    int32            `json:"total,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *ListStarredDecksResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListDecks enumerates decks visible to the caller.
func (c *Client) ListDecks(ctx context.Context, request *ListDecksRequest, opts ...RequestOpt) (*ListDecksResponse, error) {
	if request == nil {
		request = &ListDecksRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/decks", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if scope := strings.TrimSpace(string(request.Scope)); scope != "" && scope != string(DeckListScopeUnspecified) {
		query.Set("scope", scope)
	}
	if userID := strings.TrimSpace(request.UserID); userID != "" {
		query.Set("userId", userID)
	}
	if heroID := strings.TrimSpace(request.HeroID); heroID != "" {
		query.Set("heroId", heroID)
	}
	if format := strings.TrimSpace(request.Format); format != "" {
		query.Set("format", format)
	}
	if name := strings.TrimSpace(request.Name); name != "" {
		query.Set("name", name)
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	if subjectID := strings.TrimSpace(request.SubjectID); subjectID != "" {
		query.Set("subjectId", subjectID)
	}
	req.URL.RawQuery = query.Encode()

	response := &ListDecksResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateDeck provisions a new deck.
func (c *Client) CreateDeck(ctx context.Context, request *CreateDeckRequest, opts ...RequestOpt) (*CreateDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/decks", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CreateDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CloneDeck clones a deck from an existing deck version.
func (c *Client) CloneDeck(ctx context.Context, request *CloneDeckRequest, opts ...RequestOpt) (*CloneDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(struct {
		SourceDeckVersionID  string           `json:"sourceDeckVersionId,omitempty"`
		Name                 string           `json:"name,omitempty"`
		Visibility           *VisibilityLevel `json:"visibility,omitempty"`
		DeckID               string           `json:"deckId,omitempty"`
		InitialVersionName   string           `json:"initialVersionName,omitempty"`
		InitialDeckVersionID string           `json:"initialDeckVersionId,omitempty"`
	}{
		SourceDeckVersionID:  strings.TrimSpace(request.SourceDeckVersionID),
		Name:                 strings.TrimSpace(request.Name),
		Visibility:           request.Visibility,
		DeckID:               strings.TrimSpace(request.DeckID),
		InitialVersionName:   strings.TrimSpace(request.InitialVersionName),
		InitialDeckVersionID: strings.TrimSpace(request.InitialDeckVersionID),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/decks:clone", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CloneDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// SearchDecks searches for decks matching provided filters.
func (c *Client) SearchDecks(ctx context.Context, request *SearchDecksRequest, opts ...RequestOpt) (*SearchDecksResponse, error) {
	if request == nil {
		request = &SearchDecksRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/decks/search", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if heroID := strings.TrimSpace(request.HeroID); heroID != "" {
		query.Set("heroId", heroID)
	}
	if format := strings.TrimSpace(request.Format); format != "" {
		query.Set("format", format)
	}
	if term := strings.TrimSpace(request.SearchTerm); term != "" {
		query.Set("searchTerm", term)
	}
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	if deckKind := strings.TrimSpace(string(request.DeckKind)); deckKind != "" && deckKind != string(DeckKindUnspecified) {
		query.Set("deckKind", deckKind)
	}
	if sourceKind := strings.TrimSpace(string(request.SourceKind)); sourceKind != "" && sourceKind != string(DeckSourceKindUnspecified) {
		query.Set("sourceKind", sourceKind)
	}
	if sourceReference := strings.TrimSpace(request.SourceReference); sourceReference != "" {
		query.Set("sourceReference", sourceReference)
	}
	req.URL.RawQuery = query.Encode()

	response := &SearchDecksResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeck retrieves a deck by ID.
func (c *Client) GetDeck(ctx context.Context, request *GetDeckRequest, opts ...RequestOpt) (*GetDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeckAccess returns the caller's effective permission for a deck.
func (c *Client) GetDeckAccess(ctx context.Context, request *GetDeckAccessRequest, opts ...RequestOpt) (*GetDeckAccessResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/access", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetDeckAccessResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// StopDeckShare removes all explicit shares for a deck.
func (c *Client) StopDeckShare(ctx context.Context, request *StopDeckShareRequest, opts ...RequestOpt) (*StopDeckShareResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	body, err := jsonBody(struct{}{})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/decks/%s/access:stop", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &StopDeckShareResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListDeckAccessGrants lists subjects explicitly granted access to a deck.
func (c *Client) ListDeckAccessGrants(ctx context.Context, request *ListDeckAccessGrantsRequest, opts ...RequestOpt) (*ListDeckAccessGrantsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/permissions", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	req.URL.RawQuery = query.Encode()

	response := &ListDeckAccessGrantsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteDeck removes a deck by ID.
func (c *Client) DeleteDeck(ctx context.Context, request *DeleteDeckRequest, opts ...RequestOpt) (*DeleteDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateDeck applies changes to an existing deck.
func (c *Client) UpdateDeck(ctx context.Context, request *UpdateDeckRequest, opts ...RequestOpt) (*UpdateDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	if request.Name == nil && request.Author == nil && request.ActiveDeckVersionID == nil {
		return nil, errors.New("at least one field must be set")
	}
	body, err := jsonBody(struct {
		Name                *string `json:"name,omitempty"`
		Author              *string `json:"author,omitempty"`
		ActiveDeckVersionID *string `json:"activeDeckVersionId,omitempty"`
	}{
		Name:                request.Name,
		Author:              request.Author,
		ActiveDeckVersionID: request.ActiveDeckVersionID,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/decks/%s", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateDeckVisibility updates only the visibility of an existing deck.
func (c *Client) UpdateDeckVisibility(ctx context.Context, request *UpdateDeckVisibilityRequest, opts ...RequestOpt) (*UpdateDeckVisibilityResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	if request.Visibility == nil {
		return nil, errors.New("visibility must not be nil")
	}

	body, err := jsonBody(struct {
		Visibility *VisibilityLevel `json:"visibility,omitempty"`
	}{
		Visibility: request.Visibility,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/decks/%s/visibility", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateDeckVisibilityResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GrantDeckAccess assigns a permission for a user to access a deck.
func (c *Client) GrantDeckAccess(ctx context.Context, request *GrantDeckAccessRequest, opts ...RequestOpt) (*GrantDeckAccessResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	subjectID := strings.TrimSpace(request.SubjectID)
	if subjectID == "" {
		return nil, errors.New("subjectID must not be empty")
	}
	permission := DeckPermission(strings.TrimSpace(string(request.Permission)))
	if permission == "" || permission == DeckPermissionUnspecified {
		return nil, errors.New("permission must be specified")
	}

	body, err := jsonBody(struct {
		DeckID     string         `json:"resourceId,omitempty"`
		SubjectID  string         `json:"subjectId,omitempty"`
		Permission DeckPermission `json:"permission,omitempty"`
	}{
		DeckID:     deckID,
		SubjectID:  subjectID,
		Permission: permission,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/decks/permissions:grant", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &GrantDeckAccessResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// RevokeDeckAccess removes a previously granted permission from a deck.
func (c *Client) RevokeDeckAccess(ctx context.Context, request *RevokeDeckAccessRequest, opts ...RequestOpt) (*RevokeDeckAccessResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}
	subjectID := strings.TrimSpace(request.SubjectID)
	if subjectID == "" {
		return nil, errors.New("subjectID must not be empty")
	}
	permission := DeckPermission(strings.TrimSpace(string(request.Permission)))
	if permission == "" || permission == DeckPermissionUnspecified {
		return nil, errors.New("permission must be specified")
	}

	body, err := jsonBody(struct {
		DeckID     string         `json:"resourceId,omitempty"`
		SubjectID  string         `json:"subjectId,omitempty"`
		Permission DeckPermission `json:"permission,omitempty"`
	}{
		DeckID:     deckID,
		SubjectID:  subjectID,
		Permission: permission,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/decks/permissions:revoke", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &RevokeDeckAccessResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// StarDeck marks a deck as starred for the current user.
func (c *Client) StarDeck(ctx context.Context, request *StarDeckRequest, opts ...RequestOpt) (*StarDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/stars", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodPost, path, nil)
	if err != nil {
		return nil, err
	}

	response := &StarDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UnstarDeck removes the starred marker from a deck for the current user.
func (c *Client) UnstarDeck(ctx context.Context, request *UnstarDeckRequest, opts ...RequestOpt) (*UnstarDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/stars", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &UnstarDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListDeckVersions enumerates versions for a deck.
func (c *Client) ListDeckVersions(ctx context.Context, request *ListDeckVersionsRequest, opts ...RequestOpt) (*ListDeckVersionsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s/versions", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	if token := strings.TrimSpace(request.NextToken); token != "" {
		query.Set("nextToken", token)
	}
	req.URL.RawQuery = query.Encode()

	response := &ListDeckVersionsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// CreateDeckVersion provisions a new deck version.
func (c *Client) CreateDeckVersion(ctx context.Context, request *CreateDeckVersionRequest, opts ...RequestOpt) (*CreateDeckVersionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	body, err := jsonBody(struct {
		Name                string `json:"name,omitempty"`
		DeckVersionID       string `json:"deckVersionId,omitempty"`
		SourceDeckVersionID string `json:"sourceDeckVersionId,omitempty"`
		SetActive           *bool  `json:"setActive,omitempty"`
	}{
		Name:                strings.TrimSpace(request.Name),
		DeckVersionID:       strings.TrimSpace(request.DeckVersionID),
		SourceDeckVersionID: strings.TrimSpace(request.SourceDeckVersionID),
		SetActive:           request.SetActive,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/decks/%s/versions", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &CreateDeckVersionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeckVersion retrieves a specific deck version.
func (c *Client) GetDeckVersion(ctx context.Context, request *GetDeckVersionRequest, opts ...RequestOpt) (*GetDeckVersionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}

	path := fmt.Sprintf("/v1/deck_versions/%s", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetDeckVersionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteDeckVersion removes a deck version.
func (c *Client) DeleteDeckVersion(ctx context.Context, request *DeleteDeckVersionRequest, opts ...RequestOpt) (*DeleteDeckVersionResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}

	path := fmt.Sprintf("/v1/deck_versions/%s", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	response := &DeleteDeckVersionResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListDeckVersionSideboardGuides lists sideboard guides for a deck version.
func (c *Client) ListDeckVersionSideboardGuides(ctx context.Context, request *ListDeckVersionSideboardGuidesRequest, opts ...RequestOpt) (*ListDeckVersionSideboardGuidesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/sideboard_guides", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &ListDeckVersionSideboardGuidesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpsertDeckVersionSideboardGuide creates or updates a sideboard guide.
func (c *Client) UpsertDeckVersionSideboardGuide(ctx context.Context, request *UpsertDeckVersionSideboardGuideRequest, opts ...RequestOpt) (*UpsertDeckVersionSideboardGuideResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}
	targetType := strings.TrimSpace(string(request.TargetType))
	if targetType == "" || targetType == string(SideboardGuideTargetTypeUnspecified) {
		return nil, errors.New("targetType must not be empty")
	}
	target := strings.TrimSpace(request.Target)
	if target == "" {
		return nil, errors.New("target must not be empty")
	}

	body, err := jsonBody(struct {
		TargetType    SideboardGuideTargetType   `json:"targetType,omitempty"`
		Target        string                     `json:"target,omitempty"`
		Guide         string                     `json:"guide,omitempty"`
		CardsToAdd    []SideboardGuideCardChange `json:"cardsToAdd,omitempty"`
		CardsToRemove []SideboardGuideCardChange `json:"cardsToRemove,omitempty"`
	}{
		TargetType:    request.TargetType,
		Target:        target,
		Guide:         request.Guide,
		CardsToAdd:    request.CardsToAdd,
		CardsToRemove: request.CardsToRemove,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/sideboard_guides", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpsertDeckVersionSideboardGuideResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// DeleteDeckVersionSideboardGuide removes a sideboard guide.
func (c *Client) DeleteDeckVersionSideboardGuide(ctx context.Context, request *DeleteDeckVersionSideboardGuideRequest, opts ...RequestOpt) (*DeleteDeckVersionSideboardGuideResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}
	targetType := strings.TrimSpace(string(request.TargetType))
	if targetType == "" || targetType == string(SideboardGuideTargetTypeUnspecified) {
		return nil, errors.New("targetType must not be empty")
	}
	target := strings.TrimSpace(request.Target)
	if target == "" {
		return nil, errors.New("target must not be empty")
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/sideboard_guides", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("targetType", targetType)
	query.Set("target", target)
	req.URL.RawQuery = query.Encode()

	response := &DeleteDeckVersionSideboardGuideResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListDeckVersionCards returns cards contained within a deck version.
func (c *Client) ListDeckVersionCards(ctx context.Context, request *ListDeckVersionCardsRequest, opts ...RequestOpt) (*ListDeckVersionCardsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/cards", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &ListDeckVersionCardsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ModifyDeckVersionCard adjusts a card entry within a deck version.
func (c *Client) ModifyDeckVersionCard(ctx context.Context, request *ModifyDeckVersionCardRequest, opts ...RequestOpt) (*ModifyDeckVersionCardResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	cardID := strings.TrimSpace(request.CardID)
	board := strings.TrimSpace(string(request.Board))
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}
	if cardID == "" {
		return nil, errors.New("cardID must not be empty")
	}
	if board == "" || board == string(BoardTypeUnspecified) {
		return nil, errors.New("boardType must not be empty")
	}

	body, err := jsonBody(struct {
		CardID   string `json:"cardId,omitempty"`
		Board    string `json:"boardType,omitempty"`
		Quantity int32  `json:"quantity,omitempty"`
	}{
		CardID:   cardID,
		Board:    board,
		Quantity: request.Quantity,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/cards", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &ModifyDeckVersionCardResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeckVersionHistory retrieves change history for a deck version.
func (c *Client) GetDeckVersionHistory(ctx context.Context, request *GetDeckVersionHistoryRequest, opts ...RequestOpt) (*GetDeckVersionHistoryResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/history", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetDeckVersionHistoryResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDeckVersionNotes fetches notes for a deck version.
func (c *Client) GetDeckVersionNotes(ctx context.Context, request *GetDeckVersionNotesRequest, opts ...RequestOpt) (*GetDeckVersionNotesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/notes", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetDeckVersionNotesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// UpdateDeckVersionNotes sets notes for a deck version.
func (c *Client) UpdateDeckVersionNotes(ctx context.Context, request *UpdateDeckVersionNotesRequest, opts ...RequestOpt) (*UpdateDeckVersionNotesResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	deckVersionID := strings.TrimSpace(request.DeckVersionID)
	if deckVersionID == "" {
		return nil, errors.New("deckVersionID must not be empty")
	}

	body, err := jsonBody(struct {
		Notes string `json:"notes,omitempty"`
	}{
		Notes: request.Notes,
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/deck_versions/%s/notes", url.PathEscape(deckVersionID))
	req, err := c.newRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &UpdateDeckVersionNotesResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchGetDecks fetches decks in bulk.
func (c *Client) BatchGetDecks(ctx context.Context, request *BatchGetDecksRequest, opts ...RequestOpt) (*BatchGetDecksResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if len(request.DeckIDs) == 0 {
		return nil, errors.New("deckIds must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/decks:batchGet", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &BatchGetDecksResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ExportDeck exports a deck snapshot including versions.
func (c *Client) ExportDeck(ctx context.Context, request *ExportDeckRequest, opts ...RequestOpt) (*ExportDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	deckID := strings.TrimSpace(request.DeckID)
	if deckID == "" {
		return nil, errors.New("deckID must not be empty")
	}

	path := fmt.Sprintf("/v1/decks/%s:export", url.PathEscape(deckID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &ExportDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ImportDeck imports a deck snapshot including versions.
func (c *Client) ImportDeck(ctx context.Context, request *ImportDeckRequest, opts ...RequestOpt) (*ImportDeckResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(struct {
		DeckID              string              `json:"deckId,omitempty"`
		Name                string              `json:"name,omitempty"`
		HeroID              string              `json:"heroId,omitempty"`
		Format              string              `json:"format,omitempty"`
		Author              string              `json:"author,omitempty"`
		Visibility          *VisibilityLevel    `json:"visibility,omitempty"`
		CreatedAt           *time.Time          `json:"createdAt,omitempty"`
		Versions            []ImportDeckVersion `json:"versions,omitempty"`
		ActiveDeckVersionID string              `json:"activeDeckVersionId,omitempty"`
	}{
		DeckID:              strings.TrimSpace(request.DeckID),
		Name:                strings.TrimSpace(request.Name),
		HeroID:              strings.TrimSpace(request.HeroID),
		Format:              strings.TrimSpace(request.Format),
		Author:              strings.TrimSpace(request.Author),
		Visibility:          request.Visibility,
		CreatedAt:           request.CreatedAt,
		Versions:            request.Versions,
		ActiveDeckVersionID: strings.TrimSpace(request.ActiveDeckVersionID),
	})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/decks:import", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &ImportDeckResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListStarredDecks lists decks starred by a specific user.
func (c *Client) ListStarredDecks(ctx context.Context, request *ListStarredDecksRequest, opts ...RequestOpt) (*ListStarredDecksResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	userID := strings.TrimSpace(request.UserID)
	if userID == "" {
		return nil, errors.New("userID must not be empty")
	}

	path := fmt.Sprintf("/v1/users/%s/stars/decks", url.PathEscape(userID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.Limit != nil && *request.Limit >= 0 {
		query.Set("limit", strconv.Itoa(int(*request.Limit)))
	}
	if request.Offset != nil && *request.Offset >= 0 {
		query.Set("offset", strconv.Itoa(int(*request.Offset)))
	}
	req.URL.RawQuery = query.Encode()

	response := &ListStarredDecksResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}
