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

// ColorIdentity represents v1ColorIdentity.
type ColorIdentity string

const (
	ColorIdentityUnspecified ColorIdentity = "COLOR_IDENTITY_UNSPECIFIED"
	ColorIdentityRed         ColorIdentity = "COLOR_IDENTITY_RED"
	ColorIdentityYellow      ColorIdentity = "COLOR_IDENTITY_YELLOW"
	ColorIdentityBlue        ColorIdentity = "COLOR_IDENTITY_BLUE"
	ColorIdentityNone        ColorIdentity = "COLOR_IDENTITY_NONE"
)

// LegalityStatus represents v1LegalityStatus.
type LegalityStatus string

const (
	LegalityStatusUnspecified  LegalityStatus = "LEGALITY_STATUS_UNSPECIFIED"
	LegalityStatusLegal        LegalityStatus = "LEGALITY_STATUS_LEGAL"
	LegalityStatusNotLegal     LegalityStatus = "LEGALITY_STATUS_NOT_LEGAL"
	LegalityStatusBanned       LegalityStatus = "LEGALITY_STATUS_BANNED"
	LegalityStatusSuspended    LegalityStatus = "LEGALITY_STATUS_SUSPENDED"
	LegalityStatusRestricted   LegalityStatus = "LEGALITY_STATUS_RESTRICTED"
	LegalityStatusLivingLegend LegalityStatus = "LEGALITY_STATUS_LIVING_LEGEND"
)

// Edition mirrors cardsv1Edition.
type Edition string

const (
	EditionUnspecified Edition = "EDITION_UNSPECIFIED"
	EditionNormal      Edition = "NORMAL"
	EditionAlpha       Edition = "ALPHA"
	EditionFirst       Edition = "FIRST"
	EditionUnlimited   Edition = "UNLIMITED"
)

// Foiling matches v1Foiling.
type Foiling string

const (
	FoilingUnspecified Foiling = "FOILING_UNSPECIFIED"
	FoilingNone        Foiling = "NONE"
	FoilingRainbow     Foiling = "RAINBOW"
	FoilingCold        Foiling = "COLD"
	FoilingGold        Foiling = "GOLD"
)

// Rarity matches v1Rarity.
type Rarity string

const (
	RarityUnspecified Rarity = "RARITY_UNSPECIFIED"
	RarityCommon      Rarity = "COMMON"
	RarityRare        Rarity = "RARE"
	RaritySuperRare   Rarity = "SUPER_RARE"
	RarityMajestic    Rarity = "MAJESTIC"
	RarityLegendary   Rarity = "LEGENDARY"
	RarityFabled      Rarity = "FABLED"
	RarityToken       Rarity = "TOKEN"
	RarityMarvel      Rarity = "MARVEL"
	RarityPromo       Rarity = "PROMO"
	RarityBasic       Rarity = "BASIC"
)

// FormatCardLegalitySummary models v1FormatCardLegalitySummary.
type FormatCardLegalitySummary struct {
	Status         LegalityStatus `json:"status,omitempty"`
	IsLegal        bool           `json:"isLegal,omitempty"`
	IsLivingLegend bool           `json:"isLivingLegend,omitempty"`
	LivingLegendAt *time.Time     `json:"livingLegendAt,omitempty"`
	Banned         bool           `json:"banned,omitempty"`
	BannedAt       *time.Time     `json:"bannedAt,omitempty"`
	Suspended      bool           `json:"suspended,omitempty"`
	SuspendedAt    *time.Time     `json:"suspendedAt,omitempty"`
	SuspendedEnd   *time.Time     `json:"suspendedEnd,omitempty"`
	Restricted     bool           `json:"restricted,omitempty"`
	RestrictedAt   *time.Time     `json:"restrictedAt,omitempty"`
}

// CardLegalitySummary aggregates format-specific legality details.
type CardLegalitySummary struct {
	Blitz              *FormatCardLegalitySummary `json:"blitz,omitempty"`
	ClassicConstructed *FormatCardLegalitySummary `json:"classicConstructed,omitempty"`
	Commoner           *FormatCardLegalitySummary `json:"commoner,omitempty"`
	UltimatePitFight   *FormatCardLegalitySummary `json:"ultimatePitFight,omitempty"`
	LivingLegend       *FormatCardLegalitySummary `json:"livingLegend,omitempty"`
	ProjectBlue        *FormatCardLegalitySummary `json:"projectBlue,omitempty"`
}

// CardSummary mirrors v1CardSummary.
type CardSummary struct {
	Identifier               string               `json:"identifier,omitempty"`
	Name                     string               `json:"name,omitempty"`
	Cost                     string               `json:"cost,omitempty"`
	Pitch                    string               `json:"pitch,omitempty"`
	Power                    string               `json:"power,omitempty"`
	Defense                  string               `json:"defense,omitempty"`
	Health                   string               `json:"health,omitempty"`
	Intelligence             string               `json:"intelligence,omitempty"`
	ColorIdentity            ColorIdentity        `json:"colorIdentity,omitempty"`
	Arcane                   string               `json:"arcane,omitempty"`
	Types                    []string             `json:"types,omitempty"`
	Keywords                 []string             `json:"keywords,omitempty"`
	AbilitiesAndEffects      []string             `json:"abilitiesAndEffects,omitempty"`
	AbilityAndEffectKeywords []string             `json:"abilityAndEffectKeywords,omitempty"`
	GrantedKeywords          []string             `json:"grantedKeywords,omitempty"`
	RemovedKeywords          []string             `json:"removedKeywords,omitempty"`
	InteractsWithKeywords    []string             `json:"interactsWithKeywords,omitempty"`
	FunctionalText           string               `json:"functionalText,omitempty"`
	IsDoubleFacedCard        bool                 `json:"isDoubleFacedCard,omitempty"`
	IsDoubleFacedFront       bool                 `json:"isDoubleFacedFront,omitempty"`
	DoubleFacedOtherID       string               `json:"doubleFacedOtherId,omitempty"`
	DefaultImageURL          string               `json:"defaultImageUrl,omitempty"`
	PitchSiblingIDs          []string             `json:"pitchSiblingIds,omitempty"`
	ReferencedCards          []string             `json:"referencedCards,omitempty"`
	CardsReferencedBy        []string             `json:"cardsReferencedBy,omitempty"`
	Legality                 *CardLegalitySummary `json:"legality,omitempty"`
}

// TCGPlayerSummary mirrors v1TCGPlayerSummary.
type TCGPlayerSummary struct {
	ProductID string `json:"productId,omitempty"`
	URL       string `json:"url,omitempty"`
}

// ProductSummary mirrors v1ProductSummary.
type ProductSummary struct {
	Identifier      string `json:"identifier,omitempty"`
	FrontCardID     string `json:"frontCardId,omitempty"`
	FrontPrintingID string `json:"frontPrintingId,omitempty"`
	BackCardID      string `json:"backCardId,omitempty"`
	BackPrintingID  string `json:"backPrintingId,omitempty"`
	IsDFC           bool   `json:"isDfc,omitempty"`
	Type            string `json:"type,omitempty"`
	CardID          string `json:"cardId,omitempty"`
	PrintingID      string `json:"printingId,omitempty"`
	ProductGroupID  string `json:"productGroupId,omitempty"`
	Name            string `json:"name,omitempty"`
	Slug            string `json:"slug,omitempty"`
	PrintedDate     string `json:"printedDate,omitempty"`
	PrintedLanguage string `json:"printedLanguage,omitempty"`
	ReleaseDate     string `json:"releaseDate,omitempty"`
	Description     string `json:"description,omitempty"`
	Quantity        int32  `json:"quantity,omitempty"`
}

// SetSummary mirrors v1SetSummary.
type SetSummary struct {
	Code        string `json:"code,omitempty"`
	Name        string `json:"name,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
}

// PrintingSummary mirrors v1PrintingSummary.
type PrintingSummary struct {
	Identifier      string            `json:"identifier,omitempty"`
	CardID          string            `json:"cardId,omitempty"`
	SetPrintingID   string            `json:"setPrintingId,omitempty"`
	PrintingName    string            `json:"printingName,omitempty"`
	Artists         []string          `json:"artists,omitempty"`
	ArtVariations   []string          `json:"artVariations,omitempty"`
	FlavorText      string            `json:"flavorText,omitempty"`
	ImageURL        string            `json:"imageUrl,omitempty"`
	SetID           string            `json:"setId,omitempty"`
	SetName         string            `json:"setName,omitempty"`
	Edition         Edition           `json:"edition,omitempty"`
	IsExpansionSlot bool              `json:"isExpansionSlot,omitempty"`
	Foiling         Foiling           `json:"foiling,omitempty"`
	Rarity          Rarity            `json:"rarity,omitempty"`
	TCGPlayer       *TCGPlayerSummary `json:"tcgPlayer,omitempty"`
	Products        []ProductSummary  `json:"products,omitempty"`
}

// DataSnapshotFile mirrors v1DataSnapshotFile.
type DataSnapshotFile struct {
	Name            string `json:"name,omitempty"`
	URL             string `json:"url,omitempty"`
	SizeBytes       string `json:"sizeBytes,omitempty"`
	SHA256          string `json:"sha256,omitempty"`
	ContentEncoding string `json:"contentEncoding,omitempty"`
	ContentType     string `json:"contentType,omitempty"`
	RowCount        int32  `json:"rowCount,omitempty"`
}

// DataSnapshotManifest mirrors v1DataSnapshotManifest.
type DataSnapshotManifest struct {
	Version       string             `json:"version,omitempty"`
	SchemaVersion int32              `json:"schemaVersion,omitempty"`
	CreatedAt     string             `json:"createdAt,omitempty"`
	MinAppVersion string             `json:"minAppVersion,omitempty"`
	Files         []DataSnapshotFile `json:"files,omitempty"`
}

// SearchCardsRequest models query filters for card search.
type SearchCardsRequest struct {
	SearchTerm           string
	Class                string
	Type                 string
	Subtype              string
	Talent               string
	Cost                 string
	Defense              string
	Pitch                string
	Power                string
	Health               string
	Intelligence         string
	Arcane               string
	ColorIdentity        string
	BlitzLegal           *bool
	BlitzBanned          *bool
	BlitzSuspended       *bool
	BlitzLivingLegend    *bool
	CCLegal              *bool
	CCBanned             *bool
	CCSuspended          *bool
	CCLivingLegend       *bool
	CommonerLegal        *bool
	CommonerBanned       *bool
	CommonerSuspended    *bool
	UPFBanned            *bool
	LLBanned             *bool
	LLRestricted         *bool
	ProjectBlueLegal     *bool
	ProjectBlueBanned    *bool
	ProjectBlueSuspended *bool
	IsDoubleFaced        *bool
	PageSize             *int32
	NextToken            string
}

// SearchCardsResponse mirrors v1SearchCardsResponse.
type SearchCardsResponse struct {
	Summaries []CardSummary    `json:"summaries,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *SearchCardsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetCardRequest identifies a card to fetch.
type GetCardRequest struct {
	CardID string `json:"-"`
}

// GetCardResponse contains the requested card summary.
type GetCardResponse struct {
	Summary  *CardSummary     `json:"summary,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetCardResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetCardsRequest fetches multiple cards by ID.
type BatchGetCardsRequest struct {
	CardIDs      []string `json:"cardIds,omitempty"`
	AllowPartial bool     `json:"allowPartial,omitempty"`
}

// BatchGetCardsResponse holds bulk card summaries.
type BatchGetCardsResponse struct {
	Cards       map[string]CardSummary `json:"cards,omitempty"`
	NotFoundIDs []string               `json:"notFoundIds,omitempty"`
	Metadata    ResponseMetadata       `json:"-"`
}

func (r *BatchGetCardsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListPrintingsRequest enumerates printings for a card.
type ListPrintingsRequest struct {
	CardID    string `json:"-"`
	PageSize  *int32
	NextToken string
}

// ListPrintingsResponse mirrors v1ListPrintingsResponse.
type ListPrintingsResponse struct {
	Summaries []PrintingSummary `json:"summaries,omitempty"`
	NextToken string            `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata  `json:"-"`
}

func (r *ListPrintingsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListPrintingsForSetNumberRequest enumerates printings by set number.
type ListPrintingsForSetNumberRequest struct {
	SetNumber string `json:"-"`
}

// ListPrintingsForSetNumberResponse mirrors v1ListPrintingsForSetNumberResponse.
type ListPrintingsForSetNumberResponse struct {
	Summaries []PrintingSummary `json:"summaries,omitempty"`
	Metadata  ResponseMetadata  `json:"-"`
}

func (r *ListPrintingsForSetNumberResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetPrintingRequest identifies a printing to fetch.
type GetPrintingRequest struct {
	PrintingID string `json:"-"`
}

// GetPrintingResponse contains a printing summary.
type GetPrintingResponse struct {
	Summary  *PrintingSummary `json:"summary,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetPrintingResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetPrintingsRequest fetches printings in bulk.
type BatchGetPrintingsRequest struct {
	PrintingIDs  []string `json:"printingIds,omitempty"`
	AllowPartial bool     `json:"allowPartial,omitempty"`
}

// BatchGetPrintingsResponse returns bulk printing summaries.
type BatchGetPrintingsResponse struct {
	Printings   map[string]PrintingSummary `json:"printings,omitempty"`
	NotFoundIDs []string                   `json:"notFoundIds,omitempty"`
	Metadata    ResponseMetadata           `json:"-"`
}

func (r *BatchGetPrintingsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetProductRequest identifies a product to fetch.
type GetProductRequest struct {
	ProductID string `json:"-"`
}

// GetProductResponse contains a product summary.
type GetProductResponse struct {
	Summary  *ProductSummary  `json:"summary,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetProductResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetProductsRequest fetches products in bulk.
type BatchGetProductsRequest struct {
	ProductIDs   []string `json:"productIds,omitempty"`
	AllowPartial bool     `json:"allowPartial,omitempty"`
}

// BatchGetProductsResponse returns bulk product summaries.
type BatchGetProductsResponse struct {
	Products    map[string]ProductSummary `json:"products,omitempty"`
	NotFoundIDs []string                  `json:"notFoundIds,omitempty"`
	Metadata    ResponseMetadata          `json:"-"`
}

func (r *BatchGetProductsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListProductsRequest enumerates products with optional filters.
type ListProductsRequest struct {
	Type           string
	SetCode        string
	ProductGroupID string
	CardID         string
	PrintingID     string
	PageSize       *int32
	NextToken      string
}

// ListProductsResponse mirrors v1ListProductsResponse.
type ListProductsResponse struct {
	Summaries []ProductSummary `json:"summaries,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListProductsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetSetRequest identifies a set by code.
type GetSetRequest struct {
	SetCode string `json:"-"`
}

// GetSetResponse returns a set summary.
type GetSetResponse struct {
	Summary  *SetSummary      `json:"summary,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetSetResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListSetsRequest enumerates sets.
type ListSetsRequest struct {
	PageSize  *int32
	NextToken string
}

// ListSetsResponse mirrors v1ListSetsResponse.
type ListSetsResponse struct {
	Summaries []SetSummary     `json:"summaries,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListSetsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// BatchGetSetsRequest fetches sets in bulk.
type BatchGetSetsRequest struct {
	SetCodes     []string `json:"setCodes,omitempty"`
	AllowPartial bool     `json:"allowPartial,omitempty"`
}

// BatchGetSetsResponse returns bulk set summaries.
type BatchGetSetsResponse struct {
	Sets          map[string]SetSummary `json:"sets,omitempty"`
	NotFoundCodes []string              `json:"notFoundCodes,omitempty"`
	Metadata      ResponseMetadata      `json:"-"`
}

func (r *BatchGetSetsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetDataSnapshotRequest fetches card data snapshots.
type GetDataSnapshotRequest struct {
	SchemaVersion *int32
	Version       string
}

// GetDataSnapshotResponse returns snapshot metadata.
type GetDataSnapshotResponse struct {
	Manifest *DataSnapshotManifest `json:"manifest,omitempty"`
	Metadata ResponseMetadata      `json:"-"`
}

func (r *GetDataSnapshotResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// SearchCards performs a filtered card search.
func (c *Client) SearchCards(ctx context.Context, request *SearchCardsRequest, opts ...RequestOpt) (*SearchCardsResponse, error) {
	if request == nil {
		request = &SearchCardsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/cards", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "searchTerm", request.SearchTerm)
	setQueryString(query, "class", request.Class)
	setQueryString(query, "type", request.Type)
	setQueryString(query, "subtype", request.Subtype)
	setQueryString(query, "talent", request.Talent)
	setQueryString(query, "cost", request.Cost)
	setQueryString(query, "defense", request.Defense)
	setQueryString(query, "pitch", request.Pitch)
	setQueryString(query, "power", request.Power)
	setQueryString(query, "health", request.Health)
	setQueryString(query, "intelligence", request.Intelligence)
	setQueryString(query, "arcane", request.Arcane)
	setQueryString(query, "colorIdentity", request.ColorIdentity)

	setQueryBool(query, "blitzLegal", request.BlitzLegal)
	setQueryBool(query, "blitzBanned", request.BlitzBanned)
	setQueryBool(query, "blitzSuspended", request.BlitzSuspended)
	setQueryBool(query, "blitzLivingLegend", request.BlitzLivingLegend)
	setQueryBool(query, "ccLegal", request.CCLegal)
	setQueryBool(query, "ccBanned", request.CCBanned)
	setQueryBool(query, "ccSuspended", request.CCSuspended)
	setQueryBool(query, "ccLivingLegend", request.CCLivingLegend)
	setQueryBool(query, "commonerLegal", request.CommonerLegal)
	setQueryBool(query, "commonerBanned", request.CommonerBanned)
	setQueryBool(query, "commonerSuspended", request.CommonerSuspended)
	setQueryBool(query, "upfBanned", request.UPFBanned)
	setQueryBool(query, "llBanned", request.LLBanned)
	setQueryBool(query, "llRestricted", request.LLRestricted)
	setQueryBool(query, "projectBlueLegal", request.ProjectBlueLegal)
	setQueryBool(query, "projectBlueBanned", request.ProjectBlueBanned)
	setQueryBool(query, "projectBlueSuspended", request.ProjectBlueSuspended)
	setQueryBool(query, "isDoubleFaced", request.IsDoubleFaced)

	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	setQueryString(query, "nextToken", request.NextToken)
	req.URL.RawQuery = query.Encode()

	response := &SearchCardsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetCard retrieves a card summary by ID.
func (c *Client) GetCard(ctx context.Context, request *GetCardRequest, opts ...RequestOpt) (*GetCardResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CardID) == "" {
		return nil, errors.New("cardID must not be empty")
	}

	path := fmt.Sprintf("/v1/cards/%s", url.PathEscape(request.CardID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetCardResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchGetCards fetches multiple cards in one call.
func (c *Client) BatchGetCards(ctx context.Context, request *BatchGetCardsRequest, opts ...RequestOpt) (*BatchGetCardsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/cards:batchGet", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &BatchGetCardsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListPrintings lists printings for a specific card.
func (c *Client) ListPrintings(ctx context.Context, request *ListPrintingsRequest, opts ...RequestOpt) (*ListPrintingsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.CardID) == "" {
		return nil, errors.New("cardID must not be empty")
	}

	path := fmt.Sprintf("/v1/cards/%s/printings", url.PathEscape(request.CardID))
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

	response := &ListPrintingsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListPrintingsForSetNumber lists printings for a set number.
func (c *Client) ListPrintingsForSetNumber(ctx context.Context, request *ListPrintingsForSetNumberRequest, opts ...RequestOpt) (*ListPrintingsForSetNumberResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.SetNumber) == "" {
		return nil, errors.New("setNumber must not be empty")
	}

	path := fmt.Sprintf("/v1/printings/set/%s", url.PathEscape(request.SetNumber))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &ListPrintingsForSetNumberResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetPrinting fetches a printing by ID.
func (c *Client) GetPrinting(ctx context.Context, request *GetPrintingRequest, opts ...RequestOpt) (*GetPrintingResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.PrintingID) == "" {
		return nil, errors.New("printingID must not be empty")
	}

	path := fmt.Sprintf("/v1/printings/%s", url.PathEscape(request.PrintingID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetPrintingResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchGetPrintings fetches multiple printings at once.
func (c *Client) BatchGetPrintings(ctx context.Context, request *BatchGetPrintingsRequest, opts ...RequestOpt) (*BatchGetPrintingsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/printings:batchGet", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &BatchGetPrintingsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetProduct fetches a product summary.
func (c *Client) GetProduct(ctx context.Context, request *GetProductRequest, opts ...RequestOpt) (*GetProductResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.ProductID) == "" {
		return nil, errors.New("productID must not be empty")
	}

	path := fmt.Sprintf("/v1/products/%s", url.PathEscape(request.ProductID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetProductResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchGetProducts fetches multiple products.
func (c *Client) BatchGetProducts(ctx context.Context, request *BatchGetProductsRequest, opts ...RequestOpt) (*BatchGetProductsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/products:batchGet", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &BatchGetProductsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListProducts lists product summaries.
func (c *Client) ListProducts(ctx context.Context, request *ListProductsRequest, opts ...RequestOpt) (*ListProductsResponse, error) {
	if request == nil {
		request = &ListProductsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/products", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryString(query, "type", request.Type)
	setQueryString(query, "setCode", request.SetCode)
	setQueryString(query, "productGroupId", request.ProductGroupID)
	setQueryString(query, "cardId", request.CardID)
	setQueryString(query, "printingId", request.PrintingID)
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	setQueryString(query, "nextToken", request.NextToken)
	req.URL.RawQuery = query.Encode()

	response := &ListProductsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetSet fetches a set summary.
func (c *Client) GetSet(ctx context.Context, request *GetSetRequest, opts ...RequestOpt) (*GetSetResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.SetCode) == "" {
		return nil, errors.New("setCode must not be empty")
	}

	path := fmt.Sprintf("/v1/sets/%s", url.PathEscape(request.SetCode))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetSetResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// ListSets lists set summaries.
func (c *Client) ListSets(ctx context.Context, request *ListSetsRequest, opts ...RequestOpt) (*ListSetsResponse, error) {
	if request == nil {
		request = &ListSetsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/sets", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.PageSize != nil && *request.PageSize > 0 {
		query.Set("pageSize", strconv.Itoa(int(*request.PageSize)))
	}
	setQueryString(query, "nextToken", request.NextToken)
	req.URL.RawQuery = query.Encode()

	response := &ListSetsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// BatchGetSets fetches multiple sets.
func (c *Client) BatchGetSets(ctx context.Context, request *BatchGetSetsRequest, opts ...RequestOpt) (*BatchGetSetsResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/sets:batchGet", body)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	response := &BatchGetSetsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

// GetDataSnapshot retrieves card data snapshot metadata.
func (c *Client) GetDataSnapshot(ctx context.Context, request *GetDataSnapshotRequest, opts ...RequestOpt) (*GetDataSnapshotResponse, error) {
	if request == nil {
		request = &GetDataSnapshotRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/api/v1/cards/data-snapshots", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	if request.SchemaVersion != nil {
		query.Set("schemaVersion", strconv.Itoa(int(*request.SchemaVersion)))
	}
	setQueryString(query, "version", request.Version)
	req.URL.RawQuery = query.Encode()

	response := &GetDataSnapshotResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}

	return response, nil
}

func setQueryString(values url.Values, key, value string) {
	if v := strings.TrimSpace(value); v != "" {
		values.Set(key, v)
	}
}

func setQueryBool(values url.Values, key string, value *bool) {
	if value == nil {
		return
	}
	values.Set(key, strconv.FormatBool(*value))
}
