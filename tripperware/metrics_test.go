package tripperware_test

import (
	"github.com/agiledragon/gomonkey/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

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
	var requestTimeDuration = 10 * time.Millisecond
	var resp = &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: 30,
	}
	var baseTime = time.Unix(513216000, 0)

	req = httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	// mock roundTripper calls
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)
	// assert recorder calls
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), 1).Once()
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), -1).Once()
	recorderMock.On("ObserveHTTPRequestDuration", req.Context(), req.URL.String(), requestTimeDuration, http.MethodGet, "2xx")
	recorderMock.On("ObserveHTTPResponseSize", req.Context(), req.URL.String(), resp.ContentLength, http.MethodGet, "2xx")
	// mock time.Now method in order to return always the same time whenever the test is launched
	patch := gomonkey.NewPatches()
	patch.ApplyFunc(time.Now, func() time.Time { return baseTime })
	patch.ApplyFunc(time.Since, func(since time.Time) time.Duration {
		assert.Equal(t, baseTime, since)
		return requestTimeDuration
	})
	defer patch.Reset()

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
	var requestTimeDuration = 10 * time.Millisecond
	var resp = &http.Response{
		Status:        "OK",
		StatusCode:    http.StatusOK,
		ContentLength: -1,
	}
	expectedContentLength := int64(0)
	var baseTime = time.Unix(513216000, 0)

	req = httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	// mock roundTripper calls
	roundTripperMock.On("RoundTrip", req).Return(resp, nil)
	// assert recorder calls
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), 1).Once()
	recorderMock.On("AddInflightRequests", req.Context(), req.URL.String(), -1).Once()
	recorderMock.On("ObserveHTTPRequestDuration", req.Context(), req.URL.String(), requestTimeDuration, http.MethodGet, "2xx")
	recorderMock.On("ObserveHTTPResponseSize", req.Context(), req.URL.String(), expectedContentLength, http.MethodGet, "2xx")
	// mock time.Now method in order to return always the same time whenever the test is launched
	patch := gomonkey.NewPatches()
	patch.ApplyFunc(time.Now, func() time.Time { return baseTime })
	patch.ApplyFunc(time.Since, func(since time.Time) time.Duration {
		assert.Equal(t, baseTime, since)
		return requestTimeDuration
	})
	defer patch.Reset()

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
