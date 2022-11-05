package bigmodels

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type panguInfo struct {
	endpoints chan string
}

func newPanGuInfo(cfg *Config) panguInfo {
	v := panguInfo{
		endpoints: make(chan string, 1),
	}

	v.endpoints <- cfg.EndpointsOfPangu

	return v
}

func (s *service) PanGu(question string) (answer string, err error) {
	s.doIfFree(s.panguInfo.endpoints, func(e string) error {
		answer, err = s.sendReqToPangu(e, question)

		return err
	})

	return
}

func (s *service) sendReqToPangu(endpoint, question string) (answer string, err error) {
	t, err := s.token()
	if err != nil {
		return
	}

	body := []byte(fmt.Sprintf(`{"question":"%s"}`, question))

	req, err := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	var r struct {
		Result string `json:"result"`
	}

	if _, err = s.hc.ForwardTo(req, &r); err != nil {
		return
	}

	if v := strings.Split(r.Result, "\n"); len(v) == 2 && v[0] == question {
		answer = v[1]
	}

	// TODO: failed msg

	return
}
