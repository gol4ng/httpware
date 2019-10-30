package httpware

import "net/http"

// Middleware represents an http server middleware
// it wraps an http.Handler with another one
type Middleware func(http.Handler) http.Handler

// Append will add given middlewares after existing one
// t1.Append(t2, t3) => [t1, t2, t3]
// t1.Append(t2, t3).DecorateHandler(<yourHandler>) == t1(t2(t3(<yourHandler>)))
func (m Middleware) Append(middlewares ...Middleware) Middlewares {
	return append([]Middleware{m}, middlewares...)
}

// Prepend will add given middlewares before existing one
// t1.Prepend(t2, t3) => [t2, t3, t1]
// t1.Prepend(t2, t3).DecorateHandler(<yourHandler>) == t2(t3(t1(<yourHandler>)))
func (m Middleware) Prepend(middlewares ...Middleware) Middlewares {
	return append(middlewares, m)
}

// [t1, t2, t3].DecorateHandler(<yourHandler>) == t1(t2(t3(<yourHandler>)))
type Middlewares []Middleware

// DecorateHandler will decorate a given http.Handler with the given middlewares created by MiddlewareStack()
func (m Middlewares) DecorateHandler(handler http.Handler) http.Handler {
	mLen := len(m)
	for i := mLen - 1; i >= 0; i-- {
		handler = m[i](handler)
	}
	return handler
}

// DecorateHandler will decorate a given http.HandlerFunc with the given middleware collection created by MiddlewareStack()
func (m Middlewares) DecorateHandlerFunc(handler http.HandlerFunc) http.Handler {
	return m.DecorateHandler(handler)
}

// Append will add given middleware after existing one
// [t1, t2].Append(t3, t4) => [t1, t2, t3, t4]
// [t1, t2].Append(t3, t4).DecorateHandler(<yourHandler>) == t1(t2(t3(t4(<yourHandler>))))
func (m Middlewares) Append(middleware ...Middleware) Middlewares {
	return append(m, middleware...)
}

// Prepend will add given middleware before existing one
// [t1, t2].Prepend(t3, t4) => [t3, t4, t1, t2]
// [t1, t2].Prepend(t3, t4).DecorateHandler(<yourHandler>) == t3(t4(t1(t2(<yourHandler>))))
func (m Middlewares) Prepend(middleware ...Middleware) Middlewares {
	return append(middleware, m...)
}

// MiddlewareStack allows you to stack multiple middleware collection in a specific order
func MiddlewareStack(middlewares ...Middleware) Middlewares {
	return middlewares
}
