package integration

import (
	"context"
	"log"
	"os"
	"strings"
	"testing"
	"time"

	client "github.com/pitchstack-gg/pitchstack-go/client/v1"
	"github.com/stretchr/testify/require"
)

const (
	placeholderPrefix = "REPLACE_ME"

	envAPIBaseURL        = "PITCHSTACK_API_BASE_URL"
	envCardID            = "PITCHSTACK_CARD_ID"
	envCardSearchTerm    = "PITCHSTACK_CARD_SEARCH_TERM"
	envPrintingID        = "PITCHSTACK_PRINTING_ID"
	envProductID         = "PITCHSTACK_PRODUCT_ID"
	envCollectionItemID  = "PITCHSTACK_COLLECTION_ITEM_ID"
	envNonexistentCardID = "PITCHSTACK_NONEXISTENT_CARD_ID"

	defaultAPIBaseURL = "https://gamma-api.pitchstack.gg"
)

type testConfig struct {
	APIBaseURL        string
	APIToken          string
	UserID            string
	CollectionID      string
	DeckID            string
	DeckVersionID     string
	CardID            string
	CardSearchTerm    string
	PrintingID        string
	ProductID         string
	CollectionItemID  string
	TagResourceID     string
	TagResourceType   client.ResourceType
	TagKey            string
	TagValue          string
	ExpectedUsername  string
	PricingSource     string
	PricingCurrency   string
	NonexistentCardID string
}

var (
	integrationClient *client.Client
	cfg               testConfig
)

func TestMain(m *testing.M) {
	cfg = loadTestConfig()

	opts := []client.ClientOpt{
		client.WithBaseURL(cfg.APIBaseURL),
		client.WithStaticBearerToken(cfg.APIToken),
		client.WithHTTPTimeout(20 * time.Second),
	}

	c, err := client.NewClient(opts...)
	if err != nil {
		log.Fatalf("create client: %v", err)
	}
	integrationClient = c

	os.Exit(m.Run())
}

func loadTestConfig() testConfig {
	return testConfig{
		APIBaseURL:        getEnvOrDefault(envAPIBaseURL, defaultAPIBaseURL),
		CardID:            getEnvOrDefault(envCardID, "sink-below-r"),
		CardSearchTerm:    getEnvOrDefault(envCardSearchTerm, "sink"),
		PrintingID:        getEnvOrDefault(envPrintingID, "WTR215"),
		ProductID:         getEnvOrDefault(envProductID, "WTR215"),
		NonexistentCardID: getEnvOrDefault(envNonexistentCardID, "hahalol"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value != "" {
		return value
	}
	return defaultValue
}

func requireString(t testing.TB, value, name string) string {
	t.Helper()
	if isPlaceholder(value) {
		t.Skipf("%s not configured; set %s", name, name)
	}
	return value
}

func contextWithTimeout(t testing.TB) context.Context {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	t.Cleanup(cancel)
	return ctx
}

func ensureMetadata(t testing.TB, metadata client.ResponseMetadata) {
	t.Helper()
	require.NotZero(t, metadata.StatusCode, "expected response metadata status code to be populated")
	if metadata.RequestID == "" {
		t.Log("response metadata missing request id; continuing")
	}
}

func isPlaceholder(value string) bool {
	return value == "" || strings.HasPrefix(value, placeholderPrefix)
}
