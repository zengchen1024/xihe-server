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

func resourceIdFilter(identity string) bson.M {
	return bson.M{
		fieldId: identity,
	}
}

func subfieldOfItems(k string) string {
	return fieldItems + "." + k
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

func updateResourceLike(collection string, r *repositories.ResourceIndexDO, num int) error {
	updated := false
	f := func(ctx context.Context) error {
		b, err := cli.updateArrayElemCount(
			ctx, collection, fieldItems, fieldLikeCount, num,
			resourceOwnerFilter(r.Owner), resourceIdFilter(r.Id),
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

func getResourceById(collection, owner, rid string, result interface{}) error {
	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, collection, fieldItems,
			resourceOwnerFilter(owner), resourceIdFilter(rid),
			nil, result,
		)
	}

	return withContext(f)
}

func getResourceSummary(collection, owner, rId string, result interface{}) error {
	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, collection, fieldItems,
			resourceOwnerFilter(owner),
			resourceIdFilter(rId),
			bson.M{
				subfieldOfItems(fieldName):     1,
				subfieldOfItems(fieldRepoId):   1,
				subfieldOfItems(fieldRepoType): 1,
			},
			result,
		)
	}

	return withContext(f)
}

func getResourceSummaryByName(collection, owner, name string, result interface{}) error {
	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, collection, fieldItems,
			resourceOwnerFilter(owner),
			resourceNameFilter(name),
			bson.M{
				subfieldOfItems(fieldId):       1,
				subfieldOfItems(fieldRepoId):   1,
				subfieldOfItems(fieldRepoType): 1,
			},
			result,
		)
	}

	return withContext(f)
}

func getResourceByName(collection, owner, name string, result interface{}) error {
	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, collection, fieldItems,
			resourceOwnerFilter(owner), resourceNameFilter(name),
			nil, result,
		)
	}

	return withContext(f)
}

func sortByUpdateTime() bson.M {
	return bson.M{subfieldOfItems(fieldUpdatedAt): -1}
}

func sortByFirstLetter() bson.M {
	return bson.M{
		subfieldOfItems(fieldFirstLetter): 1,
	}
}

func sortByDownloadCount() bson.M {
	return bson.M{subfieldOfItems(fieldDownloadCount): -1}
}

func insertResource(collection, owner, name string, doc bson.M) error {
	docFilter := resourceOwnerFilter(owner)

	appendElemMatchToFilter(
		fieldItems, false, resourceNameFilter(name), docFilter,
	)

	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, collection, fieldItems, docFilter, doc,
		)
	}

	return withContext(f)
}

func updateResourceProperty(
	collection string, obj *repositories.ResourceToUpdateDO,
	property interface{},
) error {
	doc, err := genDoc(property)
	if err != nil {
		return err
	}

	updated := false

	f := func(ctx context.Context) error {
		updated, err = cli.updateArrayElem(
			ctx, collection, fieldItems,
			resourceOwnerFilter(obj.Owner),
			resourceIdFilter(obj.Id),
			doc, obj.Version, obj.UpdatedAt,
		)

		return err
	}

	if withContext(f); err != nil {
		return err
	}

	if !updated {
		return repositories.NewErrorConcurrentUpdating(
			errors.New("no update"),
		)
	}

	return nil
}

func listResourceWithoutSort(
	collection, owner string,
	do *repositories.ResourceListDO,
	fields []string, result interface{},
) error {
	fieldItemsRef := "$" + fieldItems

	project := bson.M{
		fieldItems: bson.M{"$filter": bson.M{
			"input": fieldItemsRef,
			"cond": func() bson.M {
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
			}(),
		}},
	}

	keep := bson.M{}
	for _, item := range fields {
		keep[subfieldOfItems(item)] = 1
	}

	pipeline := bson.A{
		bson.M{"$match": resourceOwnerFilter(owner)},
		bson.M{"$project": project},
		bson.M{"$project": keep},
	}

	return withContext(func(ctx context.Context) error {
		col := cli.collection(collection)
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, result)
	})
}

func globalListResourceWithoutSort(
	collection string,
	do *repositories.GlobalResourceListDO,
	fields []string, result interface{},
) error {
	fieldItemsRef := "$" + fieldItems

	project := bson.M{
		fieldItems: bson.M{"$filter": bson.M{
			"input": fieldItemsRef,
			"cond": func() bson.M {
				conds := bson.A{}

				if len(do.Tags) > 0 {
					for _, tag := range do.Tags {
						conds = append(conds, valueInCondForArrayElem(
							fieldTags, tag,
						))
					}
				}

				if do.Name != "" {
					conds = append(conds, matchCondForArrayElem(
						fieldName, do.Name,
					))
				}

				return condForArrayElem(conds)
			}(),
		}},
		fieldOwner: 1,
	}

	keep := bson.M{fieldOwner: 1}
	for _, item := range fields {
		keep[subfieldOfItems(item)] = 1
	}

	pipeline := bson.A{
		bson.M{"$project": project},
		bson.M{"$project": keep},
	}

	return withContext(func(ctx context.Context) error {
		col := cli.collection(collection)
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, result)
	})
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

func listUsersResourcesSummary(collection string, opts map[string][]string, result interface{}) error {
	n := len(opts)
	users := make([]string, n)
	names := make([]string, 0, n)
	n = 0
	for k, v := range opts {
		users[n] = k
		names = append(names, v...)
		n++
	}

	f := func(ctx context.Context) error {
		return cli.getArraysElemsByCustomizedCond(
			ctx, collection,
			bson.M{fieldOwner: bson.M{"$in": users}},
			map[string]func() bson.M{
				fieldItems: func() bson.M {
					return inCondForArrayElem(fieldName, names)
				},
			},
			bson.M{
				fieldOwner:                     1,
				subfieldOfItems(fieldId):       1,
				subfieldOfItems(fieldName):     1,
				subfieldOfItems(fieldRepoId):   1,
				subfieldOfItems(fieldRepoType): 1,
			}, result,
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
	arrayFilter := resourceIdFilter(do.Id)

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

func updateReverselyRelatedResource(
	collection, field string, add bool,
	do *repositories.ReverselyRelatedResourceInfoDO,
) error {
	doc := bson.M{
		field: bson.M{
			fieldRId:    do.Promoter.Id,
			fieldROwner: do.Promoter.Owner,
		},
	}

	docFilter := resourceOwnerFilter(do.Resource.Owner)
	arrayFilter := resourceIdFilter(do.Resource.Id)

	f := func(ctx context.Context) error {
		op := ""
		if add {
			op = mongoCmdPush
		} else {
			op = mongoCmdPull
		}

		_, err := cli.modifyArrayElemWithoutVersion(
			ctx, collection, fieldItems,
			docFilter, arrayFilter, doc, op,
		)

		return err
	}

	return withContext(f)
}
