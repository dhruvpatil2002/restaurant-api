package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
    "restaurant-backend/internal/models"
   
	"restaurant-backend/internal/middleware"
	"restaurant-backend/internal/service"
)





// ---------- handler struct ----------

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

// ---------- request types ----------

// handler/auth_handler.go

type RegisterRequest struct {
    Name     string          `json:"name"`
    Email    string          `json:"email"`
    Password string          `json:"password"`
    Role     *models.UserRole `json:"role"`  // ← Use pointer to distinguish "not provided" vs "empty"
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    var req RegisterRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
        return
    }
    
    // Log the raw role value for debugging
    var rawRole string
    if req.Role != nil {
        rawRole = string(*req.Role)
    } else {
        rawRole = "(nil)"
    }
    log.Printf("📥 Register payload: name=%s, email=%s, role=%q", req.Name, req.Email, rawRole)

    // Validate required fields
    if req.Name == "" || req.Email == "" || req.Password == "" {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name, email, password required"})
        return
    }

    // Validate role if provided
    if req.Role != nil {
        switch *req.Role {
        case models.RoleCustomer, models.RoleStaff, models.RoleOwner, models.RoleAdmin:
            // valid
        default:
            writeJSON(w, http.StatusBadRequest, map[string]string{
                "error": fmt.Sprintf("invalid role: %q", *req.Role),
            })
            return
        }
    }

    // Pass role to service (nil means use default)
    var role models.UserRole
    if req.Role != nil {
        role = *req.Role
    }

    user, err := h.authSvc.Register(r.Context(), service.RegisterInput{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
        Role:     role,
    })
    if err != nil {
        switch {
        case errors.Is(err, service.ErrEmailAlreadyTaken):
            writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
        case errors.Is(err, service.ErrEmailInvalid):
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        case errors.Is(err, service.ErrInvalidRole):
            writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
        default:
            writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
        }
        return
    }

    log.Printf("✅ User created: id=%s, role=%q", user.ID, user.Role)

    writeJSON(w, http.StatusCreated, map[string]any{
        "message": "user created",
        "user_id": user.ID,
        "role":    user.Role,
    })
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid json"})
		return
	}

	if req.Email == "" || req.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password required"})
		return
	}

	out, err := h.authSvc.Login(r.Context(), service.LoginInput{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "login failed"})
		return
	}

	// Set refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    out.RefreshTokenRaw,
		Path:     "/auth/refresh",
		MaxAge:   int(out.RefreshToken.ExpiresAt.Sub(time.Now()).Seconds()),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	writeJSON(w, http.StatusOK, map[string]any{
		"access_token": out.AccessToken,
		"user": map[string]any{
			"id":    out.User.ID,
			"name":  out.User.Name,
			"email": out.User.Email,
			"role":  out.User.Role,
		},
	})
}

func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c, err := r.Cookie("refresh_token")
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "missing refresh token"})
		return
	}

	out, err := h.authSvc.Refresh(r.Context(), c.Value)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh token"})
		return
	}

	maxAge := int(time.Until(out.RefreshToken.ExpiresAt).Seconds())
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    out.RefreshTokenRaw,
		Path:     "/auth/refresh",
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	writeJSON(w, http.StatusOK, map[string]string{
		"access_token": out.AccessToken,
	})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	c, err := r.Cookie("refresh_token")
	if err == nil {
		_ = h.authSvc.Logout(r.Context(), c.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/auth/refresh",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})

	writeJSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ✅ Fixed syntax error: type assertion in one line
	userID, ok := r.Context().Value(middleware.UserIDKey).(uuid.UUID)
	if !ok {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "user_id missing in context"})
		return
	}

	user, err := h.authSvc.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]any{
		"id":            user.ID,
		"name":          user.Name,
		"email":         user.Email,
		"role":          user.Role,
		"restaurant_id": user.RestaurantID,
	})
}