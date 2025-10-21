package client

import (
	"fmt"
	"net/http"
	"net/url"
)

// RequestOpt customizes an outbound HTTP request prior to transmission.
type RequestOpt func(*http.Request) error

// WithRequestHeader sets or replaces a header on the outbound request.
func WithRequestHeader(key, value string) RequestOpt {
	return func(req *http.Request) error {
		if req == nil {
			return fmt.Errorf("request is nil")
		}
		req.Header.Set(key, value)
		return nil
	}
}

// WithQueryParam adds or replaces a query string value on the request URL.
func WithQueryParam(key, value string) RequestOpt {
	return func(req *http.Request) error {
		if req == nil {
			return fmt.Errorf("request is nil")
		}
		query, err := url.ParseQuery(req.URL.RawQuery)
		if err != nil {
			return fmt.Errorf("parse query string: %w", err)
		}
		query.Set(key, value)
		req.URL.RawQuery = query.Encode()
		return nil
	}
}
