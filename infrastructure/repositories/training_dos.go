package repositories

import "github.com/opensourceways/xihe-server/domain"

type UserTrainingDO struct {
	Id        string
	Owner     string
	ProjectId string

	TrainingDO

	CreatedAt int64
}

type TrainingDO struct {
	ProjectName   string
	ProjectRepoId string

	Name string
	Desc string

	CodeDir  string
	BootFile string

	Hypeparameters []KeyValueDO
	Env            []KeyValueDO
	Inputs         []InputDO

	Compute ComputeDO
}

type KeyValueDO struct {
	Key   string
	Value string
}

type ComputeDO struct {
	Type    string
	Version string
	Flavor  string
}

type InputDO struct {
	Key    string
	User   string
	Type   string
	RepoId string
	File   string
}

func (impl training) toUserTrainingDO(ut *domain.UserTraining) UserTrainingDO {
	t := &ut.Training
	c := &t.Compute

	do := UserTrainingDO{
		Id:        ut.Id,
		Owner:     ut.Owner.Account(),
		ProjectId: ut.ProjectId,
		CreatedAt: ut.CreatedAt,

		TrainingDO: TrainingDO{
			ProjectName:   t.ProjectName.ProjName(),
			ProjectRepoId: t.ProjectRepoId,
			Name:          t.Name.TrainingName(),

			CodeDir:  t.CodeDir.Directory(),
			BootFile: t.BootFile.FilePath(),

			Hypeparameters: impl.toKeyValueDOs(t.Hypeparameters),
			Env:            impl.toKeyValueDOs(t.Env),
			Inputs:         impl.toInputDOs(t.Inputs),

			Compute: ComputeDO{
				Type:    c.Type.ComputeType(),
				Version: c.Version.ComputeVersion(),
				Flavor:  c.Flavor.ComputeFlavor(),
			},
		},
	}

	if t.Desc != nil {
		do.TrainingDO.Desc = t.Desc.TrainingDesc()
	}

	return do
}

func (impl training) toKeyValueDOs(kv []domain.KeyValue) []KeyValueDO {
	n := len(kv)
	if n == 0 {
		return nil
	}

	r := make([]KeyValueDO, n)

	for i := range kv {
		r[i].Key = kv[i].Key.CustomizedKey()

		if kv[i].Value != nil {
			r[i].Value = kv[i].Value.CustomizedValue()
		}
	}

	return r
}

func (impl training) toInputDOs(v []domain.Input) []InputDO {
	n := len(v)
	if n == 0 {
		return nil
	}

	r := make([]InputDO, n)

	for i := range v {
		item := &v[i].Value

		r[i] = InputDO{
			Key:    v[i].Key.CustomizedKey(),
			User:   item.User.Account(),
			Type:   item.Type.ResourceType(),
			File:   item.File,
			RepoId: item.RepoId,
		}
	}

	return r
}

func (impl training) toTrainingSummary(do *TrainingSummaryDO, t *domain.TrainingSummary) (
	err error,
) {
	if t.Name, err = domain.NewTrainingName(do.Name); err != nil {
		return
	}

	if t.Desc, err = domain.NewTrainingDesc(do.Desc); err != nil {
		return
	}

	t.JobId = do.JobId
	t.Endpoint = do.Endpoint
	t.CreatedAt = do.CreatedAt
	t.JobDetail.Status = do.Status
	t.JobDetail.Duration = do.Duration

	return
}

type TrainingSummaryDO struct {
	Name      string
	Desc      string
	JobId     string
	Status    string
	Endpoint  string
	Duration  int
	CreatedAt int64
}
