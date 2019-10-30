package httpware_test

import (
	"fmt"
	"github.com/gol4ng/httpware/correlation_id"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gol4ng/httpware"
)

func getMiddleware(t *testing.T, i *int, iBefore int, iAfter int) httpware.Middleware {
	return httpware.Middleware(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				assert.Equal(t, iAfter, *i)
				*i++
			}()
			assert.Equal(t, iBefore, *i)
			*i++
			h.ServeHTTP(w, r)
		})
	})
}

func TestMiddleware_Append(t *testing.T) {
	req := &http.Request{}
	responseWriterMock := &httptest.ResponseRecorder{}
	responseBody := "fake response"

	i := new(int)
	*i = 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, 3, *i)
		*i++
		assert.IsType(t, responseWriterMock, w)
		assert.Equal(t, req, r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	})

	middleware := getMiddleware(t, i, 0, 6)

	middleware.Append(
		// the middleware will be add here
		getMiddleware(t, i, 1, 5),
		getMiddleware(t, i, 2, 4),
	).DecorateHandler(handler).ServeHTTP(responseWriterMock, req)
}

func TestMiddleware_Prepend(t *testing.T) {
	req := &http.Request{}
	responseWriterMock := &httptest.ResponseRecorder{}
	responseBody := "fake response"

	i := new(int)
	*i = 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, 3, *i)
		*i++
		assert.IsType(t, responseWriterMock, w)
		assert.Equal(t, req, r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	})

	middleware := getMiddleware(t, i, 2, 4)

	middleware.Prepend(
		getMiddleware(t, i, 0, 6),
		getMiddleware(t, i, 1, 5),
		// the middleware will be add here
	).DecorateHandler(handler).ServeHTTP(responseWriterMock, req)
}

func TestMiddlewares_DecorateHandler(t *testing.T) {
	req := &http.Request{}
	responseWriterMock := &httptest.ResponseRecorder{}
	responseBody := "fake response"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.IsType(t, responseWriterMock, w)
		assert.Equal(t, req, r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	})

	stack := httpware.MiddlewareStack(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, req, r)
			assert.IsType(t, responseWriterMock, w)
			assert.Equal(t, req, r)

			h.ServeHTTP(w, r)
		})
	})

	stack.DecorateHandler(handler).ServeHTTP(responseWriterMock, req)
}

func TestMiddlewares_DecorateHandlerFunc(t *testing.T) {
	req := &http.Request{}
	responseWriterMock := &httptest.ResponseRecorder{}
	responseBody := "fake response"

	handler := func(w http.ResponseWriter, r *http.Request) {
		assert.IsType(t, responseWriterMock, w)
		assert.Equal(t, req, r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	}

	stack := httpware.MiddlewareStack(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, req, r)
			assert.IsType(t, responseWriterMock, w)
			assert.Equal(t, req, r)

			h.ServeHTTP(w, r)
		})
	})

	stack.DecorateHandlerFunc(handler).ServeHTTP(responseWriterMock, req)
}

func TestMiddlewares_Append(t *testing.T) {
	req := &http.Request{}
	responseWriterMock := &httptest.ResponseRecorder{}
	responseBody := "fake response"

	i := new(int)
	*i = 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, 4, *i)
		*i++
		assert.IsType(t, responseWriterMock, w)
		assert.Equal(t, req, r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	})

	middlewares := httpware.MiddlewareStack(
		getMiddleware(t, i, 0, 8),
		getMiddleware(t, i, 1, 7),
	)

	middlewares.Append(
		// the middlewares will be add here
		getMiddleware(t, i, 2, 6),
		getMiddleware(t, i, 3, 5),
	).DecorateHandler(handler).ServeHTTP(responseWriterMock, req)
}

func TestMiddlewares_Prepend(t *testing.T) {
	req := &http.Request{}
	responseWriterMock := &httptest.ResponseRecorder{}
	responseBody := "fake response"

	i := new(int)
	*i = 0
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, 4, *i)
		*i++
		assert.IsType(t, responseWriterMock, w)
		assert.Equal(t, req, r)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(responseBody))
	})

	middlewares := httpware.MiddlewareStack(
		getMiddleware(t, i, 2, 6),
		getMiddleware(t, i, 3, 5),
	)

	middlewares.Prepend(
		getMiddleware(t, i, 0, 8),
		getMiddleware(t, i, 1, 7),
		// the middlewares will be add here
	).DecorateHandler(handler).ServeHTTP(responseWriterMock, req)
}

// =====================================================================================================================
// =============================== use those examples when declaring an http SERVER ====================================
// =====================================================================================================================

func ExampleMiddlewareStack() {
	// create a middleware that adds a requestId header on each http-server request
	addCustomResponseHeader := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writer.Header().Add("custom-response-header", "wonderful header value")
			h.ServeHTTP(writer, req)
		})
	}
	// create a middleware that logs the response header on each call
	logResponseHeaders := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			fmt.Println("http response headers : ", writer.Header())
			h.ServeHTTP(writer, req)
		})
	}
	// create the middleware stack
	stack := httpware.MiddlewareStack(
		addCustomResponseHeader,
		logResponseHeaders,
	)
	// create a server
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		request.Header.Get(correlation_id.HeaderName)
		writer.WriteHeader(http.StatusOK)
	})

	// apply the middlewares on the server
	// note: this part is normally done on `http.ListenAndServe(":<serverPort>", stack.DecorateHandler(srv))`
	h := stack.DecorateHandler(srv)

	// fake a request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	//Output:
	//http response headers :  map[Custom-Response-Header:[wonderful header value]]
}
