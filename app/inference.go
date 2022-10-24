package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type InferenceCreateCmd struct {
	ProjectId     string
	ProjectName   domain.ProjName
	ProjectOwner  domain.Account
	ProjectRepoId string

	InferenceDir domain.Directory
	BootFile     domain.FilePath
	ModelRef     domain.ResourceRef
}

func (cmd *InferenceCreateCmd) Validate() error {
	m := &cmd.ModelRef
	b := cmd.ProjectId == "" ||
		cmd.ProjectName == nil ||
		cmd.ProjectOwner == nil ||
		cmd.ProjectRepoId == "" ||
		cmd.InferenceDir == nil ||
		cmd.BootFile == nil ||
		m.User == nil ||
		m.RepoId == "" ||
		m.File == ""

	if b {
		return errors.New("invalid cmd")
	}

	return nil
}

func (cmd *InferenceCreateCmd) toInference(v *domain.Infereance, lastCommit string) {
	v.ModelRef = cmd.ModelRef
	v.ProjectId = cmd.ProjectId
	v.LastCommit = lastCommit
	v.ProjectOwner = cmd.ProjectOwner
	v.ProjectRepoId = cmd.ProjectRepoId
}

type inferenceService struct {
	p      platform.RepoFile
	repo   repository.Inference
	sender message.Sender
}

type InferenceDTO struct {
	Error     string `json:"error"`
	AccessURL string `json:"access_url"`

	currentInstanceId  string `json:"-"`
	CreatingInstanceId string `json:"-"`
}

func (dto *InferenceDTO) hasResult() bool {
	return dto.currentInstanceId != "" || dto.CreatingInstanceId != ""
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
				Id:             dto.currentInstanceId,
				InferenceIndex: index,
			})
		}

		return
	}

	v := new(domain.Infereance)
	cmd.toInference(v, sha)

	if dto.CreatingInstanceId, err = s.repo.Save(v, version); err == nil {
		err = s.sender.CreateInference(&domain.InferenceInfo{
			Id:             dto.CreatingInstanceId,
			InferenceIndex: index,
		})

		return
	}

	if repository.IsErrorDuplicateCreating(err) {
		dto, _, err = s.check(&index)
	}

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
			dto.currentInstanceId = item.Id

			return
		}

		if item.IsCreating() {
			dto.CreatingInstanceId = item.Id

			return
		}

		if target == nil || item.Expiry > target.Expiry {
			target = item
		}
	}

	if target == nil {
		return
	}

	if e, n := target.Expiry, utils.Now(); n < e && n+3600 < e {
		dto.AccessURL = target.AccessURL
		dto.currentInstanceId = target.Id
	}

	return
}
