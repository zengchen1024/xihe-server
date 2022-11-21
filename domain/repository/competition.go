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
	Get(*domain.CompetitionIndex, domain.Account) (
		CompetitionInfo, domain.CompetitorSummary, error,
	)

	GetTeam(*domain.CompetitionIndex, domain.Account) ([]domain.Competitor, error)

	GetResult(*domain.CompetitionIndex) (
		order domain.CompetitionScoreOrder,
		teams []domain.CompetitionTeam,
		results []domain.CompetitionSubmission, err error,
	)

	GetSubmisstions(*domain.CompetitionIndex, domain.Account) (
		domain.CompetitionRepo, []domain.CompetitionSubmission, error,
	)

	SaveSubmission(*domain.CompetitionIndex, *domain.CompetitionSubmission) (string, error)

	UpdateSubmission(*domain.CompetitionIndex, *domain.CompetitionSubmissionInfo) error

	GetCompetitorAndSubmission(*domain.CompetitionIndex, domain.Account) (
		bool, []domain.CompetitionSubmissionInfo, error,
	)

	SaveCompetitor(*domain.CompetitionIndex, *domain.CompetitorInfo) error
}
