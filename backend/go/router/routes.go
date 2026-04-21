package router

import (
	"net/http"

	"github.com/ebarrosjr/igreja-aberta/backend/go/handler"
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

		r.POST("/login", handler.LoginHandler)
	}
}
