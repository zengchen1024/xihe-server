package controller

import (
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

var (
	apiConfig     APIConfig
	encryptHelper utils.SymmetricEncryption
	log           *logrus.Entry
)

func Init(cfg *APIConfig, l *logrus.Entry) error {
	log = l
	apiConfig = *cfg

	e, err := utils.NewSymmetricEncryption(cfg.EncryptionKey, "")
	if err != nil {
		return err
	}

	encryptHelper = e

	return nil
}

type APIConfig struct {
	TokenKey                 string `json:"token_key"                   required:"true"`
	TokenExpiry              int64  `json:"token_expiry"                required:"true"`
	EncryptionKey            string `json:"encryption_key"              required:"true"`
	DefaultPassword          string `json:"default_password"            required:"true"`
	MaxTrainingRecordNum     int    `json:"max_training_record_num"     required:"true"`
	MinExpiryForInference    int    `json:"min_expiry_for_inference"`
	InferenceDir             string `json:"inference_dir"`
	InferenceBootFile        string `json:"inference_boot_file"`
	InferenceTimeout         int    `json:"inference_timeout"`
	MaxPictureSizeToDescribe int64  `json:"-"`
	MaxPictureSizeToVQA      int64  `json:"-"`
}

func (cfg *APIConfig) SetDefault() {
	if cfg.MinExpiryForInference <= 0 {
		cfg.MinExpiryForInference = 3600
	}

	if cfg.InferenceDir == "" {
		cfg.InferenceDir = "inference"
	}

	if cfg.InferenceBootFile == "" {
		cfg.InferenceBootFile = "inference/app.py"
	}

	if cfg.InferenceTimeout <= 0 {
		cfg.InferenceTimeout = 300
	}
}

func (cfg *APIConfig) Validate() (err error) {
	if _, err = domain.NewPassword(cfg.DefaultPassword); err != nil {
		return
	}

	if _, err = domain.NewDirectory(cfg.InferenceDir); err != nil {
		return
	}

	_, err = domain.NewFilePath(cfg.InferenceBootFile)

	return
}
