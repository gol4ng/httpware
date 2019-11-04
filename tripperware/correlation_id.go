package tripperware

import (
	"context"
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/correlation_id"
)

// CorrelationId tripperware gets request id header if provided or generates a request id
// It will add the request ID to request context
func CorrelationId(options ...correlation_id.Option) httpware.Tripperware {
	config := correlation_id.GetConfig(options...)
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
			if v, ok := req.Context().Value(config.HeaderName).(string); ok {
				req.Header.Add(config.HeaderName, v)
				return next.RoundTrip(req)
			}

			var id string
			if req.Header != nil {
				id = req.Header.Get(config.HeaderName)
			}

			if id == "" {
				id = config.IdGenerator(req)
				// add requestId header to current request
				req.Header.Add(config.HeaderName, id)
			}
			r := req.WithContext(context.WithValue(req.Context(), config.HeaderName, id))
			return next.RoundTrip(r)
		})
	}
}
