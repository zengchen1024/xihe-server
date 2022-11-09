package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/platform"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
)

type ModelCreateCmd struct {
	Owner    domain.Account
	Name     domain.ResourceName
	Desc     domain.ResourceDesc
	RepoType domain.RepoType
	Protocol domain.ProtocolName
}

func (cmd *ModelCreateCmd) Validate() error {
	b := cmd.Owner != nil &&
		cmd.Name != nil &&
		cmd.Desc != nil &&
		cmd.RepoType != nil &&
		cmd.Protocol != nil

	if !b {
		return errors.New("invalid cmd of creating model")
	}

	return nil
}

func (cmd *ModelCreateCmd) toModel() domain.Model {
	now := utils.Now()

	return domain.Model{
		Owner:     cmd.Owner,
		Protocol:  cmd.Protocol,
		CreatedAt: now,
		UpdatedAt: now,
		ModelModifiableProperty: domain.ModelModifiableProperty{
			Name:     cmd.Name,
			Desc:     cmd.Desc,
			RepoType: cmd.RepoType,
			Tags:     []string{cmd.Protocol.ProtocolName()},
			TagKinds: []string{},
		},
	}
}

type ModelsDTO struct {
	Total  int               `json:"total"`
	Models []ModelSummaryDTO `json:"models"`
}

type ModelSummaryDTO struct {
	Id            string   `json:"id"`
	Owner         string   `json:"owner"`
	Name          string   `json:"name"`
	Desc          string   `json:"desc"`
	Tags          []string `json:"tags"`
	UpdatedAt     string   `json:"updated_at"`
	LikeCount     int      `json:"like_count"`
	DownloadCount int      `json:"download_count"`
}

type ModelDTO struct {
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

type ModelDetailDTO struct {
	ModelDTO

	RelatedDatasets []ResourceDTO `json:"related_datasets"`
	RelatedProjects []ResourceDTO `json:"related_projects"`
}

type ModelService interface {
	Create(*ModelCreateCmd, platform.Repository) (ModelDTO, error)
	Update(*domain.Model, *ModelUpdateCmd, platform.Repository) (ModelDTO, error)
	GetByName(domain.Account, domain.ResourceName, bool) (ModelDetailDTO, error)
	List(domain.Account, *ResourceListCmd) (ModelsDTO, error)
	ListGlobal(*GlobalResourceListCmd) (GlobalModelsDTO, error)

	AddLike(*domain.ResourceIndex) error
	RemoveLike(*domain.ResourceIndex) error

	AddRelatedDataset(*domain.Model, *domain.ResourceIndex) error
	RemoveRelatedDataset(*domain.Model, *domain.ResourceIndex) error

	AddRelatedProject(*domain.ReverselyRelatedResourceInfo) error
	RemoveRelatedProject(*domain.ReverselyRelatedResourceInfo) error

	SetTags(*domain.Model, *ResourceTagsUpdateCmd) error
}

func NewModelService(
	user repository.User,
	repo repository.Model,
	proj repository.Project,
	dataset repository.Dataset,
	activity repository.Activity,
	pr platform.Repository,
	sender message.Sender,
) ModelService {
	return modelService{
		repo:     repo,
		activity: activity,
		sender:   sender,
		rs: resourceService{
			user:    user,
			model:   repo,
			project: proj,
			dataset: dataset,
		},
	}
}

type modelService struct {
	repo repository.Model
	//pr       platform.Repository
	activity repository.Activity
	rs       resourceService
	sender   message.Sender
}

func (s modelService) Create(cmd *ModelCreateCmd, pr platform.Repository) (dto ModelDTO, err error) {
	pid, err := pr.New(&platform.RepoOption{
		Name:     cmd.Name,
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

	// add activity
	ua := genActivityForCreatingResource(
		m.Owner, domain.ResourceTypeModel, m.Id,
	)
	// ignore the error
	_ = s.activity.Save(&ua)

	return
}

func (s modelService) GetByName(
	owner domain.Account, name domain.ResourceName,
	allowPrivacy bool,
) (dto ModelDetailDTO, err error) {
	v, err := s.repo.GetByName(owner, name)
	if err != nil {
		return
	}

	if !allowPrivacy && v.IsPrivate() {
		err = ErrorPrivateRepo{errors.New("private repo")}

		return
	}

	d, err := s.rs.listDatasets(v.RelatedDatasets)
	if err != nil {
		return
	}
	dto.RelatedDatasets = d

	d, err = s.rs.listProjects(v.RelatedProjects)
	if err != nil {
		return
	}
	dto.RelatedProjects = d

	s.toModelDTO(&v, &dto.ModelDTO)

	return
}

func (s modelService) List(owner domain.Account, cmd *ResourceListCmd) (
	dto ModelsDTO, err error,
) {
	option := cmd.toResourceListOption()

	var v repository.UserModelsInfo

	if cmd.SortType == nil {
		v, err = s.repo.ListAndSortByUpdateTime(owner, &option)
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

	items := v.Models

	if err != nil || len(items) == 0 {
		return
	}

	dtos := make([]ModelSummaryDTO, len(items))
	for i := range items {
		s.toModelSummaryDTO(&items[i], &dtos[i])
	}

	dto.Total = v.Total
	dto.Models = dtos

	return
}

func (s modelService) toModelDTO(m *domain.Model, dto *ModelDTO) {
	*dto = ModelDTO{
		Id:            m.Id,
		Owner:         m.Owner.Account(),
		Name:          m.Name.ResourceName(),
		Desc:          m.Desc.ResourceDesc(),
		Protocol:      m.Protocol.ProtocolName(),
		RepoType:      m.RepoType.RepoType(),
		RepoId:        m.RepoId,
		Tags:          m.Tags,
		CreatedAt:     utils.ToDate(m.CreatedAt),
		UpdatedAt:     utils.ToDate(m.UpdatedAt),
		LikeCount:     m.LikeCount,
		DownloadCount: m.DownloadCount,
	}
}

func (s modelService) toModelSummaryDTO(p *domain.ModelSummary, dto *ModelSummaryDTO) {
	*dto = ModelSummaryDTO{
		Id:            p.Id,
		Owner:         p.Owner.Account(),
		Name:          p.Name.ResourceName(),
		Desc:          p.Desc.ResourceDesc(),
		Tags:          p.Tags,
		UpdatedAt:     utils.ToDate(p.UpdatedAt),
		LikeCount:     p.LikeCount,
		DownloadCount: p.DownloadCount,
	}
}
