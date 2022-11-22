package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type AIQuestionMapper interface {
	GetCompetitorAndSubmission(string, string) (bool, []int, error)

	SaveCompetitor(string, *CompetitorInfoDO) error
}

func NewAIQuestionRepository(mapper AIQuestionMapper) repository.AIQuestion {
	return aiquestion{mapper}
}

type aiquestion struct {
	mapper AIQuestionMapper
}

func (impl aiquestion) GetCompetitorAndScores(qid string, competitor domain.Account) (
	bool, []int, error,
) {
	return impl.mapper.GetCompetitorAndSubmission(qid, competitor.Account())
}

func (impl aiquestion) SaveCompetitor(qid string, competitor *domain.CompetitorInfo) error {
	do := new(CompetitorInfoDO)
	toCompetitorInfoDO(competitor, do)

	return impl.mapper.SaveCompetitor(qid, do)
}
