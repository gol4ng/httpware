package middleware_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/middleware"
	"github.com/gol4ng/httpware/v2/mocks"
	"github.com/gol4ng/httpware/v2/rate_limit"
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
	limiter := rate_limit.NewTokenBucket(1*time.Second, 1)
	defer limiter.Stop()

	port := ":9105"
	// we recommend to use MiddlewareStack to simplify managing all wanted middlewares
	// caution middleware order matters
	stack := httpware.MiddlewareStack(
		middleware.RateLimit(limiter),
	)

	srv := http.NewServeMux()
	srv.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {})
	go func() {
		if err := http.ListenAndServe(port, stack.DecorateHandler(srv)); err != nil {
			panic(err)
		}
	}()

	resp, _ := http.Get("http://localhost" + port)
	fmt.Println(resp.StatusCode)

	resp, _ = http.Get("http://localhost" + port)
	fmt.Println(resp.StatusCode)

	time.Sleep(2 * time.Second)
	resp, _ = http.Get("http://localhost" + port)
	fmt.Println(resp.StatusCode)
	// Output:
	//200
	//429
	//200
}
