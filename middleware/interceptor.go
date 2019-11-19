package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/interceptor"
)

// Interceptor middleware allow multiple req.Body read and allow to set callback before and after roundtrip
func Interceptor(options ...Option) httpware.Middleware {
	config := NewConfig(options...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writerInterceptor := NewResponseWriterInterceptor(writer)

			req.Body = interceptor.NewCopyReadCloser(req.Body)
			config.CallbackBefore(writerInterceptor, req)
			defer func() {
				config.CallbackAfter(writerInterceptor, req)
			}()

			next.ServeHTTP(writerInterceptor, req)
		})
	}
}

type Config struct {
	CallbackBefore func(*ResponseWriterInterceptor, *http.Request)
	CallbackAfter  func(*ResponseWriterInterceptor, *http.Request)
}

func (c *Config) apply(options ...Option) *Config {
	for _, option := range options {
		option(c)
	}
	return c
}

// NewConfig returns a new interceptor middleware configuration with all options applied
func NewConfig(options ...Option) *Config {
	config := &Config{
		CallbackBefore: func(_ *ResponseWriterInterceptor, _ *http.Request) {},
		CallbackAfter:  func(_ *ResponseWriterInterceptor, _ *http.Request) {},
	}
	return config.apply(options...)
}

// Option defines a interceptor middleware configuration option
type Option func(*Config)

// WithBefore will configure CallbackBefore interceptor option
func WithBefore(callbackBefore func(*ResponseWriterInterceptor, *http.Request)) Option {
	return func(config *Config) {
		config.CallbackBefore = callbackBefore
	}
}

// WithAfter will configure CallbackAfter interceptor option
func WithAfter(callbackAfter func(*ResponseWriterInterceptor, *http.Request)) Option {
	return func(config *Config) {
		config.CallbackAfter = callbackAfter
	}
}
