package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ProjectCreateCmd struct {
	Owner     string
	Name      domain.ProjName
	Desc      domain.ProjDesc
	Type      domain.RepoType
	CoverId   domain.CoverId
	Protocol  domain.ProtocolName
	Training  domain.TrainingSDK
	Inference domain.InferenceSDK
}

func (cmd *ProjectCreateCmd) Validate() error {
	b := cmd.Owner != "" &&
		cmd.Name != nil &&
		cmd.Type != nil &&
		cmd.CoverId != nil &&
		cmd.Protocol != nil &&
		cmd.Training != nil &&
		cmd.Inference != nil

	if !b {
		return errors.New("invalid cmd of creating project")
	}

	return nil
}

func (cmd *ProjectCreateCmd) toProject() domain.Project {
	return domain.Project{
		Owner:     cmd.Owner,
		Name:      cmd.Name,
		Desc:      cmd.Desc,
		Type:      cmd.Type,
		CoverId:   cmd.CoverId,
		Protocol:  cmd.Protocol,
		Training:  cmd.Training,
		Inference: cmd.Inference,
	}
}

type ProjectDTO struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Type      string `json:"type"`
	CoverId   string `json:"cover_id"`
	Protocol  string `json:"protocol"`
	Training  string `json:"training"`
	Inference string `json:"inference"`
}

type ProjectService interface {
	Create(cmd *ProjectCreateCmd) (ProjectDTO, error)
	Update(p *domain.Project, cmd *ProjectUpdateCmd) (ProjectDTO, error)
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

	return
}

func (s projectService) toProjectDTO(p *domain.Project, dto *ProjectDTO) {
	dto.Id = p.Id
}
