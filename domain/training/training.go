package training

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Training interface {
	CreateJob(endpoint string, t *domain.UserTraining) (string, error)
	DeleteJob(endpoint, jobId string) error
	GetJob(endpoint, jobId string) (domain.TrainingInfo, error)
	TerminateJob(endpoint, jobId string) error
	GetLogURL(endpoint, jobId string) (string, error)
}
