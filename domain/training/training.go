package training

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Training interface {
	CreateJob(endpoint string, info *domain.TrainingIndex, t *domain.TrainingConfig) (domain.JobInfo, error)
	DeleteJob(endpoint, jobId string) error
	TerminateJob(endpoint, jobId string) error
	GetLogDownloadURL(endpoint, jobId string) (string, error)
	IsJobDone(status string) bool
}
