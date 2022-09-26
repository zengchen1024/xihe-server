package mongodb

import (
	"context"
	"errors"

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
) ([]repositories.ProjectDO, error) {
	return col.listResource(owner, func() ([]dProject, error) {
		var v []dProject

		err := listResourceAndSortByUpdateTime(col.collectionName, owner, do, &v)

		return v, err
	})
}

func (col project) ListAndSortByFirtLetter(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ProjectDO, error) {
	return col.listResource(owner, func() ([]dProject, error) {
		var v []dProject

		err := listResourceAndSortByFirtLetter(col.collectionName, owner, do, &v)

		return v, err
	})
}

func (col project) ListAndSortByDownloadCount(
	owner string, do *repositories.ResourceListDO,
) ([]repositories.ProjectDO, error) {
	return col.listResource(owner, func() ([]dProject, error) {
		var v []dProject

		err := listResourceAndSortByDownloadCount(col.collectionName, owner, do, &v)

		return v, err
	})
}

func (col project) listResource(
	owner string, f func() ([]dProject, error),
) (r []repositories.ProjectDO, err error) {
	v, err := f()
	if err != nil {
		return
	}

	if len(v) == 0 {
		return
	}

	items := v[0].Items
	r = make([]repositories.ProjectDO, len(items))
	for i := range items {
		col.toProjectSummary(owner, &items[i], &r[i])
	}

	return
}

func (col project) toProjectSummary(owner string, item *projectItem, do *repositories.ProjectDO) {
	*do = repositories.ProjectDO{
		Id:            item.Id,
		Owner:         owner,
		Name:          item.Name,
		Desc:          item.Desc,
		Type:          item.Type,
		CoverId:       item.CoverId,
		Protocol:      item.Protocol,
		Training:      item.Training,
		RepoType:      item.RepoType,
		RepoId:        item.RepoId,
		Tags:          item.Tags, // TODO need this?
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		Version:       item.Version,
		LikeCount:     item.LikeCount,
		ForkCount:     item.ForkCount,
		DownloadCount: item.DownloadCount,
	}
}
