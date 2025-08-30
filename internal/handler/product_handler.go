package handler

import (
	"net/http"
	"strconv"

	"github.com/Incipe-win/ai-tshirt-shop/internal/repository"
	"github.com/Incipe-win/ai-tshirt-shop/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var productService *service.ProductService

func InitProductHandler(db *gorm.DB) {
	productRepo := repository.NewProductRepository(db)
	productService = service.NewProductService(productRepo)
}

// CreateProduct godoc
// @Summary 创建商品
// @Description 创建新的创意商品模板
// @Tags products
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.CreateProductRequest true "创建商品请求参数"
// @Success 201 {object} service.ProductResponse "商品创建成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /products [post]
func CreateProduct(c *gin.Context) {
	var req service.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	product, err := productService.CreateProduct(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create product",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"product": product,
		"message": "Product created successfully",
	})
}

// GetAllProducts godoc
// @Summary 获取所有商品
// @Description 获取所有可用的创意商品模板
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "商品列表"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /products [get]
func GetAllProducts(c *gin.Context) {
	products, err := productService.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch products",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"count":    len(products),
		"message":  "Products fetched successfully",
	})
}

// GetProductByID godoc
// @Summary 根据ID获取商品
// @Description 根据商品ID获取特定的创意商品模板信息
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "商品ID"
// @Success 200 {object} service.ProductResponse "商品信息"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 404 {object} map[string]interface{} "商品不存在"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /products/{id} [get]
func GetProductByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID",
		})
		return
	}

	product, err := productService.GetProductByID(uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Product not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch product",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": product,
		"message": "Product fetched successfully",
	})
}

// GetProductsByCategory godoc
// @Summary 根据分类获取商品
// @Description 根据商品分类获取创意商品列表
// @Tags products
// @Accept json
// @Produce json
// @Param category query string true "商品分类"
// @Success 200 {object} map[string]interface{} "商品列表"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /products/category [get]
func GetProductsByCategory(c *gin.Context) {
	category := c.Query("category")
	if category == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Category parameter is required",
		})
		return
	}

	products, err := productService.GetProductsByCategory(category)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch products",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"category": category,
		"count":    len(products),
		"message":  "Products fetched successfully",
	})
}