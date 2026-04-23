package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/ebarrosjr/igreja-aberta/backend/go/database"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type createUserRequest struct {
	CongregationID int    `json:"congregationId" binding:"required"`
	Name           string `json:"name" binding:"required"`
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=8"`
	Status         string `json:"status" binding:"required"`
}

type updateUserRequest struct {
	CongregationID *int    `json:"congregationId"`
	Name           *string `json:"name"`
	Email          *string `json:"email" binding:"omitempty,email"`
	Password       *string `json:"password" binding:"omitempty,min=8"`
	Status         *string `json:"status"`
}

func ListUsersHandler(c *gin.Context) {
	users, err := database.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao listar usuarios"))
		return
	}

	responseUsers := make([]userResponse, 0, len(users))
	for _, user := range users {
		userCopy := user
		responseUsers = append(responseUsers, newUserResponse(&userCopy))
	}

	c.JSON(http.StatusOK, JSONResponse(http.StatusOK, gin.H{
		"user": responseUsers,
	}, "Usuarios listados com sucesso"))
}

func CreateUserHandler(c *gin.Context) {
	var request createUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Dados do usuario invalidos"))
		return
	}
	if request.CongregationID <= 0 {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "CongregationId deve ser maior que zero"))
		return
	}
	if strings.TrimSpace(request.Name) == "" {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Nome do usuario invalido"))
		return
	}
	if strings.TrimSpace(request.Password) == "" {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Senha do usuario invalida"))
		return
	}

	if !database.IsValidUserStatus(request.Status) {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Status do usuario invalido"))
		return
	}

	email := strings.ToLower(strings.TrimSpace(request.Email))
	exists, err := database.UserExistsByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao validar email do usuario"))
		return
	}
	if exists {
		c.JSON(http.StatusConflict, JSONResponse(http.StatusConflict, gin.H{}, "Ja existe um usuario com este email"))
		return
	}

	hashedPassword, err := database.HashPassword(request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao criptografar senha"))
		return
	}

	user, err := database.CreateUser(
		request.CongregationID,
		strings.TrimSpace(request.Name),
		email,
		hashedPassword,
		strings.TrimSpace(request.Status),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao criar usuario"))
		return
	}

	c.JSON(http.StatusCreated, JSONResponse(http.StatusCreated, gin.H{
		"user": newUserResponse(user),
	}, "Usuario criado com sucesso"))
}

func GetUserHandler(c *gin.Context) {
	user, err := findUserByParamID(c.Param("id"))
	if err != nil {
		handleUserLookupError(c, err)
		return
	}

	c.JSON(http.StatusOK, JSONResponse(http.StatusOK, gin.H{
		"user": newUserResponse(user),
	}, "Usuario encontrado com sucesso"))
}

func UpdateUserHandler(c *gin.Context) {
	user, err := findUserByParamID(c.Param("id"))
	if err != nil {
		handleUserLookupError(c, err)
		return
	}

	var request updateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Dados do usuario invalidos"))
		return
	}

	if request.CongregationID != nil {
		if *request.CongregationID <= 0 {
			c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "CongregationId deve ser maior que zero"))
			return
		}
		user.CongregationId = *request.CongregationID
	}
	if request.Name != nil {
		name := strings.TrimSpace(*request.Name)
		if name == "" {
			c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Nome do usuario invalido"))
			return
		}
		user.Name = name
	}
	if request.Email != nil {
		email := strings.ToLower(strings.TrimSpace(*request.Email))
		if email == "" {
			c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Email do usuario invalido"))
			return
		}
		exists, err := database.UserExistsByEmailExceptID(email, uint(user.Id))
		if err != nil {
			c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao validar email do usuario"))
			return
		}
		if exists {
			c.JSON(http.StatusConflict, JSONResponse(http.StatusConflict, gin.H{}, "Ja existe um usuario com este email"))
			return
		}
		user.Email = email
	}
	if request.Password != nil {
		password := strings.TrimSpace(*request.Password)
		if password == "" {
			c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Senha do usuario invalida"))
			return
		}
		hashedPassword, err := database.HashPassword(password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao criptografar senha"))
			return
		}
		user.Password = hashedPassword
	}
	if request.Status != nil {
		status := database.NormalizeUserStatus(*request.Status)
		if !database.IsValidUserStatus(status) {
			c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Status do usuario invalido"))
			return
		}
		user.Status = status
	}

	if err := database.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao atualizar usuario"))
		return
	}

	c.JSON(http.StatusOK, JSONResponse(http.StatusOK, gin.H{
		"user": newUserResponse(user),
	}, "Usuario atualizado com sucesso"))
}

func DeleteUserHandler(c *gin.Context) {
	user, err := findUserByParamID(c.Param("id"))
	if err != nil {
		handleUserLookupError(c, err)
		return
	}

	if err := database.DeleteUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, JSONResponse(http.StatusInternalServerError, gin.H{}, "Erro ao remover usuario"))
		return
	}

	c.JSON(http.StatusOK, JSONResponse(http.StatusOK, gin.H{
		"user": newUserResponse(user),
	}, "Usuario removido com sucesso"))
}

func findUserByParamID(rawID string) (*database.User, error) {
	id, err := strconv.ParseUint(strings.TrimSpace(rawID), 10, 64)
	if err != nil {
		return nil, err
	}

	user, err := database.GetUserByID(uint(id))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func handleUserLookupError(c *gin.Context, err error) {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusNotFound, JSONResponse(http.StatusNotFound, gin.H{}, "Usuario nao encontrado"))
		return
	}

	c.JSON(http.StatusBadRequest, JSONResponse(http.StatusBadRequest, gin.H{}, "Identificador de usuario invalido"))
}
