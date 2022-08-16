package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
	"github.com/opensourceways/xihe-server/infrastructure/mq"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
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

	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	cfg, err := config.LoadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	// TODO verify log
	log := logrus.NewEntry(logrus.StandardLogger())

	if err := mq.Init(nil, log); err != nil {
		log.Fatalf("initialize mq failed, err:%v", err)
	}

	defer mq.Exit(log)

	m := &cfg.Mongodb
	if err := mongodb.Initialize(m.MongodbConn, m.DBName); err != nil {
		logrus.Fatalf("initialize mongodb failed, err:%s", err.Error())
	}

	defer mongodb.Close()

	initDomainConfig(cfg)

	run(cfg, log)
}

func initDomainConfig(cfg *config.Config) {
	r := &cfg.Resource
	u := &cfg.User

	domain.Init(domain.Config{
		Resource: domain.ResourceConfig{
			MaxNameLength: r.MaxNameLength,
			MinNameLength: r.MinNameLength,
			MaxDescLength: r.MaxDescLength,

			Covers:           sets.NewString(r.Covers...),
			Protocols:        sets.NewString(r.Protocols...),
			ProjectType:      sets.NewString(r.ProjectType...),
			TrainingPlatform: sets.NewString(r.TrainingPlatform...),
		},

		User: domain.UserConfig{
			MaxNicknameLength: u.MaxNicknameLength,
			MaxBioLength:      u.MaxBioLength,
		},
	})
}

func run(cfg *config.Config, log *logrus.Entry) {
	h := handler{
		user: app.NewUserService(
			repositories.NewUserRepository(
				mongodb.NewUserMapper(cfg.Mongodb.UserCollection),
			),
			nil,
		),
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	ctx, done := context.WithCancel(context.Background())
	// it seems that it will be ok even if invoking 'done' twice.
	defer done()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			logrus.Info("receive done. exit normally")
			return
		case <-sig:
			logrus.Info("receive exit signal")
			done()
			return
		}
	}(ctx)

	if err := mq.Subscribe(ctx, h, log); err != nil {
		log.Errorf("subscribe failed, err:%v", err)
	}
}
