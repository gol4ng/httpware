package auth

import (
	"context"
)

type AuthFunc func(ctx context.Context) (context.Context, error)
