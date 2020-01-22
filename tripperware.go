package httpware

import (
	"net/http"
)

// NopTripperware just return given http.RoundTripper
func NopTripperware(next http.RoundTripper) http.RoundTripper {
	return next
}

// RoundTripFunc wraps a func to make it into an http.RoundTripper. Similar to http.HandleFunc.
type RoundTripFunc func(*http.Request) (*http.Response, error)

// RoundTrip implements RoundTripper interface
func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// Tripperware represents an http client-side middleware (roundTripper middleware).
type Tripperware func(http.RoundTripper) http.RoundTripper

// RoundTrip implements RoundTripper interface
func (t Tripperware) RoundTrip(req *http.Request) (*http.Response, error) {
	return t(http.DefaultTransport).RoundTrip(req)
}

// DecorateClient will decorate a given http.Client with the tripperware
// will return a clone of client if clone arg is true
func (t Tripperware) DecorateClient(client *http.Client, clone bool) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}
	if clone {
		c := *client
		client = &c
	}
	if client.Transport == nil {
		client.Transport = http.DefaultTransport
	}
	client.Transport = t(client.Transport)
	return client
}

// Append will add given tripperwares after existing one
// t1.Append(t2, t3) == [t1, t2, t3]
// t1.Append(t2, t3).DecorateRoundTripper(<yourTripper>) == t1(t2(t3(<yourTripper>)))
func (t Tripperware) Append(tripperwares ...Tripperware) Tripperwares {
	return append([]Tripperware{t}, tripperwares...)
}

// AppendIf will add given tripperwares after existing one if condition=true
// t1.AppendIf(true, t2, t3) == [t1, t2, t3]
// t1.AppendIf(false, t2, t3) == [t1]
// t1.AppendIf(true, t2, t3).DecorateRoundTripper(<yourTripper>) == t1(t2(t3(<yourTripper>)))
func (t Tripperware) AppendIf(condition bool, tripperwares ...Tripperware) Tripperwares {
	return (&Tripperwares{t}).AppendIf(condition, tripperwares...)
}

// Prepend will add given tripperwares before existing one
// t1.Prepend(t2, t3) => [t2, t3, t1]
// t1.Prepend(t2, t3).DecorateRoundTripper(<yourTripper>) == t2(t3(t1(<yourTripper>)))
func (t Tripperware) Prepend(tripperwares ...Tripperware) Tripperwares {
	return append(tripperwares, t)
}

// PrependIf will add given tripperwares before existing one if condition=true
// t1.PrependIf(true, t2, t3) => [t2, t3, t1]
// t1.PrependIf(false, t2, t3) => [t1]
// t1.PrependIf(true, t2, t3).DecorateRoundTripper(<yourTripper>) == t2(t3(t1(<yourTripper>)))
func (t Tripperware) PrependIf(condition bool, tripperwares ...Tripperware) Tripperwares {
	return (&Tripperwares{t}).PrependIf(condition, tripperwares...)
}

// [t1, t2, t3].DecorateRoundTripper(roundTripper) == t1(t2(t3(roundTripper)))
type Tripperwares []Tripperware

// RoundTrip implements RoundTripper interface
// it will decorate the http-client request and use the default `http.DefaultTransport` RoundTripper
// use `TripperwareStack(<yourTripperwares>).Decorate(<yourTripper>)` if you don't want to use `http.DefaultTransport`
func (t Tripperwares) RoundTrip(req *http.Request) (*http.Response, error) {
	return t.DecorateRoundTripper(http.DefaultTransport).RoundTrip(req)
}

// DecorateClient will decorate a given http.Client with the tripperware collection
// will return a clone of client if clone arg is true
func (t Tripperwares) DecorateClient(client *http.Client, clone bool) *http.Client {
	if client == nil {
		client = http.DefaultClient
	}
	if clone {
		c := *client
		client = &c
	}
	client.Transport = t.DecorateRoundTripper(client.Transport)
	return client
}

// DecorateRoundTripper will decorate a given http.RoundTripper with the tripperware collection
func (t Tripperwares) DecorateRoundTripper(tripper http.RoundTripper) http.RoundTripper {
	if tripper == nil {
		tripper = http.DefaultTransport
	}
	tLen := len(t)
	for i := tLen - 1; i >= 0; i-- {
		tripper = t[i](tripper)
	}
	return tripper
}

// DecorateRoundTripFunc will decorate a given RoundTripFunc with the tripperware collection
func (t Tripperwares) DecorateRoundTripFunc(tripper RoundTripFunc) http.RoundTripper {
	if tripper == nil {
		return t.DecorateRoundTripper(http.DefaultTransport)
	}
	return t.DecorateRoundTripper(tripper)
}

// Append will add given tripperwares after existing one
// [t1, t2].Append(true, t3, t4) == [t1, t2, t3, t4]
// [t1, t2].Append(t3, t4).DecorateRoundTripper(<yourTripper>) == t1(t2(t3(t4(<yourTripper>))))
func (t *Tripperwares) Append(tripperwares ...Tripperware) Tripperwares {
	*t = append(*t, tripperwares...)
	return *t
}

// Appendif will add given tripperwares after existing one if condition=true
// [t1, t2].AppendIf(true, t3, t4) == [t1, t2, t3, t4]
// [t1, t2].AppendIf(false, t3, t4) == [t1, t2]
// [t1, t2].AppendIf(true, t3, t4).DecorateRoundTripper(<yourTripper>) == t1(t2(t3(t4(<yourTripper>))))
func (t *Tripperwares) AppendIf(condition bool, tripperwares ...Tripperware) Tripperwares {
	if condition {
		t.Append(tripperwares...)
	}
	return *t
}

// Prepend will add given tripperwares before existing one
// [t1, t2].Prepend(t3, t4) == [t3, t4, t1, t2]
// [t1, t2].Prepend(t3, t4).DecorateRoundTripper(<yourTripper>) == t3(t4(t1(t2(<yourTripper>))))
func (t *Tripperwares) Prepend(tripperwares ...Tripperware) Tripperwares {
	*t = append(tripperwares, *t...)
	return *t
}

// Prependif will add given tripperwares before existing one if condition=true
// [t1, t2].PrependIf(true, t3, t4) == [t1, t2, t3, t4]
// [t1, t2].PrependIf(false, t3, t4) == [t1, t2]
// [t1, t2].PrependIf(true, t3, t4).DecorateRoundTripper(<yourTripper>) == t1(t2(t3(t4(<yourTripper>))))
func (t *Tripperwares) PrependIf(condition bool, tripperwares ...Tripperware) Tripperwares {
	if condition {
		t.Prepend(tripperwares...)
	}
	return *t
}

// TripperwareStack allows to stack multi tripperware in order to decorate an http roundTripper
func TripperwareStack(tripperwares ...Tripperware) Tripperwares {
	return tripperwares
}
