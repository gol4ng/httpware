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
	"github.com/gol4ng/httpware/v2/mocks"
	"github.com/stretchr/testify/assert"
)

func credentialFinderMock(_ *http.Request) auth.Credential {
	return "my_credential"
}

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

func TestDefaultAuthFunc(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)

	nexCtx, err := middleware.DefaultAuthFunc(credentialFinderMock, nil, request)
	assert.NoError(t, err)
	assert.Equal(t, auth.Credential("my_credential"), auth.CredentialFromContext(nexCtx))
}

func TestDefaultAuthFunc_WithAuthenticator(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)

	authenticator := &mocks.Authenticator{}
	authenticator.On("Authenticate", "my_credential").Return("my_authenticate_credential", nil)

	nexCtx, err := middleware.DefaultAuthFunc(credentialFinderMock, authenticator, request)
	assert.NoError(t, err)
	assert.Equal(t, auth.Credential("my_authenticate_credential"), auth.CredentialFromContext(nexCtx))
	authenticator.AssertExpectations(t)
}

func TestDefaultAuthFunc_WithAuthenticator_Error(t *testing.T) {
	request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)

	err := errors.New("my_authenticate_error")
	authenticator := &mocks.Authenticator{}
	authenticator.On("Authenticate", "my_credential").Return("my_authenticate_credential", err)

	nexCtx, err := middleware.DefaultAuthFunc(credentialFinderMock, authenticator, request)
	assert.EqualError(t, err, "my_authenticate_error")
	assert.Equal(t, nil, auth.CredentialFromContext(nexCtx))
	authenticator.AssertExpectations(t)
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

	authMiddleware := middleware.Authentication(nil, middleware.WithAuthenticateFunc(func(_ middleware.CredentialFinder, _ auth.Authenticator, req *http.Request) (context.Context, error) {
		return req.Context(), nil
	}))

	authMiddleware(handler).ServeHTTP(nil, request)
	assert.True(t, handlerCalled)
	assert.Equal(t, request.Context(), innerContext)
	assert.Equal(t, nil, auth.CredentialFromContext(innerContext))
}

func TestAuthentication_Error(t *testing.T) {
	var innerErr error
	request, _ := http.NewRequest(http.MethodGet, "http://fake-addr", nil)

	handlerCalled := false
	handler := http.HandlerFunc(func(_ http.ResponseWriter, innerRequest *http.Request) {
		handlerCalled = true
	})

	authMiddleware := middleware.Authentication(
		nil,
		middleware.WithAuthenticateFunc(func(_ middleware.CredentialFinder, _ auth.Authenticator, req *http.Request) (context.Context, error) {
			return req.Context(), errors.New("my_authenticate_error")
		}),
		middleware.WithErrorHandler(func(err error, _ http.ResponseWriter, _ *http.Request) bool {
			innerErr = err
			return true
		}),
	)

	authMiddleware(handler).ServeHTTP(nil, request)
	assert.False(t, handlerCalled)
	assert.EqualError(t, innerErr, "my_authenticate_error")
}

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

			authMiddleware := middleware.Authentication(nil)

			authMiddleware(handler).ServeHTTP(nil, request)

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

	authenticator := &mocks.Authenticator{}
	authenticator.On("Authenticate", "").Return("my_authenticated_credential", errors.New("my_authenticated_error"))
	authMiddleware := middleware.Authentication(authenticator)

	authMiddleware(handler).ServeHTTP(recorder, request)

	assert.False(t, handlerCalled)
	assert.Equal(t, nil, auth.CredentialFromContext(innerContext))
	assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	assert.Equal(t, "my_authenticated_error\n", recorder.Body.String())
}
