package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type DatasetMapper interface {
	Insert(DatasetDO) (string, error)
}

func NewDatasetRepository(mapper DatasetMapper) repository.Dataset {
	return dataset{mapper}
}

type dataset struct {
	mapper DatasetMapper
}

func (impl dataset) Save(m *domain.Dataset) (r domain.Dataset, err error) {
	if m.Id != "" {
		return
	}

	do := DatasetDO{
		Owner:    m.Owner,
		Name:     m.Name.ProjName(),
		Desc:     m.Desc.ProjDesc(),
		RepoType: m.RepoType.RepoType(),
		Protocol: m.Protocol.ProtocolName(),
	}

	v, err := impl.mapper.Insert(do)
	if err != nil {
		err = convertError(err)
	} else {
		r = *m
		r.Id = v
	}

	return
}

type DatasetDO struct {
	Id       string
	Owner    string
	Name     string
	Desc     string
	RepoType string
	Protocol string
	Tags     []string
	Version  int
}
