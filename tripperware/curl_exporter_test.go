package tripperware_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gol4ng/httpware/v4/exporter"
	"github.com/gol4ng/httpware/v4/mocks"
	"github.com/gol4ng/httpware/v4/tripperware"
)

func TestCurlExporter(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", ioutil.NopCloser(strings.NewReader(url.Values{
		"mykey": {"myvalue"},
	}.Encode())))

	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	roundTripperMock.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil)

	called := false
	mockExporter := func(cmd *exporter.Cmd, err error) {
		called = true
		assert.Equal(t, "curl -X 'GET' 'http://fake-addr' -d 'mykey=myvalue'", cmd.String())
		assert.Nil(t, err)
	}
	resp2, err := tripperware.CurlExporter(mockExporter)(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, resp, resp2)
	assert.True(t, called)
}
