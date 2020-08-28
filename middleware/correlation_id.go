package middleware

import (
	"context"
	"net/http"

	"github.com/gol4ng/httpware/v3"
	"github.com/gol4ng/httpware/v3/correlation_id"
)

// CorrelationId middleware get request id header if provided or generate a request id
// It will add the request ID to request context and add it to response header to
func CorrelationId(options ...correlation_id.Option) httpware.Middleware {
	config := correlation_id.GetConfig(options...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			id := req.Header.Get(config.HeaderName)
			if id == "" {
				id = config.IdGenerator(req)
				// set requestId header to current request
				req.Header.Set(config.HeaderName, id)
			}
			// add the request id to the current context request
			r := req.WithContext(context.WithValue(req.Context(), config.HeaderName, id))
			// add it to the response headers
			writer.Header().Set(config.HeaderName, id)
			next.ServeHTTP(writer, r)
		})
	}
}
