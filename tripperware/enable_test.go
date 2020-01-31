package tripperware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/mocks"
	"github.com/gol4ng/httpware/v2/tripperware"
	"github.com/stretchr/testify/assert"
)

func TestEnable(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)

	tests := []struct {
		enable           bool
		expectedExecuted bool
	}{
		{
			enable:           true,
			expectedExecuted: true,
		},
		{
			enable:           false,
			expectedExecuted: false,
		},
	}

	executed := false
	dummyTripperware := func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			executed = true
			return next.RoundTrip(request)
		})
	}

	for k, test := range tests {
		executed = false
		t.Run(fmt.Sprintf("test %d (%v)", k, test), func(t *testing.T) {
			resp2, err := tripperware.Enable(test.enable, dummyTripperware, )(roundTripperMock).RoundTrip(req)

			assert.Nil(t, err)
			assert.Equal(t, resp, resp2)
			assert.Equal(t, test.expectedExecuted, executed)
		})
	}
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleEnable() {
	port := ":9003"

	enableDummyTripperware := true //false
	dummyTripperware := func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			request.Header.Set("FakeHeader", "this header is set when not /home url")
			return next.RoundTrip(request)
		})
	}

	// create http client using the tripperwareStack as RoundTripper
	client := http.Client{
		Transport: tripperware.Enable(enableDummyTripperware, dummyTripperware),
	}

	// create a server in order to show it work
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("server receive request with request:", request.Header.Get("FakeHeader"))
	})

	go func() {
		if err := http.ListenAndServe(port, srv); err != nil {
			panic(err)
		}
	}()

	_, _ = client.Get("http://localhost" + port + "/")

	// Output:
	//server receive request with request: this header is set when not /home url
}