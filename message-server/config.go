package main

import (
	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/config"
)

type configuration struct {
	MaxRetry int `json:"max_retry"`

	Resource config.Resource `json:"resource" required:"true"`
	Mongodb  config.Mongodb  `json:"mongodb" required:"true"`
	User     config.User     `json:"user"`
	MQ       config.MQ       `json:"mq" required:"true"`
}

func (cfg *configuration) getMQConfig() mq.MQConfig {
	return mq.MQConfig{
		Addresses: cfg.MQ.ParseAddress(),
	}
}

func (cfg *configuration) SetDefault() {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}

	items := []interface{}{
		&cfg.Resource,
		&cfg.Mongodb,
		&cfg.User,
		&cfg.MQ,
	}
	for _, i := range items {
		if f, ok := i.(config.ConfigSetDefault); ok {
			f.SetDefault()
		}
	}
}

func (cfg *configuration) Validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := []interface{}{
		&cfg.Resource,
		&cfg.Mongodb,
		&cfg.User,
		&cfg.MQ,
	}
	for _, i := range items {
		if f, ok := i.(config.ConfigValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}
