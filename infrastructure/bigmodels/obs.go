package bigmodels

import (
	"io"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

func initOBS(cfg *OBSConfig) (s obsService, err error) {
	cli, err := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	if err != nil {
		return
	}

	s.cli = cli
	s.bucket = cfg.Bucket

	return
}

type obsService struct {
	cli *obs.ObsClient

	bucket string
}

func (s *obsService) createObject(f io.Reader, path string) error {
	input := &obs.PutObjectInput{}
	input.Bucket = s.bucket
	input.Key = path
	input.Body = f

	_, err := s.cli.PutObject(input)

	return err
}

func (s *service) UploadFile(f io.Reader, path string) error {
	return s.obs.createObject(f, path)
}
