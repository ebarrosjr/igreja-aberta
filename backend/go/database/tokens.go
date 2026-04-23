package database

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
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

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	key, err := getJWTKey()
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(key), nil
	}, jwt.WithIssuer(getJWTIssuer()), jwt.WithAudience(getJWTAudience()))
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func GetUserIDFromClaims(claims jwt.MapClaims) (uint, error) {
	subject, ok := claims["sub"].(string)
	if !ok || strings.TrimSpace(subject) == "" {
		return 0, fmt.Errorf("invalid token subject")
	}

	id, err := strconv.ParseUint(subject, 10, 64)
	if err != nil {
		return 0, err
	}

	return uint(id), nil
}

func CreateRefreshToken(userID int, refreshToken string, expiresAt time.Time) error {
	token := &RefreshToken{
		UserID:    userID,
		TokenHash: HashToken(refreshToken),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now().UTC(),
	}

	return DB.Create(token).Error
}

func GetValidRefreshToken(refreshToken string) (*RefreshToken, error) {
	var token RefreshToken
	now := time.Now().UTC()

	result := DB.Where("token_hash = ? AND revoked_at IS NULL AND expires_at > ?", HashToken(refreshToken), now).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &token, nil
}

func GetValidRefreshTokenByUserID(refreshToken string, userID int) (*RefreshToken, error) {
	var token RefreshToken
	now := time.Now().UTC()

	result := DB.Where("token_hash = ? AND user_id = ? AND revoked_at IS NULL AND expires_at > ?", HashToken(refreshToken), userID, now).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &token, nil
}

func RevokeRefreshToken(token *RefreshToken) error {
	now := time.Now().UTC()
	token.RevokedAt = &now
	return DB.Save(token).Error
}

func RevokeRefreshTokenByRawToken(refreshToken string) error {
	token, err := GetValidRefreshToken(refreshToken)
	if err != nil || token == nil {
		return err
	}

	return RevokeRefreshToken(token)
}

func RevokeAllUserRefreshTokens(userID int) error {
	now := time.Now().UTC()
	return DB.Model(&RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now).Error
}

func CreatePasswordResetToken(userID int) (string, *PasswordResetToken, error) {
	resetToken, err := GenerateOpaqueToken()
	if err != nil {
		return "", nil, err
	}

	record := &PasswordResetToken{
		UserID:    userID,
		TokenHash: HashToken(resetToken),
		ExpiresAt: time.Now().UTC().Add(time.Duration(getPasswordResetTokenMinutes()) * time.Minute),
		CreatedAt: time.Now().UTC(),
	}

	if err := DB.Create(record).Error; err != nil {
		return "", nil, err
	}

	return resetToken, record, nil
}

func GetValidPasswordResetToken(resetToken string) (*PasswordResetToken, error) {
	var token PasswordResetToken
	now := time.Now().UTC()

	result := DB.Where("token_hash = ? AND used_at IS NULL AND expires_at > ?", HashToken(resetToken), now).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}

	return &token, nil
}

func MarkPasswordResetTokenAsUsed(token *PasswordResetToken) error {
	now := time.Now().UTC()
	token.UsedAt = &now
	return DB.Save(token).Error
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

func getPasswordResetTokenMinutes() int {
	value := os.Getenv("PASSWORD_RESET_TOKEN_MINUTES")
	minutes, err := strconv.Atoi(value)
	if err != nil || minutes <= 0 {
		return 30
	}

	return minutes
}
