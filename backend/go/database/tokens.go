package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type RefreshToken struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    int        `json:"user_id" gorm:"not null;index"`
	TokenHash string     `json:"token_hash" gorm:"size:64;not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	RevokedAt *time.Time `json:"revoked_at"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;autoCreateTime"`
}

type PasswordResetToken struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	UserID    int        `json:"user_id" gorm:"not null;index"`
	TokenHash string     `json:"token_hash" gorm:"size:64;not null;uniqueIndex"`
	ExpiresAt time.Time  `json:"expires_at" gorm:"not null"`
	UsedAt    *time.Time `json:"used_at"`
	CreatedAt time.Time  `json:"created_at" gorm:"not null;autoCreateTime"`
}

type AuthResponse struct {
	Token                 string    `json:"token"`
	ExpiresAt             time.Time `json:"expires_at"`
	RefreshToken          string    `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time `json:"refresh_token_expires_at"`
}

func CreateAuthResponse(user *User) (*AuthResponse, error) {
	expiresAt := time.Now().UTC().Add(time.Duration(getAccessTokenMinutes()) * time.Minute)
	refreshToken, err := GenerateOpaqueToken()
	if err != nil {
		return nil, err
	}

	token, err := GenerateJWT(user, expiresAt)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token:                 token,
		ExpiresAt:             expiresAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: time.Now().UTC().AddDate(0, 0, getRefreshTokenDays()),
	}, nil
}

func GenerateOpaqueToken() (string, error) {
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(bytes), nil
}

func HashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func GenerateJWT(user *User, expiresAt time.Time) (string, error) {
	key, err := getJWTKey()
	if err != nil {
		return "", err
	}

	claims := jwt.MapClaims{
		"sub":             fmt.Sprintf("%d", user.Id),
		"email":           user.Email,
		"congregation_id": fmt.Sprintf("%d", user.CongregationId),
		"iss":             getJWTIssuer(),
		"aud":             getJWTAudience(),
		"exp":             expiresAt.Unix(),
		"iat":             time.Now().UTC().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(key))
}

func getJWTKey() (string, error) {
	key := os.Getenv("JWT_KEY")
	if key == "" {
		key = "CHANGE_ME_TO_A_32_BYTE_OR_LONGER_SECRET_KEY"
	}

	if len([]byte(key)) < 32 {
		return "", fmt.Errorf("JWT_KEY must be at least 32 bytes")
	}

	return key, nil
}

func getJWTIssuer() string {
	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		return "igreja-aberta"
	}

	return issuer
}

func getJWTAudience() string {
	audience := os.Getenv("JWT_AUDIENCE")
	if audience == "" {
		return "igreja-aberta-api"
	}

	return audience
}

func getAccessTokenMinutes() int {
	value := os.Getenv("JWT_ACCESS_TOKEN_MINUTES")
	minutes, err := strconv.Atoi(value)
	if err != nil || minutes <= 0 {
		return 15
	}

	return minutes
}

func getRefreshTokenDays() int {
	value := os.Getenv("JWT_REFRESH_TOKEN_DAYS")
	days, err := strconv.Atoi(value)
	if err != nil || days <= 0 {
		return 7
	}

	return days
}
