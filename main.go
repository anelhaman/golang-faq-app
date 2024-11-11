package main

import (
	"golang-qa-app/handlers"
	"golang-qa-app/interfaces"
	"golang-qa-app/services"
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
	Path string `yaml:"path"`
	Type string `yaml:"type"` // "csv" or "excel"
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

	qaService := services.NewQAService(AmountAnswer)

	// Initialize handlers based on config
	for _, file := range config.Files {
		var handler interface{}
		switch file.Type {
		case "csv":
			handler = handlers.NewCSVHandler(file.Path)
		case "excel":
			handler = handlers.NewExcelHandler(file.Path)
		}

		if handler != nil {
			qaService.AddSource(handler.(interfaces.QuestionAnswerSource))
		}
	}

	// Load questions from all sources
	if err := qaService.LoadAllSources(); err != nil {
		log.Fatalf("Error loading questions: %v", err)
	}

	// Set up Gin Gonic API
	router := gin.Default()

	router.GET("/answer", func(c *gin.Context) {
		question := c.Query("question")
		result, err := qaService.FindBestAnswer(question)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusOK, result)
		}
	})

	router.Run(":8080")
}
