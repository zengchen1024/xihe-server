package domain

type Dataset struct {
	Id string

	Owner    Account
	Name     DatasetName
	Desc     ProjDesc
	RepoType RepoType
	Protocol ProtocolName

	Tags []string

	RepoId string

	Version int

	// following fileds is not under the controlling of version
	LikeCount int
}
