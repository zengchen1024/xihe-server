package mongodb

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewEvaluateMapper(name string) repositories.EvaluateMapper {
	return evaluate{name}
}

func evaluateDocFilter(owner, projectId, trainingId string) bson.M {
	return bson.M{
		fieldPId:   projectId,
		fieldTId:   trainingId,
		fieldOwner: owner,
	}
}

type evaluate struct {
	collectionName string
}

func (col evaluate) newDoc(do *repositories.EvaluateDO) error {
	docFilter := evaluateDocFilter(do.ProjectOwner, do.ProjectId, do.TrainingId)

	doc := bson.M{
		fieldOwner:   do.ProjectOwner,
		fieldPId:     do.ProjectId,
		fieldTId:     do.TrainingId,
		fieldItems:   bson.A{},
		fieldVersion: 0,
	}

	f := func(ctx context.Context) error {
		_, err := cli.newDocIfNotExist(
			ctx, col.collectionName, docFilter, doc,
		)

		return err
	}

	if err := withContext(f); err != nil && isDBError(err) {
		return err
	}

	return nil
}

func (col evaluate) Insert(do *repositories.EvaluateDO, version int) (
	identity string, err error,
) {
	identity, err = col.insert(do, version)
	if err == nil || !isDocNotExists(err) {
		return
	}

	// doc is not exist or duplicate insert

	if err = col.newDoc(do); err == nil {
		identity, err = col.insert(do, version)
		if err != nil && isDocNotExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}
	}

	return
}

func (col evaluate) insert(do *repositories.EvaluateDO, version int) (identity string, err error) {
	identity = newId()

	v := evaluateItem{
		Id:                identity,
		Type:              do.Type,
		MomentumScope:     do.Params.MomentumScope,
		BatchSizeScope:    do.Params.BatchSizeScope,
		LearningRateScope: do.Params.LearningRateScope,
	}

	doc, err := genDoc(v)
	if err != nil {
		return "", err
	}

	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, col.collectionName,
			evaluateDocFilter(do.ProjectOwner, do.ProjectId, do.TrainingId),
			bson.M{fieldItems: doc}, mongoCmdPush, version,
		)
	}

	err = withContext(f)

	return
}

func (col evaluate) UpdateDetail(
	index *repositories.EvaluateIndexDO,
	detail *repositories.EvaluateDetailDO,
) error {
	data := evaluateItem{
		Expiry:    detail.Expiry,
		Error:     detail.Error,
		AccessURL: detail.AccessURL,
	}

	doc, err := genDoc(data)
	if err != nil {
		return err
	}

	logrus.Debugf(
		"update evaluate(%s/%s) to %v",
		index.Project.Id, index.Id, doc,
	)

	f := func(ctx context.Context) error {
		_, err := cli.modifyArrayElemWithoutVersion(
			ctx, col.collectionName, fieldItems,
			col.docFilter(index),
			resourceIdFilter(index.Id), doc, mongoCmdSet,
		)

		return err
	}

	return withContext(f)
}

func (col evaluate) docFilter(index *repositories.EvaluateIndexDO) bson.M {
	return evaluateDocFilter(
		index.Project.Owner, index.Project.Id, index.TrainingId,
	)
}

func (col evaluate) GetStandardEvaluateParms(index *repositories.EvaluateIndexDO) (
	do repositories.StandardEvaluateParmsDO, err error,
) {
	var v []dEvaluate

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			col.docFilter(index),
			resourceIdFilter(index.Id),
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

	item := &v[0].Items[0]
	do.MomentumScope = item.MomentumScope
	do.BatchSizeScope = item.BatchSizeScope
	do.LearningRateScope = item.LearningRateScope

	return
}

func (col evaluate) Get(index *repositories.EvaluateIndexDO) (
	r repositories.EvaluateSummaryDO, err error,
) {
	var v []dEvaluate

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldItems,
			col.docFilter(index),
			resourceIdFilter(index.Id),
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

	col.toEvaluateSummaryDO(&v[0].Items[0], &r)

	return
}

func (col evaluate) List(index *repositories.ResourceIndexDO, trainingId string) (
	[]repositories.EvaluateSummaryDO, int, error,
) {
	var v dEvaluate

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName,
			evaluateDocFilter(index.Owner, index.Id, trainingId),
			bson.M{
				fieldVersion: 1,
				fieldItems:   1,
			}, &v)
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			return nil, 0, nil
		}

		return nil, 0, err
	}

	t := v.Items
	r := make([]repositories.EvaluateSummaryDO, len(t))

	for i := range t {
		col.toEvaluateSummaryDO(&t[i], &r[i])
	}

	return r, v.Version, nil
}

func (col evaluate) toEvaluateSummaryDO(doc *evaluateItem, r *repositories.EvaluateSummaryDO) {
	r.Id = doc.Id
	r.Error = doc.Error
	r.Expiry = doc.Expiry
	r.AccessURL = doc.AccessURL
}
