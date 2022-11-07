package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type SearchDTO struct {
	Model   ResourceSearchDTO `json:"model"`
	Project ResourceSearchDTO `json:"project"`
	Dataset ResourceSearchDTO `json:"dataset"`
}

type ResourceSearchDTO struct {
	Top   []ResourceSummaryDTO `json:"top"`
	Total int                  `json:"total"`
}

type ResourceSummaryDTO struct {
	Owner string `json:"owner"`
	Name  string `json:"name"`
}

func newResourceSearchOption(name string) repository.ResourceSearchOption {
	option := repository.ResourceSearchOption{
		Name:   name,
		TopNum: 3,
	}
	option.RepoType, _ = domain.NewRepoType(domain.RepoTypePublic)

	return option
}

type SearchService interface {
	Search(name string) (dto SearchDTO)
}

func NewSearchService(
	model repository.Model,
	project repository.Project,
	dataset repository.Dataset,
) SearchService {
	return searchService{
		model:   model,
		project: project,
		dataset: dataset,
	}
}

type searchService struct {
	model   repository.Model
	project repository.Project
	dataset repository.Dataset
}

func (s searchService) Search(name string) (dto SearchDTO) {
	option := newResourceSearchOption(name)

	v, err := s.search(&option, s.project.Search)
	if err == nil {
		dto.Project = v
	}

	v, err = s.search(&option, s.model.Search)
	if err == nil {
		dto.Model = v
	}

	v, err = s.search(&option, s.dataset.Search)
	if err == nil {
		dto.Dataset = v
	}

	return
}

func (s searchService) search(
	option *repository.ResourceSearchOption,
	f func(*repository.ResourceSearchOption) (repository.ResourceSearchResult, error),
) (
	dto ResourceSearchDTO, err error,
) {
	v, err := f(option)
	if err != nil {
		return
	}

	items := make([]ResourceSummaryDTO, len(v.Top))
	for i := range v.Top {
		item := &v.Top[i]

		items[i].Owner = item.Owner.Account()
		items[i].Name = item.Name.ResourceName()
	}

	dto.Top = items
	dto.Total = v.Total

	return
}
