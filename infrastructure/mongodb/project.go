package mongodb

import (
	"context"
	"errors"

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

func projectItemFilterById(identity string) bson.M {
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

func (col project) New(owner string) error {
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

	if err := withContext(f); err != nil && !errors.Is(err, errDocExists) {
		return err
	}

	return nil
}

func (col project) Insert(do repositories.ProjectDO) (identity string, err error) {
	identity = newId()

	docObj := projectItem{
		Id:       identity,
		Name:     do.Name,
		Desc:     do.Desc,
		Type:     do.Type,
		CoverId:  do.CoverId,
		Protocol: do.Protocol,
		Training: do.Training,
		RepoType: do.RepoType,
		Tags:     do.Tags,
	}

	doc, err := genDoc(docObj)
	if err != nil {
		return
	}

	docFilter := projectDocFilter(do.Owner)

	appendElemMatchToFilter(
		fieldItems, false,
		projectItemFilter(do.Name), docFilter,
	)

	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, col.collectionName,
			fieldItems, docFilter, doc,
		)
	}

	err = withContext(f)

	if errors.Is(err, errDocNotExists) {
		err = repositories.NewErrorDuplicateCreating(err)
	}

	return
}

func (col project) Update(string, repositories.ProjectDO) error {
	return nil
}

func (col project) Get(owner, identity string) (do repositories.ProjectDO, err error) {
	var v []dProject

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			projectDocFilter(owner), projectItemFilterById(identity),
			nil, &v, // TODO project to projectItems
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toPorjectDO(owner, &v[0].Items[0], &do)

	return
}

func (col project) toPorjectDO(owner string, item *projectItem, do *repositories.ProjectDO) {
	*do = repositories.ProjectDO{
		Id:       item.Id,
		Owner:    owner,
		Name:     item.Name,
		Desc:     item.Desc,
		Type:     item.Type,
		CoverId:  item.CoverId,
		Protocol: item.Protocol,
		Training: item.Training,
		RepoType: item.RepoType,
		Tags:     item.Tags,
	}
}
