package tripperware_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gol4ng/httpware/v2/auth"
	"github.com/gol4ng/httpware/v2/mocks"
	"github.com/gol4ng/httpware/v2/tripperware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthenticationForwarder(t *testing.T) {
	tests := []struct {
		context               context.Context
		expectedAuthorization string
	}{
		{
			context:               context.TODO(),
			expectedAuthorization: "",
		},
		{
			context:               auth.CredentialToContext(context.TODO(), "my_credential"),
			expectedAuthorization: "my_credential",
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			roundTripperMock := &mocks.RoundTripper{}
			request := httptest.NewRequest(http.MethodGet, "http://fake-addr", nil)
			request = request.WithContext(tt.context)

			roundTripperMock.On("RoundTrip", mock.AnythingOfType("*http.Request")).Return(nil, nil).Run(func(args mock.Arguments) {
				innerReq := args.Get(0).(*http.Request)
				assert.Equal(t, tt.expectedAuthorization, innerReq.Header.Get(auth.AuthorizationHeader))
				assert.Equal(t, tt.expectedAuthorization, innerReq.Header.Get(auth.XAuthorizationHeader))
			})

			_, _ = tripperware.AuthenticationForwarder()(roundTripperMock).RoundTrip(request)
			roundTripperMock.AssertExpectations(t)
		})
	}
}
