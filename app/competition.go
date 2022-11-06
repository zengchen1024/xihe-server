package app

import (
	"io"

	"github.com/opensourceways/xihe-server/domain"
)

type CompetitionService interface {
	Get(string, competitor domain.Account) (CompetitionDTO, error)
	List(domain.CompetitionStatus) ([]CompetitionSummaryDTO, error)

	// get the phase first, then check if can submit,
	// check the role of submitter
	Submit(cid string, fileName string, file io.Reader) error

	ListSubmitts(cid string, competitor domain.Account) (CompetitionResultDTO, error)

	GetTeam(cid string, competitor domain.Account) (CompetitionTeamDTO, error)

	GetRankingList(cid string, phase domain.CompetitionPhase) ([]RankingDTO, error)
}
