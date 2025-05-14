package httpware

import (
	"net/http"
)

var (
	NopHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
)

func StatusHandler(code int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
	})
}
