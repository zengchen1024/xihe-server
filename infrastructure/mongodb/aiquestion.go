package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewAIQuestionMapper(name string) repositories.AIQuestionMapper {
	return aiquestion{name}
}

type aiquestion struct {
	collectionName string
}

func (col aiquestion) GetCompetitorAndSubmission(qid, competitor string) (
	IsCompetitor bool, scores []int, err error) {
	docFilter, err := objectIdFilter(qid)
	if err != nil {
		return
	}

	var v []dAIQuestion

	f := func(ctx context.Context) error {
		filter := bson.M{fieldAccount: competitor}

		return cli.getArraysElem(
			ctx, col.collectionName,
			docFilter,
			map[string]bson.M{
				fieldCompetitors: filter,
				fieldSubmissions: filter,
			},
			bson.M{
				fieldCompetitors + "." + fieldAccount: 1,
				fieldSubmissions + "." + fieldScore:   1,
			},
			&v,
		)
	}

	if err = withContext(f); err != nil || len(v) == 0 {
		return
	}

	item := &v[0]

	if len(item.Competitors) == 0 {
		return
	}

	IsCompetitor = true

	items := item.Submissions
	scores = make([]int, len(items))
	for i := range items {
		scores[i] = items[i].Score
	}

	return
}

func (col aiquestion) SaveCompetitor(qid string, do *repositories.CompetitorInfoDO) error {
	v := new(DCompetitorInfo)
	toCompetitorInfoDOC(v, do)

	doc, err := genDoc(v)
	if err != nil {
		return err
	}

	docFilter, err := objectIdFilter(qid)
	if err != nil {
		return err
	}

	appendElemMatchToFilter(
		fieldCompetitors, false, bson.M{fieldAccount: do.Account}, docFilter,
	)

	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, col.collectionName, fieldCompetitors, docFilter, doc,
		)
	}

	return withContext(f)
}
