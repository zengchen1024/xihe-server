package mongodb

import (
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func (col project) AddLike(owner, pid string) error {
	return updateResourceLike(col.collectionName, owner, pid, 1)
}

func (col project) RemoveLike(owner, pid string) error {
	return updateResourceLike(col.collectionName, owner, pid, -1)
}

func (col project) AddRelatedModel(do *repositories.RelatedResourceDO) error {
	return updateRelatedResource(col.collectionName, fieldModels, true, do)
}

func (col project) RemoveRelatedModel(do *repositories.RelatedResourceDO) error {
	return updateRelatedResource(col.collectionName, fieldModels, false, do)
}

func (col project) AddRelatedDataset(do *repositories.RelatedResourceDO) error {
	return updateRelatedResource(col.collectionName, fieldDatasets, true, do)
}

func (col project) RemoveRelatedDataset(do *repositories.RelatedResourceDO) error {
	return updateRelatedResource(col.collectionName, fieldDatasets, false, do)
}
