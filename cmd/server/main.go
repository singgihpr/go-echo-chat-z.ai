package main

import (
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"backend-ai/internal/config"
	"backend-ai/internal/handler"
	"backend-ai/internal/service"
)

func main() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Initialize Echo
	e := echo.New()

	// Register validator
	e.Validator = &CustomValidator{validator: validator.New()}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Initialize services
	zhipuService := service.NewZhipuService(cfg)

	// Initialize handlers
	chatHandler := handler.NewChatHandler(zhipuService)

	// Routes
	e.Static("/", "public")

	e.GET("/api/info", func(c echo.Context) error { // <-- DIUBAH PATH-NYA
		return c.JSON(http.StatusOK, map[string]string{
			"message": "Welcome to Echo Zhipu AI Chat API",
			"version": "1.0.0",
		})
	})

	e.POST("/api/chat", chatHandler.Chat)

	// Start server
	log.Printf("Starting server on port %s", cfg.Port)
	if err := e.Start(":" + cfg.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// CustomValidator implements echo.Validator interface
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
