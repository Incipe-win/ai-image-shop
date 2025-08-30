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
)

func init() {
	aiService = service.NewAIService()

	if err := os.MkdirAll("./static/images", 0755); err != nil {
		logger.Error("Failed to create static images directory", err)
	}
}

func InitDesignRepository(db *gorm.DB) {
	designRepository = repository.NewDesignRepository(db)
}

// GenerateDesign godoc
// @Summary 生成AI设计
// @Description 使用AI根据提示词生成T恤设计图案
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
		Style:    "", // 可以根据需要从请求中获取风格信息
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
// @Description 获取当前用户的所有设计作品
// @Tags designs
// @Accept json
// @Produce json
// @Security BearerAuth
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

	logger.Info("Fetching designs for user", "userID", userID)

	// 从数据库获取用户的设计作品
	designs, err := designRepository.FindByUserID(userID)
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
