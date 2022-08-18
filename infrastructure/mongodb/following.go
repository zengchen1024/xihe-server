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

// following
func (col user) AddFollowing(owner, account string) error {
	return col.addFollow(owner, account, fieldFollowing)
}

func (col user) RemoveFollowing(owner, account string) error {
	return col.removeFollow(owner, account, fieldFollowing)
}

func (col user) ListFollowing(owner string) ([]repositories.FollowUserInfoDO, error) {
	return col.listFollow(owner, fieldFollowing)
}

// follower
func (col user) AddFollower(owner, account string) error {
	return col.addFollow(owner, account, fieldFollower)
}

func (col user) RemoveFollower(owner, account string) error {
	return col.removeFollow(owner, account, fieldFollower)
}

func (col user) ListFollower(owner string) ([]repositories.FollowUserInfoDO, error) {
	return col.listFollow(owner, fieldFollower)
}

// helper
func (col user) addFollow(owner, account, field string) error {
	f := func(ctx context.Context) error {
		return cli.addToSimpleArray(
			ctx, col.collectionName, field,
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

func (col user) removeFollow(owner, account, field string) error {
	f := func(ctx context.Context) error {
		return cli.removeFromSimpleArray(
			ctx, col.collectionName, field,
			userDocFilterByAccount(owner), account,
		)
	}

	return withContext(f)
}

func (col user) listFollow(owner, field string) ([]repositories.FollowUserInfoDO, error) {
	var u dUser

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, userDocFilterByAccount(owner),
			bson.M{field: 1}, &u,
		)
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return nil, err
	}

	var v []string
	switch field {
	case fieldFollower:
		v = u.Follower
	case fieldFollowing:
		v = u.Following
	}

	return col.listFollows(v)
}

func (col user) listFollows(accounts []string) ([]repositories.FollowUserInfoDO, error) {
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

	r := make([]repositories.FollowUserInfoDO, len(v))
	for i := range v {
		item := &v[i]

		r[i] = repositories.FollowUserInfoDO{
			Bio:      item.Bio,
			Account:  item.Name,
			AvatarId: item.AvatarId,
		}
	}

	return r, nil
}
