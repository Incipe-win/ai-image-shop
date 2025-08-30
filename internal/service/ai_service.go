package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Incipe-win/ai-tshirt-shop/pkg/logger"
	"github.com/spf13/viper"
)

type AIService struct {
	client *http.Client
}

type SeedreamAIRequest struct {
	Model           string `json:"model"`
	Prompt          string `json:"prompt"`
	Response_format string `json:"response_format"`
	Size            string `json:"size"`
	Guidance_scale  int    `json:"guidance_scale"`
	Watermark       bool   `json:"watermark"`
}

type SeedreamAIResponse struct {
	Model   string `json:"model"`
	Created int    `json:"created"`
	Data    []Data `json:"data"`
	Usage   Usage  `json:"usage"`
}

type Data struct {
	B64_json string `json:"b64_json"`
}

type Usage struct {
	Generated_images int `json:"generated_images"`
	Output_tokens    int `json:"output_tokens"`
	Total_tokens     int `json:"total_tokens"`
}

type AIGenerateRequest struct {
	Prompt   string `json:"prompt" binding:"required"`
	Category string `json:"category"`
	Style    string `json:"style"`
}

type AIGenerateResponse struct {
	ImageURL string `json:"image_url"`
	Message  string `json:"message"`
}

type PublishToShopRequest struct {
	DesignID    uint    `json:"design_id" binding:"required"`
	ProductName string  `json:"product_name" binding:"required"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,min=0.01"`
	Material    string  `json:"material"`
}

func NewAIService() *AIService {
	return &AIService{
		client: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (s *AIService) GenerateImage(prompt string) (string, error) {
	apiKey := viper.GetString("ai.api_key")
	if apiKey == "" {
		return "", fmt.Errorf("AI API key not configured")
	}

	requestBody := SeedreamAIRequest{
		Model:           "doubao-seedream-3-0-t2i-250415",
		Prompt:          prompt,
		Response_format: "b64_json",
		Size:            "1024x1024",
		Guidance_scale:  3,
		Watermark:       false,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	logger.Debug("Sending AI request", "body", string(jsonData))

	req, err := http.NewRequest("POST", "https://ark.cn-beijing.volces.com/api/v3/images/generations", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		logger.Debug("AI request error", "error", err)
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Debug("AI response read error", "error", err)
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		logger.Debug("AI response error", "status", resp.StatusCode)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse SeedreamAIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		logger.Debug("AI response unmarshal error", "error", err)
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(apiResponse.Data) == 0 {
		return "", fmt.Errorf("no images generated")
	}

	return apiResponse.Data[0].B64_json, nil
}

func (s *AIService) DecodeAndSaveImage(base64Data string, filename string) error {
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		logger.Error("Failed to decode base64 image", err)
		return fmt.Errorf("failed to decode base64 image: %v", err)
	}

	return saveImageFile(imageData, filename)
}

func saveImageFile(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		logger.Error("Failed to create image file", err)
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		logger.Error("Failed to write image file", err)
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
