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

	"github.com/spf13/viper"
)

type AIService struct {
	client *http.Client
}

type StabilityAIRequest struct {
	TextPrompts []TextPrompt `json:"text_prompts"`
	CfgScale    float64      `json:"cfg_scale"`
	Height      int          `json:"height"`
	Width       int          `json:"width"`
	Samples     int          `json:"samples"`
	Steps       int          `json:"steps"`
}

type TextPrompt struct {
	Text   string  `json:"text"`
	Weight float64 `json:"weight"`
}

type StabilityAIResponse struct {
	Artifacts []Artifact `json:"artifacts"`
}

type Artifact struct {
	Base64       string `json:"base64"`
	Seed         int    `json:"seed"`
	FinishReason string `json:"finishReason"`
}

type AIGenerateRequest struct {
	Prompt string `json:"prompt" binding:"required"`
}

type AIGenerateResponse struct {
	ImageURL string `json:"image_url"`
	Message  string `json:"message"`
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

	requestBody := StabilityAIRequest{
		TextPrompts: []TextPrompt{
			{
				Text:   prompt,
				Weight: 1.0,
			},
		},
		CfgScale: 7.0,
		Height:   1024,
		Width:    1024,
		Samples:  1,
		Steps:    30,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api.stability.ai/v1/generation/stable-diffusion-xl-1024-v1-0/text-to-image", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var apiResponse StabilityAIResponse
	if err := json.Unmarshal(body, &apiResponse); err != nil {
		return "", fmt.Errorf("failed to unmarshal response: %v", err)
	}

	if len(apiResponse.Artifacts) == 0 {
		return "", fmt.Errorf("no images generated")
	}

	return apiResponse.Artifacts[0].Base64, nil
}

func (s *AIService) DecodeAndSaveImage(base64Data string, filename string) error {
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return fmt.Errorf("failed to decode base64 image: %v", err)
	}

	return saveImageFile(imageData, filename)
}

func saveImageFile(data []byte, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}