package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"restaurant-backend/internal/config"
	"restaurant-backend/internal/models"
	"restaurant-backend/internal/repository"
)

// Sentinel errors
var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrInvalidToken        = errors.New("invalid token")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrEmailAlreadyTaken   = errors.New("email already registered")
	ErrEmailInvalid        = errors.New("invalid email format")
	ErrInvalidRole         = errors.New("invalid role")
)

const (
	accessTokenIssuer   = "restaurant-backend"
	refreshTokenIssuer  = "restaurant-backend-refresh"
	tokenAudience       = "restaurant-api"
	tokenTypeAccess     = "access"
	tokenTypeRefresh    = "refresh"
)

// =====================================================
// AUTH SERVICE
// =====================================================

type AuthService struct {
	userRepo *repository.UserRepository
	rtRepo   *repository.RefreshTokenRepository
	cfg      *config.Config
}

func NewAuthService(
	userRepo *repository.UserRepository,
	rtRepo *repository.RefreshTokenRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{
		userRepo: userRepo,
		rtRepo:   rtRepo,
		cfg:      cfg,
	}
}

// =====================================================
// PASSWORD HELPERS
// =====================================================

func HashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password must be at least 6 characters")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CheckPassword(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// =====================================================
// EMAIL VALIDATION
// =====================================================

func validateEmail(email string) error {
	// Simple email format check
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	if matched, _ := regexp.MatchString(emailRegex, email); !matched {
		return ErrEmailInvalid
	}
	return nil
}

// =====================================================
// REGISTER
// =====================================================

type RegisterInput struct {
	Name     string
	Email    string
	Password string
	Role     models.UserRole
}

// Add this sentinel error

// service/auth_service.go

func (s *AuthService) Register(ctx context.Context, in RegisterInput) (*models.User, error) {
    in.Email = strings.ToLower(strings.TrimSpace(in.Email))
    in.Name = strings.TrimSpace(in.Name)
    
    // Normalize role
    roleStr := strings.ToLower(strings.TrimSpace(string(in.Role)))
    log.Printf("🔄 Service: received role=%q, normalized to %q", string(in.Role), roleStr)
    
    // Default to customer if empty
    if roleStr == "" {
        roleStr = string(models.RoleCustomer)
    }
    
    // Validate against allowed roles
    switch roleStr {
    case string(models.RoleCustomer):
        in.Role = models.RoleCustomer
    case string(models.RoleStaff):
        in.Role = models.RoleStaff
    case string(models.RoleOwner):
        in.Role = models.RoleOwner
    case string(models.RoleAdmin):
        in.Role = models.RoleAdmin
    default:
        log.Printf("❌ Invalid role provided: %q", roleStr)
        return nil, ErrInvalidRole
    }

    passwordHash, err := HashPassword(in.Password)
    if err != nil {
        return nil, err
    }

    user := &models.User{
        ID:           uuid.New(),
        Name:         in.Name,
        Email:        in.Email,
        PasswordHash: passwordHash,
        Role:         in.Role,
    }
    
    log.Printf("💾 Saving user with role=%q (type: %T)", user.Role, user.Role)
    
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, fmt.Errorf("failed to create user: %w", err)
    }
    
    log.Printf("✅ User saved: id=%s, role=%q", user.ID, user.Role)
    return user, nil
}


// =====================================================
// LOGIN
// =====================================================

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	User            *models.User
	AccessToken     string
	RefreshTokenRaw string
	RefreshToken    *models.RefreshToken
}

func (s *AuthService) Login(ctx context.Context, in LoginInput) (*LoginOutput, error) {
	in.Email = strings.ToLower(strings.TrimSpace(in.Email))

	user, err := s.userRepo.GetByEmail(ctx, in.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if !CheckPassword(in.Password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	// Generate access token
	accessToken, err := s.newAccessToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create access token: %w", err)
	}

	// Generate refresh token
	refreshRaw, refreshToken, err := s.newRefreshTokenRaw(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	// Store hashed refresh token
	if err := s.rtRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to store refresh token: %w", err)
	}

	return &LoginOutput{
		User:            user,
		AccessToken:     accessToken,
		RefreshTokenRaw: refreshRaw,
		RefreshToken:    refreshToken,
	}, nil
}

// =====================================================
// REFRESH
// =====================================================

type RefreshOutput struct {
	AccessToken     string
	RefreshTokenRaw string
	RefreshToken    *models.RefreshToken
}

func (s *AuthService) Refresh(ctx context.Context, rawToken string) (*RefreshOutput, error) {
    rawToken = strings.TrimSpace(rawToken)
    if rawToken == "" {
        return nil, ErrInvalidRefreshToken
    }

    tokenHash := hashToken(rawToken)

    // Find token in database (not revoked)
    refreshToken, err := s.rtRepo.GetByHashAndNotRevoked(ctx, tokenHash)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrInvalidRefreshToken
        }
        return nil, fmt.Errorf("failed to get refresh token: %w", err)
    }

    // Check expiration (redundant with JWT but safe)
    if time.Now().After(refreshToken.ExpiresAt) {
        _ = s.rtRepo.Revoke(ctx, refreshToken) // best effort
        return nil, ErrInvalidRefreshToken
    }

    // Validate JWT signature and claims
    claims, err := s.parseRefreshToken(rawToken)
    if err != nil {
        _ = s.rtRepo.Revoke(ctx, refreshToken) // revoke on invalid token
        return nil, ErrInvalidRefreshToken
    }

    // Ensure token belongs to the user
    if claims.UserID != refreshToken.UserID {
        _ = s.rtRepo.Revoke(ctx, refreshToken)
        return nil, ErrInvalidRefreshToken
    }

    // Fetch the user
    user, err := s.userRepo.GetByID(ctx, refreshToken.UserID)
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // user deleted – revoke token
            _ = s.rtRepo.Revoke(ctx, refreshToken)
            return nil, ErrInvalidRefreshToken
        }
        return nil, fmt.Errorf("failed to get user: %w", err)
    }

    // =================================================
    // ROTATE REFRESH TOKEN (safer ordering - create before revoke)
    // =================================================
    
    // 1. Generate new tokens FIRST (before revoking old one)
    newAccessToken, err := s.newAccessToken(user)
    if err != nil {
        return nil, fmt.Errorf("failed to create new access token: %w", err)
    }

    newRefreshRaw, newRefreshToken, err := s.newRefreshTokenRaw(user)
    if err != nil {
        return nil, fmt.Errorf("failed to create new refresh token: %w", err)
    }

    // 2. Store new refresh token
    if err := s.rtRepo.Create(ctx, newRefreshToken); err != nil {
        // New token wasn't stored, but old token is still valid
        // Client can retry refresh with the same token
        return nil, fmt.Errorf("failed to store new refresh token: %w", err)
    }

    // 3. Revoke old token (only after new one is safely stored)
    if err := s.rtRepo.Revoke(ctx, refreshToken); err != nil {
        // This is less critical - old token will expire naturally
        // or be revoked on next refresh. Log it for monitoring.
        log.Printf("⚠️ Failed to revoke old refresh token: %v", err)
    }

    return &RefreshOutput{
        AccessToken:     newAccessToken,
        RefreshTokenRaw: newRefreshRaw,
        RefreshToken:    newRefreshToken,
    }, nil
}

// =====================================================
// LOGOUT
// =====================================================

func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	rawToken = strings.TrimSpace(rawToken)
	if rawToken == "" {
		return nil
	}

	tokenHash := hashToken(rawToken)
	refreshToken, err := s.rtRepo.GetByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil // already logged out or not found
		}
		return fmt.Errorf("failed to get refresh token: %w", err)
	}

	if refreshToken.Revoked {
		return nil
	}

	return s.rtRepo.Revoke(ctx, refreshToken)
}

// =====================================================
// GET USER
// =====================================================

func (s *AuthService) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return user, nil
}

// =====================================================
// JWT CLAIMS
// =====================================================

type AuthClaims struct {
	UserID    uuid.UUID       `json:"user_id"`
	Role      models.UserRole `json:"role"`
	TokenType string          `json:"token_type"`
	jwt.RegisteredClaims
}

// =====================================================
// ACCESS TOKEN CREATION
// =====================================================

func (s *AuthService) newAccessToken(user *models.User) (string, error) {
	now := time.Now()
	claims := AuthClaims{
		UserID:    user.ID,
		Role:      user.Role,
		TokenType: tokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    accessTokenIssuer,
			Audience:  jwt.ClaimStrings{tokenAudience},
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.JWTAccessExpiry)),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.cfg.JWTSigningKey)
}

// =====================================================
// REFRESH TOKEN CREATION
// =====================================================

func (s *AuthService) newRefreshTokenRaw(user *models.User) (string, *models.RefreshToken, error) {
	now := time.Now()
	expiry := now.Add(s.cfg.JWTRefreshExpiry)

	claims := AuthClaims{
		UserID:    user.ID,
		Role:      user.Role,
		TokenType: tokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    refreshTokenIssuer,
			Audience:  jwt.ClaimStrings{tokenAudience},
			Subject:   user.ID.String(),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiry),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	raw, err := token.SignedString(s.cfg.JWTSigningKey)
	if err != nil {
		return "", nil, err
	}

	tokenHash := hashToken(raw)
	refreshToken := &models.RefreshToken{
		UserID:    user.ID,
		TokenHash: tokenHash,
		ExpiresAt: expiry,
		Revoked:   false,
		CreatedAt: now,
	}
	return raw, refreshToken, nil
}

// =====================================================
// ACCESS TOKEN VALIDATION
// =====================================================

func (s *AuthService) ParseAndValidateAccessToken(tokenString string) (*AuthClaims, error) {
	if tokenString == "" {
		return nil, ErrInvalidToken
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&AuthClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.cfg.JWTSigningKey, nil
		},
		jwt.WithIssuer(accessTokenIssuer),
		jwt.WithAudience(tokenAudience),
	)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if !token.Valid {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	if claims.TokenType != tokenTypeAccess {
		return nil, ErrInvalidToken
	}
	if claims.UserID == uuid.Nil {
		return nil, ErrInvalidToken
	}
	if claims.Subject != claims.UserID.String() {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

// =====================================================
// REFRESH TOKEN VALIDATION
// =====================================================

func (s *AuthService) parseRefreshToken(tokenString string) (*AuthClaims, error) {
	tokenString = strings.TrimSpace(tokenString)
	if tokenString == "" {
		return nil, ErrInvalidRefreshToken
	}

	token, err := jwt.ParseWithClaims(
		tokenString,
		&AuthClaims{},
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return s.cfg.JWTSigningKey, nil
		},
		jwt.WithIssuer(refreshTokenIssuer),
		jwt.WithAudience(tokenAudience),
	)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	if !token.Valid {
		return nil, ErrInvalidRefreshToken
	}

	claims, ok := token.Claims.(*AuthClaims)
	if !ok {
		return nil, ErrInvalidRefreshToken
	}

	if claims.TokenType != tokenTypeRefresh {
		return nil, ErrInvalidRefreshToken
	}
	if claims.UserID == uuid.Nil {
		return nil, ErrInvalidRefreshToken
	}
	if claims.Subject != claims.UserID.String() {
		return nil, ErrInvalidRefreshToken
	}
	return claims, nil
}

// =====================================================
// TOKEN HASH
// =====================================================

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}