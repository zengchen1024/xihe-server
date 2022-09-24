package app

import (
	"errors"
	"fmt"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
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

func (s resourceService) list(resources []*domain.ResourceObject) (
	dtos []ResourceDTO, err error,
) {
	users, projects, datasets, models := s.toOptions(resources)

	return s.listResources(users, projects, datasets, models, len(resources))
}

func (s resourceService) listModels(resources []domain.ResourceIndex) (
	dtos []ResourceDTO, err error,
) {
	if len(resources) == 0 {
		return
	}

	users, options := s.singleResourceOptions(resources)

	return s.listResources(users, nil, nil, options, len(resources))
}

func (s resourceService) listDatasets(resources []domain.ResourceIndex) (
	dtos []ResourceDTO, err error,
) {
	if len(resources) == 0 {
		return
	}

	users, options := s.singleResourceOptions(resources)

	return s.listResources(users, nil, options, nil, len(resources))
}

func (s resourceService) singleResourceOptions(resources []domain.ResourceIndex) (
	users []domain.Account,
	options []repository.UserResourceListOption,
) {
	ul := make(map[string]domain.Account)
	ro := make(map[string][]string)

	for i := range resources {
		item := resources[i]

		account := item.Owner.Account()
		if _, ok := ul[account]; !ok {
			ul[account] = item.Owner
		}

		s.store(&item, ro)
	}

	options = s.toOptionList(ro, ul)

	users = s.userMapToList(ul)

	return
}

func (s resourceService) toOptions(resources []*domain.ResourceObject) (
	users []domain.Account,
	projects []repository.UserResourceListOption,
	datasets []repository.UserResourceListOption,
	models []repository.UserResourceListOption,
) {
	ul := make(map[string]domain.Account)
	po := make(map[string][]string)
	do := make(map[string][]string)
	mo := make(map[string][]string)

	for i := range resources {
		item := resources[i]

		account := item.Owner.Account()
		if _, ok := ul[account]; !ok {
			ul[account] = item.Owner
		}

		switch item.Type.ResourceType() {
		case domain.ResourceProject:
			s.store(&item.ResourceIndex, po)

		case domain.ResourceModel:
			s.store(&item.ResourceIndex, mo)

		case domain.ResourceDataset:
			s.store(&item.ResourceIndex, do)
		}
	}

	projects = s.toOptionList(po, ul)
	datasets = s.toOptionList(do, ul)
	models = s.toOptionList(mo, ul)
	users = s.userMapToList(ul)

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
			Id:        p.Id,
			Name:      p.Name.ProjName(),
			Type:      domain.ResourceProject,
			Desc:      p.Desc.ResourceDesc(),
			CoverId:   p.CoverId.CoverId(),
			UpdateAt:  utils.ToDate(p.UpdatedAt),
			LikeCount: p.LikeCount,
			//DownloadCount
			ForkCount: p.ForkCount,
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
			Id:   d.Id,
			Name: d.Name.ModelName(),
			Type: domain.ResourceModel,
			Desc: d.Desc.ResourceDesc(),
			/*
				UpdateAt

				LikeCount
				DownloadCount
				ForkCount
			*/
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
			Id:   d.Id,
			Name: d.Name.DatasetName(),
			Type: domain.ResourceDataset,
			Desc: d.Desc.ResourceDesc(),
			/*
				UpdateAt

				LikeCount
				DownloadCount
				ForkCount
			*/
		}

		if u, ok := userInfos[d.Owner.Account()]; ok {
			v.Owner.Name = u.Account.Account()
			v.Owner.AvatarId = u.AvatarId.AvatarId()
		}

		dtos[i] = v
	}
}

func (s resourceService) store(v *domain.ResourceIndex, m map[string][]string) {
	a := v.Owner.Account()

	if p, ok := m[a]; !ok {
		m[a] = []string{v.Id}
	} else {
		m[a] = append(p, v.Id)
	}
}

func (s resourceService) toOptionList(
	m map[string][]string, users map[string]domain.Account,
) []repository.UserResourceListOption {

	if len(m) == 0 {
		return nil
	}

	r := make([]repository.UserResourceListOption, len(m))

	i := 0
	for k, v := range m {
		r[i] = repository.UserResourceListOption{
			Owner: users[k],
			Ids:   v,
		}

		i++
	}

	return r
}

func (s resourceService) userMapToList(m map[string]domain.Account) []domain.Account {
	if len(m) == 0 {
		return nil
	}

	r := make([]domain.Account, len(m))

	i := 0
	for _, u := range m {
		r[i] = u
		i++
	}

	return r
}
