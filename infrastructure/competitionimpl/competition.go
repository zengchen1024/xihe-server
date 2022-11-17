package competitionimpl

import (
	"io"

	"github.com/opensourceways/xihe-server/domain/competition"
)

var cs *service

func Init(cfg *Config) error {
	obs, err := initOBS(&cfg.OBS)
	if err != nil {
		return err
	}

	cs = &service{
		obs: obs,
	}

	return nil
}

func NewCompetitionService() competition.Competition {
	return cs
}

type service struct {
	obs obsService
}

func (s *service) UploadSubmissionFile(data io.Reader, path string) error {
	return s.obs.createObject(data, path)
}
