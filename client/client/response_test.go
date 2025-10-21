package client

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type metadataCapture struct {
	Metadata ResponseMetadata
}

func (m *metadataCapture) setMetadata(metadata ResponseMetadata) {
	m.Metadata = metadata
}

var _ metadataSetter = (*metadataCapture)(nil)

func TestResponseMetadata(t *testing.T) {
	t.Run("when metadata applied, then fields are stored for debugging", func(t *testing.T) {
		capture := &metadataCapture{}
		meta := ResponseMetadata{
			StatusCode: http.StatusCreated,
			Header:     http.Header{"X-Trace": []string{"trace"}},
			RequestID:  "req-123",
		}

		capture.setMetadata(meta)
		require.Equal(t, http.StatusCreated, capture.Metadata.StatusCode)
		require.Equal(t, "trace", capture.Metadata.Header.Get("X-Trace"))
		require.Equal(t, "req-123", capture.Metadata.RequestID)
	})
}
