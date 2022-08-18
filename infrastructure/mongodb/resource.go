package mongodb

import (
	"context"

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

//TODO check and delete if ok
func (col project) ListUsersProjects1(opts map[string][]string) (
	r []repositories.ProjectDO, err error,
) {
	n := len(opts)
	users := make([]string, n)
	ids := make([]string, 0, n)
	n = 0
	for k, v := range opts {
		users[n] = k
		ids = append(ids, v...)
		n++
	}

	var v []dProject

	f := func(ctx context.Context) error {
		return cli.getArraysElemsByCustomizedCond(
			ctx, col.collectionName,
			bson.M{fieldOwner: bson.M{"$in": users}},
			map[string]func() bson.M{
				fieldItems: func() bson.M {
					return inCondForArrayElem(fieldId, ids)
				},
			},
			nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 {
		return
	}

	r = make([]repositories.ProjectDO, 0, len(v))
	for i := range v {
		owner := v[i].Owner
		items := v[i].Items

		for j := range items {
			col.toProjectDO(owner, &items[j], &r[i])
		}
	}

	return
}
