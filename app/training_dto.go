package app

import (
	"errors"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

type TrainingCreateCmd struct {
	User      domain.Account
	ProjectId string

	domain.TrainingConfig
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

		if v.Key == nil || v.User == nil || v.Type == nil || v.RepoId == "" {
			return errors.New("invalide input")
		}
	}

	return nil
}

func (cmd *TrainingCreateCmd) toTrainingConfig() *domain.TrainingConfig {
	return &cmd.TrainingConfig
}

type TrainingSummaryDTO struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	IsDone    bool   `json:"is_done"`
	Duration  int    `json:"duration"`
}

func (s trainingService) toTrainingSummaryDTO(
	t *domain.TrainingSummary, dto *TrainingSummaryDTO, done bool,
) {
	status := t.JobDetail.Status
	if status == "" {
		status = trainingStatusScheduling
	}

	*dto = TrainingSummaryDTO{
		Id:        t.Id,
		Name:      t.Name.TrainingName(),
		Status:    status,
		IsDone:    done,
		Duration:  t.JobDetail.Duration,
		CreatedAt: utils.ToDate(t.CreatedAt),
	}

	if t.Desc != nil {
		dto.Desc = t.Desc.TrainingDesc()
	}
}

type TrainingDTO struct {
	Id        string `json:"id"`
	ProjectId string `json:"project_id"`

	Name string `json:"name"`
	Desc string `json:"desc"`

	IsDone    bool       `json:"is_done"`
	Status    string     `json:"status"`
	Duration  int        `json:"duration"`
	CreatedAt string     `json:"created_at"`
	Compute   ComputeDTO `json:"compute"`

	JobEndpoint string `json:"-"`
	JobId       string `json:"-"`
}

type ComputeDTO struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Flavor  string `json:"flavor"`
}

func (s trainingService) toTrainingDTO(ut *domain.UserTraining) TrainingDTO {
	t := &ut.TrainingConfig
	detail := &ut.JobDetail
	c := &t.Compute

	status := detail.Status
	if status == "" {
		status = trainingStatusScheduling
	}

	v := TrainingDTO{
		Id:        ut.Id,
		ProjectId: ut.ProjectId,

		Name:      t.Name.TrainingName(),
		IsDone:    s.isJobDone(detail.Status),
		Status:    status,
		Duration:  detail.Duration,
		CreatedAt: utils.ToDate(ut.CreatedAt),
		Compute: ComputeDTO{
			Type:    c.Type.ComputeType(),
			Flavor:  c.Flavor.ComputeFlavor(),
			Version: c.Version.ComputeVersion(),
		},

		JobEndpoint: ut.Job.Endpoint,
		JobId:       ut.Job.JobId,
	}

	if t.Desc != nil {
		v.Desc = t.Desc.TrainingDesc()
	}

	return v
}
