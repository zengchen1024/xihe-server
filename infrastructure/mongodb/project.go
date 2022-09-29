package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func projectDocFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
	}
}

func projectItemFilter(name string) bson.M {
	return bson.M{
		fieldName: name,
	}
}

func arrayFilterById(identity string) bson.M {
	return bson.M{
		fieldId: identity,
	}
}

func NewProjectMapper(name string) repositories.ProjectMapper {
	return project{name}
}

type project struct {
	collectionName string
}

func (col project) newDoc(owner string) error {
	docFilter := projectDocFilter(owner)

	doc := bson.M{
		fieldOwner: owner,
		fieldItems: bson.A{},
	}

	f := func(ctx context.Context) error {
		_, err := cli.newDocIfNotExist(
			ctx, col.collectionName, docFilter, doc,
		)

		return err
	}

	if err := withContext(f); err != nil && isDBError(err) {
		return err
	}

	return nil
}

func (col project) Insert(do repositories.ProjectDO) (identity string, err error) {
	if identity, err = col.insert(do); err == nil || !isDocNotExists(err) {
		return
	}

	// doc is not exist or duplicate insert

	if err = col.newDoc(do.Owner); err == nil {
		if identity, err = col.insert(do); err != nil && isDocNotExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}
	}

	return
}

func (col project) insert(do repositories.ProjectDO) (identity string, err error) {
	identity = newId()

	do.Id = identity
	doc, err := col.toProjectDoc(&do)
	if err != nil {
		return
	}
	doc[fieldVersion] = 0
	doc[fieldLikeCount] = 0
	doc[fieldModels] = bson.A{}
	doc[fieldDatasets] = bson.A{}

	err = insertResource(col.collectionName, do.Owner, do.Name, doc)

	return
}

func (col project) UpdateProperty(do *repositories.ProjectPropertyDO) error {
	p := &ProjectPropertyItem{
		Name:     do.Name,
		FL:       do.FL,
		Desc:     do.Desc,
		CoverId:  do.CoverId,
		RepoType: do.RepoType,
		Tags:     do.Tags,
	}

	return updateResourceProperty(col.collectionName, &do.ResourceToUpdateDO, p)
}

func (col project) Get(owner, identity string) (do repositories.ProjectDO, err error) {
	var v []dProject

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			projectDocFilter(owner), arrayFilterById(identity),
			bson.M{fieldItems: 1}, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toProjectDO(owner, &v[0].Items[0], &do)

	return
}

func (col project) GetByName(owner, name string) (do repositories.ProjectDO, err error) {
	var v []dProject

	if err = getResourceByName(col.collectionName, owner, name, &v); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toProjectDO(owner, &v[0].Items[0], &do)

	return
}

func (col project) List(owner string, do *repositories.ResourceListDO) (
	[]repositories.ProjectSummaryDO, int, error,
) {
	return col.listResource(owner, do, nil)
}

func (col project) ListUsersProjects(opts map[string][]string) (
	r []repositories.ProjectSummaryDO, err error,
) {
	var v []dProject

	err = listUsersResources(col.collectionName, opts, &v)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]repositories.ProjectSummaryDO, 0, len(v))

	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		dos := make([]repositories.ProjectSummaryDO, len(items))
		for j := range items {
			col.toProjectSummaryDO(owner, &items[j], &dos[j])
		}

		r = append(r, dos...)
	}

	return
}

func (col project) toProjectDoc(do *repositories.ProjectDO) (bson.M, error) {
	docObj := projectItem{
		Id:        do.Id,
		Type:      do.Type,
		Protocol:  do.Protocol,
		Training:  do.Training,
		RepoId:    do.RepoId,
		CreatedAt: do.CreatedAt,
		UpdatedAt: do.UpdatedAt,
		ProjectPropertyItem: ProjectPropertyItem{
			FL:       do.FL,
			Name:     do.Name,
			Desc:     do.Desc,
			CoverId:  do.CoverId,
			RepoType: do.RepoType,
			Tags:     do.Tags,
		},
	}

	return genDoc(docObj)
}

func (col project) toProjectDO(owner string, item *projectItem, do *repositories.ProjectDO) {
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
		Tags:          item.Tags,
		CreatedAt:     item.CreatedAt,
		UpdatedAt:     item.UpdatedAt,
		Version:       item.Version,
		LikeCount:     item.LikeCount,
		ForkCount:     item.ForkCount,
		DownloadCount: item.DownloadCount,

		RelatedModels:   toResourceIndexDO(item.RelatedModels),
		RelatedDatasets: toResourceIndexDO(item.RelatedDatasets),
	}
}
