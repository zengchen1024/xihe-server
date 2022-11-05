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
	Ask(domain.Question, string) (string, error)
	UploadFile(f io.Reader, path string) error
	PanGu(string) (string, error)
	LuoJia(string) (string, error)
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

func (s bigModelService) Ask(q domain.Question, f string) (string, error) {
	// TODO check the content of question to see if it is legal

	return s.fm.Ask(q, f)
}

func (s bigModelService) UploadFile(f io.Reader, path string) error {
	return s.fm.UploadFile(f, path)
}

func (s bigModelService) PanGu(q string) (string, error) {
	// TODO check the content of question to see if it is legal

	return s.fm.PanGu(q)
}

func (s bigModelService) LuoJia(q string) (string, error) {
	// TODO check the content of question to see if it is legal

	return s.fm.LuoJia(q)
}
