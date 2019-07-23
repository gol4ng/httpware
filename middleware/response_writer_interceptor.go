package middleware

import "net/http"

type responseWriterInterceptor struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (w *responseWriterInterceptor) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *responseWriterInterceptor) Write(p []byte) (int, error) {
	w.bytesWritten += len(p)
	return w.ResponseWriter.Write(p)
}

func NewResponseWriterInterceptor(writer http.ResponseWriter) *responseWriterInterceptor {
	return &responseWriterInterceptor{
		statusCode:     http.StatusServiceUnavailable,
		ResponseWriter: writer,
	}
}
