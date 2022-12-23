package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type WuKongPictureMapper interface {
	List(string) ([]WuKongPictureDO, int, error)
	Insert(string, *WuKongPictureDO, int) (string, error)
	Delete(string, string) error
	Get(string, string) (WuKongPictureDO, error)
}

func NewWuKongPictureRepository(mapper WuKongPictureMapper) repository.WuKongPicture {
	return wukongPicture{mapper}
}

type wukongPicture struct {
	mapper WuKongPictureMapper
}

func (impl wukongPicture) List(user domain.Account) (
	[]domain.WuKongPicture, int, error,
) {
	v, version, err := impl.mapper.List(user.Account())
	if err != nil {
		return nil, 0, err
	}

	r := make([]domain.WuKongPicture, len(v))
	for i := range v {
		if r[i], err = v[i].toWuKongPicture(); err != nil {
			return nil, 0, err
		}
	}

	return r, version, nil
}

func (impl wukongPicture) Save(p *domain.UserWuKongPicture, version int) (string, error) {
	if p.Id != "" {
		return "", errors.New("must be a new picture")
	}

	do := WuKongPictureDO{
		OBSPath:   p.OBSPath,
		CreatedAt: p.CreatedAt,
	}
	do.Style = p.Style
	do.Desc = p.Desc.WuKongPictureDesc()

	v, err := impl.mapper.Insert(p.User.Account(), &do, version)
	if err != nil {
		return "", convertError(err)
	}

	return v, nil
}

func (impl wukongPicture) Delete(user domain.Account, pid string) error {
	if err := impl.mapper.Delete(user.Account(), pid); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl wukongPicture) Get(user domain.Account, pid string) (domain.WuKongPicture, error) {
	v, err := impl.mapper.Get(user.Account(), pid)
	if err != nil {
		return domain.WuKongPicture{}, convertError(err)
	}

	return v.toWuKongPicture()
}

type WuKongPictureDO struct {
	Id        string
	OBSPath   string
	CreatedAt string

	WuKongPictureMetaDO
}

func (do *WuKongPictureDO) toWuKongPicture() (r domain.WuKongPicture, err error) {
	r.Id = do.Id
	r.OBSPath = do.OBSPath
	r.CreatedAt = do.CreatedAt

	r.WuKongPictureMeta, err = do.toWuKongPictureMeta()

	return
}
