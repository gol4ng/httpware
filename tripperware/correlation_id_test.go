package tripperware_test

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gol4ng/httpware/v3/correlation_id"
	"github.com/gol4ng/httpware/v3/mocks"
	"github.com/gol4ng/httpware/v3/tripperware"
)

func TestMain(m *testing.M) {
	correlation_id.DefaultIdGenerator = correlation_id.NewRandomIdGenerator(
		rand.New(correlation_id.NewLockedSource(rand.NewSource(1))),
	)
	os.Exit(m.Run())
}

func TestCorrelationId(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	roundTripperMock.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil).Run(func(args mock.Arguments) {
		innerReq := args.Get(0).(*http.Request)
		assert.Len(t, innerReq.Header.Get(correlation_id.HeaderName), 10)
		assert.Equal(t, req.Header.Get(correlation_id.HeaderName), innerReq.Header.Get(correlation_id.HeaderName))
	})

	resp2, err := tripperware.CorrelationId()(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, resp, resp2)
	assert.Equal(t, "p1LGIehp1s", req.Header.Get(correlation_id.HeaderName))
}

func TestCorrelationId_AlreadyInContext(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	req = req.WithContext(context.WithValue(req.Context(), correlation_id.HeaderName, "my_already_exist_correlation_id"))

	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	roundTripperMock.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil).Run(func(args mock.Arguments) {
		innerReq := args.Get(0).(*http.Request)
		assert.Equal(t, req, innerReq)
		assert.Len(t, innerReq.Header.Get(correlation_id.HeaderName), 31)
		assert.Equal(t, req.Header.Get(correlation_id.HeaderName), innerReq.Header.Get(correlation_id.HeaderName))
	})

	resp2, err := tripperware.CorrelationId()(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, resp, resp2)
	assert.Equal(t, "my_already_exist_correlation_id", req.Header.Get(correlation_id.HeaderName))
}

func TestCorrelationIdCustom(t *testing.T) {
	roundTripperMock := &mocks.RoundTripper{}
	req := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	resp := &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	roundTripperMock.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(resp, nil).Run(func(args mock.Arguments) {
		innerReq := args.Get(0).(*http.Request)
		assert.Equal(t, "my_fake_correlation", innerReq.Header.Get(correlation_id.HeaderName))
	})

	resp2, err := tripperware.CorrelationId(
		correlation_id.WithIdGenerator(func(request *http.Request) string {
			return "my_fake_correlation"
		}),
	)(roundTripperMock).RoundTrip(req)
	assert.Nil(t, err)
	assert.Equal(t, resp, resp2)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleCorrelationId() {
	port := ":9001"

	// create http client using the tripperwareStack as RoundTripper
	client := http.Client{
		Transport: tripperware.CorrelationId(
			correlation_id.WithHeaderName("my-personal-header-name"),
			correlation_id.WithIdGenerator(func(request *http.Request) string {
				return "my-fixed-request-id"
			}),
		),
	}

	// create a server in order to show it work
	srv := http.NewServeMux()
	srv.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Println("server receive request with request id:", request.Header.Get("my-personal-header-name"))
	})

	go func() {
		if err := http.ListenAndServe(port, srv); err != nil {
			panic(err)
		}
	}()

	_, _ = client.Get("http://localhost" + port + "/")

	// Output: server receive request with request id: my-fixed-request-id
}
