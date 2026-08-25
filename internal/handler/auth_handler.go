package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"restaurant-backend/internal/service"

	"restaurant-backend/internal/models"

	"gorm.io/gorm"
)

type AuthHandler struct {
    authSvc *service.AuthService
}

type RegisterRequest struct {
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
    Role     string `json:"role"`
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
    return &AuthHandler{authSvc: authSvc}
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

    if req.Name == "" || req.Email == "" || req.Password == "" {
        writeJSON(w, http.StatusBadRequest, map[string]string{"error": "name, email, password required"})
        return
    }

    role := "customer"
    if req.Role != "" {
        role = req.Role
    }

    user, err := h.authSvc.Register(service.RegisterInput{
        Name:     req.Name,
        Email:    req.Email,
        Password: req.Password,
        Role:     models.UserRole(role),
    })
    if err != nil {
        if err.Error() == "email already registered" {
            writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
            return
        }
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create user"})
        return
    }

    writeJSON(w, http.StatusCreated, map[string]any{
        "message": "user created",
        "user_id": user.ID,
    })
}

type LoginRequest struct {
    Email    string `json:"email"`
    Password string `json:"password"`
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

    out, err := h.authSvc.Login(service.LoginInput{
        Email:    req.Email,
        Password: req.Password,
    })
    if err != nil {
        if err == service.ErrInvalidCredentials {
            writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
            return
        }
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "login failed"})
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    out.RefreshTokenRaw,
        Path:     "/auth/refresh",
        MaxAge:   int(out.RefreshToken.ExpiresAt.Sub(time.Now()).Seconds()),
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    })

    writeJSON(w, http.StatusOK, map[string]string{"access_token": out.AccessToken})
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

    out, err := h.authSvc.Refresh(c.Value)
    if err != nil {
        writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid or expired refresh token"})
        return
    }

    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    out.RefreshTokenRaw,
        Path:     "/auth/refresh",
        MaxAge:   int(out.RefreshToken.ExpiresAt.Sub(time.Now()).Seconds()),
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    })

    writeJSON(w, http.StatusOK, map[string]string{"access_token": out.AccessToken})
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
        return
    }

    c, err := r.Cookie("refresh_token")
    if err == nil {
        _ = h.authSvc.Logout(c.Value)
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

    userIDRaw := r.Context().Value("user_id")
    userID, ok := userIDRaw.(int64)
    if !ok {
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "user_id missing"})
        return
    }

    user, err := h.authSvc.GetUserByID(userID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            writeJSON(w, http.StatusNotFound, map[string]string{"error": "user not found"})
            return
        }
        writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get user"})
        return
    }

    writeJSON(w, http.StatusOK, map[string]any{
        "id":    user.ID,
        "name":  user.Name,
        "email": user.Email,
        "role":  user.Role,
    })
}
