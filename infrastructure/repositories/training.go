package repositories

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type TrainingMapper interface {
	Insert(*UserTrainingDO, int) (string, error)
	Delete(*TrainingInfoDO) error
	Get(*TrainingInfoDO) (TrainingDetailDO, error)
	GetTrainingConfig(*TrainingInfoDO) (TrainingConfigDO, error)
	List(user, projectId string) ([]TrainingSummaryDO, int, error)
	UpdateJobInfo(*TrainingInfoDO, *TrainingJobInfoDO) error
	GetJobInfo(*TrainingInfoDO) (TrainingJobInfoDO, error)
	UpdateJobDetail(*TrainingInfoDO, *TrainingJobDetailDO) error
}

func NewTrainingRepository(mapper TrainingMapper) repository.Training {
	return training{mapper}
}

type training struct {
	mapper TrainingMapper
}

func (impl training) Save(ut *domain.UserTraining, version int) (string, error) {
	if ut.Id != "" {
		return "", errors.New("must be a new project")
	}

	do := impl.toUserTrainingDO(ut)

	v, err := impl.mapper.Insert(&do, version)
	if err != nil {
		return "", convertError(err)
	}

	return v, nil
}

func (impl training) Delete(info *domain.TrainingInfo) error {
	do := impl.toTrainingInfoDo(info)

	if err := impl.mapper.Delete(&do); err != nil {
		return convertError(err)
	}

	return nil
}

func (impl training) Get(info *domain.TrainingInfo) (domain.UserTraining, error) {
	do := impl.toTrainingInfoDo(info)

	v, err := impl.mapper.Get(&do)
	if err != nil {
		return domain.UserTraining{}, convertError(err)
	}

	return v.toUserTraining()
}

func (impl training) GetTrainingConfig(info *domain.TrainingInfo) (domain.TrainingConfig, error) {
	do := impl.toTrainingInfoDo(info)

	v, err := impl.mapper.GetTrainingConfig(&do)
	if err != nil {
		return domain.TrainingConfig{}, convertError(err)
	}

	return v.toTrainingConfig()
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

func (impl training) GetJob(info *domain.TrainingInfo) (job domain.JobInfo, err error) {
	t := impl.toTrainingInfoDo(info)

	do, err := impl.mapper.GetJobInfo(&t)
	if err != nil {
		err = convertError(err)

		return
	}

	return do, nil
}

func (impl training) UpdateJobDetail(info *domain.TrainingInfo, detail *domain.JobDetail) error {
	do := impl.toTrainingInfoDo(info)

	if err := impl.mapper.UpdateJobDetail(&do, detail); err != nil {
		return convertError(err)
	}

	return nil
}
