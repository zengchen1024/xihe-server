package bigmodels

import (
	"errors"
	"net/url"
	"strings"
)

type Config struct {
	User         string `json:"user"            required:"true"`
	Password     string `json:"password"        required:"true"`
	Project      string `json:"project"         required:"true"`
	AuthEndpoint string `json:"auth_endpoint"   required:"true"`

	EndpointOfDescribingPicture string `json:"endpoint_of_describing_picture" required:"true"`
	EndpointsOfSinglePicture    string `json:"endpoints_of_signle_picture"    required:"true"`
	EndpointOfMultiplePictures  string `json:"endpoints_of_multiple_pictures" required:"true"`

	endpointsOfSinglePicture []string
}

func (cfg *Config) Validate() error {
	if _, err := url.Parse(cfg.EndpointOfDescribingPicture); err != nil {
		return errors.New("invalid url for describing picture")
	}

	if _, err := url.Parse(cfg.EndpointOfMultiplePictures); err != nil {
		return errors.New("invalid url for generating multiple pictures")
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
