package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/ebarrosjr/igreja-aberta/backend/go/database"
	"github.com/ebarrosjr/igreja-aberta/backend/go/middleware"
	"github.com/gin-gonic/gin"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type forgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type resetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required,min=8"`
	ConfirmPassword string `json:"confirmPassword" binding:"required"`
}

type logoutRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

func LoginHandler(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Dados de login invalidos"))
		return
	}

	email := strings.ToLower(strings.TrimSpace(request.Email))
	user, err := database.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao buscar usuario"))
		return
	}

	if user == nil || !database.IsActiveStatus(user.Status) || !database.VerifyPassword(request.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, JSONResponse(http.StatusUnauthorized, gin.H{}, "Email ou senha invalidos"))
		return
	}

	authResponse, err := issueAuthTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao gerar autenticacao"))
		return
	}

	now := time.Now().UTC()
	user.LastLoginAt = &now
	user.UpdatedAt = now

	if err := database.DB.Save(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao atualizar usuario"))
		return
	}

	c.JSON(http.StatusOK, JSONResponse(http.StatusOK, gin.H{
		"user":                     newUserResponse(user),
		"token":                    authResponse.Token,
		"expires_at":               authResponse.ExpiresAt,
		"refresh_token":            authResponse.RefreshToken,
		"refresh_token_expires_at": authResponse.RefreshTokenExpiresAt,
	}, "Logado com sucesso"))
}

func RefreshTokenHandler(c *gin.Context) {
	var request refreshTokenRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Dados de refresh token invalidos"))
		return
	}

	refreshToken := strings.TrimSpace(request.RefreshToken)
	tokenRecord, err := database.GetValidRefreshToken(refreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao validar refresh token"))
		return
	}
	if tokenRecord == nil {
		c.JSON(http.StatusUnauthorized, JSONResponse(http.StatusUnauthorized, gin.H{}, "Refresh token invalido"))
		return
	}

	user, err := database.GetUserByID(uint(tokenRecord.UserID))
	if err != nil || user == nil || !database.IsActiveStatus(user.Status) {
		c.JSON(http.StatusUnauthorized, JSONResponse(http.StatusUnauthorized, gin.H{}, "Usuario do refresh token nao encontrado"))
		return
	}

	if err := database.RevokeRefreshToken(tokenRecord); err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao atualizar refresh token"))
		return
	}

	authResponse, err := issueAuthTokens(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao gerar autenticacao"))
		return
	}

	c.JSON(http.StatusOK, JSONResponse(http.StatusOK, gin.H{
		"user":                     newUserResponse(user),
		"token":                    authResponse.Token,
		"expires_at":               authResponse.ExpiresAt,
		"refresh_token":            authResponse.RefreshToken,
		"refresh_token_expires_at": authResponse.RefreshTokenExpiresAt,
	}, "Token renovado com sucesso"))
}

func ForgotPasswordHandler(c *gin.Context) {
	var request forgotPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Dados de recuperacao invalidos"))
		return
	}

	email := strings.ToLower(strings.TrimSpace(request.Email))
	user, err := database.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao processar recuperacao de senha"))
		return
	}

	responseData := gin.H{}
	if user != nil && database.IsActiveStatus(user.Status) {
		resetToken, _, err := database.CreatePasswordResetToken(user.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao processar recuperacao de senha"))
			return
		}

		// Em ambiente real, este token deve ser enviado por email e nao retornado pela API.
		responseData["token"] = resetToken
	}

	c.JSON(http.StatusOK, JSONResponse(http.StatusOK, responseData, "Se o email existir, enviaremos as instrucoes para redefinicao de senha"))
}

func ResetPasswordHandler(c *gin.Context) {
	var request resetPasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Dados de redefinicao invalidos"))
		return
	}

	if request.Password != request.ConfirmPassword {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "As senhas informadas nao conferem"))
		return
	}

	resetToken := strings.TrimSpace(request.Token)
	tokenRecord, err := database.GetValidPasswordResetToken(resetToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao validar token de redefinicao"))
		return
	}
	if tokenRecord == nil {
		c.JSON(http.StatusUnauthorized, JSONResponse(http.StatusUnauthorized, gin.H{}, "Token de redefinicao invalido"))
		return
	}

	user, err := database.GetUserByID(uint(tokenRecord.UserID))
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, JSONResponse(http.StatusUnauthorized, gin.H{}, "Usuario do token de redefinicao nao encontrado"))
		return
	}

	hashedPassword, err := database.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao redefinir senha"))
		return
	}

	now := time.Now().UTC()
	user.Password = hashedPassword
	user.UpdatedAt = now
	if user.Status == "" {
		user.Status = database.UserStatusActive
	}

	if err := database.DB.Save(user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao redefinir senha"))
		return
	}

	if err := database.MarkPasswordResetTokenAsUsed(tokenRecord); err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao atualizar token de redefinicao"))
		return
	}

	if err := database.RevokeAllUserRefreshTokens(user.Id); err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao invalidar sessoes do usuario"))
		return
	}

	c.JSON(http.StatusOK, JSONResponse(http.StatusOK, gin.H{
		"user": newUserResponse(user),
	}, "Senha redefinida com sucesso"))
}

func LogoutHandler(c *gin.Context) {
	var request logoutRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Dados de logout invalidos"))
		return
	}

	authenticatedUser, ok := middleware.GetAuthenticatedUser(c)
	if !ok || authenticatedUser == nil {
		c.JSON(http.StatusUnauthorized, JSONResponse(http.StatusUnauthorized, gin.H{}, "Usuario autenticado nao encontrado"))
		return
	}

	refreshToken := strings.TrimSpace(request.RefreshToken)
	tokenRecord, err := database.GetValidRefreshTokenByUserID(refreshToken, authenticatedUser.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao processar logout"))
		return
	}
	if tokenRecord == nil {
		c.JSON(http.StatusUnauthorized, JSONResponse(http.StatusUnauthorized, gin.H{}, "Refresh token invalido"))
		return
	}

	if err := database.RevokeRefreshToken(tokenRecord); err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao processar logout"))
		return
	}

	c.JSON(http.StatusOK, JSONResponse(http.StatusOK, gin.H{}, "Logout realizado com sucesso"))
}

func issueAuthTokens(user *database.User) (*database.AuthResponse, error) {
	authResponse, err := database.CreateAuthResponse(user)
	if err != nil {
		return nil, err
	}

	if err := database.CreateRefreshToken(user.Id, authResponse.RefreshToken, authResponse.RefreshTokenExpiresAt); err != nil {
		return nil, err
	}

	return authResponse, nil
}
