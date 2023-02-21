package domain

import (
	"fmt"

	"github.com/opensourceways/xihe-server/utils"
)

type WorkIndex struct {
	PlayerId      string
	CompetitionId string
}

func NewWorkIndex(cid, pid string) WorkIndex {
	return WorkIndex{
		PlayerId:      pid,
		CompetitionId: cid,
	}
}

type Work struct {
	WorkIndex

	PlayerName string

	Repo        string
	Final       []Submission
	Preliminary []Submission
}

func NewWork(cid string, p *Player) Work {
	return Work{
		WorkIndex:  NewWorkIndex(cid, p.Id),
		PlayerName: p.Name(),
	}
}

func (w *Work) Submissions(phase CompetitionPhase) []Submission {
	if phase.IsPreliminary() {
		return w.Preliminary
	}

	if phase.IsFinal() {
		return w.Final
	}

	return nil
}

func (w *Work) BestOne(phase CompetitionPhase, order CompetitionScoreOrder) (
	r *Submission,
) {
	submissions := w.Submissions(phase)
	for i := range submissions {
		item := &submissions[i]

		if !item.IsSuccess() {
			continue
		}

		if r == nil || order.IsBetterThanB(item.Score, r.Score) {
			r = item
		}
	}

	return
}

func (w *Work) submissionObsPathPrefix(phase CompetitionPhase) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		w.CompetitionId,
		phase.CompetitionPhase(),
		w.PlayerId,
	)
}

func (w *Work) HasSubmittedToday(phase CompetitionPhase) bool {
	today := utils.Date()
	submissions := w.Submissions(phase)
	for i := range submissions {
		if utils.ToDate(submissions[i].SubmitAt) == today {
			return true
		}
	}

	return false
}

func (w *Work) NewSubmissionMessage(s *PhaseSubmission) SubmissionMessage {
	return SubmissionMessage{
		Index:   w.WorkIndex,
		Phase:   s.Phase,
		Id:      s.Submission.Id,
		OBSPath: s.Submission.OBSPath,
	}
}

func (w *Work) UpdateSubmission(info *CompetitionSubmissionInfo) *Submission {
	submissions := w.Submissions(info.Phase)
	for i := range submissions {
		if item := &submissions[i]; item.Id == info.Id {
			item.Status = info.Status
			item.Score = info.Score

			return &submissions[i]
		}
	}

	return nil
}
