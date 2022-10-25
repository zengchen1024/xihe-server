package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/inference"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type InferenceIndex = domain.InferenceIndex

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
	v.ProjectId = cmd.ProjectId
	v.LastCommit = lastCommit
	v.ProjectName = cmd.ProjectName
	v.ProjectOwner = cmd.ProjectOwner
}

type InferenceService interface {
	Create(*UserInfo, *InferenceCreateCmd) (InferenceDTO, string, error)
	Get(info *InferenceIndex) (InferenceDTO, error)

	CreateInferenceInstance(*domain.InferenceInfo) error
	ExtendExpiryForInstance(*domain.InferenceInfo) error
}

func NewInferenceService(
	p platform.RepoFile,
	repo repository.Inference,
	sender message.Sender,
	manager inference.Inference,
	minExpiryForInference int64,
	survivalTimeForInstance int64,
) InferenceService {
	return inferenceService{
		p:                       p,
		repo:                    repo,
		sender:                  sender,
		manager:                 manager,
		minExpiryForInference:   minExpiryForInference,
		survivalTimeForInstance: survivalTimeForInstance,
	}
}

type inferenceService struct {
	p                       platform.RepoFile
	repo                    repository.Inference
	sender                  message.Sender
	manager                 inference.Inference
	minExpiryForInference   int64
	survivalTimeForInstance int64
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
	if err != nil || dto.hasResult() {
		if dto.canReuseCurrent() {
			instance.Id = dto.InstanceId
			err = s.sender.ExtendInferenceExpiry(&instance.InferenceInfo)
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
	index := domain.ResourceIndex{
		Owner: instance.ProjectOwner,
		Id:    instance.ProjectId,
	}
	v, version, err := s.repo.FindInstances(&index, instance.LastCommit)
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

func (s inferenceService) CreateInferenceInstance(info *domain.InferenceInfo) error {
	return s.manager.Create(info)
}

func (s inferenceService) ExtendExpiryForInstance(info *domain.InferenceInfo) error {
	return s.manager.ExtendExpiry(info, utils.Now()+s.survivalTimeForInstance)
}
