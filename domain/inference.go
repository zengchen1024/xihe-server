package domain

type Infereance struct {
	Id string

	InferenceIndex

	ProjectName   ProjName
	ProjectRepoId string

	ModelRef ResourceRef

	// following fileds is not under the controlling of version
	InferenceDetail
}

type InferenceIndex struct {
	ProjectId    string
	LastCommit   string
	ProjectOwner Account
}

type InferenceInfo struct {
	Id           string
	ProjectId    string
	ProjectOwner Account
}

type InferenceDetail struct {
	// Expiry stores the time when the inference instance will exit.
	Expiry int64

	// Error stores the message when the reference instance starts failed
	Error string

	// AccessURL stores the url to access the inference service.
	AccessURL string
}
