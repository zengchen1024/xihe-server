package gitlab

import (
	"path/filepath"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

func initOBS(config *Config) (s obsService, err error) {
	cfg := &config.OBS
	cli, err := obs.New(cfg.AccessKey, cfg.SecretKey, cfg.Endpoint)
	if err != nil {
		return
	}

	s.cli = cli
	s.bucket = cfg.Bucket
	s.lfsPath = config.LFSPath
	s.downloadExpiry = config.DownloadExpiry

	return
}

type obsService struct {
	cli *obs.ObsClient

	bucket         string
	lfsPath        string
	downloadExpiry int
}

func (s *obsService) GenObjectDownloadURL(p string) (string, error) {
	input := &obs.CreateSignedUrlInput{}
	input.Method = obs.HttpMethodGet
	input.Bucket = s.bucket
	input.Key = filepath.Join(s.lfsPath, p)
	input.Expires = s.downloadExpiry

	output, err := s.cli.CreateSignedUrl(input)
	if err != nil {
		return "", err
	}

	return output.SignedUrl, nil
}
