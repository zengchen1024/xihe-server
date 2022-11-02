package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type EvaluateSummary struct {
	Id string

	domain.EvaluateDetail
}

type Evaluate interface {
	Save(*domain.Evaluate, int) (string, error)
	GetStandardEvaluateParms(*domain.EvaluateIndex) (domain.StandardEvaluateParms, error)
	UpdateDetail(*domain.EvaluateIndex, *domain.EvaluateDetail) error
	FindInstance(*domain.EvaluateIndex) (EvaluateSummary, error)
	FindInstances(index *domain.ResourceIndex, trainingId string) ([]EvaluateSummary, int, error)
}
