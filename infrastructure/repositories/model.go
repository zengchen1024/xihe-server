package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ModelMapper interface {
	Insert(ModelDO) (string, error)
}

func NewModelRepository(mapper ModelMapper) repository.Model {
	return model{mapper}
}

type model struct {
	mapper ModelMapper
}

func (impl model) Save(m *domain.Model) (r domain.Model, err error) {
	if m.Id != "" {
		return
	}

	do := ModelDO{
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

type ModelDO struct {
	Id       string
	Owner    string
	Name     string
	Desc     string
	RepoType string
	Protocol string
	Tags     []string
	Version  int
}
