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

func (d *Dataset) IsPrivate() bool {
	return d.RepoType.RepoType() == RepoTypePrivate
}

type DatasetModifiableProperty struct {
	Name     DatasetName
	Desc     ResourceDesc
	RepoType RepoType
	Tags     []string
}
