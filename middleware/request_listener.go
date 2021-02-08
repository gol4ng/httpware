package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware/v4"
)

func RequestListener(listeners ...func(*http.Request)) httpware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			for _, listener := range listeners {
				listener(request)
			}
			next.ServeHTTP(writer, request)
		})
	}
}
