package middleware

import (
	"fmt"
	"net/http"
)

// RequireRole is a middleware that checks if the user has the required role
func RequireRole(requiredRole string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authenticated, role := CheckAuth(r)
			if !authenticated {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if role != requiredRole {
				http.Error(w, fmt.Sprintf("Access denied for role: %s", role), http.StatusForbidden)
				return
			}

			// Continue to next handler if role matches
			next.ServeHTTP(w, r)
		})
	}
}
