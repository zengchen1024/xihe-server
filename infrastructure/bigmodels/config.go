package bigmodels

import (
	"errors"
	"net/url"
	"strings"
)

type Config struct {
	OBS          OBSConfig `json:"obs"             required:"true"`
	User         string    `json:"user"            required:"true"`
	Password     string    `json:"password"        required:"true"`
	Project      string    `json:"project"         required:"true"`
	AuthEndpoint string    `json:"auth_endpoint"   required:"true"`

	MaxPictureSizeToDescribe int64 `json:"max_picture_size_to_describe"`
	MaxPictureSizeToVQA      int64 `json:"max_picture_size_to_vqa"`

	EndpointsOfSinglePicture    string `json:"endpoints_of_signle_picture"    required:"true"`
	EndpointOfDescribingPicture string `json:"endpoint_of_describing_picture" required:"true"`
	EndpointOfMultiplePictures  string `json:"endpoint_of_multiple_pictures" required:"true"`
	EndpointOfVQA               string `json:"endpoint_of_vqa" required:"true"`

	endpointsOfSinglePicture []string
}

func (cfg *Config) SetDefault() {
	if cfg.MaxPictureSizeToDescribe <= 0 {
		cfg.MaxPictureSizeToDescribe = 200 << 10
	}

	if cfg.MaxPictureSizeToVQA <= 0 {
		cfg.MaxPictureSizeToVQA = 200 << 10
	}
}

func (cfg *Config) Validate() error {
	if _, err := url.Parse(cfg.EndpointOfDescribingPicture); err != nil {
		return errors.New("invalid url for describing picture")
	}

	if _, err := url.Parse(cfg.EndpointOfMultiplePictures); err != nil {
		return errors.New("invalid url for generating multiple pictures")
	}

	if _, err := url.Parse(cfg.EndpointOfVQA); err != nil {
		return errors.New("invalid url for vqa")
	}

	v := strings.Split(
		strings.Trim(cfg.EndpointsOfSinglePicture, ","), ",",
	)

	for _, i := range v {
		if _, err := url.Parse(i); err != nil {
			return errors.New("invalid url for generating single picture")
		}
	}

	if len(v) == 0 {
		return errors.New("missing endpoints for generating single picture")
	}

	cfg.endpointsOfSinglePicture = v

	return nil
}

type OBSConfig struct {
	Endpoint  string `json:"Endpoint"   required:"true"`
	AccessKey string `json:"access_key" required:"true"`
	SecretKey string `json:"secret_key" required:"true"`
	Bucket    string `json:"bucket"     required:"true"`
}
