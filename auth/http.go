package auth

import (
	"net/http"
)

const (
	AuthorizationHeader  = "Authorization"
	XAuthorizationHeader = "X-Authorization"
)

func FromHeader(request *http.Request) CredentialProvider {
	return func() Credential {
		return ExtractFromHeader(request)
	}
}

func ExtractFromHeader(request *http.Request) Credential {
	if request == nil {
		return ""
	}

	tokenHeader := request.Header.Get(AuthorizationHeader)
	if tokenHeader == "" {
		tokenHeader = request.Header.Get(XAuthorizationHeader)
	}

	return tokenHeader
}

func AddHeader(request *http.Request) CredentialSetter {
	return func(credential Credential) {
		if request == nil {
			return
		}
		if creds, ok := credential.(string); ok {
			request.Header.Set(AuthorizationHeader, creds)
			request.Header.Set(XAuthorizationHeader, creds)
		}
	}
}
