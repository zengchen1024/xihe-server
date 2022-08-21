package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ActivityDTO struct {
	Type     string      `json:"type"`
	Time     string      `json:"time"`
	Resource ResourceDTO `json:"resource"`
}

type ActivityService interface {
	List(domain.Account) ([]ActivityDTO, error)
}

func NewActivityService(
	repo repository.Activity,
	user repository.User,
	model repository.Model,
	project repository.Project,
	dataset repository.Dataset,
) ActivityService {
	return activityService{
		repo: repo,
		rs: resourceService{
			user:    user,
			model:   model,
			project: project,
			dataset: dataset,
		},
	}
}

type activityService struct {
	repo repository.Activity
	rs   resourceService
}

func (s activityService) List(owner domain.Account) (
	dtos []ActivityDTO, err error,
) {
	activities, err := s.repo.Find(owner, repository.ActivityFindOption{})
	if err != nil || len(activities) == 0 {
		return
	}

	objs := make([]*domain.ResourceObj, len(activities))
	for i := range activities {
		objs[i] = &activities[i].ResourceObj
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

	dtos = make([]ActivityDTO, len(activities))
	for i := range activities {
		item := &activities[i]

		r, ok := rm[item.String()]
		if !ok {
			return nil, errors.New("no matched resource")
		}

		dtos[i] = ActivityDTO{
			Type:     item.Type.ActivityType(),
			Time:     item.Time,
			Resource: *r,
		}
	}

	return
}
