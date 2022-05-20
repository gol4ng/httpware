package middleware

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/felixge/httpsnoop"
	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/metrics"
)

func Metrics(recorder metrics.Recorder, options ... metrics.Option) httpware.Middleware {
	config := metrics.NewConfig(recorder, options...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			handlerName := config.IdentifierProvider(req)
			if config.MeasureInflightRequests {
				config.Recorder.AddInflightRequests(req.Context(), handlerName, 1)
				defer config.Recorder.AddInflightRequests(req.Context(), handlerName, -1)
			}

			httpMetrics := httpsnoop.Metrics{}

			defer func() {
				code := strconv.Itoa(httpMetrics.Code)
				if !config.SplitStatus {
					code = fmt.Sprintf("%dxx", httpMetrics.Code/100)
				}

				config.Recorder.ObserveHTTPRequestDuration(req.Context(), handlerName, httpMetrics.Duration, req.Method, code)

				if config.ObserveResponseSize {
					config.Recorder.ObserveHTTPResponseSize(req.Context(), handlerName, httpMetrics.Written, req.Method, code)
				}
			}()

			httpMetrics.CaptureMetrics(writer, func(writer http.ResponseWriter) {
				next.ServeHTTP(writer, req)
			})
		})
	}
}
