package main

import (
	"github.com/opensourceways/xihe-training-center/controller"
	"github.com/opensourceways/xihe-training-center/sdk"

	"github.com/opensourceways/xihe-server/domain"
)

type trainingJob struct {
}

func (t trainingJob) CreateJob(endpoint string, ut *domain.UserTraining) (string, error) {
	opt := sdk.TrainingCreateOption{
		User:           ut.Owner.Account(),
		ProjectId:      ut.ProjectId,
		Name:           ut.Name.TrainingName(),
		CodeDir:        ut.CodeDir.Directory(),
		BootFile:       ut.BootFile.FilePath(),
		LogDir:         ut.LogDir.Directory(),
		Compute:        t.toCompute(&ut.Compute),
		Env:            t.toKeyValue(ut.Env),
		Inputs:         t.toKeyValue(ut.Inputs),
		Outputs:        t.toKeyValue(ut.Outputs),
		Hypeparameters: t.toKeyValue(ut.Hypeparameters),
	}

	if ut.Desc != nil {
		opt.Desc = ut.Desc.TrainingDesc()
	}

	cli := sdk.NewTrainingCenter(endpoint)

	return cli.CreateTraining(&opt)
}

func (t trainingJob) DeleteJob(endpoint, jobId string) error {
	cli := sdk.NewTrainingCenter(endpoint)

	return cli.DeleteTraining(jobId)
}

func (t trainingJob) TerminateJob(endpoint, jobId string) error {
	cli := sdk.NewTrainingCenter(endpoint)

	return cli.TerminateTraining(jobId)
}

func (t trainingJob) GetLogURL(endpoint, jobId string) (string, error) {
	cli := sdk.NewTrainingCenter(endpoint)

	v, err := cli.GetLog(jobId)
	if err != nil {
		return "", err
	}

	return v.LogURL, nil
}

func (t trainingJob) GetJob(endpoint, jobId string) (r domain.TrainingInfo, err error) {
	cli := sdk.NewTrainingCenter(endpoint)

	v, err := cli.GetTraining(jobId)
	if err != nil {
		return
	}

	r.Duration = v.Duration
	r.Status = v.Status

	return
}

func (t trainingJob) toCompute(c *domain.Compute) controller.Compute {
	return controller.Compute{
		Type:    c.Type.ComputeType(),
		Version: c.Version.ComputeVersion(),
		Flavor:  c.Flavor.ComputeFlavor(),
	}
}

func (t trainingJob) toKeyValue(kv []domain.KeyValue) []controller.KeyValue {
	if len(kv) == 0 {
		return nil
	}

	r := make([]controller.KeyValue, len(kv))

	for i := range kv {
		r[i] = controller.KeyValue{
			Key:   kv[i].Key.CustomizedKey(),
			Value: kv[i].Value.CustomizedValue(),
		}
	}

	return r
}
