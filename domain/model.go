package domain

type Model struct {
	Id    string
	Owner string

	Name     ProjName
	Desc     ProjDesc
	RepoType RepoType
	Protocol ProtocolName

	Tags []string

	Version int
}
