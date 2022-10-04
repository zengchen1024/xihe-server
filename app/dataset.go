package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
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
	now := utils.Now()

	return domain.Dataset{
		Owner:     cmd.Owner,
		Protocol:  cmd.Protocol,
		CreatedAt: now,
		UpdatedAt: now,
		DatasetModifiableProperty: domain.DatasetModifiableProperty{
			Name:     cmd.Name,
			Desc:     cmd.Desc,
			RepoType: cmd.RepoType,
			Tags:     []string{cmd.Protocol.ProtocolName()},
		},
	}
}

type DatasetsDTO struct {
	Total    int                 `json:"total"`
	Datasets []DatasetSummaryDTO `json:"models"`
}

type DatasetSummaryDTO struct {
	Id            string   `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Tags          []string `json:"tags"`
	UpdatedAt     string   `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	DownloadCount int      `json:"download_count"`
}

type DatasetDTO struct {
	Id            string   `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Protocol      string   `json:"protocol"`
	RepoType      string   `json:"repo_type"`
	RepoId        string   `json:"repo_id"`
	Tags          []string `json:"tags"`
	CreatedAt     string   `json:"created_at"`
	UpdatedAt     string   `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	DownloadCount int      `json:"download_count"`
}

type DatasetDetailDTO struct {
	DatasetDTO
}

type DatasetService interface {
	Create(*DatasetCreateCmd, platform.Repository) (DatasetDTO, error)
	Update(*domain.Dataset, *DatasetUpdateCmd, platform.Repository) (DatasetDTO, error)
	GetByName(domain.Account, domain.DatasetName, bool) (DatasetDetailDTO, error)
	List(domain.Account, *ResourceListCmd) (DatasetsDTO, error)

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	SetTags(*domain.Dataset, *ResourceTagsUpdateCmd) error
}

func NewDatasetService(
	user repository.User,
	repo repository.Dataset,
	proj repository.Project,
	model repository.Model,
	activity repository.Activity,
	pr platform.Repository,
) DatasetService {
	return datasetService{
		repo:     repo,
		activity: activity,
		rs: resourceService{
			user:    user,
			model:   model,
			project: proj,
			dataset: repo,
		},
	}
}

type datasetService struct {
	repo repository.Dataset
	//pr       platform.Repository
	activity repository.Activity
	rs       resourceService
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
	allowPrivacy bool,
) (dto DatasetDetailDTO, err error) {
	v, err := s.repo.GetByName(owner, name)
	if err != nil {
		return
	}

	if !allowPrivacy && v.IsPrivate() {
		err = ErrorPrivateRepo{errors.New("private repo")}

		return
	}

	s.toDatasetDTO(&v, &dto.DatasetDTO)

	return
}

func (s datasetService) List(owner domain.Account, cmd *ResourceListCmd) (
	dto DatasetsDTO, err error,
) {
	option := cmd.toResourceListOption()

	var v repository.UserDatasetsInfo

	if cmd.SortType == nil {
		v, err = s.repo.List(owner, &option)
	} else {
		switch cmd.SortType.SortType() {
		case domain.SortTypeUpdateTime:
			v, err = s.repo.ListAndSortByUpdateTime(owner, &option)

		case domain.SortTypeFirstLetter:
			v, err = s.repo.ListAndSortByFirstLetter(owner, &option)

		case domain.SortTypeDownloadCount:
			v, err = s.repo.ListAndSortByDownloadCount(owner, &option)
		}
	}

	items := v.Datasets

	if err != nil || len(items) == 0 {
		return
	}

	dtos := make([]DatasetSummaryDTO, len(items))
	for i := range items {
		s.toDatasetSummaryDTO(&items[i], &dtos[i])
	}

	dto.Total = v.Total
	dto.Datasets = dtos

	return
}

func (s datasetService) toDatasetDTO(d *domain.Dataset, dto *DatasetDTO) {
	*dto = DatasetDTO{
		Id:            d.Id,
		Owner:         d.Owner.Account(),
		Name:          d.Name.DatasetName(),
		Desc:          d.Desc.ResourceDesc(),
		Protocol:      d.Protocol.ProtocolName(),
		RepoType:      d.RepoType.RepoType(),
		RepoId:        d.RepoId,
		Tags:          d.Tags,
		CreatedAt:     utils.ToDate(d.CreatedAt),
		UpdatedAt:     utils.ToDate(d.UpdatedAt),
		LikeCount:     d.LikeCount,
		DownloadCount: d.DownloadCount,
	}
}

func (s datasetService) toDatasetSummaryDTO(p *domain.DatasetSummary, dto *DatasetSummaryDTO) {
	*dto = DatasetSummaryDTO{
		Id:            p.Id,
		Owner:         p.Owner.Account(),
		Name:          p.Name.DatasetName(),
		Desc:          p.Desc.ResourceDesc(),
		Tags:          p.Tags,
		UpdatedAt:     utils.ToDate(p.UpdatedAt),
		LikeCount:     p.LikeCount,
		DownloadCount: p.DownloadCount,
	}
}
