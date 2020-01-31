package httpware

import "net/http"

// NopMiddleware just return given http.Handler
func NopMiddleware(next http.Handler) http.Handler {
	return next
}

// Middleware represents an http server middleware
// it wraps an http.Handler with another one
type Middleware func(http.Handler) http.Handler

// Append will add given middlewares after existing one
// t1.Append(t2, t3) => [t1, t2, t3]
// t1.Append(t2, t3).DecorateHandler(<yourHandler>) == t1(t2(t3(<yourHandler>)))
func (m Middleware) Append(middlewares ...Middleware) Middlewares {
	return append([]Middleware{m}, middlewares...)
}

// AppendIf will add given middlewares after existing one if condition=true
// t1.AppendIf(true, t2, t3) => [t1, t2, t3]
// t1.AppendIf(false, t2, t3) => [t1]
// t1.AppendIf(true, t2, t3).DecorateHandler(<yourHandler>) == t1(t2(t3(<yourHandler>)))
func (m Middleware) AppendIf(condition bool, middlewares ...Middleware) Middlewares {
	return (&Middlewares{m}).AppendIf(condition, middlewares...)
}

// Prepend will add given middlewares before existing one
// t1.Prepend(t2, t3) => [t2, t3, t1]
// t1.Prepend(t2, t3).DecorateHandler(<yourHandler>) == t2(t3(t1(<yourHandler>)))
func (m Middleware) Prepend(middlewares ...Middleware) Middlewares {
	return append(middlewares, m)
}

// PrependIf will add given middlewares before existing one if condition=true
// t1.PrependIf(true, t2, t3) => [t2, t3, t1]
// t1.PrependIf(false, t2, t3) => [t1]
// t1.PrependIf(true, t2, t3).DecorateHandler(<yourHandler>) == t2(t3(t1(<yourHandler>)))
func (m Middleware) PrependIf(condition bool, middlewares ...Middleware) Middlewares {
	return (&Middlewares{m}).PrependIf(condition, middlewares...)
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
func (m *Middlewares) Append(middleware ...Middleware) Middlewares {
	*m = append(*m, middleware...)
	return *m
}

// AppendIf will add given middleware after existing one if condition=true
// [t1, t2].AppendIf(true, t3, t4) => [t1, t2, t3, t4]
// [t1, t2].AppendIf(false, t3, t4) => [t1, t2]
// [t1, t2].AppendIf(t3, t4).DecorateHandler(<yourHandler>) == t1(t2(t3(t4(<yourHandler>))))
func (m *Middlewares) AppendIf(condition bool, middleware ...Middleware) Middlewares {
	if condition {
		*m = append(*m, middleware...)
	}
	return *m
}

// Prepend will add given middleware before existing one
// [t1, t2].Prepend(t3, t4) => [t3, t4, t1, t2]
// [t1, t2].Prepend(t3, t4).DecorateHandler(<yourHandler>) == t3(t4(t1(t2(<yourHandler>))))
func (m *Middlewares) Prepend(middleware ...Middleware) Middlewares {
	*m = append(middleware, *m...)
	return *m
}

// PrependIf will add given middleware before existing one if condition=true
// [t1, t2].PrependIf(true, t3, t4) => [t3, t4, t1, t2]
// [t1, t2].PrependIf(false, t3, t4) => [t1, t2]
// [t1, t2].PrependIf(true, t3, t4).DecorateHandler(<yourHandler>) == t3(t4(t1(t2(<yourHandler>))))
func (m *Middlewares) PrependIf(condition bool, middleware ...Middleware) Middlewares {
	if condition {
		*m = append(middleware, *m...)
	}
	return *m
}

// MiddlewareStack allows you to stack multiple middleware collection in a specific order
func MiddlewareStack(middlewares ...Middleware) Middlewares {
	return middlewares
}
