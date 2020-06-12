package auth

type Authenticator interface {
	Authenticate(Credential) (Credential, error)
}

type AuthenticatorFunc func(Credential) (Credential, error)

func (a AuthenticatorFunc) Authenticate(credential Credential) (Credential, error) {
	return a(credential)
}
