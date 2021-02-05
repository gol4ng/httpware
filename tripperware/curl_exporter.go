package tripperware

import (
	"net/http"

	"github.com/gol4ng/httpware/v4"
	"github.com/gol4ng/httpware/v4/exporter"
)

func CurlExporter(export func(*exporter.Cmd, error)) httpware.Tripperware {
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(request *http.Request) (*http.Response, error) {
			export(exporter.GetCurlCommand(request))
			return next.RoundTrip(request)
		})
	}
}
