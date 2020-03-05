package middleware

import (
	"context"
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/auth"
)

// Authentication middleware delegate the authentication process to a authFunc configured
func Authentication(options ...AuthOption) httpware.Middleware {
	config := NewAuthConfig(options...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			newCtx, err := config.authFunc(req)
			if err != nil && config.errorHandler(err, writer, req) {
				return
			}

			next.ServeHTTP(writer, req.WithContext(newCtx))
		})
	}
}

type authFunc func(req *http.Request) (context.Context, error)
type errorHandler func(err error, writer http.ResponseWriter, req *http.Request) bool

// AuthOption defines a interceptor middleware configuration option
type AuthOption func(*AuthConfig)

type AuthConfig struct {
	authFunc     authFunc
	errorHandler errorHandler
}

func (o *AuthConfig) apply(options ...AuthOption) {
	for _, option := range options {
		option(o)
	}
}

func NewAuthConfig(options ...AuthOption) *AuthConfig {
	opts := &AuthConfig{
		authFunc:     DefaultAuthFunc,
		errorHandler: DefaultErrorHandler,
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

// WithAuthFunc will configure authFunc option
func WithAuthFunc(authFunc authFunc) AuthOption {
	return func(config *AuthConfig) {
		config.authFunc = authFunc
	}
}

// WithErrorHandler will configure errorHandler option
func WithErrorHandler(errorHandler errorHandler) AuthOption {
	return func(config *AuthConfig) {
		config.errorHandler = errorHandler
	}
}
