// middleware/jwt.go
package middleware

import (
	"context"
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/badtripdude/hackathon-sonnik-backend/services/user/pkg/auth"
	errs "github.com/badtripdude/hackathon-sonnik-backend/services/user/pkg/errors"
)

func JWTMiddleware(pubKey *rsa.PublicKey) func(next http.HandlerFunc) http.Handler {
	return func(next http.HandlerFunc) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "authorization header required", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				http.Error(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenStr := parts[1]

			claims, err := auth.ParseAndValidateAccessToken(tokenStr, pubKey)
			if err != nil {
				if err == errs.ErrAccessTokenExpired {
					http.Error(w, "access token expired", http.StatusUnauthorized)
					return
				}
				http.Error(w, "invalid access token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", claims.Subject)
			ctx = context.WithValue(ctx, "userEmail", claims.Email)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
