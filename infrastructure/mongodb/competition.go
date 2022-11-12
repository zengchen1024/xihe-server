package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewCompetitionMapper(name string) repositories.CompetitionMapper {
	return competition{name}
}

type competition struct {
	collectionName string
}

func (col competition) indexToDocFilter(index *repositories.CompetitionIndexDO) bson.M {
	return bson.M{
		fieldId:    index.Id,
		fieldPhase: index.Phase,
	}
}

func (col competition) Get(index *repositories.CompetitionIndexDO, competitor string) (
	r repositories.CompetitionDO, isCompetitor bool, err error,
) {
	fieldCount := "competitors_count"
	fieldCompetitor := "competitor"

	var v []struct {
		DCompetition `bson:",inline"`

		Competitor       []dCompetitor `bson:"competitor"`
		CompetitorsCount int           `bson:"competitors_count"`
	}

	fields := bson.M{}
	if competitor != "" {
		fields[fieldCompetitor] = bson.M{
			"$filter": bson.M{
				"input": "$" + fieldCompetitors,
				"cond":  eqCondForArrayElem(fieldAccount, competitor),
			},
		}
	}

	err = col.get(col.indexToDocFilter(index), fieldCount, fields, &v)
	if err != nil {
		return
	}

	if len(v) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	item := &v[0]
	col.toCompetitionDO(&item.DCompetition, &r)
	r.CompetitorsCount = item.CompetitorsCount

	isCompetitor = len(item.Competitor) > 0

	return
}

func (col competition) List(opt *repositories.CompetitionListOptionDO) (
	r []repositories.CompetitionSummaryDO, err error,
) {
	fieldCount := "competitors_count"
	fieldCompetitor := "competitor"

	var v []struct {
		DCompetition `bson:",inline"`

		Competitor       []dCompetitor `bson:"competitor"`
		CompetitorsCount int           `bson:"competitors_count"`
	}

	fields := bson.M{}
	if opt.Competitor != "" {
		fields[fieldCompetitor] = bson.M{
			"$filter": bson.M{
				"input": "$" + fieldCompetitors,
				"cond":  eqCondForArrayElem(fieldAccount, opt.Competitor),
			},
		}
	}

	filter := bson.M{
		fieldPhase: opt.Phase,
	}
	if opt.Status != "" {
		filter[fieldStatus] = opt.Status
	}

	err = col.get(filter, fieldCount, fields, &v)
	if err != nil || len(v) == 0 {
		return
	}

	r = make([]repositories.CompetitionSummaryDO, 0, len(v))

	for i := range v {
		item := &v[i]

		if opt.Competitor != "" && len(item.Competitor) == 0 {
			continue
		}

		col.toCompetitionSummaryDO(&item.DCompetition, &r[i])

		r[i].CompetitorsCount = item.CompetitorsCount
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

	key := "$" + fieldCompetitors
	fields[fieldCount] = bson.M{
		"$cond": bson.M{
			"if":   bson.M{"$isArray": key},
			"then": bson.M{"$size": key},
			"else": 0,
		},
	}

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

func (col competition) GetTeam(index *repositories.CompetitionIndexDO, user string) (
	[]repositories.CompetitorDO, error,
) {
	filter := col.indexToDocFilter(index)

	member, err := col.getCompetitor(filter, user)
	if err != nil {
		return nil, err
	}

	if member.TeamId == "" {
		do := repositories.CompetitorDO{}
		col.toCompetitorDO(&member, "", &do)

		return []repositories.CompetitorDO{do}, nil
	}

	team, members, err := col.getTeam(filter, member.TeamId)
	if err != nil {
		return nil, err
	}

	r := make([]repositories.CompetitorDO, len(members))

	for i := range members {
		col.toCompetitorDO(&members[i], team.Name, &r[i])
	}

	return r, nil
}

func (col competition) getCompetitor(docFilter bson.M, user string) (
	r dCompetitor, err error,
) {
	var v []DCompetition

	f := func(ctx context.Context) error {
		return cli.getArrayElem(
			ctx, col.collectionName, fieldCompetitors,
			docFilter, bson.M{fieldAccount: user}, nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 {
		err = errDocNotExists

		return
	}

	items := v[0].Competitors
	if len(items) == 0 {
		err = errDocNotExists
	} else {
		r = items[0]
	}

	return
}

func (col competition) getTeam(docFilter bson.M, tid string) (
	r dTeam, members []dCompetitor, err error,
) {
	var v []DCompetition

	f := func(ctx context.Context) error {
		return cli.getArraysElem(
			ctx, col.collectionName, docFilter,
			map[string]bson.M{
				fieldTeams:       {fieldId: tid},
				fieldCompetitors: {fieldTId: tid},
			}, nil, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 {
		err = errDocNotExists

		return
	}

	teams := v[0].Teams
	members = v[0].Competitors

	if len(members) == 0 || len(teams) == 0 {
		err = errDocNotExists
	} else {
		r = teams[0]
	}

	return
}

func (col competition) GetResult(index *repositories.CompetitionIndexDO) (
	smallerOk bool,
	teams []repositories.CompetitionTeamDO,
	results []repositories.CompetitionResultDO, err error,
) {
	var v []DCompetition

	f := func(ctx context.Context) error {
		return cli.getDoc(
			ctx, col.collectionName, col.indexToDocFilter(index),
			bson.M{
				fieldOrder:       1,
				fieldTeams:       1,
				fieldSubmissions: 1,
			}, &v,
		)
	}

	if err = withContext(f); err != nil {
		return
	}

	if len(v) == 0 {
		err = errDocNotExists

		return
	}

	rs := v[0].Submissions
	if len(rs) == 0 {
		return
	}

	results = make([]repositories.CompetitionResultDO, len(rs))
	for i := range rs {
		col.toCompetitionResultDO(&rs[i], &results[i])
	}

	smallerOk = v[0].SmallerOk

	if ts := v[0].Teams; len(ts) > 0 {
		teams = make([]repositories.CompetitionTeamDO, len(ts))
		for i := range ts {
			col.toCompetitionTeamDO(&ts[i], &teams[i])
		}
	}

	return
}
