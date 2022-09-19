package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type DatasetCreateCmd struct {
	Owner    domain.Account
	Name     domain.DatasetName
	Desc     domain.ResourceDesc
	RepoType domain.RepoType
	Protocol domain.ProtocolName
}

func (cmd *DatasetCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Name != nil &&
		cmd.Desc != nil &&
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
		Protocol: cmd.Protocol,
		DatasetModifiableProperty: domain.DatasetModifiableProperty{
			Name:     cmd.Name,
			Desc:     cmd.Desc,
			RepoType: cmd.RepoType,
		},
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
	Create(*DatasetCreateCmd, platform.Repository) (DatasetDTO, error)
	Update(*domain.Dataset, *DatasetUpdateCmd, platform.Repository) (DatasetDTO, error)
	GetByName(domain.Account, domain.DatasetName) (DatasetDTO, error)
	List(domain.Account, *ResourceListCmd) ([]DatasetDTO, error)

	AddLike(domain.Account, string) error
	RemoveLike(domain.Account, string) error
}

func NewDatasetService(
	repo repository.Dataset, activity repository.Activity, pr platform.Repository,
) DatasetService {
	return datasetService{repo: repo, activity: activity}
}

type datasetService struct {
	repo repository.Dataset
	//pr       platform.Repository
	activity repository.Activity
}

func (s datasetService) Create(cmd *DatasetCreateCmd, pr platform.Repository) (dto DatasetDTO, err error) {
	pid, err := pr.New(&platform.RepoOption{
		Name:     cmd.Name,
		RepoType: cmd.RepoType,
	})
	if err != nil {
		return
	}

	v := cmd.toDataset()
	v.RepoId = pid

	d, err := s.repo.Save(&v)
	if err != nil {
		return
	}

	s.toDatasetDTO(&d, &dto)

	// add activity
	ua := genActivityForCreatingResource(
		d.Owner, domain.ResourceTypeDataset, d.Id,
	)
	// ignore the error
	_ = s.activity.Save(&ua)

	return
}

func (s datasetService) GetByName(
	owner domain.Account, name domain.DatasetName,
) (dto DatasetDTO, err error) {
	v, err := s.repo.GetByName(owner, name)
	if err != nil {
		return
	}

	s.toDatasetDTO(&v, &dto)

	return
}

func (s datasetService) List(owner domain.Account, cmd *ResourceListCmd) (
	dtos []DatasetDTO, err error,
) {
	v, err := s.repo.List(owner, cmd.toResourceListOption())
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
		Name:     d.Name.DatasetName(),
		Desc:     d.Desc.ResourceDesc(),
		Protocol: d.Protocol.ProtocolName(),
		RepoType: d.RepoType.RepoType(),
		RepoId:   d.RepoId,
		Tags:     d.Tags,
	}
}
