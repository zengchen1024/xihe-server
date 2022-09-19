package domain

import (
	"errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
)

func ResourceTypeByName(n string) (ResourceType, error) {
	if strings.HasPrefix(n, ResourceProject) {
		return ResourceTypeProject, nil
	}

	if strings.HasPrefix(n, ResourceDataset) {
		return ResourceTypeDataset, nil
	}

	if strings.HasPrefix(n, ResourceModel) {
		return ResourceTypeModel, nil
	}

	return nil, errors.New("unknow resource")
}

// ResourceObj
type ResourceObj struct {
	ResourceOwner Account
	ResourceType  ResourceType
	ResourceId    string
}

func (r *ResourceObj) String() string {
	return fmt.Sprintf(
		"%s_%s_%s",
		r.ResourceOwner.Account(),
		r.ResourceType.ResourceType(),
		r.ResourceId,
	)
}

type ResourceIndex struct {
	ResourceOwner Account
	ResourceId    string
}

type RelatedResources []ResourceIndex

func (r RelatedResources) Has(owner Account, rid string) bool {
	v := sets.NewString()

	for i := range ([]ResourceIndex)(r) {
		v.Insert(
			r[i].ResourceOwner.Account() + r[i].ResourceId,
		)
	}

	return v.Has(owner.Account() + rid)
}
