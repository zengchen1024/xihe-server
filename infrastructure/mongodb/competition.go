package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

type competition struct {
	collectionName string
}

func (col competition) Get(cid, user string) (
	r repositories.CompetitionDO, isCompetitor bool, err error,
) {
	fieldCount := "competitors_count"
	fieldIsCompetitor := "is_competitor"

	var v []struct {
		DCompetition `bson:",inline"`

		IsCompetitor     bool `bson:"is_competitor"`
		CompetitorsCount int  `bson:"competitors_count"`
	}

	filter := bson.M{fieldId: cid}

	fields := bson.M{}
	if user != "" {
		fields[fieldIsCompetitor] = bson.M{
			"$in": bson.A{user, "$" + fieldCompetitors},
		}
	}

	err = col.get(filter, fieldCount, fields, &v)
	if err != nil || len(v) == 0 {
		return
	}

	item := &v[0]
	col.toCompetitionDO(&item.DCompetition, &r)
	r.CompetitorsCount = item.CompetitorsCount

	isCompetitor = item.IsCompetitor

	return
}

func (col competition) List(status, phase string) (
	r []repositories.CompetitionSummaryDO, err error,
) {
	fieldCount := "competitors_count"

	var v []struct {
		DCompetition `bson:",inline"`

		CompetitorsCount int `bson:"competitors_count"`
	}

	filter := bson.M{
		fieldPhase: phase,
	}
	if status != "" {
		filter[fieldStatus] = status
	}

	err = col.get(filter, fieldCount, nil, &v)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]repositories.CompetitionSummaryDO, len(v))

	for i := range v {
		col.toCompetitionSummaryDO(&v[i].DCompetition, &r[i])

		r[i].CompetitorsCount = v[i].CompetitorsCount
	}

	return
}

func (col competition) get(
	filter bson.M, fieldCount string,
	fields bson.M, result interface{},
) error {
	if fields == nil {
		fields = bson.M{}
	}

	fields[fieldCount] = bson.M{"$size": "$" + fieldCompetitors}

	f := func(ctx context.Context) error {
		pipeline := bson.A{
			bson.M{"$match": filter},
			bson.M{"$addFields": fields},
			bson.M{"$project": bson.M{
				fieldTeams:       0,
				fieldSubmissions: 0,
				fieldCompetitors: 0,
			}},
		}

		cursor, err := cli.collection(col.collectionName).Aggregate(ctx, pipeline)
		if err != nil {
			return err
		}

		return cursor.All(ctx, result)
	}

	return withContext(f)
}
