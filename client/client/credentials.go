package client

import "context"

// Credential represents an HTTP header/value pair used for authentication.
type Credential struct {
	Header string
	Value  string
}

// CredentialProvider returns credentials for a request.
type CredentialProvider interface {
	Credential(ctx context.Context) (*Credential, error)
}

// CredentialProviderFunc adapts a function into a CredentialProvider.
type CredentialProviderFunc func(ctx context.Context) (*Credential, error)

// CredentialProviderFunc implements the CredentialProvider interface.
func (f CredentialProviderFunc) Credential(ctx context.Context) (*Credential, error) {
	return f(ctx)
}
