package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func userDocFilterByAccount(account string) bson.M {
	return bson.M{
		fieldName: account,
	}
}

func (col user) AddFollowing(owner, account string) error {
	f := func(ctx context.Context) error {
		return cli.addToSimpleArray(
			ctx, col.collectionName, fieldFollowing,
			userDocFilterByAccount(owner), account,
		)
	}

	if err := withContext(f); err != nil {
		if isDocExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}

		return err
	}

	return nil
}

func (col user) RemoveFollowing(owner, account string) error {
	f := func(ctx context.Context) error {
		return cli.removeFromSimpleArray(
			ctx, col.collectionName, fieldFollowing,
			userDocFilterByAccount(owner), account,
		)
	}

	return withContext(f)
}

func (col user) ListFollowing(owner string) ([]repositories.UserInfoDO, error) {
	var v dUser

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, userDocFilterByAccount(owner),
			bson.M{fieldFollowing: 1}, &v,
		)
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return nil, err
	}

	return col.listFollowing(v.Following)
}

func (col user) listFollowing(accounts []string) ([]repositories.UserInfoDO, error) {
	var v []dUser

	f := func(ctx context.Context) error {
		filter := bson.M{
			fieldName: bson.M{
				"$in": accounts,
			},
		}

		return cli.getDocs(
			ctx, col.collectionName, filter,
			bson.M{
				fieldBio:      1,
				fieldName:     1,
				fieldAvatarId: 1,
			}, &v,
		)
	}

	if err := withContext(f); err != nil {
		return nil, err
	}

	if len(v) == 0 {
		return nil, nil
	}

	r := make([]repositories.UserInfoDO, len(v))
	for i := range v {
		item := &v[i]

		r[i] = repositories.UserInfoDO{
			Bio:      item.Bio,
			Account:  item.Name,
			AvatarId: item.AvatarId,
		}
	}

	return r, nil

}
