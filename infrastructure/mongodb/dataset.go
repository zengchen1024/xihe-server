package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func datasetDocFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
	}
}

func datasetItemFilter(name string) bson.M {
	return bson.M{
		fieldName: name,
	}
}

func NewDatasetMapper(name string) repositories.DatasetMapper {
	return dataset{name}
}

type dataset struct {
	collectionName string
}

func (col dataset) newDoc(owner string) error {
	docFilter := datasetDocFilter(owner)

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

func (col dataset) Insert(do repositories.DatasetDO) (identity string, err error) {
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

func (col dataset) insert(do repositories.DatasetDO) (identity string, err error) {
	identity = newId()

	do.Id = identity
	doc, err := col.toDatasetDoc(&do)
	if err != nil {
		return
	}
	doc[fieldVersion] = 0
	doc[fieldLikeCount] = 0

	err = insertResource(col.collectionName, do.Owner, do.Name, doc)

	return
}

func (col dataset) UpdateProperty(do *repositories.DatasetPropertyDO) error {
	p := &DatasetPropertyItem{
		FL:       do.FL,
		Name:     do.Name,
		Desc:     do.Desc,
		RepoType: do.RepoType,
		Tags:     do.Tags,
	}

	return updateResourceProperty(col.collectionName, &do.ResourceToUpdateDO, p)
}

func (col dataset) Get(owner, identity string) (do repositories.DatasetDO, err error) {
	var v []dDataset

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			datasetDocFilter(owner), arrayFilterById(identity),
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

	col.toDatasetDO(owner, &v[0].Items[0], &do)

	return
}

func (col dataset) GetByName(owner, name string) (do repositories.DatasetDO, err error) {
	var v []dDataset

	if err = getResourceByName(col.collectionName, owner, name, &v); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toDatasetDO(owner, &v[0].Items[0], &do)

	return
}

func (col dataset) List(owner string, do repositories.ResourceListDO) (
	r []repositories.DatasetDO, err error,
) {
	var v []dDataset

	err = listResource(
		col.collectionName, owner, &do, nil, col.summaryFields(), &v,
	)
	if err != nil {
		return
	}

	if len(v) == 0 {
		return
	}

	items := v[0].Items
	r = make([]repositories.DatasetDO, len(items))
	for i := range items {
		col.toDatasetDO(owner, &items[i], &r[i])
	}

	return
}

func (col dataset) ListUsersDatasets(opts map[string][]string) (
	r []repositories.DatasetDO, err error,
) {
	var v []dDataset

	err = listUsersResources(col.collectionName, opts, &v)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]repositories.DatasetDO, 0, len(v))

	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		dos := make([]repositories.DatasetDO, len(items))
		for j := range items {
			col.toDatasetDO(owner, &items[j], &dos[j])
		}

		r = append(r, dos...)
	}

	return
}

func (col dataset) toDatasetDoc(do *repositories.DatasetDO) (bson.M, error) {
	docObj := datasetItem{
		Id:       do.Id,
		RepoId:   do.RepoId,
		Protocol: do.Protocol,
		DatasetPropertyItem: DatasetPropertyItem{
			FL:       do.FL,
			Name:     do.Name,
			Desc:     do.Desc,
			RepoType: do.RepoType,
			Tags:     do.Tags,
		},
	}

	return genDoc(docObj)
}

func (col dataset) toDatasetDO(owner string, item *datasetItem, do *repositories.DatasetDO) {
	*do = repositories.DatasetDO{
		Id:        item.Id,
		Owner:     owner,
		Name:      item.Name,
		Desc:      item.Desc,
		Protocol:  item.Protocol,
		RepoType:  item.RepoType,
		RepoId:    item.RepoId,
		Tags:      item.Tags,
		Version:   item.Version,
		LikeCount: item.LikeCount,
	}
}
