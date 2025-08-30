package handler

import (
	"net/http"
	"time"

	"github.com/Incipe-win/ai-tshirt-shop/internal/middleware"
	"github.com/Incipe-win/ai-tshirt-shop/pkg/logger"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)

func InitRouter(env string) *gin.Engine {
	if env == "production" {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
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
	// 静态文件服务
	r.Static("/static", "./static")
	r.Static("/images", "./static/images")

	// Swagger文档路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 前端页面路由
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	api := r.Group("/api/v1")
	{
		api.GET("/health", healthCheck)
		api.GET("/tshirts", getTshirts)

		// 用户认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", Register)
			auth.POST("/login", Login)
			auth.POST("/refresh", RefreshToken)
		}

		// 需要JWT认证的设计生成路由
		protected := api.Group("/designs")
		protected.Use(middleware.JWTMiddleware())
		{
			protected.POST("/generate", GenerateDesign)
			protected.GET("/my-designs", GetUserDesigns)
		}
	}
}

// HealthCheck godoc
// @Summary 健康检查
// @Description 检查API服务状态
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "{“status”: “ok”, “time”: “2023-01-01T00:00:00Z”}"
// @Router /health [get]
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	})
}

// GetTshirts godoc
// @Summary 获取T恤列表
// @Description 获取所有可用的T恤产品列表
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "{“message”: “T-shirts endpoint”}"
// @Router /tshirts [get]
func getTshirts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "T-shirts endpoint",
	})
}
