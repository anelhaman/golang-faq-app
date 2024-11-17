package initializer

import (
	"golang-faq-app/config"
	"golang-faq-app/handlers"
	"golang-faq-app/interfaces"
	"golang-faq-app/services"
)

// InitializeApp loads the configuration, initializes the service, and prepares the application
func InitializeApp() (*config.Config, *services.FAQService, error) {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, err
	}

	// Ensure at least one answer is allowed
	amountAnswer := cfg.MaxAnswers
	if amountAnswer <= 1 {
		amountAnswer = 1
	}

	// Initialize FAQService
	faqService := services.NewFAQService(amountAnswer)

	// Add sources based on the configuration
	for _, file := range cfg.Files {
		var handler interfaces.QuestionAnswerSource

		// Check if file has a URL or a path to determine if it's remote or local
		if file.URL != "" { // Remote file
			switch file.Type {
			case "csv":
				handler = handlers.NewRemoteCSVHandler(file.URL)
			case "excel":
				handler = handlers.NewRemoteExcelHandler(file.URL)
			}
		} else { // Local file
			switch file.Type {
			case "csv":
				handler = handlers.NewCSVHandler(file.Path)
			case "excel":
				handler = handlers.NewExcelHandler(file.Path)
			}
		}

		// Add handler to the service if it's not nil
		if handler != nil {
			faqService.AddSource(handler)
		}
	}

	// Load all questions from sources
	if err := faqService.LoadAllSources(); err != nil {
		return nil, nil, err
	}

	return &cfg, faqService, nil
}
