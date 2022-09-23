package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

func (impl dataset) AddLike(owner domain.Account, rid string) error {
	err := impl.mapper.AddLike(owner.Account(), rid)
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl dataset) RemoveLike(owner domain.Account, rid string) error {
	err := impl.mapper.RemoveLike(owner.Account(), rid)
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl dataset) UpdateProperty(info *repository.DatasetPropertyUpdateInfo) error {
	p := &info.Property

	do := DatasetPropertyDO{
		Id:       info.Id,
		Owner:    info.Owner.Account(),
		Version:  info.Version,
		Name:     p.Name.DatasetName(),
		Desc:     p.Desc.ResourceDesc(),
		RepoType: p.RepoType.RepoType(),
		Tags:     p.Tags,
	}

	if err := impl.mapper.UpdateProperty(&do); err != nil {
		return convertError(err)
	}

	return nil
}

type DatasetPropertyDO struct {
	Id      string
	Owner   string
	Version int

	Name     string
	Desc     string
	RepoType string
	Tags     []string
}
