package handler

import (
	"github.com/gin-gonic/gin"
)

var stack = gin.H{
	"stack": "Go",
	"db":    "SQLite",
}

func JSONResponse(code int, data any, message string) gin.H {
	response := gin.H{
		"code":    code,
		"data":    data,
		"message": message,
	}

	for key, value := range stack {
		response[key] = value
	}

	return response
}
