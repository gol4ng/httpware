package auth_test

import (
	"context"
	"testing"

	"github.com/gol4ng/httpware/v2/auth"
	"github.com/stretchr/testify/assert"
)

func TestCredentialFromContext(t *testing.T) {
	ctx := context.Background()
	newCtx := auth.CredentialToContext(ctx, "foo")

	cred := auth.CredentialFromContext(newCtx)
	assert.Equal(t, "foo", string(cred))
}
