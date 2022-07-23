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
	docObj := dUser{
		Name:                    do.Account,
		Email:                   do.Email,
		Bio:                     do.Bio,
		AvatarId:                do.AvatarId,
		PlatformToken:           do.Platform.Token,
		PlatformUserId:          do.Platform.UserId,
		PlatformUserNamespaceId: do.Platform.NamespaceId,
	}

	doc, err := genDoc(docObj)
	if err != nil {
		return
	}

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

func (col user) Get(identity string) (do repositories.UserDO, err error) {
	return
}
