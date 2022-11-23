package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type AIQuestion interface {
	GetCompetitorAndScores(string, domain.Account) (bool, []int, error)

	SaveCompetitor(string, *domain.CompetitorInfo) error

	GetQuestions(pool string, choice, complition []int) (
		[]domain.ChoiceQuestion, []domain.CompletionQuestion, error,
	)

	GetSubmission(sid string, user domain.Account, date string) (
		domain.QuestionSubmission, error,
	)

	SaveSubmission(sid string, v *domain.QuestionSubmission) (string, error)
}
