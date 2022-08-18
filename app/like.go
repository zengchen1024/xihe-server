package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type LikeCreateCmd struct {
	ResourceOwner domain.Account
	ResourceType  domain.ResourceType
	ResourceId    string
}

func (cmd *LikeCreateCmd) Validate() error {
	b := cmd.ResourceOwner != nil &&
		cmd.ResourceType != nil &&
		cmd.ResourceId != ""

	if !b {
		return errors.New("invalid cmd")
	}

	return nil
}

type LikeRemoveCmd = LikeCreateCmd

type LikeDTO struct {
	Owner struct {
		Name     string `json:"name"`
		AvatarId string `json:"avatar_id"`
	} `json:"owner"`

	Name     string `json:"name"`
	Desc     string `json:"description"`
	CoverId  string `json:"cover_id"`
	UpdateAt string `json:"update_at"`

	LikeCount     int `json:"like_count"`
	DownloadCount int `json:"download_count"`
	ForkCount     int `json:"fork_count"`
}

type LikeService interface {
	Create(domain.Account, LikeCreateCmd) error
	Delete(domain.Account, LikeRemoveCmd) error
	List(domain.Account) ([]LikeDTO, error)
}

func NewLikeService(
	repo repository.Like,
	user repository.User,
	model repository.Model,
	project repository.Project,
	dataset repository.Dataset,
	sender message.Sender,
) LikeService {
	return likeService{
		repo:    repo,
		user:    user,
		model:   model,
		project: project,
		dataset: dataset,
		sender:  sender,
	}
}

type likeService struct {
	repo    repository.Like
	user    repository.User
	model   repository.Model
	project repository.Project
	dataset repository.Dataset
	sender  message.Sender
}

func (s likeService) Create(owner domain.Account, cmd LikeCreateCmd) error {
	v := domain.UserLike{
		Owner: owner,
		Like: domain.Like{
			ResourceOwner: cmd.ResourceOwner,
			ResourceType:  cmd.ResourceType,
			ResourceId:    cmd.ResourceId,
		},
	}

	if err := s.repo.Save(&v); err != nil {
		return err
	}

	// TODO: activity

	// send event
	return s.sender.AddLike(v.Like)
}

func (s likeService) Delete(owner domain.Account, cmd LikeRemoveCmd) error {
	v := domain.UserLike{
		Owner: owner,
		Like: domain.Like{
			ResourceOwner: cmd.ResourceOwner,
			ResourceType:  cmd.ResourceType,
			ResourceId:    cmd.ResourceId,
		},
	}

	if err := s.repo.Remove(&v); err != nil {
		return err
	}

	// send event
	return s.sender.RemoveLike(v.Like)
}

func (s likeService) List(owner domain.Account) (
	dtos []LikeDTO, err error,
) {
	v, err := s.repo.Find(owner, repository.LikeFindOption{})
	if err != nil || len(v) == 0 {
		return
	}

	users, projects, datasets, models := s.toOptions(v)

	return s.listLike(users, projects, datasets, models, len(v))
}

func (s likeService) toOptions(likes []domain.Like) (
	users []domain.Account,
	projects []repository.UserResourceListOption,
	datasets []repository.UserResourceListOption,
	models []repository.UserResourceListOption,
) {
	users1 := map[string]domain.Account{}
	projects1 := map[string]*repository.UserResourceListOption{}
	datasets1 := map[string]*repository.UserResourceListOption{}
	models1 := map[string]*repository.UserResourceListOption{}

	set := func(v *domain.Like, m map[string]*repository.UserResourceListOption) {
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

	for i := range likes {
		item := &likes[i]

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

func (s likeService) listLike(
	users []domain.Account,
	projects []repository.UserResourceListOption,
	datasets []repository.UserResourceListOption,
	models []repository.UserResourceListOption,
	total int,
) (
	dtos []LikeDTO, err error,
) {
	allUsers, err := s.user.FindUsers(users)
	if err != nil {
		return
	}

	userInfos := make(map[string]*domain.UserInfo)
	for i := range allUsers {
		item := &allUsers[i]
		userInfos[item.Account.Account()] = item
	}

	dtos = make([]LikeDTO, total)
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
			s.projectToLikeDTO(userInfos, all, r)
			r = r[len(all):]
		}
	}

	if len(models) > 0 {
		all, err := s.model.FindUserModels(models)
		if err != nil {
			return nil, err
		}

		s.modelToLikeDTO(userInfos, all, r)
		r = r[len(all):]
	}

	if n := len(datasets); n > 0 {
		all, err := s.dataset.FindUserDatasets(datasets)
		if err != nil {
			return nil, err
		}

		s.datasetToLikeDTO(userInfos, all, r)
	}

	return
}

func (s likeService) projectToLikeDTO(
	userInfos map[string]*domain.UserInfo,
	projects []domain.Project, dtos []LikeDTO,
) {
	for i := range projects {
		p := &projects[i]

		v := LikeDTO{
			Name:    p.Name.ProjName(),
			CoverId: p.CoverId.CoverId(),
			/*
				UpdateAt string `json:"update_at"`

				LikeCount     int `json:"like_count"`
				DownloadCount int `json:"download_count"`
				ForkCount     int `json:"fork_count"`
			*/
		}

		if p.Desc != nil {
			v.Desc = p.Desc.ProjDesc()
		}

		if u, ok := userInfos[p.Owner.Account()]; ok {
			v.Owner.Name = u.Account.Account()
			v.Owner.AvatarId = u.AvatarId.AvatarId()
		}

		dtos[i] = v
	}
}

func (s likeService) modelToLikeDTO(
	userInfos map[string]*domain.UserInfo,
	data []domain.Model, dtos []LikeDTO,
) {
	for i := range data {
		d := &data[i]

		v := LikeDTO{
			Name: d.Name.ModelName(),
			/*
				UpdateAt string `json:"update_at"`

				LikeCount     int `json:"like_count"`
				DownloadCount int `json:"download_count"`
				ForkCount     int `json:"fork_count"`
			*/
		}

		if d.Desc != nil {
			v.Desc = d.Desc.ProjDesc()
		}

		if u, ok := userInfos[d.Owner.Account()]; ok {
			v.Owner.Name = u.Account.Account()
			v.Owner.AvatarId = u.AvatarId.AvatarId()
		}

		dtos[i] = v
	}
}

func (s likeService) datasetToLikeDTO(
	userInfos map[string]*domain.UserInfo,
	data []domain.Dataset, dtos []LikeDTO,
) {
	for i := range data {
		d := &data[i]

		v := LikeDTO{
			Name: d.Name.DatasetName(),
			/*
				UpdateAt string `json:"update_at"`

				LikeCount     int `json:"like_count"`
				DownloadCount int `json:"download_count"`
				ForkCount     int `json:"fork_count"`
			*/
		}

		if d.Desc != nil {
			v.Desc = d.Desc.ProjDesc()
		}

		if u, ok := userInfos[d.Owner.Account()]; ok {
			v.Owner.Name = u.Account.Account()
			v.Owner.AvatarId = u.AvatarId.AvatarId()
		}

		dtos[i] = v
	}
}
