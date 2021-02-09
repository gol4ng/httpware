package middleware_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/middleware"
	"github.com/gol4ng/httpware/v4/mocks"
	"github.com/gol4ng/httpware/v4/rate_limit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimit(t *testing.T) {
	rateLimiterMock := &mocks.RateLimiter{}
	rateLimiterMock.On("Allow", mock.AnythingOfType("*http.Request")).Return(errors.New("failed"))

	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	responseWriter := httptest.NewRecorder()

	executed := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		executed = true
	})

	middleware.RateLimit(rateLimiterMock)(handler).ServeHTTP(responseWriter, req)

	assert.False(t, executed)
	assert.Equal(t, http.StatusTooManyRequests, responseWriter.Result().StatusCode)

	content, err := ioutil.ReadAll(responseWriter.Result().Body)
	assert.NoError(t, err)
	assert.Equal(t, "failed\n", string(content))

	rateLimiterMock.AssertExpectations(t)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleRateLimit() {
	// Example Need a random ephemeral port (to have a free port)
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	limiter := rate_limit.NewTokenBucket(1*time.Second, 1)
	defer limiter.Stop()

	// we recommend to use MiddlewareStack to simplify managing all wanted middlewares
	// caution middleware order matters
	stack := httpware.MiddlewareStack(
		middleware.RateLimit(limiter),
	)

	srv := &http.Server{
		Handler: stack.DecorateHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {}),
	}
	go func() {
		if err := srv.Serve(ln); err != nil {
			panic(err)
		}
	}()

	resp, _ := http.Get("http://" + ln.Addr().String())
	fmt.Println(resp.StatusCode)

	resp, _ = http.Get("http://" + ln.Addr().String())
	fmt.Println(resp.StatusCode)

	time.Sleep(2 * time.Second)
	resp, _ = http.Get("http://" + ln.Addr().String())
	fmt.Println(resp.StatusCode)
	// Output:
	//200
	//429
	//200
}
