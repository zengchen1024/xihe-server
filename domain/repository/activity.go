package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type ActivityFindOption struct {
}

type Activity interface {
	Save(*domain.UserActivity) error
	Find(domain.Account, ActivityFindOption) ([]domain.Activity, error)
}
