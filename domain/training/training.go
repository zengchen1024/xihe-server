package training

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Training interface {
	CreateJob(endpoint string, info *domain.TrainingIndex, t *domain.TrainingConfig) (domain.JobInfo, error)
	DeleteJob(endpoint, jobId string) error
	TerminateJob(endpoint, jobId string) error
	GetLogPreviewURL(endpoint, jobId string) (string, error)
	IsJobDone(status string) bool
	GetFileDownloadURL(endpoint, file string) (string, error)
}
