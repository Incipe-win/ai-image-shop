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
		api.GET("/creatives", getCreatives)

		// 用户认证路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", Register)
			auth.POST("/login", Login)
			auth.POST("/refresh", RefreshToken)
		}

		// 商品路由
		products := api.Group("/products")
		{
			products.GET("/", GetAllProducts)
			products.GET("/:id", GetProductByID)
			products.GET("/category", GetProductsByCategory)
			products.POST("/", CreateProduct) // 管理员功能，后续可加权限控制
		}

		// 需要JWT认证的路由
		protected := api.Group("")
		protected.Use(middleware.JWTMiddleware())
		{
			// 设计路由
			designs := protected.Group("/designs")
			{
				designs.POST("/generate", GenerateDesign)
				designs.GET("/my-designs", GetUserDesigns)
				designs.POST("/publish", PublishDesignToShop)
			}

			// 购物车路由
			cart := protected.Group("/cart")
			{
				cart.POST("/add", AddToCart)
				cart.GET("/", GetCart)
				cart.PUT("/:id", UpdateCartItem)
				cart.DELETE("/:id", RemoveFromCart)
				cart.DELETE("/clear", ClearCart)
			}

			// 订单路由
			orders := protected.Group("/orders")
			{
				orders.POST("/", CreateOrder)
				orders.GET("/", GetUserOrders)
				orders.GET("/:id", GetOrderByID)
				orders.GET("/sn/:order_sn", GetOrderByOrderSN)
				orders.PUT("/:id/status", UpdateOrderStatus)
			}
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

// GetCreatives godoc
// @Summary 获取创意产品列表
// @Description 获取所有可用的创意产品列表
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "{"message": "Creative products endpoint"}"
// @Router /creatives [get]
func getCreatives(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Creative products endpoint",
	})
}
