package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestClientNews(t *testing.T) {
	t.Run("list recommended articles", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/news/recommended", r.URL.Path)
			q := r.URL.Query()
			require.Equal(t, "20", q.Get("pageSize"))
			require.Equal(t, "tok-1", q.Get("nextToken"))
			require.Equal(t, "en-US", q.Get("locale"))
			require.Equal(t, "ios", q.Get("platform"))
			require.NoError(t, json.NewEncoder(w).Encode(ListRecommendedArticlesResponse{
				Items: []RecommendedArticle{{Article: &NewsArticle{ArticleID: "a-1"}}},
			}))
		}))
		t.Cleanup(server.Close)

		pageSize := int32(20)
		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.ListRecommendedArticles(context.Background(), &ListRecommendedArticlesRequest{
			PageSize:  &pageSize,
			NextToken: "tok-1",
			Locale:    "en-US",
			Platform:  "ios",
		})
		require.NoError(t, err)
		require.Len(t, resp.Items, 1)
	})

	t.Run("get article", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, http.MethodGet, r.Method)
			require.Equal(t, "/v1/news/a-1", r.URL.Path)
			require.NoError(t, json.NewEncoder(w).Encode(GetArticleResponse{Article: &NewsArticle{ArticleID: "a-1"}}))
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		resp, err := client.GetArticle(context.Background(), &GetArticleRequest{ArticleID: "a-1"})
		require.NoError(t, err)
		require.Equal(t, "a-1", resp.Article.ArticleID)

		_, err = client.GetArticle(context.Background(), &GetArticleRequest{})
		require.Error(t, err)
	})

	t.Run("track impression and click", func(t *testing.T) {
		var impressions, clicks int
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/v1/news/a-1:trackImpression":
				require.Equal(t, http.MethodPost, r.Method)
				impressions++
				var payload TrackArticleImpressionRequest
				require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
				require.Equal(t, "a-1", payload.ArticleID)
				require.NoError(t, json.NewEncoder(w).Encode(TrackArticleImpressionResponse{Accepted: true}))
			case "/v1/news/a-1:trackClick":
				require.Equal(t, http.MethodPost, r.Method)
				clicks++
				var payload TrackArticleClickRequest
				require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
				require.Equal(t, "a-1", payload.ArticleID)
				require.Equal(t, "https://example.com", payload.DestinationURL)
				require.NoError(t, json.NewEncoder(w).Encode(TrackArticleClickResponse{Accepted: true}))
			default:
				t.Fatalf("unexpected path %s", r.URL.Path)
			}
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		_, err := client.TrackArticleImpression(context.Background(), &TrackArticleImpressionRequest{ArticleID: "a-1"})
		require.NoError(t, err)
		_, err = client.TrackArticleClick(context.Background(), &TrackArticleClickRequest{ArticleID: "a-1", DestinationURL: "https://example.com"})
		require.NoError(t, err)
		require.Equal(t, 1, impressions)
		require.Equal(t, 1, clicks)
	})

	t.Run("admin source controls", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/v1/admin/news/sources/s-1:disable":
				require.Equal(t, http.MethodPost, r.Method)
				require.NoError(t, json.NewEncoder(w).Encode(DisableSourceResponse{Source: &NewsSource{SourceID: "s-1"}}))
			case "/v1/admin/news/sources/s-1:enable":
				require.Equal(t, http.MethodPost, r.Method)
				require.NoError(t, json.NewEncoder(w).Encode(EnableSourceResponse{Source: &NewsSource{SourceID: "s-1"}}))
			case "/v1/admin/news/sources/s-1:ingestNow":
				require.Equal(t, http.MethodPost, r.Method)
				var payload RunSourceIngestionRequest
				require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
				require.Equal(t, "s-1", payload.SourceID)
				require.NoError(t, json.NewEncoder(w).Encode(RunSourceIngestionResponse{SourceID: "s-1", Queued: true}))
			default:
				t.Fatalf("unexpected path %s", r.URL.Path)
			}
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		_, err := client.DisableSource(context.Background(), &DisableSourceRequest{SourceID: "s-1"})
		require.NoError(t, err)
		_, err = client.EnableSource(context.Background(), &EnableSourceRequest{SourceID: "s-1"})
		require.NoError(t, err)
		_, err = client.RunSourceIngestion(context.Background(), &RunSourceIngestionRequest{SourceID: "s-1"})
		require.NoError(t, err)
	})

	t.Run("manual article lifecycle methods", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/v1/admin/news/articles":
				require.Equal(t, http.MethodPost, r.Method)
				var payload CreateManualArticleRequest
				require.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
				require.Equal(t, "Hello", payload.Title)
				require.NoError(t, json.NewEncoder(w).Encode(CreateManualArticleResponse{Article: &NewsArticle{ArticleID: "a-1"}}))
			case "/v1/admin/news/articles/a-1":
				require.Equal(t, http.MethodPatch, r.Method)
				require.NoError(t, json.NewEncoder(w).Encode(UpdateManualArticleResponse{Article: &NewsArticle{ArticleID: "a-1"}}))
			case "/v1/admin/news/articles/a-1:publish":
				require.Equal(t, http.MethodPost, r.Method)
				require.NoError(t, json.NewEncoder(w).Encode(PublishManualArticleResponse{Article: &NewsArticle{ArticleID: "a-1"}}))
			case "/v1/admin/news/articles/a-1:unpublish":
				require.Equal(t, http.MethodPost, r.Method)
				require.NoError(t, json.NewEncoder(w).Encode(UnpublishManualArticleResponse{Article: &NewsArticle{ArticleID: "a-1"}}))
			case "/v1/admin/news/articles/a-1:setOverride":
				require.Equal(t, http.MethodPost, r.Method)
				require.NoError(t, json.NewEncoder(w).Encode(SetArticleOverrideResponse{Article: &NewsArticle{ArticleID: "a-1"}}))
			default:
				t.Fatalf("unexpected path %s", r.URL.Path)
			}
		}))
		t.Cleanup(server.Close)

		client := newTestClient(t, WithBaseURL(server.URL), WithHTTPClient(server.Client()))
		_, err := client.CreateManualArticle(context.Background(), &CreateManualArticleRequest{Title: "Hello"})
		require.NoError(t, err)

		title := "Updated"
		_, err = client.UpdateManualArticle(context.Background(), &UpdateManualArticleRequest{ArticleID: "a-1", Title: &title})
		require.NoError(t, err)
		_, err = client.PublishManualArticle(context.Background(), &PublishManualArticleRequest{ArticleID: "a-1"})
		require.NoError(t, err)
		_, err = client.UnpublishManualArticle(context.Background(), &UnpublishManualArticleRequest{ArticleID: "a-1"})
		require.NoError(t, err)
		_, err = client.SetArticleOverride(context.Background(), &SetArticleOverrideRequest{ArticleID: "a-1"})
		require.NoError(t, err)
	})
}
