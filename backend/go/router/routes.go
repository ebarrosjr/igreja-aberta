package router

import (
	"net/http"

	"github.com/ebarrosjr/igreja-aberta/backend/go/handler"
	"github.com/ebarrosjr/igreja-aberta/backend/go/middleware"
	gin "github.com/gin-gonic/gin"
)

func InitializeRoutes(router *gin.Engine) {

	r := router.Group("/api")
	{
		r.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, handler.JSONResponse(
				http.StatusOK,
				"Bem vindo à API da Igreja Aberta em GoLang!",
				"pong",
			))
		})

		auth := r.Group("/Auth")
		auth.POST("/login", handler.LoginHandler)
		auth.POST("/refresh-token", handler.RefreshTokenHandler)
		auth.POST("/forgot-password", handler.ForgotPasswordHandler)
		auth.POST("/reset-password", handler.ResetPasswordHandler)
		auth.Use(middleware.AuthMiddleware())
		{
			auth.POST("/logout", handler.LogoutHandler)
		}

		users := r.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("", handler.ListUsersHandler)
			users.POST("", handler.CreateUserHandler)
			users.GET("/:id", handler.GetUserHandler)
			users.PUT("/:id", handler.UpdateUserHandler)
			users.DELETE("/:id", handler.DeleteUserHandler)
		}
	}
}
