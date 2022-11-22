package challengeimpl

type Config struct {
	CompetitionSuccessScore  int      `json:"Competition_success_score"   required:"true"`
	CompetitionSuccessStatus string   `json:"Competition_success_status"   required:"true"`
	Competitions             []string `json:"competitions"                 required:"true"`
	AIQuestion               string   `json:"ai_question"                  required:"true"`
}
