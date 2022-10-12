package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
)

type TrainingMapper interface {
	Insert(*UserTrainingDO) (string, error)
	List(user, projectId string) ([]TrainingSummaryDO, int, error)
}

type training struct {
	mapper TrainingMapper
}

func (impl training) Save(ut *domain.UserTraining) (string, error) {
	if ut.Id != "" {
		return "", errors.New("must be a new project")
	}

	do := impl.toUserTrainingDO(ut)

	v, err := impl.mapper.Insert(&do)
	if err != nil {
		return "", convertError(err)
	}

	return v, nil
}

func (impl training) List(user domain.Account, projectId string) (
	r []domain.TrainingSummary, version int, err error,
) {
	v, version, err := impl.mapper.List(user.Account(), projectId)
	if err != nil {
		return
	}

	return
}
