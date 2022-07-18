package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ModelCreateCmd struct {
	Owner    string
	Name     domain.ProjName
	Desc     domain.ProjDesc
	RepoType domain.RepoType
	Protocol domain.ProtocolName
}

func (cmd *ModelCreateCmd) Validate() error {
	b := cmd.Owner != "" &&
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
		Name:     cmd.Name,
		Desc:     cmd.Desc,
		RepoType: cmd.RepoType,
		Protocol: cmd.Protocol,
	}
}

type ModelDTO struct {
	Id       string   `json:"id"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Protocol string   `json:"protocol"`
	RepoType string   `json:"repo_type"`
	Tags     []string `json:"tags"`
}

type ModelService interface {
	Create(*ModelCreateCmd) (ModelDTO, error)
	Get(string, string) (ModelDTO, error)
}

func NewModelService(repo repository.Model) ModelService {
	return modelService{repo}
}

type modelService struct {
	repo repository.Model
}

func (s modelService) Create(cmd *ModelCreateCmd) (dto ModelDTO, err error) {
	m := cmd.toModel()

	v, err := s.repo.Save(&m)
	if err != nil {
		return
	}

	s.toModelDTO(&v, &dto)

	// TODO send event

	return
}

func (s modelService) toModelDTO(m *domain.Model, dto *ModelDTO) {
	*dto = ModelDTO{
		Id:       m.Id,
		Owner:    m.Owner,
		Name:     m.Name.ProjName(),
		Desc:     m.Desc.ProjDesc(),
		Protocol: m.Protocol.ProtocolName(),
		RepoType: m.RepoType.RepoType(),
		Tags:     m.Tags,
	}
}

func (s modelService) Get(owner, modelId string) (dto ModelDTO, err error) {
	v, err := s.repo.Get(owner, modelId)
	if err != nil {
		return
	}

	s.toModelDTO(&v, &dto)

	return
}
