package tripperware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/mocks"
	"github.com/gol4ng/httpware/request_id"
	"github.com/gol4ng/httpware/tripperware"
)

func TestRequestId(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	roundTripperMock.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil).Run(func(args mock.Arguments) {
		innerReq := args.Get(0).(*http.Request)
		assert.True(t, len(innerReq.Header.Get(request_id.HeaderName)) == 10)
	})

	resp2, err := tripperware.RequestId(request_id.NewConfig())(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, resp, resp2)
}

func TestRequestIdCustom(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	roundTripperMock.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil).Run(func(args mock.Arguments) {
		innerReq := args.Get(0).(*http.Request)
		assert.Equal(t, "my_fake_request_id", innerReq.Header.Get(request_id.HeaderName))
	})

	config := request_id.NewConfig()
	config.IdGenerator = func(request *http.Request) string {
		return "my_fake_request_id"
	}

	resp2, err := tripperware.RequestId(config)(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, resp, resp2)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleRequestId() {
	port := ":5005"
	config := request_id.NewConfig()
	// you can override default header name
	config.HeaderName = "my-personal-header-name"
	// you can override default id generator
	config.IdGenerator = func(request *http.Request) string {
		return "my-generated-id"
	}

	// we recommend to use MiddlewareStack to simplify managing all wanted middleware
	// caution middleware order matter
	stack := httpware.TripperwareStack(
		tripperware.RequestId(config),
	)

	// create http client using the tripperwareStack as RoundTripper
	client := http.Client{
		Transport: stack,
	}

	// create a server in order to show it work
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("server receive request with request id:", request.Header.Get(config.HeaderName))
	})

	go func() {
		if err := http.ListenAndServe(port, srv); err != nil {
			panic(err)
		}
	}()

	_, _ = client.Get("http://localhost"+port+"/")

	// Output: server receive request with request id: my-generated-id
}