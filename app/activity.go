package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
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

	total := len(activities)
	objs := make([]*domain.ResourceObject, total)
	orders := make([]orderByTime, total)
	for i := range activities {
		item := &activities[i]

		objs[i] = &item.ResourceObject
		orders[i] = orderByTime{t: item.Time, p: i}
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

	dtos = make([]ActivityDTO, total)
	err = sortAndSet(orders, func(i, j int) error {
		item := &activities[i]

		r, ok := rm[item.String()]
		if !ok {
			return errors.New("no matched resource")
		}

		dtos[j] = ActivityDTO{
			Type:     item.Type.ActivityType(),
			Time:     utils.ToDate(item.Time),
			Resource: *r,
		}

		return nil
	})

	return
}

func genActivityForCreatingResource(
	owner domain.Account, t domain.ResourceType, rid string,
) domain.UserActivity {
	obj := domain.ResourceObject{Type: t}
	obj.Owner = owner
	obj.Id = rid

	return domain.UserActivity{
		Owner: owner,
		Activity: domain.Activity{
			Type:           domain.ActivityTypeCreate,
			Time:           utils.Now(),
			ResourceObject: obj,
		},
	}
}

func genActivityForDeletingResource(
	s *domain.ResourceSummary, t domain.ResourceType,
) domain.UserActivity {
	obj := domain.ResourceObject{Type: t}
	obj.ResourceIndex = s.ResourceIndex()

	return domain.UserActivity{
		Owner: s.Owner,
		Activity: domain.Activity{
			Type:           domain.ActivityTypeDelete,
			Time:           utils.Now(),
			ResourceObject: obj,
		},
	}
}
