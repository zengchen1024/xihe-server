package challenge

import "github.com/opensourceways/xihe-server/domain"

type ChallengeInfo struct {
	Competition []string
	AIQuestion  string
}

type Challenge interface {
	GetChallenge() ChallengeInfo
	CalcCompetitionScore([]domain.CompetitionSubmissionInfo) int

	GenAIQuestionNums() (choice, completion []int)
}
