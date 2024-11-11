# Golang QA App

`golang-qa-app` is a question-answering API built in Go. It reads questions and answers from multiple CSV or Excel files and responds with the answer that most closely matches a given query. The application includes confidence scoring and timestamps in responses.

## Features

- Reads from multiple CSV or Excel files, specified in a YAML configuration file.
- Compares the input query with stored questions and returns the best-matched answer along with a confidence score.
- Returns responses as a JSON API via Gin Gonic.
- Supports Go Modules and is designed using an interface-based architecture with separate files for scalability and readability.

## Project Structure

```plaintext
golang-qa-app/
├── config/
│   └── config.yaml         # Configuration file with paths to data files
├── data/
│   ├── questions1.csv      # Example CSV data file
│   └── questions2.xlsx     # Example Excel data file
├── handlers/
│   ├── csv_handler.go      # CSV file handler
│   └── excel_handler.go    # Excel file handler
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
- Go 1.18+
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
GET /answer

Description: Finds the best answer for a given question.

Request:

query (string, required): The question to find the best answer for.


Example:

```
GET http://localhost:8080/answer?query=What is Golang?
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

## License
This project is licensed under the MIT License.
