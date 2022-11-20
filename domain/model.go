package domain

type Model struct {
	Id string

	Owner    Account
	Protocol ProtocolName

	ModelModifiableProperty

	RepoId string

	RelatedDatasets RelatedResources

	CreatedAt int64
	UpdatedAt int64

	Version int

	// following fileds is not under the controlling of version
	LikeCount       int
	DownloadCount   int
	RelatedProjects RelatedResources
}

func (m *Model) MaxRelatedResourceNum() int {
	return config.MaxRelatedResourceNum
}

func (m *Model) IsPrivate() bool {
	return m.RepoType.RepoType() == RepoTypePrivate
}

func (p *Model) ResourceIndex() ResourceIndex {
	return ResourceIndex{
		Owner: p.Owner,
		Id:    p.Id,
	}
}

func (p *Model) ResourceObject() ResourceObject {
	return ResourceObject{
		Type:          ResourceTypeModel,
		ResourceIndex: p.ResourceIndex(),
	}
}

func (p *Model) RelatedResources() []ResourceObjects {
	r := make([]ResourceObjects, 0, 2)

	if len(p.RelatedProjects) > 0 {
		r = append(r, ResourceObjects{
			Type:    ResourceTypeProject,
			Objects: p.RelatedProjects,
		})
	}

	if len(p.RelatedDatasets) > 0 {
		r = append(r, ResourceObjects{
			Type:    ResourceTypeDataset,
			Objects: p.RelatedDatasets,
		})
	}

	return r
}

type ModelModifiableProperty struct {
	Name     ResourceName
	Desc     ResourceDesc
	RepoType RepoType
	Tags     []string
	TagKinds []string
}

type ModelSummary struct {
	Id            string
	Owner         Account
	Name          ResourceName
	Desc          ResourceDesc
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}
