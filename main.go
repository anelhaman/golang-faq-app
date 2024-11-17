package main

import (
	"golang-faq-app/handlers"

	"golang-faq-app/initializer"
	"golang-faq-app/router"
	"log"
)

func main() {
	// Load configuration and initialize dependencies
	config, faqService, err := initializer.InitializeApp()
	if err != nil {
		log.Fatalf("Error during initialization: %v", err)
	}

	// Initialize the handler
	answerHandler := handlers.NewAnswerHandler(faqService)

	// Set up the router
	appRouter := router.SetupRouter(answerHandler, config)

	// Start the server
	log.Println("Starting server on :8080")
	appRouter.Run(":8080")
}
