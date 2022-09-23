package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func likeDocFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
	}
}

func NewLikeMapper(name string) repositories.LikeMapper {
	return like{name}
}

type like struct {
	collectionName string
}

func (col like) Insert(owner string, do repositories.LikeDO) (err error) {
	if err = col.insert(owner, do); err == nil || !isDocNotExists(err) {
		return
	}

	// doc is not exist or duplicate insert

	if err = newResourceDoc(col.collectionName, owner); err == nil {
		if err = col.insert(owner, do); err != nil && isDocNotExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}
	}

	return
}

func (col like) insert(owner string, do repositories.LikeDO) error {
	doc, err := col.toLikeDoc(&do)
	if err != nil {
		return err
	}

	obj, _ := genDoc(toResourceObj(&do.ResourceObjDO))
	docFilter := likeDocFilter(owner)
	appendElemMatchToFilter(fieldItems, false, obj, docFilter)

	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, col.collectionName,
			fieldItems, docFilter, doc,
		)
	}

	return withContext(f)
}

func (col like) Delete(owner string, do repositories.LikeDO) error {
	doc, err := genDoc(toResourceObj(&do.ResourceObjDO))
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		return cli.removeFromSimpleArray(
			ctx, col.collectionName, fieldItems,
			likeDocFilter(owner), doc,
		)
	}

	return withContext(f)
}

func (col like) List(owner string, opt repositories.LikeListDO) (
	r []repositories.LikeDO, err error,
) {
	var v dLike

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, likeDocFilter(owner), nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return
	}

	items := v.Items
	r = make([]repositories.LikeDO, len(items))
	for i := range items {
		col.toLikeDO(&items[i], &r[i])
	}

	return
}

func (col like) HasLike(owner string, do *repositories.ResourceObjDO) (b bool, err error) {
	doc, err := genDoc(toResourceObj(do))
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		b, err = cli.isArrayElemExists(
			ctx, col.collectionName, fieldItems,
			likeDocFilter(owner), doc,
		)

		return nil
	}

	withContext(f)

	return
}

func (col like) toLikeDoc(do *repositories.LikeDO) (bson.M, error) {
	docObj := likeItem{
		CreatedAt:   do.CreatedAt,
		ResourceObj: toResourceObj(&do.ResourceObjDO),
	}

	return genDoc(docObj)
}

func (col like) toLikeDO(item *likeItem, do *repositories.LikeDO) {
	*do = repositories.LikeDO{
		CreatedAt:     item.CreatedAt,
		ResourceObjDO: toResourceObjDO(&item.ResourceObj),
	}
}
