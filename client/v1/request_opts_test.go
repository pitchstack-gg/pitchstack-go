package client

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWithRequestHeader(t *testing.T) {
	t.Run("when applied to request, then header is set", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
		err := WithRequestHeader("X-Test", "value")(req)
		require.NoError(t, err)
		require.Equal(t, "value", req.Header.Get("X-Test"))
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		err := WithRequestHeader("X-Test", "value")(nil)
		require.Error(t, err)
	})
}

func TestWithQueryParam(t *testing.T) {
	t.Run("when applied to request, then query param added", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "https://example.com", nil)
		err := WithQueryParam("foo", "bar")(req)
		require.NoError(t, err)

		values, err := url.ParseQuery(req.URL.RawQuery)
		require.NoError(t, err)
		require.Equal(t, "bar", values.Get("foo"))
	})

	t.Run("when request is nil, then returns error", func(t *testing.T) {
		err := WithQueryParam("foo", "bar")(nil)
		require.Error(t, err)
	})
}
