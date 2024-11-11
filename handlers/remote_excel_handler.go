package handlers

import (
	"bytes"
	"log"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/xuri/excelize/v2"
)

type RemoteExcelHandler struct {
	url       string
	questions map[string]string
}

func NewRemoteExcelHandler(url string) *RemoteExcelHandler {
	return &RemoteExcelHandler{url: url, questions: make(map[string]string)}
}

func (r *RemoteExcelHandler) LoadQuestions() error {
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

	// Open the Excel file from the buffer
	f, err := excelize.OpenReader(buf)
	if err != nil {
		// Log the error but skip this handler and return nil to continue processing other files
		log.Printf("Failed to open Excel file from URL: %s, error: %v\n", r.url, err)
		return nil
	}

	// Assume the questions are in the first sheet
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		// Log the error but skip this handler and return nil to continue processing other files
		log.Printf("Failed to get rows from Excel file URL: %s, error: %v\n", r.url, err)
		return nil
	}

	// Store the questions and answers (assuming 1st column is question, 2nd is answer)
	for _, row := range rows {
		if len(row) >= 2 {
			r.questions[row[0]] = row[1]
		}
	}

	return nil
}

func (r *RemoteExcelHandler) FindAnswer(query string) (string, error) {
	for question, answer := range r.questions {
		if strings.Contains(strings.ToLower(question), strings.ToLower(query)) {
			return answer, nil
		}
	}
	return "", nil
}

// Implement GetQuestions to match the interface
func (r *RemoteExcelHandler) GetQuestions() map[string]string {
	return r.questions
}
