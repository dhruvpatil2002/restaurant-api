package middleware

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"

	"restaurant-backend/internal/models"
	"restaurant-backend/internal/service"
)

// =====================================================
// CONTEXT KEYS (exported)
// =====================================================

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	RoleKey   contextKey = "role"
)

// =====================================================
// CONTEXT HELPERS (recommended for handlers)
// =====================================================

// GetUserIDFromContext extracts the user UUID from the request context.
// Returns the UUID and a boolean indicating success.
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return uuid.Nil, false
	}
	id, ok := val.(uuid.UUID)
	return id, ok
}

// GetRoleFromContext extracts the user role from the request context.
// Returns the role and a boolean indicating success.
func GetRoleFromContext(ctx context.Context) (models.UserRole, bool) {
	val := ctx.Value(RoleKey)
	if val == nil {
		return "", false
	}
	role, ok := val.(models.UserRole)
	return role, ok
}

// =====================================================
// AUTH MIDDLEWARE
// =====================================================

func AuthMiddleware(
	authSvc *service.AuthService,
) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			authHeader := strings.TrimSpace(
				r.Header.Get("Authorization"),
			)

			if authHeader == "" {
				writeJSONUnauth(
					w,
					"authorization header required",
				)
				return
			}

			parts := strings.Fields(authHeader)
			if len(parts) != 2 ||
				!strings.EqualFold(parts[0], "Bearer") {
				writeJSONUnauth(
					w,
					"invalid authorization header",
				)
				return
			}

			tokenString := strings.TrimSpace(parts[1])
			if tokenString == "" {
				writeJSONUnauth(
					w,
					"access token required",
				)
				return
			}

			claims, err :=
				authSvc.ParseAndValidateAccessToken(
					tokenString,
				)
			if err != nil {
				writeJSONUnauth(
					w,
					"invalid or expired access token",
				)
				return
			}

			// Inject user ID and role into context
			ctx := r.Context()
			ctx = context.WithValue(ctx, UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, RoleKey, claims.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// =====================================================
// REQUIRE ROLE
// =====================================================

func RequireRole(
	allowedRoles ...models.UserRole,
) func(http.Handler) http.Handler {

	// Pre‑normalise allowed roles for faster comparison
	allowedSet := make(map[models.UserRole]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		// Use as‑is – roles are already typed constants
		allowedSet[r] = struct{}{}
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(
			w http.ResponseWriter,
			r *http.Request,
		) {
			role, ok := GetRoleFromContext(r.Context())
			if !ok {
				writeJSONForbidden(
					w,
					"role not found in context",
				)
				return
			}

			if _, allowed := allowedSet[role]; allowed {
				next.ServeHTTP(w, r)
				return
			}

			writeJSONForbidden(
				w,
				"insufficient permissions",
			)
		})
	}
}

// =====================================================
// LOGGING
// =====================================================

func Logging(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		start := time.Now()

		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(
			wrapped,
			r,
		)

		log.Printf(
			"[%s] %s %s %d %v",
			r.Method,
			r.URL.Path,
			getClientIP(r),
			wrapped.statusCode,
			time.Since(start),
		)
	})
}

// =====================================================
// RESPONSE WRITER (for logging status)
// =====================================================

type responseWriter struct {
	http.ResponseWriter
	statusCode  int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(
	code int,
) {
	if rw.wroteHeader {
		return
	}
	rw.statusCode = code
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(
	data []byte,
) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(data)
}

// =====================================================
// RECOVERY
// =====================================================

func Recovery(
	next http.Handler,
) http.Handler {

	return http.HandlerFunc(func(
		w http.ResponseWriter,
		r *http.Request,
	) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf(
					"panic recovered: %v",
					err,
				)
				writeJSONError(
					w,
					"internal server error",
					http.StatusInternalServerError,
				)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// =====================================================
// JSON HELPERS
// =====================================================

func writeJSONUnauth(
	w http.ResponseWriter,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnauthorized)
	_ = json.NewEncoder(w).Encode(
		map[string]string{"error": message},
	)
}

func writeJSONForbidden(
	w http.ResponseWriter,
	message string,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	_ = json.NewEncoder(w).Encode(
		map[string]string{"error": message},
	)
}

func writeJSONError(
	w http.ResponseWriter,
	message string,
	statusCode int,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(
		map[string]string{"error": message},
	)
}

// =====================================================
// CLIENT IP
// =====================================================

func getClientIP(
	r *http.Request,
) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[0])
		}
	}
	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil {
		return host
	}
	return r.RemoteAddr
}