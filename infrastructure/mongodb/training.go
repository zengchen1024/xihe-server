package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func trainingDocFilter(owner, projectId string) bson.M {
	return bson.M{
		fieldOwner: owner,
		fieldPId:   projectId,
	}
}

type training struct {
	collectionName string
}

func (col training) newDoc(do *repositories.UserTrainingDO) error {
	docFilter := trainingDocFilter(do.Owner, do.ProjectId)

	doc := bson.M{
		fieldOwner:   do.Owner,
		fieldPId:     do.ProjectId,
		fieldName:    do.ProjectName,
		fieldRId:     do.ProjectRepoId,
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

func (col training) Insert(do *repositories.UserTrainingDO, version int) (
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

func (col training) insert(do *repositories.UserTrainingDO, version int) (identity string, err error) {
	identity = newId()
	do.Id = identity

	doc, err := col.toTrainingDoc(do)
	if err != nil {
		return
	}

	docFilter := trainingDocFilter(do.Owner, do.ProjectId)

	f := func(ctx context.Context) error {
		return cli.updateDoc(
			ctx, col.collectionName, docFilter,
			bson.M{fieldItems: doc}, mongoCmdPush, version,
		)
	}

	err = withContext(f)

	return
}

func (col training) List(user, projectId string) ([]repositories.TrainingSummaryDO, int, error) {
	var v dTraining

	field := func(k string) string {
		return fieldItems + k
	}

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName,
			trainingDocFilter(user, projectId),
			bson.M{
				fieldVersion:          1,
				field(fieldJob):       1,
				field(fieldName):      1,
				field(fieldDesc):      1,
				field(fieldDetail):    1,
				field(fieldCreatedAt): 1,
			}, &v)
	}

	if err := withContext(f); err != nil {
		if isDocNotExists(err) {
			return nil, 0, nil
		}

		return nil, 0, err
	}

	t := v.Items
	r := make([]repositories.TrainingSummaryDO, len(t))

	for i := range t {
		col.toTrainingSummary(&t[i], &r[i])
	}

	return r, v.Version, nil
}

func (col training) toTrainingSummary(t *trainingItem, s *repositories.TrainingSummaryDO) {
	*s = repositories.TrainingSummaryDO{
		Name:      t.Name,
		Desc:      t.Desc,
		JobId:     t.Job.JobId,
		Status:    t.JobDetail.Status,
		Duration:  t.JobDetail.Duration,
		CreatedAt: t.CreatedAt,
		Endpoint:  t.Job.Endpoint,
	}
}
