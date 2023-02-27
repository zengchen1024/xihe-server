package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
)

func (s *competitionService) Apply(cid string, cmd *CompetitorApplyCmd) (code string, err error) {
	competition, err := s.repo.FindCompetition(cid)
	if err != nil {
		return
	}

	if competition.IsOver() {
		err = errors.New("competition is over")

		return
	}

	if competition.IsFinal() {
		err = errors.New("apply on final phase")

		return
	}

	p := cmd.toPlayer(cid)
	if err = s.playerRepo.AddPlayer(&p, 0); err != nil {
		if repoerr.IsErrorDuplicateCreating(err) {
			code = errorCompetitorExists
		}
	}

	return
}

func (s *competitionService) CreateTeam(cid string, cmd *CompetitionTeamCreateCmd) error {
	p, version, err := s.playerRepo.FindPlayer(cid, cmd.User)
	if err != nil {
		return err
	}

	if err := p.CreateTeam(cmd.Name); err != nil {
		return err
	}

	return s.playerRepo.AddPlayer(&p, version)
}

func (s *competitionService) JoinTeam(cid string, cmd *CompetitionTeamJoinCmd) (code string, err error) {
	me, pv, err := s.playerRepo.FindPlayer(cid, cmd.User)
	if err != nil {
		return
	}

	team, version, err := s.playerRepo.FindPlayer(cid, cmd.Leader)
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			code = errorNoCorrespondingTeam
		}

		return
	}

	if err = me.JoinTo(&team); err != nil {
		if domain.IsErrorTeamMembersEnough(err) {
			code = errorTeamMembersEnough
		}

		return
	}

	err = s.playerRepo.AddMember(
		repository.PlayerVersion{
			Player:  &team,
			Version: version,
		},
		repository.PlayerVersion{
			Player:  &me,
			Version: pv,
		},
	)

	return
}

func (s *competitionService) GetMyTeam(cid string, user types.Account) (
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

	dto.Name = p.Name()

	m := p.Members()
	members := make([]CompetitionTeamMemberDTO, p.CompetitorsCount())
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

func (s *competitionService) ChangeTeamName(cid string, cmd *CmdToChangeCompetitionTeamName) error {
	p, version, err := s.playerRepo.FindPlayer(cid, cmd.User)
	if err != nil {
		return err
	}

	if err = p.ChangeTeamName(cmd.Name); err != nil {
		return err
	}

	return s.playerRepo.SaveTeamName(&p, version)
}

func (s *competitionService) TransferLeader(cid string, cmd *CmdToTransferTeamLeader) error {
	p, version, err := s.playerRepo.FindPlayer(cid, cmd.Leader)
	if err != nil {
		return err
	}

	if err = p.TransferLeader(cmd.User); err != nil {
		return err
	}

	return s.playerRepo.SavePlayer(&p, version)
}
