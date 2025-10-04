package handler

import (
	"net/http"

	"github.com/danirisdiandita/pdf-engine/internal/model"
	"github.com/gin-gonic/gin"
)

func Generate(c *gin.Context) {
	var payload model.PDFRequest
	c.JSON(http.StatusOK, gin.H{"status": "ok", "result": payload, "language": ""})
}
