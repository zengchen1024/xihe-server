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

	Version int

	// following fileds is not under the controlling of version
	LikeCount int
}

func (p *Project) MaxRelatedResourceNum() int {
	return config.Resource.MaxRelatedResourceNum
}

type ProjectModifiableProperty struct {
	Name     ProjName
	Desc     ResourceDesc
	CoverId  CoverId
	RepoType RepoType
	Tags     []string
}
