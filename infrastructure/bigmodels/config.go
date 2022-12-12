package bigmodels

import (
	"errors"
	"net/url"
	"strings"
)

type Config struct {
	OBS        OBSConfig   `json:"obs"             required:"true"`
	Cloud      CloudConfig `json:"cloud"           required:"true"`
	WuKong     WuKong      `json:"wukong"          required:"true"`
	Endpoints  Endpoints   `json:"endpoints"       required:"true"`
	Moderation Moderation  `json:"moderation"      required:"true"`

	MaxPictureSizeToDescribe int64 `json:"max_picture_size_to_describe"`
	MaxPictureSizeToVQA      int64 `json:"max_picture_size_to_vqa"`
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
	return cfg.Endpoints.validate()
}

type OBSConfig struct {
	Endpoint     string `json:"endpoint"           required:"true"`
	AccessKey    string `json:"access_key"         required:"true"`
	SecretKey    string `json:"secret_key"         required:"true"`
	VQABucket    string `json:"vqa_bucket"         required:"true"`
	LuoJiaBucket string `json:"luo_jia_bucket"     required:"true"`
}

type CloudConfig struct {
	User         string `json:"user"            required:"true"`
	Password     string `json:"password"        required:"true"`
	Project      string `json:"project"         required:"true"`
	AuthEndpoint string `json:"auth_endpoint"   required:"true"`
}

type Endpoints struct {
	VQA              string `json:"vqa"                required:"true"`
	Pangu            string `json:"pangu"              required:"true"`
	LuoJia           string `json:"luojia"             required:"true"`
	CodeGeex         string `json:"codegeex"           required:"true"`
	DescPicture      string `json:"desc_picture"       required:"true"`
	SinglePicture    string `json:"signle_picture"     required:"true"`
	MultiplePictures string `json:"multiple_pictures"  required:"true"`
}

func (e *Endpoints) validate() (err error) {
	if _, err = e.parse(e.VQA); err != nil {
		return
	}

	if _, err = e.parse(e.Pangu); err != nil {
		return
	}

	if _, err = e.parse(e.LuoJia); err != nil {
		return
	}

	if _, err = e.parse(e.DescPicture); err != nil {
		return
	}

	if _, err = e.parse(e.SinglePicture); err != nil {
		return
	}

	_, err = e.parse(e.MultiplePictures)

	return
}

func (e *Endpoints) parse(s string) ([]string, error) {
	v := strings.Split(strings.Trim(s, ","), ",")

	for _, i := range v {
		if _, err := url.Parse(i); err != nil {
			return nil, errors.New("invalid url")
		}
	}

	if len(v) == 0 {
		return nil, errors.New("missing endpoints")
	}

	return v, nil
}

type Moderation struct {
	Endpoint   string `json:"endpoint"       required:"true"`
	AccessKey  string `json:"access_key"     required:"true"`
	SecretKey  string `json:"secret_key"     required:"true"`
	IAMEndpint string `json:"iam_endpoint"   required:"true"`
	Region     string `json:"region"         required:"true"`
}

type WuKong struct {
	SampleId    string `json:"sample_id"     required:"true"`
	SampleNum   int    `json:"sample_num"    required:"true"`
	SampleCount int    `json:"sample_count"  required:"true"`
}
