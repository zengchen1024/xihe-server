package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/finetune"
	"github.com/opensourceways/xihe-server/domain/message"
	"github.com/opensourceways/xihe-server/domain/repository"
	"github.com/opensourceways/xihe-server/utils"
	"github.com/sirupsen/logrus"
)

type FinetuneIndex = domain.FinetuneIndex
type FinetuneConfig = domain.FinetuneConfig
type FinetuneJobDetail = domain.FinetuneJobDetail

type FinetuneService interface {
	Create(*FinetuneCreateCmd) (string, string, error)
	List(user domain.Account) ([]FinetuneSummaryDTO, error)
	Delete(*FinetuneIndex) error
	Terminate(*FinetuneIndex) error
	GetJobInfo(*FinetuneIndex) (FinetuneJobDTO, string, error)
}

func NewFinetuneService(
	fs finetune.Finetune,
	repo repository.Finetune,
	sender message.Sender,
) FinetuneService {
	return finetuneService{
		fs:     fs,
		repo:   repo,
		sender: sender,
	}
}

type finetuneService struct {
	fs     finetune.Finetune
	repo   repository.Finetune
	sender message.Sender
}

func (s finetuneService) isJobDone(status string) bool {
	return status != "" && (s.fs.IsJobDone(status) || status == trainingStatusScheduleFailed)
}

func (s finetuneService) Create(cmd *FinetuneCreateCmd) (
	fid string, code string, err error,
) {
	user := cmd.User

	v, version, err := s.repo.List(user)
	if err != nil {
		return
	}

	if len(v) >= appConfig.FinetuneMaxNum {
		code = ErrorFinetuneExccedMaxNum
		err = errors.New("exceed max record num")

		return
	}

	for i := range v {
		if !s.isJobDone(v[i].Status) {
			code = ErrorFinetuneRunningJobExists
			err = errors.New("a job is running")

			return
		}
	}

	t := domain.UserFinetune{
		CreatedAt:      utils.Now(),
		FinetuneConfig: *cmd.toFinetuneConfig(),
	}
	t.Owner = user

	if fid, err = s.repo.Save(&t, version); err != nil {
		return
	}

	// send message
	err1 := s.sender.CreateFinetune(&FinetuneIndex{
		Owner: user,
		Id:    fid,
	})
	if err1 != nil {
		logrus.Errorf("send message of creating finetune failed, err:%s", err.Error())
	}

	return
}

func (s finetuneService) List(user domain.Account) ([]FinetuneSummaryDTO, error) {
	v, _, err := s.repo.List(user)
	if err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]FinetuneSummaryDTO, len(v))
	for i := range v {
		s.toFinetuneSummaryDTO(&v[i], &r[i])
	}

	return r, nil
}

func (s finetuneService) Delete(info *FinetuneIndex) error {
	job, err := s.repo.GetJob(info)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			return nil
		}
		return err
	}

	if job.JobId != "" {
		if err = s.fs.DeleteJob(job.Endpoint, job.JobId); err != nil {
			// ignore 404
			return err
		}
	}

	return s.repo.Delete(info)
}

func (s finetuneService) Terminate(info *FinetuneIndex) error {
	job, err := s.repo.GetJob(info)
	if err != nil || job.JobId == "" {
		return err
	}

	if !s.fs.CanTerminate(job.Status) {
		return errors.New("can't terminate now")
	}

	return s.fs.TerminateJob(job.Endpoint, job.JobId)
}

func (s finetuneService) GetJobInfo(index *FinetuneIndex) (
	dto FinetuneJobDTO, code string, err error,
) {
	job, err := s.repo.GetJob(index)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			code = ErrorFinetuneNotFound
		}

		return
	}

	dto.IsDone = s.isJobDone(job.Status)

	if job.Endpoint != "" && job.JobId != "" {
		dto.LogPreviewURL, err = s.fs.GetLogPreviewURL(job.Endpoint, job.JobId)
	}

	return
}

// FinetuneInternalService
type FinetuneInternalService interface {
	UpdateJobDetail(*FinetuneIndex, *FinetuneJobDetail) error
}

type finetuneInternalService struct {
	repo repository.Finetune
}

func (s finetuneInternalService) UpdateJobDetail(info *FinetuneIndex, v *FinetuneJobDetail) error {
	return s.repo.UpdateJobDetail(info, v)
}

// FinetuneMessageService
type FinetuneMessageService interface {
	CreateFinetuneJob(*FinetuneIndex, string, bool) (bool, error)
}

type finetuneMessageService struct {
	fs   finetune.Finetune
	repo repository.Finetune
}

func (s finetuneMessageService) CreateFinetuneJob(
	info *FinetuneIndex, endpoint string, lastChance bool,
) (retry bool, err error) {
	retry, err = s.createFinetuneJob(info, endpoint)
	if err == nil {
		return
	}

	if lastChance {
		s.repo.UpdateJobDetail(info, &FinetuneJobDetail{
			Status: trainingStatusScheduleFailed,
			Error:  err.Error(),
		})
	}

	return
}

func (s finetuneMessageService) createFinetuneJob(info *FinetuneIndex, endpoint string) (
	retry bool, err error,
) {
	data, err := s.repo.Get(info)
	if err != nil {
		if repository.IsErrorResourceNotExists(err) {
			err = errors.New("finetune is not exist.")
		} else {
			retry = true
		}

		return
	}

	if data.Job.JobId != "" {
		return
	}

	v, err := s.fs.CreateJob(endpoint, info, &data.FinetuneConfig)
	if err != nil {
		// TODO maybe can't retry based on error code
		retry = true

		return
	}

	if err1 := s.repo.SaveJob(info, &v); err1 != nil {
		logrus.Errorf(
			"create finetune(%s) job(%s) successfully, but save db err:%s",
			info.Id, v.JobId, err1.Error(),
		)
	}

	return
}
