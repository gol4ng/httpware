package httpware

import "net/http"

// Middleware represents an http server middleware
// it wraps an http.Handler with another one
type Middleware func(http.Handler) http.Handler
type Middlewares []Middleware

// DecorateHandler will decorate a given http.Handler with the given middlewares created by MiddlewareStack()
func (m Middlewares) DecorateHandler(handler http.Handler) http.Handler {
	for _, middleware := range m {
		handler = middleware(handler)
	}
	return handler
}

// DecorateHandler will decorate a given http.HandlerFunc with the given middleware collection created by MiddlewareStack()
func (m Middlewares) DecorateHandlerFunc(handler http.HandlerFunc) http.Handler {
	return m.DecorateHandler(handler)
}

// MiddlewareStack allows you to stack multiple middleware collection in a specific order
func MiddlewareStack(middlewares ...Middleware) Middlewares {
	return middlewares
}
