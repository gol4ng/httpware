package httpware

import "net/http"

// RoundTripperFunc wraps a func to make it into a http.RoundTripper. Similar to http.HandleFunc.
type RoundTripperFunc func(*http.Request) (*http.Response, error)

func (f RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Tripperware is a signature for all http client-side middleware.
type Tripperware func(http.RoundTripper) http.RoundTripper
type Tripperwares []Tripperware

func (t Tripperwares) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.Decorate(http.DefaultTransport).RoundTrip(req)
}

func (t Tripperwares) DecorateFunc(tripper RoundTripperFunc) http.RoundTripper {
	return DecorateRoundTripperFunc(tripper, t...)
}

func (t Tripperwares) Decorate(tripper http.RoundTripper) http.RoundTripper {
	return DecorateRoundTripper(tripper, t...)
}

func DecorateRoundTripperFunc(tripper RoundTripperFunc, tripperwares ...Tripperware) http.RoundTripper {
	if tripper == nil {
		return DecorateRoundTripper(http.DefaultTransport, tripperwares...)
	}
	return DecorateRoundTripper(tripper, tripperwares...)
}

func DecorateRoundTripper(tripper http.RoundTripper, tripperwares ...Tripperware) http.RoundTripper {
	if tripper == nil {
		tripper = http.DefaultTransport
	}
	for _, tripperware := range tripperwares {
		tripper = tripperware(tripper)
	}
	return tripper
}

func TripperwareStack(tripperwares ...Tripperware) Tripperwares {
	return tripperwares
}
