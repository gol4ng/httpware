package tripperware

import (
	"github.com/gol4ng/httpware/v4/request_listener"
	"net/http"

	"github.com/gol4ng/httpware/v4"
)

func RequestListener(listeners ...func(*http.Request)) httpware.Tripperware {
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			for _, listener := range listeners {
				listener(request)
			}
			return next.RoundTrip(request)
		})
	}
}

func CurlLogDumper() httpware.Tripperware {
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			request_listener.CurlLogDumper(request)
			return next.RoundTrip(request)
		})
	}
}
