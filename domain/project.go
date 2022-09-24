package domain

type Project struct {
	Id string

	Owner    Account
	Type     ProjType
	Protocol ProtocolName
	Training TrainingPlatform

	ProjectModifiableProperty

	RepoId string

	RelatedModels   RelatedResources
	RelatedDatasets RelatedResources

	CreatedAt int64
	UpdatedAt int64

	Version int

	// following fileds is not under the controlling of version
	LikeCount int
	ForkCount int
}

func (p *Project) MaxRelatedResourceNum() int {
	return config.MaxRelatedResourceNum
}

type ProjectModifiableProperty struct {
	Name     ProjName
	Desc     ResourceDesc
	CoverId  CoverId
	RepoType RepoType
	Tags     []string
}
