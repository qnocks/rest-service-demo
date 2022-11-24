package http

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) backup(c *gin.Context) {
	backupFilename, err := h.service.Backup()
	if err != nil {
		buildErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.File(backupFilename)
}
