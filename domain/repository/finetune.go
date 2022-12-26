package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Finetune interface {
	Save(*domain.UserFinetune, int) (string, error)
	Get(*domain.FinetuneIndex) (domain.UserFinetune, error)
	Delete(*domain.FinetuneIndex) error
	List(user domain.Account) ([]domain.FinetuneSummary, int, error)

	SaveJob(*domain.FinetuneIndex, *domain.FinetuneJobInfo) error
	GetJob(*domain.FinetuneIndex) (domain.FinetuneJobInfo, error)

	UpdateJobDetail(*domain.FinetuneIndex, *domain.FinetuneJobDetail) error
}
