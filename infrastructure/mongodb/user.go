package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func userDocFilter(account, email string) bson.M {
	return bson.M{
		"$or": bson.A{
			bson.M{
				fieldName: account,
			},
			bson.M{
				fieldEmail: email,
			},
		},
	}
}

func NewUserMapper(name string) repositories.UserMapper {
	return user{name}
}

type user struct {
	collectionName string
}

func (col user) Insert(do repositories.UserDO) (identity string, err error) {
	doc, err := col.toUserDoc(&do)
	if err != nil {
		return
	}
	doc[fieldVersion] = 0
	doc[fieldFollower] = bson.A{}
	doc[fieldFollowing] = bson.A{}

	f := func(ctx context.Context) error {
		v, err := cli.newDocIfNotExist(
			ctx, col.collectionName,
			userDocFilter(do.Account, do.Email), doc,
		)

		identity = v

		return err
	}

	if err = withContext(f); err != nil && isDocExists(err) {
		err = repositories.NewErrorDuplicateCreating(err)
	}

	return
}

func (col user) Update(do repositories.UserDO) (err error) {
	doc, err := col.toUserDoc(&do)
	if err != nil {
		return
	}

	filter, err := objectIdFilter(do.Id)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, col.collectionName,
			filter, doc, do.Version,
		)
	}

	if err = withContext(f); err != nil && isDocNotExists(err) {
		err = repositories.NewErrorConcurrentUpdating(err)
	}

	return
}

func (col user) GetByAccount(account string) (do repositories.UserDO, err error) {
	do, _, err = col.GetByFollower(account, "")

	return
}

func (col user) GetByFollower(account, follower string) (
	do repositories.UserDO, isFollower bool, err error,
) {
	var v []dUser

	f := func(ctx context.Context) error {
		fields := bson.M{
			fieldFollowerCount:  bson.M{"$size": "$" + fieldFollower},
			fieldFollowingCount: bson.M{"$size": "$" + fieldFollowing},
		}

		if follower != "" {
			fields[fieldIsFollower] = bson.M{
				"$in": bson.A{follower, "$" + fieldFollower},
			}
		}

		pipeline := bson.A{
			bson.M{"$match": userDocFilterByAccount(account)},
			bson.M{"$addFields": fields},
			bson.M{"$project": bson.M{
				fieldFollowing: 0,
				fieldFollower:  0,
			}},
		}

		cursor, err := cli.collection(col.collectionName).Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 {
		err = repositories.NewErrorDataNotExists(errors.New("no user"))

		return
	}

	item := &v[0]
	col.toUserDO(item, &do)

	do.FollowerCount = item.FollowerCount
	do.FollowingCount = item.FollowingCount

	if follower != "" {
		isFollower = item.IsFollower
	}

	return
}

func (col user) ListUsers(accounts []string) ([]repositories.UserInfoDO, error) {
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
			Account:  item.Name,
			AvatarId: item.AvatarId,
		}
	}

	return r, nil
}

func (col user) toUserDoc(do *repositories.UserDO) (bson.M, error) {
	docObj := dUser{
		Name:                    do.Account,
		Email:                   do.Email,
		Bio:                     do.Bio,
		AvatarId:                do.AvatarId,
		PlatformToken:           do.Platform.Token,
		PlatformUserId:          do.Platform.UserId,
		PlatformUserNamespaceId: do.Platform.NamespaceId,
	}

	return genDoc(docObj)
}

func (col user) toUserDO(u *dUser, do *repositories.UserDO) {
	*do = repositories.UserDO{
		Id:       u.Id.Hex(),
		Email:    u.Email,
		Account:  u.Name,
		Bio:      u.Bio,
		AvatarId: u.AvatarId,
		Version:  u.Version,
	}

	do.Platform.Token = u.PlatformToken
	do.Platform.UserId = u.PlatformUserId
	do.Platform.NamespaceId = u.PlatformUserNamespaceId
}
