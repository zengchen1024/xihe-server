package domain

type UserTraining struct {
	Id        string
	Owner     Account
	ProjectId string

	Training

	CreatedAt int64

	// following fileds is not under the controlling of version
	Job       JobInfo
	JobDetail JobDetail
}

type Training struct {
	ProjectName   ProjName
	ProjectRepoId string

	Name TrainingName
	Desc TrainingDesc

	CodeDir  Directory
	BootFile FilePath

	Hypeparameters []KeyValue
	Env            []KeyValue
	Inputs         []Input

	Compute Compute
}

type Compute struct {
	Type    ComputeType
	Version ComputeVersion
	Flavor  ComputeFlavor
}

type KeyValue struct {
	Key   CustomizedKey
	Value CustomizedValue
}

type Input struct {
	Key   CustomizedKey
	Value ResourceInput
}

type ResourceInput struct {
	User   Account
	Type   ResourceType
	RepoId string
	File   string
}

type JobInfo struct {
	Endpoint  string
	JobId     string
	LogDir    string
	OutputDir string
}

type JobDetail struct {
	Status   string
	Duration int
}

type TrainingSummary struct {
	Name      TrainingName
	Desc      TrainingDesc
	JobDetail JobDetail
	CreatedAt int64
}
