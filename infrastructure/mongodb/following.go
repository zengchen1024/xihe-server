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
	r []repositories.FollowerUserInfoDO, total int, err error,
) {
	v, err := col.getFollows(do.User, fieldFollowing)
	if err != nil {
		return
	}

	items := v.Following
	total = len(items)
	if total == 0 {
		return
	}

	items = col.getPageItems(items, do)
	if len(items) == 0 {
		return
	}

	if b := do.User == do.Follower; b || do.Follower == "" {
		r, err = col.listFollowsDirectly(items, b)
	} else {
		r, err = col.listFollows(do.Follower, items)
	}

	return r, total, err
}

// follower
func (col user) AddFollower(user, follower string) error {
	return col.addFollow(user, follower, fieldFollower)
}

func (col user) RemoveFollower(user, follower string) error {
	return col.removeFollow(user, follower, fieldFollower)
}

func (col user) ListFollower(do *repositories.FollowerUsersInfoListDO) (
	r []repositories.FollowerUserInfoDO, total int, err error,
) {
	v, err := col.getFollows(do.User, fieldFollower)
	if err != nil {
		return
	}

	items := v.Follower
	total = len(items)
	if total == 0 {
		return
	}

	items = col.getPageItems(items, do)
	if len(items) == 0 {
		return
	}

	if do.Follower == "" {
		r, err = col.listFollowsDirectly(items, false)
	} else {
		r, err = col.listFollows(do.Follower, items)
	}

	return r, total, err
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

func (col user) getFollows(user, field string) (v DUser, err error) {
	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, userDocFilterByAccount(user),
			bson.M{field: 1}, &v,
		)
	}

	if err = withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}
	}

	return
}

func (col user) getPageItems(items []string, do *repositories.FollowerUsersInfoListDO) []string {
	if do.CountPerPage <= 0 {
		return items
	}

	total := len(items)

	if do.PageNum <= 1 {
		if total > do.CountPerPage {
			return items[:do.CountPerPage]
		}

		return items
	}

	skip := do.CountPerPage * (do.PageNum - 1)
	if skip >= total {
		return nil
	}

	if n := total - skip; n > do.CountPerPage {
		return items[skip : skip+do.CountPerPage]
	}

	return items[skip:]
}

func (col user) listFollowsDirectly(accounts []string, isFollower bool) (
	[]repositories.FollowerUserInfoDO, error,
) {
	var v []DUser

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

	if err := withContext(f); err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]repositories.FollowerUserInfoDO, len(v))
	for i := range v {
		item := &v[i]

		r[i] = repositories.FollowerUserInfoDO{
			Bio:        item.Bio,
			Account:    item.Name,
			AvatarId:   item.AvatarId,
			IsFollower: isFollower,
		}
	}

	return r, nil
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
