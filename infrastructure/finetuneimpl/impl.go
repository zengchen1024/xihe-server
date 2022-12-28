package finetuneimpl

import (
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/finetune"
)

func NewFinetune(cfg *Config) finetune.Finetune {
	return &finetuneImpl{
		doneStatus:         sets.NewString(cfg.JobDoneStatus...),
		canTerminateStatus: sets.NewString(cfg.CanTerminateStatus...),
	}
}

type finetuneImpl struct {
	// cli
	doneStatus         sets.String
	canTerminateStatus sets.String
}

func (impl *finetuneImpl) IsJobDone(status string) bool {
	return impl.doneStatus.Has(status)
}

func (impl *finetuneImpl) CanTerminate(status string) bool {
	return impl.canTerminateStatus.Has(status)
}

func (impl *finetuneImpl) CreateJob(info *domain.FinetuneIndex, t *domain.FinetuneConfig) (
	job domain.FinetuneJobInfo, err error,
) {
	return
}

func (impl *finetuneImpl) DeleteJob(jobId string) error {
	return nil
}

func (impl *finetuneImpl) TerminateJob(jobId string) error {
	return nil
}

func (impl *finetuneImpl) GetLogPreviewURL(jobId string) (string, error) {
	return "", nil
}
