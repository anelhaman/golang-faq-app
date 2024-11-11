package handlers

import (
	"bytes"
	"encoding/csv"
	"log"
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
		// Log the error but skip this handler and return nil to continue processing other files
		log.Printf("Failed to fetch file from URL: %s, error: %v\n", r.url, err)
		return nil
	}

	// Check for successful response
	if resp.StatusCode() != 200 {
		// Log the error but skip this handler and return nil to continue processing other files
		log.Printf("Failed to fetch file from URL: %s, status code: %d\n", r.url, resp.StatusCode())
		return nil
	}

	// Read the response body into a buffer
	buf := bytes.NewBuffer(resp.Body())

	// Parse the CSV file from the buffer
	reader := csv.NewReader(buf)
	records, err := reader.ReadAll()
	if err != nil {
		// Log the error but skip this handler and return nil to continue processing other files
		log.Printf("Failed to parse CSV file from URL: %s, error: %v\n", r.url, err)
		return nil
	}

	// Store the questions and answers (assuming 1st column is question, 2nd is answer)
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
