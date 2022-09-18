package domain

type Dataset struct {
	Id string

	Owner    Account
	Protocol ProtocolName

	DatasetModifiableProperty

	RepoId string

	Version int

	// following fileds is not under the controlling of version
	LikeCount int
}

type DatasetModifiableProperty struct {
	Name     DatasetName
	Desc     ProjDesc
	RepoType RepoType
	Tags     []string
}
