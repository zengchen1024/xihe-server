package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ChoiceQuestionDO = domain.ChoiceQuestion
type CompletionQuestionDO = domain.CompletionQuestion

type QuestionSubmissionDO struct {
	Id      string
	Account string
	Date    string
	Status  string
	Expiry  int64
	Score   int
	Times   int
	Version int
}

func (do *QuestionSubmissionDO) toQuestionSubmission(v *domain.QuestionSubmission) (err error) {
	*v = domain.QuestionSubmission{
		Id:      do.Id,
		Date:    do.Date,
		Status:  do.Status,
		Expiry:  do.Expiry,
		Score:   do.Score,
		Times:   do.Times,
		Version: do.Version,
	}

	v.Account, err = domain.NewAccount(do.Account)

	return
}

func (impl aiquestion) toQuestionSubmissionDO(
	q *domain.QuestionSubmission,
	do *QuestionSubmissionDO,
) {
	*do = QuestionSubmissionDO{
		Id:      q.Id,
		Account: q.Account.Account(),
		Date:    q.Date,
		Status:  q.Status,
		Expiry:  q.Expiry,
		Score:   q.Score,
		Times:   q.Times,
		Version: q.Version,
	}
}
