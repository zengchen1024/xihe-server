package mongodb

import (
	"context"
	"errors"

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

func (col dataset) New(owner string) error {
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

	if err := withContext(f); err != nil && !errors.Is(err, errDocExists) {
		return err
	}

	return nil
}

func (col dataset) Insert(do repositories.DatasetDO) (identity string, err error) {
	identity = newId()

	docObj := datasetItem{
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

	docFilter := datasetDocFilter(do.Owner)

	appendElemMatchToFilter(
		fieldItems, false,
		datasetItemFilter(do.Name), docFilter,
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
