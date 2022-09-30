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

// OBSFile
type OBSFile interface {
	OBSFile() string
}

func NewOBSFile(v string) (OBSFile, error) {
	// TODO check format
	return obsFile(v), nil
}

type obsFile string

func (s obsFile) OBSFile() string {
	return string(s)
}
