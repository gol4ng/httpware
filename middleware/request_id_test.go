package middleware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/middleware"
	"github.com/gol4ng/httpware/request_id"
)

func TestRequestId(t *testing.T) {
	var handlerReq *http.Request
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	responseWriter := &httptest.ResponseRecorder{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// not equal because req.WithContext create another request object
		assert.NotEqual(t, req, r)
		handlerReq = r
	})

	middleware.RequestId(request_id.NewConfig())(handler).ServeHTTP(responseWriter, req)
	respHeaderValue := responseWriter.Header().Get(request_id.HeaderName)
	reqContextValue := handlerReq.Context().Value(request_id.HeaderName).(string)
	assert.True(t, len(respHeaderValue) == 10)
	assert.True(t, len(reqContextValue) == 10)
	assert.True(t, respHeaderValue == reqContextValue)
}

func TestRequestIdCustom(t *testing.T) {
	var handlerReq *http.Request
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	responseWriter := &httptest.ResponseRecorder{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// not equal because req.WithContext create another request object
		assert.NotEqual(t, req, r)
		handlerReq = r
	})
	config := request_id.NewConfig()
	config.GenerateId = func(request *http.Request) string {
		return "my_fake_request_id"
	}

	middleware.RequestId(config)(handler).ServeHTTP(responseWriter, req)
	headerValue := responseWriter.Header().Get(request_id.HeaderName)
	reqContextValue := handlerReq.Context().Value(request_id.HeaderName).(string)
	assert.Equal(t, "my_fake_request_id", headerValue)
	assert.Equal(t, "my_fake_request_id", reqContextValue)
	assert.True(t, headerValue == reqContextValue)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleRequestId() {
	port := ":5001"

	config := request_id.NewConfig()
	// you can override default header name
	config.HeaderName = "my-personal-header-name"
	// you can override default id generator
	config.GenerateId = func(request *http.Request) string {
		return "my-fixed-request-id"
	}

	// we recommend to use MiddlewareStack to simplify managing all wanted middleware
	// caution middleware order matter
	stack := httpware.MiddlewareStack(
		middleware.RequestId(config),
	)

	srv := http.NewServeMux()
	go func() {
		if err := http.ListenAndServe(port, stack.DecorateHandler(srv)); err != nil {
			panic(err)
		}
	}()

	resp, err := http.Get("http://localhost" + port)
	fmt.Printf("%v %v\n", resp.Header.Get(config.HeaderName), err)

	//Output: my-fixed-request-id <nil>
}
