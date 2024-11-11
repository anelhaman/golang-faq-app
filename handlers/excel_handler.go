package handlers

import (
	"log"
	"os"
	"strings"

	"github.com/xuri/excelize/v2"
)

type ExcelHandler struct {
	path      string
	questions map[string]string
}

func NewExcelHandler(path string) *ExcelHandler {
	return &ExcelHandler{path: path, questions: make(map[string]string)}
}

func (e *ExcelHandler) LoadQuestions() error {
	// Check if the file exists before opening
	if _, err := os.Stat(e.path); os.IsNotExist(err) {
		// Log the error and skip this file
		log.Printf("File not found: %s, skipping...\n", e.path)
		return nil
	}

	// Open the Excel file
	f, err := excelize.OpenFile(e.path)
	if err != nil {
		// Log the error and skip this file
		log.Printf("Failed to open Excel file: %s, error: %v\n", e.path, err)
		return nil
	}

	// Assume the questions are in the first sheet
	rows, err := f.GetRows(f.GetSheetName(0))
	if err != nil {
		// Log the error and skip this file
		log.Printf("Failed to get rows from Excel file: %s, error: %v\n", e.path, err)
		return nil
	}

	// Store the questions and answers (assuming 1st column is question, 2nd is answer)
	for _, row := range rows {
		if len(row) >= 2 {
			e.questions[row[0]] = row[1]
		}
	}

	return nil
}

func (e *ExcelHandler) FindAnswer(query string) (string, error) {
	for question, answer := range e.questions {
		if strings.Contains(strings.ToLower(question), strings.ToLower(query)) {
			return answer, nil
		}
	}
	return "", nil
}

// Implement GetQuestions to match the interface
func (e *ExcelHandler) GetQuestions() map[string]string {
	return e.questions
}
