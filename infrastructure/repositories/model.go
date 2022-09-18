package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ModelMapper interface {
	Insert(ModelDO) (string, error)
	Update(ModelDO) error
	Get(string, string) (ModelDO, error)
	GetByName(string, string) (ModelDO, error)
	List(string, ResourceListDO) ([]ModelDO, error)
	ListUsersModels(map[string][]string) ([]ModelDO, error)

	AddLike(string, string) error
	RemoveLike(string, string) error
}

func NewModelRepository(mapper ModelMapper) repository.Model {
	return model{mapper}
}

type model struct {
	mapper ModelMapper
}

func (impl model) Save(m *domain.Model) (r domain.Model, err error) {
	if m.Id != "" {
		if err = impl.mapper.Update(impl.toModelDO(m)); err != nil {
			err = convertError(err)
		} else {
			r = *m
			r.Version += 1
		}

		return
	}

	v, err := impl.mapper.Insert(impl.toModelDO(m))
	if err != nil {
		err = convertError(err)
	} else {
		r = *m
		r.Id = v
	}

	return
}

func (impl model) Get(owner domain.Account, identity string) (r domain.Model, err error) {
	v, err := impl.mapper.Get(owner.Account(), identity)
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toModel(&r)
	}

	return
}

func (impl model) GetByName(owner domain.Account, name domain.ModelName) (
	r domain.Model, err error,
) {
	v, err := impl.mapper.GetByName(owner.Account(), name.ModelName())
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toModel(&r)
	}

	return
}

func (impl model) List(owner domain.Account, option repository.ResourceListOption) (
	r []domain.Model, err error,
) {
	do := ResourceListDO{
		Name: option.Name,
	}
	if option.RepoType != nil {
		do.RepoType = option.RepoType.RepoType()
	}

	v, err := impl.mapper.List(owner.Account(), do)
	if err != nil {
		err = convertError(err)

		return
	}

	r = make([]domain.Model, len(v))
	for i := range v {
		if err = v[i].toModel(&r[i]); err != nil {
			return
		}
	}

	return
}

func (impl model) FindUserModels(opts []repository.UserResourceListOption) (
	[]domain.Model, error,
) {
	do := make(map[string][]string)
	for i := range opts {
		do[opts[i].Owner.Account()] = opts[i].Ids
	}

	v, err := impl.mapper.ListUsersModels(do)
	if err != nil {
		return nil, convertError(err)
	}

	r := make([]domain.Model, len(v))
	for i := range v {
		if err = v[i].toModel(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl model) toModelDO(m *domain.Model) ModelDO {
	return ModelDO{
		Id:       m.Id,
		Owner:    m.Owner.Account(),
		Name:     m.Name.ModelName(),
		Desc:     m.Desc.ResourceDesc(),
		RepoType: m.RepoType.RepoType(),
		Protocol: m.Protocol.ProtocolName(),
		Tags:     m.Tags,
		RepoId:   m.RepoId,
		Version:  m.Version,
	}
}

type ModelListDO struct {
	Name     string
	RepoType string
}

type ModelDO struct {
	Id        string
	Owner     string
	Name      string
	Desc      string
	Protocol  string
	RepoType  string
	RepoId    string
	Tags      []string
	Version   int
	LikeCount int
}

func (do *ModelDO) toModel(r *domain.Model) (err error) {
	r.Id = do.Id

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Name, err = domain.NewModelName(do.Name); err != nil {
		return
	}

	if r.Desc, err = domain.NewResourceDesc(do.Desc); err != nil {
		return
	}

	if r.RepoType, err = domain.NewRepoType(do.RepoType); err != nil {
		return
	}

	if r.Protocol, err = domain.NewProtocolName(do.Protocol); err != nil {
		return
	}

	r.RepoId = do.RepoId
	r.Tags = do.Tags
	r.Version = do.Version
	r.LikeCount = do.LikeCount

	return
}
