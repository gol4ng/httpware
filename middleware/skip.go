package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware/v3"
	"github.com/gol4ng/httpware/v3/skip"
)

// Skip middleware is used to conditionnaly activate a middleware in function of request
// See Enable middleware to conditionnaly add middleware to a stack
func Skip(condition skip.Condition, middleware httpware.Middleware) httpware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			if condition(req) {
				next.ServeHTTP(writer, req)
				return
			}
			middleware(next).ServeHTTP(writer, req)
		})
	}
}
