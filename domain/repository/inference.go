package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type InferenceSummary struct {
	Id string

	domain.InferenceDetail
}

type Inference interface {
	Save(*domain.Infereance, int) (string, error)
	UpdateExpiry(string, int64) error
	UpdateDetail(string, *domain.InferenceDetail) error
	FindInstance(*domain.InferenceInfo) (InferenceSummary, error)
	FindInstances(*domain.InferenceIndex) ([]InferenceSummary, int, error)
}
