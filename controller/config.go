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
	MaxPictureSizeToDescribe int64  `json:"-"`
	MaxPictureSizeToVQA      int64  `json:"-"`
}

func (cfg *APIConfig) Validate() error {
	_, err := domain.NewPassword(cfg.DefaultPassword)

	return err
}
