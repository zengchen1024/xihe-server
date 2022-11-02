package evaluateimpl

import (
	"github.com/opensourceways/xihe-inference-evaluate/sdk"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/evaluate"
)

func NewEvaluate(cfg *Config) evaluate.Evaluate {
	v := sdk.NewInferenceEvaluate(cfg.ContainerManagerEndpoint)

	return evaluateImpl{
		cli: &v,
	}
}

type evaluateImpl struct {
	cli *sdk.InferenceEvaluate
}

func (impl evaluateImpl) Create(info *evaluate.EvaluateInfo) error {
	switch info.EvaluateType {
	case domain.EvaluateTypeCustom:
		opt := &sdk.CustomEvaluateCreateOption{
			AimPath: info.OBSPath,
		}
		// TODO survival time
		opt.User = info.Project.Owner.Account()
		opt.ProjectId = info.Project.Id
		opt.TrainingId = info.TrainingId
		opt.EvaluateId = info.Id

		return impl.cli.CreateCustomEvaluate(opt)

	case domain.EvaluateTypeStandard:
		opt := &sdk.StandardEvaluateCreateOption{
			LogPath:           info.OBSPath,
			MomentumScope:     ([]string)(info.MomentumScope),
			BatchSizeScope:    ([]string)(info.BatchSizeScope),
			LearningRateScope: ([]string)(info.LearningRateScope),
		}
		opt.User = info.Project.Owner.Account()
		opt.ProjectId = info.Project.Id
		opt.TrainingId = info.TrainingId
		opt.EvaluateId = info.Id

		return impl.cli.CreateStandardEvaluate(opt)

	default:
		return nil
	}
}
