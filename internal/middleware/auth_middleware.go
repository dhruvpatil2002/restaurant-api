package middleware

import (
    "context"
    "net/http"
   "encoding/json"

    "restaurant-backend/internal/service"
)

type contextKey string

const (
    UserIDKey contextKey = "user_id"
    RoleKey   contextKey = "role"
)

func AuthMiddleware(authSvc *service.AuthService) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                writeJSONUnauth(w)
                return
            }

            tokenStr := authHeader
            if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
                tokenStr = authHeader[7:]
            } else {
                writeJSONUnauth(w)
                return
            }

            claims, err := authSvc.ParseAndValidateAccessToken(tokenStr)
            if err != nil {
                writeJSONUnauth(w)
                return
            }

            ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
            ctx = context.WithValue(ctx, RoleKey, claims.Role)

            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

func writeJSONUnauth(w http.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusUnauthorized)
    _ = json.NewEncoder(w).Encode(map[string]string{"error": "unauthorized"})
}