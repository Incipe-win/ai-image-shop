package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Incipe-win/ai-tshirt-shop/internal/model"
	"github.com/Incipe-win/ai-tshirt-shop/internal/repository"
	"github.com/Incipe-win/ai-tshirt-shop/internal/service"
	"github.com/Incipe-win/ai-tshirt-shop/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	aiService        *service.AIService
	designRepository *repository.DesignRepository
	productRepository *repository.ProductRepository
	userRepository   *repository.UserRepository
)

func init() {
	aiService = service.NewAIService()

	if err := os.MkdirAll("./static/images", 0755); err != nil {
		logger.Error("Failed to create static images directory", err)
	}
}

func InitDesignRepository(db *gorm.DB) {
	designRepository = repository.NewDesignRepository(db)
	productRepository = repository.NewProductRepository(db)
	userRepository = repository.NewUserRepository(db)
}

func getUserByID(userID uint) (*model.User, error) {
	return userRepository.FindByID(userID)
}

// GenerateDesign godoc
// @Summary 生成AI设计
// @Description 使用AI根据提示词生成创意设计图案
// @Tags designs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.AIGenerateRequest true "AI生成请求参数"
// @Success 200 {object} service.AIGenerateResponse "设计生成成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /designs/generate [post]
func GenerateDesign(c *gin.Context) {
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

	var req service.AIGenerateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	if req.Prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Prompt is required",
		})
		return
	}

	logger.Info("Generating design for user", "userID", userID, "prompt", req.Prompt)

	base64Image, err := aiService.GenerateImage(req.Prompt)
	if err != nil {
		logger.Error("Failed to generate image", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to generate image",
			"details": err.Error(),
		})
		return
	}

	imageID := uuid.New().String()
	filename := fmt.Sprintf("%s.png", imageID)
	filepath := filepath.Join("./static/images", filename)

	err = aiService.DecodeAndSaveImage(base64Image, filepath)
	if err != nil {
		logger.Error("Failed to save generated image", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to save generated image",
			"details": err.Error(),
		})
		return
	}

	imageURL := fmt.Sprintf("/images/%s", filename)

	// 保存设计信息到数据库
	design := &model.Design{
		UserID:   userID,
		Prompt:   req.Prompt,
		ImageURL: imageURL,
		Style:    req.Style,
		Category: req.Category,
	}

	if err := designRepository.Create(design); err != nil {
		logger.Error("Failed to save design to database", err)
		// 不返回错误，因为图片生成已经成功
	}

	logger.Info("Design generated successfully", "userID", userID, "imageURL", imageURL, "designID", design.ID)

	c.JSON(http.StatusOK, service.AIGenerateResponse{
		ImageURL: imageURL,
		Message:  "Design generated successfully",
	})
}

// GetUserDesigns godoc
// @Summary 获取用户设计
// @Description 获取当前用户的所有设计作品，支持按分类筛选
// @Tags designs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param category query string false "分类筛选 (poster, sticker, canvas, tshirt)"
// @Success 200 {object} map[string]interface{} "用户设计列表"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /designs/my-designs [get]
func GetUserDesigns(c *gin.Context) {
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

	// 获取分类筛选参数
	category := c.Query("category")
	
	logger.Info("Fetching designs for user", "userID", userID, "category", category)

	// 从数据库获取用户的设计作品
	var designs []model.Design
	var err error
	
	if category != "" {
		designs, err = designRepository.FindByUserIDAndCategory(userID, category)
	} else {
		designs, err = designRepository.FindByUserID(userID)
	}
	if err != nil {
		logger.Error("Failed to fetch designs from database", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to fetch designs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"designs": designs,
		"message": "Successfully fetched user designs",
	})
}

// PublishDesignToShop godoc
// @Summary 发布设计到商店
// @Description 将设计作品发布到创意商店供其他用户购买
// @Tags designs
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body service.PublishToShopRequest true "发布到商店请求参数"
// @Success 201 {object} map[string]interface{} "发布成功"
// @Failure 400 {object} map[string]interface{} "请求参数错误"
// @Failure 401 {object} map[string]interface{} "未授权"
// @Failure 404 {object} map[string]interface{} "设计不存在"
// @Failure 500 {object} map[string]interface{} "内部服务器错误"
// @Router /designs/publish [post]
func PublishDesignToShop(c *gin.Context) {
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

	var req service.PublishToShopRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request format",
			"details": err.Error(),
		})
		return
	}

	// 验证设计是否存在且属于当前用户
	design, err := designRepository.FindByID(req.DesignID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Design not found",
		})
		return
	}

	if design.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied",
		})
		return
	}

	// 获取用户信息
	user, err := getUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user information",
		})
		return
	}

	// 创建商品
	product := &model.Product{
		Name:         req.ProductName,
		Description:  req.Description,
		BasePrice:    req.Price,
		Category:     design.Category,
		Material:     req.Material,
		Brand:        "创意工坊",
		IsActive:     true,
		DesignID:     &design.ID,
		CreatorID:    &userID,
		CreatorName:  user.Username,
		DesignPrompt: design.Prompt,
		DesignStyle:  design.Style,
		ImageURL:     design.ImageURL,
	}

	if err := productRepository.Create(product); err != nil {
		logger.Error("Failed to create product", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to publish to shop",
			"details": err.Error(),
		})
		return
	}

	logger.Info("Design published to shop successfully", "userID", userID, "designID", req.DesignID, "productID", product.ID)

	c.JSON(http.StatusCreated, gin.H{
		"product": product,
		"message": "Design published to shop successfully",
	})
}
