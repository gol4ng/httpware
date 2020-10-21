package middleware_test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware/v3"
	"github.com/gol4ng/httpware/v3/correlation_id"
	"github.com/gol4ng/httpware/v3/middleware"
)

func TestCorrelationId(t *testing.T) {
	correlation_id.DefaultIdGenerator = correlation_id.NewRandomIdGenerator(
		rand.New(correlation_id.NewLockedSource(rand.NewSource(1))),
	)

	var handlerReq *http.Request
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	responseWriter := &httptest.ResponseRecorder{}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// not equal because req.WithContext create another request object
		assert.NotEqual(t, req, r)
		assert.Equal(t, "p1LGIehp1s", r.Header.Get(correlation_id.HeaderName))
		handlerReq = r
	})

	middleware.CorrelationId()(handler).ServeHTTP(responseWriter, req)
	respHeaderValue := responseWriter.Header().Get(correlation_id.HeaderName)
	reqContextValue := handlerReq.Context().Value(correlation_id.HeaderName).(string)
	assert.Equal(t, "p1LGIehp1s", req.Header.Get(correlation_id.HeaderName))
	assert.True(t, len(respHeaderValue) == 10)
	assert.True(t, len(reqContextValue) == 10)
	assert.True(t, respHeaderValue == reqContextValue)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleCorrelationId() {
	port := ":9103"
	// we recommend to use MiddlewareStack to simplify managing all wanted middlewares
	// caution middleware order matters
	stack := httpware.MiddlewareStack(
		middleware.CorrelationId(
			correlation_id.WithHeaderName("my-personal-header-name"),
			correlation_id.WithIdGenerator(func(request *http.Request) string {
				return "my-fixed-request-id"
			}),
		),
	)

	srv := http.NewServeMux()
	go func() {
		if err := http.ListenAndServe(port, stack.DecorateHandler(srv)); err != nil {
			panic(err)
		}
	}()

	resp, err := http.Get("http://localhost" + port)
	if resp != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("%s: %v\n", "my-personal-header-name", resp.Header.Get("my-personal-header-name"))
	}

	//Output: my-personal-header-name: my-fixed-request-id
}
