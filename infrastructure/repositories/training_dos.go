package repositories

import "github.com/opensourceways/xihe-server/domain"

type UserTrainingDO struct {
	Id        string
	Owner     string
	ProjectId string

	TrainingConfigDO

	CreatedAt int64
}

type TrainingConfigDO struct {
	ProjectName   string
	ProjectRepoId string

	Name string
	Desc string

	CodeDir  string
	BootFile string

	Hypeparameters []KeyValueDO
	Env            []KeyValueDO
	Inputs         []InputDO
	EnableAim      bool
	EnableOutput   bool

	Compute ComputeDO
}

func (do *TrainingConfigDO) toTrainingConfig() (t domain.TrainingConfig, err error) {
	if t.ProjectName, err = domain.NewResourceName(do.ProjectName); err != nil {
		return
	}

	t.ProjectRepoId = do.ProjectRepoId

	if t.Name, err = domain.NewTrainingName(do.Name); err != nil {
		return
	}

	if t.Desc, err = domain.NewTrainingDesc(do.Desc); err != nil {
		return
	}

	if t.CodeDir, err = domain.NewDirectory(do.CodeDir); err != nil {
		return
	}

	if t.BootFile, err = domain.NewFilePath(do.BootFile); err != nil {
		return
	}

	if t.Compute, err = do.Compute.toCompute(); err != nil {
		return
	}

	if t.Hypeparameters, err = do.toKeyValues(do.Hypeparameters); err != nil {
		return
	}

	if t.Env, err = do.toKeyValues(do.Env); err != nil {
		return
	}

	if t.Inputs, err = do.toInputs(); err != nil {
		return
	}

	t.EnableOutput = do.EnableOutput
	t.EnableAim = do.EnableAim

	return
}

func (do *TrainingConfigDO) toKeyValues(kv []KeyValueDO) (r []domain.KeyValue, err error) {
	if len(kv) == 0 {
		return
	}

	r = make([]domain.KeyValue, len(kv))

	for i := range kv {
		if r[i], err = kv[i].toKeyValue(); err != nil {
			return
		}
	}

	return
}

func (do *TrainingConfigDO) toInputs() (r []domain.Input, err error) {
	v := do.Inputs
	if len(v) == 0 {
		return
	}

	r = make([]domain.Input, len(v))

	for i := range v {
		if r[i], err = v[i].toInput(); err != nil {
			return
		}
	}

	return
}

type KeyValueDO struct {
	Key   string
	Value string
}

func (kv *KeyValueDO) toKeyValue() (r domain.KeyValue, err error) {
	if r.Key, err = domain.NewCustomizedKey(kv.Key); err != nil {
		return
	}

	r.Value, err = domain.NewCustomizedValue(kv.Value)

	return
}

type ComputeDO struct {
	Type    string
	Flavor  string
	Version string
}

func (do *ComputeDO) toCompute() (r domain.Compute, err error) {
	if r.Type, err = domain.NewComputeType(do.Type); err != nil {
		return
	}

	if r.Version, err = domain.NewComputeVersion(do.Version); err != nil {
		return
	}

	r.Flavor, err = domain.NewComputeFlavor(do.Flavor)

	return
}

type InputDO struct {
	Key    string
	User   string
	Type   string
	RepoId string
	File   string
}

func (do *InputDO) toInput() (r domain.Input, err error) {
	if r.Key, err = domain.NewCustomizedKey(do.Key); err != nil {
		return
	}

	if r.User, err = domain.NewAccount(do.User); err != nil {
		return
	}

	if r.Type, err = domain.NewResourceType(do.Type); err != nil {
		return
	}

	r.RepoId = do.RepoId
	r.File = do.File

	return
}

func (impl training) toUserTrainingDO(ut *domain.UserTraining) UserTrainingDO {
	t := &ut.TrainingConfig
	c := &t.Compute

	do := UserTrainingDO{
		Id:        ut.Id,
		Owner:     ut.Owner.Account(),
		ProjectId: ut.ProjectId,
		CreatedAt: ut.CreatedAt,

		TrainingConfigDO: TrainingConfigDO{
			Name:          t.Name.TrainingName(),
			ProjectName:   t.ProjectName.ResourceName(),
			ProjectRepoId: t.ProjectRepoId,

			CodeDir:  t.CodeDir.Directory(),
			BootFile: t.BootFile.FilePath(),

			Hypeparameters: impl.toKeyValueDOs(t.Hypeparameters),
			Env:            impl.toKeyValueDOs(t.Env),
			Inputs:         impl.toInputDOs(t.Inputs),
			EnableAim:      t.EnableAim,
			EnableOutput:   t.EnableOutput,

			Compute: ComputeDO{
				Type:    c.Type.ComputeType(),
				Flavor:  c.Flavor.ComputeFlavor(),
				Version: c.Version.ComputeVersion(),
			},
		},
	}

	if t.Desc != nil {
		do.TrainingConfigDO.Desc = t.Desc.TrainingDesc()
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
		item := &v[i]

		r[i] = InputDO{
			Key:    item.Key.CustomizedKey(),
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

	t.Id = do.Id
	t.JobId = do.JobId
	t.Endpoint = do.Endpoint
	t.CreatedAt = do.CreatedAt
	t.Error = do.Error
	t.Status = do.Status
	t.Duration = do.Duration

	return
}

func (impl training) toTrainingInfoDo(info *domain.TrainingIndex) TrainingIndexDO {
	return TrainingIndexDO{
		User:       info.Project.Owner.Account(),
		ProjectId:  info.Project.Id,
		TrainingId: info.TrainingId,
	}
}

type TrainingSummaryDO struct {
	Id        string
	Name      string
	Desc      string
	JobId     string
	Error     string
	Status    string
	Endpoint  string
	Duration  int
	CreatedAt int64
}

type TrainingIndexDO struct {
	User       string
	ProjectId  string
	TrainingId string
}

type TrainingJobInfoDO = domain.JobInfo
type TrainingJobDetailDO = domain.JobDetail

type TrainingDetailDO struct {
	TrainingConfigDO

	Job       TrainingJobInfoDO
	JobDetail TrainingJobDetailDO
	CreatedAt int64
}

func (do *TrainingDetailDO) toUserTraining() (ut domain.UserTraining, err error) {
	if ut.TrainingConfig, err = do.TrainingConfigDO.toTrainingConfig(); err != nil {
		return
	}

	ut.Job = do.Job
	ut.JobDetail = do.JobDetail
	ut.CreatedAt = do.CreatedAt

	return
}
