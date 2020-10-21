package tripperware_test

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gol4ng/httpware/v3/mocks"
	"github.com/gol4ng/httpware/v3/rate_limit"
	"github.com/gol4ng/httpware/v3/tripperware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimit(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:     "OK",
		StatusCode: http.StatusOK,
	}

	roundTripperMock.On("RoundTrip", request).Times(2).Return(resp, nil)

	rateLimiterMock := &mocks.RateLimiter{}
	rateLimiterMock.On("Dec", mock.AnythingOfType("*http.Request")).Return()

	i := 0
	rateLimiterMock.On("Inc", mock.AnythingOfType("*http.Request")).Run(func(mock.Arguments) {
		i++
	})

	call := rateLimiterMock.On("Allow", mock.AnythingOfType("*http.Request"))
	call.Run(func(mock.Arguments) {
		if i == 1 {
			call.Return(errors.New(rate_limit.RequestLimitReachedErr))
			return
		}

		call.Return(nil)
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

func TestNewConfig(t *testing.T) {
	res, err := tripperware.NewRateLimitConfig().ErrorCallback(nil, fmt.Errorf("default error"))
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(
		t,
		"default error",
		err.Error(),
	)
}

func TestConfig_Options(t *testing.T) {
	config := tripperware.NewRateLimitConfig(
		tripperware.WithRateLimitErrorCallback(func(*http.Request, error) (*http.Response, error) {
			return nil, fmt.Errorf("error from callback")
		}),
	)

	res, err := config.ErrorCallback(nil, nil)
	assert.Nil(t, res)
	assert.NotNil(t, err)
	assert.Equal(t, "error from callback", err.Error())
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleRateLimit() {
	// Example Need a random ephemeral port (to have a free port)
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	rl := rate_limit.NewTokenBucket(1*time.Second, 1)
	defer rl.Stop()

	client := http.Client{Transport: tripperware.RateLimit(rl)}

	srv := &http.Server{
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			fmt.Println("server receive request")
		}),
	}
	go func() {
		if err := srv.Serve(ln); err != nil {
			panic(err)
		}
	}()

	_, err = client.Get("http://" + ln.Addr().String())
	fmt.Println(err)

	_, err = client.Get("http://" + ln.Addr().String())
	fmt.Println(errors.Unwrap(err))

	time.Sleep(2 * time.Second)
	_, err = client.Get("http://" + ln.Addr().String())
	fmt.Println(err)
	// Output:
	//server receive request
	//<nil>
	//request limit reached
	//server receive request
	//<nil>
}
