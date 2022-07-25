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

	f := func(ctx context.Context) error {
		v, err := cli.newDocIfNotExist(
			ctx, col.collectionName,
			userDocFilter(do.Account, do.Email), doc,
		)

		identity = v

		return err
	}

	if err = withContext(f); err != nil && errors.Is(err, errDocExists) {
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

	if err = withContext(f); err != nil && errors.Is(err, errDocNotExists) {
		err = repositories.NewErrorConcurrentUpdating(err)
	}

	return
}

func (col user) Get(uid string) (do repositories.UserDO, err error) {
	v, err := objectIdFilter(uid)
	if err != nil {
		return
	}

	return col.get(v)
}

func (col user) GetByAccount(account string) (do repositories.UserDO, err error) {
	return col.get(bson.M{fieldName: account})
}

func (col user) get(filter bson.M) (do repositories.UserDO, err error) {
	var v dUser

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, filter, nil, &v,
		)
	}

	err = withContext(f)

	if err == nil {
		col.toUserDO(&v, &do)

		return
	}

	if errors.Is(err, errDocNotExists) {
		err = repositories.NewErrorDataNotExists(err)
	}

	return
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
