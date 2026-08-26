package middleware

import (
	"net/http"
)

func RequireRole(allowedRoles ...string) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			roleValue := r.Context().Value(RoleKey)

			role, ok := roleValue.(string)

			if !ok || role == "" {
				writeJSONForbidden(w)
				return
			}

			// Check whether user's role is allowed
			for _, allowedRole := range allowedRoles {

				if role == allowedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			writeJSONForbidden(w)
		})
	}
}

func writeJSONForbidden(w http.ResponseWriter) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)

	w.Write([]byte(`{"error":"forbidden"}`))
}