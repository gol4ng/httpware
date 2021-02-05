package middleware_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gol4ng/httpware/v4/exporter"
	"github.com/gol4ng/httpware/v4/middleware"
	"github.com/stretchr/testify/assert"
)

func TestCurlExporter(t *testing.T) {
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
	mockExporter := func(cmd *exporter.Cmd, err error) {
		called = true
		assert.Equal(t, "curl -X 'GET' 'http://fake-addr' -d 'mykey=myvalue'", cmd.String())
		assert.Nil(t, err)
	}

	middleware.CurlExporter(mockExporter)(handler).ServeHTTP(responseWriter, req)
	assert.True(t, called)
	assert.True(t, handlerCalled)
}
