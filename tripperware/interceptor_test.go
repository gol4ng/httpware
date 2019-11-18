package tripperware_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/mocks"
	"github.com/gol4ng/httpware/v2/tripperware"
)

func TestInterceptor(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodPost, "http://fake-addr", bytes.NewBufferString("my_fake_body"))
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	roundTripperMock.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil).Run(func(args mock.Arguments) {
		innerReq := args.Get(0).(*http.Request)
		reqData, err := ioutil.ReadAll(innerReq.Body)
		assert.Nil(t, err)
		assert.Equal(t, "my_fake_body", string(reqData))
	})

	resp2, err := tripperware.Interceptor(
		tripperware.WithBefore(func(request *http.Request) {
			reqData, err := ioutil.ReadAll(request.Body)
			assert.Nil(t, err)
			assert.Equal(t, "my_fake_body", string(reqData))
		}),
		tripperware.WithAfter(func(response *http.Response, request *http.Request) {

		}),
	)(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, resp, resp2)

	reqData, err := ioutil.ReadAll(req.Body)
	assert.Nil(t, err)
	assert.Equal(t, "my_fake_body", string(reqData))
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleInterceptor() {
	// we recommend to use TripperwareStack to simplify managing all wanted tripperware
	// caution tripperware order matter
	stack := httpware.TripperwareStack(
		tripperware.Interceptor(
			tripperware.WithBefore(func(request *http.Request) {
				reqData, err := ioutil.ReadAll(request.Body)
				fmt.Println("before callback", string(reqData), err)
			}),
			tripperware.WithAfter(func(response *http.Response, request *http.Request) {
				reqData, err := ioutil.ReadAll(request.Body)
				fmt.Println("after callback", string(reqData), err)
			}),
		),
	)

	// create http client using the tripperwareStack as RoundTripper
	client := http.Client{
		Transport: stack,
	}

	_, _ = client.Post("fake-address.foo", "plain/text", bytes.NewBufferString("my_fake_body"))

	//Output:
	//before callback my_fake_body <nil>
	//after callback my_fake_body <nil>
}
