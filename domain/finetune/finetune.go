package finetune

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Finetune interface {
	CreateJob(info *domain.FinetuneIndex, t *domain.FinetuneConfig) (domain.FinetuneJobInfo, error)
	DeleteJob(jobId string) error
	TerminateJob(jobId string) error
	GetLogPreviewURL(jobId string) (string, error)
	IsJobDone(status string) bool
	CanTerminate(status string) bool
}
