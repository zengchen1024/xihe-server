package app

import (
	"errors"
	"io"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type CodeGeexCmd bigmodel.CodeGeexReq

func (cmd *CodeGeexCmd) Validate() error {
	b := cmd.Content != "" &&
		cmd.Lang != "" &&
		cmd.ResultNum > 0

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

type LuoJiaRecordDTO struct {
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
}

type BigModelService interface {
	DescribePicture(io.Reader, string, int64) (string, error)
	GenPicture(domain.Account, string) (string, error)
	GenPictures(domain.Account, string) ([]string, error)
	Ask(domain.Question, string) (string, error)
	VQAUploadPicture(io.Reader, domain.Account, string) error
	LuoJiaUploadPicture(io.Reader, domain.Account) error
	PanGu(string) (string, error)
	CodeGeex(*CodeGeexCmd) (string, error)
	LuoJia(domain.Account) (string, error)
	ListLuoJiaRecord(domain.Account) ([]LuoJiaRecordDTO, error)
}

func NewBigModelService(
	fm bigmodel.BigModel,
	luojia repository.LuoJia,
) BigModelService {
	return bigModelService{
		fm:     fm,
		luojia: luojia,
	}
}

type bigModelService struct {
	fm bigmodel.BigModel

	luojia repository.LuoJia
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

func (s bigModelService) VQAUploadPicture(f io.Reader, user domain.Account, fileName string) error {
	return s.fm.VQAUploadPicture(f, user, fileName)
}

func (s bigModelService) LuoJiaUploadPicture(f io.Reader, user domain.Account) error {
	return s.fm.LuoJiaUploadPicture(f, user)
}

func (s bigModelService) PanGu(q string) (string, error) {
	// TODO check the content of question to see if it is legal

	return s.fm.PanGu(q)
}

func (s bigModelService) LuoJia(user domain.Account) (v string, err error) {
	// TODO check the content of question to see if it is legal

	if v, err = s.fm.LuoJia(user.Account()); err != nil {
		return
	}

	record := domain.UserLuoJiaRecord{User: user}
	record.CreatedAt = utils.Now()

	s.luojia.Save(&record)

	return
}

func (s bigModelService) ListLuoJiaRecord(user domain.Account) (
	dtos []LuoJiaRecordDTO, err error,
) {
	v, err := s.luojia.List(user)
	if err != nil || len(v) == 0 {
		return
	}

	dtos = append(dtos, LuoJiaRecordDTO{
		CreatedAt: utils.ToDate(v[0].CreatedAt),
		Id:        v[0].Id,
	})

	return
}

func (s bigModelService) CodeGeex(cmd *CodeGeexCmd) (string, error) {
	return s.fm.CodeGeex((*bigmodel.CodeGeexReq)(cmd))
}
