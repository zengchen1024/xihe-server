package inferenceimpl

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/inference"
)

func NewInference(cfg *Config) inference.Inference {
	return inferenceImpl{}
}

type inferenceImpl struct {
}

func (impl inferenceImpl) Create(*domain.InferenceInfo) error {
	return nil
}

func (impl inferenceImpl) ExtendExpiry(*domain.InferenceInfo, int64) error {
	return nil
}
