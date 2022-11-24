package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"team-task/internal/dto"
)

func (h *Handler) set(c *gin.Context) {
	var requestBody dto.UserGrade
	if err := c.BindJSON(&requestBody); err != nil {
		buildErrorResponse(c, http.StatusBadRequest, "error parsing request body")
		return
	}

	userGrade, err := h.service.Set(requestBody)
	if err != nil {
		buildErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, userGrade)
}

func (h *Handler) get(c *gin.Context) {
	userID := c.Query("user_id")

	if len(userID) == 0 {
		buildErrorResponse(c, http.StatusBadRequest, "missing parameter [user_id]")
		return
	}

	userGrade, err := h.service.Get(userID)
	if err != nil {
		buildErrorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, userGrade)
}
