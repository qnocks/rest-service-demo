package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type errorResponse struct {
	Status    int    `json:"status"`
	Error     string `json:"error"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

func buildErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, errorResponse{
		Status:    statusCode,
		Error:     http.StatusText(statusCode),
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}
