package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware"
	"github.com/gol4ng/httpware/interceptor"
)

func Interceptor(callbackBeforeFunc func(*ResponseWriterInterceptor, *http.Request), callbackAfterFunc func(*ResponseWriterInterceptor, *http.Request)) httpware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			writerInterceptor := NewResponseWriterInterceptor(writer)

			req.Body = interceptor.NewCopyReadCloser(req.Body)
			callbackBeforeFunc(writerInterceptor, req)
			defer func() {
				callbackAfterFunc(writerInterceptor, req)
			}()

			next.ServeHTTP(writerInterceptor, req)
		})
	}
}
