package service

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"

    "restaurant-backend/internal/config"
    "restaurant-backend/internal/models"
    "restaurant-backend/internal/repository"
)

var (
    ErrInvalidCredentials = errors.New("invalid credentials")
    ErrInvalidToken       = errors.New("invalid token")
)

type AuthService struct {
    userRepo       *repository.UserRepository
    rtRepo         *repository.RefreshTokenRepository
    cfg            *config.Config
}

func NewAuthService(userRepo *repository.UserRepository, rtRepo *repository.RefreshTokenRepository, cfg *config.Config) *AuthService {
    return &AuthService{
        userRepo: userRepo,
        rtRepo:   rtRepo,
        cfg:      cfg,
    }
}

// Password helpers

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Register

type RegisterInput struct {
    Name     string
    Email    string
    Password string
    Role     models.UserRole
}

func (s *AuthService) Register(in RegisterInput) (*models.User, error) {
    if _, err := s.userRepo.GetByEmail(in.Email); err == nil {
        return nil, errors.New("email already registered")
    }

    hash, err := HashPassword(in.Password)
    if err != nil {
        return nil, err
    }

    user := &models.User{
        Name:         in.Name,
        Email:        in.Email,
        PasswordHash: hash,
        Role:         in.Role,
    }

    if err := s.userRepo.Create(user); err != nil {
        return nil, err
    }

    return user, nil
}

// Login

type LoginInput struct {
    Email    string
    Password string
}

type LoginOutput struct {
    User       *models.User
    AccessToken string
    RefreshTokenRaw string
    RefreshToken *models.RefreshToken
}

func (s *AuthService) Login(in LoginInput) (*LoginOutput, error) {
    user, err := s.userRepo.GetByEmail(in.Email)
    if err != nil {
        return nil, ErrInvalidCredentials
    }

    if !CheckPassword(in.Password, user.PasswordHash) {
        return nil, ErrInvalidCredentials
    }

    accessToken, err := s.newAccessToken(user)
    if err != nil {
        return nil, err
    }

    refreshRaw, rt, err := s.newRefreshTokenRaw(user)
    if err != nil {
        return nil, err
    }

    if err := s.rtRepo.Create(rt); err != nil {
        return nil, err
    }

    return &LoginOutput{
        User:            user,
        AccessToken:     accessToken,
        RefreshTokenRaw: refreshRaw,
        RefreshToken:    rt,
    }, nil
}

// Refresh

type RefreshOutput struct {
    AccessToken     string
    RefreshTokenRaw string
    RefreshToken    *models.RefreshToken
}

func (s *AuthService) Refresh(rawToken string) (*RefreshOutput, error) {
    _, err := s.ParseAndValidateAccessToken(rawToken)
    if err != nil {
        return nil, ErrInvalidToken
    }

    h := sha256.Sum256([]byte(rawToken))
    tokenHash := hex.EncodeToString(h[:])

    rt, err := s.rtRepo.GetByHashAndNotRevoked(tokenHash)
    if err != nil {
        return nil, ErrInvalidToken
    }

    if time.Now().After(rt.ExpiresAt) {
        return nil, ErrInvalidToken
    }

    // Revoke old
    _ = s.rtRepo.Revoke(rt)

    user, err := s.userRepo.GetByID(rt.UserID)
    if err != nil {
        return nil, err
    }

    newAccess, _ := s.newAccessToken(user)
    newRefreshRaw, newRT, _ := s.newRefreshTokenRaw(user)

    if err := s.rtRepo.Create(newRT); err != nil {
        return nil, err
    }

    return &RefreshOutput{
        AccessToken:     newAccess,
        RefreshTokenRaw: newRefreshRaw,
        RefreshToken:    newRT,
    }, nil
}

// Logout

func (s *AuthService) Logout(rawToken string) error {
    h := sha256.Sum256([]byte(rawToken))
    tokenHash := hex.EncodeToString(h[:])

    rt, err := s.rtRepo.GetByHash(tokenHash)
    if err != nil {
        return nil // ignore if not found
    }

    _ = s.rtRepo.Revoke(rt)
    return nil
}

// Get user by ID (for /me)

func (s *AuthService) GetUserByID(id int64) (*models.User, error) {
    return s.userRepo.GetByID(id)
}

// JWT helpers

type AuthClaims struct {
    UserID int64        `json:"user_id"`
    Role   models.UserRole `json:"role"`
    jwt.RegisteredClaims
}

func (s *AuthService) newAccessToken(user *models.User) (string, error) {
    now := time.Now()
    claims := AuthClaims{
        UserID: user.ID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    "https://api.yourrestaurant.com",
            Audience:  jwt.ClaimStrings{"restaurant-api"},
            Subject:   fmt.Sprintf("%d", user.ID),
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(now.Add(s.cfg.JWTAccessExpiry)),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.cfg.JWTSigningKey)
}

func (s *AuthService) newRefreshTokenRaw(user *models.User) (string, *models.RefreshToken, error) {
    now := time.Now()
    exp := now.Add(s.cfg.JWTRefreshExpiry)

    claims := AuthClaims{
        UserID: user.ID,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            Issuer:    "https://api.yourrestaurant.com",
            Audience:  jwt.ClaimStrings{"restaurant-api"},
            Subject:   fmt.Sprintf("%d", user.ID),
            IssuedAt:  jwt.NewNumericDate(now),
            ExpiresAt: jwt.NewNumericDate(exp),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    raw, err := token.SignedString(s.cfg.JWTSigningKey)
    if err != nil {
        return "", nil, err
    }

    h := sha256.Sum256([]byte(raw))
    tokenHash := hex.EncodeToString(h[:])

    rt := &models.RefreshToken{
        UserID:    user.ID,
        TokenHash: tokenHash,
        ExpiresAt: exp,
        Revoked:   false,
        CreatedAt: now,
    }

    return raw, rt, nil
}

func (s *AuthService) ParseAndValidateAccessToken(tokenStr string) (*AuthClaims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &AuthClaims{}, func(token *jwt.Token) (any, error) {
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
        return s.cfg.JWTSigningKey, nil
    })
    if err != nil || !token.Valid {
        return nil, ErrInvalidToken
    }

    claims, ok := token.Claims.(*AuthClaims)
    if !ok {
        return nil, ErrInvalidToken
    }

    if claims.UserID <= 0 {
        return nil, ErrInvalidToken
    }

    return claims, nil
}
