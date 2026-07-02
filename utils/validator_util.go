package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ValidationError(ctx *gin.Context, err error) {
	ctx.JSON(http.StatusBadRequest, gin.H{
		"error": err.Error(),
	})
}

func SuccessResponse(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, data)
}
