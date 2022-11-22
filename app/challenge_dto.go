package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
)

type CompetitorApplyCmd domain.CompetitorInfo

func (cmd *CompetitorApplyCmd) Validate() error {
	b := cmd.Account != nil &&
		cmd.Name != nil &&
		cmd.Email != nil &&
		cmd.Identity != nil

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *CompetitorApplyCmd) toCompetitor() *domain.CompetitorInfo {
	return (*domain.CompetitorInfo)(cmd)
}

type ChallengeCompetitorInfoDTO struct {
	IsCompetitor bool `json:"is_competitor"`
	Score        int  `json:"score"`
}

type ChoiceQuestionDTO struct {
	Desc    string   `json:"desc"`
	Options []string `json:"options"`
}

type AIQuestionDTO struct {
	Times       int                 `json:"times"`
	Choices     []ChoiceQuestionDTO `json:"choices"`
	Completions []string            `json:"completions"`
	Answers     string              `json:"answers"`
}
