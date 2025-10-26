package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"backend-ai/internal/config"
)

type ZhipuService struct {
	apiKey string
	apiURL string
	client *http.Client
}

// Struktur request untuk API v4 (kompatibel OpenAI)
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	TopP        float64   `json:"top_p,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Struktur response dari API v4 (kompatibel OpenAI)
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

func NewZhipuService(cfg *config.Config) *ZhipuService {
	return &ZhipuService{
		apiKey: cfg.ZhipuAPIKey,
		apiURL: cfg.ZhipuAPIURL,
		client: &http.Client{
			Timeout: 90 * time.Second,
		},
	}
}

func (s *ZhipuService) Chat(messages []Message) (*ChatResponse, error) {
	// Gunakan model GLM-4.5-Flash untuk API v4
	requestBody := ChatRequest{
		Model:       "GLM-4.5-Flash",
		Messages:    messages,
		Temperature: 0.7,
		TopP:        0.9,
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	req, err := http.NewRequest("POST", s.apiURL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.apiKey))

	// Log untuk debugging
	log.Printf("Sending request to URL: %s", req.URL.String())
	log.Printf("Request Body: %s", string(jsonBody))

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var chatResponse ChatResponse
	if err := json.Unmarshal(body, &chatResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Cek jika API mengembalikan error (format error API v4 bisa berbeda)
	if resp.StatusCode != http.StatusOK {
		// Coba parse error response
		var errorResp map[string]interface{}
		if json.Unmarshal(body, &errorResp) == nil {
			if errMsg, ok := errorResp["error"].(map[string]interface{}); ok {
				if msg, ok := errMsg["message"].(string); ok {
					return nil, fmt.Errorf("API error: %s", msg)
				}
			}
		}
		return nil, fmt.Errorf("API error: %s", string(body))
	}

	if len(chatResponse.Choices) == 0 {
		return nil, fmt.Errorf("no response from ai")
	}

	return &chatResponse, nil
}
