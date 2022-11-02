package inferenceimpl

import (
	"github.com/opensourceways/xihe-inference-evaluate/sdk"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/inference"
)

func NewInference(cfg *Config) inference.Inference {
	v := sdk.NewInferenceEvaluate(cfg.ContainerManagerEndpoint)

	return inferenceImpl{
		cli: &v,
	}
}

type inferenceImpl struct {
	cli *sdk.InferenceEvaluate
}

func (impl inferenceImpl) Create(info *inference.InferenceInfo) error {
	opt := sdk.InferenceCreateOption{
		UserToken:    info.UserToken,
		LastCommit:   info.LastCommit,
		ProjectName:  info.ProjectName.ResourceName(),
		SurvivalTime: info.SurvivalTime,
	}
	opt.User = info.Project.Owner.Account()
	opt.ProjectId = info.Project.Id
	opt.InferenceId = info.Id

	return impl.cli.CreateInference(&opt)
}

func (impl inferenceImpl) ExtendSurvivalTime(index *domain.InferenceIndex, timeToExtend int) error {
	opt := sdk.InferenceUpdateOption{
		TimeToExtend: timeToExtend,
	}
	opt.User = index.Project.Owner.Account()
	opt.ProjectId = index.Project.Id
	opt.InferenceId = index.Id

	return impl.cli.ExtendExpiryOfInference(&opt)
}
