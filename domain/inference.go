package domain

type Infereance struct {
	Id string

	InferenceIndex

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
	Id string

	InferenceIndex
}

type InferenceDetail struct {
	// Expiry stores the time when the inference instance will exit.
	Expiry int64

	// Error stores the message when the reference instance starts failed
	Error string

	// AccessURL stores the url to access the inference service.
	AccessURL string
}

func (d *InferenceDetail) IsCreating() bool {
	return d.Expiry == 0 && d.Error == "" && d.AccessURL == ""
}
