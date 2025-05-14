package tripperware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/mock"

	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/metrics"
	"github.com/gol4ng/httpware/v4/metrics/prometheus"
	"github.com/gol4ng/httpware/v4/mocks"
	"github.com/gol4ng/httpware/v4/tripperware"
)

func TestMetrics(t *testing.T) {
	var recorderMock = &mocks.Recorder{}
	var roundTripperMock = &mocks.RoundTripper{}
	var req *http.Request
	var resp = &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}

	req = httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	// mock roundTripper calls
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)
	// assert recorder calls
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), 1).Once()
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), -1).Once()
	recorderMock.On("ObserveHTTPRequestDuration", req.Context(), req.URL.String(), mock.AnythingOfType("time.Duration"), http.MethodGet, "2xx")
	recorderMock.On("ObserveHTTPResponseSize", req.Context(), req.URL.String(), resp.ContentLength, http.MethodGet, "2xx")

	// create metrics httpClient middleware
	stack := httpware.TripperwareStack(
		tripperware.Metrics(recorderMock),
	)
	_, _ = stack.DecorateRoundTripper(roundTripperMock).RoundTrip(req)
}

func TestMetricsContentLengthUnknown(t *testing.T) {
	var recorderMock = &mocks.Recorder{}
	var roundTripperMock = &mocks.RoundTripper{}
	var req *http.Request
	var resp = &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: -1,
	}
	expectedContentLength := int64(0)

	req = httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	// mock roundTripper calls
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)
	// assert recorder calls
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), 1).Once()
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), -1).Once()
	recorderMock.On("ObserveHTTPRequestDuration", req.Context(), req.URL.String(), mock.AnythingOfType("time.Duration"), http.MethodGet, "2xx")
	recorderMock.On("ObserveHTTPResponseSize", req.Context(), req.URL.String(), expectedContentLength, http.MethodGet, "2xx")

	// create metrics httpClient middleware
	stack := httpware.TripperwareStack(
		tripperware.Metrics(recorderMock),
	)
	_, _ = stack.DecorateRoundTripper(roundTripperMock).RoundTrip(req)
}

// =====================================================================================================================
// ========================================= EXAMPLES ==================================================================
// =====================================================================================================================

func ExampleMetrics() {
	recorder := prometheus.NewRecorder(prometheus.Config{}).RegisterOn(nil)

	// we recommend to use MiddlewareStack to simplify managing all wanted middleware
	// caution middleware order matter
	stack := httpware.TripperwareStack(
		tripperware.Metrics(recorder, metrics.WithIdentifierProvider(func(req *http.Request) string {
			return req.URL.Host + " -> " + req.URL.Path
		})),
	)

	// create http client using the tripperwareStack as RoundTripper
	client := http.Client{
		Transport: stack,
	}

	_, _ = client.Get("fake-address.foo")

	//Output:
}
