package mongodb

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func projectDocFilter(owner string) bson.M {
	return bson.M{
		fieldOwner: owner,
	}
}

func projectItemFilter(name string) bson.M {
	return bson.M{
		fieldName: name,
	}
}

func NewProjectMapper(name string) repositories.ProjectMapper {
	return project{name}
}

type project struct {
	collectionName string
}

func (col project) New(owner string) error {
	docFilter := projectDocFilter(owner)
	doc := bson.M{
		fieldOwner: owner,
		fieldItems: bson.A{},
	}

	f := func(ctx context.Context) error {
		_, err := cli.newDocIfNotExist(
			ctx, col.collectionName, docFilter, doc,
		)

		return err
	}

	if err := withContext(f); err != nil && !errors.Is(err, errDocExists) {
		return err
	}

	return nil
}

func (col project) Insert(do repositories.ProjectDO) (string, error) {
	docObj := projectItem{
		Id:       newId(),
		Name:     do.Name,
		Desc:     do.Desc,
		Type:     do.Type,
		CoverId:  do.CoverId,
		Protocol: do.Protocol,
		Training: do.Training,
		RepoType: do.RepoType,
		Tags:     do.Tags,
	}

	doc, err := genDoc(docObj)
	if err != nil {
		return "", err
	}

	docFilter := projectDocFilter(do.Owner)

	appendElemMatchToFilter(
		fieldItems, false,
		projectItemFilter(do.Name), docFilter,
	)

	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, col.collectionName,
			fieldItems, docFilter, doc,
		)
	}

	return docObj.Id, withContext(f)
}

func (col project) Update(string, repositories.ProjectDO) error {
	return nil
}

func (col project) Get(string) (p repositories.ProjectDO, err error) {
	return
}
