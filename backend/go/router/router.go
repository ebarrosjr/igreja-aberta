package router

import gin "github.com/gin-gonic/gin"

func InitializeRouter() {
	r := gin.Default()
	InitializeRoutes(r)

	r.Run(":8080")
}
