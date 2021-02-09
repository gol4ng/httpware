package tripperware

import (
	"net/http"

	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/auth"
)

func AuthenticationForwarder(options ...AuthOption) httpware.Tripperware {
	config := NewAuthConfig(options...)
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			config.credentialForwarder(req)
			return next.RoundTrip(req)
		})
	}
}

type credentialForwarder func(req *http.Request)

// AuthOption defines a interceptor tripperware configuration option
type AuthOption func(*AuthConfig)

type AuthConfig struct {
	credentialForwarder credentialForwarder
}

func (o *AuthConfig) apply(options ...AuthOption) {
	for _, option := range options {
		option(o)
	}
}

func NewAuthConfig(options ...AuthOption) *AuthConfig {
	opts := &AuthConfig{
		credentialForwarder: DefaultCredentialForwarder,
	}
	opts.apply(options...)
	return opts
}

func DefaultCredentialForwarder(req *http.Request) {
	auth.AddHeader(req)(auth.CredentialFromContext(req.Context()))
}

// WithCredentialForwarder will configure credentialForwarder option
func WithCredentialForwarder(authFunc credentialForwarder) AuthOption {
	return func(config *AuthConfig) {
		config.credentialForwarder = authFunc
	}
}
