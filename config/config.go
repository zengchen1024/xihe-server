package config

import (
	"github.com/huaweicloud/golangsdk"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

func LoadConfig(path string) (*Config, error) {
	v := new(Config)

	if err := utils.LoadFromYaml(path, v); err != nil {
		return nil, err
	}

	v.setDefault()

	if err := v.validate(); err != nil {
		return nil, err
	}

	return v, nil
}

type Config struct {
	DefaultPassword string `json:"default_password" required:"true"`
	EncryptionKey   string `json:"encryption_key" required:"true"`

	Authing AuthingService `json:"authing_service" required:"true"`
	Mongodb Mongodb        `json:"mongodb" required:"true"`
	Gitlab  Gitlab         `json:"gitlab" required:"true"`
	API     API            `json:"api" required:"true"`
}

func (cfg *Config) setDefault() {
}

func (cfg *Config) validate() error {
	if _, err := golangsdk.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	_, err := domain.NewPassword(cfg.DefaultPassword)

	return err
}

type Mongodb struct {
	MongodbConn       string `json:"mongodb_conn" required:"true"`
	DBName            string `json:"mongodb_db" required:"true"`
	ProjectCollection string `json:"project_collection" required:"true"`
	ModelCollection   string `json:"model_collection" required:"true"`
	DatasetCollection string `json:"dataset_collection" required:"true"`
	UserCollection    string `json:"user_collection" required:"true"`
	LoginCollection   string `json:"login_collection" required:"true"`
}

type AuthingService struct {
	APPId  string `json:"app_id" required:"true"`
	Secret string `json:"secret" required:"true"`
}

type Gitlab struct {
	Endpoint  string `json:"endpoint" required:"true"`
	RootToken string `json:"root_token" required:"true"`
}

type API struct {
	APITokenExpiry int64  `json:"api_token_expiry" required:"true"`
	APITokenKey    string `json:"api_token_key" required:"true"`
}
