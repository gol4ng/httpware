package middleware

import (
	"net/http"
)

type ResponseWriterInterceptor struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (w *ResponseWriterInterceptor) WriteHeader(statusCode int) {
	w.StatusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *ResponseWriterInterceptor) Write(p []byte) (int, error) {
	w.Body = append(w.Body, p...)
	return w.ResponseWriter.Write(p)
}

func NewResponseWriterInterceptor(writer http.ResponseWriter) *ResponseWriterInterceptor {
	return &ResponseWriterInterceptor{
		StatusCode:     http.StatusServiceUnavailable,
		ResponseWriter: writer,
	}
}
