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
		expectedCredential string
	}{
		{
			context:            nil,
			expectedCredential: "",
		},
		{
			context:            context.Background(),
			expectedCredential: "",
		},
		{
			context:            context.WithValue(context.Background(), credentialContextKey, "not a credential"),
			expectedCredential: "",
		},
		{
			context:            CredentialToContext(context.Background(), Credential("my_value")),
			expectedCredential: "my_value",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprint(i), func(t *testing.T) {
			assert.Equal(t, Credential(tt.expectedCredential), CredentialFromContext(tt.context))
		})
	}
}
