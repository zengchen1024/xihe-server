package app

import (
	//"io"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type CompetitionService interface {
	Get(cid string, competitor domain.Account) (CompetitionDTO, error)
	List(domain.CompetitionStatus) ([]CompetitionSummaryDTO, error)

	/*
		// get the phase first, then check if can submit,
		// check the role of submitter
		Submit(cid string, fileName string, file io.Reader) error

		ListSubmitts(cid string, competitor domain.Account) (CompetitionResultDTO, error)

		GetTeam(cid string, competitor domain.Account) (CompetitionTeamDTO, error)

		GetRankingList(cid string, phase domain.CompetitionPhase) ([]RankingDTO, error)
	*/
}

func NewCompetitionService(repo repository.Competition) CompetitionService {
	return competitionService{
		repo: repo,
	}
}

type competitionService struct {
	repo repository.Competition
}

func (s competitionService) Get(cid string, competitor domain.Account) (
	dto CompetitionDTO, err error,
) {
	v, b, err := s.repo.Get(cid, competitor)
	if err != nil {
		return
	}

	s.toCompetitionDTO(&v.Competition, &dto)

	dto.CompetitorCount = v.CompetitorCount

	if !b {
		dto.DatasetURL = ""
	}

	return
}

func (s competitionService) List(status domain.CompetitionStatus) (
	dtos []CompetitionSummaryDTO, err error,
) {
	v, err := s.repo.List(status, domain.CompetitionPhasePreliminary)
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]CompetitionSummaryDTO, len(v))

	for i := range v {
		s.toCompetitionSummaryDTO(&v[i].CompetitionSummary, &dtos[i])

		dtos[i].CompetitorCount = v[i].CompetitorCount
	}

	return
}
