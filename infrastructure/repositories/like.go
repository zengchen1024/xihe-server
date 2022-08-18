package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type LikeMapper interface {
	Insert(string, LikeDO) error
	Delete(string, LikeDO) error
	List(string, LikeListDO) ([]LikeDO, error)
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

func (impl like) toLikeDO(v *domain.Like) LikeDO {
	return LikeDO{
		ResourceOwner: v.ResourceOwner.Account(),
		ResourceType:  v.ResourceType.ResourceType(),
		ResourceId:    v.ResourceId,
	}
}

type LikeListDO struct {
}

type LikeDO struct {
	ResourceOwner string
	ResourceType  string
	ResourceId    string
}

func (do *LikeDO) toLike(r *domain.Like) (err error) {
	if r.ResourceOwner, err = domain.NewAccount(do.ResourceOwner); err != nil {
		return
	}

	if r.ResourceType, err = domain.NewResourceType(do.ResourceType); err != nil {
		return
	}

	r.ResourceId = do.ResourceId

	return
}
