package handlers

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
)

type CSVHandler struct {
	path      string
	questions map[string]string
}

func NewCSVHandler(path string) *CSVHandler {
	return &CSVHandler{path: path, questions: make(map[string]string)}
}

func (c *CSVHandler) LoadQuestions() error {
	// Check if the file exists before opening
	if _, err := os.Stat(c.path); os.IsNotExist(err) {
		// Log the error and skip this file
		log.Printf("File not found: %s, skipping...\n", c.path)
		return nil
	}

	// Open the CSV file
	file, err := os.Open(c.path)
	if err != nil {
		// Log the error and skip this file
		log.Printf("Failed to open CSV file: %s, error: %v\n", c.path, err)
		return nil
	}
	defer file.Close()

	// Parse the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		// Log the error and skip this file
		log.Printf("Failed to read CSV file: %s, error: %v\n", c.path, err)
		return nil
	}

	// Store the questions and answers (assuming 1st column is question, 2nd is answer)
	for _, record := range records {
		if len(record) >= 2 {
			c.questions[record[0]] = record[1]
		}
	}

	return nil
}

func (c *CSVHandler) FindAnswer(query string) (string, error) {
	for question, answer := range c.questions {
		if strings.Contains(strings.ToLower(question), strings.ToLower(query)) {
			return answer, nil
		}
	}
	return "", nil
}

// Implement GetQuestions to match the interface
func (c *CSVHandler) GetQuestions() map[string]string {
	return c.questions
}
