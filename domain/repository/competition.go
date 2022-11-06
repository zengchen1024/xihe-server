package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type CompetitionSummary struct {
	domain.CompetitionSummary
	CompetitorCount int
}

type CompetitionInfo struct {
	domain.Competition
	CompetitorCount int
}

type Competition interface {
	List(domain.CompetitionStatus) ([]CompetitionSummary, error)
	Get(string, domain.Account) (CompetitionInfo, bool, error)

	// list all the record on different phase
	//GetCompetitor(cid string, competitor domain.Account)
}
