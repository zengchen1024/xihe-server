package domain

// Question
type Question interface {
	Question() string
}

func NewQuestion(v string) (Question, error) {
	// TODO check format

	return question(v), nil
}

type question string

func (s question) Question() string {
	return string(s)
}
