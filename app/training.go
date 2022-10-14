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

const trainingCreatedFailed = "create_failed"

type TrainingService interface {
	Create(cmd *TrainingCreateCmd) (string, error)
	Recreate(info *domain.TrainingInfo) (string, error)
	List(user domain.Account, projectId string) ([]TrainingSummaryDTO, error)
	Get(info *domain.TrainingInfo) (TrainingDTO, error)
	Delete(info *domain.TrainingInfo) error
	Terminate(info *domain.TrainingInfo) error
	GetLogDownloadURL(info *domain.TrainingInfo) (string, error)
	CreateTrainingJob(info *domain.TrainingInfo, endpoint string) error
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
	return status != "" && (s.train.IsJobDone(status) || status == trainingCreatedFailed)
}

func (s trainingService) Create(cmd *TrainingCreateCmd) (string, error) {
	return s.create(cmd.User, cmd.ProjectId, cmd.toTrainingConfig())
}

func (s trainingService) Recreate(info *domain.TrainingInfo) (string, error) {
	v, err := s.repo.GetTrainingConfig(info)
	if err != nil {
		return "", err
	}

	return s.create(info.User, info.ProjectId, &v)
}

func (s trainingService) create(
	user domain.Account, projectId string, config *domain.TrainingConfig,
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
	err = s.sender.CreateTraining(&domain.TrainingInfo{
		User:       user,
		ProjectId:  projectId,
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
		item := &v[i]

		done := s.isJobDone(item.JobDetail.Status)
		if !done {
			detail, err := s.getAndUpdateJobDetail(
				&domain.TrainingInfo{
					User:       user,
					ProjectId:  projectId,
					TrainingId: item.Id,
				},
				item.Endpoint, item.JobId,
				item.JobDetail.Status,
			)

			if err == nil {
				item.JobDetail = detail
			}
		}

		s.toTrainingSummaryDTO(&v[i], &r[i], done)
	}

	return r, nil
}

func (s trainingService) Get(info *domain.TrainingInfo) (TrainingDTO, error) {
	data, err := s.repo.Get(info)
	if err != nil {
		return TrainingDTO{}, err
	}

	if s.isJobDone(data.JobDetail.Status) {
		return s.toTrainingDTO(&data), nil
	}

	detail, err := s.getAndUpdateJobDetail(
		info,
		data.Job.Endpoint, data.Job.JobId,
		data.JobDetail.Status,
	)
	if err == nil {
		data.JobDetail = detail
	}

	return s.toTrainingDTO(&data), nil
}

func (s trainingService) getAndUpdateJobDetail(
	info *domain.TrainingInfo, endpoint, jobId, old string,
) (r domain.JobDetail, err error) {
	r, err = s.train.GetJob(endpoint, jobId)
	if err != nil {
		return
	}

	if r.Status != old {
		_ = s.repo.UpdateJobDetail(info, &r)
	}

	return
}

func (s trainingService) Delete(info *domain.TrainingInfo) error {
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

func (s trainingService) Terminate(info *domain.TrainingInfo) error {
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

func (s trainingService) GetLogDownloadURL(info *domain.TrainingInfo) (string, error) {
	job, err := s.repo.GetJob(info)
	if err != nil {
		return "", err
	}

	if job.JobId != "" {
		return s.train.GetLogDownloadURL(job.Endpoint, job.JobId)
	}

	return "", nil
}

func (s trainingService) CreateTrainingJob(info *domain.TrainingInfo, endpoint string) error {
	data, err := s.repo.Get(info)
	if err != nil {
		return err
	}

	if data.Job.JobId != "" {
		return nil
	}

	v, err := s.train.CreateJob(endpoint, info.User, &data.TrainingConfig)
	if err != nil {
		return s.repo.UpdateJobDetail(
			info,
			&domain.JobDetail{
				Status: trainingCreatedFailed,
			},
		)
	}

	return s.repo.SaveJob(info, &v)
}
