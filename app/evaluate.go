package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/evaluate"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type EvaluateIndex = domain.EvaluateIndex
type EvaluateDetail = domain.EvaluateDetail

type CustomEvaluateCreateCmd struct {
	ProjectOwner domain.Account
	ProjectId    string
	TrainingId   string
	AimPath      string
}

func (cmd *CustomEvaluateCreateCmd) Validate() error {
	b := cmd.ProjectId == "" ||
		cmd.ProjectOwner == nil ||
		cmd.TrainingId == "" ||
		cmd.AimPath != ""

	if b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *CustomEvaluateCreateCmd) toEvaluate(v *domain.Evaluate) {
	v.EvaluateType = domain.EvaluateTypeCustom
	v.Project.Id = cmd.ProjectId
	v.TrainingId = cmd.TrainingId
	v.Project.Owner = cmd.ProjectOwner
}

// standard
type StandardEvaluateCreateCmd struct {
	ProjectOwner      domain.Account
	ProjectId         string
	TrainingId        string
	LogPath           string
	MomentumScope     domain.EvaluateScope
	BatchSizeScope    domain.EvaluateScope
	LearningRateScope domain.EvaluateScope
}

func (cmd *StandardEvaluateCreateCmd) Validate() error {
	b := cmd.ProjectId == "" ||
		cmd.ProjectOwner == nil ||
		cmd.TrainingId == "" ||
		cmd.LogPath != ""

	if b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *StandardEvaluateCreateCmd) toEvaluate(v *domain.Evaluate) {
	v.EvaluateType = domain.EvaluateTypeStandard
	v.Project.Id = cmd.ProjectId
	v.TrainingId = cmd.TrainingId
	v.Project.Owner = cmd.ProjectOwner

	v.MomentumScope = cmd.MomentumScope
	v.BatchSizeScope = cmd.BatchSizeScope
	v.LearningRateScope = cmd.LearningRateScope
}

type EvaluateService interface {
	CreateCustom(*CustomEvaluateCreateCmd) (EvaluateDTO, error)
	CreateStandard(*StandardEvaluateCreateCmd) (EvaluateDTO, error)
	Get(info *EvaluateIndex) (EvaluateDTO, error)
}

func NewEvaluateService(
	repo repository.Evaluate,
	sender message.Sender,
	minExpiryForEvaluate int64,
) EvaluateService {
	return evaluateService{
		repo:                 repo,
		sender:               sender,
		minExpiryForEvaluate: minExpiryForEvaluate,
	}
}

type evaluateService struct {
	repo                 repository.Evaluate
	sender               message.Sender
	minExpiryForEvaluate int64
}

type EvaluateDTO struct {
	Error      string
	AccessURL  string
	InstanceId string
}

func (dto *EvaluateDTO) hasResult() bool {
	return dto.InstanceId != ""
}

func (s evaluateService) CreateCustom(cmd *CustomEvaluateCreateCmd) (
	dto EvaluateDTO, err error,
) {
	instance := new(domain.Evaluate)
	cmd.toEvaluate(instance)

	return s.create(instance, cmd.AimPath)
}

func (s evaluateService) CreateStandard(cmd *StandardEvaluateCreateCmd) (
	EvaluateDTO, error,
) {
	instance := new(domain.Evaluate)
	cmd.toEvaluate(instance)

	return s.create(instance, cmd.LogPath)
}

func (s evaluateService) Get(index *EvaluateIndex) (dto EvaluateDTO, err error) {
	v, err := s.repo.FindInstance(index)

	dto.Error = v.Error
	dto.AccessURL = v.AccessURL
	dto.InstanceId = v.Id

	return
}

func (s evaluateService) create(instance *domain.Evaluate, path string) (
	dto EvaluateDTO, err error,
) {
	dto, version, err := s.check(&instance.EvaluateIndex)
	if err != nil || dto.hasResult() {
		return
	}

	// TODO: limit the max evluate times in a day.
	if dto.InstanceId, err = s.repo.Save(instance, version); err == nil {
		instance.Id = dto.InstanceId

		err = s.sender.CreateEvaluate(&message.EvaluateInfo{
			EvaluateIndex: instance.EvaluateIndex,
			OBSPath:       path,
			Type:          instance.EvaluateType,
		})

		return
	}

	if repository.IsErrorDuplicateCreating(err) {
		dto, _, err = s.check(&instance.EvaluateIndex)
	}

	return
}

func (s evaluateService) check(instance *domain.EvaluateIndex) (
	dto EvaluateDTO, version int, err error,
) {
	v, version, err := s.repo.FindInstances(&instance.Project, instance.TrainingId)
	if err != nil || len(v) == 0 {
		return
	}

	var target *repository.EvaluateSummary

	for i := range v {
		item := &v[i]

		if item.Error != "" {
			dto.Error = item.Error
			dto.InstanceId = item.Id

			return
		}

		if item.Expiry == 0 && item.AccessURL == "" {
			dto.InstanceId = item.Id

			return
		}

		if target == nil || item.Expiry > target.Expiry {
			target = item
		}
	}

	if target == nil {
		return
	}

	e, n := target.Expiry, utils.Now()
	if n < e && n+s.minExpiryForEvaluate <= e {
		dto.AccessURL = target.AccessURL
		dto.InstanceId = target.Id
	}

	return
}

type EvaluateInternalService interface {
	UpdateDetail(*EvaluateIndex, *EvaluateDetail) error
}

func NewEvaluateInternalService(repo repository.Evaluate) EvaluateInternalService {
	return evaluateInternalService{
		repo: repo,
	}
}

type evaluateInternalService struct {
	repo repository.Evaluate
}

func (s evaluateInternalService) UpdateDetail(index *EvaluateIndex, detail *EvaluateDetail) error {
	return s.repo.UpdateDetail(index, detail)
}

type EvaluateMessageService interface {
	CreateEvaluateInstance(*message.EvaluateInfo) error
}

func NewEvaluateMessageService(
	repo repository.Evaluate,
	manager evaluate.Evaluate,
	survivalTimeForInstance int,
) EvaluateMessageService {
	return evaluateMessageService{
		repo:                    repo,
		manager:                 manager,
		survivalTimeForInstance: survivalTimeForInstance,
	}
}

type evaluateMessageService struct {
	repo                    repository.Evaluate
	manager                 evaluate.Evaluate
	survivalTimeForInstance int
}

func (s evaluateMessageService) CreateEvaluateInstance(info *message.EvaluateInfo) error {
	switch info.Type {
	case domain.EvaluateTypeCustom:
		return s.manager.Create(&evaluate.EvaluateInfo{
			Evaluate: &domain.Evaluate{
				EvaluateIndex: info.EvaluateIndex,
				EvaluateType:  info.Type,
			},
			OBSPath:      info.OBSPath,
			SurvivalTime: s.survivalTimeForInstance,
		})

	case domain.EvaluateTypeStandard:
		p, err := s.repo.GetStandardEvaluateParms(&info.EvaluateIndex)
		if err != nil {
			return err
		}

		return s.manager.Create(&evaluate.EvaluateInfo{
			Evaluate: &domain.Evaluate{
				EvaluateIndex:         info.EvaluateIndex,
				EvaluateType:          info.Type,
				StandardEvaluateParms: p,
			},
			OBSPath:      info.OBSPath,
			SurvivalTime: s.survivalTimeForInstance,
		})

	default:
		return nil
	}
}
