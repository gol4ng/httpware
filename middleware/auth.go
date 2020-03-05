package middleware

import (
	"net/http"

	"github.com/gol4ng/httpware/v2"
	"github.com/gol4ng/httpware/v2/auth"
)

func AuthenticationMiddleware(authFunc auth.AuthFunc) httpware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
			newCtx, err := authFunc(auth.CredentialToContext(req.Context(), auth.FromHeader(req)()))
			if err != nil {
				http.Error(writer, err.Error(), http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(writer, req.WithContext(newCtx))
		})
	}
}
