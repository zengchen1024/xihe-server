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

func (do *TrainingDO) toTraining() (t domain.Training, err error) {
	if t.ProjectName, err = domain.NewProjName(do.ProjectName); err != nil {
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

	t.Inputs, err = do.toInputs()

	return
}

func (do *TrainingDO) toKeyValues(kv []KeyValueDO) (r []domain.KeyValue, err error) {
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

func (do *TrainingDO) toInputs() (r []domain.Input, err error) {
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
	Version string
	Flavor  string
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

	v := &r.Value

	if v.User, err = domain.NewAccount(do.User); err != nil {
		return
	}

	if v.Type, err = domain.NewResourceType(do.Type); err != nil {
		return
	}

	v.RepoId = do.RepoId
	v.File = do.File

	return
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

func (impl training) toTrainingInfoDo(info *domain.TrainingInfo) TrainingInfoDO {
	return TrainingInfoDO{
		User:       info.User.Account(),
		ProjectId:  info.ProjectId,
		TrainingId: info.TrainingId,
	}
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

type TrainingInfoDO struct {
	User       string
	ProjectId  string
	TrainingId string
}

type TrainingJobInfoDO = domain.JobInfo
type TrainingJobDetailDO = domain.JobDetail

type TrainingDetailDO struct {
	TrainingDO

	Job       TrainingJobInfoDO
	JobDetail TrainingJobDetailDO
	CreatedAt int64
}

func (do *TrainingDetailDO) toUserTraining() (ut domain.UserTraining, err error) {
	if ut.Training, err = do.TrainingDO.toTraining(); err != nil {
		return
	}

	ut.Job = do.Job
	ut.JobDetail = do.JobDetail
	ut.CreatedAt = do.CreatedAt

	return
}
