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
	"github.com/opensourceways/xihe-server/infrastructure/message"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
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
	log := logrus.NewEntry(logrus.StandardLogger())

	o := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	// cfg
	cfg, err := config.LoadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	// mq
	if err := message.Init(cfg.MQ.Addresses, log); err != nil {
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
	initDomainConfig(cfg)

	// run
	run(newHandler(cfg, log), log)
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

func newHandler(cfg *config.Config, log *logrus.Entry) *handler {
	return &handler{
		log:      log,
		maxRetry: cfg.MaxRetry,
		user: app.NewUserService(
			repositories.NewUserRepository(
				mongodb.NewUserMapper(cfg.Mongodb.UserCollection),
			),
			nil, nil,
		),

		project: app.NewProjectService(
			repositories.NewProjectRepository(
				mongodb.NewProjectMapper(cfg.Mongodb.ProjectCollection),
			),
			nil,
		),

		dataset: app.NewDatasetService(
			repositories.NewDatasetRepository(
				mongodb.NewDatasetMapper(cfg.Mongodb.DatasetCollection),
			),
			nil,
		),

		model: app.NewModelService(
			repositories.NewModelRepository(
				mongodb.NewModelMapper(cfg.Mongodb.ModelCollection),
			),
			nil,
		),
	}
}

func run(h *handler, log *logrus.Entry) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	defer wg.Wait()

	called := false
	ctx, done := context.WithCancel(context.Background())

	defer func() {
		if !called {
			called = true
			done()
		}
	}()

	wg.Add(1)
	go func(ctx context.Context) {
		defer wg.Done()

		select {
		case <-ctx.Done():
			log.Info("receive done. exit normally")
			return

		case <-sig:
			log.Info("receive exit signal")
			done()
			called = true
			return
		}
	}(ctx)

	if err := message.Subscribe(ctx, h, log); err != nil {
		log.Errorf("subscribe failed, err:%v", err)
	}
}
