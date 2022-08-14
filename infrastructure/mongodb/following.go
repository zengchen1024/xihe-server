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

func (col user) AddFollowing(owner, account string) (do repositories.UserDO, err error) {
	f := func(ctx context.Context) error {
		return cli.addToSimpleArray(
			ctx, col.collectionName, fieldFollowing,
			userDocFilterByAccount(owner), account,
		)
	}

	if err = withContext(f); err != nil {
		if isDocExists(err) {
			err = repositories.NewErrorConcurrentUpdating(err)
		}

		return
	}

	return col.GetByAccount(owner)
}

func (col user) RemoveFollowing(owner, account string) (do repositories.UserDO, err error) {
	f := func(ctx context.Context) error {
		return cli.removeFromSimpleArray(
			ctx, col.collectionName, fieldFollowing,
			userDocFilterByAccount(owner), account,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	return col.GetByAccount(owner)
}
