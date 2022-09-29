package domain

type Dataset struct {
	Id string

	Owner    Account
	Protocol ProtocolName

	DatasetModifiableProperty

	RepoId string

	CreatedAt int64
	UpdatedAt int64

	Version int

	// following fileds is not under the controlling of version
	LikeCount     int
	DownloadCount int
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
