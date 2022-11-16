package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/infrastructure/repositories"
	"github.com/opensourceways/xihe-server/utils"
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
	r repositories.CompetitionDO, info repositories.CompetitorInfoDO, err error,
) {
	var v []competitionInfo

	if err = col.get(col.indexToDocFilter(index), competitor, &v); err != nil {
		return
	}

	if len(v) == 0 {
		err = repositories.NewErrorDataNotExists(errDocNotExists)

		return
	}

	item := &v[0]
	col.toCompetitionDO(&item.DCompetition, &r)
	r.CompetitorsCount = item.CompetitorsCount

	if len(item.Competitor) > 0 {
		info.IsCompetitor = true

		if c := item.Competitor[0]; c.TeamId != "" {
			info.TeamId = c.TeamId
			info.TeamRole = c.TeamRole
		}
	}

	return
}

func (col competition) List(opt *repositories.CompetitionListOptionDO) (
	r []repositories.CompetitionSummaryDO, err error,
) {
	filter := bson.M{
		fieldPhase: opt.Phase,
	}
	if opt.Status != "" {
		filter[fieldStatus] = opt.Status
	}

	var v []competitionInfo

	err = col.get(filter, opt.Competitor, &v)
	if err != nil || len(v) == 0 {
		return
	}

	b := opt.Competitor != ""
	j := 0
	r = make([]repositories.CompetitionSummaryDO, len(v))

	for i := range v {
		item := &v[i]

		if b && len(item.Competitor) == 0 {
			continue
		}

		col.toCompetitionSummaryDO(&item.DCompetition, &r[j])

		r[j].CompetitorsCount = item.CompetitorsCount
		j++
	}

	r = r[:j]

	return
}

type competitionInfo struct {
	DCompetition `bson:",inline"`

	Competitor       []dCompetitor `bson:"competitor"`
	CompetitorsCount int           `bson:"competitors_count"`
}

func (col competition) get(
	filter bson.M, competitor string, result *[]competitionInfo,
) error {
	key := "$" + fieldCompetitors
	fieldCount := "competitors_count"
	fieldCompetitor := "competitor"

	fields := bson.M{}
	if competitor != "" {
		fields[fieldCompetitor] = bson.M{
			"$filter": bson.M{
				"input": key,
				"cond":  eqCondForArrayElem(fieldAccount, competitor),
			},
		}
	}

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
	results []repositories.CompetitionSubmissionDO, err error,
) {
	var v DCompetition

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

	rs := v.Submissions
	if len(rs) == 0 {
		return
	}

	results = make([]repositories.CompetitionSubmissionDO, len(rs))
	for i := range rs {
		col.toCompetitionSubmissionDO(&rs[i], &results[i])
	}

	smallerOk = v.SmallerOk

	if ts := v.Teams; len(ts) > 0 {
		teams = make([]repositories.CompetitionTeamDO, len(ts))
		for i := range ts {
			col.toCompetitionTeamDO(&ts[i], &teams[i])
		}
	}

	return
}

func (col competition) GetSubmisstions(index *repositories.CompetitionIndexDO, competitor string) (
	repo repositories.CompetitionRepoDO,
	results []repositories.CompetitionSubmissionDO, err error,
) {
	filter := col.indexToDocFilter(index)

	member, err := col.getCompetitor(filter, competitor)
	if err != nil {
		return
	}

	if member.TeamId == "" {
		return col.getResultOfCompetitor(filter, bson.M{
			fieldAccount: member.Account,
		})
	}

	return col.getResultOfCompetitor(filter, bson.M{
		fieldTId: member.TeamId,
	})
}

func (col competition) getResultOfCompetitor(docFilter, resultFilter bson.M) (
	repo repositories.CompetitionRepoDO,
	results []repositories.CompetitionSubmissionDO, err error,
) {
	var v []DCompetition

	f := func(ctx context.Context) error {
		return cli.getArraysElem(
			ctx, col.collectionName, docFilter,
			map[string]bson.M{
				fieldRepos:       resultFilter,
				fieldSubmissions: resultFilter,
			}, nil, &v)
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

	results = make([]repositories.CompetitionSubmissionDO, len(rs))
	for i := range rs {
		col.toCompetitionSubmissionDO(&rs[i], &results[i])
	}

	if ts := v[0].Repos; len(ts) > 0 {
		col.toCompetitionRepoDO(&ts[0], &repo)
	}

	return
}

func (col competition) SaveSubmission(
	index *repositories.CompetitionIndexDO,
	do *repositories.CompetitionSubmissionDO,
) (string, error) {
	dateTag := utils.ToDate(do.SubmitAt)
	v := new(dSubmission)
	do.Id = newId()
	col.toSubmissionDoc(do, v, dateTag)

	doc, err := genDoc(v)
	if err != nil {
		return "", err
	}

	docFilter := col.indexToDocFilter(index)

	sf := bson.M{fieldDate: dateTag}
	if do.TeamId != "" {
		sf[fieldTId] = do.TeamId
	} else {
		sf[fieldAccount] = do.Individual
	}
	appendElemMatchToFilter(fieldSubmissions, false, sf, docFilter)

	f := func(ctx context.Context) error {
		return cli.pushArrayElem(
			ctx, col.collectionName, fieldSubmissions, docFilter, doc,
		)
	}

	if err = withContext(f); err != nil {
		if isDocNotExists(err) {
			err = repositories.NewErrorDuplicateCreating(err)
		}
	}

	return do.Id, err
}
