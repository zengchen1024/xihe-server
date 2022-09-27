package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

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

func (col project) IncreaseFork(owner, rid string) (err error) {
	updated := false

	f := func(ctx context.Context) error {
		updated, err = cli.updateArrayElemCount(
			ctx, col.collectionName, fieldItems, fieldForkCount, 1,
			resourceOwnerFilter(owner), arrayFilterById(rid),
		)

		return nil
	}

	if withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return
	}

	if !updated {
		err = repositories.NewErrorDataNotExists(errors.New("no update"))
	}

	return
}

func (col project) ListAndSortByUpdateTime(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ProjectSummaryDO, error) {
	return col.listResource(owner, do, sortByUpdateTime())
}

func (col project) ListAndSortByFirstLetter(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ProjectSummaryDO, error) {
	return col.listResource(owner, do, sortByFirstLetter())
}

func (col project) ListAndSortByDownloadCount(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ProjectSummaryDO, error) {
	return col.listResource(owner, do, sortByDownloadCount())
}

func (col project) listResource(
	owner string, do *repositories.ResourceListDO, sort bson.M,
) (r []repositories.ProjectSummaryDO, err error) {
	var v []dProject

	err = listResource(
		col.collectionName, owner, do, sort, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 {
		return
	}

	items := v[0].Items
	r = make([]repositories.ProjectSummaryDO, len(items))

	for i := range items {
		col.toProjectSummary(owner, &items[i], &r[i])
	}

	return
}

func (col project) summaryFields() []string {
	return []string{
		fieldId, fieldName, fieldDesc, fieldCoverId, fieldTags,
		fieldUpdatedAt, fieldLikeCount, fieldForkCount, fieldDownloadCount,
	}
}

func (col project) toProjectSummary(owner string, item *projectItem, do *repositories.ProjectSummaryDO) {
	*do = repositories.ProjectSummaryDO{
		Id:            item.Id,
		Owner:         owner,
		Name:          item.Name,
		Desc:          item.Desc,
		CoverId:       item.CoverId,
		Tags:          item.Tags,
		UpdatedAt:     item.UpdatedAt,
		LikeCount:     item.LikeCount,
		ForkCount:     item.ForkCount,
		DownloadCount: item.DownloadCount,
	}
}
