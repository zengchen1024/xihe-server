package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
)

func (s competitionService) Apply(cid string, cmd *CompetitorApplyCmd) (code string, err error) {
	competition, err := s.repo.FindCompetition(cid)
	if err != nil {
		return
	}

	if competition.IsOver() {
		code = errorIsOver
		err = errors.New("competition is over")

		return
	}

	if competition.IsFinal() {
		code = errorApplyOnFinalPhase
		err = errors.New("apply on final phase")

		return
	}

	p := cmd.toPlayer(cid)
	if err = s.playerRepo.SavePlayer(&p, 0); err != nil {
		if repoerr.IsErrorDuplicateCreating(err) {
			code = errorCompetitorExists
		}
	}

	return
}

func (s competitionService) CreateTeam(cid string, cmd *CompetitionTeamCreateCmd) (
	code string, err error,
) {
	p, version, err := s.playerRepo.FindPlayer(cid, cmd.User)
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			code = errorNotACompetitor
		}

		return
	}

	if err = p.CreateTeam(cmd.Name); err == nil {
		err = s.playerRepo.SavePlayer(&p, version)
	}

	return
}

func (s competitionService) GetTeam(cid string, user types.Account) (
	dto CompetitionTeamDTO, code string, err error,
) {
	p, _, err := s.playerRepo.FindPlayer(cid, user)
	if err != nil {
		return
	}

	if !p.IsATeam() {
		code = errorNotATeam
		err = errors.New("not a team")

		return
	}

	dto.Name = p.Team.Name.TeamName()

	m := p.Team.Members
	members := make([]CompetitionTeamMemberDTO, len(m)+1)
	for i := range m {
		item := &m[i]
		members[i+1] = CompetitionTeamMemberDTO{
			Name:  item.Name.CompetitorName(),
			Email: item.Email.Email(),
		}
	}

	leader := &p.Leader
	members[0] = CompetitionTeamMemberDTO{
		Name:  leader.Name.CompetitorName(),
		Email: leader.Email.Email(),
		Role:  domain.TeamLeaderRole(),
	}

	dto.Members = members

	return
}
