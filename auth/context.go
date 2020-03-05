package auth

import (
	"context"
)

var credentialContextKey struct{}

func CredentialToContext(ctx context.Context, credential Credential) context.Context {
	return context.WithValue(ctx, credentialContextKey, credential)
}

func CredentialFromContext(ctx context.Context) Credential {
	if ctx == nil {
		return ""
	}
	value := ctx.Value(credentialContextKey)
	if value == nil {
		return ""
	}
	credential, ok := value.(Credential)
	if !ok {
		return ""
	}

	return credential
}
