package domain

type Project struct {
	Id    string
	Owner string

	Name      ProjName
	Desc      ProjDesc
	Type      RepoType
	CoverId   CoverId
	Protocol  ProtocolName
	Training  TrainingSDK
	Inference InferenceSDK
	Tags      []string

	Version int
}
