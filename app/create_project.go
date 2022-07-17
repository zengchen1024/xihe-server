package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ProjectCreateCmd struct {
	Owner    string
	Name     domain.ProjName
	Desc     domain.ProjDesc
	Type     domain.ProjType
	CoverId  domain.CoverId
	RepoType domain.RepoType
	Protocol domain.ProtocolName
	Training domain.TrainingPlatform
}

func (cmd *ProjectCreateCmd) Validate() error {
	b := cmd.Owner != "" &&
		cmd.Name != nil &&
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
	return domain.Project{
		Owner:    cmd.Owner,
		Name:     cmd.Name,
		Desc:     cmd.Desc,
		Type:     cmd.Type,
		CoverId:  cmd.CoverId,
		RepoType: cmd.RepoType,
		Protocol: cmd.Protocol,
		Training: cmd.Training,
	}
}

type ProjectDTO struct {
	Id       string   `json:"id"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Type     string   `json:"type"`
	CoverId  string   `json:"cover_id"`
	Protocol string   `json:"protocol"`
	Training string   `json:"training"`
	RepoType string   `json:"repo_type"`
	Tags     []string `json:"tags"`
}

type ProjectService interface {
	Get(string, string) (ProjectDTO, error)
	Create(*ProjectCreateCmd) (ProjectDTO, error)
	List(string, *ProjectListCmd) ([]ProjectDTO, error)
	Update(*domain.Project, *ProjectUpdateCmd) (ProjectDTO, error)
}

func NewProjectService(repo repository.Project) ProjectService {
	return projectService{repo}
}

type projectService struct {
	repo repository.Project
}

func (s projectService) Create(cmd *ProjectCreateCmd) (dto ProjectDTO, err error) {
	p := cmd.toProject()

	v, err := s.repo.Save(&p)
	if err != nil {
		return
	}

	s.toProjectDTO(&v, &dto)

	// TODO send event

	return
}

func (s projectService) toProjectDTO(p *domain.Project, dto *ProjectDTO) {
	*dto = ProjectDTO{
		Id:       p.Id,
		Owner:    p.Owner,
		Name:     p.Name.ProjName(),
		Desc:     p.Desc.ProjDesc(),
		Type:     p.Type.ProjType(),
		CoverId:  p.CoverId.CoverId(),
		Protocol: p.Protocol.ProtocolName(),
		Training: p.Training.TrainingPlatform(),
		RepoType: p.RepoType.RepoType(),
		Tags:     p.Tags,
	}
}

func (s projectService) Get(owner, projectId string) (dto ProjectDTO, err error) {
	v, err := s.repo.Get(owner, projectId)
	if err != nil {
		return
	}

	s.toProjectDTO(&v, &dto)

	return
}
