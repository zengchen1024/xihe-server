package mongodb

import (
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func (col model) AddLike(owner, rid string) error {
	return updateResourceLike(col.collectionName, owner, rid, 1)
}

func (col model) RemoveLike(owner, rid string) error {
	return updateResourceLike(col.collectionName, owner, rid, -1)
}

func (col model) AddRelatedDataset(do *repositories.RelatedResourceDO) error {
	return updateRelatedResource(col.collectionName, fieldDatasets, true, do)
}

func (col model) RemoveRelatedDataset(do *repositories.RelatedResourceDO) error {
	return updateRelatedResource(col.collectionName, fieldDatasets, false, do)
}

func (col model) summaryFields() []string {
	return []string{
		fieldId, fieldName, fieldDesc, fieldTags,
		// fieldUpdatedAt,
		fieldLikeCount,
		// fieldDownloadCount,
	}
}

func (col model) toModelSummary(owner string, item *modelItem, do *repositories.ModelDO) {
	*do = repositories.ModelDO{
		Id:    item.Id,
		Owner: owner,
		Name:  item.Name,
		Desc:  item.Desc,
		Tags:  item.Tags,
		//UpdatedAt:     item.UpdatedAt,
		LikeCount: item.LikeCount,
		//DownloadCount: item.DownloadCount,
	}
}
