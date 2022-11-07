package bigmodels

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain/bigmodel"
)

type codegeexInfo struct {
	endpoints chan string
	ak        string
	sk        string
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

	v.ak = cfg.CodeGeex.AK
	v.sk = cfg.CodeGeex.SK

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
		Prompt:    question.Content,
		N:         question.ResultNum,
		Lang:      question.Lang,
		Apikey:    s.codegeexInfo.ak,
		Apisecret: s.codegeexInfo.sk,
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

	req.Header.Set("Content-Type", "application/json")

	r := new(codegeexResp)
	if _, err = s.hc.ForwardTo(req, r); err != nil {
		return
	}

	if r.Status != 0 {
		err = errors.New("failed")

		return
	}

	if len(r.Result.OutPut.Code) > 0 {
		answer = r.Result.OutPut.Code[0]

		return
	}

	if question.Lang == "Python" {
		answer = `\n# Code generation finished, modify this comment to continue the generation.`
	} else {
		answer = `\n// Code generation finished, modify this comment to continue the generation.`
	}

	return
}

type codegeexReq struct {
	Prompt    string `json:"prompt"`
	N         int    `json:"n"`
	Lang      string `json:"lang"`
	Apikey    string `json:"apikey"`
	Apisecret string `json:"apisecret"`
}

type codegeexResp struct {
	Status int `json:"status"`
	Result struct {
		OutPut struct {
			Code []string `json:"code"`
		} `json:"output"`
	} `json:"result"`
}
