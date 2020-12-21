package auth

import (
	"context"
)

type Authenticator interface {
	Authenticate(context.Context, Credential) (Credential, error)
}

type AuthenticatorFunc func(context.Context, Credential) (Credential, error)

func (a AuthenticatorFunc) Authenticate(ctx context.Context, credential Credential) (Credential, error) {
	return a(ctx, credential)
}
