package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type LikeMapper interface {
	Insert(string, LikeDO) error
	Delete(string, LikeDO) error
	List(string, LikeListDO) ([]LikeDO, error)
	HasLike(string, *ResourceObjectDO) (bool, error)
}

func NewLikeRepository(mapper LikeMapper) repository.Like {
	return like{mapper}
}

type like struct {
	mapper LikeMapper
}

func (impl like) Save(ul *domain.UserLike) error {
	err := impl.mapper.Insert(ul.Owner.Account(), impl.toLikeDO(&ul.Like))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl like) Remove(ul *domain.UserLike) error {
	err := impl.mapper.Delete(ul.Owner.Account(), impl.toLikeDO(&ul.Like))
	if err != nil {
		err = convertError(err)
	}

	return err
}

func (impl like) Find(owner domain.Account, opt repository.LikeFindOption) (
	[]domain.Like, error,
) {
	v, err := impl.mapper.List(owner.Account(), LikeListDO{})
	if err != nil {
		if isErrorDataNotExists(err) {
			return nil, nil
		}

		return nil, convertError(err)
	}

	if len(v) == 0 {
		return nil, nil
	}

	r := make([]domain.Like, len(v))
	for i := range v {
		if err := v[i].toLike(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func (impl like) HasLike(owner domain.Account, obj *domain.ResourceObject) (bool, error) {
	do := toResourceObjectDO(obj)

	b, err := impl.mapper.HasLike(owner.Account(), &do)
	if err != nil {
		return false, convertError(err)
	}

	return b, nil
}

func (impl like) toLikeDO(v *domain.Like) LikeDO {
	return LikeDO{
		CreatedAt:        v.CreatedAt,
		ResourceObjectDO: toResourceObjectDO(&v.ResourceObject),
	}
}

type LikeListDO struct {
}

type LikeDO struct {
	CreatedAt int64

	ResourceObjectDO
}

func (do *LikeDO) toLike(r *domain.Like) (err error) {
	r.CreatedAt = do.CreatedAt

	return do.ResourceObjectDO.toResourceObject(&r.ResourceObject)
}
