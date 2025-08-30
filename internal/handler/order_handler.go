package handler

import (
	"net/http"
	"strconv"

	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"github.com/Incipe-win/ai-tshirt-shop/internal/repository"
	"github.com/Incipe-win/ai-tshirt-shop/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var orderService *service.OrderService

func InitOrderHandler(db *gorm.DB) {
	orderRepo := repository.NewOrderRepository(db)
	cartRepo := repository.NewCartRepository(db)
	orderService = service.NewOrderService(orderRepo, cartRepo)
}

// CreateOrder godoc
// @Summary 创建订单
// @Description 从购物车创建订单，将购物车中的商品转换为订单
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateOrderRequest true "创建订单请求参数"
// @Success 201 {object} service.OrderResponse "订单创建成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /orders [post]
func CreateOrder(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	var req service.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	order, err := orderService.CreateOrder(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to create order",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"order":   order,
		"message": "Order created successfully",
	})
}

// GetUserOrders godoc
// @Summary 获取用户订单列表
// @Description 获取当前用户的所有订单
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "订单列表"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /orders [get]
func GetUserOrders(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	orders, err := orderService.GetOrdersByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch orders",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"orders":  orders,
		"count":   len(orders),
		"message": "Orders fetched successfully",
	})
}

// GetOrderByID godoc
// @Summary 根据ID获取订单
// @Description 根据订单ID获取订单详细信息
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Success 200 {object} service.OrderResponse "订单详情"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "订单不存在"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /orders/{id} [get]
func GetOrderByID(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	order, err := orderService.GetOrderByID(userID, uint(orderID))
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}
		if err.Error() == "order does not belong to user" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch order",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order":   order,
		"message": "Order fetched successfully",
	})
}

// GetOrderByOrderSN godoc
// @Summary 根据订单号获取订单
// @Description 根据订单号获取订单详细信息
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param order_sn path string true "订单号"
// @Success 200 {object} service.OrderResponse "订单详情"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "订单不存在"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /orders/sn/{order_sn} [get]
func GetOrderByOrderSN(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	orderSN := c.Param("order_sn")
	if orderSN == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Order SN is required",
		})
		return
	}

	order, err := orderService.GetOrderByOrderSN(userID, orderSN)
	if err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}
		if err.Error() == "order does not belong to user" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch order",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"order":   order,
		"message": "Order fetched successfully",
	})
}

// UpdateOrderStatus godoc
// @Summary 更新订单状态
// @Description 更新订单状态（支付、发货、完成等）
// @Tags orders
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "订单ID"
// @Param request body map[string]string true "状态更新请求参数"
// @Success 200 {object} map[string]interface{} "状态更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "订单不存在"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /orders/{id}/status [put]
func UpdateOrderStatus(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User ID not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID format",
		})
		return
	}

	idStr := c.Param("id")
	orderID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 验证状态值
	status := model.OrderStatus(req.Status)
	validStatuses := []model.OrderStatus{
		model.OrderStatusPending,
		model.OrderStatusPaid,
		model.OrderStatusShipped,
		model.OrderStatusCompleted,
		model.OrderStatusCancelled,
	}

	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order status",
		})
		return
	}

	if err := orderService.UpdateOrderStatus(userID, uint(orderID), status); err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Order not found",
			})
			return
		}
		if err.Error() == "order does not belong to user" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update order status",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
		"status":  status,
	})
}