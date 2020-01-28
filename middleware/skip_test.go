package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/middleware"
	"github.com/stretchr/testify/assert"
)

func TestSkip(t *testing.T) {
	tests := []struct {
		conditionResult  bool
		expectedExecuted bool
	}{
		{
			conditionResult:  true,
			expectedExecuted: false,
		},
		{
			conditionResult:  false,
			expectedExecuted: true,
		},
	}

	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	responseWriter := &httptest.ResponseRecorder{}

	handler := http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		assert.Equal(t, request, req)
		writer.WriteHeader(http.StatusOK)
	})

	executed := false
	dummyMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			executed = true
			next.ServeHTTP(writer, request)
		})
	}

	for k, test := range tests {
		t.Run(fmt.Sprintf("test %d (%v)", k, test), func(t *testing.T) {
			middleware.Skip(func(request *http.Request) bool {
				return test.conditionResult
			}, dummyMiddleware)(handler).ServeHTTP(responseWriter, request)

			assert.Equal(t, test.expectedExecuted, executed)
		})
	}
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleSkip() {
	port := ":9902"

	dummyMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			request.Header.Set("FakeHeader", "this header is set when not /home url")
			next.ServeHTTP(writer, request)
		})
	}
	stack := httpware.MiddlewareStack(
		middleware.Skip(func(request *http.Request) bool {
			return request.URL.Path == "/home"
		}, dummyMiddleware),
	)

	// create a server in order to show it work
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Printf("server receive request %s with request: %s\n", request.URL.Path, request.Header.Get("FakeHeader"))
	})

	go func() {
		if err := http.ListenAndServe(port, stack.DecorateHandler(srv)); err != nil {
			panic(err)
		}
	}()

	_, _ = http.Get("http://localhost" + port + "/")
	_, _ = http.Get("http://localhost" + port + "/home")

	// Output:
	//server receive request / with request: this header is set when not /home url
	//server receive request /home with request:
}
