package middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gol4ng/httpware/v3"
	"github.com/gol4ng/httpware/v3/metrics"
)

func Metrics(recorder metrics.Recorder, options ... metrics.Option) httpware.Middleware {
	config := metrics.NewConfig(recorder, options...)
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writerInterceptor := NewResponseWriterInterceptor(writer)
			handlerName := config.IdentifierProvider(req)
			if config.MeasureInflightRequests {
				config.Recorder.AddInflightRequests(req.Context(), handlerName, 1)
				defer config.Recorder.AddInflightRequests(req.Context(), handlerName, -1)
			}

			start := time.Now()
			defer func() {
				code := strconv.Itoa(writerInterceptor.StatusCode)
				if !config.SplitStatus {
					code = fmt.Sprintf("%dxx", writerInterceptor.StatusCode/100)
				}

				config.Recorder.ObserveHTTPRequestDuration(req.Context(), handlerName, time.Since(start), req.Method, code)

				if config.ObserveResponseSize {
					config.Recorder.ObserveHTTPResponseSize(req.Context(), handlerName, int64(len(writerInterceptor.Body)), req.Method, code)
				}
			}()

			next.ServeHTTP(writerInterceptor, req)
		})
	}
}
