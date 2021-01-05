package middleware_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gol4ng/httpware/v4/auth"
	"github.com/gol4ng/httpware/v4/middleware"
	"github.com/gol4ng/httpware/v4/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDefaultCredentialFinder(t *testing.T) {
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
			request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
			request.Header.Set(auth.AuthorizationHeader, tt.authorizationHeader)
			request.Header.Set(auth.XAuthorizationHeader, tt.xAuthorizationHeader)

			assert.Equal(t, auth.Credential(tt.expectedCredential), middleware.DefaultCredentialFinder(request))
		})
	}
}

func TestDefaultErrorHandler(t *testing.T) {
	request, _ := http.NewRequest("", "", nil)
	response := httptest.NewRecorder()

	middleware.DefaultErrorHandler(errors.New("my_fake_error"), response, request)

	assert.Equal(t, http.StatusUnauthorized, response.Code)
	assert.Equal(t, "my_fake_error\n", response.Body.String())
}

func TestAuthentication(t *testing.T) {
	var innerContext context.Context
	request, _ := http.NewRequest(http.MethodGet, "http://fake-addr", nil)

	handlerCalled := false
	handler := http.HandlerFunc(func(_ http.ResponseWriter, innerRequest *http.Request) {
		handlerCalled = true
		innerContext = innerRequest.Context()
	})

	authMiddleware := middleware.Authentication(func(request *http.Request) (*http.Request, error) {
		newCtx := auth.CredentialToContext(request.Context(), "my_allowed_credential")
		return request.WithContext(newCtx), nil
	})

	authMiddleware(handler).ServeHTTP(nil, request)
	assert.True(t, handlerCalled)
	assert.NotEqual(t, request.Context(), innerContext)
	assert.Equal(t, "my_allowed_credential", auth.CredentialFromContext(innerContext))
}

func TestAuthentication_WithSuccessMiddleware(t *testing.T) {
	var innerContext context.Context
	request, _ := http.NewRequest(http.MethodGet, "http://fake-addr", nil)

	handlerCalled := false
	handler := http.HandlerFunc(func(_ http.ResponseWriter, innerRequest *http.Request) {
		handlerCalled = true
		innerContext = innerRequest.Context()
	})

	authMiddleware := middleware.Authentication(
		func(req *http.Request) (*http.Request, error) {
			return req, nil
		},
		middleware.WithSuccessMiddleware(func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
				assert.Nil(t, writer)
				assert.Equal(t, request, req)
				// we not call next handler for example
			})
		}),
	)

	authMiddleware(handler).ServeHTTP(nil, request)
	assert.False(t, handlerCalled)
	assert.NotEqual(t, request.Context(), innerContext)
	assert.Equal(t, nil, auth.CredentialFromContext(innerContext))
}

func TestAuthentication_WithErrorHandler(t *testing.T) {
	var innerErr error
	request, _ := http.NewRequest(http.MethodGet, "http://fake-addr", nil)

	handlerCalled := false
	handler := http.HandlerFunc(func(_ http.ResponseWriter, innerRequest *http.Request) {
		handlerCalled = true
	})

	authMiddleware := middleware.Authentication(
		func(req *http.Request) (*http.Request, error) {
			return req, errors.New("my_authenticate_error")
		},
		middleware.WithErrorHandler(func(err error, _ http.ResponseWriter, _ *http.Request) bool {
			innerErr = err
			return true
		}),
	)

	authMiddleware(handler).ServeHTTP(nil, request)
	assert.False(t, handlerCalled)
	assert.EqualError(t, innerErr, "my_authenticate_error")
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

	authMiddleware := middleware.Authentication(func(request *http.Request) (*http.Request, error) {
		return request, errors.New("my_authenticated_error")
	})

	authMiddleware(handler).ServeHTTP(recorder, request)

	assert.False(t, handlerCalled)
	assert.Equal(t, nil, auth.CredentialFromContext(innerContext))
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "my_authenticated_error\n", recorder.Body.String())
}

func TestNewAuthenticateFunc(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	request.Header.Set("Authorization", "my_credential")

	authenticator := &mocks.Authenticator{}
	authenticator.On("Authenticate", "my_credential").Return("my_authenticate_credential", nil)

	authenticateFunc := middleware.NewAuthenticateFunc(authenticator)

	newRequest, err := authenticateFunc(request)
	assert.NoError(t, err)
	assert.Equal(t, "my_authenticate_credential", auth.CredentialFromContext(newRequest.Context()))

	authenticator.AssertExpectations(t)
}

func TestNewAuthenticateFunc_WithCredentialFinder(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)

	authenticator := &mocks.Authenticator{}
	authenticator.On("Authenticate", "my_credential_finder_value").Return("my_authenticate_credential", nil)

	authenticateFunc := middleware.NewAuthenticateFunc(
		authenticator,
		middleware.WithCredentialFinder(func(r *http.Request) auth.Credential {
			return "my_credential_finder_value"
		}),
	)

	newRequest, err := authenticateFunc(request)
	assert.NoError(t, err)
	assert.Equal(t, "my_authenticate_credential", auth.CredentialFromContext(newRequest.Context()))

	authenticator.AssertExpectations(t)
}

func TestNewAuthenticateFunc_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
	request.Header.Set("Authorization", "my_credential")

	err := errors.New("my_authenticate_error")
	authenticator := &mocks.Authenticator{}
	authenticator.On("Authenticate", "my_credential").Return("my_authenticate_credential", err)

	authenticateFunc := middleware.NewAuthenticateFunc(authenticator)

	newRequest, err := authenticateFunc(request)
	assert.EqualError(t, err, "my_authenticate_error")
	assert.Nil(t, auth.CredentialFromContext(newRequest.Context()))

	authenticator.AssertExpectations(t)
}
