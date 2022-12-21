package repositories

import "github.com/opensourceways/xihe-server/domain/repository"

type WuKongPicturesDO = repository.WuKongPictures
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
	return impl.mapper.ListPictures(sid, opt)
}
