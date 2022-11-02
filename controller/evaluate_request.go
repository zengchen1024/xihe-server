package controller

import (
	"github.com/opensourceways/xihe-server/domain"
)

type EvaluateCreateRequest struct {
	Type              string               `json:"type"`
	MomentumScope     domain.EvaluateScope `json:"momentum_scope"`
	BatchSizeScope    domain.EvaluateScope `json:"batch_size_scope"`
	LearningRateScope domain.EvaluateScope `json:"learning_rate_scope"`
}
