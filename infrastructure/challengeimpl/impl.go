package challengeimpl

import (
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
