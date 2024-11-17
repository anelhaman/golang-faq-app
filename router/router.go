package router

import (
	"golang-faq-app/config"
	"golang-faq-app/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter initializes and configures the Gin router
func SetupRouter(answerHandler *handlers.AnswerHandler, config *config.Config) *gin.Engine {

	router := gin.Default()

	// Add middleware for CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Encoding")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Define routes
	router.POST("/answer", answerHandler.Handle)

	return router
}
