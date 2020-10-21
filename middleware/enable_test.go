package middleware_test

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gol4ng/httpware/v3"
	"github.com/gol4ng/httpware/v3/middleware"
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
	// Example Need a random ephemeral port (to have a free port)
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		log.Fatal(err)
	}

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
	srv := &http.Server{
		Handler: stack.DecorateHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			fmt.Println("server receive request with request:", request.Header.Get("FakeHeader"))
		}),
	}
	go func() {
		if err := srv.Serve(ln); err != nil {
			panic(err)
		}
	}()

	_, _ = http.Get("http://" + ln.Addr().String())

	// Output:
	//server receive request with request: this header is set when not /home url
}
