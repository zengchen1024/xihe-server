package bigmodels

import (
	"bytes"
	"net/http"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain/bigmodel"
)

type codegeexInfo struct {
	endpoints chan string
}

func newCodeGeexInfo(cfg *Config) codegeexInfo {
	ce := &cfg.Endpoints

	es, _ := ce.parse(ce.CodeGeex)

	v := codegeexInfo{
		endpoints: make(chan string, len(es)),
	}

	for _, e := range es {
		v.endpoints <- e
	}

	return v
}

func (s *service) CodeGeex(question *bigmodel.CodeGeexReq) (answer string, err error) {
	s.doIfFree(s.codegeexInfo.endpoints, func(e string) error {
		answer, err = s.sendReqToCodeGeex(e, question)

		return err
	})

	return
}

func (s *service) sendReqToCodeGeex(
	endpoint string, question *bigmodel.CodeGeexReq,
) (answer string, err error) {
	opt := codegeexReq{
		Samples: question.Content,
		Lang:    question.Lang,
	}

	body, err := utils.JsonMarshal(&opt)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	t, err := s.token()
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	r := new(codegeexResp)
	if _, err = s.hc.ForwardTo(req, r); err == nil {
		answer = r.Result
	}

	return
}

type codegeexReq struct {
	Samples string `json:"samples"`
	Lang    string `json:"language"`
}

type codegeexResp struct {
	Result string `json:"result"`
}
