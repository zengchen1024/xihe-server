package domain

type CompetitionSummary struct {
	Id         string
	Name       CompetitionName
	Desc       CompetitionDesc
	Host       CompetitionHost
	Bonus      CompetitionBonus
	Status     CompetitionStatus
	Duration   CompetitionDuration
	Poster     URL
	ScoreOrder CompetitionScoreOrder
}

type Competition struct {
	CompetitionSummary

	Doc        URL
	Forum      URL
	DatasetDoc URL
	DatasetURL URL

	Phase   CompetitionPhase
	Enabled bool
}

type Competitor struct {
	Account  Account
	Name     CompetitorName
	City     City
	Email    Email
	Phone    Phone
	Identity CompetitionIdentity
	Province Province
	Detail   map[string]string

	Team     CompetitionTeam
	TeamRole TeamRole
}

type CompetitionTeam struct {
	Id   string
	Name TeamName
}

type CompetitionSubmission struct {
	Id string

	TeamId     string
	Individual Account

	SubmitAt int64
	OBSPath  string
	Status   string
	Score    float32
}

type CompetitionRepo struct {
	TeamId     string
	Individual Account

	Owner Account
	Repo  ResourceName
}

func (r *CompetitionSubmission) IsTeamWork() bool {
	return r.TeamId != ""
}

func (r *CompetitionSubmission) Key() string {
	if r.TeamId != "" {
		return r.TeamId
	}

	return r.Individual.Account()
}

type CompetitionIndex struct {
	Id    string
	Phase CompetitionPhase
}

type CompetitionScoreOrder interface {
	IsBetterThanB(a, b float32) bool
}

func NewCompetitionScoreOrder(b bool) CompetitionScoreOrder {
	return smallerIsBetter(b)
}

type smallerIsBetter bool

func (order smallerIsBetter) IsBetterThanB(a, b float32) bool {
	if order {
		return a <= b
	}

	return a >= b
}
