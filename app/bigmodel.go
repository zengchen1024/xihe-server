package app

import (
	"io"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/bigmodel"
)

type BigModelService interface {
	DescribePicture(io.Reader, string, int64) (string, error)
	GenPicture(domain.Account, string) (string, error)
	GenPictures(domain.Account, string) ([]string, error)
}

func NewBigModelService(fm bigmodel.BigModel) BigModelService {
	return bigModelService{fm}
}

type bigModelService struct {
	fm bigmodel.BigModel
}

func (s bigModelService) DescribePicture(
	picture io.Reader, name string, length int64,
) (string, error) {
	return s.fm.DescribePicture(picture, name, length)
}

func (s bigModelService) GenPicture(
	user domain.Account, desc string,
) (string, error) {
	return s.fm.GenPicture(user, desc)
}

func (s bigModelService) GenPictures(
	user domain.Account, desc string,
) ([]string, error) {
	return s.fm.GenPictures(user, desc)
}
