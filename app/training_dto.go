package app

import (
	"errors"
	"strconv"

	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/utils"
)

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

func (cmd *TrainingCreateCmd) toTraining() *domain.Training {
	return cmd.Training
}

type TrainingSummaryDTO struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Status    string `json:"status"`
	IsDone    bool   `json:"is_done"`
	Duration  int    `json:"duration"`
	CreatedAt string `json:"created_at"`
}

func (s trainingService) toTrainingSummaryDTO(
	t *domain.TrainingSummary, dto *TrainingSummaryDTO, done bool,
) {
	*dto = TrainingSummaryDTO{
		Id:        t.Id,
		Name:      t.Name.TrainingName(),
		Status:    t.JobDetail.Status,
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
	Duration  string     `json:"duration"`
	CreatedAt string     `json:"created_at"`
	Compute   ComputeDTO `json:"compute"`
}

type ComputeDTO struct {
	Type    string `json:"type"`
	Version string `json:"version"`
	Flavor  string `json:"flavor"`
}

func (s trainingService) toTrainingDTO(ut *domain.UserTraining) TrainingDTO {
	t := &ut.Training
	detail := &ut.JobDetail
	c := &t.Compute

	v := TrainingDTO{
		Id:        ut.Id,
		ProjectId: ut.ProjectId,

		Name:      t.Name.TrainingName(),
		IsDone:    s.isJobDone(detail.Status),
		Status:    detail.Status,
		Duration:  strconv.Itoa(detail.Duration),
		CreatedAt: utils.ToDate(ut.CreatedAt),
		Compute: ComputeDTO{
			Type:    c.Type.ComputeType(),
			Flavor:  c.Flavor.ComputeFlavor(),
			Version: c.Version.ComputeVersion(),
		},
	}

	if t.Desc != nil {
		v.Desc = t.Desc.TrainingDesc()
	}

	return v
}
