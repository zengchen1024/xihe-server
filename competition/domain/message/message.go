package message

import "github.com/opensourceways/xihe-server/competition/domain"

type CalcScoreMessageProducer interface {
	NotifyCalcScore(*domain.SubmissionMessage) error
}
