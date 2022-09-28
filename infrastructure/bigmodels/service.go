package bigmodels

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain/bigmodel"
)

type descOfPicture struct {
	Result struct {
		Instances struct {
			Image []string `json:"image"`
		} `json:"instances"`
	} `json:"inference_result"`
}

var fm *service

func Init(cfg *Config) {
	fm = &service{
		cfg:              *cfg,
		hc:               utils.HttpClient{MaxRetries: 3},
		singlePictures:   make(chan string, len(cfg.EndpointsOfSinglePicture)),
		multiplePictures: make(chan string, 1),
	}

	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	for _, e := range cfg.EndpointsOfSinglePicture {
		fm.singlePictures <- e
	}

	fm.multiplePictures <- cfg.EndpointOfMultiplePictures
}

func NewBigModelService() bigmodel.BigModel {
	return fm
}

type service struct {
	cfg Config

	hc utils.HttpClient

	singlePictures   chan string
	multiplePictures chan string
}

func (s *service) DescribePicture(picture io.Reader, contentType string) (string, error) {
	req, err := http.NewRequest(
		http.MethodPost, s.cfg.EndpointOfDescribingPicture, picture,
	)
	if err != nil {
		return "", err
	}

	t, err := s.token()
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", contentType)
	req.Header.Set("X-Auth-Token", t)

	desc := new(descOfPicture)

	if err = s.hc.ForwardTo(req, desc); err != nil {
		return "", err
	}

	if v := desc.Result.Instances.Image; len(v) > 0 {
		return v[0], nil
	}

	return "", errors.New("no content")
}

func (s *service) token() (string, error) {
	str := `
"auth":{
   "identity":{
      "methods":[
         "password"
      ],
      "password":{
         "user":{
            "name":"%v",
            "password":"%v",
            "domain":{
               "name":"%v"
            }
         }
      }
   },
   "scope":{
      "project":{
         "name":"%s"
      }
   }
}
	`

	cfg := &s.cfg
	body := fmt.Sprintf(str, cfg.User, cfg.Password, cfg.User, cfg.Project)

	resp, err := http.Post(
		s.cfg.AuthEndpoint, "application/json",
		strings.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	t := resp.Header.Get("x-subject-token")

	resp.Body.Close()

	return t, nil
}
