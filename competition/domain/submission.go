package domain

import (
	"fmt"
	"io"
	"strconv"

	"github.com/opensourceways/xihe-server/competition/domain/uploader"
	"github.com/opensourceways/xihe-server/utils"
)

type CompetitionSubmissionInfo struct {
	Index  WorkIndex
	Phase  CompetitionPhase
	Id     string
	Status string
	Score  float32
}

func (info *CompetitionSubmissionInfo) IsSuccess() bool {
	return info.Status == competitionSubmissionStatusSuccess
}

type SubmissionMessage struct {
	Index   WorkIndex
	Phase   CompetitionPhase
	Id      string
	OBSPath string
}

type Submission struct {
	Id       string
	SubmitAt int64
	OBSPath  string
	Status   string
	Score    float32
}

func (info *Submission) IsSuccess() bool {
	return info.Status == competitionSubmissionStatusSuccess
}

type PhaseSubmission struct {
	Phase CompetitionPhase

	Submission
}

type SubmissionService struct {
	uploader uploader.Uploader
}

func NewSubmissionService(v uploader.Uploader) SubmissionService {
	return SubmissionService{v}
}

func (s *SubmissionService) Submit(
	w *Work, phase CompetitionPhase, fileName string, data io.Reader,
) (PhaseSubmission, error) {
	now := utils.Now()

	// upload file
	obspath := fmt.Sprintf(
		"%s/%s_%s",
		w.submissionObsPathPrefix(phase),
		strconv.FormatInt(now, 10), fileName,
	)
	if err := s.uploader.UploadSubmissionFile(data, obspath); err != nil {
		return PhaseSubmission{}, err
	}

	return PhaseSubmission{
		Submission: Submission{
			SubmitAt: now,
			OBSPath:  obspath,
			Status:   "calculating",
		},
		Phase: phase,
	}, nil
}
