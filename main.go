package main

import (
	"golang-faq-app/handlers"
	"golang-faq-app/interfaces"
	"golang-faq-app/services"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v3"
)

// Define the Config struct to hold file paths
type Config struct {
	Files      []FileConfig `yaml:"files"`
	MaxAnswers int          `yaml:"max_answers"`
}

type FileConfig struct {
	Path string `yaml:"path,omitempty"` // Path for local files
	URL  string `yaml:"url,omitempty"`  // URL for remote files
	Type string `yaml:"type"`           // "csv" or "excel"
}

// Function to load configuration from config.yaml
func loadConfig() (Config, error) {
	var config Config

	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

func main() {

	config, err := loadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	AmountAnswer := config.MaxAnswers
	if AmountAnswer <= 1 {
		AmountAnswer = 1
	}

	faqService := services.NewFAQService(AmountAnswer)

	// Initialize handlers based on config
	for _, file := range config.Files {
		var handler interfaces.QuestionAnswerSource

		// Check if file has a URL or a path to determine if it's remote or local
		if file.URL != "" { // Check if there's a URL field
			switch file.Type {
			case "csv":
				handler = handlers.NewRemoteCSVHandler(file.URL)
			case "excel":
				// If you have a handler for remote Excel files, add it here
				handler = handlers.NewRemoteExcelHandler(file.URL)
			}
		} else { // If no URL, it's a local file
			switch file.Type {
			case "csv":
				handler = handlers.NewCSVHandler(file.Path)
			case "excel":
				// Assuming you have a local handler for Excel
				handler = handlers.NewExcelHandler(file.Path)
			}
		}

		// Add the handler to the service if it's not nil
		if handler != nil {
			faqService.AddSource(handler)
		}
	}

	// Load questions from all sources
	if err := faqService.LoadAllSources(); err != nil {
		log.Fatalf("Error loading questions: %v", err)
	}

	// Set up Gin Gonic API
	router := gin.Default()

	router.POST("/answer", func(c *gin.Context) {
		var requestBody struct {
			Q string `json:"q"`
		}

		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
			return
		}

		result, err := faqService.FindBestAnswer(requestBody.Q)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, result)
		}
	})

	router.Run(":8080")
}
