package httpware

import "net/http"

// RoundTripFunc wraps a func to make it into an http.RoundTripper. Similar to http.HandleFunc.
type RoundTripFunc func(*http.Request) (*http.Response, error)

// RoundTrip implements RoundTripper interface
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Tripperware represents an http client-side middleware (roundTripper middleware).
type Tripperware func(http.RoundTripper) http.RoundTripper

// DecorateClient will decorate a given http.Client with the tripperware
func (t Tripperware) DecorateClient(client *http.Client) *http.Client {
	if client.Transport == nil {
		client.Transport = http.DefaultTransport
	}
	client.Transport = t(client.Transport)
	return client
}

type Tripperwares []Tripperware

// RoundTrip implements RoundTripper interface
// it will decorate the http-client request and use the default `http.DefaultTransport` RoundTripper
// use `TripperwareStack(<yourTripperwares>).Decorate(<yourTripper>)` if you don't want to use `http.DefaultTransport`
func (t Tripperwares) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.DecorateRoundTripper(http.DefaultTransport).RoundTrip(req)
}

// DecorateClient will decorate a given http.Client with the given tripperware collection
func (t Tripperwares) DecorateClient(client *http.Client) *http.Client {
	if client.Transport == nil {
		client.Transport = http.DefaultTransport
	}
	client.Transport = t.DecorateRoundTripper(client.Transport)
	return client
}

// DecorateRoundTripper will decorate a given http.RoundTripper with the given tripperware collection
func (t Tripperwares) DecorateRoundTripper(tripper http.RoundTripper) http.RoundTripper {
	if tripper == nil {
		tripper = http.DefaultTransport
	}
	for _, tripperware := range t {
		tripper = tripperware(tripper)
	}
	return tripper
}

// DecorateRoundTripFunc will decorate a given RoundTripFunc with the given tripperware collection
func (t Tripperwares) DecorateRoundTripFunc(tripper RoundTripFunc) http.RoundTripper {
	if tripper == nil {
		return t.DecorateRoundTripper(http.DefaultTransport)
	}
	return t.DecorateRoundTripper(tripper)
}

// TripperwareStack allows to stack multi tripperware in order to decorate an http roundTripper
func TripperwareStack(tripperwares ...Tripperware) Tripperwares {
	return tripperwares
}
