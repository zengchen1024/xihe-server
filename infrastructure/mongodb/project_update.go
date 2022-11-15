package mongodb

import (
	"context"
	"errors"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func (col project) AddLike(p repositories.ResourceIndexDO) error {
	return updateResourceLike(col.collectionName, &p, 1)
}

func (col project) RemoveLike(p repositories.ResourceIndexDO) error {
	return updateResourceLike(col.collectionName, &p, -1)
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

func (col project) IncreaseFork(r repositories.ResourceIndexDO) (err error) {
	updated := false

	f := func(ctx context.Context) error {
		updated, err = cli.updateArrayElemCount(
			ctx, col.collectionName, fieldItems, fieldForkCount, 1,
			resourceOwnerFilter(r.Owner), resourceIdFilter(r.Id),
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
) ([]repositories.ProjectSummaryDO, int, error) {

	f := func(items []projectItem) []projectItem {
		v := make([]updateAtSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = updateAtSortData{
				id:       item.Id,
				level:    item.Level,
				index:    i,
				updateAt: item.UpdatedAt,
			}
		}

		v = updateAtSortAndPaginate(v, do.CountPerPage, do.PageNum)
		if len(v) == 0 {
			return nil
		}

		r := make([]projectItem, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listResource(owner, do, f)
}

func (col project) ListAndSortByFirstLetter(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ProjectSummaryDO, int, error) {

	f := func(items []projectItem) []projectItem {
		v := make([]firstLetterSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = firstLetterSortData{
				index:    i,
				level:    item.Level,
				letter:   item.FL,
				updateAt: item.UpdatedAt,
			}
		}

		v = firstLetterSortAndPaginate(v, do.CountPerPage, do.PageNum)
		if len(v) == 0 {
			return nil
		}

		r := make([]projectItem, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listResource(owner, do, f)
}

func (col project) ListAndSortByDownloadCount(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ProjectSummaryDO, int, error) {

	f := func(items []projectItem) []projectItem {
		v := make([]downloadSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = downloadSortData{
				index:    i,
				level:    item.Level,
				download: item.DownloadCount,
				updateAt: item.UpdatedAt,
			}
		}

		v = downloadSortAndPaginate(v, do.CountPerPage, do.PageNum)
		if len(v) == 0 {
			return nil
		}

		r := make([]projectItem, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listResource(owner, do, f)
}

func (col project) listResource(
	owner string,
	do *repositories.ResourceListDO,
	sortAndPagination func(items []projectItem) []projectItem,
) (r []repositories.ProjectSummaryDO, total int, err error) {
	var v []dProject

	err = listResourceWithoutSort(
		col.collectionName, owner, do, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 || len(v[0].Items) == 0 {
		return
	}

	items := v[0].Items
	total = len(items)

	items = sortAndPagination(items)
	if len(items) == 0 {
		return
	}

	r = make([]repositories.ProjectSummaryDO, len(items))
	for i := range items {
		col.toProjectSummaryDO(owner, &items[i], &r[i])
	}

	return
}

func (col project) summaryFields() []string {
	return []string{
		fieldId, fieldName, fieldDesc, fieldCoverId, fieldTags, fieldFirstLetter,
		fieldUpdatedAt, fieldLikeCount, fieldForkCount, fieldDownloadCount, fieldLevel,
	}
}

func (col project) toProjectSummaryDO(owner string, item *projectItem, do *repositories.ProjectSummaryDO) {
	*do = repositories.ProjectSummaryDO{
		Id:            item.Id,
		Owner:         owner,
		Name:          item.Name,
		Desc:          item.Desc,
		Level:         item.Level,
		CoverId:       item.CoverId,
		Tags:          item.Tags,
		UpdatedAt:     item.UpdatedAt,
		LikeCount:     item.LikeCount,
		ForkCount:     item.ForkCount,
		DownloadCount: item.DownloadCount,
	}
}
