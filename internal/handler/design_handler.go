package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Incipe-win/ai-tshirt-shop/internal/service"
	"github.com/Incipe-win/ai-tshirt-shop/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var aiService *service.AIService

func init() {
	aiService = service.NewAIService()
	
	if err := os.MkdirAll("./static/images", 0755); err != nil {
		logger.Error("Failed to create static images directory", err)
	}
}

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
			"error": "Invalid request format",
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
			"error": "Failed to generate image",
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
			"error": "Failed to save generated image",
			"details": err.Error(),
		})
		return
	}

	imageURL := fmt.Sprintf("/images/%s", filename)
	
	logger.Info("Design generated successfully", "userID", userID, "imageURL", imageURL)

	c.JSON(http.StatusOK, service.AIGenerateResponse{
		ImageURL: imageURL,
		Message:  "Design generated successfully",
	})
}