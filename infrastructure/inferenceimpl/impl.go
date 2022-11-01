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

func (impl inferenceImpl) Create(info *domain.InferenceInfo, token string) error {
	opt := sdk.InferenceCreateOption{
		UserToken:   token,
		LastCommit:  info.LastCommit,
		ProjectName: info.ProjectName.ResourceName(),
	}
	opt.User = info.Project.Owner.Account()
	opt.ProjectId = info.Project.Id
	opt.InferenceId = info.Id

	return impl.cli.CreateInference(&opt)
}

func (impl inferenceImpl) ExtendExpiry(index *domain.InferenceIndex, expiry int64) error {
	opt := sdk.InferenceUpdateOption{
		Expiry: expiry,
	}
	opt.User = index.Project.Owner.Account()
	opt.ProjectId = index.Project.Id
	opt.InferenceId = index.Id

	return impl.cli.ExtendExpiryOfInference(&opt)
}
