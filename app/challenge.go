package app

import (
	"strings"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/challenge"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type ChallengeService interface {
	Apply(*CompetitorApplyCmd) error
	GetCompetitor(domain.Account) (ChallengeCompetitorInfoDTO, error)
	GetAIQuestions(domain.Account) (AIQuestionDTO, error)
}

type challengeService struct {
	comptitions []domain.CompetitionIndex
	aiQuestion  string
	delimiter   string

	competitionRepo repository.Competition
	aiQuestionRepo  repository.AIQuestion
	helper          challenge.Challenge
	encryption      utils.SymmetricEncryption
}

func NewChallengeService(
	competitionRepo repository.Competition,
	aiQuestionRepo repository.AIQuestion,
	helper challenge.Challenge,
) ChallengeService {
	v := helper.GetChallenge()

	s := &challengeService{
		competitionRepo: competitionRepo,
		aiQuestionRepo:  aiQuestionRepo,
		helper:          helper,
		delimiter:       ",",
	}

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
		err := s.competitionRepo.SaveCompetitor(&s.comptitions[i], c)
		if err != nil {
			return err
		}
	}

	// TODO allow re-apply
	return s.aiQuestionRepo.SaveCompetitor(s.aiQuestion, c)
}

func (s *challengeService) GetCompetitor(user domain.Account) (
	ChallengeCompetitorInfoDTO, error,
) {
	dto := ChallengeCompetitorInfoDTO{}

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

	isCompetitor, submissions, err := s.competitionRepo.GetCompetitorAndSubmission(
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

func (s *challengeService) GetAIQuestions(competitor domain.Account) (dto AIQuestionDTO, err error) {
	// TODO check if can gen questions

	return s.genAIQuestions()
}

func (s *challengeService) genAIQuestions() (dto AIQuestionDTO, err error) {
	choice, completion := s.helper.GenAIQuestionNums()
	choices, completions, err := s.aiQuestionRepo.GetQuestions(choice, completion)
	if err != nil {
		return
	}

	n := len(choice)
	answers := make([]string, n+len(completion))
	dto.Choices = make([]ChoiceQuestionDTO, n)

	for i := range choices {
		item := &choices[i]

		dto.Choices[i] = ChoiceQuestionDTO{
			Desc:    item.Desc,
			Options: item.Options,
		}

		answers[i] = item.Answer
	}

	dto.Completions = make([]string, len(completion))

	for i := range completions {
		item := &completions[i]

		dto.Completions[i] = item.Desc
		answers[i+n] = item.Answer
	}

	str, err := s.encryptAnswer(answers)
	if err == nil {
		dto.Answers = str
	}

	return
}

func (s *challengeService) encryptAnswer(answers []string) (string, error) {
	str := strings.Join(answers, s.delimiter)

	v, err := s.encryption.Encrypt([]byte(str))
	if err == nil {
		return string(v), nil
	}

	return "", err
}

func (s *challengeService) decryptAnswer(str string) ([]string, error) {
	v, err := s.encryption.Decrypt([]byte(str))
	if err != nil {
		return nil, err
	}

	return strings.Split(string(v), s.delimiter), nil
}
