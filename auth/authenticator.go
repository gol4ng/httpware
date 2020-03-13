package auth

type Authenticator interface {
	Authenticate(Credential) (Credential, error)
}
