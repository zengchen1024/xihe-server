package repositories

import (
	"github.com/opensourceways/xihe-server/domain/repository"
)

type WuKongMapper interface {
	ListSamples(string, []int) ([]string, error)
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
