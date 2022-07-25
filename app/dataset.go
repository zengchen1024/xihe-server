package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
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
	RepoId   string   `json:"repo_id"`
	Tags     []string `json:"tags"`
}

type DatasetService interface {
	Create(*DatasetCreateCmd) (DatasetDTO, error)
	Get(domain.Account, string) (DatasetDTO, error)
	List(domain.Account, *DatasetListCmd) ([]DatasetDTO, error)
}

func NewDatasetService(repo repository.Dataset, pr platform.Repository) DatasetService {
	return datasetService{repo: repo, pr: pr}
}

type datasetService struct {
	repo repository.Dataset
	pr   platform.Repository
}

func (s datasetService) Create(cmd *DatasetCreateCmd) (dto DatasetDTO, err error) {
	v := cmd.toDataset()

	d, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	pid, err := s.pr.New(platform.RepoOption{
		Name: cmd.Name,
		Desc: cmd.Desc,
	})
	if err != nil {
		return
	}

	d.RepoId = pid

	d, err = s.repo.Save(&d)
	if err != nil {
		return
	}

	s.toDatasetDTO(&d, &dto)

	return
}

func (s datasetService) Get(owner domain.Account, datasetId string) (dto DatasetDTO, err error) {
	v, err := s.repo.Get(owner, datasetId)
	if err != nil {
		return
	}

	s.toDatasetDTO(&v, &dto)

	return
}

type DatasetListCmd struct {
	Name domain.ProjName
}

func (cmd *DatasetListCmd) toDatasetListOption() (
	option repository.DatasetListOption,
) {
	option.Name = cmd.Name

	return
}

func (s datasetService) List(owner domain.Account, cmd *DatasetListCmd) (
	dtos []DatasetDTO, err error,
) {
	v, err := s.repo.List(owner, cmd.toDatasetListOption())
	if err != nil || len(v) == 0 {
		return
	}

	dtos = make([]DatasetDTO, len(v))
	for i := range v {
		s.toDatasetDTO(&v[i], &dtos[i])
	}

	return
}

func (s datasetService) toDatasetDTO(d *domain.Dataset, dto *DatasetDTO) {
	*dto = DatasetDTO{
		Id:       d.Id,
		Owner:    d.Owner.Account(),
		Name:     d.Name.ProjName(),
		Protocol: d.Protocol.ProtocolName(),
		RepoType: d.RepoType.RepoType(),
		RepoId:   d.RepoId,
		Tags:     d.Tags,
	}

	if d.Desc != nil {
		dto.Desc = d.Desc.ProjDesc()
	}
}
