# Golang QA App

`golang-faq-app` is a question-answering API built in Go. It reads questions and answers from multiple CSV or Excel files and responds with the answer that most closely matches a given query. The application includes confidence scoring and timestamps in responses.

## Features

- Reads from multiple CSV or Excel files, specified in a YAML configuration file.
- Compares the input query with stored questions and returns the best-matched answer along with a confidence score.
- Returns responses as a JSON API via Gin Gonic.
- Supports Go Modules and is designed using an interface-based architecture with separate files for scalability and readability.

## Project Structure

```plaintext
golang-faq-app/
├── config/
│   └── config.yaml         # Configuration file with paths to data files
├── data/
│   ├── questions1.csv      # Example CSV data file
│   └── questions2.xlsx     # Example Excel data file
├── handlers/
│   ├── csv_handler.go      # CSV file handler
│   ├── excel_handler.go    # Excel file handler
│   ├── remote_csv_handler.go      # Remote CSV file handler
│   └── remote_excel_handler.go    # Remote Excel file handler
├── interfaces/
│   └── qa_source.go        # Interface definition for question-answer sources
├── services/
│   └── qa_service.go       # QA service with logic to find best-matched answer
├── main.go                 # Entry point for the API server
├── go.mod                  # Go Modules file
└── README.md               # Project README
```

# Getting Started
Prerequisites
- Go 1.23
- Gin Gonic
- Excelize (for handling Excel files)
- YAML v3
Install dependencies:

```
go mod tidy
```
## Configuration
Create a config/config.yaml file to specify the paths of your CSV or Excel files:

```
files:
  - path: "data/questions1.csv"
    type: "csv"
  - path: "data/questions2.xlsx"
    type: "excel"
  - url: "https://example.com/questions.csv"
    type: "csv"

max_answers: 2  # Configurable number of answers to return
```

## Running the Application
To start the application:

```
go run main.go
```
The API will run on http://localhost:8080.

## API Endpoints
POST /answer

Description: Finds the best answer for a given question.

Request:

query (string, required): The question to find the best answer for.


Example:

```
POST http://localhost:8080/answer

{
	"q": "What is Golang?"
}
```

Response:


```
{
  "answer": "Golang is an open-source programming language.",
  "confidence": 0.92,
  "timestamp": "2024-11-11T15:04:05Z"
}
```

## Code Overview
- interfaces/qa_source.go: Defines the QuestionAnswerSource interface, which includes methods to load and retrieve questions and answers.
- handlers/: Contains CSVHandler and ExcelHandler, which implement QuestionAnswerSource.
- services/qa_service.go: The core service that processes the query, calculates similarity, and provides the answer with confidence.
- main.go: Initializes the application, loads data sources, and starts the Gin server.

## Docker Setup

To run the application using Docker, follow these steps:

### Build and Start with Docker Compose

1. Ensure that you have Docker and Docker Compose installed on your machine.

2. In the root of your project directory, you will find the `docker-compose.yml` file. This file defines the necessary services to run the application, including NGINX as a reverse proxy and the Go application.

3. If you have not already, create a `config.yaml` file inside the `config/` directory with the necessary configuration (refer to the `config/config.yaml` file structure in the project).

4. Run the following command to build and start the services:

   ```
   docker-compose up --build
   ```
Once the containers are up and running, the Go application will be available at http://localhost (NGINX will forward requests to the Go application running on port 8080).

## License
This project is licensed under the MIT License.
