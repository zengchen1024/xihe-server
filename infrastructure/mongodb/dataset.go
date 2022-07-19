package mongodb

import (
	"context"
	"errors"
	"fmt"

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

func (col dataset) List(owner string, do repositories.DatasetListDO) (
	r []repositories.DatasetDO, err error,
) {
	var v []dDataset

	f := func(ctx context.Context) error {
		return cli.getArraysElemsByCustomizedCond(
			ctx, col.collectionName, datasetDocFilter(owner),
			map[string]func() bson.M{
				fieldItems: func() bson.M {
					if do.Name == "" {
						return bson.M{
							"$toBool": 1,
						}
					}

					return bson.M{
						"$regexMatch": bson.M{
							"input": fmt.Sprintf("$$this.%s", fieldName),
							"regex": do.Name,
						},
					}
				},
			},
			bson.M{fieldItems: 1}, &v,
		)
	}

	if err = withContext(f); err != nil {
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

func (col dataset) toDatasetDO(owner string, item *datasetItem, do *repositories.DatasetDO) {
	*do = repositories.DatasetDO{
		Id:       item.Id,
		Owner:    owner,
		Name:     item.Name,
		Desc:     item.Desc,
		Protocol: item.Protocol,
		RepoType: item.RepoType,
		Tags:     item.Tags,
	}
}
