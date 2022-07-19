package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type DatasetMapper interface {
	Insert(DatasetDO) (string, error)
	Get(string, string) (DatasetDO, error)
	List(string, DatasetListDO) ([]DatasetDO, error)
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
		Owner:    m.Owner.Account(),
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
		do.Name = option.Name.ProjName()
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

type DatasetListDO struct {
	Name string
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

func (do *DatasetDO) toDataset(r *domain.Dataset) (err error) {
	r.Id = do.Id

	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

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
