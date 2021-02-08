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

	"github.com/gol4ng/httpware/v4/mocks"
	"github.com/gol4ng/httpware/v4/tripperware"
)

func TestRequestListener(t *testing.T) {
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
	listenerMock := func(innerReq *http.Request) {
		called = true
		assert.Equal(t, req, innerReq)
	}
	resp2, err := tripperware.RequestListener(listenerMock)(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, resp, resp2)
	assert.True(t, called)
}
