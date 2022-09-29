package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func (col dataset) AddLike(owner, rid string) error {
	return updateResourceLike(col.collectionName, owner, rid, 1)
}

func (col dataset) RemoveLike(owner, rid string) error {
	return updateResourceLike(col.collectionName, owner, rid, -1)
}

func (col dataset) summaryFields() []string {
	return []string{
		fieldId, fieldName, fieldDesc, fieldTags,
		fieldUpdatedAt, fieldLikeCount, fieldDownloadCount,
	}
}

func (col dataset) ListAndSortByUpdateTime(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.DatasetSummaryDO, int, error) {
	return col.listResource(owner, do, sortByUpdateTime())
}

func (col dataset) ListAndSortByFirstLetter(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.DatasetSummaryDO, int, error) {
	return col.listResource(owner, do, sortByFirstLetter())
}

func (col dataset) ListAndSortByDownloadCount(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.DatasetSummaryDO, int, error) {
	return col.listResource(owner, do, sortByDownloadCount())
}

func (col dataset) listResource(
	owner string, do *repositories.ResourceListDO, sort bson.M,
) (r []repositories.DatasetSummaryDO, total int, err error) {
	var v []struct {
		Total int         `bson:"total"`
		Item  datasetItem `bson:"items"`
	}

	err = listResource(
		col.collectionName, owner, do, sort, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 {
		return
	}

	total = v[0].Total

	r = make([]repositories.DatasetSummaryDO, len(v))
	for i := range v {
		col.toDatasetSummaryDO(owner, &v[i].Item, &r[i])
	}

	return
}

func (col dataset) toDatasetSummaryDO(owner string, item *datasetItem, do *repositories.DatasetSummaryDO) {
	*do = repositories.DatasetSummaryDO{
		Id:            item.Id,
		Owner:         owner,
		Name:          item.Name,
		Desc:          item.Desc,
		Tags:          item.Tags,
		UpdatedAt:     item.UpdatedAt,
		LikeCount:     item.LikeCount,
		DownloadCount: item.DownloadCount,
	}
}
