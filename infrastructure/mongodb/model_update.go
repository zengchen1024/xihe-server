package mongodb

import (
	"go.mongodb.org/mongo-driver/bson"

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

func (col model) ListAndSortByUpdateTime(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ModelSummaryDO, int, error) {
	return col.listResource(owner, do, sortByUpdateTime())
}

func (col model) ListAndSortByFirstLetter(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ModelSummaryDO, int, error) {
	return col.listResource(owner, do, sortByFirstLetter())
}

func (col model) ListAndSortByDownloadCount(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ModelSummaryDO, int, error) {
	return col.listResource(owner, do, sortByDownloadCount())
}

func (col model) listResource(
	owner string, do *repositories.ResourceListDO, sort bson.M,
) (r []repositories.ModelSummaryDO, total int, err error) {
	var v []struct {
		Total int       `bson:"total"`
		Item  modelItem `bson:"items"`
	}

	err = listResource(
		col.collectionName, owner, do, sort, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 {
		return
	}

	total = v[0].Total

	r = make([]repositories.ModelSummaryDO, len(v))
	for i := range v {
		col.toModelSummary(owner, &v[i].Item, &r[i])
	}

	return
}

func (col model) summaryFields() []string {
	return []string{
		fieldId, fieldName, fieldDesc, fieldTags,
		fieldUpdatedAt, fieldLikeCount, fieldDownloadCount,
	}
}

func (col model) toModelSummary(owner string, item *modelItem, do *repositories.ModelSummaryDO) {
	*do = repositories.ModelSummaryDO{
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
