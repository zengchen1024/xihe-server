package domain

type Model struct {
	Id string

	Owner    Account
	Protocol ProtocolName

	ModelModifiableProperty

	RepoId string

	Version int

	// following fileds is not under the controlling of version
	LikeCount int
}

type ModelModifiableProperty struct {
	Name     ModelName
	Desc     ProjDesc
	RepoType RepoType
	Tags     []string
}
