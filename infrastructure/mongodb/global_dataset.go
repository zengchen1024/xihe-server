package mongodb

import (
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

type globalDataset struct {
	owner string
	*datasetItem
}

func (col dataset) ListGlobalAndSortByUpdateTime(do *repositories.GlobalResourceListDO) (
	[]repositories.DatasetSummaryDO, int, error,
) {

	f := func(items []globalDataset) []globalDataset {
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

		r := make([]globalDataset, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listGlobalResource(do, f)
}

func (col dataset) ListGlobalAndSortByFirstLetter(do *repositories.GlobalResourceListDO) (
	[]repositories.DatasetSummaryDO, int, error,
) {

	f := func(items []globalDataset) []globalDataset {
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

		r := make([]globalDataset, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listGlobalResource(do, f)
}

func (col dataset) ListGlobalAndSortByDownloadCount(do *repositories.GlobalResourceListDO) (
	[]repositories.DatasetSummaryDO, int, error,
) {

	f := func(items []globalDataset) []globalDataset {
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

		r := make([]globalDataset, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.listGlobalResource(do, f)
}

func (col dataset) listGlobalResource(
	do *repositories.GlobalResourceListDO,
	sortAndPagination func(items []globalDataset) []globalDataset,
) (r []repositories.DatasetSummaryDO, total int, err error) {
	var v []dDataset

	err = listGlobalResourceWithoutSort(
		col.collectionName, do, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 {
		return
	}

	items := col.toGlobalDatasets(v)

	total = len(items)

	if items = sortAndPagination(items); len(items) == 0 {
		return
	}

	r = make([]repositories.DatasetSummaryDO, len(items))
	for i := range items {
		col.toDatasetSummaryDO(items[i].owner, items[i].datasetItem, &r[i])
	}

	return
}

func (col dataset) toGlobalDatasets(v []dDataset) []globalDataset {
	n := 0
	for i := range v {
		n += len(v[i].Items)
	}

	k := 0
	r := make([]globalDataset, n)

	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		for j := range items {
			r[k] = globalDataset{
				owner:       owner,
				datasetItem: &items[j],
			}
			k++
		}
	}

	return r
}

func (col dataset) Search(do *repositories.GlobalResourceListDO, topNum int) (
	r []repositories.ResourceSummaryDO, total int, err error,
) {
	var v []dDataset

	err = listGlobalResourceWithoutSort(
		col.collectionName, do, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 {
		return
	}

	items := col.toGlobalDatasets(v)

	total = len(items)

	r = make([]repositories.ResourceSummaryDO, total)

	j := 0
	for i := range items {
		r[i].Owner = items[i].owner
		r[i].Name = items[i].datasetItem.Name

		if j++; j >= topNum {
			r = r[:topNum]

			break
		}
	}

	return
}
