package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/domain/training"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

const (
	trainingStatusScheduling     = "scheduling"
	trainingStatusScheduleFailed = "schedule_failed"
)

type JobDetail = domain.JobDetail
type TrainingIndex = domain.TrainingIndex
type TrainingConfig = domain.TrainingConfig

type TrainingService interface {
	Create(*TrainingCreateCmd) (string, error)
	Recreate(*TrainingIndex) (string, error)
	UpdateJobDetail(*TrainingIndex, *JobDetail) error
	List(user domain.Account, projectId string) ([]TrainingSummaryDTO, error)
	Get(*TrainingIndex) (TrainingDTO, error)
	Delete(*TrainingIndex) error
	Terminate(*TrainingIndex) error
	GetLogDownloadURL(*TrainingIndex) (string, error)
	CreateTrainingJob(*TrainingIndex, string, bool) (bool, error)
}

func NewTrainingService(
	log *logrus.Entry,
	train training.Training,
	repo repository.Training,
	sender message.Sender,
	maxTrainingRecordNum int,
) TrainingService {
	return trainingService{
		log:    log,
		train:  train,
		repo:   repo,
		sender: sender,

		maxTrainingRecordNum: maxTrainingRecordNum,
	}
}

type trainingService struct {
	log    *logrus.Entry
	train  training.Training
	repo   repository.Training
	sender message.Sender

	maxTrainingRecordNum int
}

func (s trainingService) isJobDone(status string) bool {
	return status != "" && (s.train.IsJobDone(status) || status == trainingStatusScheduleFailed)
}

func (s trainingService) Create(cmd *TrainingCreateCmd) (string, error) {
	return s.create(cmd.User, cmd.ProjectId, cmd.toTrainingConfig())
}

func (s trainingService) Recreate(info *TrainingIndex) (string, error) {
	v, err := s.repo.GetTrainingConfig(info)
	if err != nil {
		return "", err
	}

	return s.create(info.Project.Owner, info.Project.Id, &v)
}

func (s trainingService) create(
	user domain.Account, projectId string, config *TrainingConfig,
) (string, error) {
	v, version, err := s.repo.List(user, projectId)
	if err != nil {
		return "", err
	}

	if len(v) >= s.maxTrainingRecordNum {
		return "", ErrorExccedMaxTrainingRecord{
			errors.New("exceed max training num"),
		}
	}

	for i := range v {
		if !s.isJobDone(v[i].JobDetail.Status) {
			return "", ErrorOnlyOneRunningTraining{
				errors.New("a training is running"),
			}
		}
	}

	t := domain.UserTraining{
		Owner:          user,
		ProjectId:      projectId,
		CreatedAt:      utils.Now(),
		TrainingConfig: *config,
	}

	r, err := s.repo.Save(&t, version)
	if err != nil {
		return "", err
	}

	// send message
	err = s.sender.CreateTraining(&TrainingIndex{
		Project: domain.ResourceIndex{
			Owner: user,
			Id:    projectId,
		},
		TrainingId: r,
	})
	if err != nil {
		s.log.Errorf("send message of creating training failed, err:%s", err.Error())
	}

	return r, nil
}

func (s trainingService) List(user domain.Account, projectId string) ([]TrainingSummaryDTO, error) {
	v, _, err := s.repo.List(user, projectId)
	if err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]TrainingSummaryDTO, len(v))

	for i := range v {
		s.toTrainingSummaryDTO(&v[i], &r[i])
	}

	return r, nil
}

func (s trainingService) Get(info *TrainingIndex) (TrainingDTO, error) {
	data, err := s.repo.Get(info)
	if err != nil {
		return TrainingDTO{}, err
	}

	return s.toTrainingDTO(&data), nil
}

func (s trainingService) UpdateJobDetail(info *TrainingIndex, v *JobDetail) error {
	return s.repo.UpdateJobDetail(info, v)
}

func (s trainingService) Delete(info *TrainingIndex) error {
	job, err := s.repo.GetJob(info)
	if err != nil {
		return err
	}

	if job.JobId != "" {
		err = s.train.DeleteJob(job.Endpoint, job.JobId)
		if err != nil {
			// ignore 404
			return err
		}
	}

	return s.repo.Delete(info)
}

func (s trainingService) Terminate(info *TrainingIndex) error {
	job, err := s.repo.GetJob(info)
	if err != nil {
		return err
	}

	if job.JobId != "" {
		err = s.train.TerminateJob(job.Endpoint, job.JobId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s trainingService) GetLogDownloadURL(info *TrainingIndex) (string, error) {
	job, err := s.repo.GetJob(info)
	if err != nil {
		return "", err
	}

	if job.JobId != "" {
		return s.train.GetLogDownloadURL(job.Endpoint, job.JobId)
	}

	return "", nil
}

func (s trainingService) CreateTrainingJob(
	info *TrainingIndex, endpoint string, lastChance bool,
) (retry bool, err error) {
	retry, err = s.createTrainingJob(info, endpoint)
	if err == nil {
		return
	}

	if lastChance {
		s.repo.UpdateJobDetail(info, &JobDetail{
			Status: trainingStatusScheduleFailed,
		})
	}

	return
}

func (s trainingService) createTrainingJob(info *TrainingIndex, endpoint string) (
	retry bool, err error,
) {
	data, err := s.repo.Get(info)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			err = errorUnavailableTraining{
				errors.New("training is not exist."),
			}
		} else {
			retry = true
		}

		return
	}

	if data.Job.JobId != "" {
		return false, nil
	}

	v, err := s.train.CreateJob(endpoint, info, &data.TrainingConfig)
	if err != nil {
		retry = true

		return
	}

	if err1 := s.repo.SaveJob(info, &v); err1 != nil {
		s.log.Errorf(
			"create training(%s) job(%s) successfully, but save db err:%s",
			info.TrainingId, v.JobId, err1.Error(),
		)
	}

	return
}
