package repositories

import "github.com/opensourceways/xihe-server/domain"

type ResourceListDO struct {
	Name     string
	RepoType string
}

type ResourceObjDO struct {
	ResourceOwner string
	ResourceType  string
	ResourceId    string
}

func (do *ResourceObjDO) toResourceObj(r *domain.ResourceObj) (err error) {
	if r.ResourceOwner, err = domain.NewAccount(do.ResourceOwner); err != nil {
		return
	}

	if r.ResourceType, err = domain.NewResourceType(do.ResourceType); err != nil {
		return
	}

	r.ResourceId = do.ResourceId

	return
}

func toResourceObjDO(r *domain.ResourceObj) ResourceObjDO {
	return ResourceObjDO{
		ResourceOwner: r.ResourceOwner.Account(),
		ResourceType:  r.ResourceType.ResourceType(),
		ResourceId:    r.ResourceId,
	}
}
