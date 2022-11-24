package challenge

import "github.com/opensourceways/xihe-server/domain"

type AIQuestionInfo struct {
	AIQuestionId   string
	QuestionPoolId string
	Timeout        int // minute
	RetryTimes     int
}

type ChallengeInfo struct {
	Competition []string

	AIQuestionInfo
}

type Challenge interface {
	GetChallenge() ChallengeInfo
	CalcCompetitionScore([]domain.CompetitionSubmissionInfo) int
	CalcAIQuestionScore(result, answer []string) int

	GenAIQuestionNums() (choice, completion []int)
}
