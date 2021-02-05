package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/exporter"
)

func CurlExporter(printer func(*exporter.Cmd, error)) httpware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			printer(exporter.GetCurlCommand(request))
			next.ServeHTTP(writer, request)
		})
	}
}
