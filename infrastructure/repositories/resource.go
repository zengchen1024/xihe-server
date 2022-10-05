package repositories

import (
	"github.com/opensourceways/xihe-server/domain"
	"github.com/opensourceways/xihe-server/domain/repository"
)

type ResourceListDO struct {
	Name         string
	RepoType     string
	PageNum      int
	CountPerPage int
}

func toResourceListDO(r *repository.ResourceListOption) ResourceListDO {
	do := ResourceListDO{
		Name:         r.Name,
		PageNum:      r.PageNum,
		CountPerPage: r.CountPerPage,
	}

	if r.RepoType != nil {
		do.RepoType = r.RepoType.RepoType()
	}

	return do
}

type ResourceObjectDO struct {
	Owner string
	Type  string
	Id    string
}

func (do *ResourceObjectDO) toResourceObject(r *domain.ResourceObject) (err error) {
	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	if r.Type, err = domain.NewResourceType(do.Type); err != nil {
		return
	}

	r.Id = do.Id

	return
}

func toResourceObjectDO(r *domain.ResourceObject) ResourceObjectDO {
	return ResourceObjectDO{
		Owner: r.Owner.Account(),
		Type:  r.Type.ResourceType(),
		Id:    r.Id,
	}
}

type ResourceIndexDO struct {
	Owner string
	Id    string
}

func (do *ResourceIndexDO) toResourceIndex(r *domain.ResourceIndex) (err error) {
	if r.Owner, err = domain.NewAccount(do.Owner); err != nil {
		return
	}

	r.Id = do.Id

	return
}

func toResourceIndexDO(r *domain.ResourceIndex) ResourceIndexDO {
	return ResourceIndexDO{
		Owner: r.Owner.Account(),
		Id:    r.Id,
	}
}

func convertToResourceIndex(v []ResourceIndexDO) (r []domain.ResourceIndex, err error) {
	if len(v) == 0 {
		return
	}

	r = make([]domain.ResourceIndex, len(v))

	for i := range v {
		if err = v[i].toResourceIndex(&r[i]); err != nil {
			return
		}
	}

	return
}

func toReverselyRelatedResourceInfoDO(
	info *domain.ReverselyRelatedResourceInfo,
) ReverselyRelatedResourceInfoDO {
	return ReverselyRelatedResourceInfoDO{
		Promoter: toResourceIndexDO(info.Promoter),
		Resource: toResourceIndexDO(info.Resource),
	}
}

type ReverselyRelatedResourceInfoDO struct {
	Promoter ResourceIndexDO
	Resource ResourceIndexDO
}
