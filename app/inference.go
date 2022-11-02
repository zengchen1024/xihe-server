package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/inference"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

type InferenceIndex = domain.InferenceIndex
type InferenceDetail = domain.InferenceDetail

type InferenceCreateCmd struct {
	ProjectId    string
	ProjectName  domain.ProjName
	ProjectOwner domain.Account

	InferenceDir domain.Directory
	BootFile     domain.FilePath
}

func (cmd *InferenceCreateCmd) Validate() error {
	b := cmd.ProjectId == "" ||
		cmd.ProjectName == nil ||
		cmd.ProjectOwner == nil ||
		cmd.InferenceDir == nil ||
		cmd.BootFile == nil

	if b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *InferenceCreateCmd) toInference(v *domain.Inference, lastCommit string) {
	v.Project.Id = cmd.ProjectId
	v.LastCommit = lastCommit
	v.ProjectName = cmd.ProjectName
	v.Project.Owner = cmd.ProjectOwner
}

type InferenceService interface {
	Create(*UserInfo, *InferenceCreateCmd) (InferenceDTO, string, error)
	Get(info *InferenceIndex) (InferenceDTO, error)
}

func NewInferenceService(
	p platform.RepoFile,
	repo repository.Inference,
	sender message.Sender,
	minExpiryForInference int64,
) InferenceService {
	return inferenceService{
		p:                     p,
		repo:                  repo,
		sender:                sender,
		minExpiryForInference: minExpiryForInference,
	}
}

type inferenceService struct {
	p                     platform.RepoFile
	repo                  repository.Inference
	sender                message.Sender
	minExpiryForInference int64
}

type InferenceDTO struct {
	Error      string
	AccessURL  string
	InstanceId string
}

func (dto *InferenceDTO) hasResult() bool {
	return dto.InstanceId != ""
}

func (dto *InferenceDTO) canReuseCurrent() bool {
	return dto.AccessURL != ""
}

func (s inferenceService) Create(u *UserInfo, cmd *InferenceCreateCmd) (
	dto InferenceDTO, sha string, err error,
) {
	sha, b, err := s.p.GetDirFileInfo(u, &platform.RepoDirFile{
		RepoName: cmd.ProjectName,
		Dir:      cmd.InferenceDir,
		File:     cmd.BootFile,
	})
	if err != nil {
		return
	}

	if !b {
		err = ErrorUnavailableRepoFile{
			errors.New("no boot file"),
		}
	}

	instance := new(domain.Inference)
	cmd.toInference(instance, sha)

	dto, version, err := s.check(instance)
	if err != nil {
		return
	}

	if dto.hasResult() {
		if dto.canReuseCurrent() {
			instance.Id = dto.InstanceId
			logrus.Debugf("will reuse the inference instance(%s)", dto.InstanceId)

			if err1 := s.sender.ExtendInferenceExpiry(&instance.InferenceInfo); err1 != nil {
				logrus.Errorf(
					"extend instance(%s) failed, err:%s",
					dto.InstanceId, err1.Error(),
				)
			}
		}

		return
	}

	if dto.InstanceId, err = s.repo.Save(instance, version); err == nil {
		instance.Id = dto.InstanceId
		err = s.sender.CreateInference(&instance.InferenceInfo)

		return
	}

	if repository.IsErrorDuplicateCreating(err) {
		dto, _, err = s.check(instance)
	}

	return
}

func (s inferenceService) Get(index *InferenceIndex) (dto InferenceDTO, err error) {
	v, err := s.repo.FindInstance(index)

	dto.Error = v.Error
	dto.AccessURL = v.AccessURL
	dto.InstanceId = v.Id

	return
}

func (s inferenceService) check(instance *domain.Inference) (
	dto InferenceDTO, version int, err error,
) {
	v, version, err := s.repo.FindInstances(&instance.Project, instance.LastCommit)
	if err != nil || len(v) == 0 {
		return
	}

	var target *repository.InferenceSummary

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
	if n < e && n+s.minExpiryForInference <= e {
		dto.AccessURL = target.AccessURL
		dto.InstanceId = target.Id
	}

	return
}

type InferenceInternalService interface {
	UpdateDetail(*InferenceIndex, *InferenceDetail) error
}

func NewInferenceInternalService(repo repository.Inference) InferenceInternalService {
	return inferenceInternalService{
		repo: repo,
	}
}

type inferenceInternalService struct {
	repo repository.Inference
}

func (s inferenceInternalService) UpdateDetail(index *InferenceIndex, detail *InferenceDetail) error {
	return s.repo.UpdateDetail(index, detail)
}

type InferenceMessageService interface {
	CreateInferenceInstance(*domain.InferenceInfo) error
	ExtendExpiryForInstance(*domain.InferenceInfo) error
}

func NewInferenceMessageService(
	repo repository.Inference,
	user repository.User,
	manager inference.Inference,
	survivalTimeForInstance int,
) InferenceMessageService {
	return inferenceMessageService{
		repo:                    repo,
		user:                    user,
		manager:                 manager,
		survivalTimeForInstance: survivalTimeForInstance,
	}
}

type inferenceMessageService struct {
	repo                    repository.Inference
	user                    repository.User
	manager                 inference.Inference
	survivalTimeForInstance int
}

func (s inferenceMessageService) CreateInferenceInstance(info *domain.InferenceInfo) error {
	v, err := s.user.GetByAccount(info.Project.Owner)
	if err != nil {
		return err
	}

	return s.manager.Create(&inference.InferenceInfo{
		InferenceInfo: info,
		UserToken:     v.PlatformToken,
		SurvivalTime:  s.survivalTimeForInstance,
	})
}

func (s inferenceMessageService) ExtendExpiryForInstance(info *domain.InferenceInfo) error {
	v := utils.Now() + int64(s.survivalTimeForInstance)

	if err := s.manager.ExtendExpiry(&info.InferenceIndex, v); err != nil {
		return err
	}

	return s.repo.UpdateDetail(&info.InferenceIndex, &domain.InferenceDetail{Expiry: v})
}
