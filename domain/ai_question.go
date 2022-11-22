package domain

type AIQuestion struct {
	Id string
}

type QuestionSubmission struct {
	Id       string
	Account  Account
	SubmitAt string
	Score    int
	Times    int // like version
}

type ChoiceQuestion struct {
	Id     int
	Desc   string
	Answer string
	Option []QuestionOption
}

type QuestionOption struct {
	Id   string
	Desc string
}

type CompletionQuestion struct {
	Id     int
	Desc   string
	Answer string
}
