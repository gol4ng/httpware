package auth

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Credential_Context(t *testing.T) {
	tests := []struct {
		context            context.Context
		expectedCredential Credential
	}{
		{
			context:            nil,
			expectedCredential: nil,
		},
		{
			context:            context.Background(),
			expectedCredential: nil,
		},
		{
			context:            CredentialToContext(context.Background(), Credential("my_value")),
			expectedCredential: "my_value",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Equal(t, tt.expectedCredential, CredentialFromContext(tt.context))
		})
	}
}
