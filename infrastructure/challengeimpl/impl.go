package challengeimpl

import (
	"strings"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/challenge"
	"github.com/opensourceways/xihe-server/utils"
)

func NewChallenge(cfg *Config) challenge.Challenge {
	return &challengeImpl{*cfg}
}

type challengeImpl struct {
	cfg Config
}

func (impl *challengeImpl) GetChallenge() challenge.ChallengeInfo {
	info := &impl.cfg.AIQuestion

	return challenge.ChallengeInfo{
		Competition: impl.cfg.Competitions,

		AIQuestionInfo: challenge.AIQuestionInfo{
			AIQuestionId:   info.AIQuestionId,
			QuestionPoolId: info.QuestionPoolId,
			Timeout:        info.Timeout,
			RetryTimes:     info.RetryTimes,
		},
	}
}

func (impl *challengeImpl) CalcCompetitionScore(
	submissions []domain.CompetitionSubmissionInfo,
) int {
	for i := range submissions {
		if submissions[i].IsSuccess() {
			return impl.cfg.CompetitionSuccessScore
		}
	}

	return 0
}

func (impl *challengeImpl) CalcCompetitionScoreForAll(
	submissions []domain.CompetitionSubmission,
) map[string]int {
	r := map[string]int{}

	for i := range submissions {
		item := &submissions[i]

		if item.Individual == nil {
			continue
		}

		name := item.Individual.Account()
		if _, ok := r[name]; ok {
			continue
		}

		if item.IsSuccess() {
			r[name] = impl.cfg.CompetitionSuccessScore
		}
	}

	return r
}

func (impl *challengeImpl) GenAIQuestionNums() (choice, completion []int) {
	cfg := impl.cfg.AIQuestion

	choice = utils.GenRandoms(cfg.ChoiceQuestionsCount, cfg.ChoiceQuestionsNum)
	completion = utils.GenRandoms(cfg.CompletionQuestionsCount, cfg.CompletionQuestionsNum)

	return
}

func (impl *challengeImpl) CalcAIQuestionScore(result, answer []string) (score int) {
	cfg := impl.cfg.AIQuestion

	num := cfg.ChoiceQuestionsNum
	for i := 0; i < num; i++ {
		if result[i] == answer[i] {
			score += cfg.ChoiceQuestionsScore
		}
	}

	total := num + cfg.CompletionQuestionsNum
	for i := num; i < total; i++ {
		if impl.formatCompletionAnswer(result[i]) == impl.formatCompletionAnswer(answer[i]) {
			score += cfg.CompletionQuestionsScore
		}
	}

	return
}

func (impl *challengeImpl) formatCompletionAnswer(v string) string {
	str := strings.ReplaceAll(v, " ", "")
	str = strings.ReplaceAll(str, "'", "")
	str = strings.ReplaceAll(str, "\"", "")

	return str
}
