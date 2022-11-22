package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type AIQuestion interface {
	GetCompetitorAndScores(string, domain.Account) (bool, []int, error)

	SaveCompetitor(string, *domain.CompetitorInfo) error
}
