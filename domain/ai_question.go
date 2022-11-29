package domain

const (
	AIQuestionStatusStart = "start"
	AIQuestionStatusEnd   = "end"
)

type QuestionSubmissionInfo struct {
	Account Account
	Score   int
}

type QuestionSubmission struct {
	Id      string
	Account Account
	Date    string
	Status  string
	Expiry  int64
	Score   int
	Times   int
	Version int
}

// == account && == date && == times && status == start && now < expiry && > score
// status = end , expiry = 0, score = , times++

type ChoiceQuestion struct {
	Desc    string
	Answer  string
	Options []string
}

type CompletionQuestion struct {
	Desc   string
	Info   string
	Answer string
}
