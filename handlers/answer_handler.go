package handlers

import (
	"encoding/json"
	"fmt"
	"golang-faq-app/services"
	"golang-faq-app/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AnswerHandler struct {
	faqService *services.FAQService
}

// NewAnswerHandler creates a new instance of AnswerHandler
func NewAnswerHandler(faqService *services.FAQService) *AnswerHandler {
	return &AnswerHandler{faqService: faqService}
}

// Handle handles the /answer endpoint
func (h *AnswerHandler) Handle(c *gin.Context) {
	var requestBody struct {
		Q string `json:"q"`
	}

	// Include a warning in the response if the request is not gzip-compressed
	if c.Request.Header.Get("Content-Encoding") != "gzip" {
		c.Writer.Header().Set("X-Warning", "Request was not gzip-compressed")
	}

	// Decompress the request body
	body, err := utils.DecompressRequest(c.Request)
	if err != nil {
		c.String(http.StatusBadRequest, fmt.Sprintf("Failed to process request: %v", err))
		return
	}

	// Parse the JSON data
	if err := json.Unmarshal(body, &requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Use the FAQService to find the best answer
	result, err := h.faqService.FindBestAnswer(requestBody.Q)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, result)
	}
}
