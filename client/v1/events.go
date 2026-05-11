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

// EventSource represents events.v1.EventSource.
type EventSource string

const (
	EventSourceUnspecified EventSource = "EVENT_SOURCE_UNSPECIFIED"
	EventSourceOfficialGem EventSource = "EVENT_SOURCE_OFFICIAL_GEM"
	EventSourceCommunity   EventSource = "EVENT_SOURCE_COMMUNITY"
)

// EventStatus represents events.v1.EventStatus.
type EventStatus string

const (
	EventStatusUnspecified EventStatus = "EVENT_STATUS_UNSPECIFIED"
	EventStatusPublished   EventStatus = "EVENT_STATUS_PUBLISHED"
	EventStatusCancelled   EventStatus = "EVENT_STATUS_CANCELLED"
	EventStatusHidden      EventStatus = "EVENT_STATUS_HIDDEN"
)

// StoreSource represents events.v1.StoreSource.
type StoreSource string

const (
	StoreSourceUnspecified StoreSource = "STORE_SOURCE_UNSPECIFIED"
	StoreSourceGem         StoreSource = "STORE_SOURCE_GEM"
)

// EventCoordinates represents events.v1.Coordinates.
type EventCoordinates struct {
	Latitude  float64 `json:"latitude,omitempty"`
	Longitude float64 `json:"longitude,omitempty"`
}

// StoreAddress represents events.v1.StoreAddress.
type StoreAddress struct {
	Line1       string            `json:"line1,omitempty"`
	City        string            `json:"city,omitempty"`
	Region      string            `json:"region,omitempty"`
	Postcode    string            `json:"postcode,omitempty"`
	Country     string            `json:"country,omitempty"`
	Coordinates *EventCoordinates `json:"coordinates,omitempty"`
}

// StoreContact represents events.v1.StoreContact.
type StoreContact struct {
	Phone     string `json:"phone,omitempty"`
	Email     string `json:"email,omitempty"`
	Website   string `json:"website,omitempty"`
	Facebook  string `json:"facebook,omitempty"`
	Twitter   string `json:"twitter,omitempty"`
	Instagram string `json:"instagram,omitempty"`
}

// ArmoryDay represents events.v1.ArmoryDay.
type ArmoryDay struct {
	Day  string `json:"day,omitempty"`
	Time string `json:"time,omitempty"`
}

// EventStore represents events.v1.Store.
type EventStore struct {
	StoreID      string        `json:"storeId,omitempty"`
	Source       StoreSource   `json:"source,omitempty"`
	GemID        int64         `json:"gemId,omitempty,string"`
	GemSlug      string        `json:"gemSlug,omitempty"`
	Name         string        `json:"name,omitempty"`
	About        string        `json:"about,omitempty"`
	Address      *StoreAddress `json:"address,omitempty"`
	Contact      *StoreContact `json:"contact,omitempty"`
	ArmoryDays   []ArmoryDay   `json:"armoryDays,omitempty"`
	OnlineStore  bool          `json:"onlineStore,omitempty"`
	OpeningHours string        `json:"openingHours,omitempty"`
	CreatedAt    *time.Time    `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time    `json:"updatedAt,omitempty"`
	LastSeenAt   *time.Time    `json:"lastSeenAt,omitempty"`
}

// Event represents events.v1.Event.
type Event struct {
	EventID            string            `json:"eventId,omitempty"`
	Source             EventSource       `json:"source,omitempty"`
	Status             EventStatus       `json:"status,omitempty"`
	GemEventID         int64             `json:"gemEventId,omitempty,string"`
	StoreID            string            `json:"storeId,omitempty"`
	OrganizerStoreSlug string            `json:"organizerStoreSlug,omitempty"`
	OrganizerName      string            `json:"organizerName,omitempty"`
	Title              string            `json:"title,omitempty"`
	EventType          string            `json:"eventType,omitempty"`
	Format             string            `json:"format,omitempty"`
	StartsAt           *time.Time        `json:"startsAt,omitempty"`
	EndsAt             *time.Time        `json:"endsAt,omitempty"`
	VenueName          string            `json:"venueName,omitempty"`
	Address            string            `json:"address,omitempty"`
	Country            string            `json:"country,omitempty"`
	Coordinates        *EventCoordinates `json:"coordinates,omitempty"`
	ExternalLink       string            `json:"externalLink,omitempty"`
	Description        string            `json:"description,omitempty"`
	DescriptionHTML    string            `json:"descriptionHtml,omitempty"`
	PlayerCap          int32             `json:"playerCap,omitempty"`
	OwnerUserID        string            `json:"ownerUserId,omitempty"`
	CreatedAt          *time.Time        `json:"createdAt,omitempty"`
	UpdatedAt          *time.Time        `json:"updatedAt,omitempty"`
	LastSeenAt         *time.Time        `json:"lastSeenAt,omitempty"`
}

// EventFilterValue represents events.v1.EventFilterValue.
type EventFilterValue struct {
	Kind  string `json:"kind,omitempty"`
	Value string `json:"value,omitempty"`
	Label string `json:"label,omitempty"`
	GemID int64  `json:"gemId,omitempty,string"`
}

// GemLocatorScan represents events.v1.GemLocatorScan.
type GemLocatorScan struct {
	RunID          string     `json:"runId,omitempty"`
	Kind           string     `json:"kind,omitempty"`
	Status         string     `json:"status,omitempty"`
	PagesFetched   int32      `json:"pagesFetched,omitempty"`
	RecordsFetched int32      `json:"recordsFetched,omitempty"`
	StoresUpserted int32      `json:"storesUpserted,omitempty"`
	EventsUpserted int32      `json:"eventsUpserted,omitempty"`
	HTTPStatus     int32      `json:"httpStatus,omitempty"`
	ErrorMessage   string     `json:"errorMessage,omitempty"`
	StartedAt      *time.Time `json:"startedAt,omitempty"`
	FinishedAt     *time.Time `json:"finishedAt,omitempty"`
}

// ListEventsRequest configures event search and pagination.
type ListEventsRequest struct {
	PageSize     *int32
	NextToken    string
	Source       EventSource
	Query        string
	StartsAfter  *time.Time
	StartsBefore *time.Time
	EventType    string
	Format       string
	Country      string
	StoreID      string
	StoreSlug    string
	Latitude     *float64
	Longitude    *float64
	RadiusKM     *float64
}

// ListEventsResponse returns events.
type ListEventsResponse struct {
	Events    []Event          `json:"events,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListEventsResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetEventRequest identifies an event.
type GetEventRequest struct {
	EventID string `json:"-"`
}

// GetEventResponse returns one event.
type GetEventResponse struct {
	Event    *Event           `json:"event,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetEventResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CreateCommunityEventRequest creates a community event.
type CreateCommunityEventRequest struct {
	Title        string            `json:"title,omitempty"`
	EventType    string            `json:"eventType,omitempty"`
	Format       string            `json:"format,omitempty"`
	StartsAt     *time.Time        `json:"startsAt,omitempty"`
	EndsAt       *time.Time        `json:"endsAt,omitempty"`
	VenueName    string            `json:"venueName,omitempty"`
	Address      string            `json:"address,omitempty"`
	Country      string            `json:"country,omitempty"`
	Coordinates  *EventCoordinates `json:"coordinates,omitempty"`
	ExternalLink string            `json:"externalLink,omitempty"`
	Description  string            `json:"description,omitempty"`
	PlayerCap    *int32            `json:"playerCap,omitempty"`
}

// CreateCommunityEventResponse returns the created event.
type CreateCommunityEventResponse struct {
	Event    *Event           `json:"event,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CreateCommunityEventResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UpdateCommunityEventRequest updates a community event.
type UpdateCommunityEventRequest struct {
	EventID      string            `json:"-"`
	Title        *string           `json:"title,omitempty"`
	EventType    *string           `json:"eventType,omitempty"`
	Format       *string           `json:"format,omitempty"`
	StartsAt     *time.Time        `json:"startsAt,omitempty"`
	EndsAt       *time.Time        `json:"endsAt,omitempty"`
	VenueName    *string           `json:"venueName,omitempty"`
	Address      *string           `json:"address,omitempty"`
	Country      *string           `json:"country,omitempty"`
	Coordinates  *EventCoordinates `json:"coordinates,omitempty"`
	ExternalLink *string           `json:"externalLink,omitempty"`
	Description  *string           `json:"description,omitempty"`
	PlayerCap    *int32            `json:"playerCap,omitempty"`
}

// UpdateCommunityEventResponse returns the updated event.
type UpdateCommunityEventResponse struct {
	Event    *Event           `json:"event,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UpdateCommunityEventResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// CancelCommunityEventRequest identifies a community event to cancel.
type CancelCommunityEventRequest struct {
	EventID string `json:"-"`
}

// CancelCommunityEventResponse returns the cancelled event.
type CancelCommunityEventResponse struct {
	Event    *Event           `json:"event,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *CancelCommunityEventResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListStoresRequest configures event store search and pagination.
type ListStoresRequest struct {
	PageSize    *int32
	NextToken   string
	Query       string
	Country     string
	ArmoryDay   string
	OnlineStore *bool
	Latitude    *float64
	Longitude   *float64
	RadiusKM    *float64
}

// ListStoresResponse returns stores.
type ListStoresResponse struct {
	Stores    []EventStore     `json:"stores,omitempty"`
	NextToken string           `json:"nextToken,omitempty"`
	Metadata  ResponseMetadata `json:"-"`
}

func (r *ListStoresResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetStoreRequest identifies a store by id or slug.
type GetStoreRequest struct {
	StoreRef string `json:"-"`
}

// GetStoreResponse returns one store.
type GetStoreResponse struct {
	Store    *EventStore      `json:"store,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetStoreResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListEventFiltersResponse returns event filter values.
type ListEventFiltersResponse struct {
	Filters  []EventFilterValue `json:"filters,omitempty"`
	Metadata ResponseMetadata   `json:"-"`
}

func (r *ListEventFiltersResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// GetGemLocatorScanRequest identifies a GEM locator scan.
type GetGemLocatorScanRequest struct {
	RunID string `json:"-"`
}

// GetGemLocatorScanResponse returns one GEM locator scan.
type GetGemLocatorScanResponse struct {
	Run      *GemLocatorScan  `json:"run,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *GetGemLocatorScanResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// HideEventRequest hides an event.
type HideEventRequest struct {
	EventID string `json:"-"`
	Reason  string `json:"reason,omitempty"`
}

// HideEventResponse returns the hidden event.
type HideEventResponse struct {
	Event    *Event           `json:"event,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *HideEventResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// UnhideEventRequest unhides an event.
type UnhideEventRequest struct {
	EventID string `json:"-"`
}

// UnhideEventResponse returns the unhidden event.
type UnhideEventResponse struct {
	Event    *Event           `json:"event,omitempty"`
	Metadata ResponseMetadata `json:"-"`
}

func (r *UnhideEventResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

// ListEvents lists events with optional filters.
func (c *Client) ListEvents(ctx context.Context, request *ListEventsRequest, opts ...RequestOpt) (*ListEventsResponse, error) {
	if request == nil {
		request = &ListEventsRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/events", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryInt32(query, "pageSize", request.PageSize)
	setQueryString(query, "nextToken", request.NextToken)
	if source := strings.TrimSpace(string(request.Source)); source != "" && source != string(EventSourceUnspecified) {
		query.Set("source", source)
	}
	setQueryString(query, "query", request.Query)
	setQueryTime(query, "startsAfter", request.StartsAfter)
	setQueryTime(query, "startsBefore", request.StartsBefore)
	setQueryString(query, "eventType", request.EventType)
	setQueryString(query, "format", request.Format)
	setQueryString(query, "country", request.Country)
	setQueryString(query, "storeId", request.StoreID)
	setQueryString(query, "storeSlug", request.StoreSlug)
	setQueryFloat64(query, "latitude", request.Latitude)
	setQueryFloat64(query, "longitude", request.Longitude)
	setQueryFloat64(query, "radiusKm", request.RadiusKM)
	req.URL.RawQuery = query.Encode()

	response := &ListEventsResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetEvent fetches one event.
func (c *Client) GetEvent(ctx context.Context, request *GetEventRequest, opts ...RequestOpt) (*GetEventResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	eventID := strings.TrimSpace(request.EventID)
	if eventID == "" {
		return nil, errors.New("eventID must not be empty")
	}

	path := fmt.Sprintf("/v1/events/%s", url.PathEscape(eventID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetEventResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CreateCommunityEvent creates a community event.
func (c *Client) CreateCommunityEvent(ctx context.Context, request *CreateCommunityEventRequest, opts ...RequestOpt) (*CreateCommunityEventResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	if strings.TrimSpace(request.Title) == "" {
		return nil, errors.New("title must not be empty")
	}
	if request.StartsAt == nil {
		return nil, errors.New("startsAt must not be nil")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/v1/events", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CreateCommunityEventResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UpdateCommunityEvent updates a community event.
func (c *Client) UpdateCommunityEvent(ctx context.Context, request *UpdateCommunityEventRequest, opts ...RequestOpt) (*UpdateCommunityEventResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	eventID := strings.TrimSpace(request.EventID)
	if eventID == "" {
		return nil, errors.New("eventID must not be empty")
	}

	body, err := jsonBody(request)
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/events/%s", url.PathEscape(eventID))
	req, err := c.newRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UpdateCommunityEventResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// CancelCommunityEvent cancels a community event.
func (c *Client) CancelCommunityEvent(ctx context.Context, request *CancelCommunityEventRequest, opts ...RequestOpt) (*CancelCommunityEventResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	eventID := strings.TrimSpace(request.EventID)
	if eventID == "" {
		return nil, errors.New("eventID must not be empty")
	}

	body, err := jsonBody(struct{}{})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/events/%s:cancel", url.PathEscape(eventID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &CancelCommunityEventResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListStores lists event stores with optional filters.
func (c *Client) ListStores(ctx context.Context, request *ListStoresRequest, opts ...RequestOpt) (*ListStoresResponse, error) {
	if request == nil {
		request = &ListStoresRequest{}
	}

	req, err := c.newRequest(ctx, http.MethodGet, "/v1/event-stores", nil)
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	setQueryInt32(query, "pageSize", request.PageSize)
	setQueryString(query, "nextToken", request.NextToken)
	setQueryString(query, "query", request.Query)
	setQueryString(query, "country", request.Country)
	setQueryString(query, "armoryDay", request.ArmoryDay)
	setQueryBool(query, "onlineStore", request.OnlineStore)
	setQueryFloat64(query, "latitude", request.Latitude)
	setQueryFloat64(query, "longitude", request.Longitude)
	setQueryFloat64(query, "radiusKm", request.RadiusKM)
	req.URL.RawQuery = query.Encode()

	response := &ListStoresResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetStore fetches one event store.
func (c *Client) GetStore(ctx context.Context, request *GetStoreRequest, opts ...RequestOpt) (*GetStoreResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	storeRef := strings.TrimSpace(request.StoreRef)
	if storeRef == "" {
		return nil, errors.New("storeRef must not be empty")
	}

	path := fmt.Sprintf("/v1/event-stores/%s", url.PathEscape(storeRef))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetStoreResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ListEventFilters lists available event filter values.
func (c *Client) ListEventFilters(ctx context.Context, opts ...RequestOpt) (*ListEventFiltersResponse, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/v1/events/filters", nil)
	if err != nil {
		return nil, err
	}

	response := &ListEventFiltersResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// GetGemLocatorScan fetches a GEM locator scan.
func (c *Client) GetGemLocatorScan(ctx context.Context, request *GetGemLocatorScanRequest, opts ...RequestOpt) (*GetGemLocatorScanResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	runID := strings.TrimSpace(request.RunID)
	if runID == "" {
		return nil, errors.New("runID must not be empty")
	}

	path := fmt.Sprintf("/v1/admin/events/gem-locator-scans/%s", url.PathEscape(runID))
	req, err := c.newRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	response := &GetGemLocatorScanResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// HideEvent hides an event.
func (c *Client) HideEvent(ctx context.Context, request *HideEventRequest, opts ...RequestOpt) (*HideEventResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	eventID := strings.TrimSpace(request.EventID)
	if eventID == "" {
		return nil, errors.New("eventID must not be empty")
	}

	body, err := jsonBody(struct {
		Reason string `json:"reason,omitempty"`
	}{Reason: strings.TrimSpace(request.Reason)})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/admin/events/%s:hide", url.PathEscape(eventID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &HideEventResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// UnhideEvent unhides an event.
func (c *Client) UnhideEvent(ctx context.Context, request *UnhideEventRequest, opts ...RequestOpt) (*UnhideEventResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	eventID := strings.TrimSpace(request.EventID)
	if eventID == "" {
		return nil, errors.New("eventID must not be empty")
	}

	body, err := jsonBody(struct{}{})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}

	path := fmt.Sprintf("/v1/admin/events/%s:unhide", url.PathEscape(eventID))
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	response := &UnhideEventResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

func setQueryTime(values url.Values, key string, value *time.Time) {
	if value != nil {
		values.Set(key, value.UTC().Format(time.RFC3339))
	}
}

func setQueryFloat64(values url.Values, key string, value *float64) {
	if value != nil {
		values.Set(key, strconv.FormatFloat(*value, 'f', -1, 64))
	}
}
