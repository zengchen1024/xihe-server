package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type DatasetMapper interface {
	Insert(DatasetDO) (string, error)
	Update(DatasetDO) error
	Get(string, string) (DatasetDO, error)
	List(string, DatasetListDO) ([]DatasetDO, error)
	ListUsersDatasets(map[string][]string) ([]DatasetDO, error)
}

func NewDatasetRepository(mapper DatasetMapper) repository.Dataset {
	return dataset{mapper}
}

type dataset struct {
	mapper DatasetMapper
}

func (impl dataset) Save(d *domain.Dataset) (r domain.Dataset, err error) {
	if d.Id != "" {
		if err = impl.mapper.Update(impl.toDatasetDO(d)); err != nil {
			err = convertError(err)
		} else {
			r = *d
			r.Version += 1
		}

		return
	}

	v, err := impl.mapper.Insert(impl.toDatasetDO(d))
	if err != nil {
		err = convertError(err)
	} else {
		r = *d
		r.Id = v
	}

	return
}

func (impl dataset) Get(owner domain.Account, identity string) (r domain.Dataset, err error) {
	v, err := impl.mapper.Get(owner.Account(), identity)
	if err != nil {
		err = convertError(err)
	} else {
		err = v.toDataset(&r)
	}

	return
}

func (impl dataset) List(owner domain.Account, option repository.DatasetListOption) (
	r []domain.Dataset, err error,
) {
	do := DatasetListDO{}

	if option.Name != nil {
		do.Name = option.Name.DatasetName()
	}
	if option.RepoType != nil {
		do.RepoType = option.RepoType.RepoType()
	}

	v, err := impl.mapper.List(owner.Account(), do)
	if err != nil {
		err = convertError(err)

		return
	}

	r = make([]domain.Dataset, len(v))
	for i := range v {
		if err = v[i].toDataset(&r[i]); err != nil {
			return
		}
	}

	return
}

func (impl dataset) FindUserDatasets(opts []repository.UserResourceListOption) (
	[]domain.Dataset, error,
) {
	do := make(map[string][]string)
	for i := range opts {
		do[opts[i].Owner.Account()] = opts[i].Ids
	}

	v, err := impl.mapper.ListUsersDatasets(do)
	if err != nil {
		return nil, convertError(err)
	}

	r := make([]domain.Dataset, len(v))
	for i := range v {
		if err = v[i].toDataset(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl dataset) toDatasetDO(d *domain.Dataset) DatasetDO {
	do := DatasetDO{
		Id:       d.Id,
		Owner:    d.Owner.Account(),
		Name:     d.Name.DatasetName(),
		RepoType: d.RepoType.RepoType(),
		Protocol: d.Protocol.ProtocolName(),
		Tags:     d.Tags,
		RepoId:   d.RepoId,
	}

	if d.Desc != nil {
		do.Desc = d.Desc.ProjDesc()

	}

	return do
}

type DatasetListDO struct {
	Name     string
	RepoType string
}

type DatasetDO struct {
	Id       string
	Owner    string
	Name     string
	Desc     string
	Protocol string
	RepoType string
	RepoId   string
	Tags     []string
	Version  int
}

func (do *DatasetDO) toDataset(r *domain.Dataset) (err error) {
	r.Id = do.Id

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Name, err = domain.NewDatasetName(do.Name); err != nil {
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

	r.RepoId = do.RepoId
	r.Tags = do.Tags
	r.Version = do.Version

	return
}
