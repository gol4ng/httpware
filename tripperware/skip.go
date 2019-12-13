package tripperware

import (
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/skip"
)

func Skip(condition skip.Condition, tripperware httpware.Tripperware) httpware.Tripperware {
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			if condition(request) {
				return next.RoundTrip(request)
			}
			return tripperware(next).RoundTrip(request)
		})
	}
}
