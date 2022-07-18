package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Model interface {
	Save(*domain.Model) (domain.Model, error)
	Get(string, string) (domain.Model, error)
}
