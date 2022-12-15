package bigmodels

import "github.com/opensourceways/xihe-server/utils"

type wukongInfo struct {
	cfg      WuKong
	maxBatch int
}

func newWuKongInfo(cfg *Config) wukongInfo {
	v := &cfg.WuKong

	return wukongInfo{
		cfg:      *v,
		maxBatch: utils.LCM(v.SampleCount, v.SampleNum) / v.SampleNum,
	}
}

func (s *service) GenWuKongSampleNums(batchNum int) []int {
	cfg := &s.wukongInfo.cfg
	num := cfg.SampleNum
	count := cfg.SampleCount

	i := ((batchNum % s.wukongInfo.maxBatch) * num) % count

	r := make([]int, num)
	for j := 0; j < num; j++ {
		v := i + j
		if v >= count {
			v -= count
		}
		r[j] = v + 1
	}

	return r
}

func (s *service) GetWuKongSampleId() string {
	return s.wukongInfo.cfg.SampleId
}
