package app

import (
	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	"github.com/opensourceways/xihe-server/competition/domain/uploader"
	types "github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	repoerr "github.com/opensourceways/xihe-server/domain/repository"
)

type CompetitionListCMD struct {
	repository.CompetitionListOption

	User types.Account
}

type CompetitionService interface {
	// player
	Apply(string, *CompetitorApplyCmd) (string, error)
	CreateTeam(cid string, cmd *CompetitionTeamCreateCmd) (string, error)
	GetTeam(cid string, competitor types.Account) (CompetitionTeamDTO, string, error)

	// competition
	Get(cid string, competitor types.Account) (UserCompetitionDTO, error)
	List(*CompetitionListCMD) ([]CompetitionSummaryDTO, error)

	// work
	Submit(*CompetitionSubmitCMD) (CompetitionSubmissionDTO, string, error)
	GetSubmissions(string, types.Account) (CompetitionSubmissionsDTO, error)
	GetRankingList(string) (CompetitonRankingDTO, error)
	AddRelatedProject(*CompetitionAddReleatedProjectCMD) error
}

func NewCompetitionService(
	repo repository.Competition,
	workRepo repository.Work,
	playerRepo repository.Player,
	sender message.Sender,
	uploader uploader.SubmissionFileUploader,
) CompetitionService {
	return competitionService{
		repo:             repo,
		workRepo:         workRepo,
		playerRepo:       playerRepo,
		sender:           sender,
		submissionServie: domain.NewSubmissionService(uploader),
	}
}

type competitionService struct {
	repo             repository.Competition
	workRepo         repository.Work
	playerRepo       repository.Player
	sender           message.Sender
	submissionServie domain.SubmissionService
}

// show competition detail
func (s competitionService) Get(cid string, user types.Account) (
	dto UserCompetitionDTO, err error,
) {
	c, err := s.repo.FindCompetition(cid)
	if err != nil {
		return
	}
	s.toCompetitionDTO(&c, &dto.CompetitionDTO)

	// competitors count
	dto.CompetitorCount, err = s.playerRepo.CompetitorsCount(cid)
	if err != nil {
		return
	}

	// competitor info
	if user == nil {
		return
	}

	p, _, err := s.playerRepo.FindPlayer(cid, user)
	if err != nil {
		if repoerr.IsErrorResourceNotExists(err) {
			err = nil
			dto.DatasetURL = ""
		}

		return
	}
	dto.IsCompetitor = true

	if p.IsATeam() {
		dto.TeamId = p.Id
		dto.TeamRole = p.RoleOfCurrentCompetitor()
	}

	return
}

func (s competitionService) List(cmd *CompetitionListCMD) (
	dtos []CompetitionSummaryDTO, err error,
) {
	if cmd.User != nil {
		return s.getCompetitionsUserApplied(cmd)
	}

	return s.listCompetitions(&cmd.CompetitionListOption)
}

func (s competitionService) listCompetitions(opt *repository.CompetitionListOption) (
	dtos []CompetitionSummaryDTO, err error,
) {
	v, err := s.repo.FindCompetitions(opt)
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]CompetitionSummaryDTO, len(v))

	for i := range v {
		s.toCompetitionSummaryDTO(&v[i], &dtos[i])

		dtos[i].CompetitorCount, err = s.playerRepo.CompetitorsCount(v[i].Id)
		if err != nil {
			return
		}
	}

	return
}

func (s competitionService) getCompetitionsUserApplied(cmd *CompetitionListCMD) (
	[]CompetitionSummaryDTO, error,
) {
	v, err := s.listCompetitions(&cmd.CompetitionListOption)
	if err != nil {
		return nil, err
	}

	cs, err := s.playerRepo.FindCompetitionsUserApplied(cmd.User)
	if err != nil {
		return nil, err
	}
	m := make(map[string]bool)
	for _, item := range cs {
		m[item] = true
	}

	dtos := make([]CompetitionSummaryDTO, 0, len(v))
	for i := range v {
		if m[v[i].Id] {
			dtos = append(dtos, v[i])
		}
	}

	return dtos, nil
}
