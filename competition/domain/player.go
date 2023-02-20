package domain

import (
	"errors"

	//"github.com/opensourceways/xihe-server/competition/domain/repository"
	types "github.com/opensourceways/xihe-server/domain"
)

type Competitor struct {
	Account  types.Account
	Name     CompetitorName
	City     City
	Email    types.Email
	Phone    Phone
	Identity CompetitionIdentity
	Province Province
	Detail   map[string]string
}

type Team struct {
	Name    TeamName
	Members []Competitor
}

func (t *Team) isMember(a types.Account) bool {
	return t.indexOfMember(a) >= 0
}

func (t *Team) member(a types.Account) (Competitor, bool) {
	if i := t.indexOfMember(a); i >= 0 {
		return t.Members[i], true
	}

	return Competitor{}, false
}

func (t *Team) indexOfMember(a types.Account) int {
	for i := range t.Members {
		if item := &t.Members[i]; item.Account.Account() == a.Account() {
			return i
		}
	}

	return -1
}

func (t *Team) remove(a types.Account) error {
	i := t.indexOfMember(a)
	if i < 0 {
		return errors.New("not a member")
	}

	if i == 0 {
		t.Members = t.Members[1:]
	} else {
		if n := len(t.Members) - 1; i != n {
			t.Members[i] = t.Members[n]
		}
		t.Members = t.Members[:i]
	}

	return nil
}

func (t *Team) join(c *Competitor) error {
	if len(t.Members) >= 3 {
		return errors.New("members are enough")
	}

	if t.isMember(c.Account) {
		return errors.New("already joined")
	}

	t.Members = append(t.Members, *c)

	return nil
}

type PlayerIndex struct {
	Id            string
	CompetitionId string
}

type Player struct {
	PlayerIndex
	IsFinalist bool

	Leader Competitor
	Team   Team
	user   types.Account
}

func NewPlayer(cid string, c *Competitor) Player {
	return Player{
		PlayerIndex: PlayerIndex{
			CompetitionId: cid,
		},
		Leader: *c,
		user:   c.Account,
	}
}

func (p *Player) CompetitorsCount() int {
	return len(p.Team.Members) + 1
}

func (p *Player) IsIndividual() bool {
	return p.Team.Name == nil
}

func (p *Player) IsATeam() bool {
	return p.Team.Name != nil
}

func (p *Player) isUserTheLeader() bool {
	return p.user != nil && p.user.Account() == p.Leader.Account.Account()
}

func (p *Player) HasPermission() bool {
	return p.IsIndividual() || p.isUserTheLeader()
}

func (p *Player) Name() string {
	if p.IsIndividual() {
		return p.Leader.Account.Account()
	}

	return p.Team.Name.TeamName()
}

// TODO: this func should be called when generating a new playser.
func (p *Player) SetUser(a types.Account) {
	p.user = a
}

func (p *Player) Members() []*Competitor {
	return nil
}

func (p *Player) RoleOfCurrentCompetitor() string {
	if p.IsATeam() && p.isUserTheLeader() {
		return competitionTeamRoleLeader
	}

	return ""
}

func (p *Player) CurrentCompetitor() Competitor {
	if p.IsIndividual() {
		return p.Leader
	}

	m, _ := p.Team.member(p.user)

	return m
}

func (p *Player) CreateTeam(name TeamName) error {
	if p.IsATeam() {
		return errors.New("team is ready, no need to create team again")
	}

	p.Team.Name = name

	return nil
}

func (p *Player) ChangeTeamName(name TeamName) error {
	if !p.IsATeam() {
		return errors.New("I am not a team")
	}

	if !p.isUserTheLeader() {
		return errors.New("I am not leader")
	}

	p.Team.Name = name

	return nil
}

func (p *Player) JoinTo(team *Player) error {
	if !p.IsIndividual() {
		return errors.New("you are not an individual competitor")
	}

	if team.IsIndividual() {
		return errors.New("it is not a team")
	}

	return team.join(&p.Leader)
}

func (p *Player) join(c *Competitor) error {
	if p.Leader.Account.Account() == c.Account.Account() {
		return errors.New("invalid operation")
	}

	return p.Team.join(c)
}

func (p *Player) Quit() error {
	if !p.IsATeam() {
		return errors.New("invalid operation")
	}

	if p.isUserTheLeader() {
		return errors.New("leader can't quit directly")
	}

	return p.Team.remove(p.user)
}

func (p *Player) Delete(c types.Account) (Player, error) {
	if !p.IsATeam() {
		return Player{}, errors.New("invalid operation")
	}

	if p.isUserTheLeader() {
		return Player{}, errors.New("only leader can delete a member")
	}

	m, ok := p.Team.member(c)
	if !ok {
		return Player{}, errors.New("not a meber")
	}

	_ = p.Team.remove(c)

	return Player{Leader: m}, nil
}

func (p *Player) TransferLeader(newOne types.Account) error {
	if !p.IsATeam() {
		return errors.New("invalid operation")
	}

	if p.isUserTheLeader() {
		return errors.New("only leader can delete a member")
	}

	t := &p.Team
	i := t.indexOfMember(newOne)
	if i < 0 {
		return errors.New("not a meber")
	}

	p.Leader, t.Members[i] = t.Members[i], p.Leader

	return nil
}

/*
func (p *Player) Disband(repo repository.Competition) error {
	// get team
	// set member to individual and remove the member one by one

	if !p.IsATeam() {
		return errors.New("invalid operation")
	}

	if p.isUserTheLeader() {
		return errors.New("only leader can delete a member")
	}

	tm := p.Team.Members
	for i := range tm {
		// save member

	}



	return nil

}
*/
