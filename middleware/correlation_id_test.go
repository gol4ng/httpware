package middleware_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/correlation_id"
	"github.com/gol4ng/httpware/middleware"
)

func TestRequestId(t *testing.T) {
	defaultRand := rand.New(correlation_id.NewLockedSource(rand.NewSource(1)))
	correlation_id.DefaultIdGenerator = correlation_id.NewRandomIdGenerator(defaultRand)

	var handlerReq *http.Request
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	responseWriter := &httptest.ResponseRecorder{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// not equal because req.WithContext create another request object
		assert.NotEqual(t, req, r)
		assert.Equal(t, "p1LGIehp1s", r.Header.Get(correlation_id.HeaderName))
		handlerReq = r
	})

	middleware.RequestId(correlation_id.NewConfig())(handler).ServeHTTP(responseWriter, req)
	respHeaderValue := responseWriter.Header().Get(correlation_id.HeaderName)
	reqContextValue := handlerReq.Context().Value(correlation_id.HeaderName).(string)
	assert.Equal(t, "p1LGIehp1s", req.Header.Get(correlation_id.HeaderName))
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
		assert.Equal(t, "my_fake_correlation", r.Header.Get(correlation_id.HeaderName))
		handlerReq = r
	})
	config := correlation_id.NewConfig()
	config.IdGenerator = func(request *http.Request) string {
		return "my_fake_correlation"
	}

	middleware.RequestId(config)(handler).ServeHTTP(responseWriter, req)
	headerValue := responseWriter.Header().Get(correlation_id.HeaderName)
	reqContextValue := handlerReq.Context().Value(correlation_id.HeaderName).(string)
	assert.NotEqual(t, "", req.Header.Get(correlation_id.HeaderName))
	assert.Equal(t, "my_fake_correlation", headerValue)
	assert.Equal(t, "my_fake_correlation", reqContextValue)
	assert.True(t, headerValue == reqContextValue)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleRequestId() {
	port := ":5001"

	config := correlation_id.NewConfig()
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
