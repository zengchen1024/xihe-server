package main

import (
	"flag"
	"os"
	"strconv"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/opensourceways/xihe-grpc-protocol/grpc/inference"
	"github.com/opensourceways/xihe-grpc-protocol/grpc/server"
	"github.com/opensourceways/xihe-grpc-protocol/grpc/training"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-server/app"
	"github.com/opensourceways/xihe-server/config"
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/infrastructure/mongodb"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func (o *options) addFlags(fs *flag.FlagSet) {
	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false,
		"whether to enable debug model.",
	)
}

func gatherOptions(fs *flag.FlagSet, args ...string) options {
	var o options

	o.addFlags(fs)

	_ = fs.Parse(args)

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

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// cfg
	cfg := new(configuration)
	if err := config.LoadConfig(o.service.ConfigFile, cfg); err != nil {
		log.Fatalf("load config, err:%s", err.Error())
	}

	// mongo
	m := &cfg.Mongodb
	if err := mongodb.Initialize(m.MongodbConn, m.DBName); err != nil {
		log.Fatalf("initialize mongodb failed, err:%s", err.Error())
	}

	defer mongodb.Close()

	// training
	train := app.NewTrainingService(
		log,
		nil,
		repositories.NewTrainingRepository(
			mongodb.NewTrainingMapper(cfg.Mongodb.TrainingCollection),
		),
		nil, 0,
	)

	// inference
	inferenceService := app.NewInferenceInternalService(
		repositories.NewInferenceRepository(
			mongodb.NewInferenceMapper(cfg.Mongodb.InferenceCollection),
		),
	)

	// cfg
	cfg.initDomainConfig()

	// server
	s := server.NewServer()

	s.RegisterTrainingServer(trainingServer{train})
	s.RegisterInferenceServer(inferenceServer{inferenceService})

	if err := s.Run(strconv.Itoa(o.service.Port)); err != nil {
		log.Errorf("start server failed, err:%s", err.Error())
	}
}

type trainingServer struct {
	service app.TrainingService
}

func (t trainingServer) SetTrainingInfo(index *training.TrainingIndex, v *training.TrainingInfo) error {
	u, err := domain.NewAccount(index.User)
	if err != nil {
		return nil
	}

	return t.service.UpdateJobDetail(
		&domain.TrainingInfo{
			User:       u,
			ProjectId:  index.ProjectId,
			TrainingId: index.Id,
		},
		&app.JobDetail{
			Duration:   v.Duration,
			Status:     v.Status,
			LogPath:    v.LogPath,
			AimPath:    v.AimZipPath,
			OutputPath: v.OutputZipPath,
		},
	)
}

type inferenceServer struct {
	service app.InferenceInternalService
}

func (t inferenceServer) SetInferenceInfo(index *inference.InferenceIndex, v *inference.InferenceInfo) error {
	u, err := domain.NewAccount(index.User)
	if err != nil {
		return nil
	}

	return t.service.UpdateDetail(
		&domain.InferenceIndex{
			Project: domain.ResourceIndex{
				Owner: u,
				Id:    index.ProjectId,
			},
			Id:         index.Id,
			LastCommit: index.LastCommit,
		},
		&app.InferenceDetail{
			Error:     v.Error,
			AccessURL: v.AccessURL,
		},
	)
}
