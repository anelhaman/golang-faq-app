package handlers

import (
	"encoding/csv"
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
	file, err := os.Open(c.path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

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
