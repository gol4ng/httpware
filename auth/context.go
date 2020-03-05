package auth

import (
	"context"
)

type credentialContextKeyType struct{}

var credentialContextKey = credentialContextKeyType{}

func CredentialToContext(ctx context.Context, credential Credential) context.Context {
	return context.WithValue(ctx, credentialContextKey, credential)
}

func CredentialFromContext(ctx context.Context) Credential {
	credential, ok := ctx.Value(credentialContextKey).(Credential)
	if !ok {
		return ""
	}

	return credential
}
