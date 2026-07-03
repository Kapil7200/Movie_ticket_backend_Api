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

type ResponseBody struct {
	Message    string `json:"message"`
	StatusCode int    `json:"Status_code"`
	// DevMessage error       `json:"-"`
	DevMessage string      `json:"dev_message,omitempty"`
	Body       interface{} `json:"body,omitempty"`
}
func InternalServerErrorResponse(c *gin.Context, err error) {
	response := ResponseBody{
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
		DevMessage: err.Error(),
	}
	c.JSON(http.StatusInternalServerError, response)
}