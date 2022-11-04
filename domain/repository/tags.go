package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Tags interface {
	List(domainNames []string) ([]domain.DomainTags, error)
}
