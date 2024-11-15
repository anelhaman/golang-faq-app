package services

import (
	"errors"
	"golang-qa-app/interfaces"
	"math"
	"sort"
	"strings"
	"sync"
	"time"
	"unicode"

	gothaiwordcut "github.com/narongdejsrn/go-thaiwordcut"
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

	// Extract questions and prepare a reverse map to link questions to answers
	questions := []string{}
	questionAnswerMap := make(map[string]string)

	for question, answer := range source.GetQuestions() {
		questions = append(questions, question)
		questionAnswerMap[question] = answer
	}

	// Use the parallel similarity function
	similarities := calculateSimilarityParallel(query, questions)

	// Filter and build results based on the confidence threshold
	for question, similarity := range similarities {
		if similarity > 0.55 { // Confidence threshold
			// Round confidence to two decimal places
			roundedConfidence := math.Round(similarity*100) / 100

			results = append(results, AnswerResult{
				Answer:     questionAnswerMap[question],
				Confidence: roundedConfidence,
				Timestamp:  time.Now(),
			})
		}
	}
	return results
}

func calculateSimilarityParallel(query string, questions []string) map[string]float64 {

	var queryWords, questionWords []string

	results := make(map[string]float64)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, question := range questions { // Use questions here
		wg.Add(1)
		go func(question string) {
			defer wg.Done()

			if containsThai(query) {
				queryWords = cutThaiWord(strings.ToLower(query))       // Tokenize query
				questionWords = cutThaiWord(strings.ToLower(question)) // Tokenize question
			} else {
				queryWords = strings.Fields(strings.ToLower(query))
				questionWords = strings.Fields(strings.ToLower(question))
			}

			matchCount := 0
			wordMap := make(map[string]bool)
			for _, word := range queryWords {
				wordMap[word] = true
			}

			for _, word := range questionWords {
				if wordMap[word] {
					matchCount++
				}
			}

			similarity := float64(matchCount*2) / float64(len(queryWords)+len(questionWords))

			mu.Lock() // Protect map access
			results[question] = similarity
			mu.Unlock()
		}(question)
	}

	wg.Wait()
	return results
}

func cutThaiWord(s string) []string {

	segmenter := gothaiwordcut.Wordcut()
	segmenter.LoadDefaultDict()
	result := segmenter.Segment(s)

	return result
}

// ContainsThai checks if a string contains any Thai characters
func containsThai(input string) bool {
	for _, r := range input {
		if unicode.In(r, unicode.Thai) {
			return true
		}
	}
	return false
}
