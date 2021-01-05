package tripperware

import (
	"net/http"

	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/skip"
)

// Skip tripperware is used to conditionnaly activate a tripperware in function of request
// See Enable tripperware to conditionnaly add tripperware to a stack
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
