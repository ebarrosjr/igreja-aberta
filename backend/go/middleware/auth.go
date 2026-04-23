package middleware

import (
	"net/http"
	"strings"

	"github.com/ebarrosjr/igreja-aberta/backend/go/database"
	"github.com/gin-gonic/gin"
)

const AuthenticatedUserContextKey = "authenticatedUser"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, ok := extractBearerToken(c.GetHeader("Authorization"))
		if !ok {
			c.JSON(http.StatusUnauthorized, jsonResponse(http.StatusUnauthorized, gin.H{}, "Token de autenticacao nao informado"))
			c.Abort()
			return
		}

		claims, err := database.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, jsonResponse(http.StatusUnauthorized, gin.H{}, "Token de autenticacao invalido"))
			c.Abort()
			return
		}

		userID, err := database.GetUserIDFromClaims(claims)
		if err != nil {
			c.JSON(http.StatusUnauthorized, jsonResponse(http.StatusUnauthorized, gin.H{}, "Token de autenticacao invalido"))
			c.Abort()
			return
		}

		user, err := database.GetUserByID(userID)
		if err != nil || user == nil || !database.IsActiveStatus(user.Status) {
			c.JSON(http.StatusUnauthorized, jsonResponse(http.StatusUnauthorized, gin.H{}, "Usuario autenticado nao encontrado"))
			c.Abort()
			return
		}

		c.Set(AuthenticatedUserContextKey, user)
		c.Next()
	}
}

func GetAuthenticatedUser(c *gin.Context) (*database.User, bool) {
	value, exists := c.Get(AuthenticatedUserContextKey)
	if !exists {
		return nil, false
	}

	user, ok := value.(*database.User)
	return user, ok
}

func extractBearerToken(header string) (string, bool) {
	scheme, token, found := strings.Cut(strings.TrimSpace(header), " ")
	if !found || !strings.EqualFold(scheme, "Bearer") {
		return "", false
	}

	token = strings.TrimSpace(token)
	if token == "" {
		return "", false
	}

	return token, true
}

func jsonResponse(code int, data any, message string) gin.H {
	return gin.H{
		"code":    code,
		"data":    data,
		"message": message,
		"stack":   "Go",
		"db":      "SQLite",
	}
}
