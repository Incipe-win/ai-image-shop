package handler

import (
	"net/http"
	"time"

	"github.com/Incipe-win/ai-tshirt-shop/pkg/logger"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.Use(gin.Recovery())
	r.Use(loggingMiddleware())

	setupRoutes(r)

	return r
}

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		latency := time.Since(start)

		logger.Info("HTTP request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency", latency.String(),
		)
	}
}

func setupRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		api.GET("/health", healthCheck)
		api.GET("/tshirts", getTshirts)
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

func getTshirts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "T-shirts endpoint",
	})
}