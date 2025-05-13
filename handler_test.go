package httpware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gol4ng/httpware/v4"
	"github.com/stretchr/testify/assert"
)

func TestStatusHandler(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"OK", http.StatusOK},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest("GET", "/", nil)
			recorder := httptest.NewRecorder()

			httpware.StatusHandler(tt.statusCode).ServeHTTP(recorder, request)

			assert.Equal(t, tt.statusCode, recorder.Code)
		})
	}
}

func TestStatusHandlerWillPanic(t *testing.T) {
	request := httptest.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	assert.PanicsWithValue(t, "invalid WriteHeader code -1", func() {
		httpware.StatusHandler(-1).ServeHTTP(recorder, request)
	})
}
