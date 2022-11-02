package evaluate

import (
	"github.com/opensourceways/xihe-server/domain"
)

type EvaluateInfo struct {
	*domain.Evaluate
	OBSPath      string
	SurvivalTime int
}

type Evaluate interface {
	Create(*EvaluateInfo) error
}
