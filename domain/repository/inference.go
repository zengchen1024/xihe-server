package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type InferenceSummary struct {
	Id string

	domain.InferenceDetail
}

type Inference interface {
	Save(*domain.Inference, int) (string, error)
	UpdateExpiry(*domain.InferenceIndex, int64) error
	UpdateDetail(*domain.InferenceIndex, *domain.InferenceDetail) error
	FindInstance(*domain.InferenceIndex) (InferenceSummary, error)
	FindInstances(index *domain.ResourceIndex, lastCommit string) ([]InferenceSummary, int, error)
}
