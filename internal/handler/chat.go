package handler

import (
	"net/http"
	"strings"

	"backend-ai/internal/service"

	"github.com/labstack/echo/v4"
)

type ChatHandler struct {
	zhipuService *service.ZhipuService
}

type ChatRequest struct {
	Message string `json:"message" validate:"required"`
}

type ChatResponse struct {
	Reply string `json:"reply"`
}

func NewChatHandler(zhipuService *service.ZhipuService) *ChatHandler {
	return &ChatHandler{
		zhipuService: zhipuService,
	}
}

func (h *ChatHandler) Chat(c echo.Context) error {
	var req ChatRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Prepare messages for Zhipu AI
	messages := []service.Message{
		{
			Role:    "user",
			Content: req.Message,
		},
	}

	// Call Zhipu AI API
	resp, err := h.zhipuService.Chat(messages)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get response from AI: " + err.Error()})
	}

	// Karena response dari service v4 tidak memiliki field `Data`, kita akses langsung `resp.Choices`
	if len(resp.Choices) == 0 {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "No response from AI"})
	}

	// Return the AI response
	return c.JSON(http.StatusOK, ChatResponse{
		Reply: strings.TrimSpace(resp.Choices[0].Message.Content),
	})
}
