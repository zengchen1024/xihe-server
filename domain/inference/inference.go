package inference

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Inference interface {
	Create(*domain.InferenceInfo) error
	ExtendExpiry(*domain.InferenceInfo, int64) error
}
