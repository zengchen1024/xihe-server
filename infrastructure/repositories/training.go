package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
)

type TrainingMapper interface {
	Insert(*UserTrainingDO) (string, error)
	List(user, projectId string) ([]TrainingSummaryDO, int, error)
	UpdateJobInfo(*TrainingInfoDO, *TrainingJobInfoDO) error
	UpdateJobDetail(*TrainingInfoDO, *TrainingJobDetailDO) error
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
		err = convertError(err)

		return
	}

	if len(v) == 0 {
		return
	}

	r = make([]domain.TrainingSummary, len(v))
	for i := range v {
		if err = impl.toTrainingSummary(&v[i], &r[i]); err != nil {
			return
		}
	}

	return
}

func (impl training) SaveJob(info *domain.TrainingInfo, job *domain.JobInfo) error {
	do := impl.toTrainingInfoDo(info)

	if err := impl.mapper.UpdateJobInfo(&do, job); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl training) UpdateJobDetail(info *domain.TrainingInfo, detail *domain.JobDetail) error {
	do := impl.toTrainingInfoDo(info)

	if err := impl.mapper.UpdateJobDetail(&do, detail); err != nil {
		return convertError(err)
	}

	return nil
}
