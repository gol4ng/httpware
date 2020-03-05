package middleware

import (
	"context"
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/auth"
)

func AuthenticationMiddleware(options ...AuthOption) httpware.Middleware {
	config := NewAuthConfig(options...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			newCtx, err := config.authFunc(req)
			if err != nil && config.onError(err, writer, req) {
				return
			}

			next.ServeHTTP(writer, req.WithContext(newCtx))
		})
	}
}

type AuthFunc func(req *http.Request) (context.Context, error)
type errorHandler func(err error, writer http.ResponseWriter, req *http.Request) bool

// AuthOption defines a interceptor middleware configuration option
type AuthOption func(*AuthConfig)

type AuthConfig struct {
	authFunc AuthFunc
	onError  errorHandler
}

func (o *AuthConfig) apply(options ...AuthOption) {
	for _, option := range options {
		option(o)
	}
}

func NewAuthConfig(options ...AuthOption) *AuthConfig {
	opts := &AuthConfig{
		authFunc: DefaultAuthFunc,
		onError:  DefaultErrorHandler,
	}
	opts.apply(options...)
	return opts
}

func DefaultAuthFunc(req *http.Request) (context.Context, error) {
	return auth.CredentialToContext(req.Context(), auth.FromHeader(req)()), nil
}

func DefaultErrorHandler(err error, writer http.ResponseWriter, _ *http.Request) bool {
	http.Error(writer, err.Error(), http.StatusUnauthorized)
	return true
}
