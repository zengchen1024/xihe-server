package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ModelCreateCmd struct {
	Owner    domain.Account
	Name     domain.ModelName
	Desc     domain.ProjDesc
	RepoType domain.RepoType
	Protocol domain.ProtocolName
}

func (cmd *ModelCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Name != nil &&
		cmd.RepoType != nil &&
		cmd.Protocol != nil

	if !b {
		return errors.New("invalid cmd of creating model")
	}

	return nil
}

func (cmd *ModelCreateCmd) toModel() domain.Model {
	return domain.Model{
		Owner:    cmd.Owner,
		Protocol: cmd.Protocol,
		ModelModifiableProperty: domain.ModelModifiableProperty{
			Name:     cmd.Name,
			Desc:     cmd.Desc,
			RepoType: cmd.RepoType,
		},
	}
}

type ModelDTO struct {
	Id       string   `json:"id"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Protocol string   `json:"protocol"`
	RepoType string   `json:"repo_type"`
	RepoId   string   `json:"repo_id"`
	Tags     []string `json:"tags"`
}

type ModelService interface {
	Create(*ModelCreateCmd) (ModelDTO, error)
	Update(*domain.Model, *ModelUpdateCmd) (ModelDTO, error)
	GetByName(domain.Account, domain.ModelName) (ModelDTO, error)
	List(domain.Account, *ResourceListCmd) ([]ModelDTO, error)

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error
}

func NewModelService(repo repository.Model, pr platform.Repository) ModelService {
	return modelService{repo: repo, pr: pr}
}

type modelService struct {
	repo repository.Model
	pr   platform.Repository
}

func (s modelService) Create(cmd *ModelCreateCmd) (dto ModelDTO, err error) {
	pid, err := s.pr.New(platform.RepoOption{
		Name:     cmd.Name,
		Desc:     cmd.Desc,
		RepoType: cmd.RepoType,
	})
	if err != nil {
		return
	}

	v := cmd.toModel()
	v.RepoId = pid

	m, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	s.toModelDTO(&m, &dto)

	return
}

func (s modelService) GetByName(
	owner domain.Account, name domain.ModelName,
) (dto ModelDTO, err error) {
	v, err := s.repo.GetByName(owner, name)
	if err != nil {
		return
	}

	s.toModelDTO(&v, &dto)

	return
}

func (s modelService) List(owner domain.Account, cmd *ResourceListCmd) (
	dtos []ModelDTO, err error,
) {
	v, err := s.repo.List(owner, cmd.toResourceListOption())
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]ModelDTO, len(v))
	for i := range v {
		s.toModelDTO(&v[i], &dtos[i])
	}

	return
}

func (s modelService) toModelDTO(m *domain.Model, dto *ModelDTO) {
	*dto = ModelDTO{
		Id:       m.Id,
		Owner:    m.Owner.Account(),
		Name:     m.Name.ModelName(),
		Protocol: m.Protocol.ProtocolName(),
		RepoType: m.RepoType.RepoType(),
		RepoId:   m.RepoId,
		Tags:     m.Tags,
	}

	if m.Desc != nil {
		dto.Desc = m.Desc.ProjDesc()
	}
}
