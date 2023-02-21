package messages

import (
	"github.com/opensourceways/xihe-server/competition/domain"
)

func (s sender) NotifyCalcScore(*domain.SubmissionMessage) error {
	return nil
}
