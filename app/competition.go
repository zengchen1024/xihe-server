package app

import (
	//"io"
	"sort"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type CompetitionListCMD = repository.CompetitionListOption

type CompetitionService interface {
	Get(cid string, competitor domain.Account) (CompetitionDTO, error)
	List(*CompetitionListCMD) ([]CompetitionSummaryDTO, error)

	/*
		// get the phase first, then check if can submit,
		// check the role of submitter
		Submit(cid string, fileName string, file io.Reader) error

		ListSubmitts(cid string, competitor domain.Account) (CompetitionResultDTO, error)
	*/

	GetTeam(cid string, competitor domain.Account) (CompetitionTeamDTO, error)

	GetRankingList(cid string, phase domain.CompetitionPhase) ([]RankingDTO, error)
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
	index := domain.CompetitionIndex{
		Id:    cid,
		Phase: domain.CompetitionPhasePreliminary,
	}

	v, b, err := s.repo.Get(&index, competitor)
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

func (s competitionService) List(cmd *CompetitionListCMD) (
	dtos []CompetitionSummaryDTO, err error,
) {
	cmd.Phase = domain.CompetitionPhasePreliminary

	v, err := s.repo.List(cmd)
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

func (s competitionService) GetTeam(cid string, competitor domain.Account) (
	dto CompetitionTeamDTO, err error,
) {
	index := domain.CompetitionIndex{
		Id:    cid,
		Phase: domain.CompetitionPhasePreliminary,
	}

	v, err := s.repo.GetTeam(&index, competitor)
	if err != nil {
		return
	}

	if name := v[0].Team.Name; name != nil {
		dto.Name = name.TeamName()
	}

	members := make([]CompetitionTeamMemberDTO, len(v))
	for i := range v {
		item := &v[i]

		members[i] = CompetitionTeamMemberDTO{
			Name:  item.Name.CompetitorName(),
			Email: item.Email.Email(),
		}

		if item.TeamRole != nil {
			members[i].Role = item.TeamRole.TeamRole()
		}
	}

	return
}

func (s competitionService) GetRankingList(cid string, phase domain.CompetitionPhase) (
	dtos []RankingDTO, err error,
) {
	index := domain.CompetitionIndex{
		Id:    cid,
		Phase: domain.CompetitionPhasePreliminary,
	}

	order, teams, results, err := s.repo.GetResult(&index)
	if err != nil || len(results) == 0 {
		return
	}

	rs := map[string]*domain.CompetitionResult{}

	for i := range results {
		item := &results[i]

		k := item.Key()

		if v, ok := rs[k]; !ok || order.IsBetterThanB(item.Score, v.Score) {
			rs[k] = item
		}
	}

	// sort
	i := 0
	rl := make([]*domain.CompetitionResult, len(rs))
	for _, v := range rs {
		rl[i] = v
		i++
	}

	sort.Slice(rl, func(i, j int) bool {
		return order.IsBetterThanB(rl[i].Score, rl[j].Score)
	})

	// result
	tm := map[string]string{}
	for i := range teams {
		tm[teams[i].Id] = teams[i].Name.TeamName()
	}

	dtos = make([]RankingDTO, len(rl))
	for i := range rl {
		item := rl[i]

		dtos[i] = RankingDTO{
			Score:    item.Score,
			SubmitAt: utils.ToDate(item.SubmitAt),
		}

		if item.IsTeamWork() {
			dtos[i].TeamName = tm[item.TeamId]
		} else {
			dtos[i].TeamName = item.Individual.CompetitorName()
		}
	}

	return
}
