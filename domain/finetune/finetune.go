package finetune

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Finetune interface {
	CreateJob(endpoint string, info *domain.FinetuneIndex, t *domain.FinetuneConfig) (domain.FinetuneJobInfo, error)
	DeleteJob(endpoint, jobId string) error
	TerminateJob(endpoint, jobId string) error
	GetLogPreviewURL(endpoint, jobId string) (string, error)
	IsJobDone(status string) bool
	CanTerminate(status string) bool
}
