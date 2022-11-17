package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewLuoJiaMapper(name string) repositories.LuoJiaMapper {
	return luojia{name}
}

type luojia struct {
	collectionName string
}

func (col luojia) newDoc(owner string) error {
	docFilter := resourceOwnerFilter(owner)

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

func (col luojia) Insert(do *repositories.UserLuoJiaRecordDO) (identity string, err error) {
	if identity, err = col.insert(do); err == nil || !isDocNotExists(err) {
		return
	}

	// doc is not exist or duplicate insert

	if err = col.newDoc(do.User); err == nil {
		if identity, err = col.insert(do); err != nil && isDocNotExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}
	}

	return
}

func (col luojia) insert(do *repositories.UserLuoJiaRecordDO) (identity string, err error) {
	identity = newId()
	do.Id = identity

	item := new(luojiaItem)
	col.toLuoJiaRecordDoc(&do.LuoJiaRecordDO, item)

	doc, err := genDoc(item)
	if err != nil {
		return
	}

	docFilter := resourceOwnerFilter(do.User)

	f := func(ctx context.Context) error {
		return cli.pushElemToLimitedArray(
			ctx, col.collectionName, fieldItems, 10, docFilter, doc,
		)
	}

	err = withContext(f)

	return
}

func (col luojia) List(user string) (
	dos []repositories.LuoJiaRecordDO, err error,
) {
	var v dLuoJia

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName,
			resourceOwnerFilter(user), nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		if isDocNotExists(err) {
			err = nil
		}

		return
	}

	if len(v.Items) == 0 {
		return
	}

	dos = make([]repositories.LuoJiaRecordDO, len(v.Items))

	for i := range v.Items {
		col.toLuoJiaRecordDo(&dos[i], &v.Items[i])
	}

	return
}

func (col luojia) toLuoJiaRecordDoc(
	do *repositories.LuoJiaRecordDO, doc *luojiaItem,
) {
	doc.CreatedAt = do.CreatedAt
	doc.Id = do.Id
}

func (col luojia) toLuoJiaRecordDo(
	do *repositories.LuoJiaRecordDO, doc *luojiaItem,
) {
	do.CreatedAt = doc.CreatedAt
	do.Id = doc.Id
}
