package app

import (
	"path/filepath"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type CompetitionSummaryDTO struct {
	CompetitorCount int    `json:"count"`
	Bonus           int    `json:"bonus"`
	Id              string `json:"id"`
	Name            string `json:"name"`
	Host            string `json:"host"`
	Desc            string `json:"desc"`
	Status          string `json:"status"`
	Poster          string `json:"poster"`
	Duration        string `json:"duration"`
}

type CompetitionDTO struct {
	CompetitionSummaryDTO

	Doc        string `json:"doc"`
	DatasetDoc string `json:"dataset_doc"`
	DatasetURL string `json:"dataset_url"`
}

// ranking
type RankingDTO struct {
	Score    float32 `json:"score"`
	TeamName string  `json:"team_name"`
	SubmitAt string  `json:"submit_at"`
}

// team
type CompetitionTeamDTO struct {
	Name    string                     `json:"name"`
	Members []CompetitionTeamMemberDTO `json:"members"`
}

type CompetitionTeamMemberDTO struct {
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
}

// result
type CompetitionSubmissionsDTO struct {
	RelatedProject string                     `json:"project"`
	Details        []CompetitionSubmissionDTO `json:"details"`
}

type CompetitionSubmissionDTO struct {
	SubmitAt string  `json:"submit_at"`
	FileName string  `json:"project"`
	Status   string  `json:"status"`
	Score    float32 `json:"score"`
}

func (s competitionService) toCompetitionSubmissionDTO(
	v *domain.CompetitionSubmission, dto *CompetitionSubmissionDTO,
) {
	*dto = CompetitionSubmissionDTO{
		SubmitAt: utils.ToDate(v.SubmitAt),
		FileName: filepath.Base(v.OBSPath),
		Status:   v.Status,
		Score:    v.Score,
	}
}

func (s competitionService) toCompetitionSummaryDTO(
	c *domain.CompetitionSummary, dto *CompetitionSummaryDTO,
) {
	*dto = CompetitionSummaryDTO{
		Bonus:    c.Bonus.CompetitionBonus(),
		Id:       c.Id,
		Name:     c.Name.CompetitionName(),
		Host:     c.Host.CompetitionHost(),
		Desc:     c.Desc.CompetitionDesc(),
		Status:   c.Status.CompetitionStatus(),
		Poster:   c.Poster.URL(),
		Duration: c.Duration.CompetitionDuration(),
	}
}

func (s competitionService) toCompetitionDTO(
	c *domain.Competition, dto *CompetitionDTO,
) {
	s.toCompetitionSummaryDTO(
		&c.CompetitionSummary,
		&dto.CompetitionSummaryDTO,
	)

	dto.Doc = c.Doc.URL()
	dto.DatasetDoc = c.DatasetDoc.URL()
	dto.DatasetURL = c.DatasetURL.URL()
}
