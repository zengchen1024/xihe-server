package bigmodels

import "github.com/opensourceways/xihe-server/utils"

type wukongInfo struct {
	cfg WuKong
}

func (s *service) GenWuKongSampleNums() []int {
	cfg := &s.wukongInfo.cfg

	return utils.GenRandoms(cfg.SampleCount, cfg.SampleNum)
}

func (s *service) GetWuKongSampleId() string {
	return s.wukongInfo.cfg.SampleId
}
