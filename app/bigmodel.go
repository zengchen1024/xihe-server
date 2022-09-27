package app

import (
	"io"

	"github.com/opensourceways/xihe-server/domain/bigmodel"
)

type BigModelService interface {
	DescribePicture(io.Reader, string) (string, error)
}

func NewBigModelService(fm bigmodel.BigModel) BigModelService {
	return bigModelService{fm}
}

type bigModelService struct {
	fm bigmodel.BigModel
}

func (s bigModelService) DescribePicture(
	picture io.Reader, contentType string,
) (string, error) {
	return s.fm.DescribePicture(picture, contentType)
}
