package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type InferenceInfo = domain.InferenceInfo

type InferenceCreateCmd struct {
	ProjectId     string
	ProjectName   domain.ProjName
	ProjectOwner  domain.Account
	ProjectRepoId string

	InferenceDir domain.Directory
	BootFile     domain.FilePath
}

func (cmd *InferenceCreateCmd) Validate() error {
	b := cmd.ProjectId == "" ||
		cmd.ProjectName == nil ||
		cmd.ProjectOwner == nil ||
		cmd.ProjectRepoId == "" ||
		cmd.InferenceDir == nil ||
		cmd.BootFile == nil

	if b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *InferenceCreateCmd) toInference(v *domain.Infereance, lastCommit string) {
	v.ProjectId = cmd.ProjectId
	v.LastCommit = lastCommit
	v.ProjectOwner = cmd.ProjectOwner
	v.ProjectRepoId = cmd.ProjectRepoId
}

type InferenceService interface {
	Create(*UserInfo, *InferenceCreateCmd) (InferenceDTO, error)
	Get(info *domain.InferenceInfo) (InferenceDTO, error)
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
	dto InferenceDTO, err error,
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

	index := domain.InferenceIndex{
		ProjectId:    cmd.ProjectId,
		LastCommit:   sha,
		ProjectOwner: cmd.ProjectOwner,
	}

	dto, version, err := s.check(&index)
	if err != nil || dto.hasResult() {
		if dto.canReuseCurrent() {
			err = s.sender.ExtendExpiry(&domain.InferenceInfo{
				Id:           dto.InstanceId,
				ProjectId:    index.ProjectId,
				ProjectOwner: index.ProjectOwner,
			})
		}

		return
	}

	v := new(domain.Infereance)
	cmd.toInference(v, sha)

	if dto.InstanceId, err = s.repo.Save(v, version); err == nil {
		err = s.sender.CreateInference(&domain.InferenceInfo{
			Id:           dto.InstanceId,
			ProjectId:    index.ProjectId,
			ProjectOwner: index.ProjectOwner,
		})

		return
	}

	if repository.IsErrorDuplicateCreating(err) {
		dto, _, err = s.check(&index)
	}

	return
}

func (s inferenceService) Get(info *domain.InferenceInfo) (dto InferenceDTO, err error) {
	v, err := s.repo.FindInstance(info)

	dto.Error = v.Error
	dto.AccessURL = v.AccessURL
	dto.InstanceId = info.Id

	return
}

func (s inferenceService) check(index *domain.InferenceIndex) (
	dto InferenceDTO, version int, err error,
) {
	v, version, err := s.repo.FindInstances(index)
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
