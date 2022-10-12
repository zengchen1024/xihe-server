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

type TrainingCreateCmd struct {
	User      domain.Account
	ProjectId string
	*domain.Training
}

func (cmd *TrainingCreateCmd) Validate() error {
	err := errors.New("invalid cmd of creating training")

	b := cmd.User != nil &&
		cmd.ProjectId != "" &&
		cmd.ProjectName != nil &&
		cmd.ProjectRepoId != "" &&
		cmd.Name != nil &&
		cmd.CodeDir != nil &&
		cmd.BootFile != nil

	if !b {
		return err
	}

	c := &cmd.Compute
	if c.Flavor == nil || c.Type == nil || c.Version == nil {
		return err
	}

	f := func(kv []domain.KeyValue) error {
		for i := range kv {
			if kv[i].Key == nil {
				return err
			}
		}

		return nil
	}

	if f(cmd.Hypeparameters) != nil {
		return err
	}

	if f(cmd.Env) != nil {
		return err
	}

	for i := range cmd.Inputs {
		v := &cmd.Inputs[i]
		if v.Key == nil || cmd.checkInput(&v.Value) != nil {
			return err
		}
	}

	return nil
}

func (cmd *TrainingCreateCmd) checkInput(i *domain.ResourceInput) error {
	if i.User == nil || i.Type == nil || i.RepoId == "" {
		return errors.New("invalide input")
	}

	return nil
}

func (cmd *TrainingCreateCmd) toTraining(t *domain.UserTraining) {
	t.Owner = cmd.User
	t.ProjectId = cmd.ProjectId
	t.Training = *cmd.Training
	t.CreatedAt = utils.Now()
}

type TrainingService interface {
	CreateTrainingJob(info *domain.TrainingInfo, endpoint string) error
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
	v, version, err := s.repo.List(cmd.User, cmd.ProjectId)
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

	var t domain.UserTraining
	cmd.toTraining(&t)

	r, err := s.repo.Save(&t, version)
	if err != nil {
		return "", err
	}

	// send message
	err = s.sender.CreateTraining(&domain.TrainingInfo{
		User:       cmd.User,
		ProjectId:  cmd.ProjectId,
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

		if !s.isJobDone(item.JobDetail.Status) {
			detail, err := s.updateJobDetail(
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

		s.toTrainingSummaryDTO(&v[i], &r[i])
	}

	return r, nil
}

func (s trainingService) Get(info *domain.TrainingInfo) error {
	data, err := s.repo.Get(info)
	if err != nil {
		return err
	}

	if s.isJobDone(data.JobDetail.Status) {
		// convert data
		return nil
	}

	detail, err := s.updateJobDetail(
		info,
		data.Job.Endpoint, data.Job.JobId,
		data.JobDetail.Status,
	)
	if err == nil {
		data.JobDetail = detail
	}

	// convert

	return nil
}

func (s trainingService) updateJobDetail(
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

type TrainingSummaryDTO struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Status    string `json:"status"`
	Duration  int    `json:"duration"`
	CreatedAt string `json:"created_at"`
}

func (s trainingService) toTrainingSummaryDTO(t *domain.TrainingSummary, dto *TrainingSummaryDTO) {
	*dto = TrainingSummaryDTO{
		Id:        t.Id,
		Name:      t.Name.TrainingName(),
		Status:    t.JobDetail.Status,
		Duration:  t.JobDetail.Duration,
		CreatedAt: utils.ToDate(t.CreatedAt),
	}

	if t.Desc != nil {
		dto.Desc = t.Desc.TrainingDesc()
	}
}
