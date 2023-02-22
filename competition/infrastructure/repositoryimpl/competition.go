package repositoryimpl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewCompetitionRepo(m mongodbClient) repository.Competition {
	return competitionRepoImpl{m}
}

type competitionRepoImpl struct {
	cli mongodbClient
}

func (impl competitionRepoImpl) docFilter(cid string) bson.M {
	return bson.M{
		fieldId: cid,
	}
}

func (impl competitionRepoImpl) FindCompetition(cid string) (
	c domain.Competition, err error,
) {
	var v dCompetition

	f := func(ctx context.Context) error {
		filter := impl.docFilter(cid)

		return impl.cli.GetDoc(ctx, filter, nil, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return
	}

	err = v.toCompetition(&c)

	return
}

func (impl competitionRepoImpl) FindScoreOrder(cid string) (
	domain.CompetitionScoreOrder, error,
) {
	var v dCompetition

	f := func(ctx context.Context) error {
		filter := impl.docFilter(cid)

		return impl.cli.GetDoc(ctx, filter, bson.M{"order": 1}, &v)
	}

	if err := withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repositories.NewErrorDataNotExists(err)
		}

		return nil, err
	}

	return domain.NewCompetitionScoreOrder(v.SmallerOk), nil
}

func (impl competitionRepoImpl) FindCompetitions(opt *repository.CompetitionListOption) (
	[]domain.CompetitionSummary, error,
) {
	var v []dCompetition

	f := func(ctx context.Context) error {
		filter := bson.M{}
		if opt.Status != nil {
			filter["status"] = opt.Status.CompetitionStatus()
		}
		if len(opt.CompetitionIds) > 0 {
			filter[fieldId] = bson.M{
				"$in": opt.CompetitionIds,
			}
		}

		return impl.cli.GetDocs(ctx, filter, nil, &v)
	}

	if err := withContext(f); err != nil || len(v) == 0 {
		return nil, err
	}

	r := make([]domain.CompetitionSummary, len(v))
	for i := range v {
		if err := v[i].toCompetitionSummary(&r[i]); err != nil {
			return nil, err
		}
	}

	return r, nil
}
