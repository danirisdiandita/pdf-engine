package router

import (
	"github.com/danirisdiandita/pdf-engine/internal/config"
	"github.com/danirisdiandita/pdf-engine/internal/handler"
	"github.com/gin-gonic/gin"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.POST("/generate", handler.Generate)
		api.GET("/hello", func(c *gin.Context) {
			name := c.DefaultQuery("name", "World")
			c.JSON(200, gin.H{"message": "Hello " + name + " You are my world"})
		})
	}

	return r
}
