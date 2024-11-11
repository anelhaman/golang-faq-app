package services

import (
	"errors"
	"golang-qa-app/interfaces"
	"math"
	"sort"
	"strings"
	"time"
)

type AnswerResult struct {
	Answer     string    `json:"answer"`
	Confidence float64   `json:"confidence"`
	Timestamp  time.Time `json:"timestamp"`
}

type QAService struct {
	sources    []interfaces.QuestionAnswerSource
	maxAnswers int // Store the configurable max answers
}

func NewQAService(maxAnswers int) *QAService {
	return &QAService{
		sources:    []interfaces.QuestionAnswerSource{},
		maxAnswers: maxAnswers,
	}
}

func (s *QAService) AddSource(source interfaces.QuestionAnswerSource) {
	s.sources = append(s.sources, source)
}

func (s *QAService) LoadAllSources() error {
	for _, source := range s.sources {
		if err := source.LoadQuestions(); err != nil {
			return err
		}
	}
	return nil
}

func (s *QAService) FindBestAnswer(query string) ([]AnswerResult, error) {
	var allAnswers []AnswerResult

	// Find answers with confidence scores
	for _, source := range s.sources {
		allAnswers = append(allAnswers, s.searchWithConfidence(source, query)...)
	}

	// Sort by confidence score
	sort.Slice(allAnswers, func(i, j int) bool {
		return allAnswers[i].Confidence > allAnswers[j].Confidence
	})

	// Take the top N answers based on configured maxAnswers
	if len(allAnswers) > s.maxAnswers {
		allAnswers = allAnswers[:s.maxAnswers]
	}

	if len(allAnswers) == 0 {
		return nil, errors.New("no matching answers found")
	}

	return allAnswers, nil
}

func (s *QAService) searchWithConfidence(source interfaces.QuestionAnswerSource, query string) []AnswerResult {
	var results []AnswerResult

	for question, answer := range source.GetQuestions() {
		similarity := calculateSimilarity(query, question)
		if similarity > 0.55 { // Confidence threshold, adjust as needed

			// Round confidence to two decimal places
			roundedConfidence := math.Round(similarity*100) / 100

			results = append(results, AnswerResult{
				Answer:     answer,
				Confidence: roundedConfidence,
				Timestamp:  time.Now(),
			})
		}
	}
	return results
}

func calculateSimilarity(query, question string) float64 {
	qWords := strings.Fields(strings.ToLower(query))
	aWords := strings.Fields(strings.ToLower(question))
	matchCount := 0

	for _, qWord := range qWords {
		for _, aWord := range aWords {
			if qWord == aWord {
				matchCount++
			}
		}
	}

	return float64(matchCount*2) / float64(len(qWords)+len(aWords)) // Ratio of matches
}
