package auth

type Credential interface{}

type CredentialProvider func() Credential

type CredentialSetter func(Credential)
