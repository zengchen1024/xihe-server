package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func resourceOwnerFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
	}
}

func resourceNameFilter(name string) bson.M {
	return bson.M{
		fieldName: name,
	}
}

func toResourceObj(do *repositories.ResourceObjDO) ResourceObj {
	return ResourceObj{
		ResourceId:    do.ResourceId,
		ResourceType:  do.ResourceType,
		ResourceOwner: do.ResourceOwner,
	}
}

func toResourceObjDO(doc *ResourceObj) repositories.ResourceObjDO {
	return repositories.ResourceObjDO{
		ResourceId:    doc.ResourceId,
		ResourceType:  doc.ResourceType,
		ResourceOwner: doc.ResourceOwner,
	}
}

func newResourceDoc(collection, owner string) error {
	docFilter := resourceOwnerFilter(owner)

	doc := bson.M{
		fieldOwner: owner,
		fieldItems: bson.A{},
	}

	f := func(ctx context.Context) error {
		_, err := cli.newDocIfNotExist(
			ctx, collection, docFilter, doc,
		)

		return err
	}

	if err := withContext(f); err != nil && isDBError(err) {
		return err
	}

	return nil
}

func updateResourceLike(collection, owner, rid string, num int) error {
	updated := false
	f := func(ctx context.Context) error {
		b, err := cli.updateArrayElemCount(
			ctx, collection, fieldItems, fieldLikeCount, num,
			resourceOwnerFilter(owner), arrayFilterById(rid),
		)

		updated = b

		return err
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return err
	}

	if !updated {
		return repositories.NewErrorDataNotExists(errors.New("no update"))
	}

	return nil
}

func getResourceByName(collection, owner, name string, result interface{}) error {
	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, collection, fieldItems,
			resourceOwnerFilter(owner), resourceNameFilter(name),
			bson.M{fieldItems: 1}, result,
		)
	}

	return withContext(f)
}

func listResource(
	collection, owner string,
	do repositories.ResourceListDO, result interface{},
) error {
	f := func(ctx context.Context) error {
		return cli.getArraysElemsByCustomizedCond(
			ctx, collection, resourceOwnerFilter(owner),
			map[string]func() bson.M{
				fieldItems: func() bson.M {
					conds := bson.A{}

					if do.RepoType != "" {
						conds = append(conds, eqCondForArrayElem(
							fieldRepoType, do.RepoType,
						))
					}

					if do.Name != "" {
						conds = append(conds, matchCondForArrayElem(
							fieldName, do.Name,
						))
					}

					return condForArrayElem(conds)
				},
			},
			bson.M{fieldItems: 1}, result,
		)
	}

	return withContext(f)
}

func listUsersResources(collection string, opts map[string][]string, result interface{}) error {
	n := len(opts)
	users := make([]string, n)
	ids := make([]string, 0, n)
	n = 0
	for k, v := range opts {
		users[n] = k
		ids = append(ids, v...)
		n++
	}

	f := func(ctx context.Context) error {
		return cli.getArraysElemsByCustomizedCond(
			ctx, collection,
			bson.M{fieldOwner: bson.M{"$in": users}},
			map[string]func() bson.M{
				fieldItems: func() bson.M {
					return inCondForArrayElem(fieldId, ids)
				},
			},
			nil, result,
		)
	}

	return withContext(f)
}
