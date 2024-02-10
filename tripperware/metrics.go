package tripperware

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/httptrace"
	"strconv"
	"time"

	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/metrics"
)

func Metrics(recorder metrics.Recorder, options ... metrics.Option) httpware.Tripperware {
	config := metrics.NewConfig(recorder, options...)
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
			handlerName := config.IdentifierProvider(req)
			if config.MeasureInflightRequests {
				config.Recorder.AddInflightRequests(req.Context(), handlerName, 1)
				defer config.Recorder.AddInflightRequests(req.Context(), handlerName, -1)
			}

			start := time.Now()
			defer func() {
				statusCode := http.StatusServiceUnavailable
				contentLength := int64(0)
				if resp != nil {
					statusCode = resp.StatusCode
					// ContentLength records the length of the associated content. The
					// value -1 indicates that the length is unknown. Unless Request.Method
					// is "HEAD", values >= 0 indicate that the given number of bytes may
					// be read from Body.
					if resp.ContentLength > 0 {
						contentLength = resp.ContentLength
					}
				}
				code := strconv.Itoa(statusCode)
				if !config.SplitStatus {
					// modify status to only take first digit into account (201 -> 200; 404 -> 400; ...)
					code = fmt.Sprintf("%dxx", statusCode/100)
				}

				config.Recorder.ObserveHTTPRequestDuration(req.Context(), handlerName, time.Since(start), req.Method, code)

				if config.ObserveResponseSize {
					config.Recorder.ObserveHTTPResponseSize(req.Context(), handlerName, contentLength, req.Method, code)
				}
			}()

			return next.RoundTrip(req)
		})
	}
}

func MetricsTrace(it *metrics.InstrumentTrace) httpware.Tripperware {
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (resp *http.Response, err error) {
			start := time.Now()
			ctx := req.Context()

			trace := &httptrace.ClientTrace{
				GotConn: func(_ httptrace.GotConnInfo) {
					if it.GotConn != nil {
						it.GotConn(ctx, time.Since(start))
					}
				},
				PutIdleConn: func(err error) {
					if err != nil {
						return
					}
					if it.PutIdleConn != nil {
						it.PutIdleConn(ctx, time.Since(start))
					}
				},
				DNSStart: func(_ httptrace.DNSStartInfo) {
					if it.DNSStart != nil {
						it.DNSStart(ctx, time.Since(start))
					}
				},
				DNSDone: func(_ httptrace.DNSDoneInfo) {
					if it.DNSDone != nil {
						it.DNSDone(ctx, time.Since(start))
					}
				},
				ConnectStart: func(_, _ string) {
					if it.ConnectStart != nil {
						it.ConnectStart(ctx, time.Since(start))
					}
				},
				ConnectDone: func(_, _ string, err error) {
					if err != nil {
						return
					}
					if it.ConnectDone != nil {
						it.ConnectDone(ctx, time.Since(start))
					}
				},
				GotFirstResponseByte: func() {
					if it.GotFirstResponseByte != nil {
						it.GotFirstResponseByte(ctx, time.Since(start))
					}
				},
				Got100Continue: func() {
					if it.Got100Continue != nil {
						it.Got100Continue(ctx, time.Since(start))
					}
				},
				TLSHandshakeStart: func() {
					if it.TLSHandshakeStart != nil {
						it.TLSHandshakeStart(ctx, time.Since(start))
					}
				},
				TLSHandshakeDone: func(_ tls.ConnectionState, err error) {
					if err != nil {
						return
					}
					if it.TLSHandshakeDone != nil {
						it.TLSHandshakeDone(ctx, time.Since(start))
					}
				},
				WroteHeaders: func() {
					if it.WroteHeaders != nil {
						it.WroteHeaders(ctx, time.Since(start))
					}
				},
				Wait100Continue: func() {
					if it.Wait100Continue != nil {
						it.Wait100Continue(ctx, time.Since(start))
					}
				},
				WroteRequest: func(_ httptrace.WroteRequestInfo) {
					if it.WroteRequest != nil {
						it.WroteRequest(ctx, time.Since(start))
					}
				},
			}

			return next.RoundTrip(req.WithContext(httptrace.WithClientTrace(req.Context(), trace)))
		})
	}
}
