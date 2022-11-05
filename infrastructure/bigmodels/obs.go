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

	return
}

type obsService struct {
	cli *obs.ObsClient
}

func (s *obsService) createObject(f io.Reader, bucket, path string) error {
	input := &obs.PutObjectInput{}
	input.Bucket = bucket
	input.Key = path
	input.Body = f

	_, err := s.cli.PutObject(input)

	return err
}
