package app

import (
	"errors"
	"fmt"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ResourceDTO struct {
	Owner struct {
		Name     string `json:"name"`
		AvatarId string `json:"avatar_id"`
	} `json:"owner"`

	Id       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Desc     string `json:"description"`
	CoverId  string `json:"cover_id"`
	UpdateAt string `json:"update_at"`

	LikeCount     int `json:"like_count"`
	DownloadCount int `json:"download_count"`
	ForkCount     int `json:"fork_count"`
}

func (r *ResourceDTO) identity() string {
	return fmt.Sprintf("%s_%s_%s", r.Owner.Name, r.Type, r.Id)
}

type resourceService struct {
	user    repository.User
	model   repository.Model
	project repository.Project
	dataset repository.Dataset
}

func (s resourceService) list(resources []*domain.ResourceObj) (
	dtos []ResourceDTO, err error,
) {
	users, projects, datasets, models := s.toOptions(resources)

	return s.listResources(users, projects, datasets, models, len(resources))
}

func (s resourceService) toOptions(resources []*domain.ResourceObj) (
	users []domain.Account,
	projects []repository.UserResourceListOption,
	datasets []repository.UserResourceListOption,
	models []repository.UserResourceListOption,
) {
	users1 := map[string]domain.Account{}
	projects1 := map[string]*repository.UserResourceListOption{}
	datasets1 := map[string]*repository.UserResourceListOption{}
	models1 := map[string]*repository.UserResourceListOption{}

	set := func(v *domain.ResourceObj, m map[string]*repository.UserResourceListOption) {
		a := v.ResourceOwner.Account()
		if p, ok := m[a]; !ok {
			m[a] = &repository.UserResourceListOption{
				Owner: v.ResourceOwner,
				Ids:   []string{v.ResourceId},
			}
		} else {
			p.Ids = append(p.Ids, v.ResourceId)
		}
	}

	for i := range resources {
		item := resources[i]

		account := item.ResourceOwner.Account()
		if _, ok := users1[account]; !ok {
			users1[account] = item.ResourceOwner
		}

		switch item.ResourceType.ResourceType() {
		case domain.ResourceProject:
			set(item, projects1)

		case domain.ResourceModel:
			set(item, models1)

		case domain.ResourceDataset:
			set(item, datasets1)
		}
	}

	toList := func(m map[string]*repository.UserResourceListOption) []repository.UserResourceListOption {
		n := len(m)
		r := make([]repository.UserResourceListOption, n)

		i := 0
		for _, v := range m {
			r[i] = *v
			i++
		}

		return r
	}

	projects = toList(projects1)
	datasets = toList(datasets1)
	models = toList(models1)

	users = make([]domain.Account, len(users1))
	i := 0
	for _, u := range users1 {
		users[i] = u
		i++
	}

	return
}

func (s resourceService) listResources(
	users []domain.Account,
	projects []repository.UserResourceListOption,
	datasets []repository.UserResourceListOption,
	models []repository.UserResourceListOption,
	total int,
) (
	dtos []ResourceDTO, err error,
) {
	allUsers, err := s.user.FindUsersInfo(users)
	if err != nil {
		return
	}

	userInfos := make(map[string]*domain.UserInfo)
	for i := range allUsers {
		item := &allUsers[i]
		userInfos[item.Account.Account()] = item
	}

	dtos = make([]ResourceDTO, total)
	total = 0
	r := dtos

	if len(projects) > 0 {
		all, err := s.project.FindUserProjects(projects)
		if err != nil {
			return nil, err
		}

		n := len(all)
		if n > 0 {
			if len(r) < n {
				return nil, errors.New("unmatched size")
			}
			s.projectToResourceDTO(userInfos, all, r)
			r = r[n:]
			total += n
		}
	}

	if len(models) > 0 {
		all, err := s.model.FindUserModels(models)
		if err != nil {
			return nil, err
		}

		n := len(all)
		if n > 0 {
			if len(r) < n {
				return nil, errors.New("unmatched size")
			}
			s.modelToResourceDTO(userInfos, all, r)
			r = r[n:]
			total += n
		}
	}

	if n := len(datasets); n > 0 {
		all, err := s.dataset.FindUserDatasets(datasets)
		if err != nil {
			return nil, err
		}

		n := len(all)
		if n > 0 {
			if len(r) < n {
				return nil, errors.New("unmatched size")
			}
			s.datasetToResourceDTO(userInfos, all, r)
			total += n
		}
	}

	dtos = dtos[:total]

	return
}

func (s resourceService) projectToResourceDTO(
	userInfos map[string]*domain.UserInfo,
	projects []domain.Project, dtos []ResourceDTO,
) {
	for i := range projects {
		p := &projects[i]

		v := ResourceDTO{
			Name:    p.Name.ProjName(),
			CoverId: p.CoverId.CoverId(),
			/*
				UpdateAt

				LikeCount
				DownloadCount
				ForkCount
			*/
		}

		if p.Desc != nil {
			v.Desc = p.Desc.ResourceDesc()
		}

		if u, ok := userInfos[p.Owner.Account()]; ok {
			v.Owner.Name = u.Account.Account()
			v.Owner.AvatarId = u.AvatarId.AvatarId()
		}

		dtos[i] = v
	}
}

func (s resourceService) modelToResourceDTO(
	userInfos map[string]*domain.UserInfo,
	data []domain.Model, dtos []ResourceDTO,
) {
	for i := range data {
		d := &data[i]

		v := ResourceDTO{
			Name: d.Name.ModelName(),
			/*
				UpdateAt

				LikeCount
				DownloadCount
				ForkCount
			*/
		}

		if d.Desc != nil {
			v.Desc = d.Desc.ResourceDesc()
		}

		if u, ok := userInfos[d.Owner.Account()]; ok {
			v.Owner.Name = u.Account.Account()
			v.Owner.AvatarId = u.AvatarId.AvatarId()
		}

		dtos[i] = v
	}
}

func (s resourceService) datasetToResourceDTO(
	userInfos map[string]*domain.UserInfo,
	data []domain.Dataset, dtos []ResourceDTO,
) {
	for i := range data {
		d := &data[i]

		v := ResourceDTO{
			Name: d.Name.DatasetName(),
			/*
				UpdateAt

				LikeCount
				DownloadCount
				ForkCount
			*/
		}

		if d.Desc != nil {
			v.Desc = d.Desc.ResourceDesc()
		}

		if u, ok := userInfos[d.Owner.Account()]; ok {
			v.Owner.Name = u.Account.Account()
			v.Owner.AvatarId = u.AvatarId.AvatarId()
		}

		dtos[i] = v
	}
}
