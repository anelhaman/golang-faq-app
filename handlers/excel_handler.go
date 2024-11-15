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

	// Get all sheet names
	sheetNames := f.GetSheetList()

	// Iterate through all sheets
	for _, sheetName := range sheetNames {
		// Get rows from the current sheet
		rows, err := f.GetRows(sheetName)
		if err != nil {
			// Log the error and skip this sheet
			log.Printf("Failed to get rows from sheet: %s in Excel file: %s, error: %v\n", sheetName, e.path, err)
			continue
		}

		// Store the questions and answers (assuming 1st column is question, 2nd is answer)
		for _, row := range rows {
			if len(row) >= 2 {
				e.questions[row[0]] = row[1]
			}
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
