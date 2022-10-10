package config

import (
	"errors"
	"regexp"
	"strings"

	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/controller"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/bigmodels"
	"github.com/opensourceways/xihe-server/infrastructure/messages"
)

var reIpPort = regexp.MustCompile(`^((25[0-5]|(2[0-4]|1\d|[1-9]|)\d)\.?\b){4}:[1-9][0-9]*$`)

func LoadConfig(path string, cfg interface{}) error {
	if err := utils.LoadFromYaml(path, cfg); err != nil {
		return err
	}

	if f, ok := cfg.(ConfigSetDefault); ok {
		f.SetDefault()
	}

	if f, ok := cfg.(ConfigValidate); ok {
		if err := f.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type ConfigValidate interface {
	Validate() error
}

type ConfigSetDefault interface {
	SetDefault()
}

type Config struct {
	MaxRetry        int `json:"max_retry"`
	ActivityKeepNum int `json:"activity_keep_num"`

	Authing  AuthingService       `json:"authing_service" required:"true"`
	BigModel bigmodels.Config     `json:"bigmodel"        required:"true"`
	Mongodb  Mongodb              `json:"mongodb"         required:"true"`
	Gitlab   Gitlab               `json:"gitlab"          required:"true"`
	Domain   domain.Config        `json:"domain"          required:"true"`
	API      controller.APIConfig `json:"api"             required:"true"`
	MQ       MQ                   `json:"mq"              required:"true"`
}

func (cfg *Config) GetMQConfig() mq.MQConfig {
	return mq.MQConfig{
		Addresses: cfg.MQ.ParseAddress(),
	}
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Authing,
		&cfg.Domain,
		&cfg.Mongodb,
		&cfg.Gitlab,
		&cfg.API,
		&cfg.MQ,
		&cfg.BigModel,
	}
}

func (cfg *Config) SetDefault() {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}

	if cfg.ActivityKeepNum <= 0 {
		cfg.ActivityKeepNum = 25
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(ConfigSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *Config) Validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(ConfigValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

type Mongodb struct {
	MongodbConn        string `json:"mongodb_conn" required:"true"`
	DBName             string `json:"mongodb_db" required:"true"`
	ProjectCollection  string `json:"project_collection" required:"true"`
	ModelCollection    string `json:"model_collection" required:"true"`
	DatasetCollection  string `json:"dataset_collection" required:"true"`
	UserCollection     string `json:"user_collection" required:"true"`
	LoginCollection    string `json:"login_collection" required:"true"`
	LikeCollection     string `json:"like_collection" required:"true"`
	ActivityCollection string `json:"activity_collection" required:"true"`
	TagCollection      string `json:"tag_collection" required:"true"`
}

type AuthingService struct {
	APPId  string `json:"app_id" required:"true"`
	Secret string `json:"secret" required:"true"`
}

type Gitlab struct {
	Endpoint  string `json:"endpoint" required:"true"`
	RootToken string `json:"root_token" required:"true"`
}

type MQ struct {
	Address string          `json:"address" required:"true"`
	Topics  messages.Topics `json:"topics"  required:"true"`
}

func (cfg *MQ) Validate() error {
	if r := cfg.ParseAddress(); len(r) == 0 {
		return errors.New("invalid mq address")
	}

	return nil
}

func (cfg *MQ) ParseAddress() []string {
	v := strings.Split(cfg.Address, ",")
	r := make([]string, 0, len(v))
	for i := range v {
		if reIpPort.MatchString(v[i]) {
			r = append(r, v[i])
		}
	}

	return r
}

func (cfg *Config) InitDomainConfig() {
	domain.Init(&cfg.Domain)
}
