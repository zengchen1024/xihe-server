package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type WuKongPictureDO = domain.WuKongPicture

type WuKongPictureMapper interface {
	List(string) ([]WuKongPictureDO, int, error)
	Insert(string, WuKongPictureDO, int) (string, error)
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
	return impl.mapper.List(user.Account())
}

func (impl wukongPicture) Save(p *domain.UserWuKongPicture, version int) (string, error) {
	if p.Id != "" {
		return "", errors.New("must be a new project")
	}

	v, err := impl.mapper.Insert(p.User.Account(), p.WuKongPicture, version)
	if err != nil {
		return "", convertError(err)
	}

	return v, nil
}

func (impl wukongPicture) Delete(user domain.Account, pid string) error {
	err := impl.mapper.Delete(user.Account(), pid)
	if err != nil {
		return convertError(err)
	}

	return nil
}

func (impl wukongPicture) Get(user domain.Account, pid string) (domain.WuKongPicture, error) {
	v, err := impl.mapper.Get(user.Account(), pid)
	if err != nil {
		err = convertError(err)
	}

	return v, err
}
