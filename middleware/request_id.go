package middleware

import (
	"context"
	"net/http"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/request_id"
)

// RequestId middleware get request id header if provided or generate a request id
// It will add the request ID to request context and add it to response header to
func RequestId(config *request_id.Config) httpware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			id := req.Header.Get(config.HeaderName)
			if id == "" {
				id = config.IdGenerator(req)
			}
			// add the request id to the current context request
			r := req.WithContext(context.WithValue(req.Context(), config.HeaderName, id))
			// add it to the response headers
			writer.Header().Set(config.HeaderName, id)
			next.ServeHTTP(writer, r)
		})
	}
}
