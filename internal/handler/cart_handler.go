package handler

import (
	"net/http"
	"strconv"

	"github.com/Incipe-win/ai-tshirt-shop/internal/repository"
	"github.com/Incipe-win/ai-tshirt-shop/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var cartService *service.CartService

func InitCartHandler(db *gorm.DB) {
	cartRepo := repository.NewCartRepository(db)
	productRepo := repository.NewProductRepository(db)
	designRepo := repository.NewDesignRepository(db)
	cartService = service.NewCartService(cartRepo, productRepo, designRepo)
}

// AddToCart godoc
// @Summary 添加商品到购物车
// @Description 将指定的商品和设计组合添加到用户购物车
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.AddToCartRequest true "添加到购物车请求参数"
// @Success 200 {object} map[string]interface{} "添加成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /cart/add [post]
func AddToCart(c *gin.Context) {
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

	var req service.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := cartService.AddToCart(userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to add item to cart",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item added to cart successfully",
	})
}

// GetCart godoc
// @Summary 获取购物车
// @Description 获取当前用户的购物车内容
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} service.CartResponse "购物车内容"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /cart [get]
func GetCart(c *gin.Context) {
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

	cart, err := cartService.GetCart(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch cart",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, cart)
}

// UpdateCartItem godoc
// @Summary 更新购物车商品数量
// @Description 更新购物车中指定商品的数量
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "购物车商品ID"
// @Param request body service.UpdateCartRequest true "更新购物车请求参数"
// @Success 200 {object} map[string]interface{} "更新成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "购物车商品不存在"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /cart/{id} [put]
func UpdateCartItem(c *gin.Context) {
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
	itemID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cart item ID",
		})
		return
	}

	var req service.UpdateCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if err := cartService.UpdateCartItem(userID, uint(itemID), &req); err != nil {
		if err.Error() == "cart item not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Cart item not found",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Failed to update cart item",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart item updated successfully",
	})
}

// RemoveFromCart godoc
// @Summary 从购物车删除商品
// @Description 从购物车中删除指定的商品
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "购物车商品ID"
// @Success 200 {object} map[string]interface{} "删除成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /cart/{id} [delete]
func RemoveFromCart(c *gin.Context) {
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
	itemID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid cart item ID",
		})
		return
	}

	if err := cartService.RemoveFromCart(userID, uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to remove item from cart",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Item removed from cart successfully",
	})
}

// ClearCart godoc
// @Summary 清空购物车
// @Description 清空当前用户的整个购物车
// @Tags cart
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "清空成功"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /cart/clear [delete]
func ClearCart(c *gin.Context) {
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

	if err := cartService.ClearCart(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to clear cart",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Cart cleared successfully",
	})
}