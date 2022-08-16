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
	}

	Name     string `json:"name"`
	Desc     string `json:"description"`
	CoverId  string `json:"cover_id"`
	UpdateAt string `json:"update_at"`

	LikeCount     int `json:"like_count"`
	DownloadCount int `json:"download_count"`
	ForkCount     int `json:"fork_count"`
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

func (s likeService) RemoveLike(owner domain.Account, cmd LikeRemoveCmd) error {
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

func (s likeService) ListLike(owner domain.Account) (
	dtos []LikeDTO, err error,
) {
	v, err := s.repo.Find(owner, repository.LikeFindOption{})
	if err != nil || len(v) == 0 {
		return
	}

	users := map[string]domain.Account{}
	projects := map[string]*repository.UserResourceListOption{}
	datasets := map[string]*repository.UserResourceListOption{}
	models := map[string]*repository.UserResourceListOption{}

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

	for i := range v {
		item := &v[i]

		account := item.ResourceOwner.Account()
		if _, ok := users[account]; !ok {
			users[account] = item.ResourceOwner
		}

		switch item.ResourceType.ResourceType() {
		case domain.ResourceProject:
			set(item, projects)

		case domain.ResourceModel:
			set(item, models)

		case domain.ResourceDataset:
			set(item, datasets)
		}

	}

	// TODO get user
	opts := make([]domain.Account, 0, len(users))
	for _, u := range users {
		opts = append(opts, u)
	}

	allUsers, err := s.user.Find(repository.UserFindOption{Names: opts})
	if err != nil {
		return
	}

	userInfos := make(map[string]*domain.UserInfo)
	for i := range allUsers {
		item := &allUsers[i]
		userInfos[item.Account.Account()] = item
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

	dtos = make([]LikeDTO, len(v))
	r := dtos

	if n := len(projects); n > 0 {
		all, err := s.project.FindUserProjects(toList(projects))
		if err != nil {
			return nil, err
		}

		s.projectToLikeDTO(userInfos, all, r)
		r = r[n:]
	}

	if n := len(models); n > 0 {
		all, err := s.model.FindUserModels(toList(models))
		if err != nil {
			return nil, err
		}

		s.modelToLikeDTO(userInfos, all, r)
		r = r[n:]
	}

	if n := len(datasets); n > 0 {
		all, err := s.dataset.FindUserDatasets(toList(datasets))
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
	}
}
