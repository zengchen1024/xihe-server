package mongodb

import "github.com/opensourceways/xihe-server/infrastructure/repositories"

func (col aiquestion) toQuestionSubmissionDoc(
	do *repositories.QuestionSubmissionDO,
	doc *dQuestionSubmission,
) {
	*doc = dQuestionSubmission{
		Id:      do.Id,
		Date:    do.Date,
		Status:  do.Status,
		Account: do.Account,
		Expiry:  do.Expiry,
		Score:   do.Score,
		Times:   do.Times,
	}
}

func (col aiquestion) toQuestionSubmissionDo(
	do *repositories.QuestionSubmissionDO,
	doc *dQuestionSubmission,
) {
	*do = repositories.QuestionSubmissionDO{
		Id:      doc.Id,
		Date:    doc.Date,
		Status:  doc.Status,
		Account: doc.Account,
		Expiry:  doc.Expiry,
		Score:   doc.Score,
		Times:   doc.Times,
		Version: doc.Version,
	}
}

func (col aiquestion) toChoiceQuestionDO(
	do *repositories.ChoiceQuestionDO,
	doc *dChoiceQuestion,
) {
	*do = repositories.ChoiceQuestionDO{
		Desc:    doc.Desc,
		Answer:  doc.Answer,
		Options: doc.Options,
	}
}

func (col aiquestion) toCompletionQuestionDO(
	do *repositories.CompletionQuestionDO,
	doc *dCompletionQuestion,
) {
	*do = repositories.CompletionQuestionDO{
		Desc:   doc.Desc,
		Answer: doc.Answer,
	}
}
