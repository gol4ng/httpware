package tripperware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gol4ng/httpware/v3"
	"github.com/gol4ng/httpware/v3/mocks"
	"github.com/gol4ng/httpware/v3/tripperware"
	"github.com/stretchr/testify/assert"
)

func TestSkip(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)

	tests := []struct {
		conditionResult  bool
		expectedExecuted bool
	}{
		{
			conditionResult:  true,
			expectedExecuted: false,
		},
		{
			conditionResult:  false,
			expectedExecuted: true,
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
			resp2, err := tripperware.Skip(func(request *http.Request) bool {
				return test.conditionResult
			}, dummyTripperware)(roundTripperMock).RoundTrip(req)

			assert.Nil(t, err)
			assert.Equal(t, resp, resp2)
			assert.Equal(t, test.expectedExecuted, executed)
		})
	}
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleSkip() {
	port := ":9002"

	dummyTripperware := func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			request.Header.Set("FakeHeader", "this header is set when not /home url")
			return next.RoundTrip(request)
		})
	}

	// create http client using the tripperwareStack as RoundTripper
	client := http.Client{
		Transport: tripperware.Skip(func(request *http.Request) bool {
			return request.URL.Path == "/home"
		}, dummyTripperware),
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
	_, _ = client.Get("http://localhost" + port + "/home")

	// Output:
	//server receive request with request: this header is set when not /home url
	//server receive request with request:
}
