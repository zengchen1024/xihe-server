package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Training interface {
	Save(*domain.UserTraining) (domain.UserTraining, error)

	SetJob(*domain.Job) error
}
