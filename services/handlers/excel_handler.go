package handlers

import (
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
	file, err := excelize.OpenFile(e.path)
	if err != nil {
		return err
	}
	defer file.Close()

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return err
	}

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
