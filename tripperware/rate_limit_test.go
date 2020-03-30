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
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	roundTripperMock.On("RoundTrip", req).Return(resp, nil)

	rateLimiterMock := &mocks.RateLimiter{}

	i := 0
	rateLimiterMock.On("Inc").Run(func(args mock.Arguments) {
		i++
	})

	call := rateLimiterMock.On("IsLimitReached")
	call.Run(func(args mock.Arguments) {
		call.Return(i == 1)
	})

	tr := tripperware.RateLimit(rateLimiterMock)
	_, err := tr(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)

	_, err = tr(roundTripperMock).RoundTrip(req)
	assert.NotNil(t, err)
	assert.Equal(t, "request limit reached", err.Error())

	i = 0
	resp, err = tr(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, int64(30), resp.ContentLength)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleRateLimit() {
	rl := rate_limit.NewLeakyBucket(1*time.Second, 1)
	defer rl.Stop()

	port := ":9003"
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
	//Get http://localhost:9003/: request limit reached
	//server receive request
	//<nil>
}
