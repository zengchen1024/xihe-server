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
	DatasetDoc URL
	DatasetURL URL
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

type CompetitionResult struct {
	Id string

	TeamId     string
	Individual CompetitorName

	SubmitAt int64
	OBSPath  string
	Status   string
	Score    float32
}

func (r *CompetitionResult) IsTeamWork() bool {
	return r.TeamId != ""
}

func (r *CompetitionResult) Key() string {
	if r.TeamId != "" {
		return r.TeamId
	}

	return r.Individual.CompetitorName()
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
