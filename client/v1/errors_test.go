package client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAPIError(t *testing.T) {
	t.Run("when body contains rpc status, then fields decoded", func(t *testing.T) {
		meta := ResponseMetadata{StatusCode: 400}
		err := newAPIError(meta, []byte(`{"code":3,"message":"invalid"}`))

		apiErr, ok := err.(*APIError)
		require.True(t, ok)
		require.Equal(t, meta, apiErr.Metadata)
		require.Equal(t, 3, apiErr.Status.Code)
		require.Equal(t, "invalid", apiErr.Status.Message)
		require.Equal(t, []byte(`{"code":3,"message":"invalid"}`), apiErr.RawBody)
	})

	t.Run("when body is not json, then message falls back to raw string", func(t *testing.T) {
		meta := ResponseMetadata{StatusCode: 500}
		err := newAPIError(meta, []byte("internal error"))

		apiErr := err.(*APIError)
		require.Equal(t, "internal error", apiErr.Status.Message)
		require.Equal(t, 0, apiErr.Status.Code)
	})
}

func TestAPIErrorError(t *testing.T) {
	t.Run("when formatted, then includes status code and message", func(t *testing.T) {
		err := &APIError{
			Metadata: ResponseMetadata{StatusCode: 403},
			Status:   RPCStatus{Code: 16, Message: "permission denied"},
		}

		msg := err.Error()
		require.Contains(t, msg, "status 403")
		require.Contains(t, msg, `code=16 message="permission denied"`)
	})
}
