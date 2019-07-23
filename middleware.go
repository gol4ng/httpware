package httpware

import "net/http"

type Middleware func(http.Handler) http.Handler
type Middlewares []Middleware

func (m Middlewares) DecorateFunc(handler http.HandlerFunc) http.Handler {
	return DecorateHandlerFunc(handler, m...)
}

func (m Middlewares) Decorate(handler http.Handler) http.Handler {
	return DecorateHandler(handler, m...)
}

func DecorateHandlerFunc(handler http.HandlerFunc, middlewares ...Middleware) http.Handler {
	return DecorateHandler(handler, middlewares...)
}

func DecorateHandler(handler http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		handler = middleware(handler)
	}
	return handler
}

func MiddlewareStack(middlewares ...Middleware) Middlewares {
	return middlewares
}
