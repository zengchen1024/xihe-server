package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type Tags interface {
	List(domain.ResourceType) ([]domain.DomainTags, error)
}
