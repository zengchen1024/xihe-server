package message

import (
	"github.com/opensourceways/community-robot-lib/kafka"
	"github.com/opensourceways/community-robot-lib/mq"
	libmq "github.com/opensourceways/community-robot-lib/mq"
	"github.com/sirupsen/logrus"
)

func Init(cfg mq.MQConfig, log *logrus.Entry, topic Topics) error {
	topics = topic

	err := kafka.Init(
		libmq.Addresses(cfg.Addresses...),
		libmq.Log(log),
	)
	if err != nil {
		return err
	}

	return kafka.Connect()
}

func Exit(log *logrus.Entry) {
	if err := kafka.Disconnect(); err != nil {
		log.Errorf("exit kafka, err:%v", err)
	}
}
