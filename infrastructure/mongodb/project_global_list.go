package mongodb

import (
	"fmt"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

type globalProject struct {
	owner string
	*projectItem
}

func (col project) GlobalListAndSortByUpdateTime(do *repositories.GlobalResourceListDO) (
	[]repositories.ProjectSummaryDO, int, error,
) {

	f := func(items []globalProject) []globalProject {
		v := make([]updateAtSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = updateAtSortData{
				id:       item.Id,
				index:    i,
				updateAt: item.UpdatedAt,
			}
		}

		v = updateAtSortAndPaginate(v, do.CountPerPage, do.PageNum)
		if len(v) == 0 {
			return nil
		}

		r := make([]globalProject, len(v))
		for i := range v {
			r[i] = items[v[i].index]
		}

		return r
	}

	return col.globalListResource(do, f)
}

func (col project) GlobalListAndSortByFirstLetter(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ProjectSummaryDO, int, error) {

	f := func(items []projectItem) []projectItem {
		v := make([]firstLetterSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = firstLetterSortData{
				index:    i,
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

func (col project) GlobalListAndSortByDownloadCount(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ProjectSummaryDO, int, error) {

	f := func(items []projectItem) []projectItem {
		v := make([]downloadSortData, len(items))

		for i := range items {
			item := &items[i]

			v[i] = downloadSortData{
				index:    i,
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

func (col project) globalListResource(
	do *repositories.GlobalResourceListDO,
	sortAndPagination func(items []globalProject) []globalProject,
) (r []repositories.ProjectSummaryDO, total int, err error) {
	var v []dProject

	err = globalListResourceWithoutSort(
		col.collectionName, do, col.summaryFields(), &v,
	)

	if err != nil || len(v) == 0 {
		return
	}

	fmt.Println(v)

	items := col.toGlobalProjects(v)

	total = len(items)

	if items = sortAndPagination(items); len(items) == 0 {
		return
	}

	r = make([]repositories.ProjectSummaryDO, len(items))
	for i := range items {
		col.toProjectSummaryDO(items[i].owner, items[i].projectItem, &r[i])
	}

	return
}

func (col project) toGlobalProjects(v []dProject) []globalProject {
	n := 0
	for i := range v {
		n += len(v[i].Items)
	}

	k := 0
	r := make([]globalProject, n)

	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		for j := range items {
			r[k] = globalProject{
				owner:       owner,
				projectItem: &items[j],
			}
			k++
		}
	}

	return r
}
