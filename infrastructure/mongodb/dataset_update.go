package mongodb

import "github.com/opensourceways/xihe-server/infrastructure/repositories"

func (col dataset) AddLike(owner, rid string) error {
	return updateResourceLike(col.collectionName, owner, rid, 1)
}

func (col dataset) RemoveLike(owner, rid string) error {
	return updateResourceLike(col.collectionName, owner, rid, -1)
}

func (col dataset) summaryFields() []string {
	return []string{
		fieldId, fieldName, fieldDesc, fieldTags,
		// fieldUpdatedAt,
		fieldLikeCount,
		// fieldDownloadCount,
	}
}

func (col dataset) toDatasetSummary(owner string, item *datasetItem, do *repositories.DatasetDO) {
	*do = repositories.DatasetDO{
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
