package repository

import (
	"github.com/opensourceways/xihe-server/competition/domain"
	types "github.com/opensourceways/xihe-server/domain"
)

type CompetitionListOption struct {
	Status domain.CompetitionStatus
}

type Competition interface {
	FindCompetition(cid string) (domain.Competition, error)
	FindCompetitions(*CompetitionListOption) ([]domain.CompetitionSummary, error)

	FindScoreOrder(cid string) (domain.CompetitionScoreOrder, error)
}

type Player interface {
	//DeleteTeam, remove all the members and delete the team when removing leader.
	// As a Team
	// change team name
	SaveTeamName(*domain.Player, int) error
	//SaveLeader, change the leader
	//AddMember, add a member and unable the original
	//RemoveMember, remove the member and enable the original

	// As an Individual

	// common function
	// SavePlayer should check if the player is individual or team.
	SavePlayer(*domain.Player, int) error

	CompetitorsCount(cid string) (int, error)
	FindPlayer(cid string, a types.Account) (domain.Player, int, error)

	FindCompetitionsUserApplied(types.Account) ([]string, error)
}

type Work interface {
	SaveWork(*domain.Work) error
	SaveRepo(*domain.Work, int) error
	AddSubmission(*domain.Work, *domain.CompetitionSubmission, int) error
	SaveSubmission(*domain.Work, *domain.CompetitionSubmission, int) error

	FindWork(*domain.WorkIndex, domain.CompetitionPhase) (domain.Work, int, error)
	FindWorks(cid string) ([]domain.Work, error)
}
