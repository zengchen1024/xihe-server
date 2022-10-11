package trainingimpl

import (
	"github.com/opensourceways/xihe-training-center/sdk"

	"github.com/opensourceways/xihe-server/domain"
)

type trainingImpl struct {
}

func (impl trainingImpl) CreateJob(endpoint string, user domain.Account, t *domain.Training) (
	job domain.JobInfo, err error,
) {
	opt := sdk.TrainingCreateOption{
		User:           user.Account(),
		ProjectName:    t.ProjectName.ProjName(),
		ProjectRepoId:  t.ProjectRepoId,
		Name:           t.Name.TrainingName(),
		CodeDir:        t.CodeDir.Directory(),
		BootFile:       t.BootFile.FilePath(),
		Compute:        impl.toCompute(&t.Compute),
		Env:            impl.toKeyValue(t.Env),
		Inputs:         impl.toInput(t.Inputs),
		Hypeparameters: impl.toKeyValue(t.Hypeparameters),
	}

	if t.Desc != nil {
		opt.Desc = t.Desc.TrainingDesc()
	}

	cli := sdk.NewTrainingCenter(endpoint)

	v, err := cli.CreateTraining(&opt)
	if err != nil {
		return
	}

	job.Endpoint = endpoint
	job.JobId = v.JobId
	job.LogDir = v.LogDir
	job.OutputDir = v.OutputDir

	return
}

func (t trainingImpl) DeleteJob(endpoint, jobId string) error {
	cli := sdk.NewTrainingCenter(endpoint)

	return cli.DeleteTraining(jobId)
}

func (t trainingImpl) TerminateJob(endpoint, jobId string) error {
	cli := sdk.NewTrainingCenter(endpoint)

	return cli.TerminateTraining(jobId)
}

func (t trainingImpl) GetLogURL(endpoint, jobId string) (string, error) {
	cli := sdk.NewTrainingCenter(endpoint)

	v, err := cli.GetLog(jobId)
	if err != nil {
		return "", err
	}

	return v.LogURL, nil
}

func (t trainingImpl) GetJob(endpoint, jobId string) (r domain.JobDetail, err error) {
	cli := sdk.NewTrainingCenter(endpoint)

	v, err := cli.GetTraining(jobId)
	if err != nil {
		return
	}

	r.Duration = v.Duration
	r.Status = v.Status

	return
}

func (t trainingImpl) toCompute(c *domain.Compute) sdk.Compute {
	return sdk.Compute{
		Type:    c.Type.ComputeType(),
		Version: c.Version.ComputeVersion(),
		Flavor:  c.Flavor.ComputeFlavor(),
	}
}

func (t trainingImpl) toKeyValue(kv []domain.KeyValue) []sdk.KeyValue {
	if len(kv) == 0 {
		return nil
	}

	r := make([]sdk.KeyValue, len(kv))

	for i := range kv {
		s := ""
		if kv[i].Value != nil {
			s = kv[i].Value.CustomizedValue()
		}

		r[i] = sdk.KeyValue{
			Key:   kv[i].Key.CustomizedKey(),
			Value: s,
		}
	}

	return r
}

func (t trainingImpl) toInput(v []domain.Input) []sdk.Input {
	r := make([]sdk.Input, len(v))

	for i := range v {
		input := &v[i].Value

		r[i] = sdk.Input{
			Key: v[i].Key.CustomizedKey(),
			Value: sdk.ResourceInput{
				Owner:  input.User.Account(),
				Type:   input.Type.ResourceType(),
				RepoId: input.RepoId,
				File:   input.File,
			},
		}
	}

	return r
}
