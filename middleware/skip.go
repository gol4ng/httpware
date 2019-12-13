package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/skip"
)

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