package inference

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Inference interface {
	Create(*domain.InferenceInfo, string) error
	ExtendExpiry(*domain.InferenceIndex, int64) error
}
