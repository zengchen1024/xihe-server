package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type TrainingInfo struct {
	User       domain.Account
	ProjectId  string
	TrainingId string
}

type Training interface {
	Save(*domain.UserTraining) (domain.UserTraining, error)
	Get(*TrainingInfo) (domain.UserTraining, error)
	Delete(*TrainingInfo) error
	List(user domain.Account, projectId string) ([]domain.TrainingSummary, error)

	SaveJob(*TrainingInfo, *domain.JobInfo) error
	GetJob(*TrainingInfo) (domain.JobInfo, error)

	UpdateJobDetail(*TrainingInfo, *domain.JobDetail) error
}
