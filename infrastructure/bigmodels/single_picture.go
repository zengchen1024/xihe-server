package bigmodels

import (
	"bytes"
	"errors"
	"net/http"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain"
)

func (s *service) GenPicture(user domain.Account, desc string) (string, error) {
	select {
	case e := <-s.singlePictures:
		v, err := s.genSinglePicture(user, e, desc)
		s.singlePictures <- e

		return v, err

	default:
		return "", errors.New("busy")
	}
}

type singlePictureCreateOpt struct {
	Desc string `json:"input_text"`
	User string `json:"user_name"`
}

type singlePicture struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
	Output string `json:"output_image_url"`
}

func (p *singlePicture) picture() (string, error) {
	if p.Status == -1 {
		return "", errors.New(p.Msg)
	}

	return p.Output, nil
}

func (s *service) genSinglePicture(
	user domain.Account, endpoint, desc string,
) (v string, err error) {
	t, err := s.token()
	if err != nil {
		return
	}

	opt := singlePictureCreateOpt{
		Desc: desc,
		User: user.Account(),
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
	req.Header.Set("X-Auth-Token", t)

	r := new(singlePicture)

	if err = s.hc.ForwardTo(req, r); err != nil {
		return "", err
	}

	return r.picture()
}
