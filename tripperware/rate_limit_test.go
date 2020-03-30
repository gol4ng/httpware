package tripperware_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gol4ng/httpware/v2/mocks"
	"github.com/gol4ng/httpware/v2/rate_limit"
	"github.com/gol4ng/httpware/v2/tripperware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimit(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
	}

	roundTripperMock.On("RoundTrip", request).Times(2).Return(resp, nil)

	rateLimiterMock := &mocks.RateLimiter{}

	i := 0
	rateLimiterMock.On("Inc").Run(func(mock.Arguments) {
		i++
	})

	call := rateLimiterMock.On("IsLimitReached")
	call.Run(func(mock.Arguments) {
		call.Return(i == 1)
	})

	tr := tripperware.RateLimit(rateLimiterMock)
	_, err := tr(roundTripperMock).RoundTrip(request)
	assert.Nil(t, err)

	_, err = tr(roundTripperMock).RoundTrip(request)
	assert.EqualError(t, err, "request limit reached")

	i = 0
	resp, err = tr(roundTripperMock).RoundTrip(request)
	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	roundTripperMock.AssertExpectations(t)
	rateLimiterMock.AssertExpectations(t)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleRateLimit() {
	rl := rate_limit.NewLeakyBucket(1*time.Second, 1)
	defer rl.Stop()

	client := http.Client{Transport: tripperware.RateLimit(rl)}
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	addr := listener.Addr().String()
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("server receive request")
	})

	go func() {
		if err := http.Serve(listener, srv); err != nil {
			panic(err)
		}
	}()

	_, err = client.Get("http://" + addr + "/")
	fmt.Println(err)

	_, err = client.Get("http://" + addr + "/")
	fmt.Println(errors.Unwrap(err))

	time.Sleep(2 * time.Second)
	_, err = client.Get("http://" + addr + "/")
	fmt.Println(err)
	// Output:
	//server receive request
	//<nil>
	//request limit reached
	//server receive request
	//<nil>
}
