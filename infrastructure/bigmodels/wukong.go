package bigmodels

import (
	"bytes"
	"errors"
	"net/http"
	"strings"

	libutils "github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type wukongInfo struct {
	cfg       WuKong
	maxBatch  int
	endpoints chan string
}

func newWuKongInfo(cfg *Config) wukongInfo {
	v := &cfg.WuKong
	info := wukongInfo{
		cfg:      *v,
		maxBatch: utils.LCM(v.SampleCount, v.SampleNum) / v.SampleNum,
	}

	ce := &cfg.Endpoints
	es, _ := ce.parse(ce.WuKong)

	info.endpoints = make(chan string, len(es))
	for _, e := range es {
		info.endpoints <- e
	}

	return info
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

func (s *service) GenPicturesByWuKong(
	user domain.Account, desc []string,
) (r []string, err error) {
	s.doIfFree(s.wukongInfo.endpoints, func(e string) error {
		r, err = s.genPicturesByWuKong(e, user, desc)

		return err
	})

	return
}

func (s *service) genPicturesByWuKong(
	endpoint string, user domain.Account, desc []string,
) ([]string, error) {
	t, err := genToken(&s.wukongInfo.cfg.CloudConfig)
	if err != nil {
		return nil, err
	}

	opt := wukongRequest{
		Input: strings.Join(desc, " "),
		User:  user.Account(),
	}
	body, err := libutils.JsonMarshal(&opt)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	var r wukongResponse
	if _, err = s.hc.ForwardTo(req, &r); err != nil {
		return nil, err
	}

	if r.Status == 200 {
		return r.Output, nil
	}

	return nil, errors.New(r.Msg)
}

type wukongRequest struct {
	Input string `json:"input_text" binding:"required"`
	User  string `json:"user_name"  binding:"required"`
}

type wukongResponse struct {
	Status int      `json:"status"`
	Output []string `json:"output_image_url"`
	Msg    string   `json:"msg"`
}
