package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type DatasetCreateCmd struct {
	Owner    domain.Account
	Name     domain.ProjName
	Desc     domain.ProjDesc
	RepoType domain.RepoType
	Protocol domain.ProtocolName
}

func (cmd *DatasetCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Name != nil &&
		cmd.RepoType != nil &&
		cmd.Protocol != nil

	if !b {
		return errors.New("invalid cmd of creating dataset")
	}

	return nil
}

func (cmd *DatasetCreateCmd) toDataset() domain.Dataset {
	return domain.Dataset{
		Owner:    cmd.Owner,
		Name:     cmd.Name,
		Desc:     cmd.Desc,
		RepoType: cmd.RepoType,
		Protocol: cmd.Protocol,
	}
}

type DatasetDTO struct {
	Id       string   `json:"id"`
	Owner    string   `json:"owner"`
	Name     string   `json:"name"`
	Desc     string   `json:"desc"`
	Protocol string   `json:"protocol"`
	RepoType string   `json:"repo_type"`
	Tags     []string `json:"tags"`
}

type DatasetService interface {
	Create(*DatasetCreateCmd) (DatasetDTO, error)
	Get(domain.Account, string) (DatasetDTO, error)
}

func NewDatasetService(repo repository.Dataset) DatasetService {
	return datasetService{repo}
}

type datasetService struct {
	repo repository.Dataset
}

func (s datasetService) Create(cmd *DatasetCreateCmd) (dto DatasetDTO, err error) {
	m := cmd.toDataset()

	v, err := s.repo.Save(&m)
	if err != nil {
		return
	}

	s.toDatasetDTO(&v, &dto)

	// TODO send event

	return
}

func (s datasetService) toDatasetDTO(m *domain.Dataset, dto *DatasetDTO) {
	*dto = DatasetDTO{
		Id:       m.Id,
		Owner:    m.Owner.Account(),
		Name:     m.Name.ProjName(),
		Desc:     m.Desc.ProjDesc(),
		Protocol: m.Protocol.ProtocolName(),
		RepoType: m.RepoType.RepoType(),
		Tags:     m.Tags,
	}
}

func (s datasetService) Get(owner domain.Account, datasetId string) (dto DatasetDTO, err error) {
	v, err := s.repo.Get(owner, datasetId)
	if err != nil {
		return
	}

	s.toDatasetDTO(&v, &dto)

	return
}
