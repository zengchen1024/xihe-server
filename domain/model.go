package domain

type Model struct {
	Id string

	Owner    Account
	Protocol ProtocolName

	ModelModifiableProperty

	RepoId string

	RelatedDatasets RelatedResources

	Version int

	// following fileds is not under the controlling of version
	LikeCount int
}

func (m *Model) MaxRelatedResourceNum() int {
	return config.MaxRelatedResourceNum
}

type ModelModifiableProperty struct {
	Name     ModelName
	Desc     ResourceDesc
	RepoType RepoType
	Tags     []string
}
