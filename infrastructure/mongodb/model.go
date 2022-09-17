package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func modelDocFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
	}
}

func modelItemFilter(name string) bson.M {
	return bson.M{
		fieldName: name,
	}
}

func NewModelMapper(name string) repositories.ModelMapper {
	return model{name}
}

type model struct {
	collectionName string
}

func (col model) newDoc(owner string) error {
	docFilter := modelDocFilter(owner)

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

func (col model) Insert(do repositories.ModelDO) (identity string, err error) {
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

func (col model) insert(do repositories.ModelDO) (identity string, err error) {
	identity = newId()

	do.Id = identity
	doc, err := col.toModelDoc(&do)
	if err != nil {
		return
	}
	doc[fieldVersion] = 0
	doc[fieldLikeCount] = 0

	docFilter := modelDocFilter(do.Owner)

	appendElemMatchToFilter(
		fieldItems, false,
		modelItemFilter(do.Name), docFilter,
	)

	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, col.collectionName,
			fieldItems, docFilter, doc,
		)
	}

	err = withContext(f)

	return
}

func (col model) Update(do repositories.ModelDO) error {
	doc, err := col.toModelDoc(&do)
	if err != nil {
		return err
	}

	docFilter := modelDocFilter(do.Owner)

	updated := false

	f := func(ctx context.Context) error {
		b, err := cli.updateArrayElem(
			ctx, col.collectionName, fieldItems,
			docFilter, arrayFilterById(do.Id), doc, do.Version,
		)

		updated = b

		return err
	}

	if err := withContext(f); err != nil {
		return err
	}

	if !updated {
		return repositories.NewErrorConcurrentUpdating(errors.New("no update"))
	}

	return nil
}

func (col model) Get(owner, identity string) (do repositories.ModelDO, err error) {
	var v []dModel

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			modelDocFilter(owner), arrayFilterById(identity),
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

	col.toModelDO(owner, &v[0].Items[0], &do)

	return
}

func (col model) GetByName(owner, name string) (do repositories.ModelDO, err error) {
	var v []dModel

	if err = getResourceByName(col.collectionName, owner, name, &v); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toModelDO(owner, &v[0].Items[0], &do)

	return
}

func (col model) List(owner string, do repositories.ResourceListDO) (
	r []repositories.ModelDO, err error,
) {
	var v []dModel

	err = listResource(col.collectionName, owner, do, &v)
	if err != nil {
		return
	}

	if len(v) == 0 {
		return
	}

	items := v[0].Items
	r = make([]repositories.ModelDO, len(items))
	for i := range items {
		col.toModelDO(owner, &items[i], &r[i])
	}

	return
}

func (col model) ListUsersModels(opts map[string][]string) (
	r []repositories.ModelDO, err error,
) {
	var v []dModel

	err = listUsersResources(col.collectionName, opts, &v)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]repositories.ModelDO, 0, len(v))
	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		for j := range items {
			col.toModelDO(owner, &items[j], &r[i])
		}
	}

	return
}

func (col model) toModelDoc(do *repositories.ModelDO) (bson.M, error) {
	docObj := modelItem{
		Id:       do.Id,
		Name:     do.Name,
		Desc:     do.Desc,
		Protocol: do.Protocol,
		RepoType: do.RepoType,
		RepoId:   do.RepoId,
		Tags:     do.Tags,
	}

	return genDoc(docObj)
}

func (col model) toModelDO(owner string, item *modelItem, do *repositories.ModelDO) {
	*do = repositories.ModelDO{
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
