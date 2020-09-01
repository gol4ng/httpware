package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware/v3"
	"github.com/gol4ng/httpware/v3/auth"
)

// Authentication middleware delegate the authentication process to the AuthenticateFunc
func Authentication(authenticateFunc AuthenticateFunc, options ...AuthOption) httpware.Middleware {
	config := newAuthConfig(options...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			newReq, err := authenticateFunc(req)
			if err == nil {
				config.successMiddleware(next).ServeHTTP(writer, newReq)
				return
			}
			if config.errorHandler(err, writer, req) {
				return
			}

			next.ServeHTTP(writer, newReq)
		})
	}
}

type CredentialFinder func(r *http.Request) auth.Credential
type AuthenticateFunc func(req *http.Request) (*http.Request, error)
type ErrorHandler func(err error, writer http.ResponseWriter, req *http.Request) bool

// AuthOption defines a interceptor middleware configuration option
type AuthOption func(*AuthConfig)

type AuthConfig struct {
	errorHandler      ErrorHandler
	successMiddleware httpware.Middleware
}

func (o *AuthConfig) apply(options ...AuthOption) {
	for _, option := range options {
		option(o)
	}
}

func newAuthConfig(options ...AuthOption) *AuthConfig {
	opts := &AuthConfig{
		errorHandler:      DefaultErrorHandler,
		successMiddleware: httpware.NopMiddleware,
	}
	opts.apply(options...)
	return opts
}

func DefaultErrorHandler(err error, writer http.ResponseWriter, _ *http.Request) bool {
	http.Error(writer, err.Error(), http.StatusUnauthorized)
	return true
}

// WithErrorHandler will configure ErrorHandler option
func WithErrorHandler(errorHandler ErrorHandler) AuthOption {
	return func(config *AuthConfig) {
		config.errorHandler = errorHandler
	}
}

// WithSuccessMiddleware will configure successMiddleware option
func WithSuccessMiddleware(middleware httpware.Middleware) AuthOption {
	return func(config *AuthConfig) {
		config.successMiddleware = middleware
	}
}

// NewAuthenticateFunc is an AuthenticateFunc that find, authenticate and hydrate credentials on the request context
func NewAuthenticateFunc(authenticator auth.Authenticator, options ...AuthFuncOption) AuthenticateFunc {
	config := newAuthFuncConfig(options...)
	return func(request *http.Request) (*http.Request, error) {
		credential := config.credentialFinder(request)
		if authenticator != nil {
			creds, err := authenticator.Authenticate(credential)
			if err != nil {
				return request, err
			}
			credential = creds
		}
		return request.WithContext(auth.CredentialToContext(request.Context(), credential)), nil
	}
}

// AuthFuncOption defines a AuthenticateFunc configuration option
type AuthFuncOption func(*AuthFuncConfig)

type AuthFuncConfig struct {
	credentialFinder CredentialFinder
}

func (o *AuthFuncConfig) apply(options ...AuthFuncOption) {
	for _, option := range options {
		option(o)
	}
}

func newAuthFuncConfig(options ...AuthFuncOption) *AuthFuncConfig {
	opts := &AuthFuncConfig{
		credentialFinder: DefaultCredentialFinder,
	}
	opts.apply(options...)
	return opts
}

func DefaultCredentialFinder(request *http.Request) auth.Credential {
	return auth.FromHeader(request)()
}

// WithCredentialFinder will configure AuthenticateFunc option
func WithCredentialFinder(credentialFinder CredentialFinder) AuthFuncOption {
	return func(config *AuthFuncConfig) {
		config.credentialFinder = credentialFinder
	}
}
