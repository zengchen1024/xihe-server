package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ModelMapper interface {
	Insert(ModelDO) (string, error)
	Get(string, string) (ModelDO, error)
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

func (impl model) Get(owner, identity string) (r domain.Model, err error) {
	v, err := impl.mapper.Get(owner, identity)
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toModel(&r)
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

func (do *ModelDO) toModel(r *domain.Model) (err error) {
	r.Id = do.Id
	r.Owner = do.Owner

	if r.Name, err = domain.NewProjName(do.Name); err != nil {
		return
	}

	if r.Desc, err = domain.NewProjDesc(do.Desc); err != nil {
		return
	}

	if r.RepoType, err = domain.NewRepoType(do.RepoType); err != nil {
		return
	}

	if r.Protocol, err = domain.NewProtocolName(do.Protocol); err != nil {
		return
	}

	r.Tags = do.Tags

	r.Version = do.Version

	return
}
