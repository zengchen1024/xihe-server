package domain

import (
	"github.com/opensourceways/xihe-server/competition/domain/uploader"
)

type CompetitionSubmissionInfo struct {
	Id     string
	Status string
	Score  float32

	Index WorkIndex
	Phase CompetitionPhase
}

func (info *CompetitionSubmissionInfo) IsSuccess() bool {
	return info.Status == competitionSubmissionStatusSuccess
}

type Submission struct {
	Id       string
	SubmitAt int64
	OBSPath  string
	Status   string
	Score    float32
}

type CompetitionSubmission struct {
	Submission

	Phase CompetitionPhase
}

func (info *CompetitionSubmission) IsSuccess() bool {
	return info.Status == competitionSubmissionStatusSuccess
}

type SubmissionService struct {
	uploader uploader.Uploader
	// repo
}

func (s *SubmissionService) Submit() error {
	// upload first
	// create submission
	// save repo

	return nil
}
