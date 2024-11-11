package handlers

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type RemoteCSVHandler struct {
	url       string
	questions map[string]string
}

func NewRemoteCSVHandler(url string) *RemoteCSVHandler {
	return &RemoteCSVHandler{url: url, questions: make(map[string]string)}
}

func (r *RemoteCSVHandler) LoadQuestions() error {
	// Initialize Resty client
	client := resty.New()

	// Make the GET request using Resty
	resp, err := client.R().Get(r.url)
	if err != nil {
		return err
	}

	// Check for successful response
	if resp.StatusCode() != 200 {
		return fmt.Errorf("failed to fetch the file, status code: %d", resp.StatusCode())
	}

	// Read the response body into a buffer
	buf := bytes.NewBuffer(resp.Body())

	// Process the CSV content
	reader := csv.NewReader(buf)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	// Store the questions and answers
	for _, record := range records {
		if len(record) >= 2 {
			r.questions[record[0]] = record[1]
		}
	}

	return nil
}

func (r *RemoteCSVHandler) FindAnswer(query string) (string, error) {
	for question, answer := range r.questions {
		if strings.Contains(strings.ToLower(question), strings.ToLower(query)) {
			return answer, nil
		}
	}
	return "", nil
}

// Implement GetQuestions to match the interface
func (r *RemoteCSVHandler) GetQuestions() map[string]string {
	return r.questions
}
