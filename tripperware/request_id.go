package tripperware

import (
	"context"
	"net/http"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/request_id"
)

// RequestId tripperware gets request id header if provided or generates a request id
// It will add the request ID to request context
func RequestId(config *request_id.Config) httpware.Tripperware {
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
			var id string
			if req.Header != nil {
				id = req.Header.Get(config.HeaderName)
			}
			if id == "" {
				id = config.IdGenerator(req)
			}
			r := req.WithContext(context.WithValue(req.Context(), config.HeaderName, id))
			r.Header.Add(config.HeaderName, id)
			return next.RoundTrip(r)
		})
	}
}
