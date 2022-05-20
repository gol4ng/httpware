package middleware

import (
	"github.com/felixge/httpsnoop"
	"net/http"
)

type ResponseWriterInterceptor struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func NewResponseWriterInterceptor(writer http.ResponseWriter) *ResponseWriterInterceptor {
	rw := &ResponseWriterInterceptor{
		StatusCode: http.StatusOK,
	}
	wrapper := httpsnoop.Wrap(writer, httpsnoop.Hooks{
		WriteHeader: func(next httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return func(code int) {
				next(code)
				rw.StatusCode = code
			}
		},
		Write: func(next httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return func(p []byte) (int, error) {
				n, err := next(p)
				rw.Body = append(rw.Body, p...)
				return n, err
			}
		},
	})
	rw.ResponseWriter = wrapper
	return rw
}
