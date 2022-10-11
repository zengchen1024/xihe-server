package app

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/domain/training"
)

const trainingCreatedFailed = "create_failed"

type TrainingInfo = repository.TrainingInfo

type trainingService struct {
	train training.Training
	repo  repository.Training
}

func (s trainingService) Create() (string, error) {
	return "", nil
}

func (s trainingService) CreateTrainingJob(info *TrainingInfo, endpoint string) error {
	data, err := s.repo.Get(info)
	if err != nil {
		return err
	}

	if data.Job.JobId != "" {
		return nil
	}

	v, err := s.train.CreateJob(endpoint, info.User, &data.Training)
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

func (s trainingService) Get(info *TrainingInfo) error {
	data, err := s.repo.Get(info)
	if err != nil {
		return err
	}

	status := data.JobDetail.Status
	if status != "" && (s.train.IsJobDone(status) || status == trainingCreatedFailed) {
		// convert data
		return nil
	}

	detail, err := s.train.GetJob(data.Job.Endpoint, data.Job.JobId)
	if err != nil {
		return err
	}

	if s.train.IsJobDone(detail.Status) {
		_ = s.repo.UpdateJobDetail(info, &detail)
	}

	data.JobDetail = detail

	// convert

	return nil
}

func (s trainingService) Delete(info *TrainingInfo) error {
	data, err := s.repo.Get(info)
	if err != nil {
		return err
	}

	job := &data.Job

	if job.JobId != "" {
		err = s.train.DeleteJob(job.Endpoint, job.JobId)
		if err != nil {
			// ignore 404
			return err
		}
	}

	return s.repo.Delete(info)
}

func (s trainingService) Terminate(info *TrainingInfo) error {
	data, err := s.repo.Get(info)
	if err != nil {
		return err
	}

	job := &data.Job

	if job.JobId != "" {
		err = s.train.TerminateJob(job.Endpoint, job.JobId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s trainingService) GetLogDownloadURL(info *TrainingInfo) (string, error) {
	data, err := s.repo.Get(info)
	if err != nil {
		return "", err
	}

	if job := &data.Job; job.JobId != "" {
		return s.train.GetLogDownloadURL(job.Endpoint, job.JobId)
	}

	return "", nil
}
