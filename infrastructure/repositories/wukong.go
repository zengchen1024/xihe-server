package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type WuKongPictureListOptionDO = repository.WuKongPictureListOption

type WuKongMapper interface {
	ListSamples(string, []int) ([]string, error)
	ListPictures(string, *WuKongPictureListOptionDO) (WuKongPicturesDO, error)
}

func NewWuKongRepository(mapper WuKongMapper) repository.WuKong {
	return wukong{mapper}
}

type wukong struct {
	mapper WuKongMapper
}

func (impl wukong) ListSamples(sid string, nums []int) ([]string, error) {
	return impl.mapper.ListSamples(sid, nums)
}

func (impl wukong) ListPictures(sid string, opt *repository.WuKongPictureListOption) (
	repository.WuKongPictures, error,
) {
	v, err := impl.mapper.ListPictures(sid, opt)
	if err != nil {
		return repository.WuKongPictures{}, err
	}

	return v.toWuKongPictures()
}

type WuKongPicturesDO struct {
	Pictures []WuKongPictureInfoDO
	Total    int
}

func (do *WuKongPicturesDO) toWuKongPictures() (r repository.WuKongPictures, err error) {
	r.Total = do.Total

	r.Pictures = make([]domain.WuKongPictureInfo, len(do.Pictures))
	for i := range do.Pictures {
		r.Pictures[i], err = do.Pictures[i].toWuKongPictureInfo()
		if err != nil {
			return
		}
	}

	return
}

type WuKongPictureInfoDO struct {
	Link string

	WuKongPictureMetaDO
}

func (do *WuKongPictureInfoDO) toWuKongPictureInfo() (
	info domain.WuKongPictureInfo, err error,
) {
	info.Link = do.Link
	info.WuKongPictureMeta, err = do.toWuKongPictureMeta()

	return
}

type WuKongPictureMetaDO struct {
	Style string
	Desc  string
}

func (do *WuKongPictureMetaDO) toWuKongPictureMeta() (
	meta domain.WuKongPictureMeta, err error,
) {
	meta.Style = do.Style
	meta.Desc, err = domain.NewWuKongPictureDesc(do.Desc)

	return
}
