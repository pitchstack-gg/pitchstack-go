package integration

import (
	"testing"

	client "github.com/pitchstack-gg/pitchstack-go/client/v1"
	"github.com/stretchr/testify/require"
)

func TestSearchCardsReturnsSummaries(t *testing.T) {
	term := requireString(t, cfg.CardSearchTerm, envCardSearchTerm)

	ctx := contextWithTimeout(t)
	resp, err := integrationClient.SearchCards(ctx, &client.SearchCardsRequest{
		SearchTerm: term,
		PageSize:   int32Ptr(3),
	})
	require.NoError(t, err, "SearchCards")

	ensureMetadata(t, resp.Metadata)
	require.NotEmpty(t, resp.Summaries, "expected at least one card summary")

	if resp.NextToken != "" {
		nextResp, err := integrationClient.SearchCards(ctx, &client.SearchCardsRequest{
			SearchTerm: term,
			PageSize:   int32Ptr(3),
			NextToken:  resp.NextToken,
		})
		require.NoError(t, err, "SearchCards second page")
		require.NotEmpty(t, nextResp.Summaries, "expected second page to include summaries")
	}
}

func TestGetCardReturnsSummary(t *testing.T) {
	cardID := requireString(t, cfg.CardID, envCardID)

	ctx := contextWithTimeout(t)
	resp, err := integrationClient.GetCard(ctx, &client.GetCardRequest{CardID: cardID})
	require.NoError(t, err, "GetCard")
	require.NotNil(t, resp.Summary, "expected summary in response")

	ensureMetadata(t, resp.Metadata)

	require.Equal(t, cardID, resp.Summary.Identifier, "expected identifier to match")
}

func TestGetCardPrintingReturnsSummary(t *testing.T) {
	printingID := requireString(t, cfg.PrintingID, envPrintingID)

	ctx := contextWithTimeout(t)
	resp, err := integrationClient.GetCardPrinting(ctx, &client.GetCardPrintingRequest{PrintingID: printingID})
	require.NoError(t, err, "GetCardPrinting")
	require.NotNil(t, resp.Summary, "expected summary in response")

	ensureMetadata(t, resp.Metadata)

	require.Equal(t, printingID, resp.Summary.Identifier, "expected identifier to match")
}

func int32Ptr(v int32) *int32 {
	return &v
}
