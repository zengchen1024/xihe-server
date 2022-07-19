package domain

type Dataset struct {
	Id string

	Owner    Account
	Name     ProjName
	Desc     ProjDesc
	RepoType RepoType
	Protocol ProtocolName

	Tags []string

	Version int
}
