package challengeimpl

import (
	"math/rand"
	"time"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/challenge"
)

func NewChallenge(cfg *Config) challenge.Challenge {
	return &challengeImpl{*cfg}
}

type challengeImpl struct {
	cfg Config
}

func (impl *challengeImpl) GetChallenge() challenge.ChallengeInfo {
	return challenge.ChallengeInfo{
		Competition: impl.cfg.Competitions,
		AIQuestion:  impl.cfg.AIQuestion,
	}
}

func (impl *challengeImpl) CalcCompetitionScore(
	submissions []domain.CompetitionSubmissionInfo,
) int {
	for i := range submissions {
		if submissions[i].Status == impl.cfg.CompetitionSuccessStatus {
			return impl.cfg.CompetitionSuccessScore
		}
	}

	return 0
}

func (impl *challengeImpl) GenAIQuestionNums() (choice, completion []int) {
	cfg := impl.cfg

	choice = impl.genRandoms(cfg.ChoiceQuestionsCount, cfg.ChoiceQuestionsNum)
	completion = impl.genRandoms(cfg.CompletionQuestionsCount, cfg.CompletionQuestionsNum)

	return
}

func (impl *challengeImpl) genRandoms(max, total int) []int {
	// set seed
	rand.Seed(time.Now().UnixNano())

	min := 1
	v := max - min
	i := 0
	m := make(map[int]struct{})
	r := make([]int, total)
	for {
		n := rand.Intn(v) + min

		if _, ok := m[n]; !ok {
			m[n] = struct{}{}
			r[i] = n
			if i++; i == total {
				break
			}
		}
	}

	return r
}
