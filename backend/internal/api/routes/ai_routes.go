package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/ryanlisse/go-crypto-bot/internal/api/handlers"
	"github.com/ryanlisse/go-crypto-bot/internal/domain/ai/service"
)

// RegisterAIRoutesWithGin registers AI-related routes with Gin router
func RegisterAIRoutesWithGin(
	router *gin.Engine,
	aiSvc service.AIService,
) {
	// Create AI handlers
	chatHandler := handlers.ChatHandler(aiSvc)
	functionHandler := handlers.FunctionHandler(aiSvc)

	// Register routes
	ai := router.Group("/api/v1/ai")
	{
		// Convert http.HandlerFunc to gin.HandlerFunc
		ai.POST("/chat", func(c *gin.Context) {
			chatHandler(c.Writer, c.Request)
		})
		ai.POST("/function", func(c *gin.Context) {
			functionHandler(c.Writer, c.Request)
		})
	}
}
