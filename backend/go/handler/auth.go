package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/ebarrosjr/igreja-aberta/backend/go/database"
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		}, "Dados de login invalidos"))
		return
	}

	email := strings.ToLower(strings.TrimSpace(request.Email))
	user, err := database.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		}, "Erro ao buscar usuario"))
		return
	}

	if user == nil || user.Status != "active" || !database.VerifyPassword(request.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, JSONResponse(http.StatusUnauthorized, gin.H{}, "Email ou senha invalidos"))
		return
	}

	authResponse, err := database.CreateAuthResponse(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		}, "Erro ao gerar autenticacao"))
		return
	}

	now := time.Now().UTC()
	user.LastLoginAt = &now
	user.UpdatedAt = now

	refreshToken := &database.RefreshToken{
		UserID:    user.Id,
		TokenHash: database.HashToken(authResponse.RefreshToken),
		ExpiresAt: authResponse.RefreshTokenExpiresAt,
		CreatedAt: time.Now().UTC(),
	}

	if err := database.DB.Save(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		}, "Erro ao atualizar usuario"))
		return
	}

	if err := database.DB.Create(refreshToken).Error; err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		}, "Erro ao salvar refresh token"))
		return
	}

	c.JSON(http.StatusOK, JSONResponse(http.StatusOK, gin.H{
		"user": gin.H{
			"id":              user.Id,
			"congregation_id": user.CongregationId,
			"name":            user.Name,
			"email":           user.Email,
			"status":          user.Status,
			"last_login_at":   user.LastLoginAt,
		},
		"token":                    authResponse.Token,
		"expires_at":               authResponse.ExpiresAt,
		"refresh_token":            authResponse.RefreshToken,
		"refresh_token_expires_at": authResponse.RefreshTokenExpiresAt,
	}, "Logado com sucesso"))
}
