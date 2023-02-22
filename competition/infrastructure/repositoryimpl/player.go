package repositoryimpl

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/opensourceways/xihe-server/competition/domain"
	"github.com/opensourceways/xihe-server/competition/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
	repoerr "github.com/opensourceways/xihe-server/infrastructure/repositories"
)

func NewPlayerRepo(m mongodbClient) repository.Player {
	return playerRepoImpl{m}
}

type playerRepoImpl struct {
	cli mongodbClient
}

func (impl playerRepoImpl) playerFilter(p *domain.Player) (bson.M, error) {
	filter, err := impl.cli.ObjectIdFilter(p.Id)
	if err != nil {
		return nil, err
	}

	filter[fieldCid] = p.CompetitionId
	filter[fieldEnabled] = true

	return filter, nil
}

func (impl playerRepoImpl) docFilter(cid string, a types.Account) bson.M {
	filter := bson.M{
		fieldCid:     cid,
		fieldEnabled: true,
	}
	impl.cli.AppendElemMatchToFilter(
		fieldCompetitors, true,
		bson.M{fieldAccount: a.Account()}, filter,
	)

	return filter
}

func (impl playerRepoImpl) SavePlayer(p *domain.Player, version int) error {
	if p.IsATeam() {
		return impl.insertTeam(p, version)
	}

	return impl.insertPlayer(p)
}

func (repo playerRepoImpl) genPlayerDoc(p *domain.Player) (bson.M, error) {
	var c dCompetitor
	toCompetitorDoc(&p.Leader, &c)

	obj := dPlayer{
		CompetitionId: p.CompetitionId,
		Competitors:   []dCompetitor{c},
		Leader:        p.Leader.Account.Account(),
		Enabled:       true,
	}
	if p.IsATeam() {
		obj.TeamName = p.Team.Name.TeamName()
	}

	doc, err := genDoc(&obj)
	if err == nil {
		doc[fieldVersion] = 0
	}

	return doc, err
}

func (impl playerRepoImpl) insertPlayer(p *domain.Player) error {
	doc, err := impl.genPlayerDoc(p)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		filter := impl.docFilter(p.CompetitionId, p.Leader.Account)

		_, err := impl.cli.NewDocIfNotExist(ctx, filter, doc)

		return err
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocExists(err) {
			err = repoerr.NewErrorDuplicateCreating(err)
		}
	}

	return err
}

func (impl playerRepoImpl) insertTeam(p *domain.Player, version int) error {
	if err := impl.updateEnabledOfPlayer(p, false, version); err != nil {
		return err
	}

	return impl.insertPlayer(p)
}

func (impl playerRepoImpl) updateEnabledOfPlayer(p *domain.Player, enable bool, version int) error {
	return impl.update(p, bson.M{fieldEnabled: enable}, version)
}

func (impl playerRepoImpl) SaveTeamName(p *domain.Player, version int) error {
	return impl.update(p, bson.M{fieldTeamName: p.Team.Name.TeamName()}, version)
}

func (impl playerRepoImpl) update(p *domain.Player, doc bson.M, version int) error {
	filter, err := impl.playerFilter(p)
	if err != nil {
		return err
	}

	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(ctx, filter, doc, mongoCmdSet, version)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}
	}

	return err
}

func (impl playerRepoImpl) FindPlayer(cid string, a types.Account) (
	p domain.Player, version int, err error,
) {
	var v dPlayer

	f := func(ctx context.Context) error {
		return impl.cli.GetDoc(ctx, impl.docFilter(cid, a), nil, &v)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorDataNotExists(err)
		}
	} else {
		if err = v.toPlayer(&p); err == nil {
			p.SetCurrentUser(a)

			version = v.Version
		}
	}

	return
}

func (impl playerRepoImpl) FindCompetitionsUserApplied(a types.Account) (
	r []string, err error,
) {
	var v []dPlayer

	f := func(ctx context.Context) error {
		filter := impl.docFilter("", a)
		delete(filter, fieldCid)

		return impl.cli.GetDocs(ctx, filter, bson.M{fieldCid: 1}, &v)
	}

	if err = withContext(f); err != nil || len(v) == 0 {
		return
	}

	r = make([]string, len(v))
	for i := range v {
		r[i] = v[i].Id.Hex()
	}

	return
}

func (impl playerRepoImpl) CompetitorsCount(cid string) (int, error) {
	return 0, nil
}

// AddMember
func (impl playerRepoImpl) AddMember(
	team repository.PlayerVersion,
	member repository.PlayerVersion,
) error {
	err := impl.updateEnabledOfPlayer(member.Player, false, member.Version)
	if err != nil {
		return err
	}

	return impl.addMember(team, member.Player)
}

func (impl playerRepoImpl) addMember(
	team repository.PlayerVersion, member *domain.Player,
) error {
	filter, err := impl.playerFilter(team.Player)
	if err != nil {
		return err
	}

	var c dCompetitor
	toCompetitorDoc(&member.Leader, &c)
	doc, err := genDoc(&c)
	if err == nil {
		return err
	}

	f := func(ctx context.Context) error {
		return impl.cli.UpdateDoc(
			ctx, filter,
			bson.M{fieldCompetitors: doc}, mongoCmdPush, team.Version,
		)
	}

	if err = withContext(f); err != nil {
		if impl.cli.IsDocNotExists(err) {
			err = repoerr.NewErrorConcurrentUpdating(err)
		}
	}

	return err
}
