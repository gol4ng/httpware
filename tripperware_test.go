package httpware_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/mocks"
)

func getTripper(t *testing.T, i *int, iBefore int, iAfter int) httpware.Tripperware {
	return httpware.Tripperware(func(roundTripper http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
			defer func() {
				assert.Equal(t, iAfter, *i)
				*i++
			}()
			assert.Equal(t, iBefore, *i)
			*i++
			return roundTripper.RoundTrip(req)
		})
	})
}

func TestTripperware_RoundTrip(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	resp := &http.Response{}

	roundTripperMock := &mocks.RoundTripper{}
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)

	originalDefaultTransport := http.DefaultTransport
	http.DefaultTransport = roundTripperMock

	tripper := httpware.Tripperware(func(roundTripper http.RoundTripper) http.RoundTripper {
		assert.Equal(t, http.DefaultTransport, roundTripper)
		return roundTripper
	})

	r, err := tripper.RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, r, resp)

	http.DefaultTransport = originalDefaultTransport
}

func TestTripperware_DecorateClient(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	resp := &http.Response{}

	roundTripperMock := &mocks.RoundTripper{}
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)

	tripperware := httpware.Tripperware(func(tripper http.RoundTripper) http.RoundTripper {
		assert.Equal(t, http.DefaultTransport, tripper)
		return httpware.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, req, r)
			// we already check that tripper == http.DefaultTransport
			// so we can replace the call with the mocked one
			return roundTripperMock.RoundTrip(r)
		})
	})

	t.Run("Clone", func(tt *testing.T) {
		client := &http.Client{}
		newClient := tripperware.DecorateClient(client, true)
		assert.NotEqual(t, client, newClient)

		response, err := newClient.Do(req)
		assert.Equal(tt, resp, response)
		assert.Nil(tt, err)
	})

	t.Run("Wrap", func(tt *testing.T) {
		client := &http.Client{}
		newClient := tripperware.DecorateClient(client, false)
		assert.Equal(tt, client, newClient)

		response, err := newClient.Do(req)
		assert.Equal(tt, resp, response)
		assert.Nil(tt, err)
	})

	t.Run("Clone Nil", func(tt *testing.T) {
		newClient := tripperware.DecorateClient(nil, true)
		assert.NotEqual(tt, http.DefaultClient, newClient)

		response, err := newClient.Do(req)
		assert.Equal(tt, resp, response)
		assert.Nil(tt, err)
	})

	t.Run("Wrap Nil", func(tt *testing.T) {
		newClient := tripperware.DecorateClient(nil, false)
		assert.Equal(tt, http.DefaultClient, newClient)

		response, err := newClient.Do(req)
		assert.Equal(tt, resp, response)
		assert.Nil(tt, err)
	})

	http.DefaultClient.Transport = http.DefaultTransport
}

func TestTripperware_Append(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	resp := &http.Response{}

	i := new(int)
	*i = 0
	roundTripperMock := httpware.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		assert.Equal(t, 3, *i)
		*i++
		return resp, nil
	})

	tripper := getTripper(t, i, 0, 6)

	r, err := tripper.Append(
		// the tripper will be add here
		getTripper(t, i, 1, 5),
		getTripper(t, i, 2, 4),
	).DecorateRoundTripper(roundTripperMock).RoundTrip(req)

	assert.Nil(t, err)
	assert.Equal(t, r, resp)
}

func TestTripperware_Prepend(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	resp := &http.Response{}

	i := new(int)
	*i = 0
	roundTripperMock := httpware.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		assert.Equal(t, 3, *i)
		*i++
		return resp, nil
	})

	tripper := getTripper(t, i, 2, 4)

	r, err := tripper.Prepend(
		getTripper(t, i, 0, 6),
		getTripper(t, i, 1, 5),
		// the tripper will be add here
	).DecorateRoundTripper(roundTripperMock).RoundTrip(req)

	assert.Nil(t, err)
	assert.Equal(t, r, resp)
}

func TestTripperwares_DecorateRoundTripper(t *testing.T) {
	req := &http.Request{}
	resp := &http.Response{}

	roundTripperMock := &mocks.RoundTripper{}
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)

	stack := httpware.TripperwareStack(func(tripper http.RoundTripper) http.RoundTripper {
		assert.Equal(t, http.DefaultTransport, tripper)
		return httpware.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, req, r)
			// we already check that tripper == http.DefaultTransport
			// so we can replace the call with the mocked one
			return roundTripperMock.RoundTrip(r)
		})
	})

	response, err := stack.DecorateRoundTripper(nil).RoundTrip(req)
	assert.Equal(t, resp, response)
	assert.Nil(t, err)
}

func TestTripperwares_DecorateClient(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	resp := &http.Response{}

	roundTripperMock := &mocks.RoundTripper{}
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)

	stack := httpware.TripperwareStack(func(tripper http.RoundTripper) http.RoundTripper {
		assert.Equal(t, http.DefaultTransport, tripper)
		return httpware.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, req, r)
			// we already check that tripper == http.DefaultTransport
			// so we can replace the call with the mocked one
			return roundTripperMock.RoundTrip(r)
		})
	})

	t.Run("Clone", func(tt *testing.T) {
		client := &http.Client{}
		client2 := stack.DecorateClient(client, true)
		assert.NotEqual(t, client, client2)

		response, err := client2.Do(req)
		assert.Equal(t, resp, response)
		assert.Nil(t, err)
	})

	t.Run("Wrap", func(tt *testing.T) {
		client := &http.Client{}
		client2 := stack.DecorateClient(client, false)
		assert.Equal(t, client, client2)

		response, err := client2.Do(req)
		assert.Equal(t, resp, response)
		assert.Nil(t, err)
	})

	t.Run("Clone Nil", func(tt *testing.T) {
		client2 := stack.DecorateClient(nil, true)
		assert.NotEqual(t, http.DefaultClient, client2)

		response, err := client2.Do(req)
		assert.Equal(t, resp, response)
		assert.Nil(t, err)
	})

	t.Run("Wrap Nil", func(tt *testing.T) {
		client2 := stack.DecorateClient(nil, false)
		assert.Equal(t, http.DefaultClient, client2)

		response, err := client2.Do(req)
		assert.Equal(t, resp, response)
		assert.Nil(t, err)
	})

	http.DefaultClient.Transport = http.DefaultTransport
}

func TestTripperwares_DecorateRoundTripFunc_Default(t *testing.T) {
	req := &http.Request{}
	resp := &http.Response{}

	roundTripperMock := &mocks.RoundTripper{}
	roundTripperMock.On("RoundTrip", req).Return(resp, nil).Once()

	stack := httpware.TripperwareStack(func(tripper http.RoundTripper) http.RoundTripper {
		assert.Equal(t, http.DefaultTransport, tripper)
		return httpware.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, req, r)
			// we already check that tripper == http.DefaultTransport
			// so we can replace the call with the mocked one
			return roundTripperMock.RoundTrip(r)
		})
	})

	response, err := stack.DecorateRoundTripFunc(nil).RoundTrip(req)
	assert.Equal(t, resp, response)
	assert.Nil(t, err)
}

func TestTripperwares_DecorateRoundTripFunc(t *testing.T) {
	req := &http.Request{}
	resp := &http.Response{}

	roundTripperMock := &mocks.RoundTripper{}
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)

	roundTripperFuncCalled := false
	roundTripperFunc := func(r *http.Request) (*http.Response, error) {
		assert.Equal(t, req, r)
		roundTripperFuncCalled = true
		return resp, nil
	}

	stack := httpware.TripperwareStack(func(tripper http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(r *http.Request) (*http.Response, error) {
			assert.Equal(t, req, r)
			// we already check that tripper == http.DefaultTransport
			// so we can replace the call with the mocked one
			return tripper.RoundTrip(r)
		})
	})

	response, err := stack.DecorateRoundTripFunc(roundTripperFunc).RoundTrip(req)
	assert.Equal(t, resp, response)
	assert.Nil(t, err)
	assert.True(t, roundTripperFuncCalled)
}

func TestTripperwares_Append(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	resp := &http.Response{}

	i := new(int)
	*i = 0
	roundTripperMock := httpware.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		assert.Equal(t, 4, *i)
		*i++
		return resp, nil
	})

	trippers := httpware.TripperwareStack(
		getTripper(t, i, 0, 8),
		getTripper(t, i, 1, 7),
	)

	r, err := trippers.Append(
		// the tripper will be add here
		getTripper(t, i, 2, 6),
		getTripper(t, i, 3, 5),
	).DecorateRoundTripper(roundTripperMock).RoundTrip(req)

	assert.Nil(t, err)
	assert.Equal(t, r, resp)
}

func TestTripperwares_Prepend(t *testing.T) {
	req, _ := http.NewRequest("GET", "http://localhost/", nil)
	resp := &http.Response{}

	i := new(int)
	*i = 0
	roundTripperMock := httpware.RoundTripFunc(func(*http.Request) (*http.Response, error) {
		assert.Equal(t, 4, *i)
		*i++
		return resp, nil
	})

	trippers := httpware.TripperwareStack(
		getTripper(t, i, 2, 6),
		getTripper(t, i, 3, 5),
	)

	r, err := trippers.Prepend(
		getTripper(t, i, 0, 8),
		getTripper(t, i, 1, 7),
		// the tripper will be add here
	).DecorateRoundTripper(roundTripperMock).RoundTrip(req)

	assert.Nil(t, err)
	assert.Equal(t, r, resp)
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
		addCustomRequestHeader,
		logRequestHeaders,
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
		addCustomRequestHeader,
		logRequestHeaders,
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
