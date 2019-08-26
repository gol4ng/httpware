package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/metrics"
	"github.com/gol4ng/httpware/metrics/prometheus"
	"github.com/gol4ng/httpware/middleware"
	"github.com/gol4ng/httpware/mocks"
)

func TestMetrics(t *testing.T) {
	var recorderMock = &mocks.Recorder{}
	var responseWriterMock = &httptest.ResponseRecorder{}
	var req *http.Request
	var requestTimeDuration = 10*time.Millisecond
	var baseTime = time.Unix(513216000, 0)
	var responseBody = "fake response"

	// create fake http request
	req = httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	// create handler that set http status to 200 and write some response content
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.IsType(t, middleware.NewResponseWriterInterceptor(nil), w)
		assert.Equal(t, req, r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody)) // contentLength=13
	}
	// create metrics httpClient middleware
	stack := httpware.MiddlewareStack(
		middleware.Metrics(metrics.NewConfig(recorderMock)),
	)

	// assert recorder calls
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), 1).Once()
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), -1).Once()
	recorderMock.On("ObserveHTTPRequestDuration", req.Context(), req.URL.String(), requestTimeDuration, http.MethodGet, "2xx")
	recorderMock.On("ObserveHTTPResponseSize", req.Context(), req.URL.String(), int64(len(responseBody)), http.MethodGet, "2xx")
	// mock time.Now method in order to return always the same time whenever the test is launched
	monkey.Patch(time.Now, func() time.Time { return baseTime })
	monkey.Patch(time.Since, func(since time.Time) time.Duration {
		assert.Equal(t, baseTime, since)
		return requestTimeDuration
	})
	defer monkey.UnpatchAll()

	// call the middleware stack
	stack.DecorateHandlerFunc(handler).ServeHTTP(responseWriterMock, req)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleMetrics() {
	recorder := prometheus.
		NewRecorder(prometheus.Config{}).
		RegisterOn(nil)

	conf := metrics.NewConfig(recorder)
	conf.IdentifierProvider = func(req *http.Request) string {
		return req.URL.Host+" -> "+req.URL.Path
	}
	stack := httpware.MiddlewareStack(
		middleware.Metrics(conf),
	)

	srv := http.NewServeMux()
	stack.DecorateHandler(srv)

	go func() {
		if err := http.ListenAndServe(":3000", stack.DecorateHandler(srv)); err != nil {
			panic(err)
		}
	}()

	//Output:
}
