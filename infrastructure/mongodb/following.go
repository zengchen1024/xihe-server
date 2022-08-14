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
	// see cla

	return nil, nil
}
