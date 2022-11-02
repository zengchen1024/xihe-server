package inference

import (
	"github.com/opensourceways/xihe-server/domain"
)

type InferenceInfo struct {
	*domain.InferenceInfo
	UserToken    string
	SurvivalTime int
}

type Inference interface {
	Create(*InferenceInfo) error
	ExtendSurvivalTime(index *domain.InferenceIndex, timeToExtend int) error
}
