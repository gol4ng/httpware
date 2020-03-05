package tripperware

import (
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/auth"
)

func AuthenticationForwarder() httpware.Tripperware {
	return func(next http.RoundTripper) http.RoundTripper {
		return httpware.RoundTripFunc(func(req *http.Request) (*http.Response, error) {
			auth.AddHeader(req)(auth.CredentialFromContext(req.Context()))
			return next.RoundTrip(req)
		})
	}
}
