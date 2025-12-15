package middlewares

import (
	"backend/generate"
	"context"
	"net/http"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		token, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "No token in cookies", http.StatusUnauthorized)
			return
		}

		if token.Value == "" {
			http.Error(w, "empty token value", http.StatusUnauthorized)
			return
		}

		claims, err := generate.ValidateToken(token.Value)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), generate.UserIdKey, claims.UserId)
		ctx = context.WithValue(ctx, generate.RoleKey, claims.Role)

		next.ServeHTTP(w, r)
	})
}
