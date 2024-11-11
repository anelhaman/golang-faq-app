package interfaces

type QuestionAnswerSource interface {
	LoadQuestions() error
	FindAnswer(query string) (string, error)
	GetQuestions() map[string]string // This method should be defined here
}
