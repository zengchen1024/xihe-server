package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func userDocFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
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
		Bio:         do.Bio,
		Email:       do.Email,
		Account:     do.Account,
		Password:    do.Password,
		Nickname:    do.Nickname,
		AvatarId:    do.AvatarId,
		PhoneNumber: do.PhoneNumber,
	}

	doc, err := genDoc(docObj)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		v, err := cli.newDocIfNotExist(
			ctx, col.collectionName, userDocFilter(do.Account), doc,
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
