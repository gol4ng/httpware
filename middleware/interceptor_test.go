package middleware_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/middleware"
	"github.com/stretchr/testify/assert"
)

func TestInterceptor(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/foo", bytes.NewReader([]byte("bar")))
	req.Header.Add("X-Interceptor-Request-Header", "interceptor")

	responseWriter := &httptest.ResponseRecorder{}
	stack := httpware.MiddlewareStack(
		middleware.Interceptor(
			func(responseWriterInterceptor *middleware.ResponseWriterInterceptor, req *http.Request) {
				buf := new(bytes.Buffer)
				_, err := buf.ReadFrom(req.Body)
				assert.NoError(t, err)
				assert.Equal(t, "bar", buf.String())

				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/foo", req.URL.String())

				req.Header.Add("X-Interceptor-Request-Header", "interceptor")
				responseWriterInterceptor.Header().Add("X-Interceptor-Response-Header1", "interceptor1")
			},
			func(responseWriterInterceptor *middleware.ResponseWriterInterceptor, req *http.Request) {
				assert.Equal(t, http.MethodGet, req.Method)
				assert.Equal(t, "/foo", req.URL.String())
				assert.Equal(t, "interceptor", req.Header.Get("X-Interceptor-Request-Header"))

				assert.Equal(t, http.StatusAlreadyReported, responseWriterInterceptor.StatusCode)
				assert.Equal(t, "foo bar", string(responseWriterInterceptor.Body))

				assert.Equal(t, "interceptor1", responseWriterInterceptor.Header().Get("X-Interceptor-Response-Header1"))
				assert.Equal(t, "interceptor2", responseWriterInterceptor.Header().Get("X-Interceptor-Response-Header2"))

				responseWriterInterceptor.Header().Add("X-Interceptor-Response-Header3", "interceptor3")
			},
		),
	)

	stack.DecorateHandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(req.Body)
		assert.NoError(t, err)
		assert.Equal(t, "bar", buf.String())
		rw.WriteHeader(http.StatusAlreadyReported)

		_, err = rw.Write([]byte("foo bar"))
		assert.NoError(t, err)
		assert.Equal(t, "interceptor1", rw.Header().Get("X-Interceptor-Response-Header1"))

		rw.Header().Add("X-Interceptor-Response-Header2", "interceptor2")
	}).ServeHTTP(responseWriter, req)

	assert.Equal(t, "interceptor1", responseWriter.Header().Get("X-Interceptor-Response-Header1"))
	assert.Equal(t, "interceptor2", responseWriter.Header().Get("X-Interceptor-Response-Header2"))
	assert.Equal(t, "interceptor3", responseWriter.Header().Get("X-Interceptor-Response-Header3"))
}
