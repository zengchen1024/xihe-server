package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func activityDocFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
	}
}

func NewActivityMapper(name string, keep int) repositories.ActivityMapper {
	return activity{
		collectionName: name,
		keepNum:        -keep,
	}
}

type activity struct {
	collectionName string
	keepNum        int
}

func (col activity) Insert(owner string, do repositories.ActivityDO) (err error) {
	if err = col.insert(owner, do); err == nil || isDBError(err) {
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

func (col activity) insert(owner string, do repositories.ActivityDO) error {
	v := col.toActivityDoc(&do)
	doc, err := genDoc(v)
	if err != nil {
		return err
	}

	docFilter := activityDocFilter(owner)
	resource, _ := genDoc(v.ResourceObj)
	appendElemMatchToFilter(fieldItems, false, resource, docFilter)

	f := func(ctx context.Context) error {
		r, err := cli.collection(col.collectionName).UpdateOne(
			ctx, docFilter,
			bson.M{"$push": bson.M{fieldItems: bson.M{
				"$each":  bson.A{doc},
				"$slice": col.keepNum,
			}}},
		)
		if err != nil {
			return dbError{err}
		}

		if r.MatchedCount == 0 {
			return errDocNotExists
		}

		return nil

	}

	return withContext(f)
}

func (col activity) List(owner string, opt repositories.ActivityListDO) (
	r []repositories.ActivityDO, err error,
) {
	var v dActivity

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, activityDocFilter(owner), nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return
	}

	items := v.Items
	r = make([]repositories.ActivityDO, len(items))
	for i := range items {
		col.toActivityDO(&items[i], &r[i])
	}

	return
}

func (col activity) toActivityDoc(do *repositories.ActivityDO) activityItem {
	return activityItem{
		Type: do.Type,
		ResourceObj: ResourceObj{
			ResourceId:    do.ResourceId,
			ResourceType:  do.ResourceType,
			ResourceOwner: do.ResourceOwner,
		},
	}
}

func (col activity) toActivityDO(item *activityItem, do *repositories.ActivityDO) {
	*do = repositories.ActivityDO{
		Type:          item.Type,
		ResourceId:    item.ResourceId,
		ResourceType:  item.ResourceType,
		ResourceOwner: item.ResourceOwner,
	}
}
