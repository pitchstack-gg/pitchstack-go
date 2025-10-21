package client

import (
	"encoding/json"
	"fmt"
	"strings"
)

// RPCStatus models the default error payload defined by the API.
type RPCStatus struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []interface{} `json:"details,omitempty"`
}

// APIError represents a non-2xx response from the Pitchstack API.
type APIError struct {
	Metadata ResponseMetadata
	Status   RPCStatus
	RawBody  []byte
}

// Error formats the API error message.
func (e *APIError) Error() string {
	if e == nil {
		return "<nil>"
	}

	var builder strings.Builder
	builder.WriteString("pitchstack api error")

	if e.Metadata.StatusCode != 0 {
		builder.WriteString(fmt.Sprintf(" (status %d)", e.Metadata.StatusCode))
	}
	if e.Status.Code != 0 || e.Status.Message != "" {
		builder.WriteString(fmt.Sprintf(": code=%d message=%q", e.Status.Code, e.Status.Message))
	}

	return builder.String()
}

func newAPIError(metadata ResponseMetadata, body []byte) error {
	var status RPCStatus
	if len(body) > 0 {
		if err := json.Unmarshal(body, &status); err != nil {
			// Fallback to treating the body as a plain string message.
			status.Message = strings.TrimSpace(string(body))
		}
	}

	return &APIError{
		Metadata: metadata,
		Status:   status,
		RawBody:  body,
	}
}
