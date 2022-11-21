package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/challenge"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ChallengeService interface {
	Apply(*CompetitorApplyCmd) error
	GetCompetitor(domain.Account) (ChallengeCompetitorInfo, error)
}

type challengeService struct {
	comptitions []domain.CompetitionIndex
	aiQuestion  string

	comptitionRepo repository.Competition
	aiQuestionRepo repository.AIQuestion
	helper         challenge.Challenge
}

func NewChallengeService(
	helper challenge.Challenge,
) ChallengeService {
	v := helper.GetChallenge()

	s := &challengeService{}

	s.comptitions = make([]domain.CompetitionIndex, len(v.Competition))

	for i, cid := range v.Competition {
		s.comptitions[i] = domain.CompetitionIndex{
			Id:    cid,
			Phase: domain.CompetitionPhasePreliminary,
		}
	}

	s.aiQuestion = v.AIQuestion

	return s
}

func (s *challengeService) Apply(cmd *CompetitorApplyCmd) error {
	c := cmd.toCompetitor()
	for i := range s.comptitions {
		// TODO allow re-apply
		err := s.comptitionRepo.SaveCompetitor(&s.comptitions[i], c)
		if err != nil {
			return err
		}
	}

	// TODO allow re-apply
	return s.aiQuestionRepo.SaveCompetitor(s.aiQuestion, c)
}

func (s *challengeService) GetCompetitor(user domain.Account) (
	ChallengeCompetitorInfo, error,
) {
	dto := ChallengeCompetitorInfo{}

	for i := range s.comptitions {
		isCompetitor, score, err := s.getCompetitorOfCompetition(
			&s.comptitions[i], user,
		)

		if err != nil || !isCompetitor {
			return dto, err
		}

		dto.Score += score
	}

	isCompetitor, score, err := s.getCompetitorOfAIQuestion(s.aiQuestion, user)

	if err == nil && isCompetitor {
		dto.IsCompetitor = true
		dto.Score += score
	}

	return dto, err
}

func (s *challengeService) getCompetitorOfCompetition(
	index *domain.CompetitionIndex, user domain.Account,
) (isCompetitor bool, score int, err error) {

	isCompetitor, submissions, err := s.comptitionRepo.GetCompetitorAndSubmission(
		index, user,
	)
	if err != nil || !isCompetitor {
		return
	}

	score = s.helper.CalcCompetitionScore(submissions)

	return
}

func (s *challengeService) getCompetitorOfAIQuestion(
	cid string, user domain.Account,
) (isCompetitor bool, score int, err error) {

	isCompetitor, scores, err := s.aiQuestionRepo.GetCompetitorAndScores(cid, user)
	if err != nil || !isCompetitor {
		return
	}

	for _, v := range scores {
		if v > score {
			score = v
		}
	}

	return
}
