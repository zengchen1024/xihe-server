package domain

type Project struct {
	Id string

	Owner    Account
	Name     ProjName
	Desc     ProjDesc
	Type     ProjType
	CoverId  CoverId
	RepoType RepoType
	Protocol ProtocolName
	Training TrainingPlatform

	Tags []string

	RepoId string

	Version int
}
