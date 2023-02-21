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

func (impl playerRepoImpl) docFilter(cid string, a types.Account) bson.M {
	filter := bson.M{
		fieldCid:     cid,
		fieldEnabled: true,
	}
	impl.cli.AppendElemMatchToFilter(
		"competitors", true,
		bson.M{"account": a.Account()}, filter,
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
	// TODO player to dCompetitor

	obj := dPlayer{
		CompetitionId: p.CompetitionId,
		Competitors:   []dCompetitor{c},
		Enabled:       true,
	}
	if p.IsATeam() {
		obj.Leader = p.Leader.Account.Account()
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
	if err := impl.updateEnabledOfPlayer(p.Id, false, version); err != nil {
		return err
	}

	return impl.insertPlayer(p)
}

func (impl playerRepoImpl) updateEnabledOfPlayer(pid string, enable bool, version int) error {
	return impl.update(pid, bson.M{fieldEnabled: enable}, version)
}

func (impl playerRepoImpl) SaveTeamName(p *domain.Player, version int) error {
	return impl.update(p.Id, bson.M{fieldTeamName: p.Team.Name.TeamName()}, version)
}

func (impl playerRepoImpl) update(pid string, doc bson.M, version int) error {
	filter, err := impl.cli.ObjectIdFilter(pid)
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

		return
	}

	version = v.Version

	// convert

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

	// convert

	return
}

func (impl playerRepoImpl) CompetitorsCount(cid string) (int, error) {
	return 0, nil
}

func (impl playerRepoImpl) AddMember(
	team repository.PlayerVersion,
	member repository.PlayerVersion,
) error {
	return nil
}
