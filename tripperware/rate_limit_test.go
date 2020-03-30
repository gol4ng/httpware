package tripperware_test

import (
	"fmt"
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

	port := ":9004"
	client := http.Client{
		Transport: tripperware.RateLimit(rl),
	}

	// create a server in order to show it work
	srv := http.NewServeMux()

	srv.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("server receive request")
	})

	go func() {
		if err := http.ListenAndServe(port, srv); err != nil {
			panic(err)
		}
	}()

	_, err := client.Get("http://localhost" + port + "/")
	fmt.Println(err)

	_, err = client.Get("http://localhost" + port + "/")
	fmt.Println(err)

	time.Sleep(2 * time.Second)
	_, err = client.Get("http://localhost" + port + "/")
	fmt.Println(err)
	// Output:
	//server receive request
	//<nil>
	//Get http://localhost:9004/: request limit reached
	//server receive request
	//<nil>
}
