package app

import (
	"errors"
	"io"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type CodeGeexDTO = bigmodel.CodeGeexResp

type CodeGeexCmd bigmodel.CodeGeexReq

func (cmd *CodeGeexCmd) Validate() error {
	if cmd.Content == "" || cmd.Lang == "" {
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
	GenPicture(domain.Account, string) (string, string, error)
	GenPictures(domain.Account, string) ([]string, string, error)
	Ask(domain.Question, string) (string, string, error)
	VQAUploadPicture(io.Reader, domain.Account, string) error
	LuoJiaUploadPicture(io.Reader, domain.Account) error
	PanGu(string) (string, string, error)
	CodeGeex(*CodeGeexCmd) (CodeGeexDTO, string, error)
	LuoJia(domain.Account) (string, error)
	ListLuoJiaRecord(domain.Account) ([]LuoJiaRecordDTO, error)
	GenWuKongSamples() ([]string, error)
}

func NewBigModelService(
	fm bigmodel.BigModel,
	luojia repository.LuoJia,
	wukong repository.WuKong,
) BigModelService {
	return bigModelService{
		fm:             fm,
		luojia:         luojia,
		wukong:         wukong,
		wukongSampleId: fm.GetWuKongSampleId(),
	}
}

type bigModelService struct {
	fm bigmodel.BigModel

	luojia repository.LuoJia
	wukong repository.WuKong

	wukongSampleId string
}

func (s bigModelService) DescribePicture(
	picture io.Reader, name string, length int64,
) (string, error) {
	return s.fm.DescribePicture(picture, name, length)
}

func (s bigModelService) GenPicture(
	user domain.Account, desc string,
) (link string, code string, err error) {
	if link, err = s.fm.GenPicture(user, desc); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) GenPictures(
	user domain.Account, desc string,
) (links []string, code string, err error) {
	if links, err = s.fm.GenPictures(user, desc); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) Ask(q domain.Question, f string) (v string, code string, err error) {
	if v, err = s.fm.Ask(q, f); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) VQAUploadPicture(f io.Reader, user domain.Account, fileName string) error {
	return s.fm.VQAUploadPicture(f, user, fileName)
}

func (s bigModelService) LuoJiaUploadPicture(f io.Reader, user domain.Account) error {
	return s.fm.LuoJiaUploadPicture(f, user)
}

func (s bigModelService) PanGu(q string) (v string, code string, err error) {
	if v, err = s.fm.PanGu(q); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) LuoJia(user domain.Account) (v string, err error) {
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

func (s bigModelService) CodeGeex(cmd *CodeGeexCmd) (dto CodeGeexDTO, code string, err error) {
	if dto, err = s.fm.CodeGeex((*bigmodel.CodeGeexReq)(cmd)); err != nil {
		code = s.setCode(err)
	}

	return
}

func (s bigModelService) setCode(err error) string {
	if err != nil && bigmodel.IsErrorSensitiveInfo(err) {
		return ErrorBigModelSensitiveInfo
	}

	return ""
}

func (s bigModelService) GenWuKongSamples() ([]string, error) {
	num := s.fm.GenWuKongSampleNums()
	if len(num) == 0 {
		return nil, nil
	}

	return s.wukong.ListSamples(s.wukongSampleId, num)
}
