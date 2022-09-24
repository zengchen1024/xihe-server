package main

import (
	"github.com/opensourceways/community-robot-lib/mq"
	"github.com/opensourceways/community-robot-lib/utils"

	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/domain"
)

type configuration struct {
	MaxRetry int `json:"max_retry"`

	Resource domain.ResourceConfig `json:"resource" required:"true"`
	Mongodb  config.Mongodb        `json:"mongodb"  required:"true"`
	User     domain.UserConfig     `json:"user"`
	MQ       config.MQ             `json:"mq"       required:"true"`
}

func (cfg *configuration) getMQConfig() mq.MQConfig {
	return mq.MQConfig{
		Addresses: cfg.MQ.ParseAddress(),
	}
}

func (cfg *configuration) configItems() []interface{} {
	return []interface{}{
		&cfg.Resource,
		&cfg.Mongodb,
		&cfg.User,
		&cfg.MQ,
	}
}

func (cfg *configuration) SetDefault() {
	if cfg.MaxRetry <= 0 {
		cfg.MaxRetry = 10
	}

	items := cfg.configItems()
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

	items := cfg.configItems()
	for _, i := range items {
		if f, ok := i.(config.ConfigValidate); ok {
			if err := f.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg *configuration) initDomainConfig() {
	domain.Init(&cfg.Resource, &cfg.User)
}
