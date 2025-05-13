package middleware_test

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/metrics"
	prom "github.com/gol4ng/httpware/v4/metrics/prometheus"
	"github.com/gol4ng/httpware/v4/middleware"
	"github.com/gol4ng/httpware/v4/mocks"
)

func TestMetrics(t *testing.T) {
	var recorderMock = &mocks.Recorder{}
	var responseWriterMock = &httptest.ResponseRecorder{}
	var req *http.Request
	var responseBody = "fake response"

	// create a fake http request
	req = httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	// create a handler that set http status to 200 and write some response content
	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, req, r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody)) // contentLength=13
	}
	// create metrics httpClient middleware
	stack := httpware.MiddlewareStack(
		middleware.Metrics(recorderMock),
	)

	// assert recorder calls
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), 1).Once()
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), -1).Once()
	recorderMock.On("ObserveHTTPResponseSize", req.Context(), req.URL.String(), int64(len(responseBody)), http.MethodGet, "2xx")
	recorderMock.On("ObserveHTTPRequestDuration", req.Context(), req.URL.String(), mock.AnythingOfType("time.Duration"), http.MethodGet, "2xx")

	// call the middleware stack
	stack.DecorateHandlerFunc(handler).ServeHTTP(responseWriterMock, req)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleMetrics() {
	// Example Need a random ephemeral port (to have a free port)
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

	recorder := prom.NewRecorder(prom.Config{}).RegisterOn(nil)

	// we recommend to use MiddlewareStack to simplify managing all wanted middleware
	// caution middleware order matter
	stack := httpware.MiddlewareStack(
		middleware.Metrics(recorder, metrics.WithIdentifierProvider(func(req *http.Request) string {
			return req.URL.Host + " -> " + req.URL.Path
		})),
	)

	// create a server in order to show it work
	mux := http.NewServeMux()
	mux.Handle("/metrics", stack.DecorateHandler(promhttp.Handler()))
	srv := &http.Server{
		Handler: mux,
	}
	go func() {
		if err := srv.Serve(ln); err != nil {
			panic(err)
		}
	}()

	//Output:
}
