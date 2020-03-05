package middleware_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gol4ng/httpware/v2/auth"
	"github.com/gol4ng/httpware/v2/middleware"
	"github.com/stretchr/testify/assert"
)

func TestAuthentication_hydrate_header(t *testing.T) {
	tests := []struct {
		authorizationHeader  string
		xAuthorizationHeader string
		expectedCredential   string
	}{
		{
			authorizationHeader:  "",
			xAuthorizationHeader: "",
			expectedCredential:   "",
		},
		{
			authorizationHeader:  "Foo",
			xAuthorizationHeader: "",
			expectedCredential:   "Foo",
		},
		{
			authorizationHeader:  "",
			xAuthorizationHeader: "Foo",
			expectedCredential:   "Foo",
		},
		{
			authorizationHeader:  "Foo",
			xAuthorizationHeader: "Bar",
			expectedCredential:   "Foo",
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s%s", tt.authorizationHeader, tt.xAuthorizationHeader), func(t *testing.T) {
			var innerContext context.Context
			request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
			request.Header.Set(auth.AuthorizationHeader, tt.authorizationHeader)
			request.Header.Set(auth.XAuthorizationHeader, tt.xAuthorizationHeader)

			handlerCalled := false
			handler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
				handlerCalled = true
				innerContext = r.Context()
			})

			middleware.Authentication()(handler).ServeHTTP(nil, request)

			assert.True(t, handlerCalled)
			assert.Equal(t, auth.Credential(tt.expectedCredential), auth.CredentialFromContext(innerContext))
		})
	}
}

func TestAuthentication_Unauthorize(t *testing.T) {
	var innerContext context.Context
	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	recorder := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		innerContext = r.Context()
	})

	middleware.Authentication(middleware.WithAuthFunc(func(req *http.Request) (context.Context, error) {
		return req.Context(), errors.New("my_authenticate_error")
	}))(handler).ServeHTTP(recorder, request)

	assert.False(t, handlerCalled)
	assert.Equal(t, auth.Credential(""), auth.CredentialFromContext(innerContext))
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
}

func TestAuthentication_Custom_Error_Handler(t *testing.T) {
	var innerContext context.Context
	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	recorder := httptest.NewRecorder()

	handlerCalled := false
	handler := http.HandlerFunc(func(_ http.ResponseWriter, r *http.Request) {
		handlerCalled = true
		innerContext = r.Context()
	})

	middleware.Authentication(
		middleware.WithAuthFunc(func(req *http.Request) (context.Context, error) {
			return req.Context(), errors.New("my_authenticate_error")
		}),
		middleware.WithErrorHandler(func(err error, writer http.ResponseWriter, req *http.Request) bool {
			_, _ = writer.Write([]byte("my_fake_response"))
			return true
		}),
	)(handler).ServeHTTP(recorder, request)

	assert.False(t, handlerCalled)
	assert.Equal(t, auth.Credential(""), auth.CredentialFromContext(innerContext))
	assert.Equal(t, "my_fake_response", recorder.Body.String())
}
