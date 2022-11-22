package domain

type AIQuestion struct {
	Id string
}

const (
	aiquestionStatusStart = "start"
	aiquestionStatusEnd   = "end"
)

type QuestionSubmission struct {
	Id      string
	Account Account
	Date    string
	Status  string
	Expiry  int64 // add 10 minutes
	Score   int
	Times   int // like version
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
	Answer string
}
