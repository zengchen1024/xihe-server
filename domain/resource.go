package domain

import (
	"errors"
	"fmt"
	"strings"
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
