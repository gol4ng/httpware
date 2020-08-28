package auth_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gol4ng/httpware/v3/auth"
	"github.com/stretchr/testify/assert"
)

func TestFromHeader(t *testing.T) {
	tests := []struct {
		request            *http.Request
		expectedCredential string
	}{
		{
			request:            nil,
			expectedCredential: "",
		},
		{
			request: &http.Request{Header: http.Header{
				"Authorization": []string{"foo"},
			},},
			expectedCredential: "foo",
		},
		{
			request: &http.Request{Header: http.Header{
				"X-Authorization": []string{"foo"},
			},},
			expectedCredential: "foo",
		},
		{
			request: &http.Request{Header: http.Header{
				"Authorization": []string{"foo"},
				"X-Authorization": []string{"bar"},
			},},
			expectedCredential: "foo",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Equal(t, auth.Credential(tt.expectedCredential), auth.FromHeader(tt.request)())
		})
	}
}

func TestAddHeader(t *testing.T) {
	req := &http.Request{
		Header: make(http.Header),
	}

	credSetter := auth.AddHeader(req)
	credSetter("foo")
	assert.Equal(t, "foo", req.Header.Get("Authorization"))
}
