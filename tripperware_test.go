package httpware_test

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/mocks"
)

func TestRoundTripFunc_RoundTrip(t *testing.T) {
	var roundTripperMock = &mocks.RoundTripper{}
	var req *http.Request
	var resp = &http.Response{
		Status:           "OK",
		StatusCode:       http.StatusOK,
		ContentLength:    30,
	}

	req = httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	// mock roundTripper calls
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)
	baseTransport := http.DefaultTransport
	http.DefaultTransport = roundTripperMock
	defer func() { http.DefaultTransport = baseTransport }()

	stack := httpware.TripperwareStack(func(tripper http.RoundTripper) http.RoundTripper {
		assert.Equal(t, http.DefaultTransport, tripper)
		return httpware.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, req, r)
			return tripper.RoundTrip(r)
		})
	})

	response, err := stack.DecorateRoundTripper(nil).RoundTrip(req)
	assert.Equal(t, resp, response)
	assert.Nil(t, err)
}

// =====================================================================================================================
// =============================== use those examples when declaring an http CLIENT ====================================
// =====================================================================================================================

func ExampleTripperwareStack_WithDefaultTransport() {
	// create a tripperware that adds a custom header on each http-client request
	addCustomRequestHeader := func(t http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Add("custom-header", "wonderful header value")
			return t.RoundTrip(req)
		})
	}
	// create a tripperware that logs the request header on each call
	logRequestHeaders := func(t http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			fmt.Println("http request headers :", req.Header)
			return t.RoundTrip(req)
		})
	}
	// create the RoundTripper stack
	//
	// /!\ note that the tripperware order is important here
	// each request will pass through `addCustomRequestHeader` before `logRequestHeaders`
	stack := httpware.TripperwareStack(
		logRequestHeaders,
		addCustomRequestHeader,
	)
	// create http client using the tripperwareStack as RoundTripper
	client := http.Client{
		Transport: stack,
	}

	_, _ = client.Get("fake-address.foo")

	//Output:
	//http request headers : map[Custom-Header:[wonderful header value]]
}

func ExampleTripperwareStack_WithCustomTransport() {
	// create a tripperware that adds a custom header on each http-client request
	addCustomRequestHeader := func(t http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			req.Header.Add("custom-header", "wonderful header value")
			return t.RoundTrip(req)
		})
	}
	// create a tripperware that logs the request header on each call
	logRequestHeaders := func(t http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			fmt.Println("http request headers :", req.Header)
			return t.RoundTrip(req)
		})
	}
	// create the RoundTripper stack
	//
	// /!\ note that the tripperware order is important here
	// each request will pass through `addCustomRequestHeader` before `logRequestHeaders`
	stack := httpware.TripperwareStack(
		logRequestHeaders,
		addCustomRequestHeader,
	)

	// http.Transport implements RoundTripper interface
	customTransport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second, // changed default timeout from 30 to 5
			KeepAlive: 5 * time.Second, // changed default keepAlive from 30 to 5
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          50, // changed MaxIdleConns from 100 to 50
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	// create http client using the tripperwareStack as RoundTripper AND use the custom transport
	client := http.Client{
		Transport: stack.DecorateRoundTripper(customTransport),
	}

	_, _ = client.Get("fake-address.foo")

	//Output:
	//http request headers : map[Custom-Header:[wonderful header value]]
}
