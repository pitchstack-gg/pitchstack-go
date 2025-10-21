package integration

import (
	"net/http"
	"testing"

	client "github.com/pitchstack-gg/pitchstack-go/client/client"
	"github.com/stretchr/testify/require"
)

func TestClientReturnsAPIErrorForMissingCard(t *testing.T) {
	ctx := contextWithTimeout(t)

	resp, err := integrationClient.GetCard(ctx, &client.GetCardRequest{
		CardID: cfg.NonexistentCardID,
	})
	require.Error(t, err, "expected error when fetching nonexistent card")
	require.Nil(t, resp, "response should be nil for errors")

	var apiErr *client.APIError
	require.ErrorAs(t, err, &apiErr, "expected *client.APIError")

	require.Equal(t, http.StatusNotFound, apiErr.Metadata.StatusCode, "expected 404 status")

	require.NotEmpty(t, apiErr.Status.Message, "expected API error message to be populated")
}
