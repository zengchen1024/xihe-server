package domain

type Project struct {
	Id string

	Owner    Account
	Type     ProjType
	Protocol ProtocolName
	Training TrainingPlatform

	ProjectModifiableProperty

	RepoId string

	Version int

	// following fileds is not under the controlling of version
	LikeCount int
}

type ProjectModifiableProperty struct {
	Name     ProjName
	Desc     ResourceDesc
	CoverId  CoverId
	RepoType RepoType
	Tags     []string
}
