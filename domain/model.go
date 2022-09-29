package domain

type Model struct {
	Id string

	Owner    Account
	Protocol ProtocolName

	ModelModifiableProperty

	RepoId string

	RelatedDatasets RelatedResources

	CreatedAt int64
	UpdatedAt int64

	Version int

	// following fileds is not under the controlling of version
	LikeCount     int
	DownloadCount int
}

func (m *Model) MaxRelatedResourceNum() int {
	return config.MaxRelatedResourceNum
}

func (m *Model) IsPrivate() bool {
	return m.RepoType.RepoType() == RepoTypePrivate
}

type ModelModifiableProperty struct {
	Name     ModelName
	Desc     ResourceDesc
	RepoType RepoType
	Tags     []string
}

type ModelSummary struct {
	Id            string
	Owner         Account
	Name          ModelName
	Desc          ResourceDesc
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}
