package tripperware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gol4ng/httpware/mocks"
	"github.com/gol4ng/httpware/request_id"
	"github.com/gol4ng/httpware/tripperware"
)

func TestRequestId(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	roundTripperMock.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil).Run(func(args mock.Arguments) {
		innerReq := args.Get(0).(*http.Request)
		assert.True(t, len(innerReq.Header.Get(request_id.HeaderName)) == 10)
	})

	resp2, err := tripperware.RequestId(request_id.NewConfig())(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, resp, resp2)
}

func TestRequestIdCustom(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	roundTripperMock.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil).Run(func(args mock.Arguments) {
		innerReq := args.Get(0).(*http.Request)
		assert.Equal(t, "my_fake_request_id", innerReq.Header.Get(request_id.HeaderName))
	})

	config := request_id.NewConfig()
	config.GenerateId = func(request *http.Request) string {
		return "my_fake_request_id"
	}

	resp2, err := tripperware.RequestId(config)(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, resp, resp2)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleRequestId() {
}
