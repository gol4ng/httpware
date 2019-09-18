package middleware_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/middleware"
	"github.com/gol4ng/httpware/request_id"
)

func TestRequestId(t *testing.T) {
	request_id.DefaultRand = rand.New(request_id.NewLockedSource(rand.NewSource(1)))
	request_id.DefaultIdGenerator = request_id.NewRandomIdGenerator(
		request_id.DefaultRand,
		10,
	)
	var handlerReq *http.Request
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	responseWriter := &httptest.ResponseRecorder{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// not equal because req.WithContext create another request object
		assert.NotEqual(t, req, r)
		assert.NotEqual(t, "", r.Header.Get(request_id.HeaderName))
		handlerReq = r
	})

	middleware.RequestId(request_id.NewConfig())(handler).ServeHTTP(responseWriter, req)
	respHeaderValue := responseWriter.Header().Get(request_id.HeaderName)
	reqContextValue := handlerReq.Context().Value(request_id.HeaderName).(string)
	assert.NotEqual(t, "", req.Header.Get(request_id.HeaderName))
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
		assert.NotEqual(t, "", r.Header.Get(request_id.HeaderName))
		handlerReq = r
	})
	config := request_id.NewConfig()
	config.IdGenerator = func(request *http.Request) string {
		return "my_fake_request_id"
	}

	middleware.RequestId(config)(handler).ServeHTTP(responseWriter, req)
	headerValue := responseWriter.Header().Get(request_id.HeaderName)
	reqContextValue := handlerReq.Context().Value(request_id.HeaderName).(string)
	assert.NotEqual(t, "", req.Header.Get(request_id.HeaderName))
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
	config.IdGenerator = func(request *http.Request) string {
		return "my-fixed-request-id"
	}

	// we recommend to use MiddlewareStack to simplify managing all wanted middlewares
	// caution middleware order matters
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
	fmt.Printf("%s: %v %v\n", config.HeaderName, resp.Header.Get(config.HeaderName), err)

	//Output: my-personal-header-name: my-fixed-request-id <nil>
}
