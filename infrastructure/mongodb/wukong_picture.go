package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewWuKongPictureMapper(name string) repositories.WuKongPictureMapper {
	return wukongPicture{
		collectionName: name,
	}
}

type wukongPicture struct {
	collectionName string
}

func (col wukongPicture) List(user string) ([]repositories.WuKongPictureDO, int, error) {
	var v dWuKongPicture

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName,
			resourceOwnerFilter(user),
			nil, &v,
		)
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			err = nil
		}

		return nil, 0, err
	}

	t := v.Items
	r := make([]repositories.WuKongPictureDO, len(t))

	for i := range t {
		col.toPictureDO(&t[i], &r[i])
	}

	return r, v.Version, nil
}

func (col wukongPicture) Insert(user string, do *repositories.WuKongPictureDO, version int) (
	identity string, err error,
) {
	identity, err = col.insert(user, do, version)
	if err == nil || !isDocNotExists(err) {
		return
	}

	// doc is not exist or duplicate insert
	if err = col.newDoc(user); err != nil {
		return
	}

	identity, err = col.insert(user, do, version)
	if err != nil && isDocNotExists(err) {
		err = repositories.NewErrorConcurrentUpdating(err)
	}

	return
}

func (col wukongPicture) newDoc(user string) error {
	docFilter := resourceOwnerFilter(user)

	doc := bson.M{
		fieldOwner:   user,
		fieldItems:   bson.A{},
		fieldVersion: 0,
	}

	f := func(ctx context.Context) error {
		_, err := cli.newDocIfNotExist(
			ctx, col.collectionName, docFilter, doc,
		)

		return err
	}

	if err := withContext(f); err != nil && !isDocExists(err) {
		return err
	}

	return nil
}

func (col wukongPicture) insert(user string, do *repositories.WuKongPictureDO, version int) (
	identity string, err error,
) {
	identity = newId()
	do.Id = identity

	doc, err := col.toPictureDoc(do)
	if err != nil {
		return
	}

	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, col.collectionName,
			resourceOwnerFilter(user),
			bson.M{fieldItems: doc}, mongoCmdPush, version,
		)
	}

	err = withContext(f)

	return
}

func (col wukongPicture) Delete(user string, pid string) error {
	f := func(ctx context.Context) error {
		return cli.pullArrayElem(
			ctx, col.collectionName, fieldItems,
			resourceOwnerFilter(user),
			resourceIdFilter(pid),
		)
	}

	return withContext(f)
}

func (col wukongPicture) Get(user string, pid string) (do repositories.WuKongPictureDO, err error) {
	var v []dWuKongPicture

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			resourceOwnerFilter(user),
			resourceIdFilter(pid),
			nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 || len(v[0].Items) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	col.toPictureDO(&v[0].Items[0], &do)

	return
}

func (col wukongPicture) toPictureDO(p *pictureItem, do *repositories.WuKongPictureDO) {
	*do = repositories.WuKongPictureDO{
		Id:        p.Id,
		OBSPath:   p.OBSPath,
		CreatedAt: p.CreatedAt,
	}

	do.Desc = p.Desc
	do.Style = p.Style
}

func (col wukongPicture) toPictureDoc(do *repositories.WuKongPictureDO) (bson.M, error) {
	return genDoc(pictureItem{
		Id:        do.Id,
		Desc:      do.Desc,
		Style:     do.Style,
		OBSPath:   do.OBSPath,
		CreatedAt: do.CreatedAt,
	})
}
