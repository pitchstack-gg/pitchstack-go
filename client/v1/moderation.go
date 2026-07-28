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

// ModerationReportReason is shared by profile, deck, and collection reports.
type ModerationReportReason string

const (
	ModerationReportReasonUnspecified      ModerationReportReason = "MODERATION_REPORT_REASON_UNSPECIFIED"
	ModerationReportReasonSpamOrScam       ModerationReportReason = "MODERATION_REPORT_REASON_SPAM_OR_SCAM"
	ModerationReportReasonImpersonation    ModerationReportReason = "MODERATION_REPORT_REASON_IMPERSONATION"
	ModerationReportReasonHarassmentOrHate ModerationReportReason = "MODERATION_REPORT_REASON_HARASSMENT_OR_HATE"
	ModerationReportReasonSexualContent    ModerationReportReason = "MODERATION_REPORT_REASON_SEXUAL_CONTENT"
	ModerationReportReasonViolence         ModerationReportReason = "MODERATION_REPORT_REASON_VIOLENCE"
	ModerationReportReasonOther            ModerationReportReason = "MODERATION_REPORT_REASON_OTHER"
)

// DeckReportTargetType identifies the deck content being reported.
type DeckReportTargetType string

const (
	DeckReportTargetTypeUnspecified    DeckReportTargetType = "DECK_REPORT_TARGET_TYPE_UNSPECIFIED"
	DeckReportTargetTypeDeck           DeckReportTargetType = "DECK_REPORT_TARGET_TYPE_DECK"
	DeckReportTargetTypeVersionNotes   DeckReportTargetType = "DECK_REPORT_TARGET_TYPE_VERSION_NOTES"
	DeckReportTargetTypeSideboardGuide DeckReportTargetType = "DECK_REPORT_TARGET_TYPE_SIDEBOARD_GUIDE"
	DeckReportTargetTypeMatch          DeckReportTargetType = "DECK_REPORT_TARGET_TYPE_MATCH"
)

// ModerationReportReceipt is the intentionally neutral response to a reporter.
type ModerationReportReceipt struct {
	CaseID      string     `json:"caseId,omitempty"`
	SubmittedAt *time.Time `json:"submittedAt,omitempty"`
	Duplicate   bool       `json:"duplicate,omitempty"`
}

type SubmitModerationReportResponse struct {
	Receipt  *ModerationReportReceipt `json:"receipt,omitempty"`
	Metadata ResponseMetadata         `json:"-"`
}

func (r *SubmitModerationReportResponse) setMetadata(metadata ResponseMetadata) {
	r.Metadata = metadata
}

type ReportProfileRequest struct {
	UserID  string
	Reason  ModerationReportReason
	Details string
}

type ReportDeckContentRequest struct {
	DeckID     string
	TargetType DeckReportTargetType
	TargetID   string
	Reason     ModerationReportReason
	Details    string
}

type ReportCollectionRequest struct {
	CollectionID string
	Reason       ModerationReportReason
	Details      string
}

func (c *Client) ReportProfile(ctx context.Context, request *ReportProfileRequest, opts ...RequestOpt) (*SubmitModerationReportResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	userID, err := required(request.UserID, "userID")
	if err != nil {
		return nil, err
	}
	return c.submitModerationReport(ctx, "/v1/users/"+url.PathEscape(userID)+"/profile:report", request.Reason, request.Details, nil, "", opts...)
}

func (c *Client) ReportDeckContent(ctx context.Context, request *ReportDeckContentRequest, opts ...RequestOpt) (*SubmitModerationReportResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	deckID, err := required(request.DeckID, "deckID")
	if err != nil {
		return nil, err
	}
	if request.TargetType == "" || request.TargetType == DeckReportTargetTypeUnspecified {
		return nil, errors.New("targetType must be specified")
	}
	if request.TargetType != DeckReportTargetTypeDeck && strings.TrimSpace(request.TargetID) == "" {
		return nil, errors.New("targetID is required for deck subresources")
	}
	return c.submitModerationReport(ctx, "/v1/decks/"+url.PathEscape(deckID)+":report", request.Reason, request.Details, &request.TargetType, request.TargetID, opts...)
}

func (c *Client) ReportCollection(ctx context.Context, request *ReportCollectionRequest, opts ...RequestOpt) (*SubmitModerationReportResponse, error) {
	if request == nil {
		return nil, errors.New("request must not be nil")
	}
	collectionID, err := required(request.CollectionID, "collectionID")
	if err != nil {
		return nil, err
	}
	return c.submitModerationReport(ctx, "/v1/collections/"+url.PathEscape(collectionID)+":report", request.Reason, request.Details, nil, "", opts...)
}

func (c *Client) submitModerationReport(ctx context.Context, path string, reason ModerationReportReason, details string, targetType *DeckReportTargetType, targetID string, opts ...RequestOpt) (*SubmitModerationReportResponse, error) {
	if reason == "" || reason == ModerationReportReasonUnspecified {
		return nil, errors.New("reason must be specified")
	}
	body, err := jsonBody(struct {
		Reason     ModerationReportReason `json:"reason"`
		Details    string                 `json:"details,omitempty"`
		TargetType *DeckReportTargetType  `json:"targetType,omitempty"`
		TargetID   string                 `json:"targetId,omitempty"`
	}{Reason: reason, Details: strings.TrimSpace(details), TargetType: targetType, TargetID: strings.TrimSpace(targetID)})
	if err != nil {
		return nil, fmt.Errorf("encode body: %w", err)
	}
	req, err := c.newRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	response := &SubmitModerationReportResponse{}
	if err := c.do(req, response, opts...); err != nil {
		return nil, err
	}
	return response, nil
}

// ReportGroupWithSharedReason adapts the shared reason model to the additive
// legacy group route retained during the group-report migration.
func (c *Client) ReportGroupWithSharedReason(ctx context.Context, groupID string, reason ModerationReportReason, details string, opts ...RequestOpt) (*GroupReportResponse, error) {
	mapped, ok := sharedGroupReportReason(reason)
	if !ok {
		return nil, errors.New("reason must be specified")
	}
	return c.ReportGroup(ctx, groupID, mapped, details, opts...)
}

func sharedGroupReportReason(reason ModerationReportReason) (GroupReportReason, bool) {
	switch reason {
	case ModerationReportReasonSpamOrScam:
		return GroupReportReasonSpamOrScam, true
	case ModerationReportReasonImpersonation:
		return GroupReportReasonImpersonation, true
	case ModerationReportReasonHarassmentOrHate:
		return GroupReportReasonHarassmentOrHate, true
	case ModerationReportReasonSexualContent:
		return GroupReportReasonSexualContent, true
	case ModerationReportReasonViolence:
		return GroupReportReasonViolence, true
	case ModerationReportReasonOther:
		return GroupReportReasonOther, true
	default:
		return "", false
	}
}
