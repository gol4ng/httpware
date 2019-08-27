package httpware_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/gol4ng/httpware"
)

// =====================================================================================================================
// =============================== use those examples when declaring an http SERVER ====================================
// =====================================================================================================================

func ExampleMiddlewareStack() {
	// create a middleware that adds a requestId header on each http-server request
	addCustomResponseHeader := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writer.Header().Add("custom-response-header", "wonderful header value")
			h.ServeHTTP(writer, req)
		})
	}
	// create a middleware that logs the response header on each call
	logResponseHeaders := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			fmt.Println("http response headers : ", writer.Header())
			h.ServeHTTP(writer, req)
		})
	}
	// create the middleware stack
	stack := httpware.MiddlewareStack(
		logResponseHeaders,
		addCustomResponseHeader,
	)
	// create a server
	srv := http.NewServeMux()
	// apply the middlewares on the server
	// note: this part is normally done on `http.ListenAndServe(":<serverPort>", stack.DecorateHandler(srv))`
	h := stack.DecorateHandler(srv)

	// fake a request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	//Output:
	//http response headers :  map[Custom-Response-Header:[wonderful header value]]
}
