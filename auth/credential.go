package auth

type Credential string

type CredentialProvider func() Credential

type CredentialSetter func(Credential)
