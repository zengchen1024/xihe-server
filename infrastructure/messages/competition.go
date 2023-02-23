package messages

import (
	"github.com/opensourceways/xihe-server/competition/domain"
)

func (s sender) NotifyCalcScore(v *domain.SubmissionMessage) error {
	return s.send(topics.Submission, v)
}
