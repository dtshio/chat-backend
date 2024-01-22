package server

import (
	"net/http"
	"strings"

	"github.com/datsfilipe/pkg/application/auth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}

		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		authParts := strings.Split(authHeader, "Bearer ")
		if len(authParts) != 2 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		token := authParts[1]

		isAuthorized := auth.VerifyToken(token)

		if !isAuthorized {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
