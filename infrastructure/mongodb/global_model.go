package mongodb

import (
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

type globalModel struct {
	owner string
	*modelItem
}

func (col model) ListGlobalAndSortByUpdateTime(do *repositories.GlobalResourceListDO) (
	[]repositories.ModelSummaryDO, int, error,
) {

	f := func(items []globalModel) []globalModel {
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

		r := make([]globalModel, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listGlobalResource(do, f)
}

func (col model) ListGlobalAndSortByFirstLetter(do *repositories.GlobalResourceListDO) (
	[]repositories.ModelSummaryDO, int, error,
) {

	f := func(items []globalModel) []globalModel {
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

		r := make([]globalModel, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listGlobalResource(do, f)
}

func (col model) ListGlobalAndSortByDownloadCount(do *repositories.GlobalResourceListDO) (
	[]repositories.ModelSummaryDO, int, error,
) {

	f := func(items []globalModel) []globalModel {
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

		r := make([]globalModel, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listGlobalResource(do, f)
}

func (col model) listGlobalResource(
	do *repositories.GlobalResourceListDO,
	sortAndPagination func(items []globalModel) []globalModel,
) (r []repositories.ModelSummaryDO, total int, err error) {
	var v []dModel

	err = listGlobalResourceWithoutSort(
		col.collectionName, do, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 {
		return
	}

	items := col.toGlobalModels(v)

	total = len(items)

	if items = sortAndPagination(items); len(items) == 0 {
		return
	}

	r = make([]repositories.ModelSummaryDO, len(items))
	for i := range items {
		col.toModelSummaryDO(items[i].owner, items[i].modelItem, &r[i])
	}

	return
}

func (col model) toGlobalModels(v []dModel) []globalModel {
	n := 0
	for i := range v {
		n += len(v[i].Items)
	}

	k := 0
	r := make([]globalModel, n)

	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		for j := range items {
			r[k] = globalModel{
				owner:     owner,
				modelItem: &items[j],
			}
			k++
		}
	}

	return r
}

func (col model) Search(do *repositories.GlobalResourceListDO, topNum int) (
	r []repositories.ResourceSummaryDO, total int, err error,
) {
	var v []dModel

	err = listGlobalResourceWithoutSort(
		col.collectionName, do, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 {
		return
	}

	items := col.toGlobalModels(v)

	total = len(items)

	if total < topNum {
		topNum = total
	}

	r = make([]repositories.ResourceSummaryDO, topNum)

	for i := range r {
		r[i].Owner = items[i].owner
		r[i].Name = items[i].modelItem.Name
	}

	return
}
