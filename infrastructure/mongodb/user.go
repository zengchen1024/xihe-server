package mongodb

import (
	"context"
	"errors"
	"sort"

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
			filter, doc, mongoCmdSet, do.Version,
		)
	}

	if err = withContext(f); err != nil && isDocNotExists(err) {
		err = repositories.NewErrorConcurrentUpdating(err)
	}

	return
}

func (col user) GetByFollower(info repositories.FollowerInfoDO) (
	do repositories.UserDO, isFollower bool, err error,
) {
	var v []struct {
		DUser `bson:",inline"`

		IsFollower     bool `bson:"is_follower"`
		FollowerCount  int  `bson:"follower_count"`
		FollowingCount int  `bson:"following_count"`
	}

	f := func(ctx context.Context) error {
		fields := bson.M{
			fieldFollowerCount:  bson.M{"$size": "$" + fieldFollower},
			fieldFollowingCount: bson.M{"$size": "$" + fieldFollowing},
		}

		if info.Follower != "" {
			fields[fieldIsFollower] = bson.M{
				"$in": bson.A{info.Follower, "$" + fieldFollower},
			}
		}

		pipeline := bson.A{
			bson.M{"$match": userDocFilterByAccount(info.User)},
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
	col.toUserDO(&item.DUser, &do)

	do.FollowerCount = item.FollowerCount
	do.FollowingCount = item.FollowingCount

	if info.Follower != "" {
		isFollower = item.IsFollower
	}

	return
}

func (col user) GetUserAvatarId(account string) (string, error) {
	var v DUser

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, bson.M{fieldName: account},
			bson.M{fieldAvatarId: 1}, &v,
		)
	}

	if err := withContext(f); err != nil {
		return "", err
	}

	return v.AvatarId, nil
}

func (col user) ListUsersInfo(accounts []string) ([]repositories.UserInfoDO, error) {
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

type userInfo struct {
	DUser `bson:",inline"`

	Count int `bson:"count"`
}

func (col user) Search(do *repositories.UserSearchOptionDO) (
	r []string, total int, err error,
) {
	key := "$" + fieldFollower
	fieldCount := "count"
	fields := bson.M{fieldCount: bson.M{
		"$cond": bson.M{
			"if":   bson.M{"$isArray": key},
			"then": bson.M{"$size": key},
			"else": 0,
		},
	}}

	pipeline := bson.A{
		bson.M{mongoCmdMatch: bson.M{
			"$expr": bson.M{
				"$regexMatch": bson.M{
					"input":   "$" + fieldName,
					"regex":   do.Name,
					"options": "i",
				},
			},
		}},
		bson.M{"$addFields": fields},
		bson.M{"$project": bson.M{
			fieldName:  1,
			fieldCount: 1,
		}},
	}

	var v []userInfo

	f := func(ctx context.Context) error {
		col := cli.collection(col.collectionName)
		cursor, err := col.Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, &v)
	}

	if err = withContext(f); err != nil || len(v) == 0 {
		return
	}

	total = len(v)

	items := make([]*userInfo, total)
	for i := range v {
		items[i] = &v[i]
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].Count >= items[j].Count
	})

	n := do.TopNum
	if total < n {
		n = total
	}

	r = make([]string, n)
	for i := range r {
		r[i] = items[i].Name
	}

	return
}

func (col user) toUserDoc(do *repositories.UserDO) (bson.M, error) {
	docObj := DUser{
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

func (col user) toUserDO(u *DUser, do *repositories.UserDO) {
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
