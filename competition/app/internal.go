package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
)

type CompetitionSubmissionUpdateCmd = domain.CompetitionSubmissionInfo

// Internal Service
type CompetitionInternalService interface {
	UpdateSubmission(*CompetitionSubmissionUpdateCmd) error
}

func NewCompetitionInternalService(repo repository.Work) CompetitionInternalService {
	return competitionInternalService{
		repo: repo,
	}
}

type competitionInternalService struct {
	repo repository.Work
}

func (s competitionInternalService) UpdateSubmission(cmd *CompetitionSubmissionUpdateCmd) error {
	w, version, err := s.repo.FindWork(cmd.Index, cmd.Phase)
	if err != nil {
		return err
	}

	submission := w.UpdateSubmission(cmd)
	if submission == nil {
		return errors.New("no corresponding submission")
	}

	v := domain.PhaseSubmission{
		Phase:      cmd.Phase,
		Submission: *submission,
	}

	return s.repo.SaveSubmission(&w, &v, version)
}
