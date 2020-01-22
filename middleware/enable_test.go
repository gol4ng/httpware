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

func TestEnable(t *testing.T) {
	tests := []struct {
		enable           bool
		expectedExecuted bool
	}{
		{
			enable:           true,
			expectedExecuted: true,
		},
		{
			enable:           false,
			expectedExecuted: false,
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
		executed = false
		t.Run(fmt.Sprintf("test %d (%v)", k, test), func(t *testing.T) {
			middleware.Enable(test.enable, dummyMiddleware)(handler).ServeHTTP(responseWriter, request)

			assert.Equal(t, test.expectedExecuted, executed)
		})
	}
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleEnable() {
	port := ":9104"

	enableDummyMiddleware := true // or false
	dummyMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			request.Header.Set("FakeHeader", "this header is set when not /home url")
			next.ServeHTTP(writer, request)
		})
	}
	stack := httpware.MiddlewareStack(
		middleware.Enable(enableDummyMiddleware, dummyMiddleware),
	)

	// create a server in order to show it work
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("server receive request with request:", request.Header.Get("FakeHeader"))
	})

	go func() {
		if err := http.ListenAndServe(port, stack.DecorateHandler(srv)); err != nil {
			panic(err)
		}
	}()

	_, _ = http.Get("http://localhost" + port + "/")

	// Output:
	//server receive request with request: this header is set when not /home url
}
