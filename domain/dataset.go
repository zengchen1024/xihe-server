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

	RelatedModels   RelatedResources
	RelatedProjects RelatedResources
}

func (d *Dataset) IsPrivate() bool {
	return d.RepoType.RepoType() == RepoTypePrivate
}

func (p *Dataset) ResourceIndex() ResourceIndex {
	return ResourceIndex{
		Owner: p.Owner,
		Id:    p.Id,
	}
}

func (p *Dataset) ResourceObject() ResourceObject {
	return ResourceObject{
		Type:          ResourceTypeDataset,
		ResourceIndex: p.ResourceIndex(),
	}
}

func (p *Dataset) RelatedResources() []ResourceObjects {
	r := make([]ResourceObjects, 0, 2)

	if len(p.RelatedProjects) > 0 {
		r = append(r, ResourceObjects{
			Type:    ResourceTypeProject,
			Objects: p.RelatedProjects,
		})
	}

	if len(p.RelatedModels) > 0 {
		r = append(r, ResourceObjects{
			Type:    ResourceTypeModel,
			Objects: p.RelatedModels,
		})
	}

	return r
}

type DatasetModifiableProperty struct {
	Name     ResourceName
	Desc     ResourceDesc
	RepoType RepoType
	Tags     []string
	TagKinds []string
}

type DatasetSummary struct {
	Id            string
	Owner         Account
	Name          ResourceName
	Desc          ResourceDesc
	Tags          []string
	UpdatedAt     int64
	LikeCount     int
	DownloadCount int
}
