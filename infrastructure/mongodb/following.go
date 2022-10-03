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
func (col user) AddFollowing(user, follower string) error {
	return col.addFollow(follower, user, fieldFollowing)
}

func (col user) RemoveFollowing(user, follower string) error {
	return col.removeFollow(follower, user, fieldFollowing)
}

func (col user) ListFollowing(do *repositories.FollowerUsersInfoListDO) (
	[]repositories.FollowerUserInfoDO, int, error,
) {
	//TODO if do.User == do.Follower , no need to check is follower

	return col.listFollow(do, fieldFollowing)
}

// follower
func (col user) AddFollower(user, follower string) error {
	return col.addFollow(user, follower, fieldFollower)
}

func (col user) RemoveFollower(user, follower string) error {
	return col.removeFollow(user, follower, fieldFollower)
}

func (col user) ListFollower(do *repositories.FollowerUsersInfoListDO) (
	[]repositories.FollowerUserInfoDO, int, error,
) {
	return col.listFollow(do, fieldFollower)
}

// helper
func (col user) addFollow(user, account, field string) error {
	f := func(ctx context.Context) error {
		return cli.addToSimpleArray(
			ctx, col.collectionName, field,
			userDocFilterByAccount(user), account,
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

func (col user) removeFollow(user, account, field string) error {
	f := func(ctx context.Context) error {
		return cli.removeFromSimpleArray(
			ctx, col.collectionName, field,
			userDocFilterByAccount(user), account,
		)
	}

	return withContext(f)
}

func (col user) listFollow(do *repositories.FollowerUsersInfoListDO, field string) (
	[]repositories.FollowerUserInfoDO, int, error,
) {
	var u DUser

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, userDocFilterByAccount(do.User),
			bson.M{field: 1}, &u,
		)
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return nil, 0, err
	}

	var v []string

	switch field {
	case fieldFollower:
		v = u.Follower

	case fieldFollowing:
		v = u.Following
	}

	r, err := col.listFollows(do.Follower, v)
	if err != nil {
		return nil, 0, err
	}

	return r, len(v), nil
}

func (col user) listFollows(follower string, accounts []string) (
	[]repositories.FollowerUserInfoDO, error,
) {
	var v []struct {
		DUser `bson:",inline"`

		IsFollower bool `bson:"is_follower"`
	}

	pipeline := bson.A{
		bson.M{"$match": bson.M{
			"$expr": bson.M{"$in": bson.A{"$" + fieldName, accounts}},
		}},
		bson.M{"$project": bson.M{
			fieldIsFollower: bson.M{
				"$in": bson.A{follower, "$" + fieldFollower},
			},
			fieldBio:      1,
			fieldName:     1,
			fieldAvatarId: 1,
		}},
	}

	err := withContext(func(ctx context.Context) error {
		col := cli.collection(col.collectionName)
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	})

	if err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]repositories.FollowerUserInfoDO, len(v))
	for i := range v {
		item := &v[i]

		r[i] = repositories.FollowerUserInfoDO{
			Bio:        item.Bio,
			Account:    item.Name,
			AvatarId:   item.AvatarId,
			IsFollower: item.IsFollower,
		}
	}

	return r, nil
}
