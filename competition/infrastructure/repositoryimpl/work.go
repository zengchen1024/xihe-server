package repositoryimpl

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	repoerr "github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func withContext(f func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second, // TODO use config
	)
	defer cancel()

	return f(ctx)
}

func NewWorkRepo(collectionName string, m Mongodb) repository.Work {
	return workRepoImpl{
		cli:            mongodbClient{m},
		collectionName: collectionName,
	}
}

type workRepoImpl struct {
	cli            mongodbClient
	collectionName string
}

func (impl workRepoImpl) docFilter(index *domain.WorkIndex) bson.M {
	return bson.M{
		fieldCid: index.CompetitionId,
		fieldPid: index.PlayerId,
	}
}

func (impl workRepoImpl) SaveWork(w *domain.Work) error {
	doc, err := genDoc(dWork{
		CompetitionId: w.CompetitionId,
		PlayerId:      w.PlayerId,
		PlayerName:    w.PlayerName,
		Final:         []dSubmission{},
		Preliminary:   []dSubmission{},
	})
	if err != nil {
		return err
	}
	doc[fieldVersion] = 0

	f := func(ctx context.Context) error {
		_, err := impl.cli.newDocIfNotExist(
			ctx, impl.collectionName,
			impl.docFilter(&w.WorkIndex), doc,
		)

		return err
	}

	if err = withContext(f); err != nil {
		if isDocExists(err) {
			err = repoerr.NewErrorDuplicateCreating(err)
		}
	}

	return err
}

func (impl workRepoImpl) SaveRepo(w *domain.Work, version int) error {
	f := func(ctx context.Context) error {
		return impl.cli.updateDoc(
			ctx, impl.collectionName,
			impl.docFilter(&w.WorkIndex),
			bson.M{fieldRepo: w.Repo}, mongoCmdSet, version,
		)
	}

	err := withContext(f)
	if err != nil {
		if isDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}
	}

	return err

}

func (impl workRepoImpl) AddSubmission(
	w *domain.Work, cs *domain.PhaseSubmission, version int,
) error {
	doc, err := genDoc(dSubmission{
		Id:       cs.Id,
		Status:   cs.Status,
		OBSPath:  cs.OBSPath,
		SubmitAt: cs.SubmitAt,
	})
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		filter := impl.docFilter(&w.WorkIndex)
		filter[fieldVersion] = version

		field := fieldPreliminary
		if cs.Phase.IsFinal() {
			field = fieldFinal
		}

		return impl.cli.pushArrayElem(
			ctx, impl.collectionName, field, filter, doc,
		)
	}

	if err = withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}
	}

	return err
}

func (impl workRepoImpl) SaveSubmission(
	w *domain.Work, submission *domain.PhaseSubmission, version int,
) error {
	doc, err := genDoc(dSubmission{
		Status: submission.Status,
		Score:  submission.Score,
	})
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		field := fieldPreliminary
		if submission.Phase.IsFinal() {
			field = fieldFinal
		}

		_, err := impl.cli.updateArrayElem(
			ctx, impl.collectionName, field,
			impl.docFilter(&w.WorkIndex), bson.M{fieldId: submission.Id},
			doc, version, 0,
		)

		return err
	}

	err = withContext(f)

	return err
}

func (impl workRepoImpl) FindWork(index domain.WorkIndex, Phase domain.CompetitionPhase) (
	w domain.Work, version int, err error,
) {
	var v dWork

	f := func(ctx context.Context) error {
		filter := impl.docFilter(&index)

		project := bson.M{}
		if Phase.IsPreliminary() {
			project[fieldFinal] = 0
		} else {
			project[fieldPreliminary] = 0
		}

		return impl.cli.getDoc(
			ctx, impl.collectionName, filter, project, &v,
		)
	}

	if err = withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repoerr.NewErrorDataNotExists(err)
		}

		return
	}

	version = v.Version

	// convert

	return
}

func (impl workRepoImpl) FindWorks(cid string) (ws []domain.Work, err error) {
	var v []dWork

	f := func(ctx context.Context) error {
		return impl.cli.getDocs(
			ctx, impl.collectionName,
			bson.M{fieldCid: cid}, nil, &v,
		)
	}

	if err = withContext(f); err != nil || len(v) == 0 {
		return
	}

	// convert

	return
}
