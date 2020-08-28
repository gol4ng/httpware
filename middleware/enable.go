package middleware

import (
	"github.com/gol4ng/httpware/v3"
)

// Enable middleware is used to conditionnaly add a middleware to a MiddlewareStack
// See Skip middleware to active a middleware in function of request
func Enable(enable bool, middleware httpware.Middleware) httpware.Middleware {
	if enable {
		return middleware
	}
	return httpware.NopMiddleware
}
