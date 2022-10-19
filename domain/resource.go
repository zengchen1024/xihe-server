package domain

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

func ResourceTypeByName(n string) (ResourceType, error) {
	v, err := NewResourceName(n)
	if err != nil {
		return nil, err
	}

	return v.ResourceType(), err
}

func NewResourceName(n string) (ResourceName, error) {
	if strings.HasPrefix(n, resourceProject) {
		return NewProjName(n)
	}

	if strings.HasPrefix(n, resourceDataset) {
		return NewDatasetName(n)
	}

	if strings.HasPrefix(n, resourceModel) {
		return NewModelName(n)
	}

	return nil, errors.New("unknow resource")
}

// ResourceObject
type ResourceObject struct {
	Type ResourceType

	ResourceIndex
}

func (r *ResourceObject) String() string {
	return fmt.Sprintf(
		"%s_%s_%s",
		r.Owner.Account(),
		r.Type.ResourceType(),
		r.Id,
	)
}

type ResourceIndex struct {
	Owner Account
	Id    string
}

type RelatedResources []ResourceIndex

func (r RelatedResources) Has(index *ResourceIndex) bool {
	v := sets.NewString()

	for i := range ([]ResourceIndex)(r) {
		v.Insert(
			r[i].Owner.Account() + r[i].Id,
		)
	}

	return v.Has(index.Owner.Account() + index.Id)
}

func (r RelatedResources) Count() int {
	return len(r)
}

type ReverselyRelatedResourceInfo struct {
	Promoter *ResourceIndex
	Resource *ResourceIndex
}

type ResourceSummary struct {
	Owner    Account
	Name     ResourceName
	Id       string
	RepoId   string
	RepoType RepoType
}

func (s *ResourceSummary) IsPrivate() bool {
	return s.RepoType.RepoType() == RepoTypePrivate
}
