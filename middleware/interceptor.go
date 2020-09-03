package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware/v3"
	"github.com/gol4ng/httpware/v3/interceptor"
)

// Interceptor middleware allow multiple req.Body read and allow to set callback before and after roundtrip
func Interceptor(options ...InterceptorOption) httpware.Middleware {
	config := NewInterceptorConfig(options...)
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

type InterceptorConfig struct {
	CallbackBefore func(*ResponseWriterInterceptor, *http.Request)
	CallbackAfter  func(*ResponseWriterInterceptor, *http.Request)
}

func (c *InterceptorConfig) apply(options ...InterceptorOption) *InterceptorConfig {
	for _, option := range options {
		option(c)
	}
	return c
}

// NewInterceptorConfig returns a new interceptor middleware configuration with all options applied
func NewInterceptorConfig(options ...InterceptorOption) *InterceptorConfig {
	config := &InterceptorConfig{
		CallbackBefore: func(_ *ResponseWriterInterceptor, _ *http.Request) {},
		CallbackAfter:  func(_ *ResponseWriterInterceptor, _ *http.Request) {},
	}
	return config.apply(options...)
}

// InterceptorOption defines a interceptor middleware configuration option
type InterceptorOption func(*InterceptorConfig)

// WithBefore will configure CallbackBefore interceptor option
func WithBefore(callbackBefore func(*ResponseWriterInterceptor, *http.Request)) InterceptorOption {
	return func(config *InterceptorConfig) {
		config.CallbackBefore = callbackBefore
	}
}

// WithAfter will configure CallbackAfter interceptor option
func WithAfter(callbackAfter func(*ResponseWriterInterceptor, *http.Request)) InterceptorOption {
	return func(config *InterceptorConfig) {
		config.CallbackAfter = callbackAfter
	}
}
