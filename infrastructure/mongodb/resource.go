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

func toResourceObject(do *repositories.ResourceObjectDO) ResourceObject {
	return ResourceObject{
		Id:    do.Id,
		Type:  do.Type,
		Owner: do.Owner,
	}
}

func toResourceObjectDO(doc *ResourceObject) repositories.ResourceObjectDO {
	return repositories.ResourceObjectDO{
		Id:    doc.Id,
		Type:  doc.Type,
		Owner: doc.Owner,
	}
}

func toResourceIndexDO(v []ResourceIndex) []repositories.ResourceIndexDO {
	if len(v) == 0 {
		return nil
	}

	r := make([]repositories.ResourceIndexDO, len(v))
	for i := range v {
		a, b := &r[i], &v[i]

		a.Id = b.Id
		a.Owner = b.Owner
	}

	return r
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
			bson.M{fieldOwner: 1, fieldItems: 1}, result,
		)
	}

	return withContext(f)
}

func updateRelatedResource(
	collection, field string, add bool,
	do *repositories.RelatedResourceDO,
) error {
	doc := bson.M{
		field: bson.M{
			fieldRId:    do.ResourceId,
			fieldROwner: do.ResourceOwner,
		},
	}

	docFilter := resourceOwnerFilter(do.Owner)
	arrayFilter := arrayFilterById(do.Id)

	updated := false
	var err error
	f := func(ctx context.Context) error {
		if add {
			updated, err = cli.pushNestedArrayElem(
				ctx, collection, fieldItems,
				docFilter, arrayFilter, doc,
				do.Version, do.UpdatedAt,
			)
		} else {
			updated, err = cli.pullNestedArrayElem(
				ctx, collection, fieldItems,
				docFilter, arrayFilter, doc,
				do.Version, do.UpdatedAt,
			)
		}

		return nil
	}

	if withContext(f); err != nil {
		return err
	}

	if !updated {
		return repositories.NewErrorConcurrentUpdating(errors.New("no update"))
	}

	return nil
}
