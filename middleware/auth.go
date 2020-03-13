package middleware

import (
	"context"
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/auth"
)

// Authentication middleware delegate the authentication process to a AuthenticateFunc configured
func Authentication(authenticator auth.Authenticator, options ...AuthOption) httpware.Middleware {
	config := NewAuthConfig(options...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			newCtx, err := config.authenticateFunc(config.credentialFinder, authenticator, req)
			if err != nil && config.errorHandler(err, writer, req) {
				return
			}

			next.ServeHTTP(writer, req.WithContext(newCtx))
		})
	}
}

type CredentialFinder func(r *http.Request) auth.Credential
type AuthenticateFunc func(credentialFinder CredentialFinder, authenticator auth.Authenticator, req *http.Request) (context.Context, error)
type ErrorHandler func(err error, writer http.ResponseWriter, req *http.Request) bool

// AuthOption defines a interceptor middleware configuration option
type AuthOption func(*AuthConfig)

type AuthConfig struct {
	credentialFinder CredentialFinder
	authenticateFunc AuthenticateFunc
	errorHandler     ErrorHandler
}

func (o *AuthConfig) apply(options ...AuthOption) {
	for _, option := range options {
		option(o)
	}
}

func NewAuthConfig(options ...AuthOption) *AuthConfig {
	opts := &AuthConfig{
		credentialFinder: DefaultCredentialFinder,
		authenticateFunc: DefaultAuthFunc,
		errorHandler:     DefaultErrorHandler,
	}
	opts.apply(options...)
	return opts
}

func DefaultCredentialFinder(request *http.Request) auth.Credential {
	return auth.FromHeader(request)()
}

func DefaultAuthFunc(credentialFinder CredentialFinder, authenticator auth.Authenticator, request *http.Request) (context.Context, error) {
	credential := credentialFinder(request)
	if authenticator != nil {
		creds, err := authenticator.Authenticate(credential)
		if err != nil {
			return request.Context(), err
		}
		credential = creds
	}
	return auth.CredentialToContext(request.Context(), credential), nil
}

func DefaultErrorHandler(err error, writer http.ResponseWriter, _ *http.Request) bool {
	http.Error(writer, err.Error(), http.StatusUnauthorized)
	return true
}

// WithCredentialFinder will configure AuthenticateFunc option
func WithCredentialFinder(credentialFinder CredentialFinder) AuthOption {
	return func(config *AuthConfig) {
		config.credentialFinder = credentialFinder
	}
}

// WithAuthenticateFunc will configure AuthenticateFunc option
func WithAuthenticateFunc(authenticateFunc AuthenticateFunc) AuthOption {
	return func(config *AuthConfig) {
		config.authenticateFunc = authenticateFunc
	}
}

// WithErrorHandler will configure ErrorHandler option
func WithErrorHandler(errorHandler ErrorHandler) AuthOption {
	return func(config *AuthConfig) {
		config.errorHandler = errorHandler
	}
}
