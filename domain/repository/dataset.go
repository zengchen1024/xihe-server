package repository

import (
	"github.com/opensourceways/xihe-server/domain"
)

type DatasetListOption struct {
	Name     domain.DatasetName
	RepoType domain.RepoType
}

type Dataset interface {
	Save(*domain.Dataset) (domain.Dataset, error)
	Get(domain.Account, string) (domain.Dataset, error)
	GetByName(domain.Account, domain.DatasetName) (domain.Dataset, error)
	List(domain.Account, DatasetListOption) ([]domain.Dataset, error)
	FindUserDatasets([]UserResourceListOption) ([]domain.Dataset, error)
}
