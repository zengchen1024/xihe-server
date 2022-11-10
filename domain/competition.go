package domain

type CompetitionSummary struct {
	Id       string
	Name     CompetitionName
	Desc     CompetitionDesc
	Host     CompetitionHost
	Bonus    CompetitionBonus
	Status   CompetitionStatus
	Duration CompetitionDuration
	Poster   URL
}

type Competition struct {
	CompetitionSummary

	Doc        URL
	DatasetDoc URL
	DatasetURL URL
}

type CompetitionOnPhase struct {
	phaseId string
	Phase   CompetitionPhase

	Competition
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

	Team     Team
	TeamRole TeamRole
}

type Team struct {
	Id   string
	Name TeamName
}

type CompetitionResult struct {
	SubmitAt int64
	OBSPath  string
	Status   string
	Score    float32
}
