package bigmodels

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/domain/bigmodel"
)

var fm *service

func Init(cfg *Config) error {
	obs, err := initOBS(&cfg.OBS)
	if err != nil {
		return err
	}

	fm = &service{
		obs: obs,
		cfg: cfg.Cloud,
		hc:  utils.NewHttpClient(3),
	}

	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	fm.vqaInfo = newVQAInfo(cfg)
	fm.panguInfo = newPanGuInfo(cfg)
	fm.luojiaInfo = newLuoJiaInfo(cfg)
	fm.pictureGenInfo = newPictureGenInfo(cfg)
	fm.pictureDescInfo = newPictureDescInfo(cfg)

	return nil
}

func NewBigModelService() bigmodel.BigModel {
	return fm
}

type service struct {
	cfg CloudConfig
	obs obsService

	hc utils.HttpClient

	vqaInfo         vqaInfo
	panguInfo       panguInfo
	luojiaInfo      luojiaInfo
	pictureGenInfo  pictureGenInfo
	pictureDescInfo pictureDescInfo
}

func (s *service) token() (string, error) {
	str := `
{
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
}
	`

	cfg := &s.cfg

	body := fmt.Sprintf(
		str, cfg.User, cfg.Password, cfg.User, cfg.Project,
	)

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

func (s *service) doIfFree(
	ec chan string,
	f func(string) error,
) error {
	select {
	case e := <-ec:
		err := f(e)
		ec <- e

		return err

	default:
		return errors.New("busy")
	}
}
