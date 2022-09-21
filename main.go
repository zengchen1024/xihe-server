package main

import (
	"flag"
	"os"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/controller"
	"github.com/opensourceways/xihe-server/infrastructure/authing"
	"github.com/opensourceways/xihe-server/infrastructure/gitlab"
	"github.com/opensourceways/xihe-server/infrastructure/message"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
	"github.com/opensourceways/xihe-server/server"
)

type options struct {
	service liboptions.ServiceOptions
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.service.AddFlags(fs)

	fs.Parse(args)
	return o
}

func main() {
	logrusutil.ComponentInit("xihe")
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	// cfg
	cfg := new(config.Config)
	if err := config.LoadConfig(o.service.ConfigFile, cfg); err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	// gitlab
	if err := gitlab.Init(cfg.Gitlab.Endpoint, cfg.Gitlab.RootToken); err != nil {
		logrus.Fatalf("initialize gitlab failed, err:%s", err.Error())
	}

	// authing
	authing.Init(cfg.Authing.APPId, cfg.Authing.Secret)

	// controller
	apiConfig := controller.APIConfig{
		EncryptionKey:  cfg.EncryptionKey,
		APITokenKey:    cfg.API.APITokenKey,
		APITokenExpiry: cfg.API.APITokenExpiry,
	}
	if err := controller.Init(apiConfig, log); err != nil {
		logrus.Fatalf("initialize api controller failed, err:%s", err.Error())
	}

	// mq
	topic := message.Topics{
		Like:      cfg.MQ.TopicLike,
		Following: cfg.MQ.TopicFollowing,
	}
	if err := message.Init(cfg.GetMQConfig(), log, topic); err != nil {
		log.Fatalf("initialize mq failed, err:%v", err)
	}

	defer message.Exit(log)

	// mongo
	m := &cfg.Mongodb
	if err := mongodb.Initialize(m.MongodbConn, m.DBName); err != nil {
		logrus.Fatalf("initialize mongodb failed, err:%s", err.Error())
	}

	defer mongodb.Close()

	// cfg
	cfg.InitDomainConfig()

	// run
	server.StartWebServer(o.service.Port, o.service.GracePeriod, cfg)
}
