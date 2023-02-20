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

	Repo        string
	PlayerName  string
	Final       []CompetitionSubmission
	Preliminary []CompetitionSubmission
}

func NewWork(cid string, p *Player) Work {
	return Work{
		WorkIndex:  NewWorkIndex(cid, p.Id),
		PlayerName: p.Name(),
	}
}

func (w *Work) Submissions(phase CompetitionPhase) []CompetitionSubmission {
	if phase.IsPreliminary() {
		return w.Preliminary
	}

	if phase.IsFinal() {
		return w.Final
	}

	return nil
}

func (w *Work) BestOne(phase CompetitionPhase, order CompetitionScoreOrder) (
	r *CompetitionSubmission,
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

func (w *Work) AddSubmission(v *CompetitionSubmission) {
	if v.Phase.IsPreliminary() {
		w.Preliminary = append(w.Preliminary, *v)
	}

	if v.Phase.IsFinal() {
		w.Final = append(w.Final, *v)
	}
}

func (w *Work) UpdateSubmission(info *CompetitionSubmissionInfo) *CompetitionSubmission {
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

func (w *Work) SubmissionObsPathPrefix(phase CompetitionPhase) string {
	return fmt.Sprintf(
		"%s/%s/%s",
		w.CompetitionId,
		phase.CompetitionPhase(),
		w.PlayerId,
	)
}

func (w *Work) HasSubmittedToday(phase CompetitionPhase, t int64) bool {
	v := utils.ToDate(t)
	submissions := w.Submissions(phase)
	for i := range submissions {
		if utils.ToDate(submissions[i].SubmitAt) == v {
			return true
		}
	}

	return false
}
