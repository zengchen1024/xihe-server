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

	Authing  AuthingService `json:"authing_service" required:"true"`
	Resource Resource       `json:"resource" required:"true"`
	Mongodb  Mongodb        `json:"mongodb" required:"true"`
	Gitlab   Gitlab         `json:"gitlab" required:"true"`
	API      API            `json:"api" required:"true"`
	User     User           `json:"user"`
}

func (cfg *Config) setDefault() {
	cfg.Resource.setdefault()
	cfg.User.setDefault()
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
	APPId    string `json:"app_id" required:"true"`
	Secret   string `json:"secret" required:"true"`
	Endpoint string `json:"endpoint" required:"true"`
}

type Gitlab struct {
	Endpoint  string `json:"endpoint" required:"true"`
	RootToken string `json:"root_token" required:"true"`
}

type API struct {
	APITokenExpiry int64  `json:"api_token_expiry" required:"true"`
	APITokenKey    string `json:"api_token_key" required:"true"`
}

type Resource struct {
	MaxNameLength int `json:"max_name_length"`
	MinNameLength int `json:"min_name_length"`
	MaxDescLength int `json:"max_desc_length"`

	Covers           []string `json:"covers" required:"true"`
	Protocols        []string `json:"protocols" required:"true"`
	ProjectType      []string `json:"project_type" required:"true"`
	TrainingPlatform []string `json:"training_platform" required:"true"`
}

func (r *Resource) setdefault() {
	if r.MaxNameLength == 0 {
		r.MaxNameLength = 50
	}

	if r.MinNameLength == 0 {
		r.MinNameLength = 5
	}

	if r.MaxDescLength == 0 {
		r.MaxDescLength = 100
	}
}

type User struct {
	MaxNicknameLength int `json:"max_nickname_length"`
	MaxBioLength      int `json:"max_bio_length"`
}

func (u *User) setDefault() {
	if u.MaxNicknameLength == 0 {
		u.MaxNicknameLength = 20
	}

	if u.MaxBioLength == 0 {
		u.MaxBioLength = 200
	}
}
