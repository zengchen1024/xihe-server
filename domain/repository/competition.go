package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type CompetitionListOption struct {
	Status     domain.CompetitionStatus
	Phase      domain.CompetitionPhase
	Competitor domain.Account
}

type CompetitionSummary struct {
	domain.CompetitionSummary
	CompetitorCount int
}

type CompetitionInfo struct {
	domain.Competition
	CompetitorCount int
}

type Competition interface {
	List(*CompetitionListOption) ([]CompetitionSummary, error)
	Get(*domain.CompetitionIndex, domain.Account) (CompetitionInfo, domain.CompetitorInfo, error)

	GetTeam(*domain.CompetitionIndex, domain.Account) ([]domain.Competitor, error)

	GetResult(*domain.CompetitionIndex) (
		order domain.CompetitionScoreOrder,
		teams []domain.CompetitionTeam,
		results []domain.CompetitionSubmission, err error,
	)

	GetSubmisstions(cid string, c domain.Account) (
		domain.CompetitionRepo, []domain.CompetitionSubmission, error,
	)

	SaveSubmission(*domain.CompetitionIndex, *domain.CompetitionSubmission) (string, error)
}
