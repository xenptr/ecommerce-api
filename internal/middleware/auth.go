package middleware

import (
	"net/http"
	"strings"

	"github.com/xenptr/ecommerce-api/internal/auth"
	"github.com/xenptr/ecommerce-api/internal/token"
)

func Auth(parser token.Parser) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userToken := r.Header.Get("Authorization")

			if !strings.HasPrefix(userToken, "Bearer ") {
				http.Error(w, "missing or invalid authorization header", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(userToken, "Bearer ")

			userID, err := parser.Parse(tokenString)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}

			r = auth.SetUserID(r, userID)

			next.ServeHTTP(w, r)
		})
	}
}
