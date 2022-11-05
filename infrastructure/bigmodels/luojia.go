package bigmodels

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
)

func (s *service) LuoJia(question string) (answer string, err error) {
	s.doIfFree(s.panguEndpoints, func(e string) error {
		answer, err = s.sendReqToLuojia(e, question)

		return err
	})

	return
}

func (s *service) sendReqToLuojia(endpoint, userName string) (answer string, err error) {
	t, err := s.token()
	if err != nil {
		return
	}

	body := []byte(fmt.Sprintf(`{"user_name":"%s"}`, userName))

	req, err := http.NewRequest(
		http.MethodPost, endpoint, bytes.NewBuffer(body),
	)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Auth-Token", t)

	var r struct {
		Msg    string `json:"msg"`
		Result string `json:"result"`
		Status int    `json:"status"`
	}

	if _, err = s.hc.ForwardTo(req, &r); err != nil {
		return
	}

	if r.Status != 200 {
		// TODO check
		err = errors.New(r.Msg)
	} else {
		answer = r.Result
	}

	return
}
