package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type ProjectCreateCmd struct {
	Owner    domain.Account
	Name     domain.ProjName
	Desc     domain.ResourceDesc
	Type     domain.ProjType
	CoverId  domain.CoverId
	RepoType domain.RepoType
	Protocol domain.ProtocolName
	Training domain.TrainingPlatform
}

func (cmd *ProjectCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Name != nil &&
		cmd.Desc != nil &&
		cmd.Type != nil &&
		cmd.CoverId != nil &&
		cmd.RepoType != nil &&
		cmd.Protocol != nil &&
		cmd.Training != nil

	if !b {
		return errors.New("invalid cmd of creating project")
	}

	return nil
}

func (cmd *ProjectCreateCmd) toProject() domain.Project {
	now := utils.Now()

	return domain.Project{
		Owner:     cmd.Owner,
		Type:      cmd.Type,
		Protocol:  cmd.Protocol,
		Training:  cmd.Training,
		CreatedAt: now,
		UpdatedAt: now,
		ProjectModifiableProperty: domain.ProjectModifiableProperty{
			Name:     cmd.Name,
			Desc:     cmd.Desc,
			CoverId:  cmd.CoverId,
			RepoType: cmd.RepoType,
		},
	}
}

type ProjectDTO struct {
	Id        string   `json:"id"`
	Owner     string   `json:"owner"`
	Name      string   `json:"name"`
	Desc      string   `json:"desc"`
	Type      string   `json:"type"`
	CoverId   string   `json:"cover_id"`
	Protocol  string   `json:"protocol"`
	Training  string   `json:"training"`
	RepoType  string   `json:"repo_type"`
	RepoId    string   `json:"repo_id"`
	Tags      []string `json:"tags"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	LikeCount int      `json:"like_count"`
	ForkCount int      `json:"fork_count"`
}

type ProjectDetailDTO struct {
	ProjectDTO

	RelatedModels   []ResourceDTO `json:"related_models"`
	RelatedDatasets []ResourceDTO `json:"related_datasets"`
}

type ProjectService interface {
	Create(*ProjectCreateCmd, platform.Repository) (ProjectDTO, error)
	GetByName(domain.Account, domain.ProjName, bool) (ProjectDetailDTO, error)
	List(domain.Account, *ResourceListCmd) ([]ProjectDTO, error)
	Update(*domain.Project, *ProjectUpdateCmd, platform.Repository) (ProjectDTO, error)
	Fork(*ProjectForkCmd, platform.Repository) (ProjectDTO, error)

	IncreaseFork(index domain.ResourceIndex) error

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error

	AddRelatedModel(*domain.Project, *domain.ResourceIndex) error
	RemoveRelatedModel(*domain.Project, *domain.ResourceIndex) error

	AddRelatedDataset(*domain.Project, *domain.ResourceIndex) error
	RemoveRelatedDataset(*domain.Project, *domain.ResourceIndex) error

	SetTags(*domain.Project, *ResourceTagsUpdateCmd) error
}

func NewProjectService(
	user repository.User,
	repo repository.Project,
	model repository.Model,
	dataset repository.Dataset,
	activity repository.Activity,
	pr platform.Repository,
	sender message.Sender,
) ProjectService {
	return projectService{
		repo:     repo,
		activity: activity,
		sender:   sender,
		rs: resourceService{
			user:    user,
			model:   model,
			project: repo,
			dataset: dataset,
		},
	}
}

type projectService struct {
	repo repository.Project
	//pr       platform.Repository
	activity repository.Activity
	sender   message.Sender
	rs       resourceService
}

func (s projectService) Create(cmd *ProjectCreateCmd, pr platform.Repository) (dto ProjectDTO, err error) {
	// step1: create repo on gitlab
	pid, err := pr.New(&platform.RepoOption{
		Name:     cmd.Name,
		RepoType: cmd.RepoType,
	})
	if err != nil {
		return
	}

	// step2: save
	v := cmd.toProject()
	v.RepoId = pid

	p, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	s.toProjectDTO(&p, &dto)

	// add activity
	ua := genActivityForCreatingResource(
		p.Owner, domain.ResourceTypeProject, p.Id,
	)
	// ignore the error
	_ = s.activity.Save(&ua)

	return
}

func (s projectService) GetByName(
	owner domain.Account, name domain.ProjName,
	allowPrivacy bool,
) (dto ProjectDetailDTO, err error) {
	v, err := s.repo.GetByName(owner, name)
	if err != nil {
		return
	}

	if !allowPrivacy && v.IsPrivate() {
		err = ErrorPrivateRepo{errors.New("private repo")}

		return
	}

	s.toProjectDTO(&v, &dto.ProjectDTO)

	m, err := s.rs.listModels(v.RelatedModels)
	if err != nil {
		return
	}
	dto.RelatedModels = m

	d, err := s.rs.listModels(v.RelatedModels)
	if err != nil {
		return
	}
	dto.RelatedDatasets = d

	return
}

type ResourceListCmd struct {
	Name     string
	RepoType domain.RepoType
}

func (cmd *ResourceListCmd) toResourceListOption() (
	option repository.ResourceListOption,
) {
	option.Name = cmd.Name
	option.RepoType = cmd.RepoType

	return
}

func (s projectService) List(owner domain.Account, cmd *ResourceListCmd) (
	dtos []ProjectDTO, err error,
) {
	v, err := s.repo.List(owner, cmd.toResourceListOption())
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]ProjectDTO, len(v))
	for i := range v {
		s.toProjectDTO(&v[i], &dtos[i])
	}

	return
}

func (s projectService) toProjectDTO(p *domain.Project, dto *ProjectDTO) {
	*dto = ProjectDTO{
		Id:        p.Id,
		Owner:     p.Owner.Account(),
		Name:      p.Name.ProjName(),
		Desc:      p.Desc.ResourceDesc(),
		Type:      p.Type.ProjType(),
		CoverId:   p.CoverId.CoverId(),
		Protocol:  p.Protocol.ProtocolName(),
		Training:  p.Training.TrainingPlatform(),
		RepoType:  p.RepoType.RepoType(),
		RepoId:    p.RepoId,
		Tags:      p.Tags,
		CreatedAt: utils.ToDate(p.CreatedAt),
		UpdatedAt: utils.ToDate(p.UpdatedAt),
		LikeCount: p.LikeCount,
		ForkCount: p.ForkCount,
	}
}
