package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Training interface {
	Save(*domain.UserTraining, int) (domain.UserTraining, error)
	Get(*domain.TrainingInfo) (domain.UserTraining, error)
	Delete(*domain.TrainingInfo) error
	List(user domain.Account, projectId string) ([]domain.TrainingSummary, int, error)

	SaveJob(*domain.TrainingInfo, *domain.JobInfo) error
	GetJob(*domain.TrainingInfo) (domain.JobInfo, error)

	UpdateJobDetail(*domain.TrainingInfo, *domain.JobDetail) error
}
