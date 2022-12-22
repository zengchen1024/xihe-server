package app

import (
	"errors"
	"io"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/bigmodel"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

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
	GenWuKongSamples(int) ([]string, error)
	WuKong(domain.Account, *WuKongCmd) (map[string]string, error)
	WuKongPictures(*WuKongPicturesListCmd) (WuKongPicturesDTO, error)
	AddLikeToWuKongPicture(cmd *WuKongPictureAddLikeCmd) (string, error)
	CancelLikeOnWuKongPicture(domain.Account, string) error
	ListLikedWuKongPictures(domain.Account) ([]UserLikedWuKongPictureDTO, error)
}

func NewBigModelService(
	fm bigmodel.BigModel,
	luojia repository.LuoJia,
	wukong repository.WuKong,
	wukongPicture repository.WuKongPicture,
) BigModelService {
	return bigModelService{
		fm:             fm,
		luojia:         luojia,
		wukong:         wukong,
		wukongPicture:  wukongPicture,
		wukongSampleId: fm.GetWuKongSampleId(),
	}
}

type bigModelService struct {
	fm bigmodel.BigModel

	luojia        repository.LuoJia
	wukong        repository.WuKong
	wukongPicture repository.WuKongPicture

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

func (s bigModelService) GenWuKongSamples(batchNum int) ([]string, error) {
	num := s.fm.GenWuKongSampleNums(batchNum)
	if len(num) == 0 {
		return nil, nil
	}

	return s.wukong.ListSamples(s.wukongSampleId, num)
}

func (s bigModelService) WuKong(
	user domain.Account, cmd *WuKongCmd,
) (links map[string]string, err error) {
	return s.fm.GenPicturesByWuKong(user, (*domain.WuKongPictureMeta)(cmd))
}

func (s bigModelService) WuKongPictures(cmd *WuKongPicturesListCmd) (
	dto WuKongPicturesDTO, err error,
) {
	return s.wukong.ListPictures(s.wukongSampleId, cmd)
}

func (s bigModelService) AddLikeToWuKongPicture(cmd *WuKongPictureAddLikeCmd) (
	code string, err error,
) {
	v, version, err := s.wukongPicture.List(cmd.User)
	if err != nil {
		return
	}
	if len(v) >= 10 { //TODO 10
		code = ErrorWuKongExccedMaxNum
		err = errors.New("exceed the max num user can add like to pictures")

		return
	}

	for i := range v {
		if v[i].OBSPath == cmd.OBSPath {
			code = ErrorWuKongDuplicateLike
			err = errors.New("the picture has been saved.")

			return
		}
	}

	meta, err := s.fm.ParseWuKongPictureMetaData(cmd.User, cmd.OBSPath)
	if err != nil {
		code = ErrorWuKongInvalidPath

		return
	}

	p, err := s.fm.MoveWuKongPictureToLikeDir(cmd.User, cmd.OBSPath)
	if err != nil {
		return
	}

	_, err = s.wukongPicture.Save(
		&domain.UserWuKongPicture{
			User: cmd.User,
			WuKongPicture: domain.WuKongPicture{
				OBSPath:           p,
				CreatedAt:         utils.Date(),
				WuKongPictureMeta: meta,
			},
		},
		version,
	)

	return
}

func (s bigModelService) CancelLikeOnWuKongPicture(user domain.Account, pid string) (err error) {
	v, err := s.wukongPicture.Get(user, pid)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			err = nil
		}

		return
	}

	if err = s.fm.DeleteWuKongPicture(v.OBSPath); err != nil {
		return
	}

	err = s.wukongPicture.Delete(user, pid)

	return
}

func (s bigModelService) ListLikedWuKongPictures(user domain.Account) (
	r []UserLikedWuKongPictureDTO, err error,
) {
	v, _, err := s.wukongPicture.List(user)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]UserLikedWuKongPictureDTO, len(v))
	for i := range v {
		item := &v[i]
		dto := &r[i]

		dto.Link, err = s.fm.GenWuKongPictureLink(item.OBSPath)
		if err != nil {
			return
		}

		dto.Id = item.Id
		dto.CreatedAt = item.CreatedAt
		dto.WuKongPictureMeta = item.WuKongPictureMeta
	}

	return
}
