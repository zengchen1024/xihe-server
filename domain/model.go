package domain

type Model struct {
	Id string

	Owner    Account
	Name     ProjName
	Desc     ProjDesc
	RepoType RepoType
	Protocol ProtocolName

	Tags []string

	RepoId string

	Version int
}
