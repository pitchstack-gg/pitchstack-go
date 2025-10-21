package client

import "net/http"

// ResponseMetadata surfaces useful response attributes for debugging.
type ResponseMetadata struct {
	StatusCode int
	Header     http.Header
	RequestID  string
}

type metadataSetter interface {
	setMetadata(ResponseMetadata)
}
