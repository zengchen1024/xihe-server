package mq

import (
	"github.com/opensourceways/community-robot-lib/kafka"
	libmq "github.com/opensourceways/community-robot-lib/mq"
	"github.com/sirupsen/logrus"
)

func Init(address []string, log *logrus.Entry) error {
	err := kafka.Init(
		libmq.Addresses(address...),
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
