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

func (col model) New(owner string) error {
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

	if err := withContext(f); err != nil && !errors.Is(err, errDocExists) {
		return err
	}

	return nil
}

func (col model) Insert(do repositories.ModelDO) (identity string, err error) {
	identity = newId()

	docObj := modelItem{
		Id:       identity,
		Name:     do.Name,
		Desc:     do.Desc,
		Protocol: do.Protocol,
		RepoType: do.RepoType,
		Tags:     do.Tags,
	}

	doc, err := genDoc(docObj)
	if err != nil {
		return
	}

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

	if errors.Is(err, errDocNotExists) {
		err = repositories.NewErrorDuplicateCreating(err)
	}

	return
}
