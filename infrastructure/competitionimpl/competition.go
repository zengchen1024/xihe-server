package competitionimpl

import "io"

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

func NewCompetitionService() *service {
	return cs
}

type service struct {
	obs obsService
}

func (s *service) UploadSubmissionFile(data io.Reader, path string) error {
	return s.obs.createObject(data, path)
}

func (s *service) Upload(data io.Reader, path string) error {
	return s.obs.createObject(data, path)
}
