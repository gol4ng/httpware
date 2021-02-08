package middleware_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gol4ng/httpware/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestRequestListener(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", ioutil.NopCloser(strings.NewReader(url.Values{
		"mykey": {"myvalue"},
	}.Encode())))
	responseWriter := &httptest.ResponseRecorder{}

	handlerCalled := false
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, req, r)
		handlerCalled = true
	})

	called := false
	mockExporter := func(innerReq *http.Request) {
		called = true
		assert.Equal(t, req, innerReq)
	}

	middleware.RequestListener(mockExporter)(handler).ServeHTTP(responseWriter, req)
	assert.True(t, called)
	assert.True(t, handlerCalled)
}
