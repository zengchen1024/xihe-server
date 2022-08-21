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
	Time     string      `json:"time"`
	Resource ResourceDTO `json:"resource"`
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
	activity repository.Activity,
	sender message.Sender,
) LikeService {
	return likeService{
		repo:     repo,
		activity: activity,
		sender:   sender,

		rs: resourceService{
			user:    user,
			model:   model,
			project: project,
			dataset: dataset,
		},
	}
}

type likeService struct {
	repo     repository.Like
	activity repository.Activity
	sender   message.Sender

	rs resourceService
}

func (s likeService) Create(owner domain.Account, cmd LikeCreateCmd) error {
	v := domain.UserLike{
		Owner: owner,
		Like: domain.Like{
			ResourceObj: domain.ResourceObj{
				ResourceOwner: cmd.ResourceOwner,
				ResourceType:  cmd.ResourceType,
				ResourceId:    cmd.ResourceId,
			},
		},
	}

	if err := s.repo.Save(&v); err != nil {
		return err
	}

	ua := domain.UserActivity{
		Owner: owner,
		Activity: domain.Activity{
			Type: domain.NewActivityTypeLike(),
			ResourceObj: domain.ResourceObj{
				ResourceOwner: cmd.ResourceOwner,
				ResourceType:  cmd.ResourceType,
				ResourceId:    cmd.ResourceId,
			},
		},
	}
	if err := s.activity.Save(&ua); err != nil {
		return err
	}

	// send event
	return s.sender.AddLike(v.Like)
}

func (s likeService) Delete(owner domain.Account, cmd LikeRemoveCmd) error {
	v := domain.UserLike{
		Owner: owner,
		Like: domain.Like{
			ResourceObj: domain.ResourceObj{
				ResourceOwner: cmd.ResourceOwner,
				ResourceType:  cmd.ResourceType,
				ResourceId:    cmd.ResourceId,
			},
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
	likes, err := s.repo.Find(owner, repository.LikeFindOption{})
	if err != nil || len(likes) == 0 {
		return
	}

	objs := make([]*domain.ResourceObj, len(likes))
	for i := range likes {
		objs[i] = &likes[i].ResourceObj
	}

	resources, err := s.rs.list(objs)
	if err != nil {
		return
	}

	rm := make(map[string]*ResourceDTO)
	for i := range resources {
		item := &resources[i]

		rm[item.identity()] = item
	}

	dtos = make([]LikeDTO, len(likes))
	for i := range likes {
		item := &likes[i]

		r, ok := rm[item.String()]
		if !ok {
			return nil, errors.New("no matched resource")
		}

		dtos[i] = LikeDTO{
			Time:     item.CreatedAt,
			Resource: *r,
		}
	}

	return
}
