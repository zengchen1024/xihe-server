package domain

type Project struct {
	Id    string
	Owner string

	Name     ProjName
	Desc     ProjDesc
	Type     ProjType
	CoverId  CoverId
	RepoType RepoType
	Protocol ProtocolName
	Training TrainingPlatform

	Tags []string

	Version int
}
